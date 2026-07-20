# 00 — Arsitektur

## Struktur Folder

```
laundry-service/
├── main.go                     # Entry point: inisialisasi DB, registry, gRPC, HTTP server
├── app/.env                    # Environment variables
├── config/                     # Konfigurasi database (3 koneksi: SSO, Akademik, Digiclass)
├── database/                   # Koneksi database (Connect, ConnectDigiclass)
├── domain/dto/                 # Data Transfer Object — struct request/response shared antar layer
├── repositories/               # Lapisan akses data (RAW SQL via GORM, per domain)
├── services/                   # Lapisan logika bisnis (orkestrasi, validasi, hashing)
├── controllers/                # Lapisan HTTP handler (Fiber) — parsing input, panggil service
├── routes/                     # Registrasi endpoint HTTP per domain
├── middlewares/                # JWT, validasi, upload S3, datatable, email, groupby, hashing
├── common/response/            # Struktur response JSON standar
├── types/                      # Legacy type (Response), masih dipakai middlewares
├── grpc/                       # gRPC server untuk ValidateToken
├── proto/                      # File hasil generate protobuf (auth_grpc.pb.go, auth.pb.go)
├── docs/                       # Swagger docs (auto-generated)
├── Dockerfile / docker-compose.yml
├── jenkinsfile
└── makefile
```

---

## Pola Arsitektur: Clean Architecture 4-Layer

```
┌─────────────────────────────────────────────────────────────┐
│                       HTTP Request                          │
│                           │                                 │
│                           ▼                                 │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐     │
│  │   routes/   │───▶│ controllers/│───▶│  services/  │     │
│  │  (Router)   │    │  (Handler)  │    │  (Business) │     │
│  └─────────────┘    └─────────────┘    └──────┬──────┘     │
│                                                │            │
│                                                ▼            │
│                                         ┌─────────────┐     │
│                                         │repositories/ │     │
│                                         │  (Database)  │     │
│                                         └──────┬──────┘     │
│                                                │            │
│                                                ▼            │
│                                           PostgreSQL        │
└─────────────────────────────────────────────────────────────┘
```

### Aturan Dependensi

- **Route** hanya tahu **Controller**
- **Controller** hanya tahu **Service**
- **Service** hanya tahu **Repository**
- **Repository** hanya tahu **Database**
- Semua layer tahu **DTO** (`domain/dto/`)
- **Tidak ada** layer atas yang tahu layer bawah secara langsung — selalu lewat **interface**

---

## Registry Pattern (Dependency Injection Manual)

Setiap layer punya **Registry** yang menyediakan akses ke domain via interface:

```
main.go
  │
  ├─ database.Connect()          → DB, DBAkademik, DBDigiclass
  ├─ repositories.NewRegistry()  → IRepositoryRegistry
  ├─ services.NewRegistry()      → IServiceRegistry
  ├─ controllers.NewRegistry()   → IControllerRegistry
  └─ routes.NewRegistry()        → IRouteRegistry.Serve()
```

### Alur Pembuatan:

```go
// 1. Repository Registry
repository := repositories.NewRepositoryRegistry(DB, DBAkademik, DBDigiclass)

// 2. Service Registry (butuh Repository)
service := services.NewServiceRegistry(repository)

// 3. Controller Registry (butuh Service)
controller := controllers.NewControllerRegistry(service)

// 4. Route Registry (butuh Controller + Fiber Router)
routes.NewRouteRegistry(controller, app).Serve()
```

### Diagram Dependency:

```
IRepositoryRegistry ◄── IServiceRegistry ◄── IControllerRegistry ◄── IRouteRegistry
       │                      │                       │                     │
  Registry{DB}           Registry{Repo}         Registry{Svc}         Registry{Ctrl,Router}
       │                      │                       │                     │
  GetAuth()             GetAuth()               GetAuth()               Serve()
  GetUser()             GetUser()               GetUser()               authRoute()
  GetRef()              GetRef()                GetRef()                userRoute()
  GetMasterApp()        GetMasterApp()          GetMasterApp()          ...dll
  GetMstMenu()          GetMstMenu()            GetMstMenu()
  GetGroupAkses()       GetGroupAkses()         GetGroupAkses()
```

