package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/entity"
)

// RegistroDiarioRepositorio implementa a interface repository.RegistroDiarioRepository
type RegistroDiarioRepositorio struct {
	db *gorm.DB
}

// NewRegistroDiarioRepositorio cria uma nova instância do RegistroDiarioRepositorio
func NewRegistroDiarioRepositorio(db *gorm.DB) *RegistroDiarioRepositorio {
	return &RegistroDiarioRepositorio{db: db}
}

// Criar insere um novo registro diário no banco de dados
func (r *RegistroDiarioRepositorio) Create(registro *entity.RegistroDiario) error {
	if registro == nil {
		return errors.New("registro diário não pode ser nulo")
	}

	result := r.db.Create(registro)
	if result.Error != nil {
		return fmt.Errorf("falha ao criar registro diário: %w", result.Error)
	}

	return nil
}

// ListarTodos retorna todos os registros diários com paginação
func (r *RegistroDiarioRepositorio) GetAll(page, limit int) ([]entity.RegistroDiario, int64, error) {
	var registros []entity.RegistroDiario
	var total int64

	offset := (page - 1) * limit

	err := r.db.Model(&entity.RegistroDiario{}).Count(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("falha ao contar registros diários: %w", err)
	}

	if err := r.db.Offset(offset).Limit(limit).Find(&registros).Error; err != nil {
		return nil, 0, fmt.Errorf("falha ao buscar registros diários: %w", err)
	}

	return registros, total, nil
}

// BuscarPorID busca um registro diário pelo ID
func (r *RegistroDiarioRepositorio) GetByID(id uint) (*entity.RegistroDiario, error) {
	var registro entity.RegistroDiario
	if err := r.db.First(&registro, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("erro no banco de dados ao buscar registro diário: %w", err)
	}
	return &registro, nil
}

// Atualizar atualiza um registro diário existente
func (r *RegistroDiarioRepositorio) Update(registro *entity.RegistroDiario) error {
	if registro == nil {
		return errors.New("registro diário não pode ser nulo")
	}

	result := r.db.Save(registro)
	if result.Error != nil {
		return fmt.Errorf("falha ao atualizar registro diário: %w", result.Error)
	}
	return nil
}

// Deletar remove um registro diário pelo ID
func (r *RegistroDiarioRepositorio) Delete(id uint) error {
	result := r.db.Delete(&entity.RegistroDiario{}, id)
	if result.Error != nil {
		return fmt.Errorf("falha ao deletar registro diário: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// ListarPorDiarioCultivoID retorna registros diários filtrados por DiarioCultivoID com paginação
func (r *RegistroDiarioRepositorio) ListarPorDiarioCultivoID(diarioCultivoID uint, page, limit int) ([]entity.RegistroDiario, int64, error) {
	var registros []entity.RegistroDiario
	var total int64

	offset := (page - 1) * limit

	err := r.db.Model(&entity.RegistroDiario{}).Where("diario_cultivo_id = ?", diarioCultivoID).Count(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("falha ao contar registros diários por diario_cultivo_id: %w", err)
	}

	if err := r.db.Where("diario_cultivo_id = ?", diarioCultivoID).Offset(offset).Limit(limit).Order("data desc").Find(&registros).Error; err != nil {
		return nil, 0, fmt.Errorf("falha ao buscar registros diários por diario_cultivo_id: %w", err)
	}

	return registros, total, nil
}
