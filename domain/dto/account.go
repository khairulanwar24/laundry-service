package dto

// RegisterForm adalah form pendaftaran akun laundry (pemilik/karyawan outlet).
type RegisterForm struct {
	Nama               string `json:"nama" form:"nama" validate:"required,max=100"`
	Email              string `json:"email" form:"email" validate:"required,email,max=255"`
	Telepon            string `json:"telepon" form:"telepon" validate:"omitempty,max=20"`
	Password           string `json:"password" form:"password" validate:"required,min=8"`
	PasswordKonfirmasi string `json:"password_confirmation" form:"password_confirmation" validate:"required,eqfield=Password"`
	Alamat             string `json:"alamat" form:"alamat" validate:"omitempty,max=500"`
}

// LoginForm adalah form login akun laundry.
type LoginForm struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required"`
}

// ForgotPasswordForm meminta OTP reset password dikirim ke email.
type ForgotPasswordForm struct {
	Email string `json:"email" form:"email" validate:"required,email"`
}

// ResetPasswordForm mengganti password memakai token + OTP dari ForgotPassword.
type ResetPasswordForm struct {
	Email              string `json:"email" form:"email" validate:"required,email"`
	Token              string `json:"token" form:"token" validate:"required"`
	Otp                string `json:"otp" form:"otp" validate:"required,min=4,max=6"`
	Password           string `json:"password" form:"password" validate:"required,min=8"`
	PasswordKonfirmasi string `json:"password_confirmation" form:"password_confirmation" validate:"required,eqfield=Password"`
}
