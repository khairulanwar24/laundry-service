-- ============================================================
-- SCHEMA LAUNDRY
-- Sistem manajemen laundry multi-outlet
-- ============================================================

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE SCHEMA IF NOT EXISTS laundry;

-- ============================================================
-- AKUN & OUTLET
-- ============================================================

CREATE TABLE IF NOT EXISTS laundry.users (
    id_user UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nama VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE,
    telepon VARCHAR(20) UNIQUE,
    password VARCHAR(255) NOT NULL,
    alamat TEXT,
    is_active BOOLEAN DEFAULT true,
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP
);

CREATE TABLE IF NOT EXISTS laundry.password_resets (
    id_password_reset UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL,
    otp VARCHAR(6) NOT NULL,
    token_hash VARCHAR(255) NOT NULL,
    otp_sent_at TIMESTAMP NOT NULL,
    otp_verified BOOLEAN DEFAULT false,
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_password_resets_email ON laundry.password_resets(email);

CREATE TABLE IF NOT EXISTS laundry.outlets (
    id_outlet UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_owner_user UUID NOT NULL,
    nama VARCHAR(120) NOT NULL,
    logo_path VARCHAR(255),
    alamat TEXT,
    telepon VARCHAR(30),
    is_active BOOLEAN DEFAULT true,
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP,
    CONSTRAINT fk_outlets_owner FOREIGN KEY (id_owner_user) REFERENCES laundry.users(id_user) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_outlets_owner_active ON laundry.outlets(id_owner_user, is_active);

CREATE TABLE IF NOT EXISTS laundry.user_outlets (
    id_user_outlet UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_user UUID NOT NULL,
    id_outlet UUID NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('owner', 'karyawan')),
    permissions_json JSONB DEFAULT '{}',
    is_active BOOLEAN DEFAULT true,
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP,
    CONSTRAINT fk_user_outlets_user FOREIGN KEY (id_user) REFERENCES laundry.users(id_user) ON DELETE CASCADE,
    CONSTRAINT fk_user_outlets_outlet FOREIGN KEY (id_outlet) REFERENCES laundry.outlets(id_outlet) ON DELETE CASCADE,
    CONSTRAINT uq_user_outlets_user_outlet UNIQUE (id_user, id_outlet)
);
CREATE INDEX IF NOT EXISTS idx_user_outlets_role_active ON laundry.user_outlets(role, is_active);
CREATE INDEX IF NOT EXISTS idx_user_outlets_outlet ON laundry.user_outlets(id_outlet);

-- ============================================================
-- PAYMENT METHOD
-- ============================================================

CREATE TABLE IF NOT EXISTS laundry.payment_methods (
    id_payment_method UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_outlet UUID NOT NULL,
    kategori VARCHAR(20) NOT NULL CHECK (kategori IN ('cash', 'transfer', 'e_wallet')),
    nama VARCHAR(120) NOT NULL,
    logo VARCHAR(255),
    nama_pemilik VARCHAR(120),
    tags JSONB DEFAULT '[]',
    is_active BOOLEAN DEFAULT true,
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP,
    CONSTRAINT fk_payment_methods_outlet FOREIGN KEY (id_outlet) REFERENCES laundry.outlets(id_outlet) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_payment_methods_outlet_active ON laundry.payment_methods(id_outlet, kategori, is_active);

-- ============================================================
-- KATALOG: SERVICE, VARIAN, PARFUM, DISKON
-- ============================================================

CREATE TABLE IF NOT EXISTS laundry.services (
    id_service UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_outlet UUID NOT NULL,
    nama VARCHAR(120) NOT NULL,
    prioritas SMALLINT NOT NULL DEFAULT 50,
    langkah_proses JSONB NOT NULL DEFAULT '["cuci","kering","setrika"]',
    is_active BOOLEAN DEFAULT true,
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP,
    CONSTRAINT fk_services_outlet FOREIGN KEY (id_outlet) REFERENCES laundry.outlets(id_outlet) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_services_outlet_active ON laundry.services(id_outlet, is_active);
CREATE INDEX IF NOT EXISTS idx_services_outlet_priority ON laundry.services(id_outlet, prioritas DESC);

CREATE TABLE IF NOT EXISTS laundry.service_variants (
    id_service_variant UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_service UUID NOT NULL,
    nama VARCHAR(120) NOT NULL,
    satuan VARCHAR(10) NOT NULL CHECK (satuan IN ('kg', 'pcs', 'meter')),
    harga_per_satuan NUMERIC(12,2) NOT NULL,
    durasi_pengerjaan_jam INTEGER NOT NULL,
    gambar_path VARCHAR(255),
    catatan TEXT,
    is_active BOOLEAN DEFAULT true,
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP,
    CONSTRAINT fk_service_variants_service FOREIGN KEY (id_service) REFERENCES laundry.services(id_service) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_service_variants_service_active ON laundry.service_variants(id_service, is_active);

CREATE TABLE IF NOT EXISTS laundry.perfumes (
    id_perfume UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_outlet UUID NOT NULL,
    nama VARCHAR(120) NOT NULL,
    catatan TEXT,
    is_active BOOLEAN DEFAULT true,
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP,
    CONSTRAINT fk_perfumes_outlet FOREIGN KEY (id_outlet) REFERENCES laundry.outlets(id_outlet) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_perfumes_outlet_active_nama ON laundry.perfumes(id_outlet, is_active, nama);

CREATE TABLE IF NOT EXISTS laundry.discounts (
    id_discount UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_outlet UUID NOT NULL,
    nama VARCHAR(120) NOT NULL,
    jenis VARCHAR(10) NOT NULL CHECK (jenis IN ('nominal', 'percent')),
    nilai NUMERIC(12,2) NOT NULL,
    catatan TEXT,
    is_active BOOLEAN DEFAULT true,
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP,
    CONSTRAINT fk_discounts_outlet FOREIGN KEY (id_outlet) REFERENCES laundry.outlets(id_outlet) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_discounts_outlet_active ON laundry.discounts(id_outlet, is_active);

-- ============================================================
-- PELANGGAN
-- ============================================================

CREATE TABLE IF NOT EXISTS laundry.customers (
    id_customer UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_outlet UUID NOT NULL,
    nama VARCHAR(255) NOT NULL,
    telepon VARCHAR(20),
    email VARCHAR(255),
    alamat TEXT,
    is_active BOOLEAN DEFAULT true,
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP,
    CONSTRAINT fk_customers_outlet FOREIGN KEY (id_outlet) REFERENCES laundry.outlets(id_outlet) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_customers_outlet_nama ON laundry.customers(id_outlet, nama);
CREATE INDEX IF NOT EXISTS idx_customers_outlet_telepon ON laundry.customers(id_outlet, telepon);
CREATE INDEX IF NOT EXISTS idx_customers_outlet_email ON laundry.customers(id_outlet, email);
CREATE INDEX IF NOT EXISTS idx_customers_outlet_active ON laundry.customers(id_outlet, is_active);

-- ============================================================
-- PESANAN (ORDER)
-- ============================================================

CREATE TABLE IF NOT EXISTS laundry.orders (
    id_order UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_outlet UUID NOT NULL,
    id_customer UUID NOT NULL,
    invoice_no VARCHAR(30) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'ANTRIAN' CHECK (status IN ('ANTRIAN','PROSES','SIAP_DIAMBIL','SELESAI','BATAL')),
    status_pembayaran VARCHAR(10) NOT NULL DEFAULT 'UNPAID' CHECK (status_pembayaran IN ('UNPAID','PAID')),
    id_payment_method UUID,
    id_perfume UUID,
    id_discount UUID,
    nilai_diskon NUMERIC(12,2) NOT NULL DEFAULT 0,
    subtotal NUMERIC(12,2) NOT NULL,
    total NUMERIC(12,2) NOT NULL,
    catatan TEXT,
    checkin_at TIMESTAMP NOT NULL DEFAULT NOW(),
    eta_at TIMESTAMP,
    selesai_at TIMESTAMP,
    batal_at TIMESTAMP,
    diambil_at TIMESTAMP,
    diambil_oleh_id_user UUID,
    dibuat_oleh_id_user UUID NOT NULL,
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP,
    CONSTRAINT fk_orders_outlet FOREIGN KEY (id_outlet) REFERENCES laundry.outlets(id_outlet) ON DELETE CASCADE,
    CONSTRAINT fk_orders_customer FOREIGN KEY (id_customer) REFERENCES laundry.customers(id_customer) ON DELETE CASCADE,
    CONSTRAINT fk_orders_payment_method FOREIGN KEY (id_payment_method) REFERENCES laundry.payment_methods(id_payment_method) ON DELETE SET NULL,
    CONSTRAINT fk_orders_perfume FOREIGN KEY (id_perfume) REFERENCES laundry.perfumes(id_perfume) ON DELETE SET NULL,
    CONSTRAINT fk_orders_discount FOREIGN KEY (id_discount) REFERENCES laundry.discounts(id_discount) ON DELETE SET NULL,
    CONSTRAINT fk_orders_collected_by FOREIGN KEY (diambil_oleh_id_user) REFERENCES laundry.users(id_user) ON DELETE SET NULL,
    CONSTRAINT fk_orders_created_by FOREIGN KEY (dibuat_oleh_id_user) REFERENCES laundry.users(id_user) ON DELETE CASCADE,
    CONSTRAINT uq_orders_outlet_invoice UNIQUE (id_outlet, invoice_no)
);
CREATE INDEX IF NOT EXISTS idx_orders_outlet_status ON laundry.orders(id_outlet, status);
CREATE INDEX IF NOT EXISTS idx_orders_outlet_created ON laundry.orders(id_outlet, tgl_insert);
CREATE INDEX IF NOT EXISTS idx_orders_outlet_eta ON laundry.orders(id_outlet, eta_at);

CREATE TABLE IF NOT EXISTS laundry.order_items (
    id_order_item UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_order UUID NOT NULL,
    id_service_variant UUID NOT NULL,
    satuan VARCHAR(10) NOT NULL CHECK (satuan IN ('kg', 'pcs', 'meter')),
    qty NUMERIC(10,2) NOT NULL,
    harga_per_satuan_snapshot NUMERIC(12,2) NOT NULL,
    total_harga NUMERIC(12,2) NOT NULL,
    catatan TEXT,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP,
    CONSTRAINT fk_order_items_order FOREIGN KEY (id_order) REFERENCES laundry.orders(id_order) ON DELETE CASCADE,
    CONSTRAINT fk_order_items_variant FOREIGN KEY (id_service_variant) REFERENCES laundry.service_variants(id_service_variant) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_order_items_order ON laundry.order_items(id_order);

CREATE TABLE IF NOT EXISTS laundry.order_status_histories (
    id_history UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_order UUID NOT NULL,
    status_dari VARCHAR(20),
    status_ke VARCHAR(20) NOT NULL,
    id_user UUID,
    catatan TEXT,
    waktu_perubahan TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_order_status_histories_order FOREIGN KEY (id_order) REFERENCES laundry.orders(id_order) ON DELETE CASCADE,
    CONSTRAINT fk_order_status_histories_user FOREIGN KEY (id_user) REFERENCES laundry.users(id_user) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_order_status_histories_order ON laundry.order_status_histories(id_order);
CREATE INDEX IF NOT EXISTS idx_order_status_histories_order_time ON laundry.order_status_histories(id_order, waktu_perubahan);

CREATE TABLE IF NOT EXISTS laundry.payments (
    id_payment UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_order UUID NOT NULL UNIQUE,
    id_payment_method UUID NOT NULL,
    jumlah NUMERIC(12,2) NOT NULL,
    dibayar_at TIMESTAMP NOT NULL,
    no_referensi VARCHAR(255),
    catatan TEXT,
    status VARCHAR(10) NOT NULL DEFAULT 'SUCCESS' CHECK (status IN ('SUCCESS','VOID')),
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP,
    CONSTRAINT fk_payments_order FOREIGN KEY (id_order) REFERENCES laundry.orders(id_order) ON DELETE CASCADE,
    CONSTRAINT fk_payments_method FOREIGN KEY (id_payment_method) REFERENCES laundry.payment_methods(id_payment_method) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_payments_method ON laundry.payments(id_payment_method);
CREATE INDEX IF NOT EXISTS idx_payments_paid_at ON laundry.payments(dibayar_at);
CREATE INDEX IF NOT EXISTS idx_payments_status ON laundry.payments(status);

-- ============================================================
-- PENGELUARAN (EXPENSE)
-- ============================================================

CREATE TABLE IF NOT EXISTS laundry.expenses (
    id_expense UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_outlet UUID NOT NULL,
    kategori VARCHAR(100) NOT NULL,
    jumlah NUMERIC(12,2) NOT NULL,
    deskripsi TEXT NOT NULL,
    tanggal DATE NOT NULL,
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP,
    CONSTRAINT fk_expenses_outlet FOREIGN KEY (id_outlet) REFERENCES laundry.outlets(id_outlet) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_expenses_outlet_tanggal ON laundry.expenses(id_outlet, tanggal);
CREATE INDEX IF NOT EXISTS idx_expenses_outlet_kategori ON laundry.expenses(id_outlet, kategori);
