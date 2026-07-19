# 04 — MstMenu (Master Menu & Modul)

## Deskripsi

Domain master menu & modul mengelola struktur navigasi untuk setiap aplikasi. Satu aplikasi punya banyak **menu**, dan satu menu punya banyak **modul** (halaman/feature). Modul inilah yang nantinya di-assign ke group akses.

Relasi:
```
master_aplikasi (1) ──▶ (N) master_menu (1) ──▶ (N) master_modul
```

Menggunakan **1 koneksi database**: `DB` (SSO) — tabel `master_menu` dan `master_modul`.

---

## Endpoint

### Menu

| Method | Path | Middleware | Handler |
|--------|------|-----------|---------|
| `GET` | `/mstmenu/menu/:id_master_aplikasi` | `JWTMiddleware`, `ValidatedParams(GetMstMenuParams)` | `GetMstMenu` |
| `GET` | `/mstmenu/menu/detail/:id_master_menu` | `JWTMiddleware`, `ValidatedParams(GetDetailMstMenuParams)` | `GetDetailMstMenu` |
| `POST` | `/mstmenu/menu` | `JWTMiddleware`, `ValidateForm(CreateMstMenuForm)` | `CreateMstMenu` |
| `PUT` | `/mstmenu/menu/:id_master_menu` | `JWTMiddleware`, `ValidatedParams(UpdateMstMenuParams)`, `ValidateForm(UpdateMstMenuForm)` | `UpdateMstMenu` |
| `DELETE` | `/mstmenu/menu/:id_master_menu` | `JWTMiddleware`, `ValidatedParams(DeleteMstMenuParams)` | `DeleteMstMenu` |

### Modul

| Method | Path | Middleware | Handler |
|--------|------|-----------|---------|
| `GET` | `/mstmenu/modul/:id_master_menu` | `JWTMiddleware`, `ValidatedParams(GetMstMenuModulParams)` | `GetMstMenuModul` |
| `GET` | `/mstmenu/modul/detail/:id_master_modul` | `JWTMiddleware`, `ValidatedParams(GetDetailMstMenuModulParams)` | `GetDetailMstMenuModul` |
| `POST` | `/mstmenu/modul` | `JWTMiddleware`, `ValidateForm(CreateMstMenuModulForm)` | `CreateMstMenuModul` |
| `PUT` | `/mstmenu/modul/:id_master_modul` | `JWTMiddleware`, `ValidatedParams(UpdateMstMenuModulParams)`, `ValidateForm(UpdateMstMenuModulForm)` | `UpdateMstMenuModul` |
| `DELETE` | `/mstmenu/modul/:id_master_modul` | `JWTMiddleware`, `ValidatedParams(DeleteMstMenuModulParams)` | `DeleteMstMenuModul` |

---

## DTO

### Menu — Form

```go
type CreateMstMenuForm struct {
    Id_Master_Aplikasi string `json:"id_master_aplikasi" validate:"required,uuid4"`
    Nama_Menu          string `json:"nama_menu"          validate:"required"`
    Deskripsi          string `json:"deskripsi"          validate:"required"`
    Order              string `json:"order"              validate:"required"`
    Icon               string `json:"icon"               validate:"required"`
}

type UpdateMstMenuForm struct {
    Nama_Menu string `json:"nama_menu" validate:"required"`
    Deskripsi string `json:"deskripsi" validate:"required"`
    Order     string `json:"order"     validate:"required,numeric"`
    Icon      string `json:"icon"      validate:"required"`
}
```

### Modul — Form

```go
type CreateMstMenuModulForm struct {
    Id_Master_Aplikasi string `json:"id_master_aplikasi" validate:"required,uuid4"`
    Id_Master_Menu     string `json:"id_master_menu"     validate:"required,uuid4"`
    Nama_Modul         string `json:"nama_modul"         validate:"required"`
    Path               string `json:"path"               validate:"required"`
    Deskripsi          string `json:"deskripsi"          validate:"required"`
    Order              string `json:"order"              validate:"required"`
    Icon               string `json:"icon"               validate:"required"`
}

type UpdateMstMenuModulForm struct {
    Id_Master_Aplikasi string `json:"id_master_aplikasi" validate:"required,uuid4"`
    Id_Master_Menu     string `json:"id_master_menu"     validate:"required,uuid4"`
    Nama_Modul         string `json:"nama_modul"         validate:"required"`
    Path               string `json:"path"               validate:"required"`
    Deskripsi          string `json:"deskripsi"          validate:"required"`
    Order              string `json:"order"              validate:"required"`
    Icon               string `json:"icon"               validate:"required"`
}
```

### Params

