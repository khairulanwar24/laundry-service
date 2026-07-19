package dto

type GetMstMenuParams struct {
	Id_Master_Aplikasi string `json:"id_master_aplikasi" form:"id_master_aplikasi" validate:"required,uuid4"`
}

type GetMstMenuForm struct {
	Limit  int    `json:"limit" form:"limit" validate:"required,numeric,oneof=10 25 50 100"  default:"10"`
	Offset int    `json:"offset"  form:"offset" validate:"numeric"`
	Order  string `json:"order"  form:"order" validate:""`
	Filter string `json:"filter"  form:"filter" validate:""`
}

type GetDetailMstMenuParams struct {
	Id_Master_Menu string `json:"id_master_menu" form:"id_master_menu" validate:"required,uuid4"`
}

type GetMstMenuModulParams struct {
	Id_Master_Menu string `json:"id_master_menu" form:"id_master_menu" validate:"required,uuid4"`
}

type GetMstMenuModulForm struct {
	Limit  int    `json:"limit" form:"limit" validate:"required,numeric,oneof=10 25 50 100"  default:"10"`
	Offset int    `json:"offset"  form:"offset" validate:"numeric"`
	Order  string `json:"order"  form:"order" validate:""`
	Filter string `json:"filter"  form:"filter" validate:""`
}

type GetDetailMstMenuModulParams struct {
	Id_Master_Modul string `json:"id_master_modul" form:"id_master_modul" validate:"required,uuid4"`
}

type CreateMstMenuForm struct {
	Id_Master_Aplikasi string `json:"id_master_aplikasi" form:"id_master_aplikasi" validate:"required,uuid4"`
	Nama_Menu          string `json:"nama_menu" form:"nama_menu" validate:"required"`
	Deskripsi          string `json:"deskripsi" form:"deskripsi" validate:"required"`
	Order              string `json:"order" form:"order" validate:"required"`
	Icon               string `json:"icon" form:"icon" validate:"required"`
}

type CreateMstMenuModulForm struct {
	Id_Master_Aplikasi string `json:"id_master_aplikasi" form:"id_master_aplikasi" validate:"required,uuid4"`
	Id_Master_Menu     string `json:"id_master_menu" form:"id_master_menu" validate:"required,uuid4"`
	Nama_Modul         string `json:"nama_modul" form:"nama_modul" validate:"required"`
	Path               string `json:"path" form:"path" validate:"required"`
	Deskripsi          string `json:"deskripsi" form:"deskripsi" validate:"required"`
	Order              string `json:"order" form:"order" validate:"required"`
	Icon               string `json:"icon" form:"icon" validate:"required"`
}

type UpdateMstMenuParams struct {
	Id_Master_Menu string `json:"id_master_menu" form:"id_master_menu" validate:"required,uuid4"`
}

type UpdateMstMenuForm struct {
	Nama_Menu string `json:"nama_menu" form:"nama_menu" validate:"required"`
	Deskripsi string `json:"deskripsi" form:"deskripsi" validate:"required"`
	Order     string `json:"order" form:"order" validate:"required,numeric"`
	Icon      string `json:"icon" form:"icon" validate:"required"`
}

type UpdateMstMenuModulParams struct {
	Id_Master_Modul string `json:"id_master_modul" form:"id_master_modul" validate:"required,uuid4"`
}

type UpdateMstMenuModulForm struct {
	Id_Master_Aplikasi string `json:"id_master_aplikasi" form:"id_master_aplikasi" validate:"required,uuid4"`
	Id_Master_Menu     string `json:"id_master_menu" form:"id_master_menu" validate:"required,uuid4"`
	Nama_Modul         string `json:"nama_modul" form:"nama_modul" validate:"required"`
	Path               string `json:"path" form:"path" validate:"required"`
	Deskripsi          string `json:"deskripsi" form:"deskripsi" validate:"required"`
	Order              string `json:"order" form:"order" validate:"required"`
	Icon               string `json:"icon" form:"icon" validate:"required"`
}

type DeleteMstMenuParams struct {
	Id_Master_Menu string `json:"id_master_menu" form:"id_master_menu" validate:"required,uuid4"`
}

type DeleteMstMenuModulParams struct {
	Id_Master_Modul string `json:"id_master_modul" form:"id_master_modul" validate:"required,uuid4"`
}
