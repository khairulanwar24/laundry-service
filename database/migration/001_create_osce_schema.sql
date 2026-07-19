-- ============================================================
-- SCHEMA OSCE
-- OSCE (Objective Structured Clinical Examination) System
-- Farmasi Unissula
-- ============================================================

-- Aktifkan ekstensi UUID jika belum aktif
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================================
-- TABEL REFERENSI / PENDUKUNG
-- ============================================================

-- Tahun Akademik (contoh: 2024/2025)
CREATE TABLE IF NOT EXISTS osce.academic_year (
    id_academic_year UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nama_academic_year VARCHAR(20) NOT NULL UNIQUE,
    tgl_mulai DATE NOT NULL,
    tgl_selesai DATE NOT NULL,
    is_active BOOLEAN DEFAULT false,
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_academic_year_active ON osce.academic_year(is_active) WHERE is_active = true;

-- Semester
CREATE TABLE IF NOT EXISTS osce.semester (
    id_semester UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_academic_year UUID NOT NULL,
    nama_semester VARCHAR(20) NOT NULL,
    tgl_mulai DATE NOT NULL,
    tgl_selesai DATE NOT NULL,
    is_active BOOLEAN DEFAULT false,
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP,
    CONSTRAINT fk_semester_academic_year FOREIGN KEY (id_academic_year) REFERENCES osce.academic_year(id_academic_year),
    CONSTRAINT uq_semester_year_name UNIQUE (id_academic_year, nama_semester)
);
CREATE INDEX IF NOT EXISTS idx_semester_active ON osce.semester(is_active) WHERE is_active = true;
CREATE INDEX IF NOT EXISTS idx_semester_year ON osce.semester(id_academic_year);

-- Program Studi (relasi ke data akademik)
CREATE TABLE IF NOT EXISTS osce.study_program (
    id_study_program UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_prodi VARCHAR(32) NOT NULL,
    kode_prodi VARCHAR(16) NOT NULL,
    nama_prodi VARCHAR(128) NOT NULL,
    jenjang VARCHAR(32),
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP,
    CONSTRAINT uq_study_program_id_prodi UNIQUE (id_prodi)
);
CREATE INDEX IF NOT EXISTS idx_study_program_kode ON osce.study_program(kode_prodi);

-- ============================================================
-- DOMAIN STATION
-- ============================================================

-- Tipe Station: anamnesis, konseling, dispensing, drug interaction, dll
CREATE TABLE IF NOT EXISTS osce.station_type (
    id_station_type UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nama_type VARCHAR(64) NOT NULL UNIQUE,
    deskripsi TEXT,
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP
);

-- Master Station Ujian
CREATE TABLE IF NOT EXISTS osce.mst_station (
    id_station UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_station_type UUID NOT NULL,
    kode_station VARCHAR(16) NOT NULL UNIQUE,
    nama_station VARCHAR(128) NOT NULL,
    deskripsi TEXT,
    durasi_default INTEGER NOT NULL DEFAULT 10,
    bobot DECIMAL(5,2) DEFAULT 1.00,
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP,
    CONSTRAINT fk_station_type FOREIGN KEY (id_station_type) REFERENCES osce.station_type(id_station_type),
    CONSTRAINT ck_durasi_positive CHECK (durasi_default > 0),
    CONSTRAINT ck_bobot_positive CHECK (bobot > 0)
);
CREATE INDEX IF NOT EXISTS idx_station_kode ON osce.mst_station(kode_station);
CREATE INDEX IF NOT EXISTS idx_station_type ON osce.mst_station(id_station_type);

-- Kompetensi yang dinilai per station
CREATE TABLE IF NOT EXISTS osce.station_competency (
    id_station_competency UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_station UUID NOT NULL,
    nama_kompetensi VARCHAR(128) NOT NULL,
    deskripsi TEXT,
    bobot DECIMAL(5,2) DEFAULT 1.00,
    urutan INTEGER DEFAULT 1,
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP,
    CONSTRAINT fk_comp_station FOREIGN KEY (id_station) REFERENCES osce.mst_station(id_station),
    CONSTRAINT ck_comp_bobot_positive CHECK (bobot > 0)
);
CREATE INDEX IF NOT EXISTS idx_competency_station ON osce.station_competency(id_station);

-- Checklist / Rubrik Detail per Station
CREATE TABLE IF NOT EXISTS osce.station_checklist (
    id_station_checklist UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_station UUID NOT NULL,
    id_station_competency UUID,
    urutan INTEGER NOT NULL DEFAULT 1,
    deskripsi_item TEXT NOT NULL,
    skor_maksimal INTEGER NOT NULL DEFAULT 1,
    bobot_item DECIMAL(5,2) DEFAULT 1.00,
    tipe_penilaian VARCHAR(32) DEFAULT 'skala',
    petunjuk_penilaian TEXT,
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP,
    CONSTRAINT fk_checklist_station FOREIGN KEY (id_station) REFERENCES osce.mst_station(id_station),
    CONSTRAINT fk_checklist_competency FOREIGN KEY (id_station_competency) REFERENCES osce.station_competency(id_station_competency),
    CONSTRAINT ck_checklist_skor_positive CHECK (skor_maksimal > 0),
    CONSTRAINT ck_checklist_tipe CHECK (tipe_penilaian IN ('skala','binary','global_rating','catatan'))
);
CREATE INDEX IF NOT EXISTS idx_checklist_station ON osce.station_checklist(id_station);
CREATE INDEX IF NOT EXISTS idx_checklist_competency ON osce.station_checklist(id_station_competency);

-- ============================================================
-- DOMAIN UJIAN (EXAM)
-- ============================================================

-- Blueprint / Konfigurasi Ujian OSCE
CREATE TABLE IF NOT EXISTS osce.exam (
    id_exam UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_academic_year UUID NOT NULL,
    id_semester UUID NOT NULL,
    id_study_program UUID NOT NULL,
    nama_ujian VARCHAR(128) NOT NULL,
    deskripsi TEXT,
    tgl_mulai DATE NOT NULL,
    tgl_selesai DATE NOT NULL,
    durasi_per_station INTEGER DEFAULT 10,
    jumlah_station INTEGER DEFAULT 0,
    batas_lulus DECIMAL(6,2),
    status_ujian VARCHAR(32) DEFAULT 'draft',
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP,
    CONSTRAINT fk_exam_academic_year FOREIGN KEY (id_academic_year) REFERENCES osce.academic_year(id_academic_year),
    CONSTRAINT fk_exam_semester FOREIGN KEY (id_semester) REFERENCES osce.semester(id_semester),
    CONSTRAINT fk_exam_prodi FOREIGN KEY (id_study_program) REFERENCES osce.study_program(id_study_program),
    CONSTRAINT ck_exam_date CHECK (tgl_selesai >= tgl_mulai),
    CONSTRAINT ck_exam_status CHECK (status_ujian IN ('draft','published','in_progress','completed','cancelled'))
);
CREATE INDEX IF NOT EXISTS idx_exam_academic_year ON osce.exam(id_academic_year);
CREATE INDEX IF NOT EXISTS idx_exam_semester ON osce.exam(id_semester);
CREATE INDEX IF NOT EXISTS idx_exam_prodi ON osce.exam(id_study_program);
CREATE INDEX IF NOT EXISTS idx_exam_status ON osce.exam(status_ujian);

-- Relasi Ujian dengan Station yang digunakan
CREATE TABLE IF NOT EXISTS osce.exam_station (
    id_exam_station UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_exam UUID NOT NULL,
    id_station UUID NOT NULL,
    urutan INTEGER NOT NULL,
    durasi INTEGER,
    bobot DECIMAL(5,2) DEFAULT 1.00,
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP,
    CONSTRAINT fk_exam_station_exam FOREIGN KEY (id_exam) REFERENCES osce.exam(id_exam),
    CONSTRAINT fk_exam_station_station FOREIGN KEY (id_station) REFERENCES osce.mst_station(id_station),
    CONSTRAINT uq_exam_station UNIQUE (id_exam, id_station),
    CONSTRAINT uq_exam_station_urutan UNIQUE (id_exam, urutan),
    CONSTRAINT ck_exam_station_urutan_positive CHECK (urutan > 0),
    CONSTRAINT ck_exam_station_bobot_positive CHECK (bobot > 0)
);
CREATE INDEX IF NOT EXISTS idx_exam_station_exam ON osce.exam_station(id_exam);

-- Sesi Ujian (tanggal, waktu)
CREATE TABLE IF NOT EXISTS osce.exam_session (
    id_exam_session UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_exam UUID NOT NULL,
    nama_sesi VARCHAR(64) NOT NULL,
    tgl_sesi DATE NOT NULL,
    jam_mulai TIME NOT NULL,
    jam_selesai TIME NOT NULL,
    lokasi VARCHAR(128),
    kuota_peserta INTEGER DEFAULT 0,
    status_sesi VARCHAR(32) DEFAULT 'scheduled',
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP,
    CONSTRAINT fk_session_exam FOREIGN KEY (id_exam) REFERENCES osce.exam(id_exam),
    CONSTRAINT ck_session_time CHECK (jam_selesai > jam_mulai),
    CONSTRAINT ck_session_status CHECK (status_sesi IN ('scheduled','in_progress','completed','cancelled'))
);
CREATE INDEX IF NOT EXISTS idx_session_exam ON osce.exam_session(id_exam);
CREATE INDEX IF NOT EXISTS idx_session_date ON osce.exam_session(tgl_sesi);

-- Aturan Rotasi Mahasiswa
CREATE TABLE IF NOT EXISTS osce.exam_rotation (
    id_exam_rotation UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_exam UUID NOT NULL,
    id_exam_session UUID NOT NULL,
    id_exam_station UUID NOT NULL,
    urutan_rotasi INTEGER NOT NULL,
    durasi INTEGER NOT NULL,
    jeda_antar_station INTEGER DEFAULT 0,
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP,
    CONSTRAINT fk_rotation_exam FOREIGN KEY (id_exam) REFERENCES osce.exam(id_exam),
    CONSTRAINT fk_rotation_session FOREIGN KEY (id_exam_session) REFERENCES osce.exam_session(id_exam_session),
    CONSTRAINT fk_rotation_station FOREIGN KEY (id_exam_station) REFERENCES osce.exam_station(id_exam_station),
    CONSTRAINT ck_rotation_durasi_positive CHECK (durasi > 0),
    CONSTRAINT uq_rotation_session_urutan UNIQUE (id_exam_session, urutan_rotasi)
);
CREATE INDEX IF NOT EXISTS idx_rotation_exam ON osce.exam_rotation(id_exam);
CREATE INDEX IF NOT EXISTS idx_rotation_session ON osce.exam_rotation(id_exam_session);

-- Jadwal Detail: Mahasiswa X di Station Y, Sesi Z
CREATE TABLE IF NOT EXISTS osce.exam_schedule (
    id_exam_schedule UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_exam UUID NOT NULL,
    id_exam_session UUID NOT NULL,
    id_exam_station UUID NOT NULL,
    id_user VARCHAR(36) NOT NULL,
    id_examiner UUID,
    urutan_masuk INTEGER,
    jam_mulai TIME,
    jam_selesai TIME,
    status_kehadiran VARCHAR(32) DEFAULT 'terjadwal',
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP,
    CONSTRAINT fk_schedule_exam FOREIGN KEY (id_exam) REFERENCES osce.exam(id_exam),
    CONSTRAINT fk_schedule_session FOREIGN KEY (id_exam_session) REFERENCES osce.exam_session(id_exam_session),
    CONSTRAINT fk_schedule_station FOREIGN KEY (id_exam_station) REFERENCES osce.exam_station(id_exam_station),
    CONSTRAINT ck_schedule_time CHECK (jam_selesai > jam_mulai),
    CONSTRAINT ck_schedule_status CHECK (status_kehadiran IN ('terjadwal','hadir','tidak_hadir','izin')),
    CONSTRAINT uq_schedule_unique UNIQUE (id_exam_session, id_exam_station, id_user)
);
CREATE INDEX IF NOT EXISTS idx_schedule_exam ON osce.exam_schedule(id_exam);
CREATE INDEX IF NOT EXISTS idx_schedule_session ON osce.exam_schedule(id_exam_session);
CREATE INDEX IF NOT EXISTS idx_schedule_user ON osce.exam_schedule(id_user);
CREATE INDEX IF NOT EXISTS idx_schedule_examiner ON osce.exam_schedule(id_examiner);

-- ============================================================
-- DOMAIN PENILAIAN
-- ============================================================

-- Penilaian dari penguji untuk satu mahasiswa di satu station
CREATE TABLE IF NOT EXISTS osce.assessment (
    id_assessment UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_exam_schedule UUID NOT NULL,
    id_examiner UUID NOT NULL,
    total_skor DECIMAL(8,2) DEFAULT 0,
    skor_terbobot DECIMAL(8,2) DEFAULT 0,
    global_rating VARCHAR(32),
    catatan_penguji TEXT,
    status_penilaian VARCHAR(32) DEFAULT 'draft',
    tgl_penilaian TIMESTAMP,
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP,
    CONSTRAINT fk_assessment_schedule FOREIGN KEY (id_exam_schedule) REFERENCES osce.exam_schedule(id_exam_schedule),
    CONSTRAINT ck_assessment_global_rating CHECK (global_rating IN ('borderline','pass','fail','excellent')),
    CONSTRAINT ck_assessment_status CHECK (status_penilaian IN ('draft','submitted','verified')),
    CONSTRAINT uq_assessment_unique UNIQUE (id_exam_schedule, id_examiner)
);
CREATE INDEX IF NOT EXISTS idx_assessment_schedule ON osce.assessment(id_exam_schedule);
CREATE INDEX IF NOT EXISTS idx_assessment_examiner ON osce.assessment(id_examiner);

-- Detail nilai per checklist item
CREATE TABLE IF NOT EXISTS osce.assessment_detail (
    id_assessment_detail UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_assessment UUID NOT NULL,
    id_station_checklist UUID NOT NULL,
    skor INTEGER NOT NULL DEFAULT 0,
    catatan_item TEXT,
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP,
    CONSTRAINT fk_detail_assessment FOREIGN KEY (id_assessment) REFERENCES osce.assessment(id_assessment),
    CONSTRAINT fk_detail_checklist FOREIGN KEY (id_station_checklist) REFERENCES osce.station_checklist(id_station_checklist),
    CONSTRAINT uq_detail_unique UNIQUE (id_assessment, id_station_checklist)
);
CREATE INDEX IF NOT EXISTS idx_detail_assessment ON osce.assessment_detail(id_assessment);

-- Hasil Akhir per Mahasiswa per Ujian
CREATE TABLE IF NOT EXISTS osce.exam_result (
    id_exam_result UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_exam UUID NOT NULL,
    id_user VARCHAR(36) NOT NULL,
    total_skor DECIMAL(8,2) DEFAULT 0,
    persentase_skor DECIMAL(6,2) DEFAULT 0,
    batas_lulus DECIMAL(6,2),
    status_kelulusan VARCHAR(32),
    catatan_hasil TEXT,
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP,
    CONSTRAINT fk_result_exam FOREIGN KEY (id_exam) REFERENCES osce.exam(id_exam),
    CONSTRAINT ck_result_status CHECK (status_kelulusan IN ('lulus','tidak_lulus','remedi','pending')),
    CONSTRAINT uq_result_unique UNIQUE (id_exam, id_user)
);
CREATE INDEX IF NOT EXISTS idx_result_exam ON osce.exam_result(id_exam);
CREATE INDEX IF NOT EXISTS idx_result_user ON osce.exam_result(id_user);

-- Detail hasil per station per mahasiswa
CREATE TABLE IF NOT EXISTS osce.exam_result_station (
    id_exam_result_station UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_exam_result UUID NOT NULL,
    id_exam_station UUID NOT NULL,
    skor_station DECIMAL(8,2) DEFAULT 0,
    skor_terbobot DECIMAL(8,2) DEFAULT 0,
    global_rating VARCHAR(32),
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP,
    CONSTRAINT fk_result_station_result FOREIGN KEY (id_exam_result) REFERENCES osce.exam_result(id_exam_result),
    CONSTRAINT fk_result_station_exam_station FOREIGN KEY (id_exam_station) REFERENCES osce.exam_station(id_exam_station),
    CONSTRAINT uq_result_station_unique UNIQUE (id_exam_result, id_exam_station)
);
CREATE INDEX IF NOT EXISTS idx_result_station_result ON osce.exam_result_station(id_exam_result);

-- ============================================================
-- DOMAIN PENGGUNA OSCE
-- ============================================================

-- Data Penguji / Examiner (relasi ke tabel users di schema public)
CREATE TABLE IF NOT EXISTS osce.examiner (
    id_examiner UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_user VARCHAR(36) NOT NULL,
    gelar_depan VARCHAR(32),
    gelar_belakang VARCHAR(32),
    keahlian VARCHAR(128),
    no_str VARCHAR(32),
    is_active BOOLEAN DEFAULT true,
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_examiner_user ON osce.examiner(id_user) WHERE status_data = true;

-- Kelompok Mahasiswa
CREATE TABLE IF NOT EXISTS osce.student_group (
    id_student_group UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_exam UUID,
    nama_kelompok VARCHAR(64) NOT NULL,
    deskripsi TEXT,
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP,
    CONSTRAINT fk_group_exam FOREIGN KEY (id_exam) REFERENCES osce.exam(id_exam)
);
CREATE INDEX IF NOT EXISTS idx_group_exam ON osce.student_group(id_exam);

-- Anggota Kelompok Mahasiswa
CREATE TABLE IF NOT EXISTS osce.student_group_member (
    id_student_group_member UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_student_group UUID NOT NULL,
    id_user VARCHAR(36) NOT NULL,
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP,
    CONSTRAINT fk_member_group FOREIGN KEY (id_student_group) REFERENCES osce.student_group(id_student_group),
    CONSTRAINT uq_member_unique UNIQUE (id_student_group, id_user)
);
CREATE INDEX IF NOT EXISTS idx_member_group ON osce.student_group_member(id_student_group);
CREATE INDEX IF NOT EXISTS idx_member_user ON osce.student_group_member(id_user);

-- Assign Examiner ke Station pada Sesi Tertentu
CREATE TABLE IF NOT EXISTS osce.examiner_assignment (
    id_examiner_assignment UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_examiner UUID NOT NULL,
    id_exam_session UUID NOT NULL,
    id_exam_station UUID NOT NULL,
    status_data BOOLEAN DEFAULT true,
    tgl_insert TIMESTAMP NOT NULL DEFAULT NOW(),
    tgl_update TIMESTAMP,
    CONSTRAINT fk_assign_examiner FOREIGN KEY (id_examiner) REFERENCES osce.examiner(id_examiner),
    CONSTRAINT fk_assign_session FOREIGN KEY (id_exam_session) REFERENCES osce.exam_session(id_exam_session),
    CONSTRAINT fk_assign_station FOREIGN KEY (id_exam_station) REFERENCES osce.exam_station(id_exam_station),
    CONSTRAINT uq_assign_unique UNIQUE (id_examiner, id_exam_session, id_exam_station)
);
CREATE INDEX IF NOT EXISTS idx_assign_session ON osce.examiner_assignment(id_exam_session);
CREATE INDEX IF NOT EXISTS idx_assign_station ON osce.examiner_assignment(id_exam_station);
