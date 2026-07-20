package dto

// CustomerForm dipakai untuk create maupun update data pelanggan outlet.
type CustomerForm struct {
	Nama     string `json:"nama" form:"nama" validate:"required,max=255"`
	Telepon  string `json:"telepon" form:"telepon" validate:"omitempty,max=20"`
	Email    string `json:"email" form:"email" validate:"omitempty,email,max=255"`
	Alamat   string `json:"alamat" form:"alamat" validate:"omitempty,max=500"`
	IsActive *bool  `json:"is_active" form:"is_active" validate:"omitempty"`
}
