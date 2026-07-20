# Dokumentasi Laundry Service

> **Backend manajemen laundry multi-outlet** — akun & autentikasi, outlet, metode pembayaran, katalog layanan, pelanggan, pesanan, dashboard, dan pengeluaran.

---

## Daftar Isi

| File | Deskripsi |
|------|-----------|
| [00-arsitektur.md](00-arsitektur.md) | Arsitektur keseluruhan, struktur folder, pola desain, diagram dependency injection, dan alur request |
| [07-infrastruktur.md](07-infrastruktur.md) | Middleware (JWT, validasi, upload S3, datatable), konfigurasi, database |

Dokumentasi per-domain (account, outlet, payment method, catalog, customer, order, dashboard, expense,
report) ditambahkan seiring domain tersebut selesai diimplementasikan.

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
