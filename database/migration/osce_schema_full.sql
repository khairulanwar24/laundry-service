-- ============================================================
-- QUERY UNTUK MEMBUAT SCHEMA & TABEL OSCE
-- Jalankan query ini langsung di database "sso" via psql/DBeaver
-- Semua nama field menggunakan Bahasa Indonesia
-- ============================================================

CREATE SCHEMA IF NOT EXISTS osce;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
SET search_path TO osce, public;

-- ============================================================
-- TABEL REFERENSI / PENDUKUNG
-- ============================================================

-- Tahun Akademik (contoh: 2024/2025)
CREATE TABLE IF NOT EXISTS osce.tahun_akademik (
    id_tahun_akademik UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nama_tahun_akademik VARCHAR(20) NOT NULL UNIQUE,
    tgl_mulai DATE NOT NULL,
    tgl_selesai DATE NOT NULL,
    status_aktif BOOLEAN DEFAULT false,
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP
);

-- Semester
CREATE TABLE IF NOT EXISTS osce.semester (
    id_semester UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_tahun_akademik UUID NOT NULL,
    nama_semester VARCHAR(20) NOT NULL,
    tgl_mulai DATE NOT NULL,
    tgl_selesai DATE NOT NULL,
    status_aktif BOOLEAN DEFAULT false,
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP,
    CONSTRAINT fk_semester_tahun_akademik FOREIGN KEY (id_tahun_akademik)
        REFERENCES osce.tahun_akademik(id_tahun_akademik),
    CONSTRAINT uq_semester_tahun_nama UNIQUE (id_tahun_akademik, nama_semester)
);

-- Program Studi
CREATE TABLE IF NOT EXISTS osce.program_studi (
    id_program_studi UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_prodi VARCHAR(32) NOT NULL,
    kode_prodi VARCHAR(16) NOT NULL,
    nama_prodi VARCHAR(128) NOT NULL,
    jenjang VARCHAR(32),
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP,
    CONSTRAINT uq_program_studi_id_prodi UNIQUE (id_prodi)
);

-- ============================================================
-- DOMAIN STATION
-- ============================================================

-- Tipe Station
CREATE TABLE IF NOT EXISTS osce.tipe_station (
    id_tipe_station UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nama_tipe VARCHAR(64) NOT NULL UNIQUE,
    deskripsi TEXT,
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP
);

-- Master Station Ujian
CREATE TABLE IF NOT EXISTS osce.mst_station (
    id_station UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_tipe_station UUID NOT NULL,
    kode_station VARCHAR(16) NOT NULL UNIQUE,
    nama_station VARCHAR(128) NOT NULL,
    deskripsi TEXT,
    durasi_default INTEGER NOT NULL DEFAULT 10 CHECK (durasi_default > 0),
    bobot DECIMAL(5,2) DEFAULT 1.00 CHECK (bobot > 0),
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP,
    CONSTRAINT fk_station_tipe FOREIGN KEY (id_tipe_station)
        REFERENCES osce.tipe_station(id_tipe_station)
);

-- Kompetensi yang dinilai per station
CREATE TABLE IF NOT EXISTS osce.kompetensi_station (
    id_kompetensi_station UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_station UUID NOT NULL,
    nama_kompetensi VARCHAR(128) NOT NULL,
    deskripsi TEXT,
    bobot DECIMAL(5,2) DEFAULT 1.00 CHECK (bobot > 0),
    urutan INTEGER DEFAULT 1,
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP,
    CONSTRAINT fk_kompetensi_station FOREIGN KEY (id_station)
        REFERENCES osce.mst_station(id_station)
);

-- Checklist / Rubrik Detail per Station
CREATE TABLE IF NOT EXISTS osce.checklist_station (
    id_checklist_station UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_station UUID NOT NULL,
    id_kompetensi_station UUID,
    urutan INTEGER NOT NULL DEFAULT 1,
    deskripsi_item TEXT NOT NULL,
    skor_maksimal INTEGER NOT NULL DEFAULT 1 CHECK (skor_maksimal > 0),
    bobot_item DECIMAL(5,2) DEFAULT 1.00,
    tipe_penilaian VARCHAR(32) DEFAULT 'skala'
        CHECK (tipe_penilaian IN ('skala','biner','penilaian_global','catatan')),
    petunjuk_penilaian TEXT,
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP,
    CONSTRAINT fk_checklist_station FOREIGN KEY (id_station)
        REFERENCES osce.mst_station(id_station),
    CONSTRAINT fk_checklist_kompetensi FOREIGN KEY (id_kompetensi_station)
        REFERENCES osce.kompetensi_station(id_kompetensi_station)
);

