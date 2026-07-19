# 01 — Auth (Autentikasi)

## Deskripsi

Domain autentikasi menangani: login user, logout, reset password (via OTP email), verifikasi OTP, ganti password, dan refresh token JWT.

Menggunakan **tiga jenis JWT** dengan secret key berbeda:
- **Access Token** (5 menit) — `jwtSecret` — otorisasi API
- **Refresh Token** (8 jam) — `jwtSecret2` — perpanjang access token (HTTP-only cookie)
- **OTP Token** (10 menit) — `jwtSecret3` — verifikasi OTP reset password (HTTP-only cookie)

> ⚠️ **Security by obscurity**: Saat reset password, jika username tidak valid, sistem tetap mengembalikan `Success: true` agar attacker tidak bisa membedakan user valid vs tidak.

---

## Endpoint

| Method | Path | Middleware | Handler |
|--------|------|-----------|---------|
| `POST` | `/auth/login` | `ValidateForm(LoginForm)` | `LoginController` |
| `POST` | `/auth/logout` | - | `LogoutController` |
| `POST` | `/auth/resetpassword` | `ValidateForm(ResetForm)` | `ResetPassword` |
| `POST` | `/auth/cekotp` | `ValidateForm(OtpForm)` | `CekOtp` |
| `POST` | `/auth/changepassword` | `ValidateForm(ChangePasswordForm)` | `ChangePassword` |
| `GET` | `/auth/gettoken` | - | `RefreshToken` |

---

## DTO

### Request

```go
// POST /auth/login
type LoginForm struct {
    Username string `json:"username" validate:"required"`
    Password string `json:"password" validate:"required"`
}

// POST /auth/resetpassword
type ResetForm struct {
    Username string `json:"username" validate:"required"`
}

// POST /auth/cekotp
type OtpForm struct {
    Otp string `json:"otp" validate:"required,numeric"`
}

// POST /auth/changepassword
type ChangePasswordForm struct {
    Username string `json:"username"`
    Password string `json:"password" validate:"required"`
}
```

### Internal

```go
type UserBlock struct {
    StillBlocked     bool    `json:"still_blocked"`
    SecondsRemaining float64 `json:"seconds_remaining"`
    IdUser           string  `json:"id_user"`
    Email            string  `json:"email"`
}
```

---

## Alur Fungsi

### 1. Login (`POST /auth/login`)

```
Controller: LoginController()
├── Ambil LoginForm dari c.Locals("validatedForm")
├── Panggil Service: service.GetAuth().Login(ctx, username, password)
│   ├── Repository: FindUserByUsername(ctx, username)
│   │   └── SQL: SELECT * FROM users WHERE username = ?
│   ├── Bandingkan: middleware.CheckPasswordHash(password, hash)
│   └── Return: Response{Success, Data: user (tanpa password)}
│
├── Jika sukses:
│   ├── Buat refresh_token: middleware.RefreshToken(username, id_user)
│   │   └── JWT 8 jam, secret=jwtSecret2
│   ├── Set HTTP-only cookie "refresh_token"
│   ├── Buat access_token: middleware.AccessToken(username, id_user)
│   │   └── JWT 5 menit, secret=jwtSecret
│   └── Response: {user, access_token}
│
└── Jika gagal: Response{Success: false, Message: "username or password is incorrect"}
```

### 2. Reset Password (`POST /auth/resetpassword`)

```
Controller: ResetPassword()
├── Ambil ResetForm, generate UUID baru
├── Panggil Service: service.GetAuth().ResetPassword(ctx, username, uuid)
│   ├── Repository: FindUserBlockByUsername(ctx, username)
│   │   └── SQL: SELECT id_user, email, tgl_lock > NOW() AS still_blocked,
│   │            EXTRACT(EPOCH FROM (tgl_lock - NOW())) AS seconds_remaining
│   │            FROM users WHERE username = ?
│   │
│   ├── Generate OTP 6 digit random
│   │
│   ├── JIKA user TIDAK valid || id_user kosong:
│   │   └── InsertPasswordResetNoUser (fallback, sembunyikan error)
│   │
│   ├── JIKA user TERBLOKIR:
│   │   └── Return: "Akun Anda terkunci selama X jam Y menit lagi"
│   │
│   ├── DeactivatePasswordResetByUser(idUser)
│   │   └── SQL: UPDATE password_reset SET status_data=false WHERE id_user=?
│   │
│   ├── InsertPasswordResetWithUser(uuid, username, idUser, otp, now)
│   │   └── SQL: INSERT INTO password_reset (...)
│   │
│   ├── Kirim email OTP: middleware.Mail(to, subject, body)
│   │   └── HTTP POST ke api.farmasiunissula.com/mail/send-email
│   │
│   └── Return: Response{Success: true}
│
├── Jika sukses:
│   ├── Buat OTP cookie: middleware.OTPResetpassword(uuid, 0)
│   │   └── JWT 10 menit, secret=jwtSecret3, menyimpan id_reset_password
│   └── Set HTTP-only cookie "OTP"
│
└── Return JSON
```

### 3. Cek OTP (`POST /auth/cekotp`)

