package repositories

import (
	"context"
	"errors"
	errWrap "field-service/common/error"
	"field-service/constants"
	errConstant "field-service/constants/error"
	errFieldSchedule "field-service/constants/error/fieldschedule"
	"field-service/domain/dto"
	"field-service/domain/models"
	"fmt"

	"gorm.io/gorm"
)

type FieldScheduleRepository struct {
	db *gorm.DB
}

type IFieldScheduleRepository interface {
	FindAllWithPagination(context.Context, *dto.FieldScheduleRequestParam) ([]models.FieldSchedule, int64, error)
	FindAllByFieldIDAndDate(context.Context, int, string) ([]models.FieldSchedule, error)
	FindByUUID(context.Context, string) (*models.FieldSchedule, error)
	FindByDateAndTimeID(context.Context, string, int, int) (*models.FieldSchedule, error)
	Create(context.Context, []models.FieldSchedule) error
	Update(context.Context, string, *models.FieldSchedule) (*models.FieldSchedule, error)
	UpdateStatus(context.Context, constants.FieldScheduleStatus, string) error
	Delete(context.Context, string) error
}

func NewFieldScheduleRepository(db *gorm.DB) IFieldScheduleRepository {
	return &FieldScheduleRepository{db: db}
}

func (f *FieldScheduleRepository) FindAllWithPagination(
	ctx context.Context,
	param *dto.FieldScheduleRequestParam,
) ([]models.FieldSchedule, int64, error) {
	var (
		fieldSchedules []models.FieldSchedule
		sort           string
		total          int64
	)

	fmt.Println("🔍 [DEBUG-REPOSITORIES] Sort Column:", param.SortColumn)
	fmt.Println("🔍 [DEBUG-REPOSITORIES] Sort Order:", param.SortOrder)
	if param.SortColumn != nil {
		fmt.Println("🔍 [DEBUG-REPOSITORIES] Sort Column ada, menggunakan:", *param.SortColumn, *param.SortOrder)
		sort = fmt.Sprintf("%s %s", *param.SortColumn, *param.SortOrder)
	} else {
		fmt.Println("🔍 [DEBUG-REPOSITORIES] Sort Column tidak ada, menggunakan default: created_at desc")
		sort = "created_at desc"
	}

	limit := param.Limit
	offset := (param.Page - 1) * limit
	err := f.db.
		WithContext(ctx).
		Preload("Field").
		Preload("Time").
		Limit(limit).
		Offset(offset).
		Order(sort).
		Find(&fieldSchedules).Error

	fmt.Println("🔍 [DEBUG-REPOSITORIES] Data field:", fieldSchedules)

	if err != nil {
		fmt.Println("❌ [ERROR-REPOSITORIES] Gagal mengambil data field:", err)
		return nil, 0, errWrap.WrapError(errConstant.ErrSQLError)
	}

	err = f.db.
		WithContext(ctx).
		Model(&fieldSchedules).
		Count(&total).
		Error

	fmt.Println("🔍 [DEBUG-REPOSITORIES] Total data field:", total)

	if err != nil {
		fmt.Println("❌ [ERROR-REPOSITORIES] Gagal menghitung total data field:", err)
		return nil, 0, errWrap.WrapError(errConstant.ErrSQLError)
	}

	fmt.Println("✅ [INFO-REPOSITORIES] Berhasil mengambil data field dengan total:", total)
	return fieldSchedules, total, nil
}

func (f *FieldScheduleRepository) FindAllByFieldIDAndDate(
	ctx context.Context,
	fieldID int, date string,
) ([]models.FieldSchedule, error) {
	var fieldSchedules []models.FieldSchedule
	fmt.Println("🔍 [DEBUG-REPOSITORIES] Mengambil semua data field")
	err := f.db.
		WithContext(ctx).
		Preload("Field").
		Preload("Time").
		Where("field_id = ?", fieldID).
		Where("date = ?", date).
		Joins("LEFT JOIN times on field_schedules.time_id = times.id").
		Order("times.start_time asc").
		Find(&fieldSchedules).
		Error

	fmt.Println("🔍 [DEBUG-REPOSITORIES] Data field:", fieldSchedules)
	if err != nil {
		fmt.Println("❌ [ERROR-REPOSITORIES] Gagal mengambil data field:", err)
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	fmt.Println("✅ [INFO-REPOSITORIES] Berhasil mengambil semua data field")
	return fieldSchedules, nil
}

func (f *FieldScheduleRepository) FindByUUID(ctx context.Context, uuid string) (*models.FieldSchedule, error) {
	var fieldSchedules models.FieldSchedule
	fmt.Println("🔍 [DEBUG-REPOSITORIES] Mengambil data field dengan UUID:", uuid)

	err := f.db.
		WithContext(ctx).
		Preload("Field").
		Preload("Time").
		Where("uuid = ?", uuid).
		First(&fieldSchedules).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("❌ [ERROR-REPOSITORIES] Data field tidak ditemukan")
			return nil, errWrap.WrapError(errFieldSchedule.ErrFieldScheduleNotFound)
		}
		fmt.Println("❌ [ERROR-REPOSITORIES] Gagal mengambil data field:", err)
		return nil, errWrap.WrapError((errConstant.ErrSQLError))
	}
	fmt.Println("✅ [INFO-REPOSITORIES] Berhasil mengambil data field dengan UUID:", uuid)
	return &fieldSchedules, nil
}

