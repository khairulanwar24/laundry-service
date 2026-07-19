package dto

// ==== master group akses ====

type CreateMstGroupAksesForm struct {
	Id_Master_Aplikasi string `json:"id_master_aplikasi" form:"id_master_aplikasi" validate:"required,uuid4"`
	Nama_Group         string `json:"nama_group" form:"nama_group" validate:"required"`
	Deskripsi          string `json:"deskripsi" form:"deskripsi" validate:"required"`
}

type GetMstGroupAksesParams struct {
	Id_Master_Aplikasi string `json:"id_master_aplikasi" form:"id_master_aplikasi" validate:"required,uuid4"`
}

type GetMstGroupAksesForm struct {
	Limit  int    `json:"limit" form:"limit" validate:"required,numeric,oneof=10 25 50 100"  default:"10"`
	Offset int    `json:"offset"  form:"offset" validate:"numeric"`
	Order  string `json:"order"  form:"order" validate:""`
	Filter string `json:"filter"  form:"filter" validate:""`
}

type GetMstGroupAksesModulParams struct {
	Id_Master_Aplikasi string `json:"id_master_aplikasi" form:"id_master_aplikasi" validate:"required,uuid4"`
	Id_Master_Group    string `json:"id_master_group" form:"id_master_aplikasi" validate:"required,uuid4"`
}

type GetMstGroupAksesModulForm struct {
	Limit  int    `json:"limit" form:"limit" validate:"required,numeric,oneof=10 25 50 100"  default:"10"`
	Offset int    `json:"offset"  form:"offset" validate:"numeric"`
	Order  string `json:"order"  form:"order" validate:""`
	Filter string `json:"filter"  form:"filter" validate:""`
}

type UpdateMstGroupAksesParams struct {
	Id_Master_Group string `json:"id_master_group" form:"id_master_group" validate:"required,uuid4"`
}

type UpdateMstGroupAksesForm struct {
	Nama_Group string `json:"nama_group" form:"nama_group" validate:"required"`
	Deskripsi  string `json:"deskripsi" form:"deskripsi" validate:"required"`
}

type DeleteMstGroupAksesParams struct {
	Id_Master_Group string `json:"id_master_group" form:"id_master_group" validate:"required,uuid4"`
}

type GetDetailMstGroupAksesParams struct {
	Id_Master_Group string `json:"id_master_group" form:"id_master_group" validate:"required,uuid4"`
}

// ==== group akses ====

type GetGroupAksesParams struct {
	Id_Master_Group string `json:"id_master_group" form:"id_master_group" validate:"required,uuid4"`
}

type GetGroupAksesForm struct {
	Limit  int    `json:"limit" form:"limit" validate:"required,numeric,oneof=10 25 50 100"  default:"10"`
	Offset int    `json:"offset"  form:"offset" validate:"numeric"`
	Order  string `json:"order"  form:"order" validate:""`
	Filter string `json:"filter"  form:"filter" validate:""`
}

type CreateGroupAksesForm struct {
	Id_Master_Aplikasi string `json:"id_master_aplikasi" form:"id_master_aplikasi" validate:"required,uuid4"`
	Id_Master_Group    string `json:"id_master_group" form:"id_master_group" validate:"required,uuid4"`
	Id_Master_Modul    string `json:"id_master_modul" form:"id_master_modul" validate:"required,uuid4"`
}

type DeleteGroupAksesParams struct {
	Id_Group_Akses string `json:"id_group_akses" form:"id_group_akses" validate:"required,uuid4"`
}

type GetGroupAksesUserMenuParams struct {
	Id_User            string `json:"id_user" form:"id_user" validate:"required,uuid4"`
	Id_Master_Aplikasi string `json:"id_master_aplikasi" form:"id_group_akses" validate:"required,uuid4"`
}

type GetGroupAksesUserAppsParams struct {
	Id_User string `json:"id_user" form:"id_user" validate:"required,uuid4"`
}

type CreateGroupAksesUserAppsForm struct {
	Id_User            string `json:"id_user" form:"id_user" validate:"required,uuid4"`
	Id_Master_Aplikasi string `json:"id_master_aplikasi" form:"id_master_aplikasi" validate:"required,uuid4"`
	Id_Master_Group    string `json:"id_master_group" form:"id_master_group" validate:"required,uuid4"`
	Status_Data        string `json:"status_data" form:"status_data" validate:"required,oneof='true' 'false' "`
}

type BulkGroupAksesUserAppsForm struct {
	IDUsers          []string `json:"id_users" validate:"required,min=1,dive,uuid"`
	IDMasterAplikasi string   `json:"id_master_aplikasi" validate:"required,uuid"`
	IDMasterGroup    string   `json:"id_master_group" validate:"required,uuid"`
	StatusData       bool     `json:"status_data" validate:"required"`
}

type UpdateGroupAksesUserAppsParams struct {
	Id_Trans_User_Group string `json:"id_trans_user_group" form:"id_trans_user_group" validate:"required,uuid4"`
}

type UpdateGroupAksesUserAppsForm struct {
	Status_Data string `json:"status_data" form:"status_data" validate:"required,oneof='true' 'false' "`
}

type DeleteGroupAksesUserAppsParams struct {
	Id_Trans_User_Group string `json:"id_trans_user_group" form:"id_trans_user_group" validate:"required,uuid4"`
}
