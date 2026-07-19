package models

import (
	"sso-service/database"
	middleware "sso-service/middlewares"
	"sso-service/types"
	"strings"
)

func CreateMstGroupAkses(id_master_aplikasi, nama_group, deskripsi string) types.Response {

	var resp types.Response
	result := database.DB.Exec(`INSERT INTO master_group
								( id_master_aplikasi
								, nama_group
								, deskripsi
								)
								VALUES
								(?, ?, ?)`, id_master_aplikasi, nama_group, deskripsi)

	if result.Error != nil {
		resp.Success = false
		resp.Message = result.Error.Error()
		resp.Data = nil
		return resp
	}

	resp.Success = true
	resp.Message = "Success"
	resp.Data = nil

	return resp
}

func GetMstGroupAkses(id_master_aplikasi string, limitParam int, offsetParam int, order, filter string) types.Response {
	var resp types.Response

	sRecursive := ``
	sTable := ` SELECT
					id_master_group, id_master_aplikasi, nama_group, deskripsi
				FROM
					master_group p where id_master_aplikasi = '` + id_master_aplikasi + `' and status_data = true`

	sFilter := ``
	if filter != "" {
		filter = "%" + strings.ToLower(filter) + "%"
		sFilter = `and LOWER(nama_group) LIKE ` + "'" + filter + "'" + ` OR LOWER(deskripsi) LIKE ` + "'" + filter + "'"
	} else {
		sFilter = ``
	}

	user := middleware.Datatables(sRecursive, sTable, order, sFilter, limitParam, offsetParam)
	resp.Success = true
	resp.Message = "Success"
	resp.Data = user

	return resp
}

func GetMstGroupAksesModul(id_master_aplikasi, id_master_group string, limitParam int, offsetParam int, order, filter string) types.Response {
	var resp types.Response

	sRecursive := `with modul as (select * from master_modul where id_master_aplikasi = '` + id_master_aplikasi + `'  and status_data = true),
group_akses as (select * from group_akses where id_master_group = '` + id_master_group + `')`
	sTable := ` select m.*, mm.nama_menu,ga.id_group_akses,ga.id_master_group from modul as m left join group_akses  as ga on m.id_master_modul = ga.id_master_modul 
	inner join master_menu as mm on m.id_master_menu = mm.id_master_menu
	where m.status_data = true and mm.status_data = true`

	sFilter := ``
	if filter != "" {
		filter = "%" + strings.ToLower(filter) + "%"
		sFilter = `and LOWER(m.nama_modul) LIKE ` + "'" + filter + "'" + ` OR LOWER(mm.nama_menu) LIKE ` + "'" + filter + "'"
	} else {
		sFilter = ``
	}

	user := middleware.Datatables(sRecursive, sTable, order, sFilter, limitParam, offsetParam)
	resp.Success = true
	resp.Message = "Success"
	resp.Data = user

	return resp
}

func GetGroupAkses(id_master_group string, limitParam int, offsetParam int, order, filter string) types.Response {
	var resp types.Response

	sRecursive := `with master_group_select as (select id_master_aplikasi,id_master_group from master_group WHERE id_master_group =  '` + id_master_group + `' and status_data = true)
, master_modul_aktif as ( select mgs.id_master_aplikasi,mgs.id_master_group,mn.id_master_menu ,mn.nama_menu,mn.order as ordermenu, mn.icon as iconmenu,mm.id_master_modul, mm.nama_modul,mm.path,mm.order,mm.icon from master_group_select  mgs inner join master_modul mm on mgs.id_master_aplikasi = mm.id_master_aplikasi and mm.status_data = true 
inner join master_menu mn on mm.id_master_menu = mn.id_master_menu)
, grup_akses as (select id_group_akses, id_master_group,id_master_aplikasi,id_master_modul,akses from group_akses where id_master_group =  '` + id_master_group + `')`
	sTable := `select mma.id_master_aplikasi,mma.id_master_group,mma.id_master_menu,mma.nama_menu,mma.ordermenu, mma.iconmenu, mma.id_master_modul, mma.nama_modul,mma.path,mma.order,mma.icon , ga.id_group_akses, ga.akses from master_modul_aktif as mma left join group_akses  as ga on mma.id_master_modul = ga.id_master_modul WHERE mma.id_master_group =  '` + id_master_group + `'`

	sFilter := ``
	if filter != "" {
		filter = "%" + strings.ToLower(filter) + "%"
		sFilter = `and LOWER(mma.nama_menu) LIKE ` + "'" + filter + "'" + ` OR LOWER( mma.nama_modul) LIKE ` + "'" + filter + "'"
	} else {
		sFilter = ``
	}

	user := middleware.Datatables(sRecursive, sTable, order, sFilter, limitParam, offsetParam)
	resp.Success = true
	resp.Message = "Success"
	resp.Data = user

	return resp
}

