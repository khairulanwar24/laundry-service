package models

import (
	"sso-service/database"
	middleware "sso-service/middlewares"
	"sso-service/types"
	"strings"
)

func GetMstMenu(id_master_aplikasi string, limitParam int, offseParam int, order, filter string) types.Response {
	var resp types.Response

	sRecursive := ``
	sTable := `SELECT id_master_menu, id_master_aplikasi, nama_menu, deskripsi, "order", icon FROM master_menu WHERE id_master_aplikasi = '` + id_master_aplikasi + `' AND status_data = true`

	sFilter := ``
	// Kondisi filter, jika filter tidak kosong maka ditambahkan ke query
	if filter != "" {
		filter = "%" + strings.ToLower(filter) + "%"
		sFilter = `AND (LOWER(nama_menu) LIKE ` + "'" + filter + "'" + ` OR LOWER(deskripsi) LIKE ` + "'" + filter + "')"
	} else {
		sFilter = ``
	}

	user := middleware.Datatables(sRecursive, sTable, order, sFilter, limitParam, offseParam)
	resp.Success = true
	resp.Message = "Success"
	resp.Data = user

	return resp
}

func GetDetailMstMenu(id_master_menu string) types.Response {
	var resp types.Response

	var user []map[string]interface{}

	result := database.DB.Raw(`Select id_master_menu, id_master_aplikasi, nama_menu, deskripsi, "order", icon FROM master_menu where id_master_menu = ?`, id_master_menu).First(&user)
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

func GetDetailMstMenuModul(id_master_modul string) types.Response {
	var resp types.Response

	var user []map[string]interface{}

	result := database.DB.Raw(` SELECT id_master_modul, id_master_menu, "order" as order, nama_modul, path, deskripsi, id_master_aplikasi, icon FROM master_modul where id_master_modul = ?`, id_master_modul).First(&user)
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

func GetMstMenuModul(id_master_menu string, limitParam int, offseParam int, order, filter string) types.Response {
	var resp types.Response

	sRecursive := ``
	sTable := `SELECT id_master_modul, id_master_menu, "order" as order, nama_modul, path, deskripsi, id_master_aplikasi, icon FROM master_modul WHERE id_master_menu = '` + id_master_menu + `' AND status_data = true`

	sFilter := ``
	// Kondisi filter, jika filter tidak kosong maka ditambahkan ke query
	if filter != "" {
		filter = "%" + strings.ToLower(filter) + "%"
		sFilter = `AND (LOWER(nama_modul) LIKE ` + "'" + filter + "'" + ` OR LOWER(path) LIKE ` + "'" + filter + "'" + ` OR LOWER(deskripsi) LIKE ` + "'" + filter + "')"
	} else {
		sFilter = ``
	}

	user := middleware.Datatables(sRecursive, sTable, order, sFilter, limitParam, offseParam)
	resp.Success = true
	resp.Message = "Success"
	resp.Data = user

	return resp
}

func CreateMstMenu(id_master_aplikasi, nama_menu, deskripsi, order, icon string) types.Response {
	var resp types.Response
	result := database.DB.Exec(`INSERT INTO master_menu (id_master_aplikasi, nama_menu, deskripsi, "order", icon) VALUES (?, ?, ?, ?, ?)`, id_master_aplikasi, nama_menu, deskripsi, order, icon)

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

func CreateMstMenuModul(id_master_aplikasi, id_master_menu, nama_modul, path, deskripsi, order, icon string) types.Response {
	var resp types.Response
	result := database.DB.Exec(`INSERT INTO master_modul ( id_master_menu, "order", nama_modul, path, deskripsi, id_master_aplikasi, icon) VALUES (?, ?, ?, ?, ?, ?, ?)`, id_master_menu, order, nama_modul, path, deskripsi, id_master_aplikasi, icon)

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

func UpdateMstMenu(id_master_menu, nama_menu, deskripsi, order, icon string) types.Response {
	var resp types.Response

	result := database.DB.Exec(`UPDATE master_menu
								SET nama_menu = ?, deskripsi = ?, "order" = ?, icon = ?
								WHERE id_master_menu = ?`, nama_menu, deskripsi, order, icon, id_master_menu)

	if result.Error != nil {
		resp.Success = false
		resp.Message = "Gagal Update master_menu"
		resp.Data = nil
		return resp
	} else if result.RowsAffected == 0 {
		resp.Success = false
		resp.Message = "master_menu tidak ada"
		resp.Data = nil
		return resp
	}

	resp.Success = true
	resp.Message = "Success"
	resp.Data = nil

	return resp
}

func UpdateMstMenuModul(id_master_modul, id_master_aplikasi, id_master_menu, nama_modul, path, deskripsi, order, icon string) types.Response {
	var resp types.Response

	result := database.DB.Exec(`UPDATE master_modul
								SET id_master_menu = ?, "order" = ?, nama_modul = ?, path = ?, deskripsi = ?, id_master_aplikasi = ?, icon = ?
								WHERE id_master_modul = ?`, id_master_menu, order, nama_modul, path, deskripsi, id_master_aplikasi, icon, id_master_modul)

	if result.Error != nil {
		resp.Success = false
		resp.Message = "Gagal Update master modul"
		resp.Data = nil
		return resp
	} else if result.RowsAffected == 0 {
		resp.Success = false
		resp.Message = "master modul tidak ada"
		resp.Data = nil
		return resp
	}

	resp.Success = true
	resp.Message = "Success"
	resp.Data = nil

	return resp
}

func DeleteMstMenu(id_master_menu string) types.Response {
	var resp types.Response

	result := database.DB.Exec(`UPDATE master_menu
								SET status_data = false
								WHERE id_master_menu = ?`, id_master_menu)

	// panic(result.RowsAffected)
	if result.Error != nil {
		resp.Success = false
		resp.Message = "Gagal Delete Master Menu"
		resp.Data = nil
		return resp
	} else if result.RowsAffected == 0 {
		resp.Success = false
		resp.Message = "Master Menu tidak ada"
		resp.Data = nil
		return resp
	}

	resp.Success = true
	resp.Message = "Success"
	resp.Data = nil

	return resp
}

func DeleteMstMenuModul(id_master_modul string) types.Response {
	var resp types.Response

	result := database.DB.Exec(`UPDATE master_modul
								SET status_data = false
								WHERE id_master_modul = ?`, id_master_modul)

	// panic(result.RowsAffected)
	if result.Error != nil {
		resp.Success = false
		resp.Message = "Gagal Delete Master Modul"
		resp.Data = nil
		return resp
	} else if result.RowsAffected == 0 {
		resp.Success = false
		resp.Message = "Master Modul tidak ada"
		resp.Data = nil
		return resp
	}

	resp.Success = true
	resp.Message = "Success"
	resp.Data = nil

	return resp
}
