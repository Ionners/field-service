package services

import (
	"context"
	"field-service/common/util"
	"field-service/constants"
	errFieldSchedule "field-service/constants/error/fieldschedule"
	"field-service/domain/dto"
	"field-service/domain/models"
	"field-service/repositories"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type FieldScheduleService struct {
	repository repositories.IRepositoryRegistry
}

type IFieldScheduleService interface {
	GetAllWithPagination(context.Context, *dto.FieldScheduleRequestParam) (*util.PaginationResult, error)
	GetAllByFieldIDAndDate(context.Context, string, string) ([]dto.FieldScheduleForBookingResponse, error)
	GetByUUID(context.Context, string) (*dto.FieldScheduleResponse, error)
	GenerateScheduleForOneMonth(context.Context, *dto.GenerateFieldScheduleForOneMonthRequest) error
	Create(context.Context, *dto.FieldScheduleRequest) error
	Update(context.Context, string, *dto.UpdateFieldScheduleRequest) (*dto.FieldScheduleResponse, error)
	UpdateStatus(context.Context, *dto.UpdateStatusFieldScheduleRequest) error
	Delete(context.Context, string) error
}

func NewFieldScheduleService(repository repositories.IRepositoryRegistry) IFieldScheduleService {
	return &FieldScheduleService{repository: repository}
}

func (f *FieldScheduleService) GetAllWithPagination(
	ctx context.Context,
	param *dto.FieldScheduleRequestParam,
) (*util.PaginationResult, error) {
	// 🚀 [DEBUG-FIELD-SCHEDULE-SERVICE] Mulai function GetAllWithPagination
	fmt.Println("🚀 [DEBUG-FIELD-SCHEDULE-SERVICE] Start GetAllWithPagination")
	fmt.Printf("📍 [[DEBUG-FIELD-SCHEDULE-SERVICE] Repository nil? %v\n", f.repository == nil)
	fmt.Printf("📥 [DEBUG-FIELD-SCHEDULE-SERVICE] Param request: %+v\n", param)

	// 1️⃣ Ambil semua fieldSchedule pakai pagination (limit + page)
	fieldSchedules, total, err := f.repository.GetFieldSchedule().FindAllWithPagination(ctx, param)
	if err != nil {
		fmt.Println("❌ [ERROR-FIELD-SCHEDULE-SERVICE] Gagal ambil data fieldSchedule:", err)
		return nil, err
	}
	fmt.Printf("✅ [INFO-FIELD-SCHEDULE-SERVICE] Data fieldSchedule ditemukan: %+v (total: %d)\n", fieldSchedules, total)
	// 📝 Catatan:
	// Ambil data schedule + total datanya berapa semua.
	// Kalau error ambil datanya, hentikan proses.

	// 2️⃣ Siapkan tempat kosong (slice) buat tampung hasil response per item
	fieldSchedulesResults := make([]dto.FieldScheduleResponse, 0, len(fieldSchedules))
	// 📝 Catatan:
	// Kita siapkan "wadah kosong" untuk hasil akhir (data response).

	// 3️⃣ Loop tiap data schedule → ubah jadi bentuk response
	for _, schedule := range fieldSchedules {
		fieldSchedulesResults = append(fieldSchedulesResults, dto.FieldScheduleResponse{
			UUID:         schedule.UUID,
			FieldName:    schedule.Field.Name,
			Date:         schedule.Date.Format("2006-01-02"),
			PricePerHour: schedule.Field.PricePerHour,
			Status:       schedule.Status.GetStatusString(),
			Time:         fmt.Sprintf("%s - %s", schedule.Time.StartTime, schedule.Time.EndTime),
			CreatedAt:    schedule.CreatedAt,
			UpdatedAt:    schedule.UpdatedAt,
		})
		fmt.Printf("🔄 [DEBUG-FIELD-SCHEDULE-SERVICE] Proses schedule: %+v\n", schedule)
	}
	fmt.Printf("✅ [INFO-FIELD-SCHEDULE-SERVICE] Hasil response per schedule siap: %+v\n", fieldSchedulesResults)
	// 📝 Catatan:
	// Setiap data mentah kita ubah jadi response yang sudah rapi (format tanggal, status, harga, dll).

	// 4️⃣ Siapkan data pagination param (isi total count, limit, page + datanya)
	pagination := &util.PaginationParam{
		Count: total,
		Limit: param.Limit,
		Page:  param.Page,
		Data:  fieldSchedulesResults,
	}
	fmt.Printf("📦 [DEBUG-FIELD-SCHEDULE-SERVICE] Pagination param: %+v\n", pagination)
	// 📝 Catatan:
	// Bungkus semua hasil + info pagination jadi 1 objek PaginationParam.

	// 5️⃣ Generate response pagination
	response := util.GeneratePagination(*pagination)
	fmt.Printf("✅ [INFO-FIELD-SCHEDULE-SERVICE] Response Pagination siap: %+v\n", response)

	// 6️⃣ Return hasil response pagination
	return &response, nil
}

func (f *FieldScheduleService) convertMontName(inputDate string) string {
	date, err := time.Parse(time.DateOnly, inputDate)
	if err != nil {
		fmt.Println("❌ [ERROR-FIELD-SCHEDULE-SERVICE] convertMontName", err)
		return ""
	}

	indoensiaMonth := map[string]string{
		"Jan": "Januari",
		"Feb": "Februari",
		"Mar": "Maret",
		"Apr": "April",
		"May": "Mei",
		"Jun": "Juni",
		"Jul": "Juli",
		"Aug": "Agustus",
		"Sep": "September",
		"Oct": "Oktober",
		"Nov": "November",
		"Dec": "Desember",
	}

	formattedDate := date.Format("02 Jan")
	//ini ambil index ke 3 (dimualinya dari 0) yang artinya akan ambil
	day := formattedDate[:3]
	month := formattedDate[3:]
	formattedDate = fmt.Sprintf("%s %s", day, indoensiaMonth[month])
	return formattedDate
}

func (f *FieldScheduleService) GetAllByFieldIDAndDate(
	ctx context.Context, uuid string, date string) ([]dto.FieldScheduleForBookingResponse, error) {
	// 🚀 [DEBUG-SERVICE] Start function GetAllByFieldIDAndDate
	fmt.Println("🚀 [DEBUG-SERVICE] Start GetAllByFieldIDAndDate")
	fmt.Println("🔍 [DEBUG-SERVICE] UUID:", uuid)
	fmt.Println("🔍 [DEBUG-SERVICE] Date:", date)

	// 1️⃣ Cek apakah field (lapangan) dengan UUID itu ada?
	field, err := f.repository.GetField().FindByUUID(ctx, uuid)
	if err != nil {
		fmt.Println("❌ [ERROR-SERVICE] Gagal ambil field:", err)
		return nil, err
	}
	fmt.Println("✅ [INFO-SERVICE] Field ditemukan:", field)
	// 📝 Catatan:
	// Kita pastikan lapangan dengan UUID yang dikirim user itu memang ada di database.
	// Kalau tidak ada (error), hentikan proses.

	// 2️⃣ Ambil semua jadwal (schedules) untuk field ID dan tanggal yang diminta
	fieldSchedules, err := f.repository.GetFieldSchedule().FindAllByFieldIDAndDate(ctx, int(field.ID), date)
	if err != nil {
		fmt.Println("❌ [ERROR-SERVICE] Gagal ambil field schedules:", err)
		return nil, err
	}
	fmt.Println("✅ [INFO-SERVICE] Field schedules ditemukan:", fieldSchedules)
	// 📝 Catatan:
	// Kita ambil semua jadwal yang sesuai lapangan + tanggal.
	// Kalau gagal ambil (error), hentikan proses.

	// 3️⃣ Siapkan tempat (slice) untuk tampung hasil response
	fieldScheduleResult := make([]dto.FieldScheduleForBookingResponse, 0, len(fieldSchedules))
	// 📝 Catatan:
	// Kita bikin "wadah kosong" buat hasil response-nya nanti.

	// 4️⃣ Looping setiap schedule → proses dan isi respons
	for _, schedule := range fieldSchedules {
		pricePerHour := float64(schedule.Field.PricePerHour)

		fmt.Println("🔍 [DEBUG-SERVICE] Processing schedule:", schedule)

		fieldScheduleResult = append(fieldScheduleResult, dto.FieldScheduleForBookingResponse{
			UUID:         schedule.UUID,
			Date:         f.convertMontName(schedule.Date.Format(time.DateOnly)),
			Time:         schedule.Time.StartTime,
			Status:       schedule.Status.GetStatusString(),
			PricePerHour: util.RupiahFormat(&pricePerHour),
		})
	}
	// 📝 Catatan:
	// Untuk setiap jadwal yang ketemu, kita ubah ke bentuk response yang lebih rapi + format harga rupiah + format tanggal + status string.

	// 5️⃣ Return hasil akhirnya
	fmt.Println("✅ [INFO-SERVICE] FieldScheduleResult:", fieldScheduleResult)
	return fieldScheduleResult, nil
}

func (f *FieldScheduleService) GetByUUID(ctx context.Context, uuid string) (*dto.FieldScheduleResponse, error) {
	fmt.Println("🔍 [DEBUG-FIELD-SCHEDULE-SERVICE] UUID:", uuid)

	// 1️⃣ Ambil fieldSchedule by UUID
	fieldSchedule, err := f.repository.GetFieldSchedule().FindByUUID(ctx, uuid)
	if err != nil {
		fmt.Println("❌ [ERROR-FIELD-SCHEDULE-SERVICE] Gagal ambil fieldSchedule", err)
		return nil, err
	}
	fmt.Println("✅ [INFO-FIELD-SCHEDULE-SERVICE] FieldSchedule ditemukan:", fieldSchedule)

	response := dto.FieldScheduleResponse{
		UUID:         fieldSchedule.UUID,
		FieldName:    fieldSchedule.Field.Name,
		PricePerHour: fieldSchedule.Field.PricePerHour,
		Date:         fieldSchedule.Date.Format(time.DateOnly),
		Status:       fieldSchedule.Status.GetStatusString(),
		Time:         fmt.Sprintf("%s - %s", fieldSchedule.Time.StartTime, fieldSchedule.Time.EndTime),
		CreatedAt:    fieldSchedule.CreatedAt,
		UpdatedAt:    fieldSchedule.UpdatedAt,
	}

	fmt.Println("✅ [INFO-FIELD-SCHEDULE-SERVICE] Response yang dikembalikan:", response)
	return &response, nil
}

func (f *FieldScheduleService) GenerateScheduleForOneMonth(ctx context.Context, request *dto.GenerateFieldScheduleForOneMonthRequest) error {
	// 🚀 Start debug
	fmt.Println("🚀 [DEBUG-FIELD-SCHEDULE-SERVICE] GenerateScheduleForOneMonth - Start")
	fmt.Printf("📥 [DEBUG-FIELD-SCHEDULE-SERVICE] Input request: %+v\n", request)

	// ✅ Step 1: Cek field (lapangan) ada atau tidak
	field, err := f.repository.GetField().FindByUUID(ctx, request.FieldID)
	if err != nil {
		fmt.Println("❌ [ERROR-FIELD-SCHEDULE-SERVICE] Gagal ambil field:", err)
		return err
	}
	fmt.Printf("✅ [INFO-FIELD-SCHEDULE-SERVICE] Field ditemukan: %+v\n", field)

	// ✅ Step 2: Ambil semua time slotnya (jam nya)
	times, err := f.repository.GetTime().FindAll(ctx)
	if err != nil {
		fmt.Println("❌ [ERROR-FIELD-SCHEDULE-SERVICE] Gagal ambil time:", err)
		return err
	}
	fmt.Printf("✅ [INFO-FIELD-SCHEDULE-SERVICE] Time ditemukan: %+v\n", len(times))

	// ✅ Step 3: Tentukan jumlah hari (30 hari dari hari besok)
	numberOfDays := 30
	Now := time.Now().Add(time.Duration(1) * 24 * time.Hour) // mulai dari besok
	fmt.Printf("📆 [DEBUG-FIELD-SCHEDULE-SERVICE] Generate schedule mulai dari besok: %s untuk %d hari\n", Now.Format("2006-01-02"), numberOfDays)

	// ✅ Step 4: Buat wadah kosong untuk menampung daftar jadwal baru
	fieldSchedules := make([]models.FieldSchedule, 0, numberOfDays)
	fmt.Println("📦 [DEBUG-FIELD-SCHEDULE-SERVICE] Wadah kosong untuk jadwal sudah disiapkan")

	// 🔄 Step 5: Loop untuk semua tanggal (besok sampai 30 hari kedepan)
	for i := 0; i < numberOfDays; i++ {
		currentDate := Now.AddDate(0, 0, i)
		fmt.Printf("🔄 [DEBUG-FIELD-SCHEDULE-SERVICE] Tanggal yang diproses: %s\n", currentDate.Format("2006-01-02"))

		// 🔄 Step 6: Loop untuk semua time slot di setiap tanggal
		for _, item := range times {
			fmt.Printf("🔄 [DEBUG-FIELD-SCHEDULE-SERVICE] Proses TimeSlot: %s (TimeID: %d)\n", item.StartTime, item.ID)

			// 7️⃣ Step 7: Cek apakah jadwal sudah ada (hindari duplikat)
			shcedule, err := f.repository.GetFieldSchedule().FindByDateAndTimeID(
				ctx, currentDate.Format(time.DateOnly),
				int(item.ID), int(field.ID),
			)
			if err != nil {
				fmt.Println("❌ [ERROR-FIELD-SCHEDULE-SERVICE] Gagal ambil schedule:", err)
				return err
			}

			if shcedule != nil {
				fmt.Println("⚠️ [WARN-FIELD-SCHEDULE-SERVICE] Schedule sudah ada, tidak boleh duplikat")
				return errFieldSchedule.ErrFieldScheduleIsExist
			}

			// ➕ Step 8: Tambahkan schedule baru ke wadahnya
			fieldSchedules = append(fieldSchedules, models.FieldSchedule{
				UUID:    uuid.New(),
				FieldID: field.ID,
				TimeID:  item.ID,
				Date:    currentDate,
				Status:  constants.Available,
			})
		}
	}
	fmt.Printf("💾 [INFO-FIELD-SCHEDULE-SERVICE] Siap simpan %d schedule baru ke database\n", len(fieldSchedules))

	// 🗃️ Step 9: Simpan ke DB
	err = f.repository.GetFieldSchedule().Create(ctx, fieldSchedules)
	if err != nil {
		fmt.Println("❌ [ERROR-FIELD-SCHEDULE-SERVICE] Gagal simpan schedule:", err)
		return err
	}

	fmt.Println("✅ [INFO-FIELD-SCHEDULE-SERVICE] FieldSchedules berhasil disimpan")

	// 🏁 End debug
	return nil
}

func (f *FieldScheduleService) Create(ctx context.Context, request *dto.FieldScheduleRequest) error {
	// 🚀 Start debug
	fmt.Println("🚀 [DEBUG-FIELD-SCHEDULE-SERVICE] Create - Start")
	fmt.Printf("📥 [DEBUG-FIELD-SCHEDULE-SERVICE] Input request: %+v\n", request)

	// ✅ Step 1: Cek field (lapangan) ada atau tidak
	field, err := f.repository.GetField().FindByUUID(ctx, request.FieldID)
	if err != nil {
		fmt.Println("❌ [ERROR-FIELD-SCHEDULE-SERVICE] Gagal ambil field:", err)
		return err
	}
	fmt.Printf("✅ [INFO-FIELD-SCHEDULE-SERVICE] Field ditemukan: %+v\n", field)

	//Step 2: Buat wadah kosong untuk menampung daftar jadwal baru
	fieldSchedules := make([]models.FieldSchedule, 0, len(request.TimeIDs))
	fmt.Println("📦 [DEBUG-FIELD-SCHEDULE-SERVICE] Menyiapkan wadah kosong untuk kumpulkan jadwal (fieldSchedules)")

	dateParsed, _ := time.Parse(time.DateOnly, request.Date)
	fmt.Println("📆 [DEBUG-FIELD-SCHEDULE-SERVICE] Parsed date:", dateParsed)

	// 🔄 Step 3: Looping TimeIDs, cek, dan siapkan jadwal
	for _, timeId := range request.TimeIDs {
		fmt.Println("🔄 [DEBUG-FIELD-SCHEDULE-SERVICE] Loop TimeID:", timeId)

		// ⏰ Ambil scheduleTime
		scheduleTime, err := f.repository.GetTime().FindByUUID(ctx, timeId)
		if err != nil {
			fmt.Println("❌ [ERROR-FIELD-SCHEDULE-SERVICE] Gagal ambil scheduleTime/ jamnya tidak ketemu:", err)
			return err
		}
		fmt.Printf("✅ [INFO-FIELD-SCHEDULE-SERVICE] scheduleTime ditemukan: %+v\n", scheduleTime)

		// 🔍 Cek apakah schedule sudah ada (duplikat)
		schedule, err := f.repository.GetFieldSchedule().FindByDateAndTimeID(ctx, request.Date, int(scheduleTime.ID), int(field.ID))
		if err != nil {
			fmt.Println("❌ [ERROR-FIELD-SCHEDULE-SERVICE] Gagal cek schedule existing:", err)
			return err
		}
		if schedule != nil {
			fmt.Println("⚠️ [WARN-FIELD-SCHEDULE-SERVICE] Schedule sudah ada, tidak boleh duplikat")
			return errFieldSchedule.ErrFieldScheduleIsExist
		}

		// ➕ Tambahkan schedule baru ke slice
		fieldSchedules = append(fieldSchedules, models.FieldSchedule{
			UUID:    uuid.New(),
			FieldID: field.ID,
			Date:    dateParsed,
			TimeID:  scheduleTime.ID,
			Status:  constants.Available,
		})
		fmt.Printf("➕ [DEBUG-FIELD-SCHEDULE-SERVICE] Schedule baru ditambahkan: %+v\n", fieldSchedules[len(fieldSchedules)-1])
	}

	// 🗃️ Step 4: Simpan ke DB
	fmt.Printf("💾 [INFO-FIELD-SCHEDULE-SERVICE] Siap simpan %d schedule baru ke database\n", len(fieldSchedules))

	err = f.repository.GetFieldSchedule().Create(ctx, fieldSchedules)
	if err != nil {
		fmt.Println("❌ [ERROR-FIELD-SCHEDULE-SERVICE] Gagal simpan schedule:", err)
		return err
	}
	fmt.Println("✅ [INFO-FIELD-SCHEDULE-SERVICE] FieldSchedules berhasil disimpan")

	// 🏁 End debug
	fmt.Println("🏁 [DEBUG-FIELD-SCHEDULE-SERVICE] Create - End sukses")
	return nil
}

func (f *FieldScheduleService) Update(
	ctx context.Context,
	uuid string,
	request *dto.UpdateFieldScheduleRequest,
) (*dto.FieldScheduleResponse, error) {
	// 🚀 Start debug
	fmt.Println("🚀 [DEBUG-FIELD-SCHEDULE-SERVICE] Update - Start")
	fmt.Printf("📥 [DEBUG-FIELD-SCHEDULE-SERVICE] UUID: %s\n", uuid)
	fmt.Printf("📥 [DEBUG-FIELD-SCHEDULE-SERVICE] Input request: %+v\n", request)

	// ✅ Step 1: Ambil data fieldSchedule berdasarkan UUID
	fieldSchedule, err := f.repository.GetFieldSchedule().FindByUUID(ctx, uuid)
	if err != nil {
		fmt.Println("❌ [ERROR-FIELD-SCHEDULE-SERVICE] Gagal ambil data fieldSchedule:", err)
		return nil, err
	}

	// ✅ Step 2: Ambil data waktu berdasarkan UUID
	scheduleTime, err := f.repository.GetTime().FindByUUID(ctx, request.TimeID)
	if err != nil {
		fmt.Println("❌ [ERROR-FIELD-SCHEDULE-SERVICE] Gagal ambil data waktu:", err)
		return nil, err
	}
	fmt.Printf("✅ [INFO-FIELD-SCHEDULE-SERVICE] Data waktunya ditemukan: %+v\n", scheduleTime)

	// ✅ Step 3: Buat query untuk cek apakah jadwal berdasarkan tanggal dan timeID sudah ada/ belum
	isTimeExist, err := f.repository.GetFieldSchedule().FindByDateAndTimeID(
		ctx,
		request.Date,
		int(scheduleTime.ID),
		int(fieldSchedule.FieldID),
	)
	if err != nil {
		fmt.Println("❌ [ERROR-FIELD-SCHEDULE-SERVICE] Gagal cek jadwal:", err)
		return nil, err
	}

	if isTimeExist != nil && request.Date != fieldSchedule.Date.Format(time.DateOnly) {
		checkDate, err := f.repository.GetFieldSchedule().FindByDateAndTimeID(
			ctx,
			request.Date,
			int(scheduleTime.ID),
			int(fieldSchedule.FieldID),
		)
		if err != nil {
			fmt.Println("❌ [ERROR-FIELD-SCHEDULE-SERVICE] Gagal cek jadwal:", err)
			return nil, err
		}

		if checkDate != nil {
			fmt.Println("⚠️ [WARN-FIELD-SCHEDULE-SERVICE] Jadwal sudah ada pada tanggal yang baru")
			return nil, errFieldSchedule.ErrFieldScheduleIsExist
		}
	}

	dateParesed, _ := time.Parse(time.DateOnly, request.Date)
	fieldResult, err := f.repository.GetFieldSchedule().Update(ctx, uuid, &models.FieldSchedule{
		Date:   dateParesed,
		TimeID: scheduleTime.ID,
	})

	if err != nil {
		fmt.Println("❌ [ERROR-FIELD-SCHEDULE-SERVICE] Gagal update fieldSchedule:", err)
		return nil, err
	}

	// ✅ Step 4: Buat response
	response := dto.FieldScheduleResponse{
		UUID:         fieldResult.UUID,
		FieldName:    fieldResult.Field.Name,
		Date:         fieldResult.Date.Format(time.DateOnly),
		PricePerHour: fieldResult.Field.PricePerHour,
		Status:       fieldResult.Status.GetStatusString(),
		Time:         fmt.Sprintf("%s - %s", scheduleTime.StartTime, scheduleTime.EndTime),
		CreatedAt:    fieldResult.CreatedAt,
		UpdatedAt:    fieldResult.UpdatedAt,
	}

	fmt.Printf("✅ [INFO-FIELD-SCHEDULE-SERVICE] Response yang dikembalikan: %+v\n", response)
	fmt.Println("🏁 [DEBUG-FIELD-SCHEDULE-SERVICE] Update - End sukses")
	return &response, nil
}

func (f *FieldScheduleService) UpdateStatus(ctx context.Context, request *dto.UpdateStatusFieldScheduleRequest) error {
	// 🚀 [DEBUG-FIELD-SCHEDULE-SERVICE] Mulai function UpdateStatus
	fmt.Println("🚀 [DEBUG-FIELD-SCHEDULE-SERVICE] Start UpdateStatus")
	fmt.Printf("📥 [DEBUG-FIELD-SCHEDULE-SERVICE] Input request: %+v\n", request)

	// 1️⃣ Loop semua FieldScheduleIDs (karena bentuknya array/list)
	for _, item := range request.FieldScheduleIDs {
		fmt.Printf("🔄 [DEBUG-FIELD-SCHEDULE-SERVICE] Proses FieldScheduleID: %s\n", item)

		// 2️⃣ Cek apakah fieldSchedule dengan ID itu ada di database
		_, err := f.repository.GetFieldSchedule().FindByUUID(ctx, item)
		if err != nil {
			fmt.Printf("❌ [ERROR-FIELD-SCHEDULE-SERVICE] Data tidak ditemukan (ID: %s): %v\n", item, err)
			return fmt.Errorf("gagal update schedule dengan ID %s: %w", item, err)
		}
		fmt.Printf("✅ [INFO-FIELD-SCHEDULE-SERVICE] Data ditemukan, lanjut update: %s\n", item)

		// 3️⃣ Update status jadi booked
		err = f.repository.GetFieldSchedule().UpdateStatus(ctx, constants.Booked, item)
		if err != nil {
			fmt.Println("❌ [ERROR-FIELD-SCHEDULE-SERVICE] Gagal update status:", err)
			return err
		}
		fmt.Printf("✅ [INFO-FIELD-SCHEDULE-SERVICE] Status berhasil diupdate jadi booked: %s\n", item)
	}
	// 4️⃣ Selesai, return success kalau semua berhasil
	fmt.Println("✅ [INFO-FIELD-SCHEDULE-SERVICE] Semua status berhasil diupdate")
	fmt.Println("🏁 [DEBUG-FIELD-SCHEDULE-SERVICE] End UpdateStatus sukses")
	return nil
}

func (f *FieldScheduleService) Delete(ctx context.Context, uuid string) error {
	// 🚀 [DEBUG-FIELD-SCHEDULE-SERVICE] Mulai function Delete
	fmt.Println("🚀 [DEBUG-FIELD-SCHEDULE-SERVICE] Start Delete")
	fmt.Printf("📥 [DEBUG-FIELD-SCHEDULE-SERVICE] UUID: %s\n", uuid)

	// 1️⃣ Cek apakah fieldSchedule dengan UUID itu ada?
	_, err := f.repository.GetFieldSchedule().FindByUUID(ctx, uuid)
	if err != nil {
		fmt.Println("❌ [ERROR-FIELD-SCHEDULE-SERVICE] Gagal ambil fieldSchedule:", err)
		return err
	}
	fmt.Printf("✅ [INFO-FIELD-SCHEDULE-SERVICE] FieldSchedule ditemukan: %s\n", uuid)
	// 📝 Catatan:
	// Kita pastikan fieldSchedule dengan UUID yang dikirim user itu memang ada di database.
	// Kalau tidak ada (error), hentikan proses.
	// Kalau ada, lanjut ke langkah berikutnya.

	// 2️⃣ Hapus fieldSchedule berdasarkan UUID
	err = f.repository.GetFieldSchedule().Delete(ctx, uuid)
	if err != nil {
		fmt.Println("❌ [ERROR-FIELD-SCHEDULE-SERVICE] Gagal hapus fieldSchedule:", err)
		return err
	}
	fmt.Printf("✅ [INFO-FIELD-SCHEDULE-SERVICE] FieldSchedule berhasil dihapus: %s\n", uuid)

	// 3️⃣ Selesai, return success
	fmt.Println("🏁 [DEBUG-FIELD-SCHEDULE-SERVICE] End Delete sukses")
	return nil
}
