# 07 — Infrastruktur

Dokumen ini mencakup semua komponen pendukung di luar domain bisnis: middleware, konfigurasi, database, dan utilitas.

---

## Middleware

Semua middleware ada di `middlewares/`. Berikut daftar dan fungsinya:

### 1. JWT (`jwt.go` + `laundry_auth.go`)

Dua secret key dan jenis token:

| Fungsi | Secret | Expiry | Return |
|--------|--------|--------|--------|
| `AccessToken(username, id_user)` | `jwtSecret` | 5 menit | `(token, 200/400)` |
| `RefreshToken(username, id_user)` | `jwtSecret2` | 8 jam | `(token, 200/400)` |

**Claims struct:**
```go
type JwtCustomClaims struct {
    Username string `json:"username"`
    Id_user  string `json:"id_user"`
    jwt.RegisteredClaims
}
```

**Verifikasi & Parsing (`jwt.go`):**

| Fungsi | Input | Output |
|--------|-------|--------|
| `JWTMiddleware(c)` | Fiber context | `c.Next()` atau 401 |
| `ParseToken(token)` | Token string | `*JwtCustomClaims` (access token) |
| `ParseRefreshToken(token)` | Token string | `*JwtCustomClaims` (refresh token) |

**Otorisasi domain laundry (`laundry_auth.go`):**

| Fungsi | Input | Output |
|--------|-------|--------|
| `LaundryJWTMiddleware(c)` | Fiber context | Simpan `c.Locals("id_user", claims.Id_user)`, atau 401 |
| `RequireOutletMember(c)` | Fiber context + param `:id_outlet` | Simpan `c.Locals("outlet_role", role)`, atau 403 |
| `RequireOutletOwner(c)` | Fiber context + param `:id_outlet` | Sama seperti di atas + wajib `role == "owner"`, atau 403 |

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

### 5. S3 Upload (`s3.go`)

```go
FileUploadToS3Middleware(c *fiber.Ctx, formFieldName, folderName string,
    allowedTypes []string, maxFileSize int64) error
```

- Upload file ke S3 IDCloudHost
- Validasi tipe file (MIME) dan ukuran maksimum
- Generate nama file unik: `{folderName}/{uuid}_{timestamp}{ext}`
- Set ACL Public Read
- Simpan nama file di `c.Locals("fileName")`

### 6. Groupby (`groupby.go`)

```go
Groupby(key string, data []map[string]interface{}) → map[string][]map[string]interface{}
```
Mengelompokkan array of map berdasarkan nilai dari key tertentu. Helper generik, dipakai bila suatu domain butuh mengelompokkan hasil query per kolom tertentu.

### 7. Middleware Dasar (`middleware.go`)

```go
Logger(c *fiber.Ctx) error  // set request ID (tidak terpakai aktif)
```

### 8. Upload lokal (`upload.go`)

```go
FileUploadMiddleware(c *fiber.Ctx, uploadPath, nameFormat string, allowedTypes []string, maxSize int64) error
```
Alternatif upload ke disk lokal (`asset/public/...`) — dipakai bila tidak memakai S3.

---

## Konfigurasi (`config/config.go`)

```go
type Config struct {
    DBHost, DBPort, DBUser, DBPassword, DBName, SSLMode string
}
```

`LoadConfig()` — satu koneksi database utama (default lokal: `localhost:5432`, db `laundry_service`, user/password `postgres`).
Menggunakan environment variable dengan fallback default value via `getEnv(key, defaultValue)` — override lewat `.env`/env var saat deploy.

---

## Database (`database/database.go`)

### Variabel Global

```go
var DB *gorm.DB
var SqlDB *sql.DB
```

### Inisialisasi

`Connect()` — koneksi tunggal ke Postgres (schema default `public.`, domain laundry memakai schema `laundry.` secara eksplisit di tiap query raw SQL — lihat `database/migration/002_create_laundry_schema.sql`).
Pool: MaxIdle 10, MaxOpen 100, MaxLifetime 1 jam.

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
2. database.Connect()  → DB, SqlDB
3. Inisialisasi Fiber app (CORS)
4. Dependency Injection:
   repository := repositories.NewRepositoryRegistry(DB)
   service    := services.NewServiceRegistry(repository)
   controller := controllers.NewControllerRegistry(service)
   routes.NewRouteRegistry(controller, app).Serve()
5. app.Listen(":3001")    // HTTP
```

---

## File Konfigurasi Lain

| File | Fungsi |
|------|--------|
| `app/.env` | Env variables (DB, AWS, APP_ENV) |
| `app/.env_dev` | Template env development |
| `app/.env_prod` | Template env production |
| `Dockerfile` | Build container (multi-stage) |
| `laundry-service.dockerfile` | Alternate dockerfile |
| `docker-compose.yml` | Local dev setup |
| `jenkinsfile` | CI/CD pipeline |
| `makefile` | Shortcut commands |
| `go.mod` | Dependencies (GORM, Fiber, JWT, bcrypt, S3, validator) |

---

## Diagram Infrastruktur

```
                           ┌──────────────┐
                           │   Nginx /    │
                           │  API Gateway │
                           └──────┬───────┘
                                  │
                             Port 3001
                              (Fiber)
                                  │
                          ┌───────┴───────┐
                          │  Middlewares   │
                          │  (JWT, Valid,  │
                          │   S3, etc.)    │
                          └───────┬───────┘
                                  │
                          ┌───────┴───────┐
                          │ Domain Laundry │
                          │ (Clean Arch)   │
                          └───────┬───────┘
                                  │
                            ┌─────┴─────┐
                            │ PostgreSQL │
                            │  schema    │
                            │  "laundry" │
                            └─────┬─────┘
                                  │
                          ┌───────┴───────┐
                          │  S3 (IDCloud) │
                          │  Logo & foto  │
                          └───────────────┘
```
