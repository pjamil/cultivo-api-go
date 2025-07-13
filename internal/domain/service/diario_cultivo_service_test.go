package service_test

import (
	"errors"
	"testing"
	"time"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/service"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/service/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockDiarioCultivoRepository é um mock para a interface DiarioCultivoRepository.
type MockDiarioCultivoRepository struct {
	mock.Mock
}

func (m *MockDiarioCultivoRepository) Create(diario *models.DiarioCultivo) error {
	args := m.Called(diario)
	return args.Error(0)
}

func (m *MockDiarioCultivoRepository) GetByID(id uint) (*models.DiarioCultivo, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DiarioCultivo), args.Error(1)
}

func (m *MockDiarioCultivoRepository) GetAll(page, limit int) ([]models.DiarioCultivo, int64, error) {
	args := m.Called(page, limit)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]models.DiarioCultivo), args.Get(1).(int64), args.Error(2)
}

func (m *MockDiarioCultivoRepository) Update(diario *models.DiarioCultivo) error {
	args := m.Called(diario)
	return args.Error(0)
}

func (m *MockDiarioCultivoRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// Métodos de associação não utilizados pelo serviço principal, mas mockados para completar a interface.
func (m *MockDiarioCultivoRepository) AddPlantas(diarioCultivoID uint, plantas []*models.Planta) error {
	return m.Called(diarioCultivoID, plantas).Error(0)
}
func (m *MockDiarioCultivoRepository) RemovePlantas(diarioCultivoID uint, plantas []*models.Planta) error {
	return m.Called(diarioCultivoID, plantas).Error(0)
}
func (m *MockDiarioCultivoRepository) AddAmbientes(diarioCultivoID uint, ambientes []*models.Ambiente) error {
	return m.Called(diarioCultivoID, ambientes).Error(0)
}
func (m *MockDiarioCultivoRepository) RemoveAmbientes(diarioCultivoID uint, ambientes []*models.Ambiente) error {
	return m.Called(diarioCultivoID, ambientes).Error(0)
}

func TestDiarioCultivoService_Criar(t *testing.T) {
	mockRepo := new(MockDiarioCultivoRepository)
	mockPlantaRepo := new(test.MockPlantaRepositorio)
	mockAmbienteRepo := new(test.MockAmbienteRepositorio)
	diarioService := service.NewDiarioCultivoService(mockRepo, mockPlantaRepo, mockAmbienteRepo)

	createDTO := &dto.CreateDiarioCultivoDTO{
		Nome:        "Meu Diário de Cultivo",
		DataInicio:  time.Now(),
		UsuarioID:   1,
		Privacidade: "privado",
	}

	mockRepo.On("Create", mock.AnythingOfType("*models.DiarioCultivo")).Run(func(args mock.Arguments) {
		diario := args.Get(0).(*models.DiarioCultivo)
		diario.ID = 1
	}).Return(nil)
	mockRepo.On("GetByID", mock.AnythingOfType("uint")).Return(&models.DiarioCultivo{Model: gorm.Model{ID: 1}, Nome: createDTO.Nome, UsuarioID: createDTO.UsuarioID}, nil)

	responseDTO, err := diarioService.Create(createDTO)

	assert.NoError(t, err)
	assert.NotNil(t, responseDTO)
	assert.Equal(t, uint(1), responseDTO.ID)
	assert.Equal(t, createDTO.Nome, responseDTO.Nome)
	mockRepo.AssertExpectations(t)
}

func TestDiarioCultivoService_BuscarPorID(t *testing.T) {
	mockRepo := new(MockDiarioCultivoRepository)
	mockPlantaRepo := new(test.MockPlantaRepositorio)
	mockAmbienteRepo := new(test.MockAmbienteRepositorio)
	diarioService := service.NewDiarioCultivoService(mockRepo, mockPlantaRepo, mockAmbienteRepo)

	diario := &models.DiarioCultivo{
		Model: gorm.Model{ID: 1},
		Nome:  "Diário Teste",
	}

	t.Run("sucesso", func(t *testing.T) {
		mockRepo.On("GetByID", uint(1)).Return(diario, nil).Once()

		responseDTO, err := diarioService.GetByID(1)

		assert.NoError(t, err)
		assert.NotNil(t, responseDTO)
		assert.Equal(t, diario.ID, responseDTO.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("não encontrado", func(t *testing.T) {
		mockRepo.On("GetByID", uint(2)).Return(nil, gorm.ErrRecordNotFound).Once()

		responseDTO, err := diarioService.GetByID(2)

		assert.Error(t, err)
		assert.Nil(t, responseDTO)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestDiarioCultivoService_ListarTodos(t *testing.T) {
	mockRepo := new(MockDiarioCultivoRepository)
	mockPlantaRepo := new(test.MockPlantaRepositorio)
	mockAmbienteRepo := new(test.MockAmbienteRepositorio)
	diarioService := service.NewDiarioCultivoService(mockRepo, mockPlantaRepo, mockAmbienteRepo)

	diarios := []models.DiarioCultivo{
		{Model: gorm.Model{ID: 1}, Nome: "Diário 1"},
		{Model: gorm.Model{ID: 2}, Nome: "Diário 2"},
	}
	total := int64(2)

	mockRepo.On("GetAll", 1, 10).Return(diarios, total, nil).Once()

	paginatedResponse, err := diarioService.GetAll(1, 10)

	assert.NoError(t, err)
	assert.NotNil(t, paginatedResponse)
	assert.Equal(t, total, paginatedResponse.Total)
	assert.Len(t, paginatedResponse.Data, 2)
	mockRepo.AssertExpectations(t)
}

func TestDiarioCultivoService_Atualizar(t *testing.T) {
	mockRepo := new(MockDiarioCultivoRepository)
	mockPlantaRepo := new(test.MockPlantaRepositorio)
	mockAmbienteRepo := new(test.MockAmbienteRepositorio)
	diarioService := service.NewDiarioCultivoService(mockRepo, mockPlantaRepo, mockAmbienteRepo)

	updateDTO := &dto.UpdateDiarioCultivoDTO{
		Nome: "Diário Atualizado",
	}

	existingDiario := &models.DiarioCultivo{
		Model: gorm.Model{ID: 1},
		Nome:  "Diário Antigo",
	}

	t.Run("sucesso", func(t *testing.T) {
		mockRepo.On("GetByID", uint(1)).Return(existingDiario, nil).Once()
		mockRepo.On("Update", mock.AnythingOfType("*models.DiarioCultivo")).Return(nil).Once()

		responseDTO, err := diarioService.Update(1, updateDTO)

		assert.NoError(t, err)
		assert.NotNil(t, responseDTO)
		assert.Equal(t, updateDTO.Nome, responseDTO.Nome)
		mockRepo.AssertExpectations(t)
	})

	t.Run("não encontrado", func(t *testing.T) {
		mockRepo.On("GetByID", uint(2)).Return(nil, gorm.ErrRecordNotFound).Once()

		responseDTO, err := diarioService.Update(2, updateDTO)

		assert.Error(t, err)
		assert.Nil(t, responseDTO)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestDiarioCultivoService_Deletar(t *testing.T) {
	mockRepo := new(MockDiarioCultivoRepository)
	mockPlantaRepo := new(test.MockPlantaRepositorio)
	mockAmbienteRepo := new(test.MockAmbienteRepositorio)
	diarioService := service.NewDiarioCultivoService(mockRepo, mockPlantaRepo, mockAmbienteRepo)

	t.Run("sucesso", func(t *testing.T) {
		mockRepo.On("Delete", uint(1)).Return(nil).Once()

		err := diarioService.Delete(1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro no repositório", func(t *testing.T) {
		expectedError := errors.New("erro ao deletar")
		mockRepo.On("Delete", uint(2)).Return(expectedError).Once()

		err := diarioService.Delete(2)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})
}
