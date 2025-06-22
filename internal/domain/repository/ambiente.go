package repository

import (
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gorm.io/gorm"
)

type AmbienteRepository interface {
	FindByID(id uint) (*models.Ambiente, error)
	// Add other methods as needed
}

type ambienteRepository struct {
	db *gorm.DB // Uncomment if you need a database connection
}

func NewAmbienteRepository(db *gorm.DB) AmbienteRepository {
	return &ambienteRepository{db: db}
}
func (r *ambienteRepository) FindByID(id uint) (*models.Ambiente, error) {
	var ambiente models.Ambiente
	if err := r.db.First(&ambiente, id).Error; err != nil {
		return nil, err
	}
	return &ambiente, nil
}
