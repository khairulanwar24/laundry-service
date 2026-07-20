// Package services (outlet) berisi logika bisnis domain outlet: kelola outlet,
// keanggotaan, dan permission karyawan.
package services

import (
	"context"
	"encoding/json"

	"laundry-service/common/response"
	"laundry-service/domain/dto"
	middleware "laundry-service/middlewares"
	"laundry-service/repositories"
)

var permissionKeys = []string{
	"create_order", "cancel_order", "create_expense", "manage_services",
	"manage_customers", "manage_employees", "view_revenue", "view_report_tx",
	"view_report_finance", "view_report_customer",
}

func defaultPermissionsForRole(role string) map[string]bool {
	if role == "owner" {
		perms := make(map[string]bool, len(permissionKeys))
		for _, k := range permissionKeys {
			perms[k] = true
		}
		return perms
	}
	// karyawan
	return map[string]bool{
		"create_order":         true,
		"cancel_order":         true,
		"create_expense":       true,
		"manage_customers":     true,
		"manage_services":      false,
		"manage_employees":     false,
		"view_revenue":         false,
		"view_report_tx":       false,
		"view_report_finance":  false,
		"view_report_customer": false,
	}
}

func mergePermissions(base, override map[string]bool) map[string]bool {
	result := make(map[string]bool, len(base))
	for k, v := range base {
		result[k] = v
	}
	for _, k := range permissionKeys {
		if v, ok := override[k]; ok {
			result[k] = v
		}
	}
	return result
}

// OutletService membungkus akses ke repository registry.
type OutletService struct {
	repository repositories.IRepositoryRegistry
}

// IOutletService adalah kontrak logika bisnis domain outlet.
type IOutletService interface {
	ListMyOutlets(ctx context.Context, idUser string) response.Response
	CreateOutlet(ctx context.Context, idUser string, form dto.CreateOutletForm) response.Response
	GetOutlet(ctx context.Context, idOutlet string) response.Response
	UpdateOutlet(ctx context.Context, idOutlet string, form dto.UpdateOutletForm) response.Response
	ListStaff(ctx context.Context, idOutlet string) response.Response
	InviteEmployee(ctx context.Context, idOutlet string, form dto.InviteEmployeeForm) response.Response
	UpdateMember(ctx context.Context, idOutlet, idMember string, form dto.UpdateMemberForm) response.Response
	RemoveMember(ctx context.Context, idOutlet, idMember string) response.Response
}

// NewOutletService membuat instance OutletService baru.
func NewOutletService(repository repositories.IRepositoryRegistry) IOutletService {
	return &OutletService{repository: repository}
}

func (s *OutletService) ListMyOutlets(ctx context.Context, idUser string) response.Response {
	data, err := s.repository.GetOutlet().ListOutletsForUser(ctx, idUser)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil daftar outlet: " + err.Error()}
	}
	return response.Response{Success: true, Message: "sukses", Data: data}
}

func (s *OutletService) CreateOutlet(ctx context.Context, idUser string, form dto.CreateOutletForm) response.Response {
	repo := s.repository.GetOutlet()

	idOutlet, err := repo.CreateOutlet(ctx, idUser, form.Nama, form.Alamat, form.Telepon, form.LogoPath)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal membuat outlet: " + err.Error()}
	}

	permsJSON, _ := json.Marshal(defaultPermissionsForRole("owner"))
	if _, err := repo.AssignUser(ctx, idOutlet, idUser, "owner", string(permsJSON)); err != nil {
		return response.Response{Success: false, Message: "Outlet dibuat tapi gagal menetapkan owner: " + err.Error()}
	}

	return response.Response{Success: true, Message: "Outlet berhasil dibuat", Data: map[string]interface{}{"id_outlet": idOutlet}}
}

func (s *OutletService) GetOutlet(ctx context.Context, idOutlet string) response.Response {
	data, err := s.repository.GetOutlet().FindOutletByID(ctx, idOutlet)
	if err != nil {
		return response.Response{Success: false, Message: "Outlet tidak ditemukan"}
	}
	return response.Response{Success: true, Message: "sukses", Data: data}
}

func (s *OutletService) UpdateOutlet(ctx context.Context, idOutlet string, form dto.UpdateOutletForm) response.Response {
	if err := s.repository.GetOutlet().UpdateOutlet(ctx, idOutlet, form.Nama, form.Alamat, form.Telepon, form.LogoPath); err != nil {
		return response.Response{Success: false, Message: "Gagal memperbarui outlet: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Outlet berhasil diperbarui"}
}

func (s *OutletService) ListStaff(ctx context.Context, idOutlet string) response.Response {
	data, err := s.repository.GetOutlet().ListStaff(ctx, idOutlet)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil daftar staff: " + err.Error()}
	}
	return response.Response{Success: true, Message: "sukses", Data: data}
}

