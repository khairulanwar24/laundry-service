package dto

// PaymentMethodForm dipakai untuk create maupun update payment method
// (semua field wajib diisi, mengikuti perilaku Laravel: PUT mengirim payload penuh).
type PaymentMethodForm struct {
	Kategori    string   `json:"kategori" form:"kategori" validate:"required,oneof=cash transfer e_wallet"`
	Nama        string   `json:"nama" form:"nama" validate:"required,max=120"`
	Logo        string   `json:"logo" form:"logo" validate:"omitempty,max=255"`
	NamaPemilik string   `json:"nama_pemilik" form:"nama_pemilik" validate:"omitempty,max=120"`
	Tags        []string `json:"tags" form:"tags" validate:"omitempty,dive,max=30"`
	IsActive    *bool    `json:"is_active" form:"is_active" validate:"omitempty"`
}
