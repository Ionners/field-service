package repositories

import (
	"context"
	"errors"
	errWrap "field-service/common/error"
	errConstant "field-service/constants/error"
	errTime "field-service/constants/error/time"
	"field-service/domain/models"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TimeRepository struct {
	db *gorm.DB
}

type ITimeRepository interface {
	FindAll(context.Context) ([]models.Time, error)
	FindByUUID(context.Context, string) (*models.Time, error)
	FindById(context.Context, int) (*models.Time, error)
	Create(context.Context, *models.Time) (*models.Time, error)
}

func NewTimeRepository(db *gorm.DB) ITimeRepository {
	return &TimeRepository{db: db}
}

func (t *TimeRepository) FindAll(ctx context.Context) ([]models.Time, error) {
	fmt.Println("🔍 [DEBUG-REPOSITORIES] Mengambil semua data waktu")
	var times []models.Time
	err := t.db.WithContext(ctx).Find(&times).Error
	if err != nil {
		fmt.Println("❌ [ERROR-REPOSITORIES] Gagal mengambil data waktu:", err)
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	fmt.Println("✅ [INFO-REPOSITORIES] Berhasil mengambil data waktu:", times)
	return times, nil
}

func (t *TimeRepository) FindByUUID(ctx context.Context, uuid string) (*models.Time, error) {
	var time models.Time
	err := t.db.WithContext(ctx).Where("uuid = ?", uuid).First(&time).Error
	if err != nil {
		fmt.Println("❌ [ERROR-REPOSITORIES] Gagal mengambil data waktu:", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("❌ [ERROR-REPOSITORIES] Data waktu tidak ditemukan")
			return nil, errWrap.WrapError(errTime.ErrTimeNotFound)
		}
		fmt.Println("❌ [ERROR-REPOSITORIES] Gagal mengambil data waktu:", err)
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	fmt.Println("✅ [INFO-REPOSITORIES] Berhasil mengambil data waktu:", time)
	return &time, nil
}

func (t *TimeRepository) FindById(ctx context.Context, id int) (*models.Time, error) {
	var time models.Time
	err := t.db.WithContext(ctx).Where("id = ?", id).First(&time).Error
	if err != nil {
		fmt.Println("❌ [ERROR-REPOSITORIES] Gagal mengambil data waktu:", err)
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	fmt.Println("✅ [INFO-REPOSITORIES] Berhasil mengambil data waktu:", time)
	return &time, nil
}

func (t *TimeRepository) Create(ctx context.Context, time *models.Time) (*models.Time, error) {
	time.UUID = uuid.New()
	err := t.db.WithContext(ctx).Create(time).Error
	if err != nil {
		fmt.Println("❌ [ERROR-REPOSITORIES] Gagal membuat data waktu:", err)
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	fmt.Println("✅ [INFO-REPOSITORIES] Berhasil membuat data waktu:", time)
	return time, nil
}
