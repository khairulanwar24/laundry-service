/*
 Navicat Premium Data Transfer

 Source Server         : sso development
 Source Server Type    : PostgreSQL
 Source Server Version : 110022 (110022)
 Source Host           : 112.78.41.222:5432
 Source Catalog        : sso
 Source Schema         : sso

 Target Server Type    : PostgreSQL
 Target Server Version : 110022 (110022)
 File Encoding         : 65001

 Date: 29/12/2024 12:55:34
*/


-- ----------------------------
-- Table structure for group_akses
-- ----------------------------
DROP TABLE IF EXISTS "sso"."group_akses";
CREATE TABLE "sso"."group_akses" (
  "id_group_akses" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "id_master_group" uuid NOT NULL,
  "id_master_aplikasi" uuid NOT NULL,
  "id_master_modul" uuid NOT NULL,
  "akses" varchar COLLATE "pg_catalog"."default"
)
;
ALTER TABLE "sso"."group_akses" OWNER TO "sso";

-- ----------------------------
-- Table structure for master_aplikasi
-- ----------------------------
DROP TABLE IF EXISTS "sso"."master_aplikasi";
CREATE TABLE "sso"."master_aplikasi" (
  "id_master_aplikasi" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "nama_aplikasi" varchar COLLATE "pg_catalog"."default" NOT NULL,
  "deskripsi" varchar COLLATE "pg_catalog"."default",
  "versi_aplikasi" varchar COLLATE "pg_catalog"."default",
  "status_data" bool DEFAULT true,
  "tgl_version" date,
  "url" varchar COLLATE "pg_catalog"."default",
  "image" varchar COLLATE "pg_catalog"."default"
)
;
ALTER TABLE "sso"."master_aplikasi" OWNER TO "sso";
COMMENT ON COLUMN "sso"."master_aplikasi"."nama_aplikasi" IS 'Nama Aplikasi';
COMMENT ON COLUMN "sso"."master_aplikasi"."deskripsi" IS 'Deskripsi Aplikasi';
COMMENT ON COLUMN "sso"."master_aplikasi"."versi_aplikasi" IS 'Versi Aplikasi';
COMMENT ON COLUMN "sso"."master_aplikasi"."image" IS 'Gambar Aplikasi Wajib PNG';
COMMENT ON TABLE "sso"."master_aplikasi" IS 'Master Aplikasi';

-- ----------------------------
-- Table structure for master_group
-- ----------------------------
DROP TABLE IF EXISTS "sso"."master_group";
CREATE TABLE "sso"."master_group" (
  "id_master_group" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "id_master_aplikasi" uuid NOT NULL,
  "nama_group" varchar COLLATE "pg_catalog"."default" NOT NULL,
  "deskripsi" varchar COLLATE "pg_catalog"."default",
  "status_data" bool DEFAULT true
)
;
ALTER TABLE "sso"."master_group" OWNER TO "sso";

-- ----------------------------
-- Table structure for master_menu
-- ----------------------------
DROP TABLE IF EXISTS "sso"."master_menu";
CREATE TABLE "sso"."master_menu" (
  "id_master_menu" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "nama_menu" varchar COLLATE "pg_catalog"."default" NOT NULL,
  "deskripsi" varchar COLLATE "pg_catalog"."default",
  "order" int4 NOT NULL,
  "id_master_aplikasi" uuid NOT NULL,
  "icon" varchar COLLATE "pg_catalog"."default",
  "status_data" bool DEFAULT true
)
;
ALTER TABLE "sso"."master_menu" OWNER TO "sso";
COMMENT ON COLUMN "sso"."master_menu"."nama_menu" IS 'nama Menu di aplikasi';
COMMENT ON COLUMN "sso"."master_menu"."deskripsi" IS 'deskripsi pada menu yang akan  tampil di page module';
COMMENT ON COLUMN "sso"."master_menu"."order" IS 'urutan aplikasi';
COMMENT ON COLUMN "sso"."master_menu"."id_master_aplikasi" IS 'id master aplikasi';
COMMENT ON COLUMN "sso"."master_menu"."icon" IS 'diambil dari component css';

-- ----------------------------
-- Table structure for master_modul
-- ----------------------------
DROP TABLE IF EXISTS "sso"."master_modul";
CREATE TABLE "sso"."master_modul" (
  "id_master_modul" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "id_master_menu" uuid NOT NULL,
  "order" int4,
  "nama_modul" varchar COLLATE "pg_catalog"."default" NOT NULL,
  "path" varchar COLLATE "pg_catalog"."default" NOT NULL,
  "deskripsi" varchar COLLATE "pg_catalog"."default",
  "id_master_aplikasi" uuid NOT NULL,
  "icon" varchar COLLATE "pg_catalog"."default",
  "status_data" bool DEFAULT true
)
;
ALTER TABLE "sso"."master_modul" OWNER TO "sso";
COMMENT ON COLUMN "sso"."master_modul"."path" IS 'path atau / aplikasi';
COMMENT ON COLUMN "sso"."master_modul"."deskripsi" IS 'deskripsi module';
COMMENT ON TABLE "sso"."master_modul" IS 'master modul';

