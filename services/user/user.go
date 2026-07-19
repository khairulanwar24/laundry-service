// Package services (user) berisi logika bisnis domain user:
// orkestrasi, hashing password, dan penyusunan response {success,message,data}.
package services

import (
	"context"

	"sso-service/common/response"
	"sso-service/domain/dto"
	middleware "sso-service/middlewares"
	"sso-service/models"
	"sso-service/repositories"
)

// UserService membungkus akses ke repository registry.
type UserService struct {
	repository repositories.IRepositoryRegistry
}

// IUserService adalah kontrak logika bisnis domain user.
type IUserService interface {
	GetUsers(order, filter string, limit, offset int) response.Response
	CreateUser(ctx context.Context, email, idPerson, jenisUser, namaLengkap, noHp, username, password, avatar string) response.Response
	GetUser(ctx context.Context, idUser string) response.Response
	UpdateUser(ctx context.Context, idUser, avatar, email, idPerson, jenisUser, namaLengkap, noHp, username string) response.Response
	UpdatePassword(ctx context.Context, idUser, password string) response.Response
	DeleteUser(ctx context.Context, idUser string) response.Response
	GetUsersDosen(ctx context.Context) response.Response
	GetDetailDosen(ctx context.Context, personID string) response.Response
	BulkCreateUsersMahasiswa(ctx context.Context, items []dto.BulkMahasiswaItem) response.Response
	GetDetailMahasiswa(ctx context.Context, idRegistrasi string) response.Response
	GetUsersMahasiswa(ctx context.Context, idProdi string) response.Response
	GetUsersMahasiswaData(ctx context.Context, idProdi, idAngkatan, order, filter string, limit, offset int) response.Response
	GenerateUserMahasiswa(ctx context.Context) response.Response
	GenerateUserDosen(ctx context.Context) response.Response
	GenerateUserTendik(ctx context.Context) response.Response
}

// NewUserService membuat instance UserService baru.
func NewUserService(repository repositories.IRepositoryRegistry) IUserService {
	return &UserService{repository: repository}
}

// GetUsers mengembalikan daftar user (datatable).
func (s *UserService) GetUsers(order, filter string, limit, offset int) response.Response {
	data := s.repository.GetUser().GetUsers(order, filter, limit, offset)
	return response.Response{Success: true, Message: "Success", Data: data}
}

// CreateUser membuat user baru: cek duplikat -> hash password -> insert.
func (s *UserService) CreateUser(ctx context.Context, email, idPerson, jenisUser, namaLengkap, noHp, username, password, avatar string) response.Response {
	count, err := s.repository.GetUser().CountByPerson(ctx, idPerson)
	if err != nil {
		return response.Response{Success: false, Message: err.Error()}
	}
	if count > 0 {
		return response.Response{Success: false, Message: "User dengan id_person sudah ada"}
	}

	passwordHash, _ := middleware.HashPassword(password)

	if err := s.repository.GetUser().Insert(ctx, email, idPerson, jenisUser, namaLengkap, noHp, username, passwordHash, avatar); err != nil {
		return response.Response{Success: false, Message: err.Error(), Data: err.Error()}
	}

	return response.Response{Success: true, Message: "Success"}
}

// GetUser mengambil detail satu user.
func (s *UserService) GetUser(ctx context.Context, idUser string) response.Response {
	user, rows, err := s.repository.GetUser().FindByID(ctx, idUser)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal Mendapatkan users"}
	} else if rows == 0 {
		return response.Response{Success: false, Message: "users tidak ada"}
	}
	return response.Response{Success: true, Message: "Success", Data: user}
}

// UpdateUser memperbarui data user (dengan / tanpa avatar).
func (s *UserService) UpdateUser(ctx context.Context, idUser, avatar, email, idPerson, jenisUser, namaLengkap, noHp, username string) response.Response {
	var rows int64
	var err error
	if avatar == "" {
		rows, err = s.repository.GetUser().UpdateWithoutAvatar(ctx, email, idPerson, jenisUser, namaLengkap, noHp, username, idUser)
	} else {
		rows, err = s.repository.GetUser().UpdateWithAvatar(ctx, avatar, email, idPerson, jenisUser, namaLengkap, noHp, username, idUser)
	}
	if err != nil {
		return response.Response{Success: false, Message: "Gagal Mendapatkan users"}
	} else if rows == 0 {
		return response.Response{Success: false, Message: "users tidak ada"}
	}
	return response.Response{Success: true, Message: "Success"}
}

// UpdatePassword memperbarui password user (menangani first_login).
func (s *UserService) UpdatePassword(ctx context.Context, idUser, password string) response.Response {
	user, err := s.repository.GetUser().FindFirstLogin(ctx, idUser)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal Mendapatkan users"}
	}

	passwordHash, _ := middleware.HashPassword(password)
	firstLogin := user[0]["first_login"].(bool)

	var rows int64
	if firstLogin {
		rows, err = s.repository.GetUser().UpdatePasswordFirstLogin(ctx, passwordHash, idUser)
	} else {
		rows, err = s.repository.GetUser().UpdatePasswordNormal(ctx, passwordHash, idUser)
	}

	if err != nil {
		return response.Response{Success: false, Message: "Gagal Update users"}
	} else if rows == 0 {
		return response.Response{Success: false, Message: "users tidak ada"}
	}
	return response.Response{Success: true, Message: "Success"}
}

