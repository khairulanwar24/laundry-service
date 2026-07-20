# 16 — Report (Laporan Transaksi)

## Deskripsi

Domain laporan transaksi outlet — sengaja dibuat **tipis** (satu endpoint) karena secara
fungsional tumpang tindih dengan listing pesanan di domain `order` (`GET /outlets/:id/orders`).
Ini meniru `ReportController` versi Laravel yang juga cuma pembungkus filter tambahan di atas query
yang mirip `OrderController::index`.

Bedanya dengan `order.List`: `report.Transactions` memfilter berdasarkan `status_pembayaran` dan
rentang tanggal `checkin_at` eksplisit (`tanggal_mulai`/`tanggal_akhir`), bukan `tab` (kategori
status yang sudah dipetakan) + pencarian teks bebas.

---

## Endpoint

| Method | Path | Middleware | Query params |
|--------|------|-----------|---------------|
| `GET` | `/outlets/:id_outlet/reports/transactions` | `LaundryJWTMiddleware`, `RequireOutletMember` | `status_pesanan`, `status_pembayaran`, `tanggal_mulai`, `tanggal_akhir`, `page`, `per_page` |

Response:
```json
{
  "success": true,
  "data": {
    "data": [{"id_order","invoice_no","status","status_pembayaran","subtotal","nilai_diskon","total","checkin_at","nama_pelanggan","telepon_pelanggan"}],
    "meta": {"current_page": 1, "per_page": 15, "total": 0}
  }
}
```

---

## Repository (`repositories/report/report.go`)

Query terhadap `laundry.orders` JOIN `laundry.customers`, dengan filter opsional:
`o.status = ?`, `o.status_pembayaran = ?`, `o.checkin_at >= tanggal_mulai::date`,
`o.checkin_at < (tanggal_akhir::date + INTERVAL '1 day')`.

---

## Dependency

Tidak ada dependency lintas domain di level kode (query independen terhadap `orders`/`customers`),
tapi secara konsep melayani kebutuhan yang sama dengan domain `order` — kalau kebutuhan filter
laporan berkembang lebih jauh, pertimbangkan menyatukan ke domain `order` alih-alih menduplikasi
lebih banyak query di sini.