```
Controller: CekOtp()
├── Ambil OtpForm
├── Verifikasi cookie "OTP": middleware.CheckjwtOTP(c)
├── Parse cookie "OTP": middleware.ParseTokenOTP(c)
│   └── Dapatkan IDResetPassword
│
├── Panggil Service: service.GetAuth().CekOtp(ctx, IDResetPassword, otp)
│   ├── Repository: FindPasswordReset(ctx, idPasswordReset)
│   │   └── SQL: SELECT username, id_user, tgl_insert, percobaan, otp
│   │            FROM password_reset WHERE id_password_reset = ? AND status_data=true
│   │
│   ├── Increment percobaan: UpdatePasswordResetPercobaan(ctx, percobaan+1, id)
│   │
│   ├── Jika id_user == nil → "OTP Salah"
│   │
│   ├── Jika percobaan > 5 → LockUser(ctx, idUser)
│   │   └── SQL: UPDATE users SET tgl_lock = NOW() + INTERVAL '30 minutes'
│   │   └── Return: "Limit percobaan OTP terpenuhi, coba lagi dalam 30 menit"
│   │
│   ├── Bandingkan OTP dari user vs OTP di DB
│   │
│   ├── Cek expired: now - tgl_insert > 10 menit → "OTP Expired"
│   │
│   └── Return: Response{Success: true, Message: "OTP Valid"}
│
├── Jika sukses:
│   ├── Generate Token Change Password:
│   │   service.GetAuth().GenerateTokenChangePassword(ctx, IDResetPassword)
│   │   └── Repository: FindPasswordReset → dapatkan username & idUser
│   │   └── middleware.AccessTokenOTP(username, idUser)
│   │       └── JWT 5 menit, secret=jwtSecret3
│   └── Set cookie "access_token" (nilai = token change password)
│
├── Jika gagal:
│   ├── Cek percobaan: service.GetAuth().CekPercobaan(ctx, IDResetPassword)
│   └── Update cookie "OTP" dengan percobaan baru
│
└── Return JSON
```

### 4. Change Password (`POST /auth/changepassword`)

```
Controller: ChangePassword()
├── Ambil ChangePasswordForm
├── Ambil Bearer token dari Authorization header
│
├── Panggil Service: service.GetAuth().ChangePassword(ctx, username, password, token)
│   ├── Parse token: middleware.ParseTokenChangePassword(token)
│   │   └── JWT, secret=jwtSecret3 → dapatkan idUser
│   │
│   ├── Repository: FindUserFirstLogin(ctx, idUser)
│   │   └── SQL: SELECT id_user, first_login FROM users WHERE id_user=?
│   │
│   ├── Hash password: middleware.HashPassword(password)
│   │
│   ├── JIKA first_login == true:
│   │   └── UpdateUserPasswordFirstLogin(passwordHash, idUser)
│   │       └── SQL: UPDATE users SET password=?, tgl_update=NOW(),
│   │                first_login=false, tgl_first_login=NOW() WHERE id_user=?
│   │
│   ├── JIKA first_login == false:
│   │   └── UpdateUserPasswordNormal(passwordHash, idUser)
│   │       └── SQL: UPDATE users SET password=?, tgl_update=NOW() WHERE id_user=?
│   │
│   └── Return: Response{Success: true}
│
└── Return JSON
```

### 5. Refresh Token (`GET /auth/gettoken`)

```
Controller: RefreshToken()
├── Cek cookie "refresh_token": middleware.CheckjwtRefresh(c)
│   └── JWT, secret=jwtSecret2
│
├── Jika valid: middleware.JWTParse(c)
│   ├── Parse refresh token dari cookie
│   ├── Dapatkan username & id_user dari claims
│   ├── Buat access_token baru: middleware.AccessToken(username, id_user)
│   └── Return: {access_token}
│
└── Return JSON
```

### 6. Logout (`POST /auth/logout`)

```
Controller: LogoutController()
└── Hapus cookie "refresh_token" (set expired ke masa lalu)
```

---

## Repository (Query SQL)

| Method | Tabel | Operasi |
|--------|-------|---------|
| `FindUserByUsername` | `users` | `SELECT` where username |
| `FindUserBlockByUsername` | `users` | `SELECT` tgl_lock, still_blocked |
| `InsertPasswordResetNoUser` | `password_reset` | `INSERT` tanpa id_user |
| `DeactivatePasswordResetByUser` | `password_reset` | `UPDATE status_data=false` |
| `InsertPasswordResetWithUser` | `password_reset` | `INSERT` dengan id_user |
| `FindPasswordReset` | `password_reset` | `SELECT` where id + status_data=true |
| `UpdatePasswordResetPercobaan` | `password_reset` | `UPDATE percobaan` |
| `LockUser` | `users` | `UPDATE tgl_lock = NOW() + 30 min` |
| `FindUserFirstLogin` | `users` | `SELECT first_login` |
| `UpdateUserPasswordFirstLogin` | `users` | `UPDATE password + first_login=false` |
| `UpdateUserPasswordNormal` | `users` | `UPDATE password` |

---

## Dependency

- `middleware.CheckPasswordHash` — verifikasi bcrypt
- `middleware.HashPassword` — generate bcrypt hash
- `middleware.AccessToken` — generate JWT access token
- `middleware.RefreshToken` — generate JWT refresh token
- `middleware.AccessTokenOTP` — generate JWT token ganti password
- `middleware.OTPResetpassword` — generate JWT OTP cookie
- `middleware.CheckjwtRefresh` — verifikasi refresh token dari cookie
- `middleware.CheckjwtOTP` — verifikasi OTP token dari cookie
- `middleware.ParseTokenOTP` — decode OTP token claims
- `middleware.ParseTokenChangePassword` — decode change password token claims
- `middleware.JWTParse` — decode refresh token, buat access token baru
- `middleware.Mail` — kirim email via HTTP API
- `common/response.Response` — struktur response
