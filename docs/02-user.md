# 02 — User

## Deskripsi

Domain user menangani manajemen user SSO secara penuh: CRUD user, melihat daftar dosen/mahasiswa dari database eksternal, bulk create mahasiswa, dan generate user otomatis (mahasiswa/dosen/tendik).

Domain ini **satu-satunya yang menggunakan 3 koneksi database sekaligus**:
- `DB` (SSO) — tabel `users`
- `DBAkademik` — tabel `list_mahasiswa`, `mahasiswa_angkatan`
- `DBDigiclass` — tabel `master_dosen`

---

## Endpoint

Semua endpoint di bawah prefix `/users` dan dilindungi `middleware.JWTMiddleware` (kecuali jika ada pengecualian khusus).

| Method | Path | Middleware | Handler |
|--------|------|-----------|---------|
| `GET` | `/users/` | `JWTMiddleware`, `ValidatedParams2(GetUsersForm)` | `GetUsers` |
| `POST` | `/users/` | `JWTMiddleware`, `ValidateForm(CreateUsersForm)` | `CreateUsers` |
| `GET` | `/users/:id_user/` | `JWTMiddleware`, `ValidatedParams(GetUserParams)` | `GetUser` |
| `PUT` | `/users/:id_user` | `JWTMiddleware`, `ValidatedParams(UpdateUsersParams)`, `ValidateForm(UpdateUsersForm)` | `UpdateUsers` |
| `PUT` | `/users/updatepassword/:id_user` | `JWTMiddleware`, `ValidatedParams(UpdatePasswordParams)`, `ValidateForm(UpdatePasswordForm)` | `UpdatePassword` |
| `DELETE` | `/users/:id_user` | `JWTMiddleware`, `ValidatedParams(DeleteUserParams)` | `DeleteUser` |
| `GET` | `/users/get_dosen` | `JWTMiddleware` | `GetUsersDosen` |
| `GET` | `/users/get_detail_dosen/:person_id` | `JWTMiddleware`, `ValidatedParams(GetIDUserParams)` | `GetDetailDosen` |
| `GET` | `/users/get_mahasiswa/:id_prodi` | `JWTMiddleware`, `ValidatedParams(GetProdiParams)` | `GetUsersMahasiswa` |
| `GET` | `/users/get_mahasiswa_data/:id_prodi` | `JWTMiddleware`, `ValidatedParams(GetProdiParams)` | `GetUsersMahasiswaData` |
| `GET` | `/users/get_detail_mahasiswa/:id_registrasi_mahasiswa` | `JWTMiddleware`, `ValidatedParams(GetMahasiswaParams)` | `GetDetailMahasiswa` |
| `POST` | `/users/bulk_mahasiswa` | `JWTMiddleware` | `BulkCreateUsersMahasiswa` |
| `POST` | `/users/generate` | `JWTMiddleware` | `GenerateUserMahasiswa` |
| `POST` | `/users/generate_dosen` | `JWTMiddleware` | `GenerateUserDosen` |
| `POST` | `/users/generate_tendik` | `JWTMiddleware` | `GenerateUserTendik` |

---

## DTO

### Request — Form

```go
type GetUsersForm struct {
    Limit  int    `json:"limit"  validate:"required,numeric,oneof=10 25 50 100" default:"10"`
    Offset int    `json:"offset" validate:"numeric"`
    Order  string `json:"order"`
    Filter string `json:"filter"`
}

type CreateUsersForm struct {
    Email        string `json:"email"        validate:"required,email"`
    Id_Person    string `json:"id_person"    validate:"required,uuid4"`
    Jenis_User   string `json:"jenis_user"   validate:"required,oneof='dosen' 'tenaga pendidik' 'mahasiswa' 'orang tua' 'perseptor'"`
    Nama_Lengkap string `json:"nama_lengkap" validate:"required"`
    No_Hp        string `json:"no_hp"        validate:"required,numeric"`
    Username     string `json:"username"     validate:"required"`
    Password     string `json:"password"     validate:"required"`
}

type UpdateUsersForm struct {
    Email        string `json:"email"        validate:"required,email"`
    Id_Person    string `json:"id_person"    validate:"required,uuid4"`
    Jenis_User   string `json:"jenis_user"   validate:"required,oneof='dosen' 'tenaga tendidik' 'mahasiswa' 'orang tua' 'perseptor'"`
    Nama_Lengkap string `json:"nama_lengkap" validate:"required"`
    No_Hp        string `json:"no_hp"        validate:"required,numeric"`
    Username     string `json:"username"     validate:"required"`
}

type UpdatePasswordForm struct {
    Password string `json:"password" validate:"required"`
}
```

