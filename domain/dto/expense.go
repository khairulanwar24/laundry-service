package dto

// ExpenseForm dipakai untuk create maupun update data pengeluaran outlet.
type ExpenseForm struct {
	Kategori  string  `json:"kategori" form:"kategori" validate:"required,max=100"`
	Jumlah    float64 `json:"jumlah" form:"jumlah" validate:"required,min=0"`
	Deskripsi string  `json:"deskripsi" form:"deskripsi" validate:"required"`
	Tanggal   string  `json:"tanggal" form:"tanggal" validate:"required,datetime=2006-01-02"`
}
