package dto

type LoginForm struct {
	Username string `json:"username" form:"username" validate:"required"`
	Password string `json:"password"  form:"password" validate:"required"`
}

type ResetForm struct {
	Username string `json:"username" form:"username" validate:"required"`
}

type ChangePasswordForm struct {
	Username string `json:"username" form:"username" validate:""`
	Password string `json:"password"  form:"password" validate:"required"`
}

type OtpForm struct {
	Otp string `json:"otp" form:"otp" validate:"required,numeric"`
}

// UserBlock menampung info status penguncian akun saat reset password.
type UserBlock struct {
	StillBlocked     bool    `json:"still_blocked"`
	SecondsRemaining float64 `json:"seconds_remaining"`
	IdUser           string  `json:"id_user"`
	Email            string  `json:"email"`
}
