package controllers

import (
	"field-service/common/response"
	"field-service/domain/dto"
	"field-service/services"
	"net/http"

	"github.com/gin-gonic/gin"
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

	// 🧲 Step 2: (DIHANDLE DI SERVICE) Ambil data dari body JSON dan bind ke struct request
	// 📌 Catatan: Binding dilakukan langsung di dalam service, bukan di controller
	// Pastikan service melakukan binding + validasi, karena controller tidak melakukannya

	// 🚀 Step 3: Panggil service untuk membuat data waktu
	result, err := t.service.GetTime().Create(c, &request)
	if err != nil {
		// 🛑 Step 4: Jika ada error, kirim response error
		response.HttpResponse(response.ParamHttpResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	// ✅ Step 5: Jika data waktu berhasil dibuat, kirim response sukses
	response.HttpResponse(response.ParamHttpResp{
		Code: http.StatusCreated,
		Gin:  c,
		Data: result,
	})
}
