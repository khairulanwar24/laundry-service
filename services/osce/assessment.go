// Package osce adalah lapisan business logic untuk domain OSCE penilaian.
package osce

import (
	"context"
	"math"
	"sso-service/common/response"
	dtos "sso-service/domain/dto/osce"
	"sso-service/repositories"
)

// OsceAssessmentService membungkus repository registry.
type OsceAssessmentService struct {
	repo repositories.IRepositoryRegistry
}

// IOsceAssessmentService adalah kontrak business logic domain penilaian.
type IOsceAssessmentService interface {
	// Penilaian
	GetFormPenilaian(ctx context.Context, idJadwal string) response.Response
	CreatePenilaian(ctx context.Context, form dtos.PenilaianForm) response.Response
	UpdatePenilaian(ctx context.Context, id string, form dtos.PenilaianUpdateForm) response.Response
	KirimPenilaian(ctx context.Context, id string) response.Response
	GetPenilaianByID(ctx context.Context, id string) response.Response

	// Hasil Ujian
	GetHasilUjian(ctx context.Context, idUjian string) response.Response
	GetHasilSaya(ctx context.Context, idUser string) response.Response
	GetHasilByID(ctx context.Context, id string) response.Response
	KalkulasiHasil(ctx context.Context, form dtos.KalkulasiHasilRequest) response.Response

	// Borderline
	KalkulasiBorderline(ctx context.Context, form dtos.BorderlineRequest) response.Response

	// Statistik
	GetStatistikUjian(ctx context.Context, idUjian string) response.Response

	// Jadwal Penguji
	GetJadwalPenguji(ctx context.Context, idPenguji string) response.Response
}

// NewOsceAssessmentService membuat instance baru.
func NewOsceAssessmentService(repo repositories.IRepositoryRegistry) IOsceAssessmentService {
	return &OsceAssessmentService{repo: repo}
}

// ===== PENILAIAN =====

func (s *OsceAssessmentService) GetFormPenilaian(ctx context.Context, idJadwal string) response.Response {
	// Ambil data jadwal + checklist
	jadwal, err := s.repo.GetOsceExam().FindJadwalByID(ctx, idJadwal)
	if err != nil {
		return response.Response{Success: false, Message: "Jadwal tidak ditemukan"}
	}

	// Cari id_station dari ujian_station
	idUjianStation, _ := jadwal["id_ujian_station"].(string)
	ujianStation, err := s.repo.GetOsceExam().FindUjianStationByID(ctx, idUjianStation)
	if err != nil {
		return response.Response{Success: false, Message: "Station ujian tidak ditemukan"}
	}

	idStation, _ := ujianStation["id_station"].(string)
	checklists, _ := s.repo.GetOsceStation().FindChecklistByStation(ctx, idStation)

	// Cek apakah sudah ada penilaian sebelumnya
	penilaianSebelumnya, _ := s.repo.GetOsceAssessment().FindPenilaianByJadwal(ctx, idJadwal)

	var detailPenilaian []map[string]interface{}
	if penilaianSebelumnya != nil {
		idPenilaian, _ := penilaianSebelumnya["id_penilaian"].(string)
		detailPenilaian, _ = s.repo.GetOsceAssessment().FindDetailPenilaian(ctx, idPenilaian)
	}

	return response.Response{
		Success: true,
		Message: "sukses",
		Data: map[string]interface{}{
			"jadwal":         jadwal,
			"ujian_station":  ujianStation,
			"checklists":     checklists,
			"penilaian":      penilaianSebelumnya,
			"details":        detailPenilaian,
		},
	}
}

