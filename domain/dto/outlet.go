package dto

// CreateOutletForm adalah form membuat outlet baru (pembuat otomatis jadi owner).
type CreateOutletForm struct {
	Nama     string `json:"nama" form:"nama" validate:"required,max=120"`
	Alamat   string `json:"alamat" form:"alamat" validate:"omitempty,max=255"`
	Telepon  string `json:"telepon" form:"telepon" validate:"omitempty,max=30"`
	LogoPath string `json:"logo_path" form:"logo_path" validate:"omitempty,url"`
}

// UpdateOutletForm adalah form memperbarui data outlet (owner only).
type UpdateOutletForm struct {
	Nama     string `json:"nama" form:"nama" validate:"required,max=120"`
	Alamat   string `json:"alamat" form:"alamat" validate:"omitempty,max=255"`
	Telepon  string `json:"telepon" form:"telepon" validate:"omitempty,max=30"`
	LogoPath string `json:"logo_path" form:"logo_path" validate:"omitempty,url"`
}

// InviteEmployeeForm mengundang/menambahkan karyawan ke outlet.
type InviteEmployeeForm struct {
	Nama               string          `json:"nama" form:"nama" validate:"required,max=120"`
	Telepon            string          `json:"telepon" form:"telepon" validate:"required,max=20"`
	Email              string          `json:"email" form:"email" validate:"required,email,max=255"`
	Password           string          `json:"password" form:"password" validate:"required,min=8"`
	PasswordKonfirmasi string          `json:"password_confirmation" form:"password_confirmation" validate:"required,eqfield=Password"`
	Role               string          `json:"role" form:"role" validate:"required,oneof=owner karyawan"`
	Alamat             string          `json:"alamat" form:"alamat" validate:"required,max=500"`
	Permissions        map[string]bool `json:"permissions" form:"permissions" validate:"omitempty"`
}

// UpdateMemberForm memperbarui role/permission anggota outlet (owner only).
// Role kosong berarti role tidak diubah — hanya permissions yang di-merge.
type UpdateMemberForm struct {
	Role        string          `json:"role" form:"role" validate:"omitempty,oneof=owner karyawan"`
	Permissions map[string]bool `json:"permissions" form:"permissions" validate:"omitempty"`
}
