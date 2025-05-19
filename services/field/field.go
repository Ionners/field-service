package services

import (
	"bytes"
	"context"
	"field-service/common/gcs"
	"field-service/common/util"
	errConstant "field-service/constants/error"
	"field-service/domain/dto"
	"field-service/domain/models"
	"field-service/repositories"
	"fmt"
	"io"
	"mime/multipart"
	"path"
	"time"
)

type FieldService struct {
	repository repositories.IRepositoryRegistry
	gcs        gcs.IGCSClient
}

type IFieldService interface {
	GetAllWithPagination(context.Context, *dto.FieldRequestParam) (*util.PaginationResult, error)
	GetAllWithoutPagination(context.Context) ([]dto.FieldResponse, error)
	GetByUUID(context.Context, string) (*dto.FieldResponse, error)
	Create(context.Context, *dto.FieldRequest) (*dto.FieldResponse, error)
	Update(context.Context, string, *dto.UpdateFieldRequest) (*dto.FieldResponse, error)
	Delete(context.Context, string) error
}

func NewFieldService(repository repositories.IRepositoryRegistry, gcs gcs.IGCSClient) IFieldService {
	return &FieldService{repository: repository, gcs: gcs}
}

func (f *FieldService) GetAllWithPagination(
	ctx context.Context,
	param *dto.FieldRequestParam,
) (*util.PaginationResult, error) {
	fmt.Println("🔍 [DEBUG-FIELD-SERVICE] GetAllWithPagination")
	fields, total, err := f.repository.GetField().FindAllWithPagination(ctx, param)
	if err != nil {
		fmt.Println("🔍 [DEBUG-FIELD-SERVICE] GetAllWithPagination", err)
		return nil, err
	}

	fieldResults := make([]dto.FieldResponse, 0)
	for _, field := range fields {
		fieldResults = append(fieldResults, dto.FieldResponse{
			UUID:         field.UUID,
			Code:         field.Code,
			Name:         field.Name,
			PricePerHour: field.PricePerHour,
			Images:       field.Images,
			CreatedAt:    field.CreatedAt,
			UpdatedAt:    field.UpdatedAt,
		})
		fmt.Println("🔍 [DEBUG-FIELD-SERVICE] GetAllWithPagination", fieldResults)
	}

	pagination := &util.PaginationParam{
		Count: total,
		Page:  param.Page,
		Limit: param.Limit,
		Data:  fieldResults,
	}
	fmt.Println("🔍 [DEBUG-FIELD-SERVICE] GetAllWithPagination", pagination)

	response := util.GeneratePagination(*pagination)
	fmt.Println("🔍 [DEBUG-FIELD-SERVICE] GetAllWithPagination", response)
	return &response, nil
}

func (f *FieldService) GetAllWithoutPagination(ctx context.Context) ([]dto.FieldResponse, error) {
	fields, err := f.repository.GetField().FindAllWithoutPagination(ctx)
	if err != nil {
		fmt.Println("🔍 [DEBUG-FIELD-SERVICE] GetAllWithoutPagination", err)
		return nil, err
	}

	fieldResults := make([]dto.FieldResponse, 0, len(fields))
	for _, field := range fields {
		fieldResults = append(fieldResults, dto.FieldResponse{
			UUID:         field.UUID,
			Name:         field.Name,
			PricePerHour: field.PricePerHour,
			Images:       field.Images,
		})
	}
	fmt.Println("🔍 [DEBUG-FIELD-SERVICE] GetAllWithoutPagination", fieldResults)
	return fieldResults, nil
}

func (f *FieldService) GetByUUID(ctx context.Context, uuid string) (*dto.FieldResponse, error) {
	fields, err := f.repository.GetField().FindByUUID(ctx, uuid)
	if err != nil {
		fmt.Println("🔍 [DEBUG-FIELD-SERVICE] GetByUUID", err)
		return nil, err
	}

	fieldResult := dto.FieldResponse{
		UUID:         fields.UUID,
		Code:         fields.Code,
		Name:         fields.Name,
		PricePerHour: fields.PricePerHour,
		Images:       fields.Images,
		CreatedAt:    fields.CreatedAt,
		UpdatedAt:    fields.UpdatedAt,
	}

	fmt.Println("🔍 [DEBUG-FIELD-SERVICE] GetByUUID", fieldResult)
	return &fieldResult, nil
}

func (f *FieldService) validateUpload(images []multipart.FileHeader) error {
	if len(images) == 0 {
		return errConstant.ErrInvalidUploadFile
	}

	for _, image := range images {
		if image.Size > 5*1024*1024 {
			fmt.Println("🔍 [DEBUG-FIELD-SERVICE] validateUpload", errConstant.ErrSizeTooBig)
			return errConstant.ErrSizeTooBig
		}
	}

	fmt.Println("🔍 [DEBUG-FIELD-SERVICE] validateUpload", "success")
	return nil
}

