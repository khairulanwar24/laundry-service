package dto

// GetData adalah DTO pagination generik (mirror dari types.GetData).
type GetData struct {
	Limit  int    `json:"limit" form:"limit" validate:"required,numeric,oneof=10 25 50 100" default:"10"`
	Offset int    `json:"offset" form:"offset" validate:"numeric"`
	Order  string `json:"order" form:"order" validate:""`
	Filter string `json:"filter" form:"filter" validate:""`
	Params string `json:"params" form:"params" validate:""`
}

type GetMasterAppByIdParams struct {
	Id_master_aplikasi string `json:"id_master_aplikasi" form:"id_master_aplikasi" validate:"required,uuid4"`
}

type CreateMasterAppForm struct {
	Nama_Aplikasi  string `json:"nama_aplikasi" form:"nama_aplikasi" validate:"required"`
	Deskripsi      string `json:"deskripsi" form:"deskripsi" validate:""`
	Versi_Aplikasi string `json:"versi_aplikasi" form:"versi_aplikasi" validate:"required"`
	Tgl_Version    string `json:"tgl_version" form:"tgl_version" validate:"required,datetime=2006-01-02"`
	Url            string `json:"url" form:"url" validate:"required,url"`
}

type UpdateMasterAppsParams struct {
	Id_master_aplikasi string `json:"id_master_aplikasi" form:"id_master_aplikasi" validate:"required,uuid4"`
}

type UpdateMasterAppsForm struct {
	Nama_Aplikasi  string `json:"nama_aplikasi" form:"nama_aplikasi" validate:"required"`
	Deskripsi      string `json:"deskripsi" form:"deskripsi" validate:"required"`
	Versi_Aplikasi string `json:"versi_aplikasi" form:"versi_aplikasi" validate:"required"`
	Tgl_Version    string `json:"tgl_version" form:"tgl_version" validate:"required"`
	URL            string `json:"url" form:"url" validate:"required"`
}

type DeleteMasterAppParams struct {
	Id_Master_Aplikasi string `json:"id_master_aplikasi" validate:"required,uuid4"`
}
