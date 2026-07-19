# 05 — Group Akses

## Deskripsi

Domain group akses adalah yang paling kompleks. Mengelola:

1. **Master Group Akses** — group (role) yang dimiliki suatu aplikasi (contoh: "Admin", "Dosen", "Mahasiswa")
2. **Group Akses (Modul)** — menentukan modul apa saja yang bisa diakses oleh suatu group
3. **Group Akses User Apps** — menautkan user ke group aplikasi tertentu (user assignment)

Relasi:
```
master_aplikasi ──▶ master_group ──▶ group_akses ◀── master_modul
                                         │
                                   trans_user_group
                                         │
                                       users
```

Menggunakan **1 koneksi database**: `DB` (SSO) — tabel `master_group`, `group_akses`, `trans_user_group`.

---

## Endpoint

### Master Group Akses

| Method | Path | Middleware | Handler |
|--------|------|-----------|---------|
| `POST` | `/groupakses/mastergroup` | `JWTMiddleware`, `ValidateForm(CreateMstGroupAksesForm)` | `CreateMstGroupAkses` |
| `GET` | `/groupakses/mastergroup/:id_master_aplikasi` | `JWTMiddleware`, `ValidatedParams(GetMstGroupAksesParams)` | `GetMstGroupAkses` |
| `GET` | `/groupakses/mastergroup/detail/:id_master_group` | `JWTMiddleware`, `ValidatedParams(GetDetailMstGroupAksesParams)` | `GetDetailMstGroupAkses` |
| `GET` | `/groupakses/mastergroup/modul/:id_master_aplikasi/:id_master_group` | `JWTMiddleware`, `ValidatedParams(GetMstGroupAksesModulParams)` | `GetMstGroupAksesModul` |
| `PUT` | `/groupakses/mastergroup/:id_master_group` | `JWTMiddleware`, `ValidatedParams(UpdateMstGroupAksesParams)`, `ValidateForm(UpdateMstGroupAksesForm)` | `UpdateMstGroupAkses` |
| `DELETE` | `/groupakses/mastergroup/:id_master_group` | `JWTMiddleware`, `ValidatedParams(DeleteMstGroupAksesParams)` | `DeleteMstGroupAkses` |

### Group Akses (Modul Assignment)

| Method | Path | Middleware | Handler |
|--------|------|-----------|---------|
| `POST` | `/groupakses/group` | `JWTMiddleware`, `ValidateForm(CreateGroupAksesForm)` | `CreateGroupAkses` |
| `GET` | `/groupakses/group/:id_master_group` | `JWTMiddleware`, `ValidatedParams(GetGroupAksesParams)` | `GetGroupAkses` |
| `DELETE` | `/groupakses/group/:id_group_akses` | `JWTMiddleware`, `ValidatedParams(DeleteGroupAksesParams)` | `DeleteGroupAkses` |

### Group Akses User Apps (User Assignment)

| Method | Path | Middleware | Handler |
|--------|------|-----------|---------|
| `GET` | `/groupakses/userapps/menu/:id_user/:id_master_aplikasi` | `JWTMiddleware`, `ValidatedParams(GetGroupAksesUserMenuParams)` | `GetGroupAksesUserMenu` |
| `GET` | `/groupakses/userapps/apps/:id_user` | `JWTMiddleware`, `ValidatedParams(GetGroupAksesUserAppsParams)` | `GetGroupAksesUserApps` |
| `POST` | `/groupakses/userapps` | `JWTMiddleware`, `ValidateForm(CreateGroupAksesUserAppsForm)` | `CreateGroupAksesUserApps` |
| `POST` | `/groupakses/userapps/bulk` | `JWTMiddleware`, `ValidateForm(BulkGroupAksesUserAppsForm)` | `BulkCreateGroupAksesUserApps` |
| `PUT` | `/groupakses/userapps/:id_trans_user_group` | `JWTMiddleware`, `ValidatedParams(UpdateGroupAksesUserAppsParams)`, `ValidateForm(UpdateGroupAksesUserAppsForm)` | `UpdateGroupAksesUserApps` |
| `DELETE` | `/groupakses/userapps/:id_trans_user_group` | `JWTMiddleware`, `ValidatedParams(DeleteGroupAksesUserAppsParams)` | `DeleteGroupAksesUserApps` |

---

## DTO

### Master Group

