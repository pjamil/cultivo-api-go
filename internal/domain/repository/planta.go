package repository

import "gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"

type PlantaRepository interface {
	Create(planta *models.Planta) error
	FindAll() ([]models.Planta, error)
	FindByID(id uint) (*models.Planta, error)
	Update(planta *models.Planta) error
	Delete(id uint) error
	FindBySpecies(species models.Especie) ([]models.Planta, error) // Novo método
	FindByStatus(status string) ([]models.Planta, error)           // Novo método
	ExistsByName(name string) bool                                 // Verifica se uma planta com o mesmo nome já existe
}
