package osce

import "time"

// ===== Tipe Station =====

type TipeStationForm struct {
	NamaTipe  string `json:"nama_tipe" form:"nama_tipe" validate:"required"`
	Deskripsi string `json:"deskripsi" form:"deskripsi"`
}

type TipeStationParams struct {
	ID string `json:"id" form:"id" validate:"required,uuid4"`
}

// ===== Master Station =====

type StationForm struct {
	IDTipeStation string  `json:"id_tipe_station" form:"id_tipe_station" validate:"required,uuid4"`
	KodeStation   string  `json:"kode_station" form:"kode_station" validate:"required"`
	NamaStation   string  `json:"nama_station" form:"nama_station" validate:"required"`
	Deskripsi     string  `json:"deskripsi" form:"deskripsi"`
	DurasiDefault int     `json:"durasi_default" form:"durasi_default" validate:"required,min=1"`
	Bobot         float64 `json:"bobot" form:"bobot" validate:"required,gt=0"`
}

type StationParams struct {
	ID string `json:"id" form:"id" validate:"required,uuid4"`
}

type StationQuery struct {
	Limit  int    `json:"limit" form:"limit" validate:"required,numeric,oneof=10 25 50 100" default:"10"`
	Offset int    `json:"offset" form:"offset" validate:"numeric"`
	Order  string `json:"order" form:"order"`
	Filter string `json:"filter" form:"filter"`
}

// ===== Kompetensi Station =====

type KompetensiStationForm struct {
	IDStation       string  `json:"id_station" form:"id_station" validate:"required,uuid4"`
	NamaKompetensi  string  `json:"nama_kompetensi" form:"nama_kompetensi" validate:"required"`
	Deskripsi       string  `json:"deskripsi" form:"deskripsi"`
	Bobot           float64 `json:"bobot" form:"bobot" validate:"gt=0"`
	Urutan          int     `json:"urutan" form:"urutan"`
}

type KompetensiStationParams struct {
	ID string `json:"id" form:"id" validate:"required,uuid4"`
}

// ===== Checklist / Rubrik =====

type ChecklistStationForm struct {
	IDStation            string  `json:"id_station" form:"id_station" validate:"required,uuid4"`
	IDKompetensiStation  string  `json:"id_kompetensi_station" form:"id_kompetensi_station" validate:"uuid4"`
	Urutan               int     `json:"urutan" form:"urutan" validate:"min=1"`
	DeskripsiItem        string  `json:"deskripsi_item" form:"deskripsi_item" validate:"required"`
	SkorMaksimal         int     `json:"skor_maksimal" form:"skor_maksimal" validate:"required,min=1"`
	BobotItem            float64 `json:"bobot_item" form:"bobot_item" validate:"gt=0"`
	TipePenilaian        string  `json:"tipe_penilaian" form:"tipe_penilaian" validate:"oneof=skala biner penilaian_global catatan"`
	PetunjukPenilaian    string  `json:"petunjuk_penilaian" form:"petunjuk_penilaian"`
}

type ChecklistStationParams struct {
	ID string `json:"id" form:"id" validate:"required,uuid4"`
}

type ChecklistReorderForm struct {
	IDChecklist string `json:"id_checklist" validate:"required,uuid4"`
	Urutan      int    `json:"urutan" validate:"required,min=1"`
}

type ChecklistReorderRequest struct {
	Items []ChecklistReorderForm `json:"items" validate:"required,min=1,dive"`
}

// ===== Response struct =====

type TipeStationResponse struct {
	ID         string     `json:"id_tipe_station"`
	NamaTipe   string     `json:"nama_tipe"`
	Deskripsi  string     `json:"deskripsi"`
	TglInsert  time.Time  `json:"tgl_insert"`
	TglUpdate  *time.Time `json:"tgl_update"`
}

type StationResponse struct {
	ID             string     `json:"id_station"`
	IDTipeStation  string     `json:"id_tipe_station"`
	NamaTipe       string     `json:"nama_tipe"`
	KodeStation    string     `json:"kode_station"`
	NamaStation    string     `json:"nama_station"`
	Deskripsi      string     `json:"deskripsi"`
	DurasiDefault  int        `json:"durasi_default"`
	Bobot          float64    `json:"bobot"`
	TglInsert      time.Time  `json:"tgl_insert"`
	TglUpdate      *time.Time `json:"tgl_update"`
}

type KompetensiResponse struct {
	ID              string     `json:"id_kompetensi_station"`
	IDStation       string     `json:"id_station"`
	NamaKompetensi  string     `json:"nama_kompetensi"`
	Deskripsi       string     `json:"deskripsi"`
	Bobot           float64    `json:"bobot"`
	Urutan          int        `json:"urutan"`
	TglInsert       time.Time  `json:"tgl_insert"`
	TglUpdate       *time.Time `json:"tgl_update"`
}

type ChecklistResponse struct {
	ID                   string     `json:"id_checklist_station"`
	IDStation            string     `json:"id_station"`
	IDKompetensiStation  string     `json:"id_kompetensi_station"`
	NamaKompetensi       string     `json:"nama_kompetensi"`
	Urutan               int        `json:"urutan"`
	DeskripsiItem        string     `json:"deskripsi_item"`
	SkorMaksimal         int        `json:"skor_maksimal"`
	BobotItem            float64    `json:"bobot_item"`
	TipePenilaian        string     `json:"tipe_penilaian"`
	PetunjukPenilaian    string     `json:"petunjuk_penilaian"`
	TglInsert            time.Time  `json:"tgl_insert"`
	TglUpdate            *time.Time `json:"tgl_update"`
}
