package repositories

import (
	fieldRepositories "field-service/repositories/field"
	fieldScheduleRepositories "field-service/repositories/fieldschedule"
	timeRepositories "field-service/repositories/time"

	"gorm.io/gorm"
)

type Registry struct {
	db *gorm.DB
}

type IRepositoryRegistry interface {
	GetField() fieldRepositories.IFieldRepository
	GetFieldSchedule() fieldScheduleRepositories.IFieldScheduleRepository
	GetTime() timeRepositories.ITimeRepository
}

func NewRepositoryRegistry(db *gorm.DB) IRepositoryRegistry {
	return &Registry{db: db}
}

func (r *Registry) GetField() fieldRepositories.IFieldRepository {
	return fieldRepositories.NewFieldRepository(r.db)
}

func (r *Registry) GetFieldSchedule() fieldScheduleRepositories.IFieldScheduleRepository {
	return fieldScheduleRepositories.NewFieldScheduleRepository(r.db)
}

func (r *Registry) GetTime() timeRepositories.ITimeRepository {
	return timeRepositories.NewTimeRepository(r.db)
}
