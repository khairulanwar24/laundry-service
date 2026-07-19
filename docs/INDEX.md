# Dokumentasi SSO Service (ios-service)

> **Microservice Single Sign-On** — autentikasi, otorisasi, manajemen user, aplikasi, menu, modul, dan group akses untuk ekosistem Farmasi UNISSULA.

---

## Daftar Isi

| File | Deskripsi |
|------|-----------|
| [00-arsitektur.md](00-arsitektur.md) | Arsitektur keseluruhan, struktur folder, pola desain, diagram dependency injection, dan alur request |
| [01-auth.md](01-auth.md) | Domain autentikasi: login, logout, reset password, OTP, ganti password, refresh token |
| [02-user.md](02-user.md) | Domain user: CRUD user, bulk create mahasiswa, generate user (mahasiswa/dosen/tendik) |
| [03-masterapp.md](03-masterapp.md) | Domain master aplikasi: kelola aplikasi yang terdaftar di SSO |
| [04-mstmenu.md](04-mstmenu.md) | Domain master menu & modul: struktur navigasi per aplikasi |
| [05-groupakses.md](05-groupakses.md) | Domain group akses: master group, akses modul, user-apps assignment |
| [06-ref.md](06-ref.md) | Domain referensi: prodi & angkatan (data lookup) |
| [07-infrastruktur.md](07-infrastruktur.md) | Middleware (JWT, validasi, upload S3, email, datatable), konfigurasi, database, gRPC |

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
| **Go 1.21+** | Bahasa pemrograman utama |
| **Fiber v2** | HTTP framework (mirip Express.js) |
| **GORM** | ORM untuk PostgreSQL |
| **PostgreSQL** | Database (3 koneksi: SSO, Akademik, Digiclass) |
| **gRPC / Protobuf** | Token validation service antar microservice |
| **JWT (golang-jwt)** | Access token, refresh token, OTP token |
| **bcrypt** | Hashing password |
| **S3 (IDCloudHost)** | Object storage untuk avatar & image |
| **RabbitMQ** | Antrean email (opsional, ada cadangan HTTP API) |
| **go-playground/validator** | Validasi form & params input |