func (f *FieldScheduleRepository) FindByDateAndTimeID(
	ctx context.Context,
	date string,
	timeID int,
	fieldID int,
) (*models.FieldSchedule, error) {
	var fieldSchedules models.FieldSchedule
	fmt.Println("🔍 [DEBUG-REPOSITORIES] Mengambil data field dengan date:", date, "timeID:", timeID, "fieldID:", fieldID)

	err := f.db.
		WithContext(ctx).
		Preload("Field").
		Preload("Time").
		Where("date = ?", date).
		Where("time_id = ?", timeID).
		Where("field_id = ?", fieldID).
		First(&fieldSchedules).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("❌ [ERROR-REPOSITORIES] Data field tidak ditemukan")
			return nil, nil
		}
		fmt.Println("❌ [ERROR-REPOSITORIES] Gagal mengambil data field:", err)
		return nil, errWrap.WrapError((errConstant.ErrSQLError))
	}

	fmt.Println("✅ [INFO-REPOSITORIES] Berhasil mengambil data field dengan date:", date, "timeID:", timeID, "fieldID:", fieldID)
	return &fieldSchedules, nil
}

func (f *FieldScheduleRepository) Create(ctx context.Context, req []models.FieldSchedule) error {
	fmt.Println("🔍 [DEBUG-REPOSITORIES] Membuat data field schedule baru")
	err := f.db.WithContext(ctx).Create(&req).Error
	if err != nil {
		fmt.Println("❌ [ERROR-REPOSITORIES] Gagal membuat data field:", err)
		return errWrap.WrapError(errConstant.ErrSQLError)
	}

	fmt.Println("✅ [INFO-REPOSITORIES] Berhasil membuat data field baru")
	return nil
}

func (f *FieldScheduleRepository) Update(
	ctx context.Context,
	uuid string,
	req *models.FieldSchedule) (*models.FieldSchedule, error) {

	fieldSchedule, err := f.FindByUUID(ctx, uuid)
	if err != nil {
		fmt.Println("❌ [ERROR-REPOSITORIES] Gagal mengambil data field:", err)
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	fieldSchedule.Date = req.Date
	fmt.Println("🔍 [DEBUG-REPOSITORIES] Memperbarui data field dengan date:", fieldSchedule.Date)
	err = f.db.WithContext(ctx).Save(&fieldSchedule).Error
	if err != nil {
		fmt.Println("❌ [ERROR-REPOSITORIES] Gagal memperbarui data field:", err)
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	fmt.Println("✅ [INFO-REPOSITORIES] Berhasil memperbarui data field dengan date:", fieldSchedule.Date)
	return fieldSchedule, nil
}

func (f *FieldScheduleRepository) UpdateStatus(
	ctx context.Context,
	status constants.FieldScheduleStatus,
	uuid string,
) error {

	fieldSchedule, err := f.FindByUUID(ctx, uuid)
	if err != nil {
		fmt.Println("❌ [ERROR-REPOSITORIES] Gagal mengambil data field:", err)
		return err
	}

	fieldSchedule.Status = status
	fmt.Println("🔍 [DEBUG-REPOSITORIES] Memperbarui data field dengan date:", fieldSchedule.Status)
	err = f.db.WithContext(ctx).Save(&fieldSchedule).Error
	if err != nil {
		fmt.Println("❌ [ERROR-REPOSITORIES] Gagal memperbarui data field:", err)
		return errWrap.WrapError(errConstant.ErrSQLError)
	}

	fmt.Println("✅ [INFO-REPOSITORIES] Berhasil memperbarui data field dengan date:", fieldSchedule.Status)
	return nil
}

func (f *FieldScheduleRepository) Delete(ctx context.Context, uuid string) error {
	fmt.Println("🔍 [DEBUG-REPOSITORIES] Menghapus data field dengan UUID:", uuid)
	err := f.db.WithContext(ctx).Where("uuid = ?", uuid).Delete(&models.FieldSchedule{}).Error
	if err != nil {
		fmt.Println("❌ [ERROR-REPOSITORIES] Gagal menghapus data field:", err)
		return errWrap.WrapError(errConstant.ErrSQLError)
	}

	fmt.Println("✅ [INFO-REPOSITORIES] Berhasil menghapus data field dengan UUID:", uuid)
	return nil
}
