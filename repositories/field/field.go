package repositories

import (
	"context"
	"errors"
	errWrap "field-service/common/error"
	errConstant "field-service/constants/error"
	"field-service/domain/dto"
	"field-service/domain/models"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FieldRepository struct {
	db *gorm.DB
}

type IFieldRepository interface {
	FindAllWithPagination(context.Context, *dto.FieldRequestParam) ([]models.Field, int64, error)
	FindAllWithoutPagination(context.Context) ([]models.Field, error)
	FindByUUID(context.Context, string) (*models.Field, error)
	Create(context.Context, *models.Field) (*models.Field, error)
	Update(context.Context, string, *models.Field) (*models.Field, error)
	Delete(context.Context, string) error
}

func NewFieldRepository(db *gorm.DB) IFieldRepository {
	return &FieldRepository{db: db}
}

func (f *FieldRepository) FindAllWithPagination(
	ctx context.Context,
	param *dto.FieldRequestParam,
) ([]models.Field, int64, error) {
	var (
		fields []models.Field
		sort   string
		total  int64
	)

	fmt.Println("🔍 [DEBUG-REPOSITORIES] Sort Column:", *param.SortColumn)
	fmt.Println("🔍 [DEBUG-REPOSITORIES] Sort Order:", *param.SortOrder)
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
		Limit(limit).
		Offset(offset).
		Order(sort).
		Find(&fields).Error

	fmt.Println("🔍 [DEBUG-REPOSITORIES] Data field:", fields)

	if err != nil {
		fmt.Println("❌ [ERROR-REPOSITORIES] Gagal mengambil data field:", err)
		return nil, 0, errWrap.WrapError(errConstant.ErrSQLError)
	}

	err = f.db.
		WithContext(ctx).
		Model(&fields).
		Count(&total).
		Error

	fmt.Println("🔍 [DEBUG-REPOSITORIES] Total data field:", total)

	if err != nil {
		fmt.Println("❌ [ERROR-REPOSITORIES] Gagal menghitung total data field:", err)
		return nil, 0, errWrap.WrapError(errConstant.ErrSQLError)
	}

	fmt.Println("✅ [INFO-REPOSITORIES] Berhasil mengambil data field dengan total:", total)
	return fields, total, nil
}

func (f *FieldRepository) FindAllWithoutPagination(ctx context.Context) ([]models.Field, error) {
	var fields []models.Field
	fmt.Println("🔍 [DEBUG-REPOSITORIES] Mengambil semua data field")
	err := f.db.
		WithContext(ctx).
		Find(&fields).
		Error

	fmt.Println("🔍 [DEBUG-REPOSITORIES] Data field:", fields)
	if err != nil {
		fmt.Println("❌ [ERROR-REPOSITORIES] Gagal mengambil data field:", err)
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	fmt.Println("✅ [INFO-REPOSITORIES] Berhasil mengambil semua data field")
	return fields, nil
}

func (f *FieldRepository) FindByUUID(ctx context.Context, uuid string) (*models.Field, error) {
	var fields models.Field
	fmt.Println("🔍 [DEBUG-REPOSITORIES] Mengambil data field dengan UUID:", uuid)

	err := f.db.
		WithContext(ctx).
		Where("uuid = ?", uuid).
		First(&fields).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("❌ [ERROR-REPOSITORIES] Data field tidak ditemukan")
			return nil, errWrap.WrapError(errConstant.ErrSQLError)
		}
		fmt.Println("❌ [ERROR-REPOSITORIES] Gagal mengambil data field:", err)
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	fmt.Println("✅ [INFO-REPOSITORIES] Berhasil mengambil data field dengan UUID:", uuid)
	return &fields, nil
}

func (f *FieldRepository) Create(ctx context.Context, req *models.Field) (*models.Field, error) {
	fmt.Println("🔍 [DEBUG-REPOSITORIES] Membuat data field baru")
	field := models.Field{
		UUID:         uuid.New(),
		Code:         req.Code,
		Name:         req.Name,
		Images:       req.Images,
		PricePerHour: req.PricePerHour,
	}

	err := f.db.WithContext(ctx).Create(&field).Error
	if err != nil {
		fmt.Println("❌ [ERROR-REPOSITORIES] Gagal membuat data field:", err)
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	fmt.Println("✅ [INFO-REPOSITORIES] Berhasil membuat data field baru")
	return &field, nil
}

func (f *FieldRepository) Update(ctx context.Context, uuid string, req *models.Field) (*models.Field, error) {
	fmt.Println("🔍 [DEBUG-REPOSITORIES] Memperbarui data field dengan UUID:", uuid)
	field := models.Field{
		Code:         req.Code,
		Name:         req.Name,
		Images:       req.Images,
		PricePerHour: req.PricePerHour,
	}

	err := f.db.WithContext(ctx).Where("uuid = ?", uuid).Updates(&field).Error
	if err != nil {
		fmt.Println("❌ [ERROR-REPOSITORIES] Gagal memperbarui data field:", err)
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	fmt.Println("✅ [INFO-REPOSITORIES] Berhasil memperbarui data field dengan UUID:", uuid)
	return &field, nil
}

func (f *FieldRepository) Delete(ctx context.Context, uuid string) error {
	fmt.Println("🔍 [DEBUG-REPOSITORIES] Menghapus data field dengan UUID:", uuid)
	err := f.db.WithContext(ctx).Where("uuid = ?", uuid).Delete(&models.Field{}).Error
	if err != nil {
		fmt.Println("❌ [ERROR-REPOSITORIES] Gagal menghapus data field:", err)
		return errWrap.WrapError(errConstant.ErrSQLError)
	}

	fmt.Println("✅ [INFO-REPOSITORIES] Berhasil menghapus data field dengan UUID:", uuid)
	return nil
}
