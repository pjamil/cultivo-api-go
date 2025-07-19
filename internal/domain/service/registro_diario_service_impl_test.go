package service_test

import (
	"errors"
	"testing"
	"time"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/entity"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/service"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockRegistroDiarioRepositorio is a mock for the RegistroDiarioRepository interface.
type MockRegistroDiarioRepositorio struct {
	mock.Mock
}

func (m *MockRegistroDiarioRepositorio) Criar(registro *entity.RegistroDiario) error {
	args := m.Called(registro)
	return args.Error(0)
}

// ListarTodos is not used by the service under test, but required by the interface.
func (m *MockRegistroDiarioRepositorio) ListarTodos(page, limit int) ([]entity.RegistroDiario, int64, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]entity.RegistroDiario), args.Get(1).(int64), args.Error(2)
}
func (m *MockRegistroDiarioRepositorio) ListarPorDiarioCultivoID(diarioID uint, page, limit int) ([]entity.RegistroDiario, int64, error) {
	args := m.Called(diarioID, page, limit)
	return args.Get(0).([]entity.RegistroDiario), args.Get(1).(int64), args.Error(2)
}

func (m *MockRegistroDiarioRepositorio) BuscarPorID(id uint) (*entity.RegistroDiario, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.RegistroDiario), args.Error(1)
}

func (m *MockRegistroDiarioRepositorio) Atualizar(registro *entity.RegistroDiario) error {
	args := m.Called(registro)
	return args.Error(0)
}

func (m *MockRegistroDiarioRepositorio) Deletar(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestRegistroDiarioService_CriarRegistro(t *testing.T) {
	mockRegistroRepo := new(MockRegistroDiarioRepositorio)
	mockDiarioRepo := new(MockDiarioCultivoRepository) // Assuming this mock implements the required interface methods
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
	mockRegistroRepo := new(MockRegistroDiarioRepositorio)
	mockDiarioRepo := new(MockDiarioCultivoRepository) // Not directly used in this test, but needed for service creation
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
