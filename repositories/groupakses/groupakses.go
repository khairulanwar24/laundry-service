// Package repositories (groupakses) adalah lapisan akses data untuk domain master group akses & group akses.
package repositories

import (
	"context"
	"strings"

	middleware "sso-service/middlewares"

	"gorm.io/gorm"
)

// GroupAksesRepository memegang koneksi database utama.
type GroupAksesRepository struct {
	db *gorm.DB
}

// IGroupAksesRepository adalah kontrak akses data domain group akses.
type IGroupAksesRepository interface {
	CreateMstGroupAkses(ctx context.Context, idMasterAplikasi, namaGroup, deskripsi string) error
	GetMstGroupAkses(idMasterAplikasi string, limit, offset int, order, filter string) map[string]interface{}
	GetMstGroupAksesModul(idMasterAplikasi, idMasterGroup string, limit, offset int, order, filter string) map[string]interface{}
	GetGroupAkses(idMasterGroup string, limit, offset int, order, filter string) map[string]interface{}
	UpdateMstGroupAkses(ctx context.Context, idMasterGroup, namaGroup, deskripsi string) (int64, error)
	GetDetailMstGroupAkses(ctx context.Context, idMasterGroup string) ([]map[string]interface{}, int64, error)
	DeleteMstGroupAkses(ctx context.Context, idMasterGroup string) (int64, error)
	CreateGroupAkses(ctx context.Context, idMasterAplikasi, idMasterGroup, idMasterModul, akses string) error
	DeleteGroupAkses(ctx context.Context, idGroupAkses string) (int64, error)
	GetGroupAksesUserMenu(ctx context.Context, idUser, idMasterAplikasi string) (map[string]interface{}, error)
	GetGroupAksesUserApps(ctx context.Context, idUser string) (map[string]interface{}, error)
	CreateGroupAksesUserApps(ctx context.Context, idUser, idMasterAplikasi, idMasterGroup, statusData string) error
	BulkCreateGroupAksesUserApps(ctx context.Context, idUsers []string, idMasterAplikasi, idMasterGroup string, statusData bool) error
	UpdateGroupAksesUserApps(ctx context.Context, idTransUserGroup, statusData string) (int64, error)
	DeleteGroupAksesUserApps(ctx context.Context, idTransUserGroup string) (int64, error)
}

// NewGroupAksesRepository membuat instance GroupAksesRepository baru.
func NewGroupAksesRepository(db *gorm.DB) IGroupAksesRepository {
	return &GroupAksesRepository{db: db}
}

// CreateMstGroupAkses menyimpan master group baru.
func (r *GroupAksesRepository) CreateMstGroupAkses(ctx context.Context, idMasterAplikasi, namaGroup, deskripsi string) error {
	return r.db.WithContext(ctx).Exec(`INSERT INTO master_group
								( id_master_aplikasi
								, nama_group
								, deskripsi
								)
								VALUES
								(?, ?, ?)`, idMasterAplikasi, namaGroup, deskripsi).Error
}

// GetMstGroupAkses membangun query datatable master group.
func (r *GroupAksesRepository) GetMstGroupAkses(idMasterAplikasi string, limit, offset int, order, filter string) map[string]interface{} {
	sRecursive := ``
	sTable := ` SELECT
					id_master_group, id_master_aplikasi, nama_group, deskripsi
				FROM
					master_group p where id_master_aplikasi = '` + idMasterAplikasi + `' and status_data = true`

	sFilter := ``
	if filter != "" {
		filter = "%" + strings.ToLower(filter) + "%"
		sFilter = `and LOWER(nama_group) LIKE ` + "'" + filter + "'" + ` OR LOWER(deskripsi) LIKE ` + "'" + filter + "'"
	}

	return middleware.Datatables(sRecursive, sTable, order, sFilter, limit, offset)
}

// GetMstGroupAksesModul membangun query datatable modul untuk suatu group.
func (r *GroupAksesRepository) GetMstGroupAksesModul(idMasterAplikasi, idMasterGroup string, limit, offset int, order, filter string) map[string]interface{} {
	sRecursive := `with modul as (select * from master_modul where id_master_aplikasi = '` + idMasterAplikasi + `'  and status_data = true),
group_akses as (select * from group_akses where id_master_group = '` + idMasterGroup + `')`
	sTable := ` select m.*, mm.nama_menu,ga.id_group_akses,ga.id_master_group from modul as m left join group_akses  as ga on m.id_master_modul = ga.id_master_modul
	inner join master_menu as mm on m.id_master_menu = mm.id_master_menu
	where m.status_data = true and mm.status_data = true`

	sFilter := ``
	if filter != "" {
		filter = "%" + strings.ToLower(filter) + "%"
		sFilter = `and LOWER(m.nama_modul) LIKE ` + "'" + filter + "'" + ` OR LOWER(mm.nama_menu) LIKE ` + "'" + filter + "'"
	}

	return middleware.Datatables(sRecursive, sTable, order, sFilter, limit, offset)
}