### Request — Params (URL path)

```go
type GetUserParams          struct { Id_User string `validate:"required,uuid4"` }
type UpdateUsersParams      struct { Id_User string `validate:"required,uuid4"` }
type UpdatePasswordParams   struct { Id_User string `validate:"required,uuid4"` }
type DeleteUserParams       struct { Id_User string `validate:"required,uuid4"` }
type GetIDUserParams        struct { Person_ID string `validate:"required,uuid4"` }
type GetProdiParams         struct { ID_Prodi string `validate:"required,uuid4"` }
type GetMahasiswaParams     struct { ID_Registrasi_Mahasiswa string `validate:"required,uuid4"` }
```

### Data Struct (internal)

```go
type Mahasiswa struct {
    Email, IdPerson, JenisUser, Nama, NoHp, Username, Password, Avatar string
}

type Pegawai struct {
    IDUser, Password string
}

type DaftarUser struct {
    PersonID, NamaLengkap, NIM string
}

type DetailUser struct {
    PersonID, NamaLengkap, Email, NoHp, Username, JenisUser string
}

type BulkMahasiswaItem struct {
    Email, IdPerson, NamaLengkap, NoHp, Username, Password string
}

type BulkCreateResult struct {
    Total, Berhasil, Gagal int
    DetailGagal            []string
}
```

---

## Alur Fungsi

### 1. Get Users (Datatable) — `GET /users/`

```
Controller: GetUsers()
├── Ambil GetUsersForm dari c.Locals("validatedForm")
│   └── Diparsing dari query string oleh ValidatedParams2
├── Service: GetUsers(order, filter, limit, offset)
│   ├── Repository: GetUsers(order, filter, limit, offset)
│   │   └── Bangun query SQL datatable
│   │   └── middleware.Datatables(sRecursive, sTable, order, sFilter, limit, offset)
│   │       ├── SELECT count(*) → recordsTotal
│   │       ├── SELECT * ... LIMIT ? OFFSET ? → data
│   │       └── SELECT count(*) filtered → recordsFiltered
│   │   └── Return: {recordsTotal, recordsFiltered, data}
│   └── Return: Response{Success: true, Data: hasil}
└── Return JSON
```

### 2. Create User — `POST /users/`

```
Controller: CreateUsers()
├── [Opsional] Upload avatar ke S3: middleware.FileUploadToS3Middleware(c, "avatar", "avatars", ...)
│   └── Jika ada file → fileName disimpan di c.Locals("fileName")
│
├── Ambil CreateUsersForm dari c.Locals("validatedForm")
├── Service: CreateUser(ctx, email, idPerson, jenisUser, namaLengkap, noHp, username, password, avatar)
│   ├── Repository: CountByPerson(ctx, idPerson)
│   │   └── SQL: SELECT COUNT(*) FROM users WHERE id_person=? AND status_data=true
│   ├── Jika count > 0 → "User dengan id_person sudah ada"
│   ├── Hash password: middleware.HashPassword(password) → bcrypt
│   ├── Repository: Insert(ctx, email, idPerson, jenisUser, namaLengkap, noHp, username, passwordHash, avatar)
│   │   └── SQL: INSERT INTO users (email, id_person, jenis_user, nama_lengkap, no_hp,
│   │            username, password, first_login, status_data, avatar) VALUES (...)
│   └── Return: Response{Success: true}
└── Return JSON
```

