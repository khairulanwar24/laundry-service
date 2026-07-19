# 06 — Ref (Referensi)

## Deskripsi

Domain referensi menyediakan data lookup (pencarian) untuk prodi dan angkatan mahasiswa. Data diambil dari **database Akademik** (bukan database SSO).

Fungsi utama: memberi data ke frontend untuk dropdown/pilihan saat membuat user mahasiswa atau filter data.

Menggunakan **1 koneksi database**: `DBAkademik` — tabel `master_prodi` dan `master_angkatan_mahasiswa`.

---

## Endpoint

| Method | Path | Middleware | Handler |
|--------|------|-----------|---------|
| `GET` | `/ref/get_prodi` | `JWTMiddleware` | `GetMasterProdi` |
| `GET` | `/ref/get_angkatan/:id_prodi` | `JWTMiddleware`, `ValidatedParams(GetAngkatanParams)` | `GetAngkatan` |

---

## DTO

```go
type GetAngkatanParams struct {
    IDProdi string `json:"id_prodi" validate:"required,uuid4"`
}
```

---

## Alur Fungsi

### 1. Get Prodi — `GET /ref/get_prodi`

```
Controller → Service.GetMasterProdi(ctx)
  └── Repository.GetMasterProdi(ctx)
      └── SQL (DBAkademik): SELECT * FROM master_prodi
      └── Return: []map[string]any
```

### 2. Get Angkatan — `GET /ref/get_angkatan/:id_prodi`

```
Controller → Service.GetAngkatan(ctx, idProdi)
  └── Repository.GetAngkatan(ctx, idProdi)
      └── SQL (DBAkademik):
          SELECT id_angkatan_mahasiswa, nama_angkatan_mahasiswa
          FROM master_angkatan_mahasiswa
          WHERE id_prodi=? AND status_data=true
          ORDER BY nama_angkatan_mahasiswa DESC
      └── Return: []map[string]any
```

---

## Repository (Query SQL)

| Method | DB | Tabel | Operasi |
|--------|----|-------|---------|
| `GetMasterProdi` | DBAkademik | `master_prodi` | `SELECT *` |
| `GetAngkatan` | DBAkademik | `master_angkatan_mahasiswa` | `SELECT` where id_prodi + status_data=true |

---

## Catatan

- Domain paling sederhana — hanya 2 endpoint read-only.
- Tidak ada operasi create/update/delete.
- **Pilot migrasi**: Domain ref adalah yang pertama dimigrasi ke clean architecture (commit `bb2055a`), menjadi template untuk 5 domain lainnya.
- Response langsung dari repository tanpa transformasi bisnis tambahan di service layer.
