# Golang untuk Pemula — Belajar dari Kode OSCE

> Dokumen ini menjelaskan konsep-konsep dasar Go (Golang) memakai contoh kode **asli** dari domain OSCE di proyek ini (`services/osce`, `repositories/osce`, `controllers/osce`, `routes/osce`, `domain/dto/osce`). Cocok dibaca kalau kamu baru pertama kali membaca kode Go dan ingin paham alur request dari HTTP masuk sampai ke database.

---

## Daftar Isi

1. [Package — cara Go mengelompokkan kode](#1-package--cara-go-mengelompokkan-kode)
2. [Struct — pengganti "class" di Go](#2-struct--pengganti-class-di-go)
3. [Struct tag — anotasi di belakang field](#3-struct-tag--anotasi-di-belakang-field)
4. [Interface — kontrak, bukan implementasi](#4-interface--kontrak-bukan-implementasi)
5. [Constructor — fungsi `New...`](#5-constructor--fungsi-new)
6. [Method — fungsi yang "menempel" ke struct](#6-method--fungsi-yang-menempel-ke-struct)
7. [Pointer (`*`) vs Value](#7-pointer--vs-value)
8. [Error handling — `if err != nil`](#8-error-handling--if-err--nil)
9. [Alur lengkap 1 request: Route → Controller → Service → Repository](#9-alur-lengkap-1-request-route--controller--service--repository)
10. [Registry pattern — dependency injection ala proyek ini](#10-registry-pattern--dependency-injection-ala-proyek-ini)
11. [Middleware — kode yang jalan sebelum handler](#11-middleware--kode-yang-jalan-sebelum-handler)
12. [`context.Context` — apa itu dan kenapa selalu ada di parameter pertama](#12-contextcontext--apa-itu-dan-kenapa-selalu-ada-di-parameter-pertama)
13. [`main.go` — titik awal semua ini dirakit](#13-maingo--titik-awal-semua-ini-dirakit)
14. [Rangkuman peta istilah](#14-rangkuman-peta-istilah)

---

## 1. Package — cara Go mengelompokkan kode

Setiap file Go diawali dengan `package <nama>`. Semua file dalam folder yang sama biasanya satu package, dan cuma yang **namanya diawali huruf besar** yang bisa dipakai dari luar package (istilahnya "exported").

Contoh dari [`services/osce/station.go`](../services/osce/station.go):

```go
// Package osce adalah lapisan business logic untuk domain OSCE station.
package osce

import (
	"context"
	"laundry-service/common/response"
	dtos "laundry-service/domain/dto/osce"
	"laundry-service/repositories"
)
```

Perhatikan: folder `services/osce` dan `repositories/osce` **sama-sama** bernama package `osce`. Ini boleh saja karena mereka ada di path import yang berbeda (`laundry-service/services/osce` vs `laundry-service/repositories/osce`). Saat dipakai bareng di file lain, salah satunya harus diberi alias supaya tidak bentrok — itu sebabnya di `routes/osce/station.go` domain DTO diimpor sebagai:

```go
import (
	"laundry-service/domain/dto/osce" // dipakai langsung sebagai osce.TipeStationForm
)
```

sedangkan di `services/osce/station.go` DTO yang sama diimpor dengan alias `dtos`:

```go
dtos "laundry-service/domain/dto/osce" // dipakai sebagai dtos.TipeStationForm
```

**Intinya:** nama package = nama folder secara konvensi, tapi path import yang membedakan satu package dengan package lain.

---

## 2. Struct — pengganti "class" di Go

Go tidak punya `class` seperti Java/PHP. Sebagai gantinya dipakai `struct` — kumpulan field, tanpa constructor bawaan, tanpa inheritance.

Contoh DTO (Data Transfer Object) di [`domain/dto/osce/station.go`](../domain/dto/osce/station.go):

```go
type StationForm struct {
	IDTipeStation string  `json:"id_tipe_station" form:"id_tipe_station" validate:"required,uuid4"`
	KodeStation   string  `json:"kode_station" form:"kode_station" validate:"required"`
	NamaStation   string  `json:"nama_station" form:"nama_station" validate:"required"`
	Deskripsi     string  `json:"deskripsi" form:"deskripsi"`
	DurasiDefault int     `json:"durasi_default" form:"durasi_default" validate:"required,min=1"`
	Bobot         float64 `json:"bobot" form:"bobot" validate:"required,gt=0"`
}
```

Ini adalah "bentuk" data yang dikirim client saat membuat station baru (misalnya lewat `POST /osce/stations`). Setiap field punya tipe data eksplisit (`string`, `int`, `float64`) — Go bahasa yang **statically typed**, jadi tipe data harus jelas sejak awal, tidak seperti JavaScript/PHP.

Struct juga dipakai untuk membungkus dependency (bukan cuma data), lihat bagian [Registry pattern](#10-registry-pattern--dependency-injection-ala-proyek-ini) di bawah.

---

## 3. Struct tag — anotasi di belakang field

Teks dalam backtick setelah tipe data (contoh: `` `json:"kode_station" form:"kode_station" validate:"required"` ``) disebut **struct tag**. Ini bukan komentar — ini metadata yang dibaca library lain lewat reflection.

Di proyek ini ada 3 tag yang sering dipakai:

| Tag | Dipakai oleh | Fungsi |
|-----|--------------|--------|
| `json:"..."` | `encoding/json` (built-in Go) | Nama field saat struct di-encode/decode ke/dari JSON |
| `form:"..."` | Fiber (`c.BodyParser`) | Nama field saat data datang dari form-data / query string |
| `validate:"..."` | `go-playground/validator` | Aturan validasi (wajib diisi, minimal, format UUID, dst) |

Jadi field `KodeStation` di Go akan:
- Dibaca dari request body sebagai `"kode_station"` (json) ATAU dari form field `kode_station`
- Divalidasi wajib ada (`required`) sebelum diproses

Kalau tag `validate` dilanggar (misal `kode_station` dikosongkan), middleware [`ValidateForm`](#11-middleware--kode-yang-jalan-sebelum-handler) otomatis menolak request dengan pesan error dalam Bahasa Indonesia (lihat [`middlewares/validation.go`](../middlewares/validation.go)) sebelum request sampai ke controller.

---

## 4. Interface — kontrak, bukan implementasi

Interface di Go adalah daftar method yang harus dipunyai suatu tipe, **tanpa** tipe itu harus secara eksplisit bilang "saya implement interface X" (beda dengan `implements` di Java). Kalau struct punya semua method yang diminta interface, otomatis dianggap "memenuhi" interface itu.

Contoh dari [`services/osce/station.go`](../services/osce/station.go):

```go
// IOsceStationService adalah kontrak business logic domain station.
type IOsceStationService interface {
	GetTipeStation(ctx context.Context, limit, offset int, order, filter string) response.Response
	GetTipeStationByID(ctx context.Context, id string) response.Response
	CreateTipeStation(ctx context.Context, form dtos.TipeStationForm) response.Response
	UpdateTipeStation(ctx context.Context, id string, form dtos.TipeStationForm) response.Response
	DeleteTipeStation(ctx context.Context, id string) response.Response
	// ... method lain
}
```

`IOsceStationService` cuma daftar **janji** — "siapa pun yang mengklaim dirinya `IOsceStationService`, harus punya method `GetTipeStation`, `CreateTipeStation`, dst dengan signature persis seperti itu."

Yang benar-benar mengimplementasikan janji itu adalah struct `OsceStationService`:

```go
type OsceStationService struct {
	repo repositories.IRepositoryRegistry
}

func (s *OsceStationService) GetTipeStation(ctx context.Context, limit, offset int, order, filter string) response.Response {
	// ...isi implementasi asli
}
```

**Kenapa dipakai interface di sini?**
1. **Controller tidak perlu tahu detail implementasi Service** — dia cukup pegang tipe `IOsceStationService`, panggil method-nya. Kalau nanti implementasi service diganti (misal untuk testing pakai mock), controller tidak perlu diubah sama sekali.
2. **Semua layer di proyek ini (Repository, Service, Controller, Route) memakai pola yang sama**: struct konkret + interface kontrak. Ini konsisten di seluruh domain — auth, user, groupakses, osce, dst.

---

## 5. Constructor — fungsi `New...`

Go tidak punya keyword `new ClassName()`. Sebagai konvensi, hampir semua struct di proyek ini dibuat lewat fungsi bernama `New<NamaStruct>` yang **mengembalikan interface**, bukan struct langsung.

```go
// NewOsceStationService membuat instance baru.
func NewOsceStationService(repo repositories.IRepositoryRegistry) IOsceStationService {
	return &OsceStationService{repo: repo}
}
```

Baca baris ini pelan-pelan:
- Menerima `repo` (dependency-nya) sebagai parameter
- Membuat `&OsceStationService{repo: repo}` — struct baru dengan field `repo` diisi, `&` berarti ambil alamat memorinya (pointer, lihat bagian 7)
- Return type-nya `IOsceStationService` (interface), **bukan** `*OsceStationService` (struct pointer) — inilah trik supaya pemanggil constructor hanya bisa memakai method yang ada di interface, tidak bisa mengintip/mengubah field internal struct

Pola yang sama berulang di setiap layer:

```go
// Repository layer — repositories/osce/station.go
func NewOsceStationRepository(db *gorm.DB) IOsceStationRepository {
	return &OsceStationRepository{db: db}
}

// Controller layer — controllers/osce/station.go
func NewOsceStationController(service services.IServiceRegistry) IOsceStationController {
	return &OsceStationController{service: service}
}

// Route layer — routes/osce/station.go
func NewOsceStationRoute(controller controllers.IControllerRegistry, router fiber.Router) IOsceStationRoute {
	return &OsceStationRoute{controller: controller, router: router}
}
```

Semua constructor menerima "layer di bawahnya" sebagai dependency: Route butuh Controller, Controller butuh Service, Service butuh Repository. Ini yang disebut **dependency injection manual** — tidak ada framework DI magic, semuanya di-wire tangan lewat parameter constructor (lihat [bagian 10](#10-registry-pattern--dependency-injection-ala-proyek-ini) dan [`main.go`](#13-maingo--titik-awal-semua-ini-dirakit)).

---

## 6. Method — fungsi yang "menempel" ke struct

Fungsi biasa di Go:

```go
func TambahDua(a int) int {
	return a + 2
}
```

Method adalah fungsi dengan **receiver** — parameter khusus sebelum nama fungsi yang menentukan struct mana yang "punya" fungsi ini:

```go
func (s *OsceStationService) GetTipeStation(ctx context.Context, limit, offset int, order, filter string) response.Response {
	data, err := s.repo.GetOsceStation().FindTipeStation(ctx, limit, offset, order, filter)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil data: " + err.Error()}
	}
	count, _ := s.repo.GetOsceStation().CountTipeStation(ctx, filter)
	return response.Response{
		Success: true, Message: "sukses",
		Data: map[string]interface{}{"data": data, "total_data": count},
	}
}
```

`(s *OsceStationService)` di sini artinya: method `GetTipeStation` ini "menempel" pada tipe `*OsceStationService`, dan di dalam method itu `s` dipakai untuk mengakses field struct (`s.repo`). Ini kurang lebih setara dengan `this`/`self` di bahasa lain.

---

## 7. Pointer (`*`) vs Value

Tanda `*` di depan tipe (`*OsceStationService`, `*gorm.DB`) berarti **pointer** — variabel yang isinya alamat memori, bukan salinan data itu sendiri.

```go
type OsceStationRepository struct {
	db *gorm.DB // pointer ke koneksi database, BUKAN salinan koneksinya
}
```

Kenapa pakai pointer di sini?
- **`*gorm.DB`**: koneksi database harus dipakai bersama (shared) oleh semua repository, bukan disalin satu-satu. Kalau tanpa `*`, tiap repository akan punya "kloningan" struct `gorm.DB` sendiri-sendiri — salah dan boros.
- **`&OsceStationService{...}`** di constructor: dengan mengembalikan pointer, method yang punya receiver `*OsceStationService` bisa mengubah field struct aslinya (bukan salinannya) kalau suatu saat dibutuhkan.

Aturan sederhana yang dipakai konsisten di proyek ini: **struct yang menyimpan dependency (koneksi DB, registry lain) selalu diakses lewat pointer.** DTO/form (`StationForm`, dst) biasanya dipakai sebagai *value* karena cuma data sekali pakai, bukan sesuatu yang dibagi.

---

## 8. Error handling — `if err != nil`

Go tidak punya `try/catch`. Fungsi yang bisa gagal biasanya mengembalikan dua nilai: hasil dan error, lalu pemanggil **wajib** mengecek error secara eksplisit.

```go
func (r *OsceStationRepository) InsertStation(ctx context.Context, form osce.StationForm) (string, error) {
	var id string
	result := r.db.WithContext(ctx).Raw(
		`INSERT INTO osce.mst_station (id_tipe_station, kode_station, nama_station, deskripsi, durasi_default, bobot)
		 VALUES (?, ?, ?, ?, ?, ?) RETURNING id_station`,
		form.IDTipeStation, form.KodeStation, form.NamaStation, form.Deskripsi, form.DurasiDefault, form.Bobot,
	).Scan(&id)
	return id, result.Error
}
```

Fungsi ini mengembalikan `(string, error)` — id station baru dan error kalau query gagal. Di layer Service, error itu wajib dicek sebelum dipakai:

```go
func (s *OsceStationService) CreateStation(ctx context.Context, form dtos.StationForm) response.Response {
	id, err := s.repo.GetOsceStation().InsertStation(ctx, form)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal membuat station: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Station berhasil dibuat", Data: map[string]string{"id": id}}
}
```

Pola `if err != nil { return ... }` ini akan kamu lihat berulang-ulang di **hampir setiap fungsi** di proyek ini. Ini gaya penulisan Go yang idiomatik — error ditangani sedini mungkin, bukan dibiarkan menjalar sampai runtime panic.

---

## 9. Alur lengkap 1 request: Route → Controller → Service → Repository

Proyek ini memakai **clean architecture** 4 lapis. Mari ikuti 1 request nyata: **`POST /osce/stations`** — membuat station baru.

```
Client (Postman/Frontend)
      │  POST /osce/stations  { "id_tipe_station": "...", "kode_station": "STA01", ... }
      ▼
┌─────────────────────────────────────────────────────────────┐
│ 1. ROUTE (routes/osce/station.go)                            │
│    Mendaftarkan path + middleware + handler mana yang dipanggil│
└─────────────────────────────────────────────────────────────┘
      │
      ▼
┌─────────────────────────────────────────────────────────────┐
│ 2. MIDDLEWARE (middlewares/validation.go)                    │
│    Parse body JSON → validasi struct tag `validate:"..."`     │
│    Kalau gagal → langsung balas 400, TIDAK lanjut ke controller│
└─────────────────────────────────────────────────────────────┘
      │  form yang sudah tervalidasi disimpan di c.Locals("validatedForm")
      ▼
┌─────────────────────────────────────────────────────────────┐
│ 3. CONTROLLER (controllers/osce/station.go)                  │
│    Ambil form dari Locals, panggil Service, kembalikan JSON   │
└─────────────────────────────────────────────────────────────┘
      │
      ▼
┌─────────────────────────────────────────────────────────────┐
│ 4. SERVICE (services/osce/station.go)                        │
│    Business logic (di sini sederhana: langsung teruskan)      │
│    Panggil Repository, bungkus hasil jadi response.Response   │
└─────────────────────────────────────────────────────────────┘
      │
      ▼
┌─────────────────────────────────────────────────────────────┐
│ 5. REPOSITORY (repositories/osce/station.go)                 │
│    Query SQL mentah lewat GORM, langsung sentuh database      │
└─────────────────────────────────────────────────────────────┘
      │
      ▼
   PostgreSQL (schema osce.*)
```

Sekarang lihat kode asli tiap lapisnya, dari bawah ke atas:

**5. Repository** — [`repositories/osce/station.go:150`](../repositories/osce/station.go)
```go
func (r *OsceStationRepository) InsertStation(ctx context.Context, form osce.StationForm) (string, error) {
	var id string
	result := r.db.WithContext(ctx).Raw(
		`INSERT INTO osce.mst_station (...) VALUES (?, ?, ?, ?, ?, ?) RETURNING id_station`,
		form.IDTipeStation, form.KodeStation, form.NamaStation, form.Deskripsi, form.DurasiDefault, form.Bobot,
	).Scan(&id)
	return id, result.Error
}
```

**4. Service** — [`services/osce/station.go:116`](../services/osce/station.go)
```go
func (s *OsceStationService) CreateStation(ctx context.Context, form dtos.StationForm) response.Response {
	id, err := s.repo.GetOsceStation().InsertStation(ctx, form)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal membuat station: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Station berhasil dibuat", Data: map[string]string{"id": id}}
}
```

**3. Controller** — [`controllers/osce/station.go:101`](../controllers/osce/station.go)
```go
func (ctrl *OsceStationController) CreateStation(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*osce.StationForm) // ambil form yang sudah divalidasi middleware
	data := ctrl.service.GetOsceStation().CreateStation(c.UserContext(), *form)
	return c.JSON(data)
}
```

**2. Middleware** dipasang langsung di baris route (bukan file terpisah untuk tiap endpoint):

**1. Route** — [`routes/osce/station.go:42`](../routes/osce/station.go)
```go
api.Post("/stations", middleware.ValidateForm(&osce.StationForm{}), r.controller.GetOsceStation().CreateStation)
```

Baris ini artinya: "kalau ada `POST /osce/stations`, jalankan dulu `ValidateForm(&osce.StationForm{})` (parse + validasi), baru kalau lolos, panggil `CreateStation` di controller."

**Detail kecil yang sering membingungkan pemula:** `form := c.Locals("validatedForm").(*osce.StationForm)` itu disebut **type assertion**. `c.Locals(...)` mengembalikan tipe `interface{}` (bisa apa saja), jadi harus "dipaksa" jadi tipe konkret `*osce.StationForm` sebelum field-nya bisa diakses.

---

## 10. Registry pattern — dependency injection ala proyek ini

Karena Go tidak punya container DI otomatis (seperti Spring di Java), proyek ini pakai pola manual bernama **Registry**: satu struct pusat per layer yang tahu cara membuat semua sub-komponen di layer itu.

Contoh [`repositories/registry.go`](../repositories/registry.go):

```go
type Registry struct {
	db          *gorm.DB
	dbAkademik  *gorm.DB
	dbDigiclass *gorm.DB
}

type IRepositoryRegistry interface {
	GetRef() refRepo.IRefRepository
	GetUser() userRepo.IUserRepository
	// ...
	GetOsceStation() osceRepo.IOsceStationRepository
	GetOsceExam() osceRepo.IOsceExamRepository
}

func NewRepositoryRegistry(db, dbAkademik, dbDigiclass *gorm.DB) IRepositoryRegistry {
	return &Registry{db: db, dbAkademik: dbAkademik, dbDigiclass: dbDigiclass}
}

func (r *Registry) GetOsceStation() osceRepo.IOsceStationRepository {
	return osceRepo.NewOsceStationRepository(r.db)
}
```

Jadi kalau Controller butuh Repository OSCE Station, dia tidak pernah menulis `osceRepo.NewOsceStationRepository(db)` langsung — dia cukup panggil `registry.GetOsceStation()`. Registry yang tahu detail cara merakitnya (termasuk koneksi DB mana yang dipakai).

Ada **3 Registry berjenjang** di proyek ini, satu untuk tiap layer:

```go
repository := repositories.NewRepositoryRegistry(database.DB, database.DBAkademik, database.DBDigiclass)
service    := services.NewServiceRegistry(repository)   // Service registry butuh Repository registry
controller := controllers.NewControllerRegistry(service) // Controller registry butuh Service registry
routes.NewRouteRegistry(controller, app).Serve()          // Route registry butuh Controller registry
```

Baris ini ada di [`main.go`](#13-maingo--titik-awal-semua-ini-dirakit) — di sinilah 4 lapisan clean architecture benar-benar "disambung" jadi satu.

---

## 11. Middleware — kode yang jalan sebelum handler

Middleware di Fiber adalah fungsi `func(c *fiber.Ctx) error` yang dipasang sebelum handler asli, dan bisa memutuskan lanjut (`c.Next()`) atau berhenti (langsung `return` dengan response error).

Contoh [`middlewares/validation.go`](../middlewares/validation.go):

```go
func ValidateForm(form interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := c.BodyParser(form); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(types.Response{
				Success: false, Message: err.Error(), Data: nil,
			})
		}

		if err := validate.Struct(form); err != nil {
			errors := make(map[string]string)
			for _, err := range err.(validator.ValidationErrors) {
				errors[err.Field()] = getErrorMessage(err)
			}
			return c.Status(fiber.StatusBadRequest).JSON(types.Response{
				Success: false, Message: "Invalid input", Data: errors,
			})
		}

		c.Locals("validatedForm", form) // simpan supaya bisa diambil controller
		return c.Next()                 // lolos validasi → lanjut ke handler berikutnya
	}
}
```

`ValidateForm` sendiri bukan middleware — dia adalah **fungsi yang mengembalikan middleware** (`fiber.Handler`). Pola "fungsi yang menghasilkan fungsi" ini disebut **higher-order function**, dan itu sebabnya di route ditulis dengan tanda kurung ganda secara implisit:

```go
middleware.ValidateForm(&osce.StationForm{})
//         └─ dipanggil sekali saat Serve(), menghasilkan fiber.Handler
//            fiber.Handler itulah yang dijalankan Fiber tiap kali ada request masuk
```

---

## 12. `context.Context` — apa itu dan kenapa selalu ada di parameter pertama

Hampir semua method di Service dan Repository punya parameter pertama `ctx context.Context`:

```go
func (s *OsceStationService) GetTipeStation(ctx context.Context, limit, offset int, order, filter string) response.Response
func (r *OsceStationRepository) FindTipeStation(ctx context.Context, limit, offset int, order, filter string) ([]map[string]interface{}, error)
```

`context.Context` adalah "amplop" yang dibawa sepanjang 1 request — isinya bisa deadline/timeout, sinyal pembatalan (kalau client disconnect di tengah jalan), atau data request-scoped lain. Di proyek ini asalnya dari Fiber:

```go
data := ctrl.service.GetOsceStation().GetTipeStation(c.UserContext(), limit, offset, order, filter)
```

`c.UserContext()` mengambil context bawaan request Fiber itu, lalu context yang sama diteruskan turun sampai ke query database:

```go
query := r.db.WithContext(ctx).Table("osce.tipe_station")...
```

`db.WithContext(ctx)` artinya kalau request dibatalkan/timeout, query database juga ikut dibatalkan — tidak menggantung sia-sia di database. Konvensi Go: **`ctx` selalu jadi parameter pertama**, dan diteruskan (bukan dibuat baru) di setiap pemanggilan fungsi turunannya.

---

## 13. `main.go` — titik awal semua ini dirakit

[`main.go`](../main.go) adalah pintu masuk aplikasi (fungsi `func main()` wajib ada persis satu di package `main`). Di sinilah semua Registry di-construct berurutan:

```go
func main() {
	database.Connect()
	database.ConnectDigiclass()

	app := fiber.New()
	app.Use(cors.New(cors.Config{ /* ... */ }))

	// Dependency injection untuk domain yang sudah dimigrasi ke pola berlapis (registry).
	repository := repositories.NewRepositoryRegistry(database.DB, database.DBAkademik, database.DBDigiclass)
	service := services.NewServiceRegistry(repository)
	controller := controllers.NewControllerRegistry(service)
	routes.NewRouteRegistry(controller, app).Serve()

	go func() {
		app.Listen(":" + httpPort) // HTTP server (Fiber) jalan di goroutine terpisah
	}()

	// gRPC server jalan di goroutine utama
	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, &grpcsso.AuthServiceServer{})
	grpcServer.Serve(listener)
}
```

Poin penting:
- **4 baris DI** (`repository := ...` sampai `routes.NewRouteRegistry(...).Serve()`) adalah urutan wajib — tiap Registry butuh Registry sebelumnya sebagai parameter constructor, persis seperti dijelaskan di [bagian 10](#10-registry-pattern--dependency-injection-ala-proyek-ini).
- `go func() { ... }()` menjalankan HTTP server di **goroutine** — thread ringan bawaan Go — supaya `main()` bisa lanjut menjalankan gRPC server di baris setelahnya tanpa saling memblokir. `go` di sini adalah keyword bahasa (bukan nama package "Go"), bukan konsep lanjutan yang perlu dikuasai penuh di awal — cukup tahu artinya "jalankan fungsi ini secara paralel, jangan tunggu selesai."

---

## 14. Rangkuman peta istilah

| Istilah Go | Analogi bahasa lain | Contoh di proyek ini |
|---|---|---|
| `package` | namespace / module | `package osce` |
| `struct` | class (tanpa method bawaan) | `type OsceStationService struct { repo ... }` |
| struct tag | decorator/annotation | `` `json:"kode_station" validate:"required"` `` |
| `interface` | contract / abstract class | `IOsceStationService` |
| `func New...(...)` | constructor | `NewOsceStationService(repo)` |
| method (`func (s *X) F()`) | instance method | `func (s *OsceStationService) GetTipeStation(...)` |
| `*T` (pointer) | reference/pass-by-reference | `*gorm.DB`, `*OsceStationService` |
| `err != nil` | try/catch | `if err != nil { return ... }` |
| `context.Context` | request-scoped object | `ctx context.Context` di semua Service/Repository |
| `go func(){}()` | thread / async task | HTTP server jalan paralel dengan gRPC di `main.go` |
| Registry pattern | DI container manual | `repositories.NewRepositoryRegistry(...)` |
| Middleware | filter/interceptor | `middleware.ValidateForm(&osce.StationForm{})` |

---

## Kalau mau belajar lebih lanjut

Cara paling efektif untuk lanjut belajar dari proyek ini: buka satu domain OSCE lain (misalnya `services/osce/exam.go` atau `services/osce/assessment.go`), lalu coba telusuri sendiri alurnya dari `routes/osce/exam.go` sampai ke `repositories/osce/exam.go` — polanya persis sama dengan yang dijelaskan di dokumen ini, cuma nama domain dan query SQL-nya yang beda. Untuk gambaran arsitektur keseluruhan (bukan level bahasa Go, tapi level desain sistem), baca [`00-arsitektur.md`](00-arsitektur.md).
