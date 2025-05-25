package routes

import (
	"field-service/clients"
	"field-service/constants"
	"field-service/controllers"
	"field-service/middlewares"

	"github.com/gin-gonic/gin"
)

type FieldScheduleRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
	client     clients.IClientRegistry
}

type IFieldScheduleRoute interface {
	Run()
}

func NewFieldScheduleRoute(controller controllers.IControllerRegistry,
	group *gin.RouterGroup, client clients.IClientRegistry) IFieldScheduleRoute {
	return &FieldScheduleRoute{
		controller: controller,
		group:      group,
		client:     client,
	}
}

func (f *FieldScheduleRoute) Run() {
	// 🛣️ Subgroup dengan prefix /field/schedule (sehingga endpoint jadi /field/schedule/...)
	group := f.group.Group("/field/schedule")

	// Middleware optional token (boleh tidak login), AuthenticateWithoutToken

	// 🛣️ [GET] Endpoint untuk mendapatkan semua field schedule berdasarkan ID dan tanggal
	group.GET("/lists/:uuid", middlewares.AuthenticateWithoutToken(), f.controller.GetFieldSchedule().GetAllByFieldIDAndDate)
	// 🛣️ [GET] Endpoint untuk update status fieldSchedule
	group.PATCH("/status", middlewares.AuthenticateWithoutToken(), f.controller.GetFieldSchedule().UpdateStatus)

	// 🔐 Middleware wajib login untuk semua route di bawah ini
	group.Use(middlewares.Authenticate())

	// 🛣️ [GET] Endpoint untuk mendapatkan semua field schedule dengan pagination
	// Middleware Authenticate() untuk memeriksa token dan role
	// Middleware CheckRole untuk memeriksa role user
	// Hanya role Admin dan Customer yang bisa mengakses endpoint ini
	group.GET("/pagination", middlewares.CheckRole([]string{
		constants.Admin,
		constants.Customer,
	}, f.client),
		f.controller.GetFieldSchedule().GetAllWithPagination)

	// 🛣️ [GET] Endpoint untuk mendapatkan field schedule berdasarkan UUID
	group.GET("/:uuid", middlewares.CheckRole([]string{
		constants.Admin,
		constants.Customer,
	}, f.client),
		f.controller.GetFieldSchedule().GetByUUID)

	// ➕ [POST] Endpoint untuk membuat field schedule baru
	// Middleware CheckRole untuk memeriksa role user
	// Hanya role Admin yang bisa mengakses endpoint ini
	group.POST("", middlewares.CheckRole([]string{
		constants.Admin,
	}, f.client),
		f.controller.GetFieldSchedule().Create)

	// 🛣️ [GET] Endpoint untuk generete schedule for one month
	group.POST("/one-month", middlewares.CheckRole([]string{
		constants.Admin,
	}, f.client),
		f.controller.GetFieldSchedule().GenerateScheduleForOneMonth)

	// 🛣️ [PUT] Endpoint untuk upate data berdasarkan uuid
	// Middleware CheckRole untuk memeriksa role user
	// Hanya role Admin yang bisa mengakses endpoint ini
	group.PUT("/:uuid", middlewares.CheckRole([]string{
		constants.Admin,
	}, f.client),
		f.controller.GetFieldSchedule().Update)

	// 🛣️ [DELETE] Endpoint untuk menghapus field schedule berdasarkan UUID (Hanya admin)
	group.DELETE("/:uuid", middlewares.CheckRole([]string{
		constants.Admin,
	}, f.client),
		f.controller.GetFieldSchedule().Delete)
}
