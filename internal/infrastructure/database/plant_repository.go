package database

import (
	"errors"
	"fmt"

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
func (r *PlantRepository) Create(plant *models.Planta) error {
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
func (r *PlantRepository) FindAll() ([]models.Planta, error) {
	var plants []models.Planta

	if err := r.db.
		Preload("Ambiente").
		Preload("Genetica").
		Preload("MeioCultivo").
		Preload("Usuario").
		Find(&plants).Error; err != nil {
		return nil, fmt.Errorf("falha ao buscar plantas: %w", err)
	}

	if len(plants) == 0 {
		return nil, errors.New("nenhuma planta encontrada")
	}

	return plants, nil
}

// FindByID busca uma planta pelo ID
func (r *PlantRepository) FindByID(id uint) (*models.Planta, error) {
	var plant models.Planta
	if err := r.db.First(&plant, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("database error: %v", err)
	}
	return &plant, nil
}

// Update atualiza uma planta existente
func (r *PlantRepository) Update(plant *models.Planta) error {
	if plant == nil {
		return errors.New("plant cannot be nil")
	}

	// Verifica se a planta existe
	_, err := r.FindByID(plant.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}
		return err
	}

	result := r.db.Save(plant)
	return result.Error
}

// Delete remove uma planta pelo ID
func (r *PlantRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Planta{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// FindBySpecies retorna plantas por espécie (método adicional)
func (r *PlantRepository) FindBySpecies(species models.Especie) ([]models.Planta, error) {
	var plants []models.Planta

	result := r.db.Where("species = ?", species).Find(&plants)
	if result.Error != nil {
		return nil, result.Error
	}

	return plants, nil
}

// FindByStatus retorna plantas por status (método adicional)
func (r *PlantRepository) FindByStatus(status string) ([]models.Planta, error) {
	var plants []models.Planta

	result := r.db.Where("status = ?", status).Find(&plants)
	if result.Error != nil {
		return nil, result.Error
	}

	return plants, nil
}

// ExistsByName verifica se uma planta com o mesmo nome já existe
func (r *PlantRepository) ExistsByName(name string) bool {
	var count int64
	result := r.db.Model(&models.Planta{}).Where("name = ?", name).Count(&count)
	if result.Error != nil {
		return false // Erro na consulta, assume que não existe
	}
	return count > 0 // Retorna true se existir pelo menos uma planta com o nome
}
