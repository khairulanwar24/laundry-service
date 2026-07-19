package dto

// GetAngkatanParams adalah DTO untuk parameter rute get_angkatan.
type GetAngkatanParams struct {
	IDProdi string `json:"id_prodi" form:"id_prodi" validate:"required,uuid4"`
}
