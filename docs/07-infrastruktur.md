# 07 — Infrastruktur

Dokumen ini mencakup semua komponen pendukung di luar 6 domain bisnis: middleware, konfigurasi, database, gRPC, dan utilitas.

---

## Middleware

Semua middleware ada di `middlewares/` (10 file). Berikut daftar dan fungsinya:

### 1. JWT (`jwt.go`)

Tiga secret key dan jenis token:

| Fungsi | Secret | Expiry | Return |
|--------|--------|--------|--------|
| `AccessToken(username, id_user)` | `jwtSecret` | 5 menit | `(token, 200/400)` |
| `RefreshToken(username, id_user)` | `jwtSecret2` | 8 jam | `(token, 200/400)` |
| `AccessTokenOTP(username, id_user)` | `jwtSecret3` | 5 menit | `(token, 200/400)` |
| `OTPResetpassword(id_reset_password, percobaan)` | `jwtSecret3` | 10 menit | `(token, 200/400)` |

**Claims struct:**
```go
type JwtCustomClaims struct {
    Username string `json:"username"`
    Id_user  string `json:"id_user"`
    jwt.RegisteredClaims
}

type JwtOTPCustomClaims struct {
    IDResetPassword string `json:"id_reset_password"`
    Percobaan       int32  `json:"percobaan"`
    jwt.RegisteredClaims
}
```

**Verifikasi & Parsing:**

| Fungsi | Input | Output |
|--------|-------|--------|
| `JWTMiddleware(c)` | Fiber context | `c.Next()` atau 401 |
| `Checkjwt(c)` | Fiber context | `bool` (Bearer header) |
| `CheckjwtRefresh(c)` | Fiber context | `bool` (cookie "refresh_token") |
| `CheckjwtOTP(c)` | Fiber context | `bool` (cookie "OTP") |
| `ValidateTokenMiddleware(token)` | Token string | `(bool, message)` — untuk gRPC |
| `JWTParse(c)` | Fiber context | `types.Response` — decode refresh token + buat access token baru |
| `ParseToken(token)` | Token string | `*JwtCustomClaims` |
| `ParseTokenOTP(c)` | Fiber context | `*JwtOTPCustomClaims` |
| `ParseTokenChangePassword(token)` | Token string | `*JwtCustomClaims` |

### 2. Validasi (`validation.go`)

Tiga middleware Fiber untuk validasi input:

| Middleware | Dari | Disimpan di | Untuk |
|-----------|------|-------------|-------|
| `ValidateForm(form)` | `c.BodyParser()` | `c.Locals("validatedForm")` | POST/PUT body |
| `ValidatedParams(form)` | `c.ParamsParser()` | `c.Locals("validatedParams")` | URL path params |
| `ValidatedParams2(form)` | `c.QueryParser()` | `c.Locals("validatedForm")` | Query string params |

Menggunakan library `go-playground/validator/v10`. `getErrorMessage()` menerjemahkan tag validasi ke pesan bahasa Indonesia.

### 3. Password (`password.go`)

```go
HashPassword(password) → (hash string, error)    // bcrypt, cost 14
CheckPasswordHash(password, hash) → bool          // bcrypt compare
```

### 4. Datatable (`datatable.go`)

Helper untuk server-side datatable (DataTables jQuery):
```go
Datatables(sRecursive, sTable, order, sFilter string, limit, offset int) → map[string]interface{}
```
Output: `{recordsTotal, recordsFiltered, data}`

> ⚠️ Masih menggunakan `database.DB` global (tidak pakai GORM WithContext), jadi belum menerima context dan tidak bisa menggunakan transaction. Perlu refactor ke depan.

### 5. Email (`mail.go`)

- `Mail(to, subject, body)` — kirim email via HTTP POST ke `https://api.farmasiunissula.com/mail/send-email`
  - Auth: Basic Auth hardcoded
  - Content-Type: application/json
- `PublishMessageToQueue(to, subject, body)` — alternatif via RabbitMQ (tidak digunakan saat ini)

### 6. S3 Upload (`s3.go`)

```go
FileUploadToS3Middleware(c *fiber.Ctx, formFieldName, folderName string,
    allowedTypes []string, maxFileSize int64) error
```

- Upload file ke S3 IDCloudHost
- Validasi tipe file (MIME) dan ukuran maksimum
- Generate nama file unik: `{folderName}/{uuid}_{timestamp}{ext}`
- Set ACL Public Read
- Simpan nama file di `c.Locals("fileName")`

### 7. Groupby (`groupby.go`)

```go
Groupby(key string, data []map[string]interface{}) → map[string][]map[string]interface{}
```
Mengelompokkan array of map berdasarkan nilai dari key tertentu. Digunakan untuk mengelompokkan modul per menu (`GetGroupAksesUserMenu`) dan app per aplikasi (`GetGroupAksesUserApps`).

### 8. Middleware Dasar (`middleware.go`)

```go
Logger(c *fiber.Ctx) error  // set request ID (tidak terpakai aktif)
```

### 9. WhatsApp (`whatsapp.go`)

File ada tetapi tidak digunakan dalam kode saat ini.

### 10. Upload (`upload.go`)