func UpdateMstGroupAkses(id_master_group, nama_group, deskripsi string) types.Response {
	var resp types.Response

	result := database.DB.Exec(`UPDATE master_group
								SET nama_group = ?, deskripsi = ?
								WHERE id_master_group = ?`, nama_group, deskripsi, id_master_group)

	if result.Error != nil {
		resp.Success = false
		resp.Message = "Gagal Update master_group"
		resp.Data = nil
		return resp
	} else if result.RowsAffected == 0 {
		resp.Success = false
		resp.Message = "master_group tidak ada"
		resp.Data = nil
		return resp
	}

	resp.Success = true
	resp.Message = "Success"
	resp.Data = nil

	return resp
}

func GetDetailMstGroupAkses(id_master_group string) types.Response {
	var resp types.Response

	var user []map[string]interface{}

	result := database.DB.Raw(`SELECT
								id_master_group,
								id_master_aplikasi
								, nama_group
								, deskripsi
								FROM  master_group where id_master_group = ? and status_data = true`, id_master_group).First(&user)
	if result.Error != nil {
		resp.Success = false
		resp.Message = "Gagal Mendapatkan Master Group"
		resp.Data = nil
		return resp
	} else if result.RowsAffected == 0 {
		resp.Success = false
		resp.Message = "Master Group tidak ada"
		resp.Data = nil
		return resp
	}

	resp.Success = true
	resp.Message = "Success"
	resp.Data = user

	return resp
}

func DeleteMstGroupAkses(id_master_group string) types.Response {
	var resp types.Response

	result := database.DB.Exec(`UPDATE master_group
								SET status_data = false
								WHERE id_master_group = ?`, id_master_group)

	// panic(result.RowsAffected)
	if result.Error != nil {
		resp.Success = false
		resp.Message = "Gagal Delete Master Group"
		resp.Data = nil
		return resp
	} else if result.RowsAffected == 0 {
		resp.Success = false
		resp.Message = "Master Group tidak ada"
		resp.Data = nil
		return resp
	}

	resp.Success = true
	resp.Message = "Success"
	resp.Data = nil

	return resp
}

func CreateGroupAkses(id_master_aplikasi, id_master_group, id_master_modul, akses string) types.Response {

	var resp types.Response
	result := database.DB.Exec(`INSERT INTO group_akses
	(  id_master_group, id_master_aplikasi, id_master_modul, akses) VALUES ( ?, ?, ?, ? )`, id_master_group, id_master_aplikasi, id_master_modul, akses)

	if result.Error != nil {
		resp.Success = false
		resp.Message = result.Error.Error()
		resp.Data = nil
		return resp
	}

	resp.Success = true
	resp.Message = "Success"
	resp.Data = nil

	return resp
}

func DeleteGroupAkses(id_group_akses string) types.Response {
	var resp types.Response

	result := database.DB.Exec(`delete from group_akses
								WHERE id_group_akses = ?`, id_group_akses)

	// panic(result.RowsAffected)
	if result.Error != nil {
		resp.Success = false
		resp.Message = "Gagal Delete Group Akses"
		resp.Data = nil
		return resp
	} else if result.RowsAffected == 0 {
		resp.Success = false
		resp.Message = "Group Akses tidak ada"
		resp.Data = nil
		return resp
	}

	resp.Success = true
	resp.Message = "Success"
	resp.Data = nil

	return resp
}

