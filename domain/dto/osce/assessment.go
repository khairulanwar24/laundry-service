package osce

import "time"

// ===== Penilaian =====

type PenilaianForm struct {
	IDJadwalUjian   string                  `json:"id_jadwal_ujian" form:"id_jadwal_ujian" validate:"required,uuid4"`
	IDPenguji       string                  `json:"id_penguji" form:"id_penguji" validate:"required,uuid4"`
	PenilaianGlobal string                  `json:"penilaian_global" form:"penilaian_global" validate:"oneof=borderline lulus tidak_lulus istimewa"`
	CatatanPenguji  string                  `json:"catatan_penguji" form:"catatan_penguji"`
	Details         []DetailPenilaianForm   `json:"details" form:"details" validate:"required,min=1,dive"`
}

type DetailPenilaianForm struct {
	IDChecklistStation string `json:"id_checklist_station" form:"id_checklist_station" validate:"required,uuid4"`
	Skor               int    `json:"skor" form:"skor" validate:"min=0"`
	CatatanItem        string `json:"catatan_item" form:"catatan_item"`
}

type PenilaianParams struct {
	ID string `json:"id" form:"id" validate:"required,uuid4"`
}

type PenilaianUpdateForm struct {
	IDJadwalUjian   string                  `json:"id_jadwal_ujian" form:"id_jadwal_ujian" validate:"required,uuid4"`
	IDPenguji       string                  `json:"id_penguji" form:"id_penguji" validate:"required,uuid4"`
	PenilaianGlobal string                  `json:"penilaian_global" form:"penilaian_global" validate:"oneof=borderline lulus tidak_lulus istimewa"`
	CatatanPenguji  string                  `json:"catatan_penguji" form:"catatan_penguji"`
	Details         []DetailPenilaianForm   `json:"details" form:"details" validate:"required,min=1,dive"`
}

// ===== Penilaian Response =====

type PenilaianResponse struct {
	ID              string                    `json:"id_penilaian"`
	IDJadwalUjian   string                    `json:"id_jadwal_ujian"`
	IDPenguji       string                    `json:"id_penguji"`
	NamaPenguji     string                    `json:"nama_penguji"`
	NamaLengkap     string                    `json:"nama_lengkap"`
	NIM             string                    `json:"nim"`
	KodeStation     string                    `json:"kode_station"`
	NamaStation     string                    `json:"nama_station"`
	TotalSkor       float64                   `json:"total_skor"`
	SkorTerbobot    float64                   `json:"skor_terbobot"`
	PenilaianGlobal string                    `json:"penilaian_global"`
	CatatanPenguji  string                    `json:"catatan_penguji"`
	StatusPenilaian string                    `json:"status_penilaian"`
	TglPenilaian    *time.Time                `json:"tgl_penilaian"`
	Details         []DetailPenilaianResponse `json:"details"`
	TglInsert       time.Time                 `json:"tgl_insert"`
	TglUpdate       *time.Time                `json:"tgl_update"`
}

type DetailPenilaianResponse struct {
	ID                  string  `json:"id_detail_penilaian"`
	IDChecklistStation  string  `json:"id_checklist_station"`
	DeskripsiItem       string  `json:"deskripsi_item"`
	SkorMaksimal        int     `json:"skor_maksimal"`
	BobotItem           float64 `json:"bobot_item"`
	Skor                int     `json:"skor"`
	CatatanItem         string  `json:"catatan_item"`
}

// ===== Hasil Ujian =====

type HasilUjianResponse struct {
	ID              string                     `json:"id_hasil_ujian"`
	IDUjian         string                     `json:"id_ujian"`
	NamaUjian       string                     `json:"nama_ujian"`
	IDUser          string                     `json:"id_user"`
	NamaLengkap     string                     `json:"nama_lengkap"`
	NIM             string                     `json:"nim"`
	TotalSkor       float64                    `json:"total_skor"`
	PersentaseSkor  float64                    `json:"persentase_skor"`
	BatasLulus      *float64                   `json:"batas_lulus"`
	StatusKelulusan string                     `json:"status_kelulusan"`
	CatatanHasil    string                     `json:"catatan_hasil"`
	Stations        []HasilUjianStationResponse `json:"stations"`
	TglInsert       time.Time                  `json:"tgl_insert"`
	TglUpdate       *time.Time                 `json:"tgl_update"`
}

type HasilUjianStationResponse struct {
	ID              string  `json:"id_hasil_ujian_station"`
	IDUjianStation  string  `json:"id_ujian_station"`
	KodeStation     string  `json:"kode_station"`
	NamaStation     string  `json:"nama_station"`
	SkorStation     float64 `json:"skor_station"`
	SkorTerbobot    float64 `json:"skor_terbobot"`
	PenilaianGlobal string  `json:"penilaian_global"`
}

// ===== Borderline Regression =====

type BorderlineResultItem struct {
	IDUjianStation string  `json:"id_ujian_station"`
	KodeStation    string  `json:"kode_station"`
	NamaStation    string  `json:"nama_station"`
	CutOffScore    float64 `json:"cut_off_score"`
	R2             float64 `json:"r2"`
	Slope          float64 `json:"slope"`
	Intercept      float64 `json:"intercept"`
}

type BorderlineResultResponse struct {
	IDUjian    string                `json:"id_ujian"`
	NamaUjian  string                `json:"nama_ujian"`
	Results    []BorderlineResultItem `json:"results"`
}

// ===== Statistik =====

type StatistikUjianResponse struct {
	IDUjian          string  `json:"id_ujian"`
	NamaUjian        string  `json:"nama_ujian"`
	TotalPeserta     int     `json:"total_peserta"`
	JumlahLulus      int     `json:"jumlah_lulus"`
	JumlahTidakLulus int     `json:"jumlah_tidak_lulus"`
	JumlahRemedi     int     `json:"jumlah_remedi"`
	PersentaseLulus  float64 `json:"persentase_lulus"`
	RataRataSkor     float64 `json:"rata_rata_skor"`
	SkorTertinggi    float64 `json:"skor_tertinggi"`
	SkorTerendah     float64 `json:"skor_terendah"`
}

// ===== Kalkulasi Request =====

type KalkulasiHasilRequest struct {
	IDUjian     string   `json:"id_ujian" form:"id_ujian" validate:"required,uuid4"`
	BatasLulus  *float64 `json:"batas_lulus" form:"batas_lulus"`
}

type BorderlineRequest struct {
	IDUjian string `json:"id_ujian" form:"id_ujian" validate:"required,uuid4"`
}