// InviteEmployee menambahkan karyawan/owner baru ke outlet. Jika email/telepon
// sudah terdaftar, akun yang sudah ada dipakai (tidak dibuat ulang).
func (s *OutletService) InviteEmployee(ctx context.Context, idOutlet string, form dto.InviteEmployeeForm) response.Response {
	accountRepo := s.repository.GetAccount()

	var idUser string
	existing, err := accountRepo.FindUserByEmail(ctx, form.Email)
	if err == nil {
		idUser, _ = existing["id_user"].(string)
	} else {
		if existsPhone, _ := accountRepo.ExistsPhone(ctx, form.Telepon); existsPhone {
			return response.Response{Success: false, Message: "Telepon sudah terdaftar dengan akun lain"}
		}
		passwordHash, err := middleware.HashPassword(form.Password)
		if err != nil {
			return response.Response{Success: false, Message: "Gagal memproses password"}
		}
		idUser, err = accountRepo.CreateUser(ctx, form.Nama, form.Email, form.Telepon, passwordHash, form.Alamat)
		if err != nil {
			return response.Response{Success: false, Message: "Gagal membuat akun karyawan: " + err.Error()}
		}
	}

	perms := mergePermissions(defaultPermissionsForRole(form.Role), form.Permissions)
	permsJSON, _ := json.Marshal(perms)

	if _, err := s.repository.GetOutlet().AssignUser(ctx, idOutlet, idUser, form.Role, string(permsJSON)); err != nil {
		return response.Response{Success: false, Message: "Gagal menambahkan karyawan ke outlet: " + err.Error()}
	}

	return response.Response{Success: true, Message: "Karyawan berhasil ditambahkan", Data: map[string]interface{}{"id_user": idUser}}
}

func (s *OutletService) UpdateMember(ctx context.Context, idOutlet, idMember string, form dto.UpdateMemberForm) response.Response {
	repo := s.repository.GetOutlet()

	member, err := repo.FindMemberByID(ctx, idOutlet, idMember)
	if err != nil {
		return response.Response{Success: false, Message: "Anggota tidak ditemukan"}
	}

	currentRole, _ := member["role"].(string)
	role := currentRole
	if form.Role != "" {
		role = form.Role
	}

	var currentPerms map[string]bool
	if raw, ok := member["permissions_json"]; ok && raw != nil {
		_ = json.Unmarshal([]byte(toJSONString(raw)), &currentPerms)
	}
	base := defaultPermissionsForRole(role)
	if form.Role == "" && currentPerms != nil {
		base = currentPerms
	}
	perms := mergePermissions(base, form.Permissions)
	permsJSON, _ := json.Marshal(perms)

	idUserOutlet, _ := member["id_user_outlet"].(string)
	if err := repo.UpdateMember(ctx, idUserOutlet, role, string(permsJSON)); err != nil {
		return response.Response{Success: false, Message: "Gagal memperbarui anggota: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Anggota berhasil diperbarui"}
}

func (s *OutletService) RemoveMember(ctx context.Context, idOutlet, idMember string) response.Response {
	repo := s.repository.GetOutlet()

	member, err := repo.FindMemberByID(ctx, idOutlet, idMember)
	if err != nil {
		return response.Response{Success: false, Message: "Anggota tidak ditemukan"}
	}

	if role, _ := member["role"].(string); role == "owner" {
		count, err := repo.CountActiveOwners(ctx, idOutlet)
		if err != nil {
			return response.Response{Success: false, Message: "Gagal memeriksa jumlah owner: " + err.Error()}
		}
		if count <= 1 {
			return response.Response{Success: false, Message: "Tidak dapat menghapus owner terakhir outlet"}
		}
	}

	idUserOutlet, _ := member["id_user_outlet"].(string)
	if err := repo.DeactivateMember(ctx, idUserOutlet); err != nil {
		return response.Response{Success: false, Message: "Gagal menghapus anggota: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Anggota berhasil dihapus dari outlet"}
}

// toJSONString menormalkan nilai kolom jsonb yang dibaca lewat map[string]interface{}
// (driver Postgres bisa mengembalikannya sebagai []byte maupun string).
func toJSONString(v interface{}) string {
	switch val := v.(type) {
	case []byte:
		return string(val)
	case string:
		return val
	default:
		return "{}"
	}
}
