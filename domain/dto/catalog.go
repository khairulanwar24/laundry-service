package dto

// ServiceForm dipakai untuk create maupun update layanan (service) outlet.
type ServiceForm struct {
	Nama          string   `json:"nama" form:"nama" validate:"required,max=120"`
	Prioritas     *int     `json:"prioritas" form:"prioritas" validate:"omitempty,min=0,max=100"`
	LangkahProses []string `json:"langkah_proses" form:"langkah_proses" validate:"omitempty,dive,oneof=cuci kering setrika"`
	IsActive      *bool    `json:"is_active" form:"is_active" validate:"omitempty"`
}

// ServiceVariantForm dipakai untuk create maupun update varian layanan.
type ServiceVariantForm struct {
	Nama                string  `json:"nama" form:"nama" validate:"required,max=120"`
	Satuan              string  `json:"satuan" form:"satuan" validate:"required,oneof=kg pcs meter"`
	HargaPerSatuan      float64 `json:"harga_per_satuan" form:"harga_per_satuan" validate:"required,min=0"`
	DurasiPengerjaanJam int     `json:"durasi_pengerjaan_jam" form:"durasi_pengerjaan_jam" validate:"required,min=1"`
	GambarPath          string  `json:"gambar_path" form:"gambar_path" validate:"omitempty,url"`
	Catatan             string  `json:"catatan" form:"catatan" validate:"omitempty"`
	IsActive            *bool   `json:"is_active" form:"is_active" validate:"omitempty"`
}

// PerfumeForm dipakai untuk create maupun update parfum outlet.
type PerfumeForm struct {
	Nama     string `json:"nama" form:"nama" validate:"required,max=120"`
	Catatan  string `json:"catatan" form:"catatan" validate:"omitempty"`
	IsActive *bool  `json:"is_active" form:"is_active" validate:"omitempty"`
}

// DiscountForm dipakai untuk create maupun update diskon outlet.
type DiscountForm struct {
	Nama     string  `json:"nama" form:"nama" validate:"required,max=120"`
	Jenis    string  `json:"jenis" form:"jenis" validate:"required,oneof=nominal percent"`
	Nilai    float64 `json:"nilai" form:"nilai" validate:"required,min=0"`
	Catatan  string  `json:"catatan" form:"catatan" validate:"omitempty"`
	IsActive *bool   `json:"is_active" form:"is_active" validate:"omitempty"`
}
