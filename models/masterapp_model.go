package models

import (
	"sso-service/database"
	middleware "sso-service/middlewares"
	"time"

	"sso-service/types"
	"strings"
)

func GetMasterApps(form *types.GetData) types.Response {
	var resp types.Response

	sRecursive := ``
	sTable := ` SELECT id_master_aplikasi 
											, nama_aplikasi
											, deskripsi
											, versi_aplikasi
											, tgl_version
											, url
											, image from
											master_aplikasi where status_data = true`

	sFilter := ``
	if form.Filter != "" {
		form.Filter = "%" + strings.ToLower(form.Filter) + "%"
		sFilter = `and LOWER(nama_aplikasi) LIKE ` + "'" + form.Filter + "'" + ` OR LOWER(deskripsi) LIKE ` + "'" + form.Filter + "'" + ` OR LOWER(versi_aplikasi) LIKE ` + "'" + form.Filter + "'" + ` OR LOWER(url) LIKE ` + "'" + form.Filter + "'"
	} else {
		sFilter = ``
	}

	// Cetak query yang terbentuk untuk debugging
	// finalQuery := sTable + sFilter
	// fmt.Println("Query yang dieksekusi:", finalQuery)

	masterapp := middleware.Datatables(
		sRecursive, sTable, form.Order, sFilter, form.Limit, form.Offset)

	// masterapp := middleware.Datatables(
	// 	sRecursive, sTable, order, sFilter, limitParam, offsetParam)

	// // Cetak data yang diambil dari database untuk debugging
	// fmt.Println("Data yang diambil dari database:", masterapp)

	resp.Success = true
	resp.Message = "Success"
	resp.Data = masterapp

	return resp

}

func GetMasterAppById(id_master_aplikasi string) types.Response {
	var resp types.Response

	var master []map[string]interface{}

	result := database.DB.Raw(`SELECT id_master_aplikasi
								, nama_aplikasi
								, deskripsi
								, versi_aplikasi
								, tgl_version
								, url
								, image FROM master_aplikasi WHERE id_master_aplikasi = ? AND status_data = true`, id_master_aplikasi).First(&master)
	if result.Error != nil {
		resp.Success = false
		resp.Message = "Gagal mendapatkan aplikasi"
		resp.Data = nil
		return resp
	} else if result.RowsAffected == 0 {
		resp.Success = false
		resp.Message = "aplikasi tidak tersedia"
		resp.Data = nil
		return resp
	}

	resp.Success = true
	resp.Message = "Success"
	resp.Data = master
	return resp
}

func CreateMasterApps(nama_aplikasi, deskripsi string, tgl_version time.Time, url, versi_aplikasi, image string) types.Response {
	var resp types.Response

	// password_hash, _ := middleware.HashPassword(password)

	result := database.DB.Exec(`INSERT INTO master_aplikasi 
								(nama_aplikasi
								, deskripsi
								, tgl_version
								, url
								, versi_aplikasi
								, image
								, status_data
								)
								VALUES
								(?, ?, ?, ?, ?, ?, ?)`, nama_aplikasi, deskripsi, tgl_version, url, versi_aplikasi, image, true)

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

func UpdateMasterApp(id_master_aplikasi, image, nama_aplikasi, deskripsi, versi_aplikasi string, tgl_version time.Time, url string) types.Response {
	var resp types.Response

	if image == "" {
		// Jika image kosong, update tanpa mengubah kolom image
		result := database.DB.Exec(`UPDATE master_aplikasi
									SET nama_aplikasi = ?, deskripsi = ?, versi_aplikasi = ?, tgl_version = ?, url = ?
									WHERE id_master_aplikasi = ?`,
			nama_aplikasi, deskripsi, versi_aplikasi, tgl_version, url, id_master_aplikasi)

		// Cek jika ada error atau tidak ada baris yang diperbarui
		if result.Error != nil {
			resp.Success = false
			resp.Message = "Gagal memperbarui aplikasi"
			resp.Data = nil
			return resp
		} else if result.RowsAffected == 0 {
			resp.Success = false
			resp.Message = "Aplikasi tidak ditemukan"
			resp.Data = nil
			return resp
		}
	} else {
		// Jika image tidak kosong, update termasuk kolom image
		result := database.DB.Exec(`UPDATE master_aplikasi
									SET image = ?, nama_aplikasi = ?, deskripsi = ?, versi_aplikasi = ?, tgl_version = ?, url = ?
									WHERE id_master_aplikasi = ?`,
			image, nama_aplikasi, deskripsi, versi_aplikasi, tgl_version, url, id_master_aplikasi)

		// Cek jika ada error atau tidak ada baris yang diperbarui
		if result.Error != nil {
			resp.Success = false
			resp.Message = "Gagal memperbarui aplikasi"
			resp.Data = nil
			return resp
		} else if result.RowsAffected == 0 {
			resp.Success = false
			resp.Message = "Aplikasi tidak ditemukan"
			resp.Data = nil
			return resp
		}
	}

	// Jika update berhasil
	resp.Success = true
	resp.Message = "Aplikasi berhasil diperbarui"
	resp.Data = nil

	return resp
}

func DeleteMasterApp(id_master_aplikasi string) types.Response {
	var resp types.Response

	// Eksekusi query untuk menghapus data master_aplikasi secara permanen
	result := database.DB.Exec(`DELETE FROM master_aplikasi
								WHERE id_master_aplikasi = ?`, id_master_aplikasi)

	// Jika terjadi error saat eksekusi query
	if result.Error != nil {
		resp.Success = false
		resp.Message = "Gagal Delete master aplikasi"
		resp.Data = nil
		return resp
	} else if result.RowsAffected == 0 {
		// Jika tidak ada data yang terpengaruh (id_master_aplikasi tidak ditemukan)
		resp.Success = false
		resp.Message = "Master aplikasi tidak ditemukan"
		resp.Data = nil
		return resp
	}

	// Jika berhasil menghapus data
	resp.Success = true
	resp.Message = "Success"
	resp.Data = nil

	return resp
}

// func UpdateMasterApps(nama_aplikasi, deskripsi, image string) types.Response {
// 	result := database.DB.Exec()
// }