```go
type CreateMstGroupAksesForm struct {
    Id_Master_Aplikasi string `json:"id_master_aplikasi" validate:"required,uuid4"`
    Nama_Group         string `json:"nama_group"         validate:"required"`
    Deskripsi          string `json:"deskripsi"          validate:"required"`
}

type UpdateMstGroupAksesForm struct {
    Nama_Group string `json:"nama_group" validate:"required"`
    Deskripsi  string `json:"deskripsi"  validate:"required"`
}
```

### Group Akses (Modul)

```go
type CreateGroupAksesForm struct {
    Id_Master_Aplikasi string `json:"id_master_aplikasi" validate:"required,uuid4"`
    Id_Master_Group    string `json:"id_master_group"    validate:"required,uuid4"`
    Id_Master_Modul    string `json:"id_master_modul"    validate:"required,uuid4"`
}
```

### User Apps

```go
type CreateGroupAksesUserAppsForm struct {
    Id_User            string `json:"id_user"            validate:"required,uuid4"`
    Id_Master_Aplikasi string `json:"id_master_aplikasi" validate:"required,uuid4"`
    Id_Master_Group    string `json:"id_master_group"    validate:"required,uuid4"`
    Status_Data        string `json:"status_data"        validate:"required,oneof='true' 'false'"`
}

type BulkGroupAksesUserAppsForm struct {
    IDUsers          []string `json:"id_users"           validate:"required,min=1,dive,uuid"`
    IDMasterAplikasi string   `json:"id_master_aplikasi" validate:"required,uuid"`
    IDMasterGroup    string   `json:"id_master_group"    validate:"required,uuid"`
    StatusData       bool     `json:"status_data"        validate:"required"`
}

type UpdateGroupAksesUserAppsForm struct {
    Status_Data string `json:"status_data" validate:"required,oneof='true' 'false'"`
}
```

### Params

```go
type GetMstGroupAksesParams        struct { Id_Master_Aplikasi string }
type GetMstGroupAksesModulParams   struct { Id_Master_Aplikasi, Id_Master_Group string }
type GetGroupAksesParams           struct { Id_Master_Group string }
type UpdateMstGroupAksesParams     struct { Id_Master_Group string }
type DeleteMstGroupAksesParams     struct { Id_Master_Group string }
type GetDetailMstGroupAksesParams  struct { Id_Master_Group string }
type DeleteGroupAksesParams        struct { Id_Group_Akses string }
type GetGroupAksesUserMenuParams   struct { Id_User, Id_Master_Aplikasi string }
type GetGroupAksesUserAppsParams   struct { Id_User string }
type UpdateGroupAksesUserAppsParams struct { Id_Trans_User_Group string }
type DeleteGroupAksesUserAppsParams struct { Id_Trans_User_Group string }
```

---

## Alur Fungsi Kunci

### 1. Get Group Akses User Menu — `GET /groupakses/userapps/menu/:id_user/:id_master_aplikasi`

Fungsi paling penting — mengembalikan menu & modul yang bisa diakses user tertentu di aplikasi tertentu:

```
Controller → Service.GetGroupAksesUserMenu(ctx, idUser, idMasterAplikasi)
  └── Repository.GetGroupAksesUserMenu(ctx, idUser, idMasterAplikasi)
      ├── SQL (JOIN 4 tabel):
      │   SELECT mn.nama_menu, mn.order AS menu_order, mn.icon AS menu_icon,
      │          mm.nama_modul, mm.path AS modul_path, mm.order AS modul_order, mm.icon AS modul_icon
      │   FROM trans_user_group tug
      │   INNER JOIN group_akses ga ON tug.id_master_group = ga.id_master_group
      │   INNER JOIN master_modul mm ON ga.id_master_modul = mm.id_master_modul
      │   INNER JOIN master_menu mn ON mm.id_master_menu = mn.id_master_menu
      │   WHERE tug.id_user=? AND tug.id_master_aplikasi=? AND tug.status_data=true
      │   GROUP BY mn.nama_menu, mn.order, mn.icon, mm.nama_modul, mm.path, mm.order, mm.icon
      │   ORDER BY mn.order, mm.order
      │
      └── middleware.Groupby("nama_menu", result)
          └── Kelompokkan modul-modul di bawah menu-nya masing-masing
          └── Return: {menu: {nama_menu: [...modul]}}
```