func (s *OsceAssessmentService) CreatePenilaian(ctx context.Context, form dtos.PenilaianForm) response.Response {
	// Ambil checklist untuk station terkait
	jadwal, err := s.repo.GetOsceExam().FindJadwalByID(ctx, form.IDJadwalUjian)
	if err != nil {
		return response.Response{Success: false, Message: "Jadwal tidak ditemukan"}
	}

	idUjianStation, _ := jadwal["id_ujian_station"].(string)
	ujianStation, err := s.repo.GetOsceExam().FindUjianStationByID(ctx, idUjianStation)
	if err != nil {
		return response.Response{Success: false, Message: "Station ujian tidak ditemukan"}
	}

	idStation, _ := ujianStation["id_station"].(string)
	checklists, err := s.repo.GetOsceStation().FindChecklistByStation(ctx, idStation)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil checklist"}
	}

	// Hitung total skor & skor terbobot
	var totalSkor float64
	var totalBobot float64
	for _, cl := range checklists {
		skorMaks, _ := cl["skor_maksimal"].(int64)
		bobotItem, _ := cl["bobot_item"].(float64)
		totalBobot += float64(skorMaks) * bobotItem
	}

	for _, detail := range form.Details {
		totalSkor += float64(detail.Skor)
	}

	skorTerbobot := float64(0)
	if totalBobot > 0 {
		skorTerbobot = totalSkor / totalBobot * 100
	}

	// Insert penilaian
	id, err := s.repo.GetOsceAssessment().InsertPenilaian(ctx, form, totalSkor, skorTerbobot)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal menyimpan penilaian: " + err.Error()}
	}

	// Insert detail
	for _, detail := range form.Details {
		_ = s.repo.GetOsceAssessment().InsertDetailPenilaian(ctx, id, detail.IDChecklistStation, detail.Skor, detail.CatatanItem)
	}

	return response.Response{
		Success: true,
		Message: "Penilaian berhasil disimpan",
		Data:    map[string]string{"id_penilaian": id},
	}
}

func (s *OsceAssessmentService) UpdatePenilaian(ctx context.Context, id string, form dtos.PenilaianUpdateForm) response.Response {
	jadwal, err := s.repo.GetOsceExam().FindJadwalByID(ctx, form.IDJadwalUjian)
	if err != nil {
		return response.Response{Success: false, Message: "Jadwal tidak ditemukan"}
	}

	idUjianStation, _ := jadwal["id_ujian_station"].(string)
	ujianStation, _ := s.repo.GetOsceExam().FindUjianStationByID(ctx, idUjianStation)
	idStation, _ := ujianStation["id_station"].(string)
	checklists, _ := s.repo.GetOsceStation().FindChecklistByStation(ctx, idStation)

	var totalSkor float64
	var totalBobot float64
	for _, cl := range checklists {
		skorMaks, _ := cl["skor_maksimal"].(int64)
		bobotItem, _ := cl["bobot_item"].(float64)
		totalBobot += float64(skorMaks) * bobotItem
	}

	for _, detail := range form.Details {
		totalSkor += float64(detail.Skor)
	}

	skorTerbobot := float64(0)
	if totalBobot > 0 {
		skorTerbobot = totalSkor / totalBobot * 100
	}

	if err := s.repo.GetOsceAssessment().UpdatePenilaian(ctx, id, form, totalSkor, skorTerbobot); err != nil {
		return response.Response{Success: false, Message: "Gagal memperbarui: " + err.Error()}
	}

	// Update detail
	_ = s.repo.GetOsceAssessment().DeleteDetailPenilaian(ctx, id)
	for _, detail := range form.Details {
		_ = s.repo.GetOsceAssessment().InsertDetailPenilaian(ctx, id, detail.IDChecklistStation, detail.Skor, detail.CatatanItem)
	}

	return response.Response{Success: true, Message: "Penilaian berhasil diperbarui"}
}

func (s *OsceAssessmentService) KirimPenilaian(ctx context.Context, id string) response.Response {
	if err := s.repo.GetOsceAssessment().KirimPenilaian(ctx, id); err != nil {
		return response.Response{Success: false, Message: "Gagal mengirim penilaian: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Penilaian berhasil dikirim"}
}

func (s *OsceAssessmentService) GetPenilaianByID(ctx context.Context, id string) response.Response {
	data, err := s.repo.GetOsceAssessment().FindPenilaianByID(ctx, id)
	if err != nil {
		return response.Response{Success: false, Message: "Data tidak ditemukan"}
	}
	// Ambil detail
	details, _ := s.repo.GetOsceAssessment().FindDetailPenilaian(ctx, id)
	data["details"] = details
	return response.Response{Success: true, Message: "sukses", Data: data}
}

// ===== HASIL UJIAN =====

func (s *OsceAssessmentService) GetHasilUjian(ctx context.Context, idUjian string) response.Response {
	data, err := s.repo.GetOsceAssessment().FindHasilByUjian(ctx, idUjian)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil data: " + err.Error()}
	}
	return response.Response{Success: true, Message: "sukses", Data: data}
}