File ada tetapi tidak digunakan — upload file ditangani langsung oleh `s3.go`.

---

## Konfigurasi (`config/config.go`)

```go
type Config struct {
    DBHost, DBPort, DBUser, DBPassword, DBName, SSLMode string
}
```

Tiga fungsi loader:
- `LoadConfig()` — Database SSO (default: `rdp.farmasiunissula.com:2345`, user `sso`)
- `LoadConfigAkademik()` — Database Akademik
- `LoadConfigDigiclass()` — Database Digiclass

Menggunakan environment variable dengan fallback default value via `getEnv(key, defaultValue)`.

---

## Database (`database/database.go`)

### Variabel Global

```go
var DB, DBAkademik, DBDigiclass *gorm.DB
var SqlDB, SqlDBAkademik, SqlDBDigiclass *sql.DB
```

### Inisialisasi

1. `Connect()` — koneksi ke database SSO + Akademik
   - Schema: `public.` (SSO), `akademik.` (Akademik)
   - Pool: MaxIdle 10, MaxOpen 100, MaxLifetime 1 jam
   - Auto-create schema `akademik` jika belum ada

2. `ConnectDigiclass()` — koneksi ke database Digiclass
   - Dipanggil terpisah dari `Connect()`
   - Pool settings sama

---

## gRPC (`grpc/server.go`)

### Service: `AuthServiceServer`

Satu RPC method:
```protobuf
service AuthService {
    rpc ValidateToken(ValidateTokenRequest) returns (ValidateTokenResponse);
}
```

**Implementasi:**
```go
func (s *AuthServiceServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
    isValid, message := middleware.ValidateTokenMiddleware(req.Token)
    return &pb.ValidateTokenResponse{IsValid: isValid, Message: message}, nil
}
```

Jalan di port `50051` (default), digunakan microservice lain untuk memvalidasi JWT token via gRPC.

> **Catatan**: gRPC server belum mengikuti pola registry — masih akses langsung ke `middleware.ValidateTokenMiddleware`. Bisa di-refactor nanti.

---

## Response Standard

### `common/response/response.go` (pakai di semua domain baru)

```go
type Response struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}
```

Helper: `Success(c, data)` → 200 OK, `Error(c, code, message)` → status custom.

### `types/response.go` (legacy, masih dipakai middlewares)

```go
type Response struct { ... }  // sama strukturnya
type GetData struct { ... }   // pagination params
```

Middleware (`validation.go`, `jwt.go`) masih menggunakan `types.Response`. Perlu diseragamkan ke `common/response` ke depannya.

---

## Entry Point (`main.go`)

```
1. godotenv.Load("app/.env")
2. database.Connect()           → DB, SqlDB, DBAkademik, SqlDBAkademik
3. database.ConnectDigiclass()  → DBDigiclass, SqlDBDigiclass
4. Inisialisasi Fiber app (CORS)
5. Dependency Injection:
   repository := repositories.NewRepositoryRegistry(DB, DBAkademik, DBDigiclass)
   service    := services.NewServiceRegistry(repository)
   controller := controllers.NewControllerRegistry(service)
   routes.NewRouteRegistry(controller, app).Serve()
6. Goroutine: app.Listen(":3001")    // HTTP
7. Main thread: grpcServer.Serve()    // gRPC :50051
```

---

## File Konfigurasi Lain

| File | Fungsi |
|------|--------|
| `app/.env` | Env variables (DB, AWS, RabbitMQ, APP_ENV) |
| `app/.env_dev` | Template env development |
| `app/.env_prod` | Template env production |
| `Dockerfile` | Build container (multi-stage) |
| `laundry-service.dockerfile` | Alternate dockerfile |
| `docker-compose.yml` | Local dev setup |
| `jenkinsfile` | CI/CD pipeline |
| `makefile` | Shortcut commands |
| `go.mod` | Dependencies (GORM, Fiber, JWT, bcrypt, S3, RabbitMQ, gRPC, protobuf) |

---

## Diagram Infrastruktur

```
                           ┌──────────────┐
                           │   Nginx /    │
                           │  API Gateway │
                           └──────┬───────┘
                                  │
                    ┌─────────────┼─────────────┐
                    │             │             │
              Port 3001      Port 50051    External API
              (Fiber)        (gRPC)        (Mail Service)
                    │             │             │
           ┌───────┴───────┐     │      ┌──────┴──────┐
           │  Middlewares   │     │      │  RabbitMQ   │
           │  (JWT, Valid,  │     │      │  (optional) │
           │   S3, etc.)    │     │      └─────────────┘
           └───────┬───────┘     │
                   │             │
           ┌───────┴───────┐     │
           │  6 Domains     │     │
           │  (Clean Arch)  │◄────┘
           └───────┬───────┘
                   │
        ┌──────────┼──────────┐
        │          │          │
   PostgreSQL  PostgreSQL  PostgreSQL
   (SSO)       (Akademik)  (Digiclass)
   public      akademik    public
        │          │          │
        └──────────┼──────────┘
                   │
           ┌───────┴───────┐
           │  S3 (IDCloud) │
           │  Avatars &    │
           │  Images       │
           └───────────────┘
```