**Output:**
```json
{
  "menu": {
    "Dashboard": [
      {"nama_modul": "Overview", "path": "/dashboard", ...}
    ],
    "Master Data": [
      {"nama_modul": "Users", "path": "/master/users", ...},
      {"nama_modul": "Aplikasi", "path": "/master/apps", ...}
    ]
  }
}
```

### 2. Get Group Akses User Apps — `GET /groupakses/userapps/apps/:id_user`

Mengembalikan daftar aplikasi yang bisa diakses user (beserta group-nya):

```
Controller → Service.GetGroupAksesUserApps(ctx, idUser)
  └── Repository.GetGroupAksesUserApps(ctx, idUser)
      ├── SQL (JOIN 3 tabel):
      │   SELECT ma.nama_aplikasi, ma.image, ma.deskripsi, ma.versi_aplikasi,
      │          ma.tgl_version, ma.url, mg.nama_group, tug.id_trans_user_group, tug.status_data
      │   FROM trans_user_group tug
      │   INNER JOIN master_group mg ON tug.id_master_group = mg.id_master_group
      │   INNER JOIN master_aplikasi ma ON mg.id_master_aplikasi = ma.id_master_aplikasi
      │   WHERE tug.id_user=? AND mg.status_data=TRUE AND ma.status_data=TRUE
      │
      └── middleware.Groupby("nama_aplikasi", result)
          └── Return: {apps: {nama_aplikasi: [...group]}}
```

### 3. Bulk Create User Apps — `POST /groupakses/userapps/bulk`

```
Controller → Service.BulkCreateGroupAksesUserApps(ctx, idUsers, idMasterApp, idMasterGroup, statusData)
  └── Repository.BulkCreateGroupAksesUserApps(ctx, idUsers, idMasterApp, idMasterGroup, statusData)
      ├── BEGIN transaction
      ├── Untuk setiap idUser: INSERT INTO trans_user_group (...)
      ├── Jika ada error → ROLLBACK
      └── COMMIT
```

---

## Repository (Query SQL)

### Tabel `master_group`

| Method | Operasi |
|--------|---------|
| `CreateMstGroupAkses` | `INSERT` (id_master_aplikasi, nama_group, deskripsi) |
| `GetMstGroupAkses` | Datatable SELECT (where id_master_aplikasi + status_data=true) |
| `GetMstGroupAksesModul` | CTE + LEFT JOIN modul & group_akses |
| `GetDetailMstGroupAkses` | `SELECT` where id_master_group |
| `UpdateMstGroupAkses` | `UPDATE nama_group, deskripsi` |
| `DeleteMstGroupAkses` | `UPDATE status_data=false` (soft delete) |

### Tabel `group_akses`

| Method | Operasi |
|--------|---------|
| `GetGroupAkses` | CTE kompleks + LEFT JOIN (modul aktif + akses) |
| `CreateGroupAkses` | `INSERT` (id_master_group, id_master_aplikasi, id_master_modul, akses) |
| `DeleteGroupAkses` | Hard `DELETE` where id_group_akses |

### Tabel `trans_user_group`

| Method | Operasi |
|--------|---------|
| `GetGroupAksesUserMenu` | JOIN 4 tabel, GROUP BY, di-group via middleware.Groupby |
| `GetGroupAksesUserApps` | JOIN 3 tabel, GROUP BY nama_aplikasi |
| `CreateGroupAksesUserApps` | `INSERT` single |
| `BulkCreateGroupAksesUserApps` | `INSERT` multiple dalam 1 transaksi |
| `UpdateGroupAksesUserApps` | `UPDATE status_data` |
| `DeleteGroupAksesUserApps` | Hard `DELETE` |

---

## Catatan

- **Query CTE (Common Table Expression)**: `GetGroupAkses` dan `GetMstGroupAksesModul` menggunakan CTE PostgreSQL untuk query yang kompleks — menggabungkan master_modul aktif dengan group_akses yang sudah ada.
- **Groupby middleware**: Hasil flat query dikelompokkan menggunakan `middleware.Groupby(key, data)` — mengelompokkan array of map berdasarkan key tertentu.
- **Bulk insert transaksional**: `BulkCreateGroupAksesUserApps` menggunakan transaksi GORM (`tx.Begin()` ... `tx.Commit()` / `tx.Rollback()`).
- **Cross-domain usage**: Repository `CreateGroupAksesUserApps` juga dipanggil dari service `user` saat generate dosen/tendik.
