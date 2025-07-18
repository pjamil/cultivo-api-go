package repository

import (
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/entity"
)

type PlantaRepositorio interface {
	Criar(planta *entity.Planta) error
	ListarTodos(page, limit int) ([]entity.Planta, int64, error)
	BuscarPorID(id uint) (*entity.Planta, error)
	Atualizar(planta *entity.Planta) error
	Deletar(id uint) error
	BuscarPorEspecie(especie entity.Especie) ([]entity.Planta, error)
	BuscarPorStatus(status string) ([]entity.Planta, error)
	ExistePorNome(nome string) bool
	CriarRegistroDiario(registro *entity.RegistroDiario) error
}