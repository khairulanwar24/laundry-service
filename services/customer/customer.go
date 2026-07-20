// Package services (customer) berisi logika bisnis domain pelanggan outlet.
package services

import (
	"context"

	"laundry-service/common/response"
	"laundry-service/domain/dto"
	"laundry-service/repositories"
)

// CustomerService membungkus akses ke repository registry.
type CustomerService struct {
	repository repositories.IRepositoryRegistry
}

// ICustomerService adalah kontrak logika bisnis domain pelanggan.
type ICustomerService interface {
	List(ctx context.Context, idOutlet, search string, isActive *bool, page, perPage int) response.Response
	Create(ctx context.Context, idOutlet string, form dto.CustomerForm) response.Response
	Get(ctx context.Context, idOutlet, idCustomer string) response.Response
	Update(ctx context.Context, idOutlet, idCustomer string, form dto.CustomerForm) response.Response
	Delete(ctx context.Context, idOutlet, idCustomer string) response.Response
	ListOrders(ctx context.Context, idOutlet, idCustomer, status, fromDate, toDate string, page, perPage int) response.Response
}

// NewCustomerService membuat instance CustomerService baru.
func NewCustomerService(repository repositories.IRepositoryRegistry) ICustomerService {
	return &CustomerService{repository: repository}
}

func normalizePage(page, perPage int) (int, int) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 15
	}
	return page, perPage
}

func (s *CustomerService) List(ctx context.Context, idOutlet, search string, isActive *bool, page, perPage int) response.Response {
	page, perPage = normalizePage(page, perPage)
	repo := s.repository.GetCustomer()

	data, err := repo.List(ctx, idOutlet, search, isActive, perPage, (page-1)*perPage)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil daftar pelanggan: " + err.Error()}
	}
	total, err := repo.Count(ctx, idOutlet, search, isActive)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal menghitung total pelanggan: " + err.Error()}
	}

	return response.Response{
		Success: true,
		Message: "sukses",
		Data: map[string]interface{}{
			"data": data,
			"meta": map[string]interface{}{
				"current_page": page,
				"per_page":     perPage,
				"total":        total,
			},
		},
	}
}

func (s *CustomerService) Create(ctx context.Context, idOutlet string, form dto.CustomerForm) response.Response {
	repo := s.repository.GetCustomer()

	if existsPhone, _ := repo.ExistsPhoneInOutlet(ctx, idOutlet, form.Telepon, ""); existsPhone {
		return response.Response{Success: false, Message: "Telepon sudah dipakai pelanggan lain di outlet ini"}
	}
	if existsEmail, _ := repo.ExistsEmailInOutlet(ctx, idOutlet, form.Email, ""); existsEmail {
		return response.Response{Success: false, Message: "Email sudah dipakai pelanggan lain di outlet ini"}
	}

	isActive := true
	if form.IsActive != nil {
		isActive = *form.IsActive
	}

	id, err := repo.Create(ctx, idOutlet, form.Nama, form.Telepon, form.Email, form.Alamat, isActive)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal membuat pelanggan: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Pelanggan berhasil dibuat", Data: map[string]interface{}{"id_customer": id}}
}

func (s *CustomerService) Get(ctx context.Context, idOutlet, idCustomer string) response.Response {
	data, err := s.repository.GetCustomer().FindByID(ctx, idOutlet, idCustomer)
	if err != nil {
		return response.Response{Success: false, Message: "Pelanggan tidak ditemukan"}
	}
	return response.Response{Success: true, Message: "sukses", Data: data}
}

func (s *CustomerService) Update(ctx context.Context, idOutlet, idCustomer string, form dto.CustomerForm) response.Response {
	repo := s.repository.GetCustomer()

	if _, err := repo.FindByID(ctx, idOutlet, idCustomer); err != nil {
		return response.Response{Success: false, Message: "Pelanggan tidak ditemukan"}
	}
	if existsPhone, _ := repo.ExistsPhoneInOutlet(ctx, idOutlet, form.Telepon, idCustomer); existsPhone {
		return response.Response{Success: false, Message: "Telepon sudah dipakai pelanggan lain di outlet ini"}
	}
	if existsEmail, _ := repo.ExistsEmailInOutlet(ctx, idOutlet, form.Email, idCustomer); existsEmail {
		return response.Response{Success: false, Message: "Email sudah dipakai pelanggan lain di outlet ini"}
	}

	isActive := true
	if form.IsActive != nil {
		isActive = *form.IsActive
	}

	if err := repo.Update(ctx, idCustomer, form.Nama, form.Telepon, form.Email, form.Alamat, isActive); err != nil {
		return response.Response{Success: false, Message: "Gagal memperbarui pelanggan: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Pelanggan berhasil diperbarui"}
}

func (s *CustomerService) Delete(ctx context.Context, idOutlet, idCustomer string) response.Response {
	repo := s.repository.GetCustomer()

	if _, err := repo.FindByID(ctx, idOutlet, idCustomer); err != nil {
		return response.Response{Success: false, Message: "Pelanggan tidak ditemukan"}
	}

	count, err := repo.CountOrders(ctx, idCustomer)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal memeriksa riwayat pesanan: " + err.Error()}
	}
	if count > 0 {
		return response.Response{Success: false, Message: "Pelanggan memiliki riwayat pesanan, nonaktifkan saja alih-alih menghapus"}
	}

	if err := repo.Delete(ctx, idCustomer); err != nil {
		return response.Response{Success: false, Message: "Gagal menghapus pelanggan: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Pelanggan berhasil dihapus"}
}

func (s *CustomerService) ListOrders(ctx context.Context, idOutlet, idCustomer, status, fromDate, toDate string, page, perPage int) response.Response {
	repo := s.repository.GetCustomer()

	if _, err := repo.FindByID(ctx, idOutlet, idCustomer); err != nil {
		return response.Response{Success: false, Message: "Pelanggan tidak ditemukan"}
	}

	page, perPage = normalizePage(page, perPage)
	data, err := repo.ListOrders(ctx, idCustomer, status, fromDate, toDate, perPage, (page-1)*perPage)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil riwayat pesanan: " + err.Error()}
	}
	return response.Response{Success: true, Message: "sukses", Data: data}
}
