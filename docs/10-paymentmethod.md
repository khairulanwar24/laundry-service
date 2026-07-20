# 10 — Payment Method

## Deskripsi

Domain metode pembayaran milik outlet (kas, transfer bank, e-wallet, dll). Dipakai saat mencatat
pembayaran pesanan di domain `order`.

---

## Endpoint

Semua di bawah `/outlets/:id_outlet/payment-methods` dengan `LaundryJWTMiddleware` +
`RequireOutletMember` di level group; create/update/delete tambahan mensyaratkan
`RequireOutletOwner` (hanya owner yang boleh mengubah).

| Method | Path | Middleware tambahan | Handler |
|--------|------|---------------------|---------|
| `GET` | `/outlets/:id_outlet/payment-methods` | - | `ListActive` |
| `POST` | `/outlets/:id_outlet/payment-methods` | `RequireOutletOwner`, `ValidateForm(PaymentMethodForm)` | `Create` |
| `PUT` | `/outlets/:id_outlet/payment-methods/:id_payment_method` | `RequireOutletOwner`, `ValidateForm(PaymentMethodForm)` | `Update` |
| `DELETE` | `/outlets/:id_outlet/payment-methods/:id_payment_method` | `RequireOutletOwner` | `Deactivate` |

`ListActive` hanya mengembalikan metode dengan `is_active = true` (tidak ada endpoint untuk
melihat yang nonaktif) — sama seperti perilaku Laravel aslinya.

---

## DTO (`domain/dto/paymentmethod.go`)

```go
type PaymentMethodForm struct {
    Kategori    string   // required, oneof: cash | transfer | e_wallet
    Nama        string   // required, max 120
    Logo        string   // opsional
    NamaPemilik string   // opsional (misal nama pemilik rekening)
    Tags        []string // opsional
    IsActive    *bool    // opsional, default true kalau kosong
}
```

Dipakai untuk **create maupun update** — PUT mengirim payload penuh (bukan partial update),
konsisten dengan `UpdatePaymentMethodRequest` versi Laravel.

---

## Alur Fungsi

```
Create/Update:
├── Default is_active = true kalau form.IsActive nil
├── Default tags = [] kalau nil, di-marshal ke JSON string lalu di-cast ?::jsonb saat INSERT/UPDATE
└── Update juga memverifikasi payment method milik outlet yang benar (FindByID(idOutlet, id))
    sebelum mengubah — kalau tidak cocok, dianggap "tidak ditemukan" (404 semantik, bukan 403)

Deactivate: soft-delete — UPDATE is_active = false, bukan DELETE baris.
```

---

## Repository (`repositories/paymentmethod/paymentmethod.go`)

Tabel: `laundry.payment_methods`. Kolom `tags` bertipe JSONB, di-passing sebagai string JSON dan
di-cast eksplisit `?::jsonb` di query raw SQL.

---

## Dependency

- `middleware.RequireOutletOwner`
