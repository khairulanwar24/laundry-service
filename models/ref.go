package models

import (
	"sso-service/database"
	"sso-service/types"
)

func GetMasterProdi() types.Response {

	var prodi []map[string]any
	var resp types.Response
	database.DBAkademik.Raw(`select * from master_prodi`).Scan(&prodi)
	resp.Data = prodi
	resp.Success = true
	resp.Message = "Success"

	return resp
}

func GetAngkatan(id_prodi string) types.Response {
	var angkatan []map[string]any
	var resp types.Response
	database.DBAkademik.Raw(
		`SELECT id_angkatan_mahasiswa, nama_angkatan_mahasiswa
		 FROM master_angkatan_mahasiswa
		 WHERE id_prodi = ? AND status_data = true
		 ORDER BY nama_angkatan_mahasiswa DESC`,
		id_prodi,
	).Scan(&angkatan)
	resp.Data = angkatan
	resp.Success = true
	resp.Message = "Success"
	return resp
}
