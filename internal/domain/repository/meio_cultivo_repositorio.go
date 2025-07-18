package repository

import (
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/entity"
)

type MeioCultivoRepositorio interface {
	Criar(meioCultivo *entity.MeioCultivo) error
	ListarTodos(page, limit int) ([]entity.MeioCultivo, int64, error)
	BuscarPorID(id uint) (*entity.MeioCultivo, error)
	Atualizar(meioCultivo *entity.MeioCultivo) error
	Deletar(id uint) error
}