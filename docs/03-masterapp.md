# 03 — Master App (Master Aplikasi)

## Deskripsi

Domain master aplikasi mengelola aplikasi yang terdaftar di SSO. Setiap aplikasi punya nama, deskripsi, versi, URL, dan gambar (icon/logo). Aplikasi inilah yang nantinya punya menu, modul, dan group akses.

Menggunakan **1 koneksi database**: `DB` (SSO) — tabel `master_aplikasi`.

---

## Endpoint

| Method | Path | Middleware | Handler |
|--------|------|-----------|---------|
| `GET` | `/masterapp/` | `JWTMiddleware`, `ValidatedParams2(GetData)` | `GetMasterApps` |
| `GET` | `/masterapp/:id_master_aplikasi` | `JWTMiddleware`, `ValidatedParams(GetMasterAppByIdParams)` | `GetMasterAppById` |
| `POST` | `/masterapp/` | `JWTMiddleware`, `ValidateForm(CreateMasterAppForm)` | `CreateMasterApp` |
| `PUT` | `/masterapp/:id_master_aplikasi` | `JWTMiddleware`, `ValidatedParams(UpdateMasterAppsParams)`, `ValidateForm(UpdateMasterAppsForm)` | `UpdateMasterApp` |
| `DELETE` | `/masterapp/:id_master_aplikasi` | `JWTMiddleware`, `ValidatedParams(DeleteMasterAppParams)` | `DeleteMasterApp` |

---

## DTO

### Request — Form

```go
type GetData struct {
    Limit  int    `json:"limit"  validate:"required,numeric,oneof=10 25 50 100" default:"10"`
    Offset int    `json:"offset" validate:"numeric"`
    Order  string `json:"order"`
    Filter string `json:"filter"`
    Params string `json:"params"`
}

type CreateMasterAppForm struct {
    Nama_Aplikasi  string `json:"nama_aplikasi"  validate:"required"`
    Deskripsi      string `json:"deskripsi"`
    Versi_Aplikasi string `json:"versi_aplikasi" validate:"required"`
    Tgl_Version    string `json:"tgl_version"    validate:"required,datetime=2006-01-02"`
    Url            string `json:"url"            validate:"required,url"`
}

type UpdateMasterAppsForm struct {
    Nama_Aplikasi  string `json:"nama_aplikasi"  validate:"required"`
    Deskripsi      string `json:"deskripsi"      validate:"required"`
    Versi_Aplikasi string `json:"versi_aplikasi" validate:"required"`
    Tgl_Version    string `json:"tgl_version"    validate:"required"`
    URL            string `json:"url"            validate:"required"`
}
```

### Request — Params

```go
type GetMasterAppByIdParams struct {
    Id_master_aplikasi string `validate:"required,uuid4"`
}

type UpdateMasterAppsParams struct {
    Id_master_aplikasi string `validate:"required,uuid4"`
}

type DeleteMasterAppParams struct {
    Id_Master_Aplikasi string `validate:"required,uuid4"`
}
```

---

## Alur Fungsi

### 1. Get All (Datatable) — `GET /masterapp/`

```
Controller → Service.GetMasterApps(order, filter, limit, offset)
  └── Repository.GetMasterApps(order, filter, limit, offset)
      ├── Bangun query SQL datatable dari tabel master_aplikasi (status_data=true)
      ├── Filter: LOWER(nama_aplikasi) LIKE / LOWER(deskripsi) LIKE / ...
      └── middleware.Datatables() → {recordsTotal, recordsFiltered, data}
```

### 2. Get By ID — `GET /masterapp/:id_master_aplikasi`

```
Controller → Service.GetMasterAppById(ctx, id)
  └── Repository.GetMasterAppById(ctx, id)
      └── SQL: SELECT id_master_aplikasi, nama_aplikasi, deskripsi, versi_aplikasi,
                tgl_version, url, image FROM master_aplikasi
                WHERE id_master_aplikasi=? AND status_data=true
```

### 3. Create — `POST /masterapp/`

```
Controller → Service.CreateMasterApp(ctx, nama, deskripsi, tglVersion, url, versi, image)
  ├── [Opsional] Upload image ke S3
  └── Repository.Create(ctx, namaAplikasi, deskripsi, tglVersion, url, versiAplikasi, image)
      └── SQL: INSERT INTO master_aplikasi (nama_aplikasi, deskripsi, tgl_version, url,
              versi_aplikasi, image, status_data) VALUES (?, ?, ?, ?, ?, ?, true)
```

### 4. Update — `PUT /masterapp/:id_master_aplikasi`

```
Controller → Service.UpdateMasterApp(ctx, id, image, nama, deskripsi, versi, tglVersion, url)
  ├── Jika image kosong → Repository.UpdateWithoutImage(...)
  │   └── SQL: UPDATE master_aplikasi SET nama_aplikasi=?, deskripsi=?, versi_aplikasi=?,
  │            tgl_version=?, url=? WHERE id_master_aplikasi=?
  └── Jika image ada   → Repository.UpdateWithImage(...)
      └── SQL: UPDATE ... SET image=?, nama_aplikasi=?, ...
```

### 5. Delete — `DELETE /masterapp/:id_master_aplikasi`

```
Controller → Service.DeleteMasterApp(ctx, id)
  └── Repository.Delete(ctx, id)
      └── SQL: DELETE FROM master_aplikasi WHERE id_master_aplikasi=?
```

---

## Repository (Query SQL)

| Method | Tabel | Operasi |
|--------|-------|---------|
| `GetMasterApps` | `master_aplikasi` | Datatable SELECT (status_data=true) |
| `GetMasterAppById` | `master_aplikasi` | `SELECT` where id + status_data=true |
| `Create` | `master_aplikasi` | `INSERT` (7 kolom + status_data) |
| `UpdateWithoutImage` | `master_aplikasi` | `UPDATE` 5 kolom |
| `UpdateWithImage` | `master_aplikasi` | `UPDATE` 6 kolom (termasuk image) |
| `Delete` | `master_aplikasi` | Hard `DELETE` |

---

## Catatan

- **Delete bersifat hard delete** (`DELETE FROM`), bukan soft delete — berbeda dengan mstmenu dan groupakses yang pakai `UPDATE status_data=false`.
- **Image upload opsional** — sama seperti avatar di user. Upload ke S3 via `middleware.FileUploadToS3Middleware`.
- **Tanggal versi** menggunakan format `YYYY-MM-DD` (divalidasi dengan tag `datetime=2006-01-02`).
