# 12 — Customer (Pelanggan)

## Deskripsi

Domain pelanggan outlet: CRUD data pelanggan serta riwayat pesanan per pelanggan.

---

## Endpoint

Semua di bawah `/outlets/:id_outlet/customers` dengan `LaundryJWTMiddleware` + `RequireOutletMember`.

| Method | Path | Query params | Handler |
|--------|------|---------------|---------|
| `GET` | `/` | `search`, `is_active`, `page`, `per_page` | `List` |
| `POST` | `/` | - | `Create` |
| `GET` | `/:id_customer` | - | `Get` |
| `PUT` | `/:id_customer` | - | `Update` |
| `DELETE` | `/:id_customer` | - | `Delete` |
| `GET` | `/:id_customer/orders` | `status`, `from_date`, `to_date`, `page`, `per_page` | `ListOrders` |

---

## DTO (`domain/dto/customer.go`)

```go
type CustomerForm struct {
    Nama     string // required, max 255
    Telepon  string // opsional, max 20
    Email    string // opsional, format email
    Alamat   string // opsional, max 500
    IsActive *bool  // opsional, default true
}
```

---

## Alur Fungsi

### Create / Update — validasi unik per outlet

```
├── ExistsPhoneInOutlet(idOutlet, telepon, excludeID) — dicek di level SERVICE, bukan constraint DB
│   (persis seperti Laravel: unique:customers,phone,...,outlet_id)
├── ExistsEmailInOutlet(idOutlet, email, excludeID) — sama untuk email
└── Insert/Update ke laundry.customers
```

`excludeID` diisi `id_customer` saat update (supaya pelanggan itu sendiri tidak dianggap konflik).

### Delete — diblokir jika sudah punya riwayat pesanan

```
Delete(idOutlet, idCustomer)
├── CountOrders(idCustomer) — hitung baris di laundry.orders
├── Jika > 0 → tolak: "Pelanggan memiliki riwayat pesanan, nonaktifkan saja alih-alih menghapus"
└── Jika 0 → DELETE (hard delete) dari laundry.customers
```

### ListOrders — riwayat pesanan pelanggan

```
ListOrders(idOutlet, idCustomer, status, fromDate, toDate, page, perPage)
└── JOIN orders + order_items + service_variants langsung di query
    (memperbaiki bug Laravel: relasi Eloquent items.variant yang salah membuat
     GET .../customers/{id}/orders crash/selalu null — di sini join-nya benar
     sehingga nama_layanan & satuan tampil dengan semestinya)
```

---

## Repository (`repositories/customer/customer.go`)

Tabel: `laundry.customers` (baca), `laundry.orders` + `laundry.order_items` +
`laundry.service_variants` (baca, untuk `ListOrders`).

| Fungsi | Keterangan |
|---|---|
| `List` / `Count` | Filter `search` (ILIKE nama/telepon/email) + `is_active`, dengan limit/offset |
| `ExistsPhoneInOutlet` / `ExistsEmailInOutlet` | `COUNT(*) ... AND id_customer <> excludeID` |
| `CountOrders` | Guard sebelum hard-delete |
| `ListOrders` | Raw SQL dengan `json_agg` untuk item per pesanan |

---

## Dependency

Dibaca oleh domain `order` (`GetCustomer().FindByID`) saat membuat pesanan, untuk memastikan
`customer_id` valid dan milik outlet yang sama.
