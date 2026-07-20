// Package repositories (user) adalah lapisan akses data untuk domain user.
// Seluruh query mentah dipindahkan ke sini dari models/user_model.go (query dijaga verbatim).
package repositories

import (
	"context"
	"strings"

	"laundry-service/domain/dto"
	middleware "laundry-service/middlewares"

	"gorm.io/gorm"
)

// UserRepository memegang 3 koneksi database yang dipakai domain user.
type UserRepository struct {
	db          *gorm.DB // database utama (sso)
	dbAkademik  *gorm.DB // database akademik
	dbDigiclass *gorm.DB // database digiclass
}

// IUserRepository adalah kontrak akses data domain user.
type IUserRepository interface {
	GetUsers(order, filter string, limit, offset int) map[string]interface{}
	CountByPerson(ctx context.Context, idPerson string) (int64, error)
	Insert(
		ctx context.Context,
		email,
		idPerson,
		jenisUser,
		namaLengkap,
		noHp,
		username,
		passwordHash,
		avatar string,
	) error
	FindByID(ctx context.Context, idUser string) ([]map[string]interface{}, int64, error)
	UpdateWithoutAvatar(
		ctx context.Context,
		email,
		idPerson,
		jenisUser,
		namaLengkap,
		noHp,
		username,
		idUser string,
	) (int64, error)
	UpdateWithAvatar(
		ctx context.Context,
		avatar,
		email,
		idPerson,
		jenisUser,
		namaLengkap,
		noHp,
		username,
		idUser string,
	) (int64, error)
	FindFirstLogin(ctx context.Context, idUser string) ([]map[string]interface{}, error)
	UpdatePasswordFirstLogin(ctx context.Context, passwordHash, idUser string) (int64, error)
	UpdatePasswordNormal(ctx context.Context, passwordHash, idUser string) (int64, error)
	Delete(ctx context.Context, idUser string) (int64, error)
	GetDosen(ctx context.Context) ([]dto.DaftarUser, error)
	GetDetailDosen(ctx context.Context, personID string) ([]dto.DetailUser, error)
	GetMahasiswa(ctx context.Context, idProdi string) ([]dto.DaftarUser, error)
	GetDetailMahasiswa(ctx context.Context, idRegistrasi string) ([]map[string]any, error)
	GetMahasiswaData(
		ctx context.Context,
		idProdi,
		idAngkatan,
		order,
		filter string,
		limit,
		offset int,
	) ([]dto.DaftarUser, int, int, error)
	GetMahasiswaSourceForGenerate(ctx context.Context) ([]dto.Mahasiswa, error)
	GetPegawaiNeedingPassword(ctx context.Context, jenisUser string) ([]dto.Pegawai, error)
}

// NewUserRepository membuat instance UserRepository baru.
func NewUserRepository(db, dbAkademik, dbDigiclass *gorm.DB) IUserRepository {
	return &UserRepository{db: db, dbAkademik: dbAkademik, dbDigiclass: dbDigiclass}
}

// GetUsers membangun query datatable & mengembalikan hasil dari helper Datatables.
func (r *UserRepository) GetUsers(order, filter string, limit, offset int) map[string]interface{} {
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
	}

	return middleware.Datatables(sRecursive, sTable, order, sFilter, limit, offset)
}

// CountByPerson menghitung user aktif berdasarkan id_person.
func (r *UserRepository) CountByPerson(ctx context.Context, idPerson string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Raw("SELECT COUNT(*) FROM users WHERE id_person = ? and status_data = true", idPerson).
		Scan(&count).Error
	return count, err
}