func GetGroupAksesUserMenu(id_user, id_master_aplikasi string) types.Response {
	var resp types.Response
	// var user []map[string]interface{}

	var menu []map[string]interface{}
	result := database.DB.Raw(`
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
	`, id_user, id_master_aplikasi).Scan(&menu)

	// panic(menu)
	menus, _ := middleware.Groupby("nama_menu", menu)

	respmenu := make(map[string]interface{})

	respmenu["menu"] = menus
	if result.Error != nil {

		resp.Success = false
		resp.Message = "Gagal Mendapatkan Group Akses"
		resp.Data = nil
		return resp
	}

	resp.Success = true
	resp.Message = "Success"
	resp.Data = respmenu

	return resp
}

func GetGroupAksesUserApps(id_user string) types.Response {
	var resp types.Response
	// var user []map[string]interface{}

	var app []map[string]interface{}
	result := database.DB.Raw(`
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
	`, id_user).Scan(&app)

	// panic(menu)
	apps, _ := middleware.Groupby("nama_aplikasi", app)

	respapps := make(map[string]interface{})

	respapps["apps"] = apps
	if result.Error != nil {

		resp.Success = false
		resp.Message = "Gagal Mendapatkan Akses Aplikasi"
		resp.Data = nil
		return resp
	}

	resp.Success = true
	resp.Message = "Success"
	resp.Data = respapps

	return resp
}

func CreateGroupAksesUserApps(id_user, id_master_aplikasi, id_master_group, status_data string) types.Response {

	var resp types.Response
	result := database.DB.Exec(`INSERT INTO trans_user_group
								( id_user
								, id_master_aplikasi
								, id_master_group
								, status_data
								)
								VALUES
								(?, ?, ?, ?)`, id_user, id_master_aplikasi, id_master_group, status_data)

	if result.Error != nil {
		resp.Success = false
		resp.Message = result.Error.Error()
		resp.Data = nil
		return resp
	}

	resp.Success = true
	resp.Message = "Success"
	resp.Data = nil

	return resp
}
func BulkCreateGroupAksesUserApps(idUsers []string,idMasterAplikasi string,idMasterGroup string,statusData bool) types.Response {

	var resp types.Response

	tx := database.DB.Begin()

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
			resp.Success = false
			resp.Message = err.Error()
			return resp
		}
	}

	tx.Commit()

	resp.Success = true
	resp.Message = "Bulk group akses berhasil disimpan"
	resp.Data = nil

	return resp
}

func UpdateGroupAksesUserApps(id_trans_user_group, status_data string) types.Response {
	var resp types.Response

	result := database.DB.Exec(`UPDATE trans_user_group
								SET status_data = ?
								WHERE id_trans_user_group = ?`, status_data, id_trans_user_group)

	if result.Error != nil {
		resp.Success = false
		resp.Message = "Gagal Update trans_user_group"
		resp.Data = nil
		return resp
	} else if result.RowsAffected == 0 {
		resp.Success = false
		resp.Message = "trans_user_group tidak ada"
		resp.Data = nil
		return resp
	}

	resp.Success = true
	resp.Message = "Success"
	resp.Data = nil

	return resp
}

func DeleteGroupAksesUserApps(id_trans_user_group string) types.Response {
	var resp types.Response

	result := database.DB.Exec(`Delete from trans_user_group
								WHERE id_trans_user_group = ?`, id_trans_user_group)

	// panic(result.RowsAffected)
	if result.Error != nil {
		resp.Success = false
		resp.Message = "Gagal Delete Master Group"
		resp.Data = nil
		return resp
	} else if result.RowsAffected == 0 {
		resp.Success = false
		resp.Message = "Master Group tidak ada"
		resp.Data = nil
		return resp
	}

	resp.Success = true
	resp.Message = "Success"
	resp.Data = nil

	return resp
}
