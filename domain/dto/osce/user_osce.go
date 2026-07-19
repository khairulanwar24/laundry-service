package osce

import "time"

// ===== Penguji =====

type PengujiForm struct {
	IDUser        string `json:"id_user" form:"id_user" validate:"required"`
	GelarDepan    string `json:"gelar_depan" form:"gelar_depan"`
	GelarBelakang string `json:"gelar_belakang" form:"gelar_belakang"`
	Keahlian      string `json:"keahlian" form:"keahlian"`
	NoSTR         string `json:"no_str" form:"no_str"`
	StatusAktif   *bool  `json:"status_aktif" form:"status_aktif"`
}

type PengujiParams struct {
	ID string `json:"id" form:"id" validate:"required,uuid4"`
}

type PengujiQuery struct {
	Limit  int    `json:"limit" form:"limit" validate:"required,numeric,oneof=10 25 50 100" default:"10"`
	Offset int    `json:"offset" form:"offset" validate:"numeric"`
	Order  string `json:"order" form:"order"`
	Filter string `json:"filter" form:"filter"`
}

type PengujiResponse struct {
	ID             string     `json:"id_penguji"`
	IDUser         string     `json:"id_user"`
	Username       string     `json:"username"`
	NamaLengkap    string     `json:"nama_lengkap"`
	Email          string     `json:"email"`
	GelarDepan     string     `json:"gelar_depan"`
	GelarBelakang  string     `json:"gelar_belakang"`
	Keahlian       string     `json:"keahlian"`
	NoSTR          string     `json:"no_str"`
	StatusAktif    bool       `json:"status_aktif"`
	TglInsert      time.Time  `json:"tgl_insert"`
	TglUpdate      *time.Time `json:"tgl_update"`
}

// ===== Kelompok Mahasiswa =====

type KelompokMahasiswaForm struct {
	IDUjian      string `json:"id_ujian" form:"id_ujian" validate:"uuid4"`
	NamaKelompok string `json:"nama_kelompok" form:"nama_kelompok" validate:"required"`
	Deskripsi    string `json:"deskripsi" form:"deskripsi"`
}

type KelompokMahasiswaParams struct {
	ID string `json:"id" form:"id" validate:"required,uuid4"`
}

type KelompokMahasiswaQuery struct {
	Limit  int    `json:"limit" form:"limit" validate:"required,numeric,oneof=10 25 50 100" default:"10"`
	Offset int    `json:"offset" form:"offset" validate:"numeric"`
	Order  string `json:"order" form:"order"`
	Filter string `json:"filter" form:"filter"`
}

type KelompokMahasiswaResponse struct {
	ID            string     `json:"id_kelompok_mahasiswa"`
	IDUjian       string     `json:"id_ujian"`
	NamaUjian     string     `json:"nama_ujian"`
	NamaKelompok  string     `json:"nama_kelompok"`
	Deskripsi     string     `json:"deskripsi"`
	JumlahAnggota int        `json:"jumlah_anggota"`
	TglInsert     time.Time  `json:"tgl_insert"`
	TglUpdate     *time.Time `json:"tgl_update"`
}

// ===== Anggota Kelompok Mahasiswa =====

type AnggotaKelompokForm struct {
	IDUser string `json:"id_user" form:"id_user" validate:"required"`
}

type AnggotaKelompokParams struct {
	ID string `json:"id" form:"id" validate:"required,uuid4"`
}

type BulkAnggotaKelompokForm struct {
	IDUsers []string `json:"id_users" form:"id_users" validate:"required,min=1"`
}

type AnggotaKelompokResponse struct {
	ID                 string `json:"id_anggota_kelompok"`
	IDKelompokMahasiswa string `json:"id_kelompok_mahasiswa"`
	IDUser             string `json:"id_user"`
	Username           string `json:"username"`
	NamaLengkap        string `json:"nama_lengkap"`
	NIM                string `json:"nim"`
	Email              string `json:"email"`
}

// ===== Penugasan Penguji =====

type PenugasanPengujiForm struct {
	IDPenguji      string `json:"id_penguji" form:"id_penguji" validate:"required,uuid4"`
	IDSesiUjian    string `json:"id_sesi_ujian" form:"id_sesi_ujian" validate:"required,uuid4"`
	IDUjianStation string `json:"id_ujian_station" form:"id_ujian_station" validate:"required,uuid4"`
}

type PenugasanPengujiParams struct {
	ID string `json:"id" form:"id" validate:"required,uuid4"`
}

type PenugasanPengujiResponse struct {
	ID              string     `json:"id_penugasan_penguji"`
	IDPenguji       string     `json:"id_penguji"`
	NamaPenguji     string     `json:"nama_penguji"`
	IDSesiUjian     string     `json:"id_sesi_ujian"`
	NamaSesi        string     `json:"nama_sesi"`
	TglSesi         string     `json:"tgl_sesi"`
	IDUjianStation  string     `json:"id_ujian_station"`
	KodeStation     string     `json:"kode_station"`
	NamaStation     string     `json:"nama_station"`
	TglInsert       time.Time  `json:"tgl_insert"`
	TglUpdate       *time.Time `json:"tgl_update"`
}