// Insert menyimpan user baru (password sudah dalam bentuk hash).
func (r *UserRepository) Insert(
	ctx context.Context,
	email,
	idPerson,
	jenisUser,
	namaLengkap,
	noHp,
	username,
	passwordHash,
	avatar string,
) error {
	return r.db.WithContext(ctx).Exec(`INSERT INTO users
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
		email, idPerson, jenisUser, namaLengkap, noHp, username, passwordHash, true, true, avatar).Error
}

// FindByID mengambil detail user berdasarkan id_user (mengembalikan data, rowsAffected, error).
func (r *UserRepository) FindByID(ctx context.Context, idUser string) ([]map[string]interface{}, int64, error) {
	var user []map[string]interface{}
	result := r.db.WithContext(ctx).Raw(`SELECT
								id_user
								, username
								, email
								, nama_lengkap
								, avatar
								, id_person
								, jenis_user
								, no_hp
								, password
								FROM  users where id_user = ? and status_data = true`, idUser).First(&user)
	return user, result.RowsAffected, result.Error
}

// UpdateWithoutAvatar memperbarui data user tanpa mengubah avatar.
func (r *UserRepository) UpdateWithoutAvatar(
	ctx context.Context,
	email,
	idPerson,
	jenisUser,
	namaLengkap,
	noHp,
	username,
	idUser string,
) (int64, error) {
	result := r.db.WithContext(ctx).Exec(`UPDATE users
								SET  email = ?, id_person = ?,jenis_user = ?, nama_lengkap = ?, no_hp = ?, username = ?
								WHERE id_user = ?`, email, idPerson, jenisUser, namaLengkap, noHp, username, idUser)
	return result.RowsAffected, result.Error
}

// UpdateWithAvatar memperbarui data user termasuk avatar.
func (r *UserRepository) UpdateWithAvatar(
	ctx context.Context,
	avatar,
	email,
	idPerson,
	jenisUser,
	namaLengkap,
	noHp,
	username,
	idUser string,
) (int64, error) {
	result := r.db.WithContext(ctx).Exec(`UPDATE users
								SET avatar = ?, email = ?, id_person = ?,jenis_user = ?, nama_lengkap = ?, no_hp = ?, username = ?
								WHERE id_user = ?`, avatar, email, idPerson, jenisUser, namaLengkap, noHp, username, idUser)
	return result.RowsAffected, result.Error
}

// FindFirstLogin mengambil status first_login user.
func (r *UserRepository) FindFirstLogin(ctx context.Context, idUser string) ([]map[string]interface{}, error) {
	var user []map[string]interface{}
	result := r.db.WithContext(ctx).Raw(`SELECT
								id_user
								, first_login
								FROM  users where id_user = ? and status_data = true`, idUser).First(&user)
	return user, result.Error
}

// UpdatePasswordFirstLogin memperbarui password saat login pertama (menandai first_login false).
func (r *UserRepository) UpdatePasswordFirstLogin(ctx context.Context, passwordHash, idUser string) (int64, error) {
	result := r.db.WithContext(ctx).Exec(`UPDATE users
								SET password = ?,tgl_update = NOW(),first_login = false,tgl_first_login = NOW()
								WHERE id_user = ?`, passwordHash, idUser)
	return result.RowsAffected, result.Error
}

// UpdatePasswordNormal memperbarui password biasa.
func (r *UserRepository) UpdatePasswordNormal(ctx context.Context, passwordHash, idUser string) (int64, error) {
	result := r.db.WithContext(ctx).Exec(`UPDATE users
								SET password = ?,tgl_update = NOW()
								WHERE id_user = ?`, passwordHash, idUser)
	return result.RowsAffected, result.Error
}

// Delete menghapus user berdasarkan id_user.
func (r *UserRepository) Delete(ctx context.Context, idUser string) (int64, error) {
	result := r.db.WithContext(ctx).Exec(`delete from users
								WHERE id_user = ?`, idUser)
	return result.RowsAffected, result.Error
}

// GetDosen mengambil daftar dosen dari database digiclass.
func (r *UserRepository) GetDosen(ctx context.Context) ([]dto.DaftarUser, error) {
	var dosen []dto.DaftarUser
	err := r.dbDigiclass.WithContext(ctx).Raw(`select id_master_dosen as person_id, nama_dosen as nama_lengkap from master_dosen`).Scan(&dosen).Error
	return dosen, err
}

// GetDetailDosen mengambil detail satu dosen dari database digiclass.
func (r *UserRepository) GetDetailDosen(ctx context.Context, personID string) ([]dto.DetailUser, error) {
	var dosen []dto.DetailUser
	err := r.dbDigiclass.WithContext(ctx).Raw(`select id_master_dosen as person_id, nama_dosen as nama_lengkap,'tenaga pendidik' as jenis_user from master_dosen where id_master_dosen = ?`, personID).Scan(&dosen).Error
	return dosen, err
}

// GetMahasiswa mengambil daftar mahasiswa per prodi dari database akademik.
func (r *UserRepository) GetMahasiswa(ctx context.Context, idProdi string) ([]dto.DaftarUser, error) {
	var mahasiswa []dto.DaftarUser
	err := r.dbAkademik.WithContext(ctx).Raw(`select id_registrasi_mahasiswa as person_id, nama_mahasiswa as nama_lengkap, nim from list_mahasiswa where id_prodi = ?`, idProdi).Scan(&mahasiswa).Error
	return mahasiswa, err
}

// GetDetailMahasiswa mengambil detail satu mahasiswa dari database akademik.
func (r *UserRepository) GetDetailMahasiswa(ctx context.Context, idRegistrasi string) ([]map[string]any, error) {
	var mahasiswa []map[string]any
	err := r.dbAkademik.WithContext(ctx).Raw(`select * from list_mahasiswa where id_registrasi_mahasiswa = ?`, idRegistrasi).Scan(&mahasiswa).Error
	return mahasiswa, err
}

// GetMahasiswaData mengambil data mahasiswa dengan paginasi + filter (server-side datatable) dari database akademik.
func (r *UserRepository) GetMahasiswaData(
	ctx context.Context,
	idProdi,
	idAngkatan,
	order,
	filter string,
	limit,
	offset int,
) ([]dto.DaftarUser, int, int, error) {
	db := r.dbAkademik.WithContext(ctx)

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
	db.Raw(`select count(*) as total from (`+sQuery+`) as t`, baseArgs...).First(&totalCount)
	total := 0
	if len(totalCount) > 0 {
		if v, ok := totalCount[0]["total"].(int64); ok {
			total = int(v)
		}
	}

	// Count filtered
	filteredArgs := append(baseArgs, filterArgs...)
	var filteredCount []map[string]interface{}
	db.Raw(`select count(*) as total from (`+sQuery+sFilter+`) as t`, filteredArgs...).First(&filteredCount)
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
	data := []dto.DaftarUser{}
	err := db.Raw(sQuery+sFilter+orderStr+" LIMIT ? OFFSET ?", dataArgs...).Scan(&data).Error

	return data, total, totalFiltered, err
}

// GetMahasiswaSourceForGenerate mengambil sumber data mahasiswa aktif untuk generate user (database akademik).
func (r *UserRepository) GetMahasiswaSourceForGenerate(ctx context.Context) ([]dto.Mahasiswa, error) {
	var mahasiswas []dto.Mahasiswa
	err := r.dbAkademik.WithContext(ctx).Raw(`
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
		Scan(&mahasiswas).Error
	return mahasiswas, err
}

// GetPegawaiNeedingPassword mengambil pegawai (Dosen/Tenaga Pendidik) yang password-nya masih sama dengan username.
func (r *UserRepository) GetPegawaiNeedingPassword(ctx context.Context, jenisUser string) ([]dto.Pegawai, error) {
	var pegawai []dto.Pegawai
	err := r.db.WithContext(ctx).Raw(`
		select
			id_user,
			password
		from users
		where jenis_user = ? and status_data = true and password=username`, jenisUser).
		Scan(&pegawai).Error
	return pegawai, err
}
