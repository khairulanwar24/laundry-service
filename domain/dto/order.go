package dto

// CreateOrderItemForm adalah satu baris item pesanan dari request client.
type CreateOrderItemForm struct {
	IDServiceVariant string  `json:"service_variant_id" form:"service_variant_id" validate:"required,uuid4"`
	Qty              float64 `json:"qty" form:"qty" validate:"required,gt=0"`
	Catatan          string  `json:"catatan" form:"catatan" validate:"omitempty,max=1000"`
}

// CreateOrderForm adalah form membuat pesanan baru.
type CreateOrderForm struct {
	IDCustomer string                `json:"customer_id" form:"customer_id" validate:"required,uuid4"`
	Items      []CreateOrderItemForm `json:"items" form:"items" validate:"required,min=1,dive"`
	IDPerfume  string                `json:"perfume_id" form:"perfume_id" validate:"omitempty,uuid4"`
	IDDiscount string                `json:"discount_id" form:"discount_id" validate:"omitempty,uuid4"`
	Catatan    string                `json:"catatan" form:"catatan" validate:"omitempty,max=2000"`
}

// ChangeStatusForm dipakai untuk mengubah status pesanan.
type ChangeStatusForm struct {
	To      string `json:"to" form:"to" validate:"required,oneof=ANTRIAN PROSES SIAP_DIAMBIL SELESAI BATAL"`
	Catatan string `json:"catatan" form:"catatan" validate:"omitempty,max=1000"`
}

// PayOrderForm dipakai untuk mencatat pembayaran pesanan.
type PayOrderForm struct {
	IDPaymentMethod string  `json:"payment_method_id" form:"payment_method_id" validate:"required,uuid4"`
	Jumlah          float64 `json:"amount" form:"amount" validate:"required,gt=0"`
	NoReferensi     string  `json:"ref_no" form:"ref_no" validate:"omitempty,max=255"`
}

// PickupForm dipakai untuk menandai pesanan diambil pelanggan.
type PickupForm struct {
	Catatan string `json:"catatan" form:"catatan" validate:"omitempty,max=1000"`
}

// OrderItemInput adalah bentuk item pesanan yang sudah dihitung/di-snapshot
// oleh service layer, dikirim ke repository untuk disimpan.
type OrderItemInput struct {
	IDServiceVariant       string
	Satuan                 string
	Qty                    float64
	HargaPerSatuanSnapshot float64
	TotalHarga             float64
	Catatan                string
}
