package controllers

import (
	errValidation "field-service/common/error"
	"field-service/common/response"
	"field-service/domain/dto"
	"field-service/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type TimeController struct {
	service services.IServiceRegistry
}

type ITimeController interface {
	GetAll(*gin.Context)
	GetByUUID(*gin.Context)
	Create(*gin.Context)
}

func NewTimeController(service services.IServiceRegistry) ITimeController {
	return &TimeController{service: service}
}

func (t *TimeController) GetAll(c *gin.Context) {
	// 🚀 Step 1: Ambil semua data waktu dari service
	result, err := t.service.GetTime().GetAll(c)
	if err != nil {
		// 🛑 Step 2: Jika ada error, kirim response error
		response.HttpResponse(response.ParamHttpResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	// ✅ Step 3: Jika semua data waktu ditemukan, kirim response sukses
	response.HttpResponse(response.ParamHttpResp{
		Code: http.StatusOK,
		Gin:  c,
		Data: result,
	})
}

func (t *TimeController) GetByUUID(c *gin.Context) {
	// 🚀 Step 1: Ambil UUID dari parameter URL
	uuid := c.Param("uuid")

	// 🚀 Step 2: Ambil data waktu berdasarkan UUID dari service
	result, err := t.service.GetTime().GetByUUID(c, uuid)
	if err != nil {
		// 🛑 Step 3: Jika ada error, kirim response error
		response.HttpResponse(response.ParamHttpResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	// ✅ Step 4: Jika data waktu ditemukan, kirim response sukses
	response.HttpResponse(response.ParamHttpResp{
		Code: http.StatusOK,
		Gin:  c,
		Data: result,
	})
}

func (t *TimeController) Create(c *gin.Context) {
	// 🧾 Step 1: Siapkan struct request untuk menampung input dari client (body JSON)
	var request dto.TimeRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		// 🛑 Step 2: Jika ada error saat binding, kirim response error
		response.HttpResponse(response.ParamHttpResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	// 📜 Step 3: Validasi input menggunakan validator
	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		// 🛑 Step 4: Jika ada error saat validasi, kirim response error
		fmt.Println("❌ [ERROR-TIME-CONTROLLER] Gagal validasi input:", err)
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errValidation.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHttpResp{
			Code:    http.StatusBadRequest,
			Err:     err,
			Message: &errMessage,
			Data:    errResponse, // Detail kesalahan input
			Gin:     c,
		})
		return
	}

	// 🚀 Step 5: Kirim request ke service untuk membuat data waktu baru
	result, err := t.service.GetTime().Create(c, &request)
	if err != nil {
		// 🛑 Step 6: Jika ada error, kirim response error
		response.HttpResponse(response.ParamHttpResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	// ✅ Step 7: Jika data waktu berhasil dibuat, kirim response sukses
	response.HttpResponse(response.ParamHttpResp{
		Code: http.StatusCreated,
		Gin:  c,
		Data: result,
	})
}