// GetGroupAkses membangun query datatable akses modul suatu group.
func (r *GroupAksesRepository) GetGroupAkses(idMasterGroup string, limit, offset int, order, filter string) map[string]interface{} {
	sRecursive := `with master_group_select as (select id_master_aplikasi,id_master_group from master_group WHERE id_master_group =  '` + idMasterGroup + `' and status_data = true)
, master_modul_aktif as ( select mgs.id_master_aplikasi,mgs.id_master_group,mn.id_master_menu ,mn.nama_menu,mn.order as ordermenu, mn.icon as iconmenu,mm.id_master_modul, mm.nama_modul,mm.path,mm.order,mm.icon from master_group_select  mgs inner join master_modul mm on mgs.id_master_aplikasi = mm.id_master_aplikasi and mm.status_data = true
inner join master_menu mn on mm.id_master_menu = mn.id_master_menu)
, grup_akses as (select id_group_akses, id_master_group,id_master_aplikasi,id_master_modul,akses from group_akses where id_master_group =  '` + idMasterGroup + `')`
	sTable := `select mma.id_master_aplikasi,mma.id_master_group,mma.id_master_menu,mma.nama_menu,mma.ordermenu, mma.iconmenu, mma.id_master_modul, mma.nama_modul,mma.path,mma.order,mma.icon , ga.id_group_akses, ga.akses from master_modul_aktif as mma left join group_akses  as ga on mma.id_master_modul = ga.id_master_modul WHERE mma.id_master_group =  '` + idMasterGroup + `'`

	sFilter := ``
	if filter != "" {
		filter = "%" + strings.ToLower(filter) + "%"
		sFilter = `and LOWER(mma.nama_menu) LIKE ` + "'" + filter + "'" + ` OR LOWER( mma.nama_modul) LIKE ` + "'" + filter + "'"
	}

	return middleware.Datatables(sRecursive, sTable, order, sFilter, limit, offset)
}

// UpdateMstGroupAkses memperbarui master group.
func (r *GroupAksesRepository) UpdateMstGroupAkses(ctx context.Context, idMasterGroup, namaGroup, deskripsi string) (int64, error) {
	result := r.db.WithContext(ctx).Exec(`UPDATE master_group
								SET nama_group = ?, deskripsi = ?
								WHERE id_master_group = ?`, namaGroup, deskripsi, idMasterGroup)
	return result.RowsAffected, result.Error
}

// GetDetailMstGroupAkses mengambil detail satu master group.
func (r *GroupAksesRepository) GetDetailMstGroupAkses(ctx context.Context, idMasterGroup string) ([]map[string]interface{}, int64, error) {
	var user []map[string]interface{}
	result := r.db.WithContext(ctx).Raw(`SELECT
								id_master_group,
								id_master_aplikasi
								, nama_group
								, deskripsi
								FROM  master_group where id_master_group = ? and status_data = true`, idMasterGroup).First(&user)
	return user, result.RowsAffected, result.Error
}

// DeleteMstGroupAkses menonaktifkan (soft delete) master group.
func (r *GroupAksesRepository) DeleteMstGroupAkses(ctx context.Context, idMasterGroup string) (int64, error) {
	result := r.db.WithContext(ctx).Exec(`UPDATE master_group
								SET status_data = false
								WHERE id_master_group = ?`, idMasterGroup)
	return result.RowsAffected, result.Error
}

// CreateGroupAkses menyimpan akses modul untuk suatu group.
func (r *GroupAksesRepository) CreateGroupAkses(ctx context.Context, idMasterAplikasi, idMasterGroup, idMasterModul, akses string) error {
	return r.db.WithContext(ctx).Exec(`INSERT INTO group_akses
	(  id_master_group, id_master_aplikasi, id_master_modul, akses) VALUES ( ?, ?, ?, ? )`, idMasterGroup, idMasterAplikasi, idMasterModul, akses).Error
}

// DeleteGroupAkses menghapus akses modul.
func (r *GroupAksesRepository) DeleteGroupAkses(ctx context.Context, idGroupAkses string) (int64, error) {
	result := r.db.WithContext(ctx).Exec(`delete from group_akses
								WHERE id_group_akses = ?`, idGroupAkses)
	return result.RowsAffected, result.Error
}

