package routes

import (
	"field-service/clients"
	"field-service/controllers"
	routesField "field-service/routes/field"
	routesFieldSchedule "field-service/routes/fieldschedule"
	routesTime "field-service/routes/time"

	"github.com/gin-gonic/gin"
)

type Registry struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
	client     clients.IClientRegistry
}

type IRegistry interface {
	Serve()
}

func NewRouteRegistry(
	controller controllers.IControllerRegistry,
	group *gin.RouterGroup,
	client clients.IClientRegistry) IRegistry {
	return &Registry{
		controller: controller,
		group:      group,
		client:     client,
	}
}

func (r *Registry) fieldRoute() routesField.IFieldRoute {
	return routesField.NewFieldRoute(r.controller, r.group, r.client)
}

func (r *Registry) fieldScheduleRoute() routesFieldSchedule.IFieldScheduleRoute {
	return routesFieldSchedule.NewFieldScheduleRoute(r.controller, r.group, r.client)
}

func (r *Registry) timeRoute() routesTime.ITimeRoute {
	return routesTime.NewTimeRoute(r.controller, r.group, r.client)
}

func (r *Registry) Serve() {
	// 🛣️ Endpoint untuk field
	r.fieldRoute().Run()

	// 🛣️ Endpoint untuk field schedule
	r.fieldScheduleRoute().Run()

	// 🛣️ Endpoint untuk time
	r.timeRoute().Run()
}
