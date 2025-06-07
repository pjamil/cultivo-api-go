package repository

import "gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"

type PlantRepository interface {
	Create(plant *models.Plant) error
	FindAll() ([]models.Plant, error)
	FindByID(id uint) (*models.Plant, error)
	Update(plant *models.Plant) error
	Delete(id uint) error
	FindBySpecies(species models.PlantSpecies) ([]models.Plant, error) // Novo método
	FindByStatus(status string) ([]models.Plant, error)                // Novo método
}
