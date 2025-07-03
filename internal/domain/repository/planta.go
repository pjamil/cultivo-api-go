package repository

import "gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"

type PlantaRepositorio interface {
	Criar(planta *models.Planta) error
	ListarTodos(page, limit int) ([]models.Planta, int64, error)
	BuscarPorID(id uint) (*models.Planta, error)
	Atualizar(planta *models.Planta) error
	Deletar(id uint) error
	BuscarPorEspecie(especie models.Especie) ([]models.Planta, error)
	BuscarPorStatus(status string) ([]models.Planta, error)
	ExistePorNome(nome string) bool
	CriarRegistroDiario(registro *models.RegistroDiario) error
}
