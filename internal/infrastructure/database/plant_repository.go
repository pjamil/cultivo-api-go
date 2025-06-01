package database

import (
	"errors"

	"gorm.io/gorm"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
)

// PlantRepository implementa a interface repository.PlantRepository
type PlantRepository struct {
	db *gorm.DB
}

// NewPlantRepository cria uma nova instância do PlantRepository
func NewPlantRepository(db *gorm.DB) *PlantRepository {
	return &PlantRepository{db: db}
}

// Create insere uma nova planta no banco de dados
func (r *PlantRepository) Create(plant *models.Plant) error {
	if plant == nil {
		return errors.New("plant cannot be nil")
	}

	result := r.db.Create(plant)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// FindAll retorna todas as plantas cadastradas
func (r *PlantRepository) FindAll() ([]models.Plant, error) {
	var plants []models.Plant

	result := r.db.Find(&plants)
	if result.Error != nil {
		return nil, result.Error
	}

	return plants, nil
}

// FindByID busca uma planta pelo ID
func (r *PlantRepository) FindByID(id uint) (*models.Plant, error) {
	var plant models.Plant

	result := r.db.First(&plant, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Retorna nil sem erro quando não encontrado
		}
		return nil, result.Error
	}

	return &plant, nil
}

// Update atualiza uma planta existente
func (r *PlantRepository) Update(plant *models.Plant) error {
	if plant == nil {
		return errors.New("plant cannot be nil")
	}

	// Verifica se a planta existe
	existingPlant, err := r.FindByID(plant.ID)
	if err != nil {
		return err
	}
	if existingPlant == nil {
		return errors.New("plant not found")
	}

	result := r.db.Save(plant)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// Delete remove uma planta pelo ID
func (r *PlantRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Plant{}, id)
	if result.Error != nil {
		return result.Error
	}

	// Verifica se algum registro foi realmente deletado
	if result.RowsAffected == 0 {
		return errors.New("no plant found with the given ID")
	}

	return nil
}

// FindBySpecies retorna plantas por espécie (método adicional)
func (r *PlantRepository) FindBySpecies(species models.PlantSpecies) ([]models.Plant, error) {
	var plants []models.Plant

	result := r.db.Where("species = ?", species).Find(&plants)
	if result.Error != nil {
		return nil, result.Error
	}

	return plants, nil
}

// FindByStatus retorna plantas por status (método adicional)
func (r *PlantRepository) FindByStatus(status string) ([]models.Plant, error) {
	var plants []models.Plant

	result := r.db.Where("status = ?", status).Find(&plants)
	if result.Error != nil {
		return nil, result.Error
	}

	return plants, nil
}