// GetGroupAksesUserMenu mengambil menu & modul yang bisa diakses user pada suatu aplikasi (dikelompokkan per menu).
func (r *GroupAksesRepository) GetGroupAksesUserMenu(ctx context.Context, idUser, idMasterAplikasi string) (map[string]interface{}, error) {
	var menu []map[string]interface{}
	result := r.db.WithContext(ctx).Raw(`
		SELECT
			mn.nama_menu AS nama_menu,
			mn.order AS menu_order,
			mn.icon AS menu_icon,
			mm.nama_modul AS nama_modul,
			mm.path AS modul_path,
			mm.order AS modul_order,
			mm.icon AS modul_icon
		FROM
			trans_user_group AS tug
			INNER JOIN group_akses AS ga ON tug.id_master_group = ga.id_master_group
			INNER JOIN master_modul AS mm ON ga.id_master_modul = mm.id_master_modul
			INNER JOIN master_menu AS mn ON mm.id_master_menu = mn.id_master_menu
		WHERE
			tug.id_user = ? AND
			tug.id_master_aplikasi = ? AND
			tug.status_data = true
		group by
			mn.nama_menu ,
			mn.order ,
			mn.icon ,
			mm.nama_modul ,
			mm.path,
			mm.order ,
			mm.icon
		ORDER BY
			mn.order,
			mm.order
	`, idUser, idMasterAplikasi).Scan(&menu)

	menus, _ := middleware.Groupby("nama_menu", menu)

	respmenu := make(map[string]interface{})
	respmenu["menu"] = menus

	return respmenu, result.Error
}

// GetGroupAksesUserApps mengambil daftar aplikasi yang bisa diakses user (dikelompokkan per aplikasi).
func (r *GroupAksesRepository) GetGroupAksesUserApps(ctx context.Context, idUser string) (map[string]interface{}, error) {
	var app []map[string]interface{}
	result := r.db.WithContext(ctx).Raw(`
		 		SELECT
                  ma.nama_aplikasi,
					  ma.image,
					  ma.deskripsi,
					  ma.versi_aplikasi,
					  ma.tgl_version,
					  ma.url,
					  ma.image,
					  mg.nama_group,
					  tug.id_trans_user_group,
					  tug.status_data
                FROM
                        trans_user_group AS tug
					inner join
						master_group as mg on tug.id_master_group = mg.id_master_group
					inner join
						master_aplikasi as ma on mg.id_master_aplikasi = ma.id_master_aplikasi
                WHERE
                        tug.id_user = ? AND
						mg.status_data = TRUE and
						ma.status_data = TRUE
	`, idUser).Scan(&app)

	apps, _ := middleware.Groupby("nama_aplikasi", app)

	respapps := make(map[string]interface{})
	respapps["apps"] = apps

	return respapps, result.Error
}

// CreateGroupAksesUserApps menautkan user ke suatu group aplikasi.
func (r *GroupAksesRepository) CreateGroupAksesUserApps(ctx context.Context, idUser, idMasterAplikasi, idMasterGroup, statusData string) error {
	return r.db.WithContext(ctx).Exec(`INSERT INTO trans_user_group
								( id_user
								, id_master_aplikasi
								, id_master_group
								, status_data
								)
								VALUES
								(?, ?, ?, ?)`, idUser, idMasterAplikasi, idMasterGroup, statusData).Error
}

// BulkCreateGroupAksesUserApps menautkan banyak user ke suatu group aplikasi dalam satu transaksi.
func (r *GroupAksesRepository) BulkCreateGroupAksesUserApps(ctx context.Context, idUsers []string, idMasterAplikasi, idMasterGroup string, statusData bool) error {
	tx := r.db.WithContext(ctx).Begin()

	for _, idUser := range idUsers {
		err := tx.Exec(`
			INSERT INTO trans_user_group
			(
				id_user,
				id_master_aplikasi,
				id_master_group,
				status_data
			)
			VALUES (?, ?, ?, ?)
		`, idUser, idMasterAplikasi, idMasterGroup, statusData).Error

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

// UpdateGroupAksesUserApps memperbarui status akses user.
func (r *GroupAksesRepository) UpdateGroupAksesUserApps(ctx context.Context, idTransUserGroup, statusData string) (int64, error) {
	result := r.db.WithContext(ctx).Exec(`UPDATE trans_user_group
								SET status_data = ?
								WHERE id_trans_user_group = ?`, statusData, idTransUserGroup)
	return result.RowsAffected, result.Error
}

// DeleteGroupAksesUserApps menghapus tautan user ke group aplikasi.
func (r *GroupAksesRepository) DeleteGroupAksesUserApps(ctx context.Context, idTransUserGroup string) (int64, error) {
	result := r.db.WithContext(ctx).Exec(`Delete from trans_user_group
								WHERE id_trans_user_group = ?`, idTransUserGroup)
	return result.RowsAffected, result.Error
}
