package dto

// ==== Request DTO (form & params) ====

type GetUsersForm struct {
	Limit  int    `json:"limit" form:"limit" validate:"required,numeric,oneof=10 25 50 100"  default:"10"`
	Offset int    `json:"offset"  form:"offset" validate:"numeric"`
	Order  string `json:"order"  form:"order" validate:""`
	Filter string `json:"filter"  form:"filter" validate:""`
}

type CreateUsersForm struct {
	Email        string `json:"email" form:"email" validate:"required,email"`
	Id_Person    string `json:"id_person" form:"id_person" validate:"required,uuid4"`
	Jenis_User   string `json:"jenis_user" form:"jenis_user" validate:"required,oneof='dosen' 'tenaga pendidik' 'mahasiswa' 'orang tua' 'perseptor'"`
	Nama_Lengkap string `json:"nama_lengkap" form:"nama_lengkap" validate:"required"`
	No_Hp        string `json:"no_hp" form:"no_hp" validate:"required,numeric"`
	Username     string `json:"username" form:"username" validate:"required"`
	Password     string `json:"password"  form:"password" validate:"required"`
}

type GetUserParams struct {
	Id_User string `json:"id_user" form:"id_user" validate:"required,uuid4"`
}

type UpdateUsersParams struct {
	Id_User string `json:"id_user" form:"id_user" validate:"required,uuid4"`
}

type UpdateUsersForm struct {
	Email        string `json:"email" form:"email" validate:"required,email"`
	Id_Person    string `json:"id_person" form:"id_person" validate:"required,uuid4"`
	Jenis_User   string `json:"jenis_user" form:"jenis_user" validate:"required,oneof='dosen' 'tenaga tendidik' 'mahasiswa' 'orang tua' 'perseptor'"`
	Nama_Lengkap string `json:"nama_lengkap" form:"nama_lengkap" validate:"required"`
	No_Hp        string `json:"no_hp" form:"no_hp" validate:"required,numeric"`
	Username     string `json:"username" form:"username" validate:"required"`
}

type UpdatePasswordParams struct {
	Id_User string `json:"id_user" form:"id_user" validate:"required,uuid4"`
}

type UpdatePasswordForm struct {
	Password string `json:"password" form:"password" validate:"required"`
}

type DeleteUserParams struct {
	Id_User string `json:"id_user" form:"id_user" validate:"required,uuid4"`
}

type GetIDUserParams struct {
	Person_ID string `json:"person_id" form:"person_id" validate:"required,uuid4"`
}

type GetProdiParams struct {
	ID_Prodi string `json:"id_prodi" form:"id_prodi" validate:"required,uuid4"`
}

type GetMahasiswaParams struct {
	ID_Registrasi_Mahasiswa string `json:"id_registrasi_mahasiswa" form:"id_registrasi_mahasiswa" validate:"required,uuid4"`
}

// ==== Data struct (hasil query / transfer antar layer) ====

type Mahasiswa struct {
	Email     string
	IdPerson  string
	JenisUser string
	Nama      string
	NoHp      string
	Username  string
	Password  string
	Avatar    string
}

type Pegawai struct {
	IDUser   string
	Password string
}

type DaftarUser struct {
	PersonID    string `json:"person_id"`
	NamaLengkap string `json:"nama_lengkap"`
	NIM         string `json:"nim"`
}

type DetailUser struct {
	PersonID    string `json:"person_id"`
	NamaLengkap string `json:"nama_lengkap"`
	Email       string `json:"email"`
	NoHp        string `json:"no_hp"`
	Username    string `json:"username"`
	JenisUser   string `json:"jenis_user"`
}

type BulkMahasiswaItem struct {
	Email       string `json:"email"`
	IdPerson    string `json:"id_person"`
	NamaLengkap string `json:"nama_lengkap"`
	NoHp        string `json:"no_hp"`
	Username    string `json:"username"`
	Password    string `json:"password"`
}

type BulkCreateResult struct {
	Total       int      `json:"total"`
	Berhasil    int      `json:"berhasil"`
	Gagal       int      `json:"gagal"`
	DetailGagal []string `json:"detail_gagal"`
}