```go
type GetMstMenuParams           struct { Id_Master_Aplikasi string `validate:"required,uuid4"` }
type GetDetailMstMenuParams     struct { Id_Master_Menu string    `validate:"required,uuid4"` }
type GetMstMenuModulParams      struct { Id_Master_Menu string    `validate:"required,uuid4"` }
type GetDetailMstMenuModulParams struct { Id_Master_Modul string  `validate:"required,uuid4"` }
type UpdateMstMenuParams        struct { Id_Master_Menu string    `validate:"required,uuid4"` }
type UpdateMstMenuModulParams   struct { Id_Master_Modul string  `validate:"required,uuid4"` }
type DeleteMstMenuParams        struct { Id_Master_Menu string    `validate:"required,uuid4"` }
type DeleteMstMenuModulParams   struct { Id_Master_Modul string  `validate:"required,uuid4"` }
```

### Form Paginasi (query string)

```go
type GetMstMenuForm struct {
    Limit  int    `json:"limit"  validate:"required,numeric,oneof=10 25 50 100" default:"10"`
    Offset int    `json:"offset" validate:"numeric"`
    Order  string `json:"order"`
    Filter string `json:"filter"`
}

type GetMstMenuModulForm struct {
    Limit  int    `json:"limit"  validate:"required,numeric,oneof=10 25 50 100" default:"10"`
    Offset int    `json:"offset" validate:"numeric"`
    Order  string `json:"order"`
    Filter string `json:"filter"`
}
```

---

## Alur Fungsi

### 1. Get Menu List (Datatable) — `GET /mstmenu/menu/:id_master_aplikasi`

```
Controller → Service.GetMstMenu(ctx, idMasterAplikasi, limit, offset, order, filter)
  └── Repository.GetMstMenu(idMasterAplikasi, limit, offset, order, filter)
      ├── Query: SELECT id_master_menu, id_master_aplikasi, nama_menu, deskripsi, "order", icon
      │          FROM master_menu WHERE id_master_aplikasi=? AND status_data=true
      └── middleware.Datatables() → {recordsTotal, recordsFiltered, data}
```

### 2. Get Menu Detail — `GET /mstmenu/menu/detail/:id_master_menu`

```
Controller → Service.GetDetailMstMenu(ctx, idMasterMenu)
  └── Repository.GetDetailMstMenu(ctx, idMasterMenu)
      └── SQL: SELECT ... FROM master_menu WHERE id_master_menu=?
```

### 3. Create Menu — `POST /mstmenu/menu`

```
Controller → Service.CreateMstMenu(ctx, idMasterAplikasi, namaMenu, deskripsi, order, icon)
  └── Repository.CreateMstMenu(ctx, idMasterAplikasi, namaMenu, deskripsi, order, icon)
      └── SQL: INSERT INTO master_menu (id_master_aplikasi, nama_menu, deskripsi, "order", icon) VALUES (...)
```

### 4. Update Menu — `PUT /mstmenu/menu/:id_master_menu`

```
Controller → Service.UpdateMstMenu(ctx, idMasterMenu, namaMenu, deskripsi, order, icon)
  └── Repository.UpdateMstMenu(ctx, idMasterMenu, namaMenu, deskripsi, order, icon)
      └── SQL: UPDATE master_menu SET nama_menu=?, deskripsi=?, "order"=?, icon=?
              WHERE id_master_menu=?
```

### 5. Delete Menu — `DELETE /mstmenu/menu/:id_master_menu`

```
Controller → Service.DeleteMstMenu(ctx, idMasterMenu)
  └── Repository.DeleteMstMenu(ctx, idMasterMenu)
      └── SQL: UPDATE master_menu SET status_data=false WHERE id_master_menu=?
      → SOFT DELETE
```

> Modul mengikuti pola yang sama: Get list (datatable by id_master_menu), Get detail, Create, Update, Delete (soft delete).

---

## Repository (Query SQL)

### Tabel `master_menu`

| Method | Operasi |
|--------|---------|
| `GetMstMenu` | Datatable SELECT (where id_master_aplikasi + status_data=true) |
| `GetDetailMstMenu` | `SELECT` where id_master_menu |
| `CreateMstMenu` | `INSERT` (5 kolom) |
| `UpdateMstMenu` | `UPDATE` 4 kolom |
| `DeleteMstMenu` | `UPDATE status_data=false` (soft delete) |

### Tabel `master_modul`

| Method | Operasi |
|--------|---------|
| `GetMstMenuModul` | Datatable SELECT (where id_master_menu + status_data=true) |
| `GetDetailMstMenuModul` | `SELECT` where id_master_modul |
| `CreateMstMenuModul` | `INSERT` (7 kolom) |
| `UpdateMstMenuModul` | `UPDATE` 7 kolom |
| `DeleteMstMenuModul` | `UPDATE status_data=false` (soft delete) |

---

## Catatan

- **Soft delete**: Menu dan modul tidak dihapus permanen — hanya `status_data` di-set `false`. Berbeda dengan `master_aplikasi` yang hard delete.
- **Kolom `order`**: Nama kolom bentrok dengan keyword SQL, jadi di-escape dengan double quote (`"order"`).
- **Paginasi**: Semua endpoint GET list menggunakan `ValidatedParams2` (query string parsing) untuk `limit`, `offset`, `order`, `filter`.
- **Path modul**: String path frontend (contoh: `/dashboard`, `/admin/users`) — digunakan untuk routing di frontend.
