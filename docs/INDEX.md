# Dokumentasi Laundry Service

> **Backend manajemen laundry multi-outlet** — akun & autentikasi, outlet, metode pembayaran, katalog layanan, pelanggan, pesanan, dashboard, dan pengeluaran.

---

## Daftar Isi

| File | Deskripsi |
|------|-----------|
| [00-arsitektur.md](00-arsitektur.md) | Arsitektur keseluruhan, struktur folder, pola desain, diagram dependency injection, dan alur request |
| [07-infrastruktur.md](07-infrastruktur.md) | Middleware (JWT, validasi, upload S3, datatable), konfigurasi, database |
| [08-account.md](08-account.md) | Domain akun: registrasi, login, profil, logout, reset password via OTP email |
| [09-outlet.md](09-outlet.md) | Domain outlet: kelola outlet, staff, undang karyawan, role & permission |
| [10-paymentmethod.md](10-paymentmethod.md) | Domain metode pembayaran outlet (cash/transfer/e-wallet) |
| [11-catalog.md](11-catalog.md) | Domain katalog: layanan, varian layanan, parfum, diskon |
| [12-customer.md](12-customer.md) | Domain pelanggan outlet + riwayat pesanan |
| [13-order.md](13-order.md) | Domain pesanan: mesin status, invoice, pembayaran, pickup — inti bisnis laundry |
| [14-dashboard.md](14-dashboard.md) | Domain dashboard: ringkasan operasional & keuangan outlet |
| [15-expense.md](15-expense.md) | Domain pengeluaran outlet |
| [16-report.md](16-report.md) | Domain laporan transaksi outlet |

---

## Cara Membaca Dokumentasi Ini

Setiap dokumen fitur mengikuti struktur yang sama:

1. **Deskripsi** — apa yang dilakukan domain ini
2. **Endpoint** — daftar lengkap route HTTP (method, path, middleware, handler)
3. **DTO (Data Transfer Object)** — struct request/response
4. **Alur Fungsi** — bagaimana request mengalir dari Route → Controller → Service → Repository
5. **Repository (Query SQL)** — query mentah yang digunakan
6. **Dependency** — apa saja yang dibutuhkan domain ini

---

## Teknologi

| Teknologi | Penggunaan |
|-----------|-----------|
| **Go 1.23+** | Bahasa pemrograman utama |
| **Fiber v2** | HTTP framework (mirip Express.js) |
| **GORM** | Akses raw SQL ke PostgreSQL |
| **PostgreSQL** | Database, domain laundry memakai schema `laundry` |
| **JWT (golang-jwt)** | Access token & refresh token |
| **bcrypt** | Hashing password |
| **S3 (IDCloudHost)** | Object storage untuk logo outlet & foto layanan |
| **go-playground/validator** | Validasi form & params input |