// DeleteUser menghapus user.
func (s *UserService) DeleteUser(ctx context.Context, idUser string) response.Response {
	rows, err := s.repository.GetUser().Delete(ctx, idUser)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal Delete users"}
	} else if rows == 0 {
		return response.Response{Success: false, Message: "users tidak ada"}
	}
	return response.Response{Success: true, Message: "Success"}
}

// GetUsersDosen mengambil daftar dosen (digiclass).
func (s *UserService) GetUsersDosen(ctx context.Context) response.Response {
	dosen, _ := s.repository.GetUser().GetDosen(ctx)
	return response.Response{Success: true, Message: "Success", Data: dosen}
}

// GetDetailDosen mengambil detail satu dosen (digiclass).
func (s *UserService) GetDetailDosen(ctx context.Context, personID string) response.Response {
	dosen, _ := s.repository.GetUser().GetDetailDosen(ctx, personID)
	return response.Response{Success: true, Message: "Success", Data: dosen}
}

// BulkCreateUsersMahasiswa membuat banyak user mahasiswa sekaligus.
func (s *UserService) BulkCreateUsersMahasiswa(ctx context.Context, items []dto.BulkMahasiswaItem) response.Response {
	result := dto.BulkCreateResult{
		Total:       len(items),
		DetailGagal: []string{},
	}

	for _, item := range items {
		r := s.CreateUser(ctx, item.Email, item.IdPerson, "mahasiswa", item.NamaLengkap, item.NoHp, item.Username, item.Password, "")
		if r.Success {
			result.Berhasil++
		} else {
			result.Gagal++
			result.DetailGagal = append(result.DetailGagal, item.Username+": "+r.Message)
		}
	}

	return response.Response{Success: true, Message: "Bulk create selesai", Data: result}
}

// GetDetailMahasiswa mengambil detail satu mahasiswa (akademik).
func (s *UserService) GetDetailMahasiswa(ctx context.Context, idRegistrasi string) response.Response {
	mahasiswa, _ := s.repository.GetUser().GetDetailMahasiswa(ctx, idRegistrasi)
	return response.Response{Success: true, Message: "Success", Data: mahasiswa}
}

// GetUsersMahasiswa mengambil daftar mahasiswa per prodi (akademik).
func (s *UserService) GetUsersMahasiswa(ctx context.Context, idProdi string) response.Response {
	mahasiswa, _ := s.repository.GetUser().GetMahasiswa(ctx, idProdi)
	return response.Response{Success: true, Message: "Success", Data: mahasiswa}
}

// GetUsersMahasiswaData mengambil data mahasiswa dengan paginasi & filter (akademik).
func (s *UserService) GetUsersMahasiswaData(ctx context.Context, idProdi, idAngkatan, order, filter string, limit, offset int) response.Response {
	data, total, totalFiltered, _ := s.repository.GetUser().GetMahasiswaData(ctx, idProdi, idAngkatan, order, filter, limit, offset)
	return response.Response{
		Success: true,
		Message: "Success",
		Data: map[string]interface{}{
			"data":            data,
			"recordsTotal":    total,
			"recordsFiltered": totalFiltered,
		},
	}
}

// GenerateUserMahasiswa membuat user untuk seluruh mahasiswa aktif sumber (akademik).
func (s *UserService) GenerateUserMahasiswa(ctx context.Context) response.Response {
	mahasiswas, err := s.repository.GetUser().GetMahasiswaSourceForGenerate(ctx)
	if err != nil {
		return response.Response{Success: false, Message: err.Error()}
	}

	for _, m := range mahasiswas {
		_ = s.CreateUser(ctx, m.Email, m.IdPerson, m.JenisUser, m.Nama, m.NoHp, m.Username, m.Password, m.Avatar)
	}

	return response.Response{Success: true, Message: "Success generate user mahasiswa"}
}

// GenerateUserDosen menyetel password dosen & memberikan group akses.
func (s *UserService) GenerateUserDosen(ctx context.Context) response.Response {
	dosen, err := s.repository.GetUser().GetPegawaiNeedingPassword(ctx, "Dosen")
	if err != nil {
		return response.Response{Success: false, Message: err.Error()}
	}

	for _, m := range dosen {
		_ = s.UpdatePassword(ctx, m.IDUser, m.Password)
		// Bridge lintas-domain: mst_group_akses belum dimigrasi ke pola berlapis.
		_ = models.CreateGroupAksesUserApps(m.IDUser, "acce4365-6835-441e-88d2-56539fb0e823", "f66ff11d-a59c-4b17-b039-3a5b3ed03508", "true")
	}

	return response.Response{Success: true, Message: "Success generate user dosen"}
}

// GenerateUserTendik menyetel password tenaga pendidik & memberikan group akses.
func (s *UserService) GenerateUserTendik(ctx context.Context) response.Response {
	tendik, err := s.repository.GetUser().GetPegawaiNeedingPassword(ctx, "Tenaga Pendidik")
	if err != nil {
		return response.Response{Success: false, Message: err.Error()}
	}

	for _, m := range tendik {
		_ = s.UpdatePassword(ctx, m.IDUser, m.Password)
		// Bridge lintas-domain: mst_group_akses belum dimigrasi ke pola berlapis.
		_ = models.CreateGroupAksesUserApps(m.IDUser, "acce4365-6835-441e-88d2-56539fb0e823", "0daae1e8-ffb2-4c9e-bf32-6dbc5a8da847", "true")
	}

	return response.Response{Success: true, Message: "Success generate user Tendik"}
}
