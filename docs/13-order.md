# 13 — Order (Pesanan)

## Deskripsi

Domain paling inti dari bisnis laundry: pembuatan pesanan, perhitungan harga & diskon, mesin
status, pencatatan pembayaran, dan pengambilan pesanan oleh pelanggan.

### Mesin status

```
ANTRIAN ──► PROSES ──► SIAP_DIAMBIL ──► SELESAI
   │           │              │
   └────────►BATAL◄───────────┘
```

Transisi valid (ditegakkan di `isValidTransition`, bukan hanya validasi input):

| Dari | Boleh ke |
|---|---|
| `ANTRIAN` | `PROSES`, `BATAL` |
| `PROSES` | `SIAP_DIAMBIL`, `BATAL` |
| `SIAP_DIAMBIL` | `SELESAI`, `BATAL` |
| `SELESAI` | (terminal) |
| `BATAL` | (terminal) |

`status_pembayaran`: `UNPAID` → `PAID` (satu arah, dibayar penuh sekaligus — tidak ada cicilan/
partial payment, sengaja mengikuti desain Laravel aslinya).

### 6 perbaikan dibanding versi Laravel

Ditemukan & diperbaiki saat porting (disepakati sebelumnya bersama user):

1. **Diskon benar-benar dihitung.** `id_discount` yang dikirim client dipakai menghitung
   `nilai_diskon` (nominal langsung, atau `subtotal * nilai/100` untuk percent, di-clamp maksimal
   sebesar subtotal) dan disimpan di kolom `orders.id_discount`. Di Laravel, `discount_id`
   divalidasi tapi lalu **dibuang begitu saja** — diskon tidak pernah mempengaruhi total.
2. **Nomor invoice aman dari race condition.** Format `JL-yymmdd0001` dibuat di dalam transaksi
   dengan `pg_advisory_xact_lock(hashtext('laundry_invoice:{outlet}:{tanggal}'))`, sehingga dua
   request buat-pesanan bersamaan untuk outlet & tanggal yang sama tidak bisa menghasilkan nomor
   invoice yang sama. Laravel men-generate nomor lewat baca-lalu-tulis tanpa lock apa pun.
3. **Setiap perubahan status mencatat pelaku & catatan.** `order_status_histories.id_user` dan
   `catatan` selalu diisi dari user yang login + body request — termasuk lewat endpoint
   `POST .../status`. Di Laravel, controller tidak meneruskan `user()->id` maupun `notes` ke
   service, jadi kolom itu selalu `NULL`.
4. **Membatalkan pesanan membatalkan pembayarannya.** Kalau status diubah ke `BATAL` dan pesanan
   itu sudah punya pembayaran `SUCCESS`, pembayaran otomatis di-`VOID`-kan dalam alur yang sama.
   Di Laravel, method `cancelOrder()` yang punya logika ini tidak pernah dipanggil oleh route
   manapun — endpoint `/status` yang sesungguhnya melewatinya begitu saja.
5. **Item pesanan menampilkan nama layanan & satuan yang benar** — lewat JOIN langsung ke
   `service_variants`/`services`, bukan relasi Eloquent yang salah (di Laravel field ini selalu
   `null` di response).
6. **Item dengan satuan `pcs` wajib qty bulat** — divalidasi eksplisit (bukan cuma longgar seperti
   sebelumnya).

Yang **tidak** diubah (memang desain Laravel, bukan bug): pembayaran hanya bisa lunas sekali jalan
(tidak ada cicilan, `amount` boleh lebih besar dari total tapi tidak boleh kurang), dan pickup tidak
mensyaratkan status pembayaran sudah `PAID`.

---

## Endpoint

Semua di bawah `/outlets/:id_outlet/orders` dengan `LaundryJWTMiddleware` + `RequireOutletMember`.

| Method | Path | Query/Body | Handler |
|--------|------|-----------|---------|
| `GET` | `/` | `tab`, `q`, `page`, `per_page` | `List` |
| `GET` | `/search` | `q` (wajib), `page`, `per_page` | `Search` |
| `POST` | `/` | `CreateOrderForm` | `Create` |
| `GET` | `/:id_order` | - | `Get` |
| `GET` | `/:id_order/history` | - | `GetHistory` |
| `POST` | `/:id_order/status` | `ChangeStatusForm` | `ChangeStatus` |
| `POST` | `/:id_order/pay` | `PayOrderForm` | `Pay` |
| `POST` | `/:id_order/pickup` | `PickupForm` | `Pickup` |

`tab` (case-insensitive, "-" dan "_" dianggap sama): `antrian`, `proses`, `siap_diambil`/`siap-ambil`,
`selesai`, `batal`, atau kosong/`all` untuk semua status. `search` (endpoint `/search`) selalu
dibatasi ke status `ANTRIAN`, `PROSES`, `SIAP_DIAMBIL` saja dan `q` wajib diisi.

---

## DTO (`domain/dto/order.go`)

