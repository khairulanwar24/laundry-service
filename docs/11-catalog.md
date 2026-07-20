# 11 — Catalog (Layanan, Varian, Parfum, Diskon)

## Deskripsi

Domain katalog mencakup 4 entitas yang saling terkait, masing-masing punya file & controller/service/
repository sendiri di dalam satu folder `catalog/` (mengikuti pola pengelompokan yang sama dengan
domain `osce` di project sebelumnya, tapi package tetap mengikuti nama layer induk):

- **Service** (`laundry.services`) — layanan laundry (misal "Cuci Kering", "Setrika")
- **ServiceVariant** (`laundry.service_variants`) — varian harga/satuan/durasi per layanan
- **Perfume** (`laundry.perfumes`) — pilihan parfum outlet
- **Discount** (`laundry.discounts`) — diskon nominal/persen yang bisa dipakai saat membuat pesanan

Semua endpoint terbuka untuk **anggota outlet aktif** (bukan owner-only) — ini bukan bug, memang
perilaku asli Laravel `ServiceCatalogController`/`PerfumeController`/`DiscountController`.

---

## Endpoint

Semua di bawah `/outlets/:id_outlet` dengan `LaundryJWTMiddleware` + `RequireOutletMember`.

| Method | Path | Handler |
|--------|------|---------|
| `GET` | `/services?q=` | `ServiceController.ListActive` |
| `POST` | `/services` | `ServiceController.Create` |
| `PUT` | `/services/:id_service` | `ServiceController.Update` |
| `DELETE` | `/services/:id_service` | `ServiceController.Deactivate` |
| `POST` | `/services/:id_service/variants` | `ServiceVariantController.Create` |
| `PUT` | `/service-variants/:id_service_variant` | `ServiceVariantController.Update` |
| `GET` | `/perfumes` | `PerfumeController.ListActive` |
| `POST` | `/perfumes` | `PerfumeController.Create` |
| `PUT` | `/perfumes/:id_perfume` | `PerfumeController.Update` |
| `GET` | `/discounts` | `DiscountController.ListActive` |
| `POST` | `/discounts` | `DiscountController.Create` |
| `PUT` | `/discounts/:id_discount` | `DiscountController.Update` |

Catatan: tidak ada endpoint `DELETE` untuk variant/perfume/discount — sama seperti versi Laravel,
nonaktifkan lewat `PUT ... {"is_active": false}`.

---

## DTO (`domain/dto/catalog.go`)

```go
type ServiceForm struct {
    Nama          string   // required, max 120
    Prioritas     *int     // opsional 0-100, default 50
    LangkahProses []string // opsional, oneof per elemen: cuci|kering|setrika, default ["cuci","kering","setrika"]
    IsActive      *bool
}

type ServiceVariantForm struct {
    Nama                string  // required
    Satuan              string  // required, oneof: kg|pcs|meter
    HargaPerSatuan      float64 // required, >= 0
    DurasiPengerjaanJam int     // required, >= 1 — dipakai menghitung ETA pesanan
    GambarPath          string
    Catatan             string
    IsActive            *bool
}

type PerfumeForm struct {
    Nama, Catatan string // Nama required
    IsActive      *bool
}

type DiscountForm struct {
    Nama    string  // required
    Jenis   string  // required, oneof: nominal|percent
    Nilai   float64 // required, >= 0
    Catatan string
    IsActive *bool
}
```

---

## Alur Fungsi

### Service: List dengan pencarian

```
ListActive(idOutlet, q)
└── SQL JOIN services + service_variants (LEFT JOIN, agregasi json_agg per service)
    → filter: q kosong ATAU nama service/varian ILIKE %q%
    → urut prioritas DESC, nama ASC
```

### Service: Deactivate (hapus)

```
Deactivate(idOutlet, idService)
├── DeactivateVariantsByService() — semua varian milik service ini di-nonaktifkan dulu
└── Deactivate() — baru service-nya di-nonaktifkan
```

### ServiceVariant: Create/Update

```
Create(idOutlet, idService, form)
└── Verifikasi service ditemukan & milik outlet (FindByID) sebelum insert varian

Update(idOutlet, idServiceVariant, form)
└── FindByID(idServiceVariant) mengembalikan id_outlet lewat JOIN ke services —
    dicocokkan manual dengan idOutlet di service layer (bukan di query)
```

---

## Repository

File terpisah per entitas dalam satu folder `repositories/catalog/`:
`service.go`, `servicevariant.go`, `perfume.go`, `discount.go`.

Tabel: `laundry.services`, `laundry.service_variants`, `laundry.perfumes`, `laundry.discounts`.

---

## Dependency

Tidak ada dependency lintas domain — murni CRUD per outlet. Dipakai (dibaca) oleh domain `order`
saat membuat pesanan (memvalidasi & mengambil harga/diskon).
