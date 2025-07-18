package repository

import (
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/entity"
)

type GeneticaRepositorio interface {
	Criar(genetica *entity.Genetica) error
	ListarTodos(page, limit int) ([]entity.Genetica, int64, error)
	BuscarPorID(id uint) (*entity.Genetica, error)
	Atualizar(genetica *entity.Genetica) error
	Deletar(id uint) error
}