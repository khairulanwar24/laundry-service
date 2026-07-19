package models

import (
	"sso-service/database"
	middleware "sso-service/middlewares"
	"sso-service/types"
	"strings"
)

func GetUsers(limitParam int, offsetParam int, order, filter string) types.Response {
	var resp types.Response

	sRecursive := ``
	sTable := ` select id_user
								, username
								, email
								, nama_lengkap
								, avatar
								, id_person
								, jenis_user
								, no_hp from 
								users where status_data = true`

	sFilter := ``
	if filter != "" {
		filter = "%" + strings.ToLower(filter) + "%"
		sFilter = `and (LOWER(username) LIKE ` + "'" + filter + "'" + ` OR LOWER(email) LIKE ` + "'" + filter + "'" + ` OR LOWER(nama_lengkap) LIKE ` + "'" + filter + "'" + ` OR LOWER(no_hp) LIKE ` + "'" + filter + "')"
	} else {
		sFilter = ``
	}

	user := middleware.Datatables(sRecursive, sTable, order, sFilter, limitParam, offsetParam)
	resp.Success = true
	resp.Message = "Success"
	resp.Data = user

	return resp
}

func CreateUsers(email, id_person, jenis_user, nama_lengkap, no_hp, username, password, avatar string) types.Response {
	var resp types.Response

	// cek apakah id_person sudah ada
	var count int64
	err := database.DB.Raw("SELECT COUNT(*) FROM users WHERE id_person = ? and status_data = true", id_person).
		Scan(&count).Error

	if err != nil {
		resp.Success = false
		resp.Message = err.Error()
		return resp
	}

	if count > 0 {
		resp.Success = false
		resp.Message = "User dengan id_person sudah ada"
		return resp
	}

	// hash password
	password_hash, _ := middleware.HashPassword(password)

	// insert data baru
	result := database.DB.Exec(`INSERT INTO users
								( email
								, id_person
								, jenis_user
								, nama_lengkap
								, no_hp
								, username
								, password
								, first_login
								, status_data
								, avatar
								)
								VALUES
								(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		email, id_person, jenis_user, nama_lengkap, no_hp, username, password_hash, true, true, avatar)

	if result.Error != nil {
		resp.Success = false
		resp.Message = result.Error.Error()
		resp.Data = result.Error.Error()
		return resp
	}

	resp.Success = true
	resp.Message = "Success"
	return resp
}

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

func GenerateUserMahasiswa() types.Response {
	var resp types.Response
	var mahasiswas []Mahasiswa

	// ambil data mahasiswa aktif
	result := database.DBAkademik.Raw(`
		select 
			'' as email, 
			id_registrasi_mahasiswa as id_person,
			'Mahasiswa' as jenis_user, 
			nama_mahasiswa as nama, 
			'' as no_hp,
			nim as username,
			nim as password,
			'' as avatar
		from list_mahasiswa
		where nama_status_mahasiswa = 'AKTIF' and id_periode_masuk = '20251'`).
		Scan(&mahasiswas)

	if result.Error != nil {
		resp.Success = false
		resp.Message = result.Error.Error()
		resp.Data = nil
		return resp
	}
	// fmt.Println(mahasiswas)

	// looping insert user baru
	for _, m := range mahasiswas {
		_ = CreateUsers(
			m.Email,
			m.IdPerson,
			m.JenisUser,
			m.Nama,
			m.NoHp,
			m.Username,
			m.Password, // nanti di-hash di CreateUsers
			m.Avatar,
		)
	}

	resp.Success = true
	resp.Message = "Success generate user mahasiswa"
	resp.Data = nil
	return resp
}

type Pegawai struct {
	IDUser   string
	Password string
}

func GenerateUserDosen() types.Response {
	var resp types.Response
	var dosen []Pegawai

	// ambil data mahasiswa aktif
	result := database.DB.Raw(`
		select 
			id_user, 
			password
		from users
		where jenis_user = 'Dosen' and status_data = true and password=username`).
		Scan(&dosen)

	if result.Error != nil {
		resp.Success = false
		resp.Message = result.Error.Error()
		resp.Data = nil
		return resp
	}

	// looping insert user baru
	for _, m := range dosen {
		_ = UpdatePassword(
			m.IDUser,
			m.Password,
		)
		_ = CreateGroupAksesUserApps(m.IDUser, "acce4365-6835-441e-88d2-56539fb0e823", "f66ff11d-a59c-4b17-b039-3a5b3ed03508", "true")
	}

	resp.Success = true
	resp.Message = "Success generate user dosen"
	resp.Data = nil
	return resp
}

func GenerateUserTendik() types.Response {
	var resp types.Response
	var dosen []Pegawai

	// ambil data mahasiswa aktif
	result := database.DB.Raw(`
		select 
			id_user, 
			password
		from users
		where jenis_user = 'Tenaga Pendidik' and status_data = true and password=username`).
		Scan(&dosen)

	if result.Error != nil {
		resp.Success = false
		resp.Message = result.Error.Error()
		resp.Data = nil
		return resp
	}

	// looping insert user baru
	for _, m := range dosen {
		_ = UpdatePassword(
			m.IDUser,
			m.Password,
		)
		_ = CreateGroupAksesUserApps(m.IDUser, "acce4365-6835-441e-88d2-56539fb0e823", "0daae1e8-ffb2-4c9e-bf32-6dbc5a8da847", "true")
	}

	resp.Success = true
	resp.Message = "Success generate user Tendik"
	resp.Data = nil
	return resp
}

func GetUser(id_user string) types.Response {
	var resp types.Response

	var user []map[string]interface{}

	result := database.DB.Raw(`SELECT
								id_user
								, username
								, email
								, nama_lengkap
								, avatar
								, id_person
								, jenis_user
								, no_hp 
								, password
								FROM  users where id_user = ? and status_data = true`, id_user).First(&user)
	if result.Error != nil {
		resp.Success = false
		resp.Message = "Gagal Mendapatkan users"
		resp.Data = nil
		return resp
	} else if result.RowsAffected == 0 {
		resp.Success = false
		resp.Message = "users tidak ada"
		resp.Data = nil
		return resp
	}

	resp.Success = true
	resp.Message = "Success"
	resp.Data = user

	return resp
}

func UpdateUsers(id_user, avatar, email, id_person, jenis_user, nama_lengkap, no_hp, username string) types.Response {
	var resp types.Response
	// var user []map[string]interface{}
	// result := ""
	if avatar == "" {
		// fmt.Println("1")
		result := database.DB.Exec(`UPDATE users
								SET  email = ?, id_person = ?,jenis_user = ?, nama_lengkap = ?, no_hp = ?, username = ?
								WHERE id_user = ?`, email, id_person, jenis_user, nama_lengkap, no_hp, username, id_user)

		if result.Error != nil {
			resp.Success = false
			resp.Message = "Gagal Mendapatkan users"
			resp.Data = nil
			return resp
		} else if result.RowsAffected == 0 {
			resp.Success = false
			resp.Message = "users tidak ada"
			resp.Data = nil
			return resp
		}
	} else {
		// fmt.Println("2")

		result := database.DB.Exec(`UPDATE users
								SET avatar = ?, email = ?, id_person = ?,jenis_user = ?, nama_lengkap = ?, no_hp = ?, username = ?
								WHERE id_user = ?`, avatar, email, id_person, jenis_user, nama_lengkap, no_hp, username, id_user)
		if result.Error != nil {
			resp.Success = false
			resp.Message = "Gagal Mendapatkan users"
			resp.Data = nil
			return resp
		} else if result.RowsAffected == 0 {
			resp.Success = false
			resp.Message = "users tidak ada"
			resp.Data = nil
			return resp
		}

	}

	// panic(result.RowsAffected)

	resp.Success = true
	resp.Message = "Success"
	resp.Data = nil

	return resp
}

func UpdatePassword(id_user, password string) types.Response {
	var resp types.Response
	var user []map[string]interface{}

	result := database.DB.Raw(`SELECT
								id_user
								, first_login
								FROM  users where id_user = ? and status_data = true`, id_user).First(&user)
	if result.Error != nil {
		resp.Success = false
		resp.Message = "Gagal Mendapatkan users"
		resp.Data = nil
		return resp
	}
	password_hash, _ := middleware.HashPassword(password)
	first_login := user[0]["first_login"].(bool)
	if first_login {
		result = database.DB.Exec(`UPDATE users
								SET password = ?,tgl_update = NOW(),first_login = false,tgl_first_login = NOW()
								WHERE id_user = ?`, password_hash, id_user)

	} else {
		result = database.DB.Exec(`UPDATE users
								SET password = ?,tgl_update = NOW()
								WHERE id_user = ?`, password_hash, id_user)
	}

	// panic(result.RowsAffected)
	if result.Error != nil {
		resp.Success = false
		resp.Message = "Gagal Update users"
		resp.Data = nil
		return resp
	} else if result.RowsAffected == 0 {
		resp.Success = false
		resp.Message = "users tidak ada"
		resp.Data = nil
		return resp
	}

	resp.Success = true
	resp.Message = "Success"
	resp.Data = nil

	return resp
}

func DeleteUser(id_user string) types.Response {
	var resp types.Response

	result := database.DB.Exec(`delete from users
								WHERE id_user = ?`, id_user)

	// panic(result.RowsAffected)
	if result.Error != nil {
		resp.Success = false
		resp.Message = "Gagal Delete users"
		resp.Data = nil
		return resp
	} else if result.RowsAffected == 0 {
		resp.Success = false
		resp.Message = "users tidak ada"
		resp.Data = nil
		return resp
	}

	resp.Success = true
	resp.Message = "Success"
	resp.Data = nil

	return resp
}

type DaftarUser struct {
	PersonID    string `json:"person_id"`
	NamaLengkap string `json:"nama_lengkap"`
	NIM         string `json:"nim"`
}

func GetUsersDosen() types.Response {
	var dosen []DaftarUser
	var resp types.Response
	database.DBDigiclass.Raw(`select id_master_dosen as person_id, nama_dosen as nama_lengkap from master_dosen`).Scan(&dosen)

	resp.Data = dosen
	resp.Success = true
	resp.Message = "Success"

	return resp

}

type DetailUser struct {
	PersonID    string `json:"person_id"`
	NamaLengkap string `json:"nama_lengkap"`
	Email       string `json:"email"`
	NoHp        string `json:"no_hp"`
	Username    string `json:"username"`
	JenisUser   string `json:"jenis_user"`
}

func GetDetailDosen(person_id string) types.Response {
	var dosen []DetailUser
	var resp types.Response
	database.DBDigiclass.Raw(`select id_master_dosen as person_id, nama_dosen as nama_lengkap,'tenaga pendidik' as jenis_user from master_dosen where id_master_dosen = ?`, person_id).Scan(&dosen)

	resp.Data = dosen
	resp.Success = true
	resp.Message = "Success"

	return resp

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

func BulkCreateUsersMahasiswa(items []BulkMahasiswaItem) types.Response {
	var resp types.Response
	result := BulkCreateResult{
		Total:       len(items),
		DetailGagal: []string{},
	}

	for _, item := range items {
		r := CreateUsers(item.Email, item.IdPerson, "mahasiswa", item.NamaLengkap, item.NoHp, item.Username, item.Password, "")
		if r.Success {
			result.Berhasil++
		} else {
			result.Gagal++
			result.DetailGagal = append(result.DetailGagal, item.Username+": "+r.Message)
		}
	}

	resp.Success = true
	resp.Message = "Bulk create selesai"
	resp.Data = result
	return resp
}

func GetDetailMahasiswa(id_registrasi_mahasiswa string) types.Response {
	var mahasiswa []map[string]any
	var resp types.Response
	database.DBAkademik.Raw(`select * from list_mahasiswa where id_registrasi_mahasiswa = ?`, id_registrasi_mahasiswa).Scan(&mahasiswa)
	resp.Data = mahasiswa
	resp.Success = true
	resp.Message = "Success"

	return resp
}

func GetUsersMahasiswa(id_prodi string) types.Response {

	var mahasiswa []DaftarUser
	var resp types.Response
	database.DBAkademik.Raw(`select id_registrasi_mahasiswa as person_id, nama_mahasiswa as nama_lengkap, nim from list_mahasiswa where id_prodi = ?`, id_prodi).Scan(&mahasiswa)
	resp.Data = mahasiswa
	resp.Success = true
	resp.Message = "Success"

	return resp
}

func GetUsersMahasiswaData(idProdi, idAngkatan, order, filter string, limit, offset int) types.Response {
	var resp types.Response

	sBase := `select lm.id_registrasi_mahasiswa as person_id, lm.nama_mahasiswa as nama_lengkap, lm.nim
	          from list_mahasiswa lm`

	sWhere := ` where lm.id_prodi = ?`
	baseArgs := []interface{}{idProdi}

	if idAngkatan != "" {
		sBase += ` inner join mahasiswa_angkatan ma on ma.id_registrasi_mahasiswa = lm.id_registrasi_mahasiswa`
		sWhere += ` and ma.id_angkatan_mahasiswa = ?`
		baseArgs = append(baseArgs, idAngkatan)
	}

	sFilter := ""
	filterArgs := []interface{}{}
	if filter != "" {
		f := "%" + strings.ToLower(filter) + "%"
		sFilter = ` and (LOWER(lm.nama_mahasiswa) LIKE ? OR LOWER(lm.nim) LIKE ?)`
		filterArgs = []interface{}{f, f}
	}

	sQuery := sBase + sWhere

	// Count total (tanpa filter teks)
	var totalCount []map[string]interface{}
	database.DBAkademik.Raw(`select count(*) as total from (`+sQuery+`) as t`, baseArgs...).First(&totalCount)
	total := 0
	if len(totalCount) > 0 {
		if v, ok := totalCount[0]["total"].(int64); ok {
			total = int(v)
		}
	}

	// Count filtered
	filteredArgs := append(baseArgs, filterArgs...)
	var filteredCount []map[string]interface{}
	database.DBAkademik.Raw(`select count(*) as total from (`+sQuery+sFilter+`) as t`, filteredArgs...).First(&filteredCount)
	totalFiltered := 0
	if len(filteredCount) > 0 {
		if v, ok := filteredCount[0]["total"].(int64); ok {
			totalFiltered = int(v)
		}
	}

	orderStr := ""
	if order != "" {
		orderStr = " order by " + order
	}

	dataArgs := append(filteredArgs, limit, offset)
	data := []DaftarUser{}
	database.DBAkademik.Raw(sQuery+sFilter+orderStr+" LIMIT ? OFFSET ?", dataArgs...).Scan(&data)

	resp.Success = true
	resp.Message = "Success"
	resp.Data = map[string]interface{}{
		"data":            data,
		"recordsTotal":    total,
		"recordsFiltered": totalFiltered,
	}

	return resp
}
