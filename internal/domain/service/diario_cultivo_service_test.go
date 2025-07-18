package service_test

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/entity"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockDiarioCultivoRepository é um mock para a interface DiarioCultivoRepository.
type MockDiarioCultivoRepository struct {
	mock.Mock
}

func (m *MockDiarioCultivoRepository) Create(diario *entity.DiarioCultivo) error {
	args := m.Called(diario)
	return args.Error(0)
}

func (m *MockDiarioCultivoRepository) GetByID(id uint) (*entity.DiarioCultivo, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.DiarioCultivo), args.Error(1)
}

func (m *MockDiarioCultivoRepository) GetAll(page, limit int) ([]entity.DiarioCultivo, int64, error) {
	args := m.Called(page, limit)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]entity.DiarioCultivo), args.Get(1).(int64), args.Error(2)
}

func (m *MockDiarioCultivoRepository) GetAllByUserID(userID uint) ([]entity.DiarioCultivo, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.DiarioCultivo), args.Error(1)
}

func (m *MockDiarioCultivoRepository) Update(diario *entity.DiarioCultivo) error {
	args := m.Called(diario)
	return args.Error(0)
}

func (m *MockDiarioCultivoRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// Métodos de associação não utilizados pelo serviço principal, mas mockados para completar a interface.
func (m *MockDiarioCultivoRepository) AddPlantas(diarioCultivoID uint, plantas []*entity.Planta) error {
	return m.Called(diarioCultivoID, plantas).Error(0)
}
func (m *MockDiarioCultivoRepository) RemovePlantas(diarioCultivoID uint, plantas []*entity.Planta) error {
	return m.Called(diarioCultivoID, plantas).Error(0)
}
func (m *MockDiarioCultivoRepository) AddAmbientes(diarioCultivoID uint, ambientes []*entity.Ambiente) error {
	return m.Called(diarioCultivoID, ambientes).Error(0)
}
func (m *MockDiarioCultivoRepository) RemoveAmbientes(diarioCultivoID uint, ambientes []*entity.Ambiente) error {
	return m.Called(diarioCultivoID, ambientes).Error(0)
}

func TestDiarioCultivoService_Criar(t *testing.T) {
	mockRepo := new(MockDiarioCultivoRepository)
	diarioService := service.NewDiarioCultivoService(mockRepo)

	createDTO := dto.CreateDiarioCultivoDTO{
		Nome:        "Meu Diário de Cultivo",
		DataInicio:  time.Now(),
		UsuarioID:   1,
		Privacidade: "privado",
	}

	mockRepo.On("Create", mock.AnythingOfType("*entity.DiarioCultivo")).Run(func(args mock.Arguments) {
		diario := args.Get(0).(*entity.DiarioCultivo)
		diario.ID = 1
	}).Return(nil)

	responseDTO, err := diarioService.CreateDiario(createDTO)

	assert.NoError(t, err)
	assert.NotNil(t, responseDTO)
	assert.Equal(t, uint(1), responseDTO.ID)
	assert.Equal(t, createDTO.Nome, responseDTO.Nome)
	mockRepo.AssertExpectations(t)
}

func TestDiarioCultivoService_BuscarPorID(t *testing.T) {
	mockRepo := new(MockDiarioCultivoRepository)
	diarioService := service.NewDiarioCultivoService(mockRepo)

	diario := &entity.DiarioCultivo{
		Model: gorm.Model{ID: 1},
		Nome:  "Diário Teste",
	}

	t.Run("sucesso", func(t *testing.T) {
		mockRepo.On("GetByID", uint(1)).Return(diario, nil).Once()

		responseDTO, err := diarioService.GetDiarioByID(1)

		assert.NoError(t, err)
		assert.NotNil(t, responseDTO)
		assert.Equal(t, diario.ID, responseDTO.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("não encontrado", func(t *testing.T) {
		mockRepo.On("GetByID", uint(2)).Return(nil, gorm.ErrRecordNotFound).Once()

		responseDTO, err := diarioService.GetDiarioByID(2)

		assert.Error(t, err)
		assert.Nil(t, responseDTO)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestDiarioCultivoService_ListarTodos(t *testing.T) {
	mockRepo := new(MockDiarioCultivoRepository)
	diarioService := service.NewDiarioCultivoService(mockRepo)

	diarios := []entity.DiarioCultivo{
		{Model: gorm.Model{ID: 1}, Nome: "Diário 1"},
		{Model: gorm.Model{ID: 2}, Nome: "Diário 2"},
	}
	total := int64(2)

	mockRepo.On("GetAll", 1, 10).Return(diarios, total, nil).Once()

	paginatedResponse, err := diarioService.GetAllDiarios(1, 10)

	assert.NoError(t, err)
	assert.NotNil(t, paginatedResponse)
	assert.Equal(t, total, paginatedResponse.Total)

	var actualDiarios []entity.DiarioCultivo
	err = json.Unmarshal(paginatedResponse.Data, &actualDiarios)
	assert.NoError(t, err)
	assert.Equal(t, diarios, actualDiarios)
	mockRepo.AssertExpectations(t)
}

func TestDiarioCultivoService_Atualizar(t *testing.T) {
	mockRepo := new(MockDiarioCultivoRepository)
	diarioService := service.NewDiarioCultivoService(mockRepo)

	updateDTO := &dto.UpdateDiarioCultivoDTO{
		Nome: "Diário Atualizado",
	}

	existingDiario := &entity.DiarioCultivo{
		Model: gorm.Model{ID: 1},
		Nome:  "Diário Antigo",
	}

	t.Run("sucesso", func(t *testing.T) {
		mockRepo.On("GetByID", uint(1)).Return(existingDiario, nil).Twice()
		mockRepo.On("Update", mock.AnythingOfType("*entity.DiarioCultivo")).Return(nil).Once()

		t.Logf("Calling diarioService.Update with ID: %d", 1)
		responseDTO, err := diarioService.UpdateDiario(1, *updateDTO)

		assert.NoError(t, err)
		assert.NotNil(t, responseDTO)
		assert.Equal(t, updateDTO.Nome, responseDTO.Nome)
		mockRepo.AssertExpectations(t)
	})

	t.Run("não encontrado", func(t *testing.T) {
		mockRepo.On("GetByID", uint(2)).Return(nil, gorm.ErrRecordNotFound).Once()

		t.Logf("Calling diarioService.Update with ID: %d", 2)
		responseDTO, err := diarioService.UpdateDiario(2, *updateDTO)

		assert.Error(t, err)
		assert.Nil(t, responseDTO)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestDiarioCultivoService_Deletar(t *testing.T) {
	mockRepo := new(MockDiarioCultivoRepository)
	diarioService := service.NewDiarioCultivoService(mockRepo)

	t.Run("sucesso", func(t *testing.T) {
		mockRepo.On("Delete", uint(1)).Return(nil).Once()

		err := diarioService.DeleteDiario(1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro no repositório", func(t *testing.T) {
		expectedError := errors.New("erro ao deletar")
		mockRepo.On("Delete", uint(2)).Return(expectedError).Once()

		err := diarioService.DeleteDiario(2)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})
}