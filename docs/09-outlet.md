# 09 — Outlet

## Deskripsi

Domain outlet: satu akun (`laundry.users`) bisa memiliki/bergabung ke banyak outlet, dengan role
`owner` atau `karyawan` per outlet (tabel penghubung `laundry.user_outlets`). Domain ini juga
mengatur undangan karyawan baru dan permission granular per anggota.

**Role & permission default** (disimpan sebagai JSONB di `user_outlets.permissions_json`):

| Permission | owner | karyawan |
|---|---|---|
| create_order | true | true |
| cancel_order | true | true |
| create_expense | true | true |
| manage_customers | true | true |
| manage_services | true | false |
| manage_employees | true | false |
| view_revenue | true | false |
| view_report_tx | true | false |
| view_report_finance | true | false |
| view_report_customer | true | false |

> Permission ini **disimpan & dikembalikan** di response tapi **tidak ditegakkan** (enforced) di
> endpoint lain — persis seperti perilaku Laravel aslinya. Otorisasi endpoint yang sesungguhnya
> memakai `RequireOutletMember` (anggota aktif) atau `RequireOutletOwner` (khusus role owner).

---

## Endpoint

| Method | Path | Middleware | Handler |
|--------|------|-----------|---------|
| `GET` | `/outlets` | `LaundryJWTMiddleware` | `ListMyOutlets` |
| `POST` | `/outlets` | `LaundryJWTMiddleware`, `ValidateForm(CreateOutletForm)` | `CreateOutlet` |
| `GET` | `/outlets/:id_outlet` | `RequireOutletMember` | `GetOutlet` |
| `PUT` | `/outlets/:id_outlet` | `RequireOutletOwner`, `ValidateForm(UpdateOutletForm)` | `UpdateOutlet` |
| `GET` | `/outlets/:id_outlet/staff` | `RequireOutletMember` | `ListStaff` |
| `POST` | `/outlets/:id_outlet/invite` | `RequireOutletOwner`, `ValidateForm(InviteEmployeeForm)` | `InviteEmployee` |
| `PUT` | `/outlets/:id_outlet/members/:id_member` | `RequireOutletOwner`, `ValidateForm(UpdateMemberForm)` | `UpdateMember` |
| `DELETE` | `/outlets/:id_outlet/members/:id_member` | `RequireOutletOwner` | `RemoveMember` |

Semua route berada di bawah group `/outlets` dengan `LaundryJWTMiddleware` (mengisi
`c.Locals("id_user")`); route yang butuh `:id_outlet` menambah `RequireOutletMember`/`RequireOutletOwner`
yang membaca param tersebut dan memvalidasi ke `laundry.user_outlets`.

---

## DTO (`domain/dto/outlet.go`)

```go
type CreateOutletForm struct {
    Nama, Alamat, Telepon, LogoPath string // Nama required max 120, sisanya opsional
}

type UpdateOutletForm struct { // sama seperti CreateOutletForm
    Nama, Alamat, Telepon, LogoPath string
}

type InviteEmployeeForm struct {
    Nama, Telepon, Email               string
    Password, PasswordKonfirmasi       string // min 8, harus sama
    Role                               string // "owner" | "karyawan"
    Alamat                             string // required
    Permissions                        map[string]bool // override opsional
}

type UpdateMemberForm struct {
    Role        string          // opsional — kosong berarti role tidak berubah
    Permissions map[string]bool // opsional — di-merge dengan default role
}
```

---

## Alur Fungsi

### 1. Create Outlet

```
Service.CreateOutlet(idUser, form)
├── repository.CreateOutlet() → INSERT laundry.outlets (id_owner_user = idUser)
├── Hitung default permission role "owner" (semua true)
└── repository.AssignUser(idOutlet, idUser, "owner", permissions)
    → INSERT laundry.user_outlets ... ON CONFLICT (id_user, id_outlet) DO UPDATE
```

### 2. Invite Employee

```
Service.InviteEmployee(idOutlet, form)
├── Cari user existing by email → repository.GetAccount().FindUserByEmail()
│   ├── Jika ADA → pakai akun itu (nama/password/alamat dari form DIABAIKAN,
│   │              persis seperti Laravel: findOrCreateByEmailOrPhone)
│   └── Jika TIDAK ADA →
│       ├── Cek telepon belum dipakai akun lain → ExistsPhone()
│       ├── Hash password, CreateUser() (via repository account, reuse lintas domain)
├── Hitung permission: default(role) di-merge dengan form.Permissions
└── AssignUser(idOutlet, idUser, role, permissions) — upsert keanggotaan
```

### 3. Update Member / Remove Member

```
UpdateMember: ambil role saat ini kalau form.Role kosong → merge permission → UpdateMember()

RemoveMember:
├── Jika anggota berrole "owner" → CountActiveOwners(idOutlet)
│   └── Jika hanya tersisa 1 owner aktif → tolak ("Tidak dapat menghapus owner terakhir")
└── DeactivateMember() → UPDATE user_outlets SET is_active = false (soft remove)
```

---

## Repository (`repositories/outlet/outlet.go`)

Tabel: `laundry.outlets`, `laundry.user_outlets`.

| Fungsi | Keterangan |
|---|---|
| `ListOutletsForUser` | Join outlets + user_outlets, filter `is_active = true` di kedua sisi |
| `AssignUser` | `INSERT ... ON CONFLICT (id_user, id_outlet) DO UPDATE` — upsert keanggotaan |
| `ListStaff` | Join ke `laundry.users`, urut owner dulu lalu nama |
| `CountActiveOwners` | Guard sebelum menghapus anggota berrole owner |

---

## Dependency

- `repositories.IRepositoryRegistry.GetAccount()` — dipakai langsung dari service outlet untuk
  mencari/membuat akun karyawan saat invite (lintas domain lewat registry, bukan HTTP call)
- `middleware.RequireOutletMember` / `RequireOutletOwner` (`middlewares/laundry_auth.go`)
- `middleware.HashPassword`