---

## Alur Request Lengkap (Contoh: Login)

```
1. HTTP POST /auth/login {username, password}
       │
       ▼
2. routes/auth/auth.go — Run()
   api.Post("/login", middleware.ValidateForm(&dto.LoginForm{}), controller.GetAuth().LoginController)
       │
       ├── Middleware ValidateForm:
       │   • BodyParser → LoginForm struct
       │   • validator.Struct() → cek tag validate
       │   • Simpan ke c.Locals("validatedForm")
       │
       ▼
3. controllers/auth/auth.go — LoginController()
   • Ambil form: c.Locals("validatedForm").(*dto.LoginForm)
   • Panggil: ctrl.service.GetAuth().Login(ctx, username, password)
       │
       ▼
4. services/auth/auth.go — Login()
   • Panggil: s.repository.GetAuth().FindUserByUsername(ctx, username)
   • Bandingkan hash password: middleware.CheckPasswordHash(password, hash)
   • Kembalikan response.Response{Success: true, Data: user}
       │
       ▼
5. repositories/auth/auth.go — FindUserByUsername()
   • Eksekusi RAW SQL via GORM:
     SELECT username, nama_lengkap, email, no_hp, avatar, id_user, id_person,
            first_login, password, jenis_user
     FROM users WHERE username = ?
   • Kembalikan map[string]interface{}
       │
       ▼
6. Kembali ke Controller — LoginController():
   • Jika sukses:
     • Buat refresh_token (JWT 8 jam, HTTP-only cookie)
     • Buat access_token  (JWT 5 menit, di response body)
     • Return JSON: {success, message, data: {user, access_token}}
```

---

## 3 Koneksi Database

| Koneksi | Schema | Penggunaan |
|---------|--------|-----------|
| `database.DB` | `public` | SSO utama: users, master_aplikasi, master_menu, master_modul, master_group, group_akses, trans_user_group, password_reset |
| `database.DBAkademik` | `akademik` | Data mahasiswa: list_mahasiswa, master_angkatan_mahasiswa, master_prodi |
| `database.DBDigiclass` | `public` | Data dosen: master_dosen |

Dominan domain hanya pakai `database.DB`, kecuali:
- **user** — pakai ketiganya (DB+DBAkademik+DBDigiclass)
- **ref** — pakai DBAkademik saja

---

## Interface vs Struct

Setiap domain di tiap layer memakai pola yang sama:

```go
// Interface (kontrak) — inilah yang di-inject ke registry
type IXxxRepository interface { ... }
type IXxxService    interface { ... }
type IXxxController interface { ... }
type IXxxRoute      interface { Run() }

// Struct privat — implementasi
type XxxRepository struct { db *gorm.DB }
type XxxService    struct { repository repositories.IRepositoryRegistry }
type XxxController struct { service    services.IServiceRegistry }
type XxxRoute      struct { controller controllers.IControllerRegistry; router fiber.Router }

// Constructor
func NewXxxRepository(db *gorm.DB) IXxxRepository { ... }
func NewXxxService(repo repositories.IRepositoryRegistry) IXxxService { ... }
func NewXxxController(svc services.IServiceRegistry) IXxxController { ... }
func NewXxxRoute(ctrl controllers.IControllerRegistry, router fiber.Router) IXxxRoute { ... }
```

---

## Konvensi Kode

- **Nama file**: lowercase, satu file per layer (`auth.go`, `user.go`)
- **Nama package di subdirektori**: sama dengan nama direktori parent (`package controllers`, `package services`, `package repositories`, `package routes`)
- **Query SQL**: RAW string via `db.Raw()` (bukan GORM query builder)
- **Parameter multi-baris**: fungsi dengan 4+ parameter ditulis satu per baris
- **Context propagation**: `ctx context.Context` selalu parameter pertama di repository & service
- **Response**: selalu `response.Response{Success, Message, Data}` dari `common/response/`
- **Error handling**: error dikembalikan ke caller, tidak di-swallow (kecuali kasus sengaja untuk keamanan seperti di reset password)
