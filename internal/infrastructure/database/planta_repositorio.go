package database

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
)

// PlantaRepositorio implementa a interface repository.PlantaRepositorio
type PlantaRepositorio struct {
	db *gorm.DB
}

// NewPlantaRepositorio cria uma nova instância do PlantaRepositorio
func NewPlantaRepositorio(db *gorm.DB) *PlantaRepositorio {
	return &PlantaRepositorio{db: db}
}

// Criar insere uma nova planta no banco de dados
func (r *PlantaRepositorio) Criar(planta *models.Planta) error {
	if planta == nil {
		return errors.New("planta não pode ser nula")
	}

	result := r.db.Create(planta)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// ListarTodos retorna todas as plantas cadastradas com paginação
func (r *PlantaRepositorio) ListarTodos(page, limit int) ([]models.Planta, int64, error) {
	var plantas []models.Planta
	var total int64

	offset := (page - 1) * limit

	err := r.db.Model(&models.Planta{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if err := r.db.
		Preload("Ambiente").
		Preload("Genetica").
		Preload("MeioCultivo").
		Preload("Usuario").
		Offset(offset).Limit(limit).Find(&plantas).Error; err != nil {
		return nil, 0, fmt.Errorf("falha ao buscar plantas: %w", err)
	}

	return plantas, total, nil
}

// BuscarPorID busca uma planta pelo ID
func (r *PlantaRepositorio) BuscarPorID(id uint) (*models.Planta, error) {
	var planta models.Planta
	if err := r.db.First(&planta, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("erro no banco de dados: %v", err)
	}
	return &planta, nil
}

// Atualizar atualiza uma planta existente
func (r *PlantaRepositorio) Atualizar(planta *models.Planta) error {
	if planta == nil {
		return errors.New("planta não pode ser nula")
	}

	// Verifica se a planta existe
	_, err := r.BuscarPorID(planta.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}
		return err
	}

	result := r.db.Save(planta)
	return result.Error
}

// Deletar remove uma planta pelo ID
func (r *PlantaRepositorio) Deletar(id uint) error {
	result := r.db.Delete(&models.Planta{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// BuscarPorEspecie retorna plantas por espécie (método adicional)
func (r *PlantaRepositorio) BuscarPorEspecie(especie models.Especie) ([]models.Planta, error) {
	var plantas []models.Planta

	result := r.db.Where("species = ?", especie).Find(&plantas)
	if result.Error != nil {
		return nil, result.Error
	}

	return plantas, nil
}

// BuscarPorStatus retorna plantas por status (método adicional)
func (r *PlantaRepositorio) BuscarPorStatus(status string) ([]models.Planta, error) {
	var plantas []models.Planta

	result := r.db.Where("status = ?", status).Find(&plantas)
	if result.Error != nil {
		return nil, result.Error
	}

	return plantas, nil
}

// ExistePorNome verifica se uma planta com o mesmo nome já existe
func (r *PlantaRepositorio) ExistePorNome(nome string) bool {
	var count int64
	result := r.db.Model(&models.Planta{}).Where("nome = ?", nome).Count(&count)
	if result.Error != nil {
		return false // Erro na consulta, assume que não existe
	}
	return count > 0 // Retorna true se existir pelo menos uma planta com o nome
}

func (r *PlantaRepositorio) CriarRegistroDiario(registro *models.RegistroDiario) error {
	if registro == nil {
		return errors.New("registro diário não pode ser nulo")
	}
	result := r.db.Create(registro)
	return result.Error
}