### 3. Get User Detail — `GET /users/:id_user/`

```
Controller: GetUser()
├── Ambil GetUserParams dari c.Locals("validatedParams")
├── Service: GetUser(ctx, idUser)
│   ├── Repository: FindByID(ctx, idUser)
│   │   └── SQL: SELECT id_user, username, email, nama_lengkap, avatar, id_person,
│   │            jenis_user, no_hp, password FROM users WHERE id_user=? AND status_data=true
│   ├── Jika error → "Gagal Mendapatkan users"
│   ├── Jika rows==0 → "users tidak ada"
│   └── Return: Response{Success: true, Data: user}
└── Return JSON
```

### 4. Update User — `PUT /users/:id_user`

```
Controller: UpdateUsers()
├── [Opsional] Upload avatar ke S3
├── Ambil UpdateUsersParams + UpdateUsersForm
├── Service: UpdateUser(ctx, idUser, avatar, email, idPerson, jenisUser, namaLengkap, noHp, username)
│   ├── JIKA avatar == "":
│   │   └── Repository: UpdateWithoutAvatar(ctx, email, idPerson, jenisUser, namaLengkap, noHp, username, idUser)
│   │       └── SQL: UPDATE users SET email=?, id_person=?, jenis_user=?, nama_lengkap=?, no_hp=?, username=?
│   │                WHERE id_user=?
│   └── JIKA avatar != "":
│       └── Repository: UpdateWithAvatar(ctx, avatar, email, idPerson, jenisUser, namaLengkap, noHp, username, idUser)
│           └── SQL: UPDATE users SET avatar=?, email=?, id_person=?, ... WHERE id_user=?
│   └── Return Response
└── Return JSON
```

### 5. Update Password — `PUT /users/updatepassword/:id_user`

```
Controller: UpdatePassword()
├── Ambil UpdatePasswordParams + UpdatePasswordForm
├── Service: UpdatePassword(ctx, idUser, password)
│   ├── Repository: FindFirstLogin(ctx, idUser)
│   │   └── SQL: SELECT id_user, first_login FROM users WHERE id_user=? AND status_data=true
│   ├── Hash password baru: middleware.HashPassword(password)
│   ├── JIKA first_login == true:
│   │   └── UpdatePasswordFirstLogin(passwordHash, idUser)
│   │       └── SQL: UPDATE users SET password=?, tgl_update=NOW(), first_login=false, tgl_first_login=NOW()
│   └── JIKA first_login == false:
│       └── UpdatePasswordNormal(passwordHash, idUser)
│           └── SQL: UPDATE users SET password=?, tgl_update=NOW()
│   └── Return Response
└── Return JSON
```

### 6. Delete User — `DELETE /users/:id_user`

```
Controller: DeleteUser()
├── Ambil DeleteUserParams
├── Service: DeleteUser(ctx, idUser)
│   └── Repository: Delete(ctx, idUser)
│       └── SQL: DELETE FROM users WHERE id_user=?
└── Return JSON
```

### 7. Bulk Create Mahasiswa — `POST /users/bulk_mahasiswa`

```
Controller: BulkCreateUsersMahasiswa()
├── BodyParser → []dto.BulkMahasiswaItem (parse manual, tidak pakai ValidateForm)
├── Service: BulkCreateUsersMahasiswa(ctx, items)
│   ├── Iterasi items, untuk setiap item panggil CreateUser()
│   │   └── (pakai jenisUser="mahasiswa", avatar="")
│   ├── Kumpulkan hasil: berhasil, gagal, detail_gagal
│   └── Return: BulkCreateResult
└── Return JSON
```

### 8. Generate User Mahasiswa — `POST /users/generate`