```go
type CreateOrderItemForm struct {
    IDServiceVariant string  // required, uuid4
    Qty              float64 // required, > 0
    Catatan          string
}

type CreateOrderForm struct {
    IDCustomer string
    Items      []CreateOrderItemForm // required, minimal 1
    IDPerfume  string                // opsional
    IDDiscount string                // opsional
    Catatan    string
}

type ChangeStatusForm struct {
    To      string // required, oneof: ANTRIAN|PROSES|SIAP_DIAMBIL|SELESAI|BATAL
    Catatan string
}

type PayOrderForm struct {
    IDPaymentMethod string
    Jumlah          float64 // required, > 0
    NoReferensi     string
}

type PickupForm struct {
    Catatan string
}

// OrderItemInput: bentuk item setelah dihitung/di-snapshot service layer,
// dikirim ke repository (bukan dari request langsung).
type OrderItemInput struct {
    IDServiceVariant       string
    Satuan                 string
    Qty                    float64
    HargaPerSatuanSnapshot float64
    TotalHarga             float64
    Catatan                string
}
```

---

## Alur Fungsi

### 1. Create Order

```
Service.Create(idOutlet, idUser, form)
├── Validasi customer ada & milik outlet → repository.GetCustomer().FindByID()
├── Jika ada perfume_id → GetPerfume().FindByID()
├── Jika ada discount_id → GetDiscount().FindByID() → ambil jenis & nilai
├── Untuk tiap item:
│   ├── GetServiceVariant().FindByID() → validasi milik outlet & aktif
│   ├── Kalau satuan == "pcs" → qty harus bilangan bulat (math.Trunc check)
│   ├── total_harga = qty * harga_per_satuan (harga di-snapshot saat ini)
│   └── Lacak durasi_pengerjaan_jam terbesar antar item → dipakai hitung ETA
├── subtotal = jumlah semua total_harga
├── nilai_diskon = nominal langsung, atau subtotal*persen/100 untuk jenis "percent"
│                  (di-clamp maksimal sebesar subtotal)
├── total = subtotal - nilai_diskon
├── eta_at = checkin_at + max(durasi_pengerjaan_jam) jam
└── repository.GetOrder().CreateOrder(...) — SATU transaksi:
    ├── pg_advisory_xact_lock per outlet+tanggal
    ├── Hitung nomor urut invoice hari itu (MAX+1) dari kolom invoice_no
    ├── INSERT laundry.orders (status=ANTRIAN, status_pembayaran=UNPAID)
    ├── INSERT laundry.order_items (satu per item)
    └── INSERT laundry.order_status_histories (from=NULL, to=ANTRIAN, id_user=pembuat)
```

### 2. Change Status

```
Service.ChangeStatus(idOutlet, idOrder, idUser, form)
├── FindByID() → ambil status saat ini
├── isValidTransition(current, form.To) → gagal → error "Transisi status tidak valid dari X ke Y"
├── UpdateStatus() → UPDATE orders SET status=?, selesai_at/batal_at diisi otomatis via CASE SQL
├── InsertStatusHistory(from, to, idUser, catatan)
└── Jika form.To == "BATAL":
    └── FindPaymentByOrder() ada & status SUCCESS → VoidPayment() (set status VOID)
```

### 3. Pay

```
Service.Pay(idOutlet, idOrder, idUser, form)
├── FindByID() pesanan
├── FindPaymentByOrder() sudah ada → tolak ("Pesanan sudah memiliki pembayaran")
├── Validasi payment_method milik outlet → GetPaymentMethod().FindByID()
├── form.Jumlah < order.total → tolak ("Jumlah pembayaran kurang dari total pesanan")
├── CreatePayment() → INSERT laundry.payments (status SUCCESS)
├── SetPaymentStatus(idOrder, "PAID")
└── Jika status pesanan masih "ANTRIAN" → otomatis UpdateStatus("PROSES")
    + InsertStatusHistory(..., idUser, "Pembayaran berhasil diproses")
```

### 4. Pickup

```
Service.Pickup(idOutlet, idOrder, idUser, form)
├── Status harus "SIAP_DIAMBIL", selain itu tolak
├── UpdateStatus("SELESAI") + InsertStatusHistory(...)
└── SetPickup(idOrder, idUser) → UPDATE diambil_at=NOW(), diambil_oleh_id_user=idUser
```

---

## Repository (`repositories/order/order.go`)

Tabel: `laundry.orders`, `laundry.order_items`, `laundry.order_status_histories`,
`laundry.payments`.

| Fungsi | Keterangan |
|---|---|
| `CreateOrder` | Satu transaksi: advisory lock → generate invoice → insert order+items+history |
| `List` / `Count` | Filter status IN (...) + pencarian invoice_no/nama/telepon pelanggan |
| `ListItems` | JOIN ke `service_variants` untuk nama layanan & satuan |
| `ListHistory` | JOIN ke `users` untuk nama pelaku perubahan status |
| `UpdateStatus` | `selesai_at`/`batal_at` diisi otomatis via `CASE WHEN` di SQL yang sama |
| `VoidPayment` | `UPDATE payments SET status='VOID' WHERE id_order=? AND status='SUCCESS'` |

---

## Dependency

Order adalah domain yang paling banyak **membaca lintas domain** lewat `IRepositoryRegistry`:
`GetCustomer()`, `GetPerfume()`, `GetDiscount()`, `GetServiceVariant()`, `GetPaymentMethod()`.
Tidak ada HTTP call antar domain — semua lewat pemanggilan langsung antar repository di dalam satu
process, sesuai pola registry yang sudah ada.