-- ============================================================
-- DOMAIN UJIAN
-- ============================================================

-- Blueprint Ujian
CREATE TABLE IF NOT EXISTS osce.ujian (
    id_ujian UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_tahun_akademik UUID NOT NULL,
    id_semester UUID NOT NULL,
    id_program_studi UUID NOT NULL,
    nama_ujian VARCHAR(128) NOT NULL,
    deskripsi TEXT,
    tgl_mulai DATE NOT NULL,
    tgl_selesai DATE NOT NULL CHECK (tgl_selesai >= tgl_mulai),
    durasi_per_station INTEGER DEFAULT 10,
    jumlah_station INTEGER DEFAULT 0,
    batas_lulus DECIMAL(6,2),
    status_ujian VARCHAR(32) DEFAULT 'draft'
        CHECK (status_ujian IN ('draft','publikasi','berlangsung','selesai','dibatalkan')),
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP,
    CONSTRAINT fk_ujian_tahun_akademik FOREIGN KEY (id_tahun_akademik)
        REFERENCES osce.tahun_akademik(id_tahun_akademik),
    CONSTRAINT fk_ujian_semester FOREIGN KEY (id_semester)
        REFERENCES osce.semester(id_semester),
    CONSTRAINT fk_ujian_prodi FOREIGN KEY (id_program_studi)
        REFERENCES osce.program_studi(id_program_studi)
);

-- Relasi Ujian ↔ Station
CREATE TABLE IF NOT EXISTS osce.ujian_station (
    id_ujian_station UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_ujian UUID NOT NULL,
    id_station UUID NOT NULL,
    urutan INTEGER NOT NULL CHECK (urutan > 0),
    durasi INTEGER,
    bobot DECIMAL(5,2) DEFAULT 1.00 CHECK (bobot > 0),
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP,
    CONSTRAINT fk_ujian_station_ujian FOREIGN KEY (id_ujian)
        REFERENCES osce.ujian(id_ujian),
    CONSTRAINT fk_ujian_station_station FOREIGN KEY (id_station)
        REFERENCES osce.mst_station(id_station),
    CONSTRAINT uq_ujian_station UNIQUE (id_ujian, id_station),
    CONSTRAINT uq_ujian_station_urutan UNIQUE (id_ujian, urutan)
);

-- Sesi Ujian
CREATE TABLE IF NOT EXISTS osce.sesi_ujian (
    id_sesi_ujian UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_ujian UUID NOT NULL,
    nama_sesi VARCHAR(64) NOT NULL,
    tgl_sesi DATE NOT NULL,
    jam_mulai TIME NOT NULL,
    jam_selesai TIME NOT NULL CHECK (jam_selesai > jam_mulai),
    lokasi VARCHAR(128),
    kuota_peserta INTEGER DEFAULT 0,
    status_sesi VARCHAR(32) DEFAULT 'terjadwal'
        CHECK (status_sesi IN ('terjadwal','berlangsung','selesai','dibatalkan')),
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP,
    CONSTRAINT fk_sesi_ujian FOREIGN KEY (id_ujian)
        REFERENCES osce.ujian(id_ujian)
);

