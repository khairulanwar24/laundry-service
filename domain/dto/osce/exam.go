package osce

import "time"

// ===== Tahun Akademik =====

type TahunAkademikForm struct {
	Nama      string `json:"nama_tahun_akademik" form:"nama_tahun_akademik" validate:"required"`
	TglMulai  string `json:"tgl_mulai" form:"tgl_mulai" validate:"required,datetime=2006-01-02"`
	TglSelesai string `json:"tgl_selesai" form:"tgl_selesai" validate:"required,datetime=2006-01-02"`
	StatusAktif *bool `json:"status_aktif" form:"status_aktif"`
}

type TahunAkademikParams struct {
	ID string `json:"id" form:"id" validate:"required,uuid4"`
}

type TahunAkademikResponse struct {
	ID        string     `json:"id_tahun_akademik"`
	Nama      string     `json:"nama_tahun_akademik"`
	TglMulai  string     `json:"tgl_mulai"`
	TglSelesai string    `json:"tgl_selesai"`
	StatusAktif bool     `json:"status_aktif"`
	TglInsert time.Time  `json:"tgl_insert"`
	TglUpdate *time.Time `json:"tgl_update"`
}

// ===== Semester =====

type SemesterForm struct {
	IDTahunAkademik string `json:"id_tahun_akademik" form:"id_tahun_akademik" validate:"required,uuid4"`
	NamaSemester    string `json:"nama_semester" form:"nama_semester" validate:"required"`
	TglMulai        string `json:"tgl_mulai" form:"tgl_mulai" validate:"required,datetime=2006-01-02"`
	TglSelesai      string `json:"tgl_selesai" form:"tgl_selesai" validate:"required,datetime=2006-01-02"`
	StatusAktif     *bool  `json:"status_aktif" form:"status_aktif"`
}

type SemesterParams struct {
	ID string `json:"id" form:"id" validate:"required,uuid4"`
}

type SemesterResponse struct {
	ID               string     `json:"id_semester"`
	IDTahunAkademik  string     `json:"id_tahun_akademik"`
	NamaTahunAkademik string    `json:"nama_tahun_akademik"`
	NamaSemester     string     `json:"nama_semester"`
	TglMulai         string     `json:"tgl_mulai"`
	TglSelesai       string     `json:"tgl_selesai"`
	StatusAktif      bool       `json:"status_aktif"`
	TglInsert        time.Time  `json:"tgl_insert"`
	TglUpdate        *time.Time `json:"tgl_update"`
}

// ===== Program Studi =====

type ProgramStudiForm struct {
	IDProdi   string `json:"id_prodi" form:"id_prodi" validate:"required"`
	KodeProdi string `json:"kode_prodi" form:"kode_prodi" validate:"required"`
	NamaProdi string `json:"nama_prodi" form:"nama_prodi" validate:"required"`
	Jenjang   string `json:"jenjang" form:"jenjang"`
}

type ProgramStudiParams struct {
	ID string `json:"id" form:"id" validate:"required,uuid4"`
}

type ProgramStudiResponse struct {
	ID        string     `json:"id_program_studi"`
	IDProdi   string     `json:"id_prodi"`
	KodeProdi string     `json:"kode_prodi"`
	NamaProdi string     `json:"nama_prodi"`
	Jenjang   string     `json:"jenjang"`
	TglInsert time.Time  `json:"tgl_insert"`
	TglUpdate *time.Time `json:"tgl_update"`
}

// ===== Ujian (Blueprint Ujian) =====

type UjianForm struct {
	IDTahunAkademik  string   `json:"id_tahun_akademik" form:"id_tahun_akademik" validate:"required,uuid4"`
	IDSemester       string   `json:"id_semester" form:"id_semester" validate:"required,uuid4"`
	IDProgramStudi   string   `json:"id_program_studi" form:"id_program_studi" validate:"required,uuid4"`
	NamaUjian        string   `json:"nama_ujian" form:"nama_ujian" validate:"required"`
	Deskripsi        string   `json:"deskripsi" form:"deskripsi"`
	TglMulai         string   `json:"tgl_mulai" form:"tgl_mulai" validate:"required,datetime=2006-01-02"`
	TglSelesai       string   `json:"tgl_selesai" form:"tgl_selesai" validate:"required,datetime=2006-01-02"`
	DurasiPerStation int      `json:"durasi_per_station" form:"durasi_per_station"`
	BatasLulus       *float64 `json:"batas_lulus" form:"batas_lulus"`
}