-- ----------------------------
-- Table structure for password_reset
-- ----------------------------
DROP TABLE IF EXISTS "sso"."password_reset";
CREATE TABLE "sso"."password_reset" (
  "id_password_reset" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "id_user" uuid NOT NULL,
  "otp" varchar(6) COLLATE "pg_catalog"."default" NOT NULL,
  "tgl_insert" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
  "status_data" bool DEFAULT true,
  "username" varchar(100) COLLATE "pg_catalog"."default",
  "percobaan" int4 DEFAULT 0
)
;
ALTER TABLE "sso"."password_reset" OWNER TO "sso";
COMMENT ON COLUMN "sso"."password_reset"."otp" IS 'token 6 digit';
COMMENT ON COLUMN "sso"."password_reset"."tgl_insert" IS 'dibuat untuk tgl dibuat \ndan berlaku token hanya 10 menit';
COMMENT ON COLUMN "sso"."password_reset"."status_data" IS 'jika insert baru semua  diupdate delete';
COMMENT ON TABLE "sso"."password_reset" IS 'untuk reset password';

-- ----------------------------
-- Table structure for session_users
-- ----------------------------
DROP TABLE IF EXISTS "sso"."session_users";
CREATE TABLE "sso"."session_users" (
  "id_session_users" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "id_user" uuid,
  "refresh_token" varchar COLLATE "pg_catalog"."default",
  "is_revoked" bool,
  "tgl_insert" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
  "expires_at" timestamp(6)
)
;
ALTER TABLE "sso"."session_users" OWNER TO "sso";
COMMENT ON COLUMN "sso"."session_users"."id_user" IS 'id user';
COMMENT ON TABLE "sso"."session_users" IS 'Session Refresh token';

-- ----------------------------
-- Table structure for trans_user_group
-- ----------------------------
DROP TABLE IF EXISTS "sso"."trans_user_group";
CREATE TABLE "sso"."trans_user_group" (
  "id_trans_user_group" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "id_user" uuid NOT NULL,
  "id_master_aplikasi" uuid NOT NULL,
  "id_master_group" uuid NOT NULL,
  "status_data" bool NOT NULL DEFAULT true
)
;
ALTER TABLE "sso"."trans_user_group" OWNER TO "sso";

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS "sso"."users";
CREATE TABLE "sso"."users" (
  "id_user" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "username" varchar COLLATE "pg_catalog"."default" NOT NULL,
  "password" text COLLATE "pg_catalog"."default" NOT NULL,
  "email" text COLLATE "pg_catalog"."default",
  "nama_lengkap" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "tgl_insert" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
  "tgl_update" timestamp(6),
  "avatar" varchar COLLATE "pg_catalog"."default",
  "status_data" bool NOT NULL DEFAULT true,
  "id_person" uuid NOT NULL,
  "jenis_user" varchar COLLATE "pg_catalog"."default" NOT NULL,
  "first_login" bool NOT NULL DEFAULT true,
  "tgl_first_login" timestamp(6),
  "no_hp" varchar(20) COLLATE "pg_catalog"."default"
)
;
ALTER TABLE "sso"."users" OWNER TO "sso";
COMMENT ON COLUMN "sso"."users"."id_user" IS 'uuid_generate_v4()';
COMMENT ON COLUMN "sso"."users"."username" IS 'username \nharus uniq';
COMMENT ON COLUMN "sso"."users"."password" IS 'pasword \nmenggunakan hash bcrypt';
COMMENT ON COLUMN "sso"."users"."email" IS 'email boleh kosong waktu create awal\nwajib setelah update';
COMMENT ON COLUMN "sso"."users"."nama_lengkap" IS 'maksimal 100 karakter';
COMMENT ON COLUMN "sso"."users"."tgl_insert" IS 'tgl insert data';
COMMENT ON COLUMN "sso"."users"."tgl_update" IS 'tgl update data';
COMMENT ON COLUMN "sso"."users"."avatar" IS 'gambar avatar user';
COMMENT ON COLUMN "sso"."users"."status_data" IS 'true = aktif\nfalse = aktif';
COMMENT ON COLUMN "sso"."users"."id_person" IS 'digunakan untuk relasi dengan tabel akademik / pegawai';
COMMENT ON COLUMN "sso"."users"."jenis_user" IS 'berisi jenis user
1. Dosen
2. Tenaga Pendidik
3. Mahasiswa
4. Orang tua
5. Perseptor';
COMMENT ON COLUMN "sso"."users"."first_login" IS 'jika belum pernah login  maka status true dan user dipaksa mengganti password untuk pertama kali login. dan setelah  itu diupdate status false';
COMMENT ON COLUMN "sso"."users"."tgl_first_login" IS 'jika sudah mengganti password pertama kali maka diupdate tanggal ini';
COMMENT ON TABLE "sso"."users" IS 'Tabel User';

-- ----------------------------
-- Primary Key structure for table group_akses
-- ----------------------------
ALTER TABLE "sso"."group_akses" ADD CONSTRAINT "group_akses_pk" PRIMARY KEY ("id_group_akses");