-- Aturan Rotasi Mahasiswa
CREATE TABLE IF NOT EXISTS osce.rotasi_ujian (
    id_rotasi_ujian UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_ujian UUID NOT NULL,
    id_sesi_ujian UUID NOT NULL,
    id_ujian_station UUID NOT NULL,
    urutan_rotasi INTEGER NOT NULL,
    durasi INTEGER NOT NULL CHECK (durasi > 0),
    jeda_antar_station INTEGER DEFAULT 0,
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP,
    CONSTRAINT fk_rotasi_ujian FOREIGN KEY (id_ujian)
        REFERENCES osce.ujian(id_ujian),
    CONSTRAINT fk_rotasi_sesi FOREIGN KEY (id_sesi_ujian)
        REFERENCES osce.sesi_ujian(id_sesi_ujian),
    CONSTRAINT fk_rotasi_station FOREIGN KEY (id_ujian_station)
        REFERENCES osce.ujian_station(id_ujian_station),
    CONSTRAINT uq_rotasi_sesi_urutan UNIQUE (id_sesi_ujian, urutan_rotasi)
);

-- Jadwal Detail Mahasiswa
CREATE TABLE IF NOT EXISTS osce.jadwal_ujian (
    id_jadwal_ujian UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_ujian UUID NOT NULL,
    id_sesi_ujian UUID NOT NULL,
    id_ujian_station UUID NOT NULL,
    id_user VARCHAR(36) NOT NULL,
    id_penguji UUID,
    urutan_masuk INTEGER,
    jam_mulai TIME,
    jam_selesai TIME CHECK (jam_selesai > jam_mulai),
    status_kehadiran VARCHAR(32) DEFAULT 'terjadwal'
        CHECK (status_kehadiran IN ('terjadwal','hadir','tidak_hadir','izin')),
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP,
    CONSTRAINT fk_jadwal_ujian FOREIGN KEY (id_ujian)
        REFERENCES osce.ujian(id_ujian),
    CONSTRAINT fk_jadwal_sesi FOREIGN KEY (id_sesi_ujian)
        REFERENCES osce.sesi_ujian(id_sesi_ujian),
    CONSTRAINT fk_jadwal_station FOREIGN KEY (id_ujian_station)
        REFERENCES osce.ujian_station(id_ujian_station),
    CONSTRAINT uq_jadwal_unique UNIQUE (id_sesi_ujian, id_ujian_station, id_user)
);

-- ============================================================
-- DOMAIN PENILAIAN
-- ============================================================

-- Penilaian dari penguji
CREATE TABLE IF NOT EXISTS osce.penilaian (
    id_penilaian UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_jadwal_ujian UUID NOT NULL,
    id_penguji UUID NOT NULL,
    total_skor DECIMAL(8,2) DEFAULT 0,
    skor_terbobot DECIMAL(8,2) DEFAULT 0,
    penilaian_global VARCHAR(32)
        CHECK (penilaian_global IN ('borderline','lulus','tidak_lulus','istimewa')),
    catatan_penguji TEXT,
    status_penilaian VARCHAR(32) DEFAULT 'draft'
        CHECK (status_penilaian IN ('draft','terkirim','terverifikasi')),
    tgl_penilaian TIMESTAMP,
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP,
    CONSTRAINT fk_penilaian_jadwal FOREIGN KEY (id_jadwal_ujian)
        REFERENCES osce.jadwal_ujian(id_jadwal_ujian),
    CONSTRAINT uq_penilaian_unique UNIQUE (id_jadwal_ujian, id_penguji)
);

-- Detail nilai per checklist item
CREATE TABLE IF NOT EXISTS osce.detail_penilaian (
    id_detail_penilaian UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_penilaian UUID NOT NULL,
    id_checklist_station UUID NOT NULL,
    skor INTEGER NOT NULL DEFAULT 0,
    catatan_item TEXT,
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP,
    CONSTRAINT fk_detail_penilaian FOREIGN KEY (id_penilaian)
        REFERENCES osce.penilaian(id_penilaian),
    CONSTRAINT fk_detail_checklist FOREIGN KEY (id_checklist_station)
        REFERENCES osce.checklist_station(id_checklist_station),
    CONSTRAINT uq_detail_unique UNIQUE (id_penilaian, id_checklist_station)
);

