# 08 — Account (Akun Laundry)

## Deskripsi

Domain akun untuk pengguna aplikasi laundry (pemilik/karyawan outlet) — **terpisah dari domain
`auth` SSO lama** (yang sudah dihapus) karena konsepnya berbeda: login pakai email+password,
disimpan di skema `laundry.users`, bukan `public.users` milik SSO Farmasi UNISSULA.

Mencakup: registrasi, login, lihat profil sendiri (`me`), logout, dan reset password via kode OTP
yang dikirim ke email.

JWT memakai dua jenis token yang **sama persis dengan mekanisme yang sudah ada** di
`middlewares/jwt.go` (tidak dibuat baru):
- **Access Token** (5 menit) — secret `jwtSecret` — dipakai di header `Authorization: Bearer`
- **Refresh Token** (8 jam) — secret `jwtSecret2`

> Karena JWT di project ini stateless (tidak ada tabel token di database), `Logout` tidak benar-benar
> mencabut token di server — hanya mengembalikan response sukses; client wajib membuang token yang
> tersimpan.

---

## Endpoint

| Method | Path | Middleware | Handler |
|--------|------|-----------|---------|
| `POST` | `/account/register` | `ValidateForm(RegisterForm)` | `Register` |
| `POST` | `/account/login` | `ValidateForm(LoginForm)` | `Login` |
| `POST` | `/account/forgot-password` | `ValidateForm(ForgotPasswordForm)` | `ForgotPassword` |
| `POST` | `/account/reset-password` | `ValidateForm(ResetPasswordForm)` | `ResetPassword` |
| `GET` | `/account/me` | `LaundryJWTMiddleware` | `Me` |
| `POST` | `/account/logout` | `LaundryJWTMiddleware` | `Logout` |

---

## DTO (`domain/dto/account.go`)

```go
type RegisterForm struct {
    Nama               string // required, max 100
    Email              string // required, email, max 255
    Telepon            string // opsional, max 20
    Password           string // required, min 8
    PasswordKonfirmasi string // required, harus sama dengan Password (json: password_confirmation)
    Alamat             string // opsional, max 500
}

type LoginForm struct {
    Email    string // required, email
    Password string // required
}

type ForgotPasswordForm struct {
    Email string // required, email
}

type ResetPasswordForm struct {
    Email              string
    Token              string // token plaintext dari response ForgotPassword
    Otp                string // 4-6 karakter
    Password           string // min 8
    PasswordKonfirmasi string // harus sama dengan Password
}
```

---

## Alur Fungsi

### 1. Register (`POST /account/register`)

```
Service.Register()
├── Cek email sudah dipakai? → repository.GetAccount().ExistsEmail()
├── Cek telepon sudah dipakai (jika diisi)? → ExistsPhone()
├── Hash password: middleware.HashPassword()
├── Insert ke laundry.users → CreateUser() (RETURNING id_user)
├── Buat access_token & refresh_token: middleware.AccessToken()/RefreshToken()
│   (parameter "username" JWT diisi dengan email)
└── Response: {id_user, nama, email, access_token, refresh_token}
```

### 2. Login (`POST /account/login`)

```
Service.Login()
├── FindUserByEmail() → ambil hash password & is_active
├── middleware.CheckPasswordHash(password, hash) → salah → "Email atau password salah"
├── is_active == false → "Akun tidak aktif"
├── Buat access_token & refresh_token
└── Response: {user (tanpa password), access_token, refresh_token}
```

### 3. Forgot Password (`POST /account/forgot-password`)

```
Service.ForgotPassword()
├── FindUserByEmail() — jika tidak ketemu/tidak aktif, tetap balas sukses generik
│   ("Jika email terdaftar, OTP telah dikirim") supaya tidak bocorkan email terdaftar/tidak
├── Generate OTP 6 digit (crypto/rand) & token acak 32 byte (hex)
├── Hash token: middleware.HashPassword(token) → disimpan sebagai token_hash
├── UpsertPasswordReset(): hapus baris lama utk email ini, insert baris baru
│   ke laundry.password_resets (email, otp, token_hash, otp_sent_at)
├── Kirim email OTP: middleware.Mail() (SMTP, lihat 07-infrastruktur.md)
└── Response: {token (plaintext, dikirim balik ke client), otp_expires_in: 600}
```

Client harus menyimpan `token` dari response ini dan mengirimkannya kembali di `reset-password`
bersama OTP yang diterima lewat email.

### 4. Reset Password (`POST /account/reset-password`)

```
Service.ResetPassword()
├── FindPasswordResetByEmail()
├── CheckPasswordHash(token, token_hash) → token salah → error
├── Bandingkan OTP (string exact match) → salah → error
├── Cek kedaluwarsa: now - otp_sent_at > 10 menit → error, hapus baris reset
├── FindUserByEmail() → dapatkan id_user
├── Hash password baru, UpdateUserPassword()
├── DeletePasswordResetByEmail()
└── Response sukses
```

---

## Repository (`repositories/account/account.go`)

Tabel: `laundry.users`, `laundry.password_resets` (lihat skema di
`database/migration/002_create_laundry_schema.sql`).

| Fungsi | Keterangan |
|--------|-----------|
| `FindUserByEmail` | `SELECT ... FROM laundry.users WHERE email = ? AND status_data = true` |
| `FindUserByID` | Ambil profil tanpa password, untuk endpoint `Me` |
| `ExistsEmail` / `ExistsPhone` | `COUNT(*)` untuk validasi unik saat registrasi |
| `CreateUser` | `INSERT ... RETURNING id_user` |
| `UpsertPasswordReset` | `DELETE` lalu `INSERT` (satu baris aktif per email) |
| `UpdateUserPassword` | `UPDATE laundry.users SET password = ?` |

---

## Dependency

- `middleware.AccessToken` / `RefreshToken` / `ParseToken` (`middlewares/jwt.go`)
- `middleware.HashPassword` / `CheckPasswordHash` (`middlewares/password.go`)
- `middleware.Mail` (`middlewares/mail.go`) — SMTP generik, lihat `07-infrastruktur.md`
- `middleware.LaundryJWTMiddleware` (`middlewares/laundry_auth.go`) — menyimpan `id_user` dari
  token ke `c.Locals("id_user")`, dipakai oleh hampir semua domain lain untuk mengetahui pelaku aksi
