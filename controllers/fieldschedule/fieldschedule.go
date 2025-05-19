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

type FieldScheduleController struct {
	service services.IServiceRegistry
}

type IFieldScheduleController interface {
	GetAllWithPagination(*gin.Context)
	GetAllByFieldIDAndDate(*gin.Context)
	GetByUUID(*gin.Context)
	Create(*gin.Context)
	Update(*gin.Context)
	UpdateStatus(*gin.Context)
	Delete(*gin.Context)
	GenerateScheduleForOneMonth(*gin.Context)
}

func NewFieldScheduleController(service services.IServiceRegistry) IFieldScheduleController {
	return &FieldScheduleController{service: service}
}

func (f *FieldScheduleController) GetAllWithPagination(c *gin.Context) {
	// 🚀 Step 1: Inisialisasi dan binding query parameter dari URL
	var params dto.FieldScheduleRequestParam
	err := c.ShouldBindQuery(&params)

	// 🧪 Debug input dari query
	fmt.Printf("📥 [DEBUG-FIELDSCHEDULE-CONTROLLER] Query Params: %+v\n", params)

	// 🛑 Step 2: Cek jika binding query gagal
	if err != nil {
		fmt.Printf("❌ [ERROR-FIELDSCHEDULE-CONTROLLER] Gagal binding query params: %v\n", err)
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
		fmt.Printf("❌ [ERROR-FIELDSCHEDULE-CONTROLLER] Validasi gagal: %v\n", err)

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
	result, err := f.service.GetFieldSchedule().GetAllWithPagination(c, &params)

	// 🛑 Step 6: Cek jika ada error dari service/ tangani error jika service gagal
	if err != nil {
		fmt.Printf("❌ [ERROR-FIELDSCHEDULE-CONTROLLER] Gagal ambil data field: %v\n", err)
		response.HttpResponse(response.ParamHttpResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	// ✅ Step 7: Jika tidak ada error, kirimkan response sukses
	fmt.Printf("✅ [INFO-FIELDSCHEDULE-CONTROLLER] Berhasil ambil data field: %+v\n", result)
	response.HttpResponse(response.ParamHttpResp{
		Code: http.StatusOK,
		Data: result,
		Gin:  c,
	})
}

func (f *FieldScheduleController) GetAllByFieldIDAndDate(c *gin.Context) {
	// 📦 Step 1: Siapkan struct untuk menampung query parameter dari URL (?date=...)
	var params dto.FieldScheduleByFieldIDAndDateRequestParam

	// 🧲 Ambil query dari URL dan simpan ke struct params
	err := c.ShouldBindQuery(&params)

	// 🧪 Debug input dari query
	fmt.Printf("📥 [DEBUG-FIELDSCHEDULE-CONTROLLER] Query Params: %+v\n", params)

	// 🛑 Step 2: Cek jika binding query gagal
	if err != nil {
		fmt.Printf("❌ [ERROR-FIELDSCHEDULE-CONTROLLER] Gagal binding query params: %v\n", err)
		response.HttpResponse(response.ParamHttpResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	// ✅ Step 3: Validasi isi struct params (misalnya: date wajib diisi, dll)
	validate := validator.New()
	err = validate.Struct(params)

	// 🛑 Step 4: Cek jika validasi gagal
	if err != nil {
		fmt.Printf("❌ [ERROR-FIELDSCHEDULE-CONTROLLER] Validasi gagal: %v\n", err)

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

	// 🚀 Step 5: Panggil service dengan UUID dari path dan tanggal dari query
	result, err := f.service.GetFieldSchedule().GetAllByFieldIDAndDate(c, c.Param("uuid"), params.Date)

	// 🛑 Step 6: Cek jika ada error saat ambil data dari service
	if err != nil {
		response.HttpResponse(response.ParamHttpResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	// ✅ Step 7: Jika tidak ada error, kirimkan response sukses
	fmt.Printf("✅ [INFO-FIELDSCHEDULE-CONTROLLER] Berhasil ambil data field: %+v\n", result)
	response.HttpResponse(response.ParamHttpResp{
		Code: http.StatusOK,
		Data: result,
		Gin:  c,
	})
}

func (f *FieldScheduleController) GetByUUID(c *gin.Context) {
	// 🚀 Step 1: Ambil UUID dari URL
	uuid := c.Param("uuid")
	fmt.Printf("🔍 [DEBUG-FIELDSCHEDULE-CONTROLLER] UUID: %s\n", uuid)

	// 🔄 Step 2: Panggil service untuk ambil data berdasarkan UUID
	result, err := f.service.GetFieldSchedule().GetByUUID(c, uuid)

	// 🛑 Step 3: Cek jika ada error saat ambil data dari service
	if err != nil {
		fmt.Printf("❌ [ERROR-FIELDSCHEDULE-CONTROLLER] Gagal ambil data field schedule: %v\n", err)
		response.HttpResponse(response.ParamHttpResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	// ✅ Step 4: Jika tidak ada error, kirimkan response sukses
	fmt.Printf("✅ [INFO-FIELDSCHEDULE-CONTROLLER] Berhasil ambil data field schedule: %+v\n", result)
	response.HttpResponse(response.ParamHttpResp{
		Code: http.StatusOK,
		Data: result,
		Gin:  c,
	})
}

func (f *FieldScheduleController) Create(c *gin.Context) {
	// 🧾 Step 1: Siapkan struct untuk menampung request dari client (body JSON)
	var params dto.FieldScheduleRequest

	// 🧲 Step 2: Ambil data dari body JSON dan simpan ke struct params
	err := c.ShouldBindJSON(&params)
	if err != nil {
		// ❌ Jika gagal binding (misalnya format JSON salah), kirim error ke client
		fmt.Printf("❌ [ERROR-FIELDSCHEDULE-CONTROLLER] Gagal binding JSON: %v\n", err)
		response.HttpResponse(response.ParamHttpResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	// ✅ Step 3: Validasi input (pastikan semua field sesuai aturan: required, format, dll)
	validate := validator.New()
	err = validate.Struct(params)
	if err != nil {
		// ❌ Jika validasi gagal, kirim pesan error dan detail kesalahannya
		fmt.Printf("❌ [ERROR-FIELDSCHEDULE-CONTROLLER] Validasi gagal: %v\n", err)
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

	// 🚀 Step 4: Panggil service untuk simpan data field schedule ke database
	err = f.service.GetFieldSchedule().Create(c, &params)
	if err != nil {
		// ❌ Jika gagal simpan (misal karena konflik jadwal atau DB error), kirim error
		fmt.Printf("❌ [ERROR-FIELDSCHEDULE-CONTROLLER] Gagal membuat field schedule: %v\n", err)
		response.HttpResponse(response.ParamHttpResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	// ✅ Step 5: Jika berhasil, kirim response sukses dengan status 201 (Created)
	response.HttpResponse(response.ParamHttpResp{
		Code: http.StatusCreated,
		Gin:  c,
	})
}

func (f *FieldScheduleController) GenerateScheduleForOneMonth(c *gin.Context) {
	// 🧾 Step 1: Siapkan struct untuk menampung request dari client (body JSON)
	var params dto.GenerateFieldScheduleForOneMonthRequest

	// 🧲 Step 2: Ambil data dari body JSON dan simpan ke struct params
	err := c.ShouldBindJSON(&params)
	if err != nil {
		// ❌ Jika gagal binding (misalnya format JSON salah), kirim error ke client
		fmt.Printf("❌ [ERROR-FIELDSCHEDULE-CONTROLLER] Gagal binding JSON: %v\n", err)
		response.HttpResponse(response.ParamHttpResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	// ✅ Step 3: Validasi input (pastikan semua field sesuai aturan: required, format, dll)
	// Validasi input menggunakan validator
	// Pastikan semua field sesuai aturan: required, format, dll
	validate := validator.New()
	err = validate.Struct(params)
	if err != nil {
		// ❌ Jika validasi gagal, kirim pesan error dan detail kesalahannya
		// Ambil pesan error dari validasi
		fmt.Printf("❌ [ERROR-FIELDSCHEDULE-CONTROLLER] Validasi gagal: %v\n", err)
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

	// 🚀 Step 4: Panggil service untuk proses generate schedule sebulan ke database
	err = f.service.GetFieldSchedule().GenerateScheduleForOneMonth(c, &params)
	if err != nil {
		// ❌ Jika gagal simpan, kirim error
		fmt.Printf("❌ [ERROR-FIELDSCHEDULE-CONTROLLER] Gagal membuat field schedule: %v\n", err)
		response.HttpResponse(response.ParamHttpResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	// ✅ Step 5: Jika berhasil, kirim response sukses dengan status 201 (Created)
	response.HttpResponse(response.ParamHttpResp{
		Code: http.StatusCreated,
		Gin:  c,
	})
}

func (f *FieldScheduleController) Update(c *gin.Context) {
	// 🧾 Step 1: Siapkan struct untuk menampung request dari client (body JSON)
	var params dto.UpdateFieldScheduleRequest

	// 🧲 Step 2: Ambil data dari body JSON dan simpan ke struct params
	err := c.ShouldBindJSON(&params)
	if err != nil {
		// ❌ Jika gagal binding (misalnya format JSON salah), kirim error ke client
		fmt.Printf("❌ [ERROR-FIELDSCHEDULE-CONTROLLER] Gagal binding JSON: %v\n", err)
		response.HttpResponse(response.ParamHttpResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	// ✅ Step 3: Validasi input (pastikan semua field sesuai aturan: required, format, dll)
	validate := validator.New()
	err = validate.Struct(params)
	if err != nil {
		// ❌ Jika validasi gagal, kirim pesan error dan detail kesalahannya
		fmt.Printf("❌ [ERROR-FIELDSCHEDULE-CONTROLLER] Validasi gagal: %v\n", err)
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

	// 🚀 Step 4: Panggil service untuk update data ke database
	result, err := f.service.GetFieldSchedule().Update(c, c.Param("uuid"), &params)
	if err != nil {
		// ❌ Jika gagal update, kirim error
		fmt.Printf("❌ [ERROR-FIELDSCHEDULE-CONTROLLER] Gagal update data field schedule: %v\n", err)
		response.HttpResponse(response.ParamHttpResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	// ✅ Step 5: Jika berhasil, kirim response sukses dengan status 200 (Ok)
	response.HttpResponse(response.ParamHttpResp{
		Code: http.StatusOK,
		Gin:  c,
		Data: result,
	})
}

func (f *FieldScheduleController) UpdateStatus(c *gin.Context) {
	// 🧾 Step 1: Siapkan struct untuk menampung request dari client (body JSON)
	var request dto.UpdateStatusFieldScheduleRequest

	// 🧲 Step 2: Ambil data dari body JSON dan simpan ke struct params
	err := c.ShouldBindJSON(&request)
	if err != nil {
		// ❌ Jika gagal binding (misalnya format JSON salah), kirim error ke client
		fmt.Printf("❌ [ERROR-FIELDSCHEDULE-CONTROLLER] Gagal binding JSON: %v\n", err)
		response.HttpResponse(response.ParamHttpResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	// ✅ Step 3: Validasi input (pastikan semua field sesuai aturan: required, format, dll)
	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		// ❌ Jika validasi gagal, kirim pesan error dan detail kesalahannya
		fmt.Printf("❌ [ERROR-FIELDSCHEDULE-CONTROLLER] Validasi gagal: %v\n", err)
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

	// 🚀 Step 4: Panggil service untuk update data ke database
	err = f.service.GetFieldSchedule().UpdateStatus(c, &request)
	if err != nil {
		// ❌ Jika gagal update, kirim error
		fmt.Printf("❌ [ERROR-FIELDSCHEDULE-CONTROLLER] Gagal update data field schedule: %v\n", err)
		response.HttpResponse(response.ParamHttpResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	// ✅ Step 5: Jika berhasil, kirim response sukses dengan status 200 (Ok)
	response.HttpResponse(response.ParamHttpResp{
		Code: http.StatusOK,
		Gin:  c,
	})
}

func (f *FieldScheduleController) Delete(c *gin.Context) {
	// 🚀 Step 1: Ambil UUID dari URL
	uuid := c.Param("uuid")
	fmt.Printf("🔍 [DEBUG-FIELDSCHEDULE-CONTROLLER] UUID: %s\n", uuid)

	// 🔄 Step 2: Panggil service untuk hapus data berdasarkan UUID
	err := f.service.GetFieldSchedule().Delete(c, uuid)

	// 🛑 Step 3: Cek jika ada error saat hapus data dari service
	if err != nil {
		fmt.Printf("❌ [ERROR-FIELDSCHEDULE-CONTROLLER] Gagal hapus data field schedule: %v\n", err)
		response.HttpResponse(response.ParamHttpResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	// ✅ Step 4: Jika tidak ada error, kirimkan response sukses
	fmt.Printf("✅ [INFO-FIELDSCHEDULE-CONTROLLER] Berhasil hapus data field schedule\n")
	response.HttpResponse(response.ParamHttpResp{
		Code: http.StatusOK,
		Gin:  c,
	})
}
