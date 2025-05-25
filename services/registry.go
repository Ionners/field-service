package services

import (
	"field-service/common/gcs"
	"field-service/repositories"
	fieldService "field-service/services/field"
	fieldScheduleService "field-service/services/fieldschedule"
	timeService "field-service/services/time"
	"fmt"
)

type Registry struct {
	repository repositories.IRepositoryRegistry
	gcs        gcs.IGCSClient
}

type IServiceRegistry interface {
	GetField() fieldService.IFieldService
	GetFieldSchedule() fieldScheduleService.IFieldScheduleService
	GetTime() timeService.ITimeService
}

func NewServiceRegistry(repository repositories.IRepositoryRegistry, gcs gcs.IGCSClient) IServiceRegistry {
	if repository == nil {
		fmt.Println("❌ [REGISTRY-SERVICE-ERROR] FieldService: repo is nil")
	}
	if repository.GetField() == nil {
		fmt.Println("❌ [REGISTRY-SERVICE-ERROR] FieldService: repo.GetField() is nil")
	}
	return &Registry{
		repository: repository,
		gcs:        gcs,
	}
}

func (r *Registry) GetField() fieldService.IFieldService {
	return fieldService.NewFieldService(r.repository, r.gcs)
}

func (r *Registry) GetFieldSchedule() fieldScheduleService.IFieldScheduleService {
	return fieldScheduleService.NewFieldScheduleService(r.repository)
}

func (r *Registry) GetTime() timeService.ITimeService {
	return timeService.NewTimeService(r.repository)
}
