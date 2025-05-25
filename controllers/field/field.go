package controllers

import (
	errValidation "field-service/common/error"
	"field-service/common/response"
	"field-service/domain/dto"
	"field-service/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type FieldController struct {
	service services.IServiceRegistry
}

type IFieldController interface {
	GetAllWithPagination(*gin.Context)
	GetAllWithoutPagination(*gin.Context)
	GetByUUID(*gin.Context)
	Create(*gin.Context)
	Update(*gin.Context)
	Delete(*gin.Context)
}

func NewFieldController(service services.IServiceRegistry) IFieldController {
	return &FieldController{service: service}
}

func (f *FieldController) GetAllWithPagination(c *gin.Context) {
	// 🚀 Step 1: Inisialisasi dan binding query parameter dari URL
	var params dto.FieldRequestParam
	err := c.ShouldBindQuery(&params)

	// 🧪 Debug input dari query
	fmt.Printf("📥 [DEBUG-FIELD-CONTROLLER] Query Params: %+v\n", params)

	// 🛑 Step 2: Cek jika binding query gagal
	if err != nil {
		fmt.Printf("❌ [ERROR-FIELD-CONTROLLER] Gagal binding query params: %v\n", err)
		response.HttpResponse(response.ParamHttpResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	// ✅ Step 3: Validasi isi struct params (misalnya: required, min, dll)
	validate := validator.New()
	err = validate.Struct(params)

	// 🛑 Step 4: Cek jika validasi gagal
	if err != nil {
		fmt.Printf("❌ [ERROR-FIELD-CONTROLLER] Validasi gagal: %v\n", err)

		// Ambil pesan error dari validasi
		// dan buat response error
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errorResponse := errValidation.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHttpResp{
			Code:    http.StatusBadRequest,
			Err:     err,
			Message: &errMessage,
			Data:    errorResponse,
			Gin:     c,
		})
		return
	}

	// 🔄 Step 5: Panggil service untuk ambil data paginasi field
	result, err := f.service.GetField().GetAllWithPagination(c, &params)

	// 🛑 Step 6: Cek jika ada error dari service/ tangani error jika service gagal
	if err != nil {
		fmt.Printf("❌ [ERROR-FIELD-CONTROLLER] Gagal ambil data field: %v\n", err)
		response.HttpResponse(response.ParamHttpResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	// ✅ Step 7: Jika tidak ada error, kirimkan response sukses
	fmt.Printf("✅ [INFO-FIELD-CONTROLLER] Berhasil ambil data field: %+v\n", result)
	response.HttpResponse(response.ParamHttpResp{
		Code: http.StatusOK,
		Data: result,
		Gin:  c,
	})
}

func (f *FieldController) GetAllWithoutPagination(c *gin.Context) {
	// 🚀 Step 1: Panggil service untuk ambil semua data field
	result, err := f.service.GetField().GetAllWithoutPagination(c)

	if err != nil {
		// ❌ Step 2: Kalau error saat ambil data, tampilkan pesan error + kirim response error ke client
		fmt.Printf("❌ [ERROR-FIELD-CONTROLLER] Gagal ambil data field: %v\n", err)
		response.HttpResponse(response.ParamHttpResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	// ✅ Step 3: Kalau sukses, kirim data ke client dengan status 200 (OK)
	response.HttpResponse(response.ParamHttpResp{
		Code: http.StatusOK,
		Data: result, // Kirim data hasil dari service
		Gin:  c,
	})
}

func (f *FieldController) GetByUUID(c *gin.Context) {
	// 🚀 Step 1: Ambil parameter UUID dari URL
	// Misalnya endpoint: /fields/1234 → maka "1234" adalah UUID yang diambil lewat c.Param("uuid")
	uuid := c.Param("uuid")

	// 📞 Step 2: Panggil service untuk ambil data field berdasarkan UUID
	result, err := f.service.GetField().GetByUUID(c, uuid)

	if err != nil {
		// ❌ Step 3: Kalau gagal ambil data, log error & kirim response error ke client
		fmt.Printf("❌ [ERROR-FIELD-CONTROLLER] Gagal ambil data field (UUID: %s): %v\n", uuid, err)
		response.HttpResponse(response.ParamHttpResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	// ✅ Step 4: Kalau sukses, kirim data ke client dengan status 200 (OK)
	response.HttpResponse(response.ParamHttpResp{
		Code: http.StatusOK,
		Data: result, // Kirim data field
		Gin:  c,
	})
}

func (f *FieldController) Create(c *gin.Context) {
	// 🧾 Step 1: Inisialisasi struct request untuk menampung input dari client
	var request dto.FieldRequest

	// 📥 Step 2: Binding data dari form multipart ke struct request
	err := c.ShouldBindWith(&request, binding.FormMultipart)
	if err != nil {
		// ❌ Step 3: Kalau gagal binding (input tidak cocok), tampilkan error dan kirim response ke client
		fmt.Printf("❌ [ERROR-FIELD-CONTROLLER] Gagal binding request: %v\n", err)
		response.HttpResponse(response.ParamHttpResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})

		return
	}

	// ✅ Debug log untuk input yang diterima
	// ✅ Debug setiap field
	fmt.Println("🔍 [CONTROLLER-DEBUG-FUNC-CREATE] Code:", request.Code)
	fmt.Println("🔍 [CONTROLLER-DEBUG-FUNC-CREATE] Name:", request.Name)
	fmt.Println("🔍 [CONTROLLER-DEBUG-FUNC-CREATE] PricePerHour:", request.PricePerHour)
	// fmt.Println("🔍 [CONTROLLER-DEBUG-FUNC-CREATE] Images:", len(request.Images))

	// ✅ Step 4: Validasi input menggunakan validator (cth: wajib diisi, panjang maksimal, dsb)

	validate := validator.New()
	err = validate.Struct(request)

	if err != nil {
		// ❌ Step 5: Kalau validasi gagal, tampilkan pesan error + detail field yang error
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errorResponse := errValidation.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHttpResp{
			Code:    http.StatusBadRequest,
			Err:     err,
			Message: &errMessage,
			Data:    errorResponse, //Detail keasalahan input
			Gin:     c,
		})
		return
	}

	// 🚀 Step 6: Kirim data request ke layer service untuk dibuat di database
	result, err := f.service.GetField().Create(c, &request)
	if err != nil {
		// ❌ Step 7: Kalau gagal saat proses simpan di service, tampilkan error
		fmt.Printf("❌ [ERROR-FIELD-CONTROLLER] Gagal buat data field: %v\n", err)
		response.HttpResponse(response.ParamHttpResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	// ✅ Step 8: Kalau sukses, kirim response sukses ke client
	response.HttpResponse(response.ParamHttpResp{
		Code: http.StatusOK,
		Data: result,
		Gin:  c,
	})
}

func (f *FieldController) Update(c *gin.Context) {
	// 🧾 Step 1: Siapkan struct untuk menampung data input dari client
	var request dto.UpdateFieldRequest

	// 📥 Step 2: Ambil data dari form (termasuk file) dan simpan ke struct request
	err := c.ShouldBindWith(&request, binding.FormMultipart)
	if err != nil {
		// ❌ Step 3: Kalau data dari client tidak valid (gagal dibaca), tampilkan error
		fmt.Printf("❌ [ERROR-FIELD-CONTROLLER] Gagal binding request: %v\n", err)
		response.HttpResponse(response.ParamHttpResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	// ✅ Debug input dari client
	fmt.Printf("📥 [DEBUG-FIELD-CONTROLLER] Input update: %+v\n", request)

	// 🧪 Step 4: Validasi input (misal field wajib diisi, format harus benar, dll)
	validate := validator.New()
	err = validate.Struct(request)

	if err != nil {
		// ❌ Step 5: Kalau validasi gagal, tampilkan detail kesalahan input
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errorResponse := errValidation.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHttpResp{
			Code:    http.StatusBadRequest,
			Err:     err,
			Message: &errMessage,
			Data:    errorResponse, //Detail keasalahan input
			Gin:     c,
		})
		return
	}

	// 🚀 Step 6: Kirim ke service untuk proses update
	// Ambil UUID dari parameter URL: /field/:uuid
	result, err := f.service.GetField().Update(c, c.Param("uuid"), &request)
	if err != nil {
		// ❌ Step 7: Kalau error saat update di service, kirim response gagal
		fmt.Printf("❌ [ERROR-FIELD-CONTROLLER] Gagal update data field: %v\n", err)
		response.HttpResponse(response.ParamHttpResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	// ✅ Step 8: Kalau sukses update, kirim data hasilnya ke client
	response.HttpResponse(response.ParamHttpResp{
		Code: http.StatusOK,
		Data: result,
		Gin:  c,
	})
}

func (f *FieldController) Delete(c *gin.Context) {
	// 🚀 Step 1: Ambil parameter UUID dari URL
	uuid := c.Param("uuid")

	// 📞 Step 2: Panggil service untuk menghapus data field berdasarkan UUID
	err := f.service.GetField().Delete(c, uuid)

	if err != nil {
		// ❌ Step 3: Kalau gagal hapus data, tampilkan pesan error + kirim response error ke client
		fmt.Printf("❌ [ERROR-FIELD-CONTROLLER] Gagal hapus data field (UUID: %s): %v\n", uuid, err)
		response.HttpResponse(response.ParamHttpResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	// ✅ Step 4: Kalau sukses, kirim response sukses ke client dengan status 200 (OK)
	response.HttpResponse(response.ParamHttpResp{
		Code: http.StatusOK,
		Data: "Data field berhasil dihapus",
		Gin:  c,
	})
	fmt.Printf("✅ [INFO-FIELD-CONTROLLER] Berhasil hapus data field (UUID: %s)\n", uuid)
}