type UjianParams struct {
	ID string `json:"id" form:"id" validate:"required,uuid4"`
}

type UjianQuery struct {
	Limit  int    `json:"limit" form:"limit" validate:"required,numeric,oneof=10 25 50 100" default:"10"`
	Offset int    `json:"offset" form:"offset" validate:"numeric"`
	Order  string `json:"order" form:"order"`
	Filter string `json:"filter" form:"filter"`
}

type UjianResponse struct {
	ID               string     `json:"id_ujian"`
	IDTahunAkademik  string     `json:"id_tahun_akademik"`
	NamaTahunAkademik string    `json:"nama_tahun_akademik"`
	IDSemester       string     `json:"id_semester"`
	NamaSemester     string     `json:"nama_semester"`
	IDProgramStudi   string     `json:"id_program_studi"`
	NamaProdi        string     `json:"nama_prodi"`
	NamaUjian        string     `json:"nama_ujian"`
	Deskripsi        string     `json:"deskripsi"`
	TglMulai         string     `json:"tgl_mulai"`
	TglSelesai       string     `json:"tgl_selesai"`
	DurasiPerStation int        `json:"durasi_per_station"`
	JumlahStation    int        `json:"jumlah_station"`
	BatasLulus       *float64   `json:"batas_lulus"`
	StatusUjian      string     `json:"status_ujian"`
	TglInsert        time.Time  `json:"tgl_insert"`
	TglUpdate        *time.Time `json:"tgl_update"`
}

// ===== Ujian Station =====

type UjianStationForm struct {
	IDUjian   string  `json:"id_ujian" form:"id_ujian" validate:"required,uuid4"`
	IDStation string  `json:"id_station" form:"id_station" validate:"required,uuid4"`
	Urutan    int     `json:"urutan" form:"urutan" validate:"required,min=1"`
	Durasi    int     `json:"durasi" form:"durasi"`
	Bobot     float64 `json:"bobot" form:"bobot" validate:"gt=0"`
}

type UjianStationParams struct {
	ID string `json:"id" form:"id" validate:"required,uuid4"`
}

type UjianStationResponse struct {
	ID          string     `json:"id_ujian_station"`
	IDUjian     string     `json:"id_ujian"`
	IDStation   string     `json:"id_station"`
	KodeStation string     `json:"kode_station"`
	NamaStation string     `json:"nama_station"`
	NamaTipe    string     `json:"nama_tipe"`
	Urutan      int        `json:"urutan"`
	Durasi      int        `json:"durasi"`
	Bobot       float64    `json:"bobot"`
	TglInsert   time.Time  `json:"tgl_insert"`
	TglUpdate   *time.Time `json:"tgl_update"`
}

// ===== Sesi Ujian =====

type SesiUjianForm struct {
	IDUjian       string `json:"id_ujian" form:"id_ujian" validate:"required,uuid4"`
	NamaSesi      string `json:"nama_sesi" form:"nama_sesi" validate:"required"`
	TglSesi       string `json:"tgl_sesi" form:"tgl_sesi" validate:"required,datetime=2006-01-02"`
	JamMulai      string `json:"jam_mulai" form:"jam_mulai" validate:"required,datetime=15:04"`
	JamSelesai    string `json:"jam_selesai" form:"jam_selesai" validate:"required,datetime=15:04"`
	Lokasi        string `json:"lokasi" form:"lokasi"`
	KuotaPeserta  int    `json:"kuota_peserta" form:"kuota_peserta"`
}

type SesiUjianParams struct {
	ID string `json:"id" form:"id" validate:"required,uuid4"`
}

