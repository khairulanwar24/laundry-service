# 14 — Dashboard

## Deskripsi

Domain analitik/ringkasan operasional & keuangan outlet — murni baca (read-only), tidak ada
mutasi data. Semua endpoint (kecuali `revenue-summary` & `profit-loss-summary`) menerima query
param opsional `?date=YYYY-MM-DD` (default hari ini di timezone lokal server).

---

## Endpoint

Semua di bawah `/outlets/:id_outlet/dashboard` dengan `LaundryJWTMiddleware` + `RequireOutletMember`.

| Method | Path | Query params |
|--------|------|---------------|
| `GET` | `/summary` | `date` |
| `GET` | `/daily-report` | `date` |
| `GET` | `/financial-summary` | `date` |
| `GET` | `/revenue-summary` | - (selalu relatif ke "sekarang") |
| `GET` | `/profit-loss-summary` | - |
| `GET` | `/customer-report` | `search`, `status` (all/active/inactive), `date_from`, `date_to` |
| `GET` | `/customer-search` | sama seperti di atas + `page`, `per_page` |

---

## Response per endpoint

### `GET /summary`
```json
{
  "date": "2026-07-21", "masuk": 0, "harus_selesai": 0, "terlambat": 0, "item_diambil": 0,
  "omset": 0, "pendapatan": 0, "pengeluaran": 0
}
```
- `masuk`: pesanan dibuat pada tanggal itu, status ≠ BATAL
- `harus_selesai`: `eta_at` jatuh pada tanggal itu, status belum SELESAI/BATAL
- `terlambat`: `eta_at` sudah lewat **dari sekarang** (bukan dari tanggal filter), status belum SELESAI/BATAL
- `item_diambil`: `diambil_at` jatuh pada tanggal itu
- `omset`: total nilai pesanan dibuat pada tanggal itu (status ≠ BATAL)
- `pendapatan`: total pembayaran SUCCESS yang `dibayar_at`-nya jatuh pada tanggal itu
- `pengeluaran`: total `laundry.expenses` pada tanggal itu

### `GET /daily-report`
```json
{
  "date": "...", "total_transaksi_masuk": 0, "total_transaksi_selesai": 0,
  "total_transaksi_batal": 0, "total_pendapatan": 0, "total_omset": 0, "total_pengeluaran": 0
}
```
`total_transaksi_selesai`/`_batal` dihitung dari **kapan transisi status terjadi**
(`order_status_histories.waktu_perubahan`), bukan kapan pesanan dibuat — beda semantik dengan
`summary` di atas, persis seperti versi Laravel.

### `GET /financial-summary`
```json
{ "date": "...", "total_keuangan": 0, "pendapatan": 0, "pengeluaran": 0, "transaksi_terakhir": [...] }
```
`transaksi_terakhir`: gabungan 5 transaksi terbaru dari pembayaran & pengeluaran hari itu (union
query, diurutkan waktu descending, dipotong ke 5 teratas).

### `GET /revenue-summary`
```json
{
  "periode": {"hari_ini": "...", "bulan_ini": "...", "tahun_ini": "..."},
  "omset": {"hari_ini": 0, "bulan_ini": 0, "tahun_ini": 0},
  "total_transaksi": {...},
  "rata_rata_transaksi": {...},
  "pertumbuhan": {
    "hari_ini": {"nominal": 0, "persentase": 0},
    "bulan_ini": {...}, "tahun_ini": {...}
  },
  "top_layanan": [{"nama_layanan": "...", "total_omset": 0, "total_transaksi": 0}],
  "top_pelanggan": [{"nama_pelanggan": "...", "total_omset": 0, "total_transaksi": 0}]
}
```
`pertumbuhan` membandingkan periode berjalan vs periode sebelumnya (hari ini vs kemarin, bulan ini
vs bulan lalu, tahun ini vs tahun lalu). Rumus (`calculateGrowth`): kalau nilai sebelumnya 0,
persentase langsung 100 (jika nilai sekarang > 0) atau 0; selain itu
`(sekarang - sebelumnya) / sebelumnya * 100`. `top_layanan`/`top_pelanggan` dihitung dari bulan
berjalan, maksimal 5 baris, diurutkan omset terbesar.

### `GET /profit-loss-summary`
```json
{
  "laba": {...}, "margin_persen": {...}, "roi_persen": {...},
  "ringkasan_keuangan": {"hari_ini": {"pendapatan","pengeluaran","laba_bersih"}, ...},
  "rincian_pengeluaran": {"bulan_ini": [{"kategori","total"}]},
  "perbandingan_bulanan": {"bulan_ini": 0, "tiga_bulan_lalu": 0},
  "metrik_utama": {"break_event_point_nominal": 0, "cost_per_transaction_nominal": 0}
}
```
- `margin_persen` = laba / omset * 100
- `roi_persen` = laba / "modal awal" (pengeluaran 3 bulan lalu, dipakai sebagai heuristik kasar —
  bukan pelacakan modal sesungguhnya, sama seperti pendekatan Laravel aslinya)
- `break_event_point_nominal` = total pengeluaran bulan berjalan pada kategori biaya tetap yang
  di-whitelist: `Sewa`, `Gaji`, `Listrik`, `Air`, `Internet`

### `GET /customer-report` & `/customer-search`
```json
{
  "total_pelanggan": {"semua": 0, "aktif": 0, "non_aktif": 0},
  "filter": {"search": "", "status": "all", "date_from": "", "date_to": ""},
  "daftar_pelanggan": [{"id","nama","nomor_telepon","email","alamat","status","tanggal_bergabung"}]
}
```
`customer-report` mengembalikan **semua** baris yang cocok filter (tidak dipaginasi).
`customer-search` menambahkan `data` (baris terpaginasi) + `meta` (`current_page`, `per_page`,
`total`) di samping `total_pelanggan` & `filter` yang sama.

---

## Repository (`repositories/dashboard/dashboard.go`)

Query agregat langsung terhadap `laundry.orders`, `laundry.payments`, `laundry.expenses`,
`laundry.order_status_histories`, `laundry.customers` — banyak memakai
`COUNT(*) FILTER (WHERE ...)` Postgres supaya beberapa metrik bisa diambil dalam satu round-trip
query (`OrdersSummaryOnDate`), bukan query terpisah per metrik.

---

## Dependency

Tidak ada dependency lintas domain di level kode — semua query langsung ke tabel di atas.