func (f *FieldService) processAndUploadImage(ctx context.Context, image multipart.FileHeader) (string, error) {
	file, err := image.Open()
	if err != nil {
		fmt.Println("🔍 [DEBUG-FIELD-SERVICE] processAndUploadImage", err)
		return "", err
	}
	defer file.Close()

	buffer := new(bytes.Buffer)
	_, err = io.Copy(buffer, file)
	if err != nil {
		fmt.Println("🔍 [DEBUG-FIELD-SERVICE] processAndUploadImage", err)
		return "", err
	}

	filename := fmt.Sprintf("images/%s-%s-%s", time.Now().Format("20060102150405"), image.Filename, path.Ext(image.Filename))
	url, err := f.gcs.UploadFile(ctx, filename, buffer.Bytes())
	if err != nil {
		fmt.Println("🔍 [DEBUG-FIELD-SERVICE] processAndUploadImage", err)
		return "", err
	}
	fmt.Println("🔍 [DEBUG-FIELD-SERVICE] processAndUploadImage", url)
	return url, nil
}

func (f *FieldService) uploadImage(ctx context.Context, images []multipart.FileHeader) ([]string, error) {
	err := f.validateUpload(images)
	if err != nil {
		fmt.Println("🔍 [DEBUG-FIELD-SERVICE] uploadImage", err)
		return nil, err
	}

	urls := make([]string, 0, len(images))
	for _, image := range images {
		url, err := f.processAndUploadImage(ctx, image)
		if err != nil {
			fmt.Println("🔍 [DEBUG-FIELD-SERVICE] uploadImage", err)
			return nil, err
		}
		fmt.Println("🔍 [DEBUG-FIELD-SERVICE] uploadImage", url)
		urls = append(urls, url)
	}

	fmt.Println("🔍 [DEBUG-FIELD-SERVICE] uploadImage", urls)
	return urls, nil
}

func (f *FieldService) Create(ctx context.Context, request *dto.FieldRequest) (*dto.FieldResponse, error) {
	imageUrl, err := f.uploadImage(ctx, request.Images)
	if err != nil {
		fmt.Println("🔍 [DEBUG-FIELD-SERVICE] Create", err)
		return nil, err
	}

	field, err := f.repository.GetField().Create(ctx, &models.Field{
		Code:         request.Code,
		Name:         request.Name,
		PricePerHour: request.PricePerHour,
		Images:       imageUrl,
	})
	if err != nil {
		return nil, err
	}

	response := dto.FieldResponse{
		UUID:         field.UUID,
		Code:         field.Code,
		Name:         field.Name,
		PricePerHour: field.PricePerHour,
		Images:       imageUrl,
		CreatedAt:    field.CreatedAt,
		UpdatedAt:    field.UpdatedAt,
	}

	fmt.Println("🔍 [DEBUG-FIELD-SERVICE] Create", response)
	return &response, nil
}

func (f *FieldService) Update(ctx context.Context, uuid string, req *dto.UpdateFieldRequest) (*dto.FieldResponse, error) {
	field, err := f.repository.GetField().FindByUUID(ctx, uuid)
	if err != nil {
		fmt.Println("🔍 [DEBUG-FIELD-SERVICE] Update", err)
		return nil, err
	}

	var imageUrl []string
	if req.Images == nil {
		imageUrl = field.Images
	} else {
		imageUrl, err = f.uploadImage(ctx, req.Images)
		if err != nil {
			fmt.Println("🔍 [DEBUG-FIELD-SERVICE] Update", err)
			return nil, err
		}
	}

	fieldResult, err := f.repository.GetField().Update(ctx, uuid, &models.Field{
		Code:         req.Code,
		Name:         req.Name,
		PricePerHour: req.PricePerHour,
		Images:       imageUrl,
	})
	if err != nil {
		return nil, err
	}

	return &dto.FieldResponse{
		UUID:         fieldResult.UUID,
		Code:         fieldResult.Code,
		Name:         fieldResult.Name,
		PricePerHour: fieldResult.PricePerHour,
		Images:       fieldResult.Images,
		CreatedAt:    fieldResult.CreatedAt,
		UpdatedAt:    fieldResult.UpdatedAt,
	}, nil
}

func (f *FieldService) Delete(ctx context.Context, uuid string) error {
	//cek dulu datanya ada atau tidak
	_, err := f.repository.GetField().FindByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	err = f.repository.GetField().Delete(ctx, uuid)
	if err != nil {
		fmt.Println("🔍 [DEBUG-FIELD-SERVICE] Delete", err)
		return err
	}

	fmt.Println("🔍 [DEBUG-FIELD-SERVICE] Delete", "success")
	return nil
}
