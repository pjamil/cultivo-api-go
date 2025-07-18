package repository

import (
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/entity"
)

type RegistroDiarioRepository interface {
	Criar(registro *entity.RegistroDiario) error
	ListarTodos(page, limit int) ([]entity.RegistroDiario, int64, error)
	BuscarPorID(id uint) (*entity.RegistroDiario, error)
	Atualizar(registro *entity.RegistroDiario) error
	Deletar(id uint) error
	ListarPorDiarioCultivoID(diarioCultivoID uint, page, limit int) ([]entity.RegistroDiario, int64, error)
}