-- Hasil Akhir per Mahasiswa per Ujian
CREATE TABLE IF NOT EXISTS osce.hasil_ujian (
    id_hasil_ujian UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_ujian UUID NOT NULL,
    id_user VARCHAR(36) NOT NULL,
    total_skor DECIMAL(8,2) DEFAULT 0,
    persentase_skor DECIMAL(6,2) DEFAULT 0,
    batas_lulus DECIMAL(6,2),
    status_kelulusan VARCHAR(32)
        CHECK (status_kelulusan IN ('lulus','tidak_lulus','remedi','tertunda')),
    catatan_hasil TEXT,
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP,
    CONSTRAINT fk_hasil_ujian FOREIGN KEY (id_ujian)
        REFERENCES osce.ujian(id_ujian),
    CONSTRAINT uq_hasil_unique UNIQUE (id_ujian, id_user)
);

-- Detail hasil per station per mahasiswa
CREATE TABLE IF NOT EXISTS osce.hasil_ujian_station (
    id_hasil_ujian_station UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_hasil_ujian UUID NOT NULL,
    id_ujian_station UUID NOT NULL,
    skor_station DECIMAL(8,2) DEFAULT 0,
    skor_terbobot DECIMAL(8,2) DEFAULT 0,
    penilaian_global VARCHAR(32),
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP,
    CONSTRAINT fk_hasil_station_hasil FOREIGN KEY (id_hasil_ujian)
        REFERENCES osce.hasil_ujian(id_hasil_ujian),
    CONSTRAINT fk_hasil_station_ujian_station FOREIGN KEY (id_ujian_station)
        REFERENCES osce.ujian_station(id_ujian_station),
    CONSTRAINT uq_hasil_station_unique UNIQUE (id_hasil_ujian, id_ujian_station)
);

-- ============================================================
-- DOMAIN PENGGUNA OSCE
-- ============================================================

-- Data Penguji
CREATE TABLE IF NOT EXISTS osce.penguji (
    id_penguji UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_user VARCHAR(36) NOT NULL,
    gelar_depan VARCHAR(32),
    gelar_belakang VARCHAR(32),
    keahlian VARCHAR(128),
    no_str VARCHAR(32),
    status_aktif BOOLEAN DEFAULT true,
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP
);

-- Kelompok Mahasiswa
CREATE TABLE IF NOT EXISTS osce.kelompok_mahasiswa (
    id_kelompok_mahasiswa UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_ujian UUID,
    nama_kelompok VARCHAR(64) NOT NULL,
    deskripsi TEXT,
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP,
    CONSTRAINT fk_kelompok_ujian FOREIGN KEY (id_ujian)
        REFERENCES osce.ujian(id_ujian)
);

-- Anggota Kelompok Mahasiswa
CREATE TABLE IF NOT EXISTS osce.anggota_kelompok (
    id_anggota_kelompok UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_kelompok_mahasiswa UUID NOT NULL,
    id_user VARCHAR(36) NOT NULL,
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP,
    CONSTRAINT fk_anggota_kelompok FOREIGN KEY (id_kelompok_mahasiswa)
        REFERENCES osce.kelompok_mahasiswa(id_kelompok_mahasiswa),
    CONSTRAINT uq_anggota_unique UNIQUE (id_kelompok_mahasiswa, id_user)
);

-- Penugasan Penguji ke Station pada Sesi Tertentu
CREATE TABLE IF NOT EXISTS osce.penugasan_penguji (
    id_penugasan_penguji UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_penguji UUID NOT NULL,
    id_sesi_ujian UUID NOT NULL,
    id_ujian_station UUID NOT NULL,
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP,
    CONSTRAINT fk_penugasan_penguji FOREIGN KEY (id_penguji)
        REFERENCES osce.penguji(id_penguji),
    CONSTRAINT fk_penugasan_sesi FOREIGN KEY (id_sesi_ujian)
        REFERENCES osce.sesi_ujian(id_sesi_ujian),
    CONSTRAINT fk_penugasan_station FOREIGN KEY (id_ujian_station)
        REFERENCES osce.ujian_station(id_ujian_station),
    CONSTRAINT uq_penugasan_unique UNIQUE (id_penguji, id_sesi_ujian, id_ujian_station)
);

-- ============================================================
-- INDEKS (performansi query)
-- ============================================================