```
Controller: GenerateUserMahasiswa()
├── Service: GenerateUserMahasiswa(ctx)
│   ├── Repository: GetMahasiswaSourceForGenerate(ctx)
│   │   └── SQL (DBAkademik): SELECT ... FROM list_mahasiswa
│   │       WHERE nama_status_mahasiswa='AKTIF' AND id_periode_masuk='20251'
│   ├── Untuk setiap mahasiswa → CreateUser()
│   │   └── password = NIM, username = NIM
│   └── Return Response
└── Return JSON
```

### 9. Generate User Dosen — `POST /users/generate_dosen`

```
Controller: GenerateUserDosen()
├── Service: GenerateUserDosen(ctx)
│   ├── Repository: GetPegawaiNeedingPassword(ctx, "Dosen")
│   │   └── SQL (DB): SELECT id_user, password FROM users
│   │       WHERE jenis_user=? AND status_data=true AND password=username
│   ├── Untuk setiap dosen:
│   │   ├── UpdatePassword(idUser, password)
│   │   └── CreateGroupAksesUserApps(idUser, appID, groupDosen, "true")
│   │       └── (memanggil repository groupakses — cross-domain)
│   └── Return Response
└── Return JSON
```

### 10. Generate User Tendik — `POST /users/generate_tendik`

```
Sama seperti GenerateUserDosen, tapi:
├── GetPegawaiNeedingPassword(ctx, "Tenaga Pendidik")
└── Group akses berbeda: groupTendik
```

---

## Repository (Query SQL)

| Method | DB | Tabel | Operasi |
|--------|----|-------|---------|
| `GetUsers` | DB | `users` | Datatable SELECT |
| `CountByPerson` | DB | `users` | `SELECT COUNT(*)` |
| `Insert` | DB | `users` | `INSERT` |
| `FindByID` | DB | `users` | `SELECT` where id_user |
| `UpdateWithoutAvatar` | DB | `users` | `UPDATE` 6 kolom |
| `UpdateWithAvatar` | DB | `users` | `UPDATE` 7 kolom |
| `FindFirstLogin` | DB | `users` | `SELECT first_login` |
| `UpdatePasswordFirstLogin` | DB | `users` | `UPDATE password + first_login` |
| `UpdatePasswordNormal` | DB | `users` | `UPDATE password` |
| `Delete` | DB | `users` | `DELETE` |
| `GetDosen` | DBDigiclass | `master_dosen` | `SELECT` all |
| `GetDetailDosen` | DBDigiclass | `master_dosen` | `SELECT` where id |
| `GetMahasiswa` | DBAkademik | `list_mahasiswa` | `SELECT` where id_prodi |
| `GetDetailMahasiswa` | DBAkademik | `list_mahasiswa` | `SELECT *` where id |
| `GetMahasiswaData` | DBAkademik | `list_mahasiswa` + `mahasiswa_angkatan` | Datatable paginasi |
| `GetMahasiswaSourceForGenerate` | DBAkademik | `list_mahasiswa` | `SELECT` aktif periode 20251 |
| `GetPegawaiNeedingPassword` | DB | `users` | `SELECT` where password=username |

---

## Catatan Khusus

- **Avatar upload**: opsional, hanya diproses jika `c.FormFile("avatar")` tidak error. Upload ke S3 IDCloudHost.
- **Cross-domain call**: `GenerateUserDosen` dan `GenerateUserTendik` memanggil `repository.GetGroupAkses().CreateGroupAksesUserApps()` — ini adalah cross-domain dependency di level service.
- **GetMahasiswaData**: query paling kompleks — join `list_mahasiswa` + `mahasiswa_angkatan`, dengan filter teks, paginasi LIMIT/OFFSET, dan menghitung `recordsTotal` vs `recordsFiltered` secara terpisah.
- **BulkCreate**: tidak pakai transaksi database — setiap item di-insert satu per satu lewat `CreateUser`. Gagal di satu item tidak membatalkan item lain.
