package service_test

import (
	"errors"
	"testing"
	"time"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/entity"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/service"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/service/test"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestRegistroDiarioService_CriarRegistro(t *testing.T) {
	mockRegistroRepo := new(test.MockRegistroDiarioRepositorio)
	mockDiarioRepo := new(test.MockDiarioCultivoRepository)
	s := service.NewRegistroDiarioService(mockRegistroRepo, mockDiarioRepo)

	diarioID := uint(1)
	createDTO := &dto.CreateRegistroDiarioDTO{
		Titulo:   "Teste Registro",
		Conteudo: "Conteúdo do registro de teste",
		Data:     time.Now(),
		Tipo:     "observacao",
	}

	t.Run("Success", func(t *testing.T) {
		// Arrange
		mockDiarioRepo.On("GetByID", diarioID).Return(&entity.DiarioCultivo{}, nil).Once()
		mockRegistroRepo.On("Create", mock.AnythingOfType("*entity.RegistroDiario")).Return(nil).Once().Run(func(args mock.Arguments) {
			registro := args.Get(0).(*entity.RegistroDiario)
			registro.ID = 1 // Simulate ID assignment
		})

		// Act
		response, err := s.CriarRegistro(diarioID, createDTO)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, uint(1), response.ID)
		assert.Equal(t, createDTO.Titulo, response.Titulo)
		mockDiarioRepo.AssertExpectations(t)
		mockRegistroRepo.AssertExpectations(t)
	})

	t.Run("Error - DiarioCultivo Not Found", func(t *testing.T) {
		// Arrange
		mockDiarioRepo.On("GetByID", diarioID).Return(nil, gorm.ErrRecordNotFound).Once()

		// Act
		response, err := s.CriarRegistro(diarioID, createDTO)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, utils.ErrNotFound, err)
		mockDiarioRepo.AssertExpectations(t)
		mockRegistroRepo.AssertNotCalled(t, "Create")
	})

	t.Run("Error - Repository Create Error", func(t *testing.T) {
		// Arrange
		mockDiarioRepo.On("GetByID", diarioID).Return(&entity.DiarioCultivo{}, nil).Once()
		mockRegistroRepo.On("Create", mock.AnythingOfType("*entity.RegistroDiario")).Return(errors.New("db error")).Once()

		// Act
		response, err := s.CriarRegistro(diarioID, createDTO)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "falha ao criar registro no diário: db error")
		mockDiarioRepo.AssertExpectations(t)
		mockRegistroRepo.AssertExpectations(t)
	})
}

func TestRegistroDiarioService_ListarRegistrosPorDiarioID(t *testing.T) {
	mockRegistroRepo := new(test.MockRegistroDiarioRepositorio)
	mockDiarioRepo := new(test.MockDiarioCultivoRepository) // Not directly used in this test, but needed for service creation
	s := service.NewRegistroDiarioService(mockRegistroRepo, mockDiarioRepo)

	diarioID := uint(1)
	page := 1
	limit := 10

	t.Run("Success - Registros Found", func(t *testing.T) {
		// Arrange
		mockRegistros := []entity.RegistroDiario{
			{Model: gorm.Model{ID: 1}, Titulo: "Reg 1", Conteudo: "Content 1", Data: time.Now(), Tipo: "observacao"},
			{Model: gorm.Model{ID: 2}, Titulo: "Reg 2", Conteudo: "Content 2", Data: time.Now(), Tipo: "evento"},
		}
		expectedTotal := int64(len(mockRegistros))
		expectedResponseDTOs := []dto.RegistroDiarioResponseDTO{
			{ID: 1, Titulo: "Reg 1", Conteudo: "Content 1", Data: mockRegistros[0].Data, Tipo: "observacao"},
			{ID: 2, Titulo: "Reg 2", Conteudo: "Content 2", Data: mockRegistros[1].Data, Tipo: "evento"},
		}

		mockRegistroRepo.On("ListarPorDiarioCultivoID", diarioID, page, limit).Return(mockRegistros, expectedTotal, nil).Once()

		// Act
		responseDTOs, total, err := s.ListarRegistrosPorDiarioID(diarioID, page, limit)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, responseDTOs)
		assert.Equal(t, expectedTotal, total)
		assert.Equal(t, expectedResponseDTOs, responseDTOs)
		mockRegistroRepo.AssertExpectations(t)
	})

	t.Run("Success - No Registros Found", func(t *testing.T) {
		// Arrange
		mockRegistroRepo.On("ListarPorDiarioCultivoID", diarioID, page, limit).Return([]entity.RegistroDiario{}, int64(0), nil).Once()

		// Act
		responseDTOs, total, err := s.ListarRegistrosPorDiarioID(diarioID, page, limit)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, responseDTOs)
		assert.Equal(t, int64(0), total)
		assert.Empty(t, responseDTOs)
		mockRegistroRepo.AssertExpectations(t)
	})

	t.Run("Error - Repository List Error", func(t *testing.T) {
		// Arrange
		mockRegistroRepo.On("ListarPorDiarioCultivoID", diarioID, page, limit).Return([]entity.RegistroDiario{}, int64(0), errors.New("db error")).Once()

		// Act
		responseDTOs, total, err := s.ListarRegistrosPorDiarioID(diarioID, page, limit)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, responseDTOs)
		assert.Equal(t, int64(0), total)
		assert.Contains(t, err.Error(), "falha ao listar registros do diário: db error")
		mockRegistroRepo.AssertExpectations(t)
	})
}