func (s *OsceAssessmentService) GetHasilSaya(ctx context.Context, idUser string) response.Response {
	// Cari semua hasil untuk user ini
	data, err := s.repo.GetOsceAssessment().FindHasilByUjianAndUser(ctx, "", idUser)
	if err != nil {
		return response.Response{Success: true, Message: "sukses", Data: []interface{}{}}
	}
	return response.Response{Success: true, Message: "sukses", Data: data}
}

func (s *OsceAssessmentService) GetHasilByID(ctx context.Context, id string) response.Response {
	data, err := s.repo.GetOsceAssessment().FindHasilByID(ctx, id)
	if err != nil {
		return response.Response{Success: false, Message: "Data tidak ditemukan"}
	}
	stations, _ := s.repo.GetOsceAssessment().FindHasilStation(ctx, id)
	data["stations"] = stations
	return response.Response{Success: true, Message: "sukses", Data: data}
}

func (s *OsceAssessmentService) KalkulasiHasil(ctx context.Context, form dtos.KalkulasiHasilRequest) response.Response {
	// Ambil semua ujian station
	ujianStations, err := s.repo.GetOsceExam().FindUjianStationByUjian(ctx, form.IDUjian)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil data station ujian"}
	}

	jadwals, err := s.repo.GetOsceExam().FindJadwalBySesi(ctx, "")
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil data jadwal"}
	}

	// Proses per mahasiswa
	type skorMahasiswa struct {
		totalSkor     float64
		totalMaksimal float64
		idUser        string
	}
	skorByUser := make(map[string]*skorMahasiswa)

	for _, j := range jadwals {
		idUser, _ := j["id_user"].(string)
		idJadwalUjian, _ := j["id_jadwal_ujian"].(string)
		idUjianStation, _ := j["id_ujian_station"].(string)

		// Cari penilaian untuk jadwal ini
		penilaian, err := s.repo.GetOsceAssessment().FindPenilaianByJadwal(ctx, idJadwalUjian)
		if err != nil || penilaian == nil {
			continue
		}

		skorTerbobot, _ := penilaian["skor_terbobot"].(float64)

		// Cari bobot station
		bobotStation := float64(1.0)
		for _, es := range ujianStations {
			if idES, _ := es["id_ujian_station"].(string); idES == idUjianStation {
				bobotStation, _ = es["bobot"].(float64)
				break
			}
		}

		if _, ok := skorByUser[idUser]; !ok {
			skorByUser[idUser] = &skorMahasiswa{idUser: idUser}
		}
		skorByUser[idUser].totalSkor += skorTerbobot * bobotStation
		skorByUser[idUser].totalMaksimal += 100 * bobotStation
	}

	// Simpan hasil
	batasLulus := float64(60.0)
	if form.BatasLulus != nil {
		batasLulus = *form.BatasLulus
	}
	diproses := 0
	for idUser, sc := range skorByUser {
		persentase := float64(0)
		if sc.totalMaksimal > 0 {
			persentase = sc.totalSkor / sc.totalMaksimal * 100
		}
		status := "tidak_lulus"
		if persentase >= batasLulus {
			status = "lulus"
		}
		_ = s.repo.GetOsceAssessment().UpsertHasilUjian(ctx, form.IDUjian, idUser, sc.totalSkor, persentase, batasLulus, status)
		diproses++
	}

	// Update juga per station hasil
	for idUser := range skorByUser {
		hasil, err := s.repo.GetOsceAssessment().FindHasilByUjianAndUser(ctx, form.IDUjian, idUser)
		if err != nil || hasil == nil {
			continue
		}
		idHasil, _ := hasil["id_hasil_ujian"].(string)
		for _, es := range ujianStations {
			idUjianStation, _ := es["id_ujian_station"].(string)
			for _, j := range jadwals {
				idJadwalUser, _ := j["id_user"].(string)
				idJadwalStation, _ := j["id_ujian_station"].(string)
				if idJadwalUser == idUser && idJadwalStation == idUjianStation {
					idJadwalUjian, _ := j["id_jadwal_ujian"].(string)
					penilaian, err := s.repo.GetOsceAssessment().FindPenilaianByJadwal(ctx, idJadwalUjian)
					if err == nil && penilaian != nil {
						skor, _ := penilaian["skor_terbobot"].(float64)
						penilaianGlobal, _ := penilaian["penilaian_global"].(string)
						bobotStation, _ := es["bobot"].(float64)
						_ = s.repo.GetOsceAssessment().UpsertHasilStation(ctx, idHasil, idUjianStation, skor, skor*bobotStation, penilaianGlobal)
					}
					break
				}
			}
		}
	}

	return response.Response{
		Success: true,
		Message: "Kalkulasi hasil selesai",
		Data:    map[string]interface{}{"total_mahasiswa_diproses": diproses},
	}
}

