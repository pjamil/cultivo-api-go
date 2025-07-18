package repository

import (
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/entity"
)

type AmbienteRepositorio interface {
	Criar(ambiente *entity.Ambiente) error
	ListarTodos(page, limit int) ([]entity.Ambiente, int64, error)
	BuscarPorID(id uint) (*entity.Ambiente, error)
	Atualizar(ambiente *entity.Ambiente) error
	Deletar(id uint) error
}