CREATE INDEX IF NOT EXISTS idx_semester_aktif ON osce.semester(status_aktif) WHERE status_aktif = true;
CREATE INDEX IF NOT EXISTS idx_semester_tahun ON osce.semester(id_tahun_akademik);
CREATE INDEX IF NOT EXISTS idx_program_studi_kode ON osce.program_studi(kode_prodi);
CREATE INDEX IF NOT EXISTS idx_station_kode ON osce.mst_station(kode_station);
CREATE INDEX IF NOT EXISTS idx_station_tipe ON osce.mst_station(id_tipe_station);
CREATE INDEX IF NOT EXISTS idx_kompetensi_station ON osce.kompetensi_station(id_station);
CREATE INDEX IF NOT EXISTS idx_checklist_station ON osce.checklist_station(id_station);
CREATE INDEX IF NOT EXISTS idx_checklist_kompetensi ON osce.checklist_station(id_kompetensi_station);
CREATE INDEX IF NOT EXISTS idx_ujian_tahun_akademik ON osce.ujian(id_tahun_akademik);
CREATE INDEX IF NOT EXISTS idx_ujian_semester ON osce.ujian(id_semester);
CREATE INDEX IF NOT EXISTS idx_ujian_prodi ON osce.ujian(id_program_studi);
CREATE INDEX IF NOT EXISTS idx_ujian_status ON osce.ujian(status_ujian);
CREATE INDEX IF NOT EXISTS idx_ujian_station_ujian ON osce.ujian_station(id_ujian);
CREATE INDEX IF NOT EXISTS idx_sesi_ujian ON osce.sesi_ujian(id_ujian);
CREATE INDEX IF NOT EXISTS idx_sesi_tgl ON osce.sesi_ujian(tgl_sesi);
CREATE INDEX IF NOT EXISTS idx_rotasi_ujian ON osce.rotasi_ujian(id_ujian);
CREATE INDEX IF NOT EXISTS idx_rotasi_sesi ON osce.rotasi_ujian(id_sesi_ujian);
CREATE INDEX IF NOT EXISTS idx_jadwal_ujian ON osce.jadwal_ujian(id_ujian);
CREATE INDEX IF NOT EXISTS idx_jadwal_sesi ON osce.jadwal_ujian(id_sesi_ujian);
CREATE INDEX IF NOT EXISTS idx_jadwal_user ON osce.jadwal_ujian(id_user);
CREATE INDEX IF NOT EXISTS idx_jadwal_penguji ON osce.jadwal_ujian(id_penguji);
CREATE INDEX IF NOT EXISTS idx_penilaian_jadwal ON osce.penilaian(id_jadwal_ujian);
CREATE INDEX IF NOT EXISTS idx_penilaian_penguji ON osce.penilaian(id_penguji);
CREATE INDEX IF NOT EXISTS idx_detail_penilaian ON osce.detail_penilaian(id_penilaian);
CREATE INDEX IF NOT EXISTS idx_hasil_ujian ON osce.hasil_ujian(id_ujian);
CREATE INDEX IF NOT EXISTS idx_hasil_user ON osce.hasil_ujian(id_user);
CREATE INDEX IF NOT EXISTS idx_hasil_station_hasil ON osce.hasil_ujian_station(id_hasil_ujian);
CREATE UNIQUE INDEX IF NOT EXISTS idx_penguji_user ON osce.penguji(id_user) WHERE status_data = true;
CREATE INDEX IF NOT EXISTS idx_kelompok_ujian ON osce.kelompok_mahasiswa(id_ujian);
CREATE INDEX IF NOT EXISTS idx_anggota_kelompok ON osce.anggota_kelompok(id_kelompok_mahasiswa);
CREATE INDEX IF NOT EXISTS idx_anggota_user ON osce.anggota_kelompok(id_user);
CREATE INDEX IF NOT EXISTS idx_penugasan_sesi ON osce.penugasan_penguji(id_sesi_ujian);
CREATE INDEX IF NOT EXISTS idx_penugasan_station ON osce.penugasan_penguji(id_ujian_station);

-- Selesai --
SELECT 'Schema OSCE berhasil dibuat! ' || COUNT(*) || ' tabel telah siap.' AS hasil
FROM information_schema.tables
WHERE table_schema = 'osce';
