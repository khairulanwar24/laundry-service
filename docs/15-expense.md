# 15 — Expense (Pengeluaran)

## Deskripsi

Domain pencatatan pengeluaran outlet (sewa, gaji, listrik, dll — kategori berupa teks bebas, bukan
enum tertutup). Dipakai juga oleh domain `dashboard` untuk menghitung pendapatan bersih, rincian
per kategori, dan break-even point.

---

## Endpoint

Semua di bawah `/outlets/:id_outlet/expenses` dengan `LaundryJWTMiddleware` + `RequireOutletMember`.

| Method | Path | Query/Body |
|--------|------|-----------|
| `GET` | `/` | `kategori`, `start_date`, `end_date`, `page`, `per_page` |
| `POST` | `/` | `ExpenseForm` |
| `GET` | `/:id_expense` | - |
| `PUT` | `/:id_expense` | `ExpenseForm` |
| `DELETE` | `/:id_expense` | - |

---

## DTO (`domain/dto/expense.go`)

```go
type ExpenseForm struct {
    Kategori  string  // required, max 100 (teks bebas)
    Jumlah    float64 // required, >= 0
    Deskripsi string  // required
    Tanggal   string  // required, format "2006-01-02"
}
```

---

## Alur Fungsi

```
Create/Update:
└── validateTanggalNotFuture(form.Tanggal) → tanggal tidak boleh di masa depan
    (mengikuti validasi StoreExpenseRequest Laravel: date must be <= today)

Delete:
└── Soft delete — UPDATE status_data = false (konsisten dengan konvensi soft-delete domain lain,
    bukan hard DELETE)
```

---

## Repository (`repositories/expense/expense.go`)

Tabel: `laundry.expenses`. Filter list memakai kombinasi `kategori` (exact match), `start_date`/
`end_date` (rentang `tanggal`), dengan pagination limit/offset.

---

## Dependency

Dibaca oleh domain `dashboard` (`GetDashboard` memanggil query terhadap `laundry.expenses` secara
langsung, bukan lewat repository expense — lihat `repositories/dashboard/dashboard.go`) untuk
`pengeluaran`, `rincian_pengeluaran`, dan `break_event_point_nominal`.