type SesiUjianResponse struct {
	ID           string     `json:"id_sesi_ujian"`
	IDUjian      string     `json:"id_ujian"`
	NamaSesi     string     `json:"nama_sesi"`
	TglSesi      string     `json:"tgl_sesi"`
	JamMulai     string     `json:"jam_mulai"`
	JamSelesai   string     `json:"jam_selesai"`
	Lokasi       string     `json:"lokasi"`
	KuotaPeserta int        `json:"kuota_peserta"`
	StatusSesi   string     `json:"status_sesi"`
	TglInsert    time.Time  `json:"tgl_insert"`
	TglUpdate    *time.Time `json:"tgl_update"`
}

// ===== Rotasi Ujian =====

type RotasiUjianForm struct {
	IDUjian          string `json:"id_ujian" form:"id_ujian" validate:"required,uuid4"`
	IDSesiUjian      string `json:"id_sesi_ujian" form:"id_sesi_ujian" validate:"required,uuid4"`
	IDUjianStation   string `json:"id_ujian_station" form:"id_ujian_station" validate:"required,uuid4"`
	UrutanRotasi     int    `json:"urutan_rotasi" form:"urutan_rotasi" validate:"required,min=1"`
	Durasi           int    `json:"durasi" form:"durasi" validate:"required,min=1"`
	JedaAntarStation int    `json:"jeda_antar_station" form:"jeda_antar_station"`
}

type RotasiUjianParams struct {
	ID string `json:"id" form:"id" validate:"required,uuid4"`
}

// ===== Jadwal Ujian =====

type JadwalUjianForm struct {
	IDUjian        string `json:"id_ujian" form:"id_ujian" validate:"required,uuid4"`
	IDSesiUjian    string `json:"id_sesi_ujian" form:"id_sesi_ujian" validate:"required,uuid4"`
	IDUjianStation string `json:"id_ujian_station" form:"id_ujian_station" validate:"required,uuid4"`
	IDUser         string `json:"id_user" form:"id_user" validate:"required"`
	IDPenguji      string `json:"id_penguji" form:"id_penguji" validate:"uuid4"`
	UrutanMasuk    int    `json:"urutan_masuk" form:"urutan_masuk"`
	JamMulai       string `json:"jam_mulai" form:"jam_mulai" validate:"datetime=15:04"`
	JamSelesai     string `json:"jam_selesai" form:"jam_selesai" validate:"datetime=15:04"`
}

type JadwalUjianParams struct {
	ID string `json:"id" form:"id" validate:"required,uuid4"`
}

type GenerateJadwalRequest struct {
	IDUjian         string `json:"id_ujian" form:"id_ujian" validate:"required,uuid4"`
	IDSesiUjian     string `json:"id_sesi_ujian" form:"id_sesi_ujian" validate:"required,uuid4"`
	IDKelompokMahasiswa string `json:"id_kelompok_mahasiswa" form:"id_kelompok_mahasiswa" validate:"required,uuid4"`
}

type JadwalUjianResponse struct {
	ID              string     `json:"id_jadwal_ujian"`
	IDUjian         string     `json:"id_ujian"`
	IDSesiUjian     string     `json:"id_sesi_ujian"`
	NamaSesi        string     `json:"nama_sesi"`
	IDUjianStation  string     `json:"id_ujian_station"`
	KodeStation     string     `json:"kode_station"`
	NamaStation     string     `json:"nama_station"`
	IDUser          string     `json:"id_user"`
	NamaLengkap     string     `json:"nama_lengkap"`
	IDPenguji       string     `json:"id_penguji"`
	NamaPenguji     string     `json:"nama_penguji"`
	UrutanMasuk     int        `json:"urutan_masuk"`
	JamMulai        string     `json:"jam_mulai"`
	JamSelesai      string     `json:"jam_selesai"`
	StatusKehadiran string     `json:"status_kehadiran"`
	TglInsert       time.Time  `json:"tgl_insert"`
	TglUpdate       *time.Time `json:"tgl_update"`
}
