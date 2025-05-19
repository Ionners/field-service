package services

import (
	"context"
	"field-service/domain/dto"
	"field-service/domain/models"
	"field-service/repositories"
	"fmt"
)

type TimeService struct {
	repository repositories.IRepositoryRegistry
}

type ITimeService interface {
	GetAll(context.Context) ([]dto.TimeResponse, error)
	GetByUUID(context.Context, string) (*dto.TimeResponse, error)
	Create(context.Context, *dto.TimeRequest) (*dto.TimeResponse, error)
}

func NewTimeService(repository repositories.IRepositoryRegistry) ITimeService {
	return &TimeService{repository: repository}
}

func (t *TimeService) GetAll(ctx context.Context) ([]dto.TimeResponse, error) {
	// 🚀 Step 1: Ambil semua data time dari repository
	fmt.Println("🚀 [DEBUG-TIME-SERVICE] Mulai GetAll")
	times, err := t.repository.GetTime().FindAll(ctx)
	if err != nil {
		fmt.Println("❌ [ERROR-TIME-SERVICE] Gagal mengambil data waktu:", err)
		return nil, err
	}
	fmt.Println("✅ [INFO-TIME-SERVICE] Berhasil mengambil data waktu:", len(times))

	// 🛠️ Step 2: Siapkan slice kosong untuk hasil response
	timeResults := make([]dto.TimeResponse, 0, len(times))

	// 🔄 Step 3: Loop tiap data untuk diubah ke format DTO dan tambahkan ke response
	for _, time := range times {
		fmt.Printf("🔍 [DEBUG-TIME-SERVICE] Proses time: %+v\n", time)
		// Ubah data time ke format DTO
		// dan tambahkan ke slice hasil response
		timeResults = append(timeResults, dto.TimeResponse{
			UUID:      time.UUID,
			StartTime: time.StartTime,
			EndTime:   time.EndTime,
			CreatedAt: time.CreatedAt,
			UpdatedAt: time.UpdatedAt,
		})
	}

	// ✅ Step 4: Return hasil final
	fmt.Println("🏁 [INFO-TIME-SERVICE] GetAll selesai, return response")
	return timeResults, nil
}

func (t *TimeService) GetByUUID(ctx context.Context, uuid string) (*dto.TimeResponse, error) {
	// 🚀 Step 1: Debug input UUID
	fmt.Println("🔍 [DEBUG-TIME-SERVICE] GetByUUID dimulai")
	fmt.Println("🆔 [DEBUG-TIME-SERVICE] UUID:", uuid)

	// 🔎 Step 2: Ambil data dari repository berdasarkan UUID
	time, err := t.repository.GetTime().FindByUUID(ctx, uuid)
	if err != nil {
		fmt.Printf("❌ [ERROR-TIME-SERVICE] Gagal mengambil data waktu dengan UUID %s: %v\n", uuid, err)
		return nil, err
	}
	fmt.Printf("✅ [INFO-TIME-SERVICE] Berhasil mengambil data waktu dengan UUID %s: %+v\n", uuid, time)

	// 🛠️ Step 3: Siapkan hasil response
	// Ubah data time ke format DTO
	// dan kembalikan hasil response
	timeResult := dto.TimeResponse{
		UUID:      time.UUID,
		StartTime: time.StartTime,
		EndTime:   time.EndTime,
		CreatedAt: time.CreatedAt,
		UpdatedAt: time.UpdatedAt,
	}
	fmt.Printf("📦 [DEBUG-TIME-SERVICE] Hasil DTO: %+v\n", timeResult)

	// ✅ Step 4: Return hasil DTO
	fmt.Println("🏁 [INFO-TIME-SERVICE] Selesai GetByUUID, return response")
	return &timeResult, nil
}

func (t *TimeService) Create(ctx context.Context, req *dto.TimeRequest) (*dto.TimeResponse, error) {
	// 🚀 Step 1: Mulai proses & debug input
	fmt.Println("🚀 [DEBUG-TIME-SERVICE] Create dimulai")
	fmt.Printf("📥 [DEBUG-TIME-SERVICE] Input: %+v\n", req)

	// 🛠️ Step 2: Persiapkan data untuk disimpan (ambil dari request)
	time := &dto.TimeRequest{
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	// 💾 Step 3: Simpan ke database melalui repository
	timeResult, err := t.repository.GetTime().Create(ctx, &models.Time{
		StartTime: time.StartTime,
		EndTime:   time.EndTime,
	})
	if err != nil {
		fmt.Printf("❌ [ERROR-TIME-SERVICE] Gagal menyimpan data waktu: %v\n", err)
		return nil, err
	}
	fmt.Printf("✅ [INFO-TIME-SERVICE] Berhasil menyimpan data waktu: %+v\n", timeResult)

	// 🛠️ Step 4: Siapkan hasil response
	// Ubah hasil DB ke format DTO
	// dan kembalikan hasil response
	response := dto.TimeResponse{
		UUID:      timeResult.UUID,
		StartTime: timeResult.StartTime,
		EndTime:   timeResult.EndTime,
		CreatedAt: timeResult.CreatedAt,
		UpdatedAt: timeResult.UpdatedAt,
	}
	fmt.Printf("📤 [DEBUG-TIME-SERVICE] Response: %+v\n", response)

	// ✅ Step 5: Return hasil response
	fmt.Println("🏁 [INFO-TIME-SERVICE] Selesai Create, return response")
	return &response, nil
}