// ===== BORDERLINE REGRESSION =====

func (s *OsceAssessmentService) KalkulasiBorderline(ctx context.Context, form dtos.BorderlineRequest) response.Response {
	// Ambil semua penilaian yang sudah terkirim
	penilaians, err := s.repo.GetOsceAssessment().FindPenilaianUntukBorderline(ctx, form.IDUjian)
	if err != nil || len(penilaians) == 0 {
		return response.Response{Success: false, Message: "Belum ada penilaian yang dikirim"}
	}

	// Kelompokkan per station
	byStation := make(map[string][]map[string]interface{})
	for _, p := range penilaians {
		idUjianStation, _ := p["id_ujian_station"].(string)
		byStation[idUjianStation] = append(byStation[idUjianStation], p)
	}

	var results []dtos.BorderlineResultItem
	for idUjianStation, items := range byStation {
		// Hitung simple linear regression: skor_checklist = a * penilaian_global_borderline + b
		var sumX, sumY, sumXY, sumX2 float64
		n := float64(len(items))

		for _, item := range items {
			skor, _ := item["skor_terbobot"].(float64)
			penilaianGlobal, _ := item["penilaian_global"].(string)
			var x float64
			if penilaianGlobal == "borderline" {
				x = 1
			}

			sumX += x
			sumY += skor
			sumXY += x * skor
			sumX2 += x * x
		}

		// Slope (a) dan Intercept (b)
		var slope, intercept, cutOff, r2 float64
		if (n*sumX2 - sumX*sumX) != 0 {
			slope = (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)
			intercept = (sumY - slope*sumX) / n

			// Cut-off skor = intercept (nilai y saat x = 0, yaitu garis borderline)
			cutOff = math.Round(intercept*100) / 100

			// R-squared sederhana
			var ssRes, ssTot, meanY float64
			meanY = sumY / n
			for _, item := range items {
				skor, _ := item["skor_terbobot"].(float64)
				penilaianGlobal, _ := item["penilaian_global"].(string)
				var x float64
				if penilaianGlobal == "borderline" {
					x = 1
				}
				yPred := slope*x + intercept
				ssRes += (skor - yPred) * (skor - yPred)
				ssTot += (skor - meanY) * (skor - meanY)
			}
			if ssTot != 0 {
				r2 = math.Round((1-ssRes/ssTot)*10000) / 10000
			}
		}

		// Ambil info station
		kodeStation := ""
		namaStation := ""
		if len(items) > 0 {
			kodeStation, _ = items[0]["kode_station"].(string)
			namaStation, _ = items[0]["nama_station"].(string)
		}

		results = append(results, dtos.BorderlineResultItem{
			IDUjianStation: idUjianStation,
			KodeStation:    kodeStation,
			NamaStation:    namaStation,
			CutOffScore:    cutOff,
			R2:             r2,
			Slope:          math.Round(slope*10000) / 10000,
			Intercept:      math.Round(intercept*10000) / 10000,
		})
	}

	return response.Response{
		Success: true,
		Message: "Borderline regression selesai",
		Data: dtos.BorderlineResultResponse{
			IDUjian:   form.IDUjian,
			Results:   results,
		},
	}
}

// ===== STATISTIK =====

func (s *OsceAssessmentService) GetStatistikUjian(ctx context.Context, idUjian string) response.Response {
	data, err := s.repo.GetOsceAssessment().FindStatistikUjian(ctx, idUjian)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil statistik"}
	}
	return response.Response{Success: true, Message: "sukses", Data: data}
}

// ===== JADWAL PENGUJI =====

func (s *OsceAssessmentService) GetJadwalPenguji(ctx context.Context, idPenguji string) response.Response {
	data, err := s.repo.GetOsceExam().FindJadwalByPenguji(ctx, idPenguji)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil data: " + err.Error()}
	}
	return response.Response{Success: true, Message: "sukses", Data: data}
}
