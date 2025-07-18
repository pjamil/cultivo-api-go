package repository

import (
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/entity"
)

type RegistroDiarioRepositorio interface {
	Create(registroDiario *entity.RegistroDiario) error
	GetByID(id uint) (*entity.RegistroDiario, error)
	GetAll(page, limit int) ([]entity.RegistroDiario, int64, error)
	Atualizar(registroDiario *entity.RegistroDiario) error
	Delete(id uint) error
}