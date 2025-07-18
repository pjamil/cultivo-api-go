package repository

import (
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/entity"
)

type DiarioCultivoRepositorio interface {
	Create(diarioCultivo *entity.DiarioCultivo) error
	GetAll(page, limit int) ([]entity.DiarioCultivo, int64, error)
	GetAllByUserID(userID uint) ([]entity.DiarioCultivo, error)
	GetByID(id uint) (*entity.DiarioCultivo, error)
	Update(diarioCultivo *entity.DiarioCultivo) error
	Delete(id uint) error
}