-- ----------------------------
-- Primary Key structure for table master_aplikasi
-- ----------------------------
ALTER TABLE "sso"."master_aplikasi" ADD CONSTRAINT "pk_master_aplikasi" PRIMARY KEY ("id_master_aplikasi");

-- ----------------------------
-- Primary Key structure for table master_group
-- ----------------------------
ALTER TABLE "sso"."master_group" ADD CONSTRAINT "master_group_pk" PRIMARY KEY ("id_master_group");

-- ----------------------------
-- Primary Key structure for table master_menu
-- ----------------------------
ALTER TABLE "sso"."master_menu" ADD CONSTRAINT "pk_master_menu" PRIMARY KEY ("id_master_menu");

-- ----------------------------
-- Primary Key structure for table master_modul
-- ----------------------------
ALTER TABLE "sso"."master_modul" ADD CONSTRAINT "master_modul_pk" PRIMARY KEY ("id_master_modul");

-- ----------------------------
-- Primary Key structure for table password_reset
-- ----------------------------
ALTER TABLE "sso"."password_reset" ADD CONSTRAINT "pk_password_reset" PRIMARY KEY ("id_password_reset");

-- ----------------------------
-- Primary Key structure for table session_users
-- ----------------------------
ALTER TABLE "sso"."session_users" ADD CONSTRAINT "pk_session_users" PRIMARY KEY ("id_session_users");

-- ----------------------------
-- Uniques structure for table users
-- ----------------------------
ALTER TABLE "sso"."users" ADD CONSTRAINT "idx_username" UNIQUE ("username", "status_data");

-- ----------------------------
-- Primary Key structure for table users
-- ----------------------------
ALTER TABLE "sso"."users" ADD CONSTRAINT "pk_users" PRIMARY KEY ("id_user");

-- ----------------------------
-- Foreign Keys structure for table group_akses
-- ----------------------------
ALTER TABLE "sso"."group_akses" ADD CONSTRAINT "fk_group_akses_master_aplikasi" FOREIGN KEY ("id_master_aplikasi") REFERENCES "sso"."master_aplikasi" ("id_master_aplikasi") ON DELETE RESTRICT ON UPDATE RESTRICT;
ALTER TABLE "sso"."group_akses" ADD CONSTRAINT "fk_group_akses_master_modul" FOREIGN KEY ("id_master_modul") REFERENCES "sso"."master_modul" ("id_master_modul") ON DELETE RESTRICT ON UPDATE RESTRICT;
ALTER TABLE "sso"."group_akses" ADD CONSTRAINT "group_akses_fk" FOREIGN KEY ("id_master_group") REFERENCES "sso"."master_group" ("id_master_group") ON DELETE RESTRICT ON UPDATE RESTRICT;

-- ----------------------------
-- Foreign Keys structure for table master_group
-- ----------------------------
ALTER TABLE "sso"."master_group" ADD CONSTRAINT "master_group_fk" FOREIGN KEY ("id_master_aplikasi") REFERENCES "sso"."master_aplikasi" ("id_master_aplikasi") ON DELETE RESTRICT ON UPDATE RESTRICT;

-- ----------------------------
-- Foreign Keys structure for table master_menu
-- ----------------------------
ALTER TABLE "sso"."master_menu" ADD CONSTRAINT "fk_master_menu_master_aplikasi" FOREIGN KEY ("id_master_aplikasi") REFERENCES "sso"."master_aplikasi" ("id_master_aplikasi") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table master_modul
-- ----------------------------
ALTER TABLE "sso"."master_modul" ADD CONSTRAINT "fk_master_modul_master_aplikasi" FOREIGN KEY ("id_master_aplikasi") REFERENCES "sso"."master_aplikasi" ("id_master_aplikasi") ON DELETE RESTRICT ON UPDATE RESTRICT;
ALTER TABLE "sso"."master_modul" ADD CONSTRAINT "fk_master_modul_master_menu" FOREIGN KEY ("id_master_menu") REFERENCES "sso"."master_menu" ("id_master_menu") ON DELETE RESTRICT ON UPDATE RESTRICT;

-- ----------------------------
-- Foreign Keys structure for table trans_user_group
-- ----------------------------
ALTER TABLE "sso"."trans_user_group" ADD CONSTRAINT "fk_trans_user_group_master_aplikasi" FOREIGN KEY ("id_master_aplikasi") REFERENCES "sso"."master_aplikasi" ("id_master_aplikasi") ON DELETE RESTRICT ON UPDATE RESTRICT;
ALTER TABLE "sso"."trans_user_group" ADD CONSTRAINT "fk_trans_user_group_master_group" FOREIGN KEY ("id_master_group") REFERENCES "sso"."master_group" ("id_master_group") ON DELETE RESTRICT ON UPDATE RESTRICT;
ALTER TABLE "sso"."trans_user_group" ADD CONSTRAINT "fk_trans_user_group_users" FOREIGN KEY ("id_user") REFERENCES "sso"."users" ("id_user") ON DELETE RESTRICT ON UPDATE RESTRICT;
