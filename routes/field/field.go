package routes

import (
	"field-service/clients"
	"field-service/constants"
	"field-service/controllers"
	"field-service/middlewares"

	"github.com/gin-gonic/gin"
)

type FieldRoute struct {
	controller controllers.IControllerRegistry //mengakses controller terkait field
	group      *gin.RouterGroup                //prefix URL untuk route field
	client     clients.IClientRegistry         // Depedency untuk pengecekan role user (auth)
}

type IFieldRoute interface {
	Run()
}

func NewFieldRoute(controller controllers.IControllerRegistry,
	group *gin.RouterGroup, client clients.IClientRegistry) IFieldRoute {
	return &FieldRoute{
		controller: controller,
		group:      group,
		client:     client,
	}
}

func (f *FieldRoute) Run() {
	// 🛣️ Subgroup dengan prefix /field (sehingga endpoint jadi /field/...)
	group := f.group.Group("/field")

	// 🛣️ [GET] Endpoint untuk mendapatkan semua field tanpa pagination
	group.GET("", middlewares.AuthenticateWithoutToken(), // Middleware optional token (boleh tidak login)
		f.controller.GetField().GetAllWithoutPagination)

	// 🛣️ [GET] Endpoint untuk mendapatkan field berdasarkan UUID
	// Middleware optional token (boleh tidak login), AuthenticateWithoutToken
	group.GET("/:uuid", middlewares.AuthenticateWithoutToken(), f.controller.GetField().GetByUUID)

	// 🔐 Middleware wajib login untuk semua route di bawah ini
	group.Use(middlewares.Authenticate())

	// 🛣️ [GET] Endpoint untuk mendapatkan semua field dengan pagination
	// Middleware Authenticate() untuk memeriksa token dan role
	// Hanya role Admin dan Customer yang bisa mengakses endpoint ini
	group.GET("/pagination", middlewares.CheckRole([]string{
		constants.Admin,
		constants.Customer,
	}, f.client),
		f.controller.GetField().GetAllWithPagination)

	// ➕ [POST] Endpoint untuk membuat field baru
	// Middleware CheckRole untuk memeriksa role user
	// Hanya role Admin yang bisa mengakses endpoint ini
	group.POST("/", middlewares.CheckRole([]string{
		constants.Admin,
	}, f.client),
		f.controller.GetField().Create)

	// 🛣️ [PUT] Endpoint untuk update data field berdasarkan UUID
	// Middleware CheckRole untuk memeriksa role user
	// Hanya role Admin yang bisa mengakses endpoint ini
	group.PUT("/:uuid", middlewares.CheckRole([]string{
		constants.Admin,
	}, f.client),
		f.controller.GetField().Update)

	// 🛣️ [DELETE] Endpoint untuk menghapus field berdasarkan UUID
	// Middleware CheckRole untuk memeriksa role user
	// Hanya role Admin yang bisa mengakses endpoint ini
	group.DELETE("/:uuid", middlewares.CheckRole([]string{
		constants.Admin,
	}, f.client),
		f.controller.GetField().Delete)
}
