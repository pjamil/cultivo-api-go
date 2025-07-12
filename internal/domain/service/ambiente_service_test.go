package service_test

import (
	"encoding/json"
	"errors"
	"testing"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockAmbienteRepositorio struct {
	mock.Mock
}

func (m *MockAmbienteRepositorio) Criar(ambiente *models.Ambiente) error {
	args := m.Called(ambiente)
	return args.Error(0)
}

func (m *MockAmbienteRepositorio) ListarTodos(page, limit int) ([]models.Ambiente, int64, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]models.Ambiente), args.Get(1).(int64), args.Error(2)
}

func (m *MockAmbienteRepositorio) BuscarPorID(id uint) (*models.Ambiente, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Ambiente), args.Error(1)
}

func (m *MockAmbienteRepositorio) Atualizar(ambiente *models.Ambiente) error {
	args := m.Called(ambiente)
	return args.Error(0)
}

func (m *MockAmbienteRepositorio) Deletar(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestAmbienteService_Criar(t *testing.T) {
	mockRepo := new(MockAmbienteRepositorio)
	service := service.NewAmbienteService(mockRepo)

	t.Run("Success", func(t *testing.T) {
												createDTO := &dto.CreateAmbienteDTO{Nome: "Ambiente Teste", Descricao: "Descrição do ambiente"}
		expectedResponse := &dto.AmbienteResponseDTO{ID: 1, Nome: "Ambiente Teste", Descricao: "Descrição do ambiente"}

		mockRepo.On("Criar", mock.AnythingOfType("*models.Ambiente")).Run(func(args mock.Arguments) {
			ambiente := args.Get(0).(*models.Ambiente)
			ambiente.ID = 1
		}).Return(nil).Once()

		response, err := service.Criar(createDTO)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, expectedResponse.Nome, response.Nome)
		assert.Equal(t, expectedResponse.ID, response.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Repository Error", func(t *testing.T) {
		createDTO := &dto.CreateAmbienteDTO{Nome: "Ambiente Teste", Descricao: "Descrição do ambiente"}

		mockRepo.On("Criar", mock.AnythingOfType("*models.Ambiente")).Return(errors.New("erro no repositório")).Once()

		response, err := service.Criar(createDTO)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.EqualError(t, err, "erro no repositório")
		mockRepo.AssertExpectations(t)
	})
}

func TestAmbienteService_ListarTodos(t *testing.T) {
	mockRepo := new(MockAmbienteRepositorio)
	service := service.NewAmbienteService(mockRepo)

	t.Run("Success - Ambientes Encontrados", func(t *testing.T) {
		expectedAmbientes := []models.Ambiente{
			{Nome: "Ambiente 1"},
			{Nome: "Ambiente 2"},
		}
		for i := range expectedAmbientes {
			expectedAmbientes[i].ID = uint(i + 1)
		}
		expectedTotal := int64(len(expectedAmbientes))
		page := 1
		limit := 10

		mockRepo.On("ListarTodos", page, limit).Return(expectedAmbientes, expectedTotal, nil).Once()

		response, err := service.ListarTodos(page, limit)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, expectedTotal, response.Total)
		assert.Equal(t, page, response.Page)
		assert.Equal(t, limit, response.Limit)

		actualAmbientesBytes, _ := json.Marshal(response.Data)
		var actualAmbientes []models.Ambiente
		json.Unmarshal(actualAmbientesBytes, &actualAmbientes)

		assert.Equal(t, expectedAmbientes, actualAmbientes)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success - Nenhum Ambiente Encontrado", func(t *testing.T) {
		page := 1
		limit := 10
		mockRepo.On("ListarTodos", page, limit).Return([]models.Ambiente{}, int64(0), nil).Once()

		response, err := service.ListarTodos(page, limit)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, int64(0), response.Total)
		assert.Empty(t, response.Data)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Repository Error", func(t *testing.T) {
		page := 1
		limit := 10
		mockRepo.On("ListarTodos", page, limit).Return([]models.Ambiente{}, int64(0), errors.New("erro no repositório")).Once()

		response, err := service.ListarTodos(page, limit)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.EqualError(t, err, "erro no repositório")
		mockRepo.AssertExpectations(t)
	})
}

func TestAmbienteService_BuscarPorID(t *testing.T) {
	mockRepo := new(MockAmbienteRepositorio)
	service := service.NewAmbienteService(mockRepo)

		t.Run("Success", func(t *testing.T) {
		ambienteID := uint(1)
		expectedAmbiente := &models.Ambiente{Nome: "Ambiente Teste"}
		expectedAmbiente.ID = ambienteID
		expectedResponse := &dto.AmbienteResponseDTO{ID: ambienteID, Nome: "Ambiente Teste"}

		mockRepo.On("BuscarPorID", ambienteID).Return(expectedAmbiente, nil).Once()

		response, err := service.BuscarPorID(ambienteID)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, expectedResponse.ID, response.ID)
		assert.Equal(t, expectedResponse.Nome, response.Nome)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		ambienteID := uint(999)

		mockRepo.On("BuscarPorID", ambienteID).Return((*models.Ambiente)(nil), gorm.ErrRecordNotFound).Once()

		response, err := service.BuscarPorID(ambienteID)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		ambienteID := uint(1)
		expectedError := errors.New("erro no repositório")

		mockRepo.On("BuscarPorID", ambienteID).Return((*models.Ambiente)(nil), expectedError).Once()

		response, err := service.BuscarPorID(ambienteID)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), expectedError.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestAmbienteService_Atualizar(t *testing.T) {
	mockRepo := new(MockAmbienteRepositorio)
	service := service.NewAmbienteService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		ambienteID := uint(1)
		updateDTO := &dto.UpdateAmbienteDTO{Nome: "Ambiente Atualizado"}
		existingAmbiente := &models.Ambiente{Nome: "Ambiente Antigo"}
		existingAmbiente.ID = ambienteID
		expectedResponse := &dto.AmbienteResponseDTO{ID: ambienteID, Nome: "Ambiente Atualizado"}

		mockRepo.On("BuscarPorID", ambienteID).Return(existingAmbiente, nil).Once()
		mockRepo.On("Atualizar", mock.AnythingOfType("*models.Ambiente")).Run(func(args mock.Arguments) {
			ambiente := args.Get(0).(*models.Ambiente)
			ambiente.Nome = updateDTO.Nome
		}).Return(nil).Once()

		response, err := service.Atualizar(ambienteID, updateDTO)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, expectedResponse.Nome, response.Nome)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		ambienteID := uint(999)
		updateDTO := &dto.UpdateAmbienteDTO{Nome: "Ambiente Atualizado"}

		mockRepo.On("BuscarPorID", ambienteID).Return((*models.Ambiente)(nil), gorm.ErrRecordNotFound).Once()

		response, err := service.Atualizar(ambienteID, updateDTO)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		mockRepo.AssertExpectations(t)
	})

		t.Run("Repository Error on Update", func(t *testing.T) {
				ambienteID := uint(1)
		updateDTO := &dto.UpdateAmbienteDTO{Nome: "Ambiente Atualizado"}
		existingAmbiente := &models.Ambiente{Nome: "Ambiente Antigo"}
		existingAmbiente.ID = ambienteID
		expectedError := errors.New("erro no repositório")

		mockRepo.On("BuscarPorID", ambienteID).Return(existingAmbiente, nil).Once()
		mockRepo.On("Atualizar", mock.AnythingOfType("*models.Ambiente")).Return(expectedError).Once()

		response, err := service.Atualizar(ambienteID, updateDTO)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), expectedError.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestAmbienteService_Deletar(t *testing.T) {
	mockRepo := new(MockAmbienteRepositorio)
	service := service.NewAmbienteService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		ambienteID := uint(1)
		mockRepo.On("Deletar", ambienteID).Return(nil).Once()

		err := service.Deletar(ambienteID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		ambienteID := uint(999)
		mockRepo.On("Deletar", ambienteID).Return(gorm.ErrRecordNotFound).Once()

		err := service.Deletar(ambienteID)

		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		ambienteID := uint(1)
		expectedError := errors.New("erro no repositório")
		mockRepo.On("Deletar", ambienteID).Return(expectedError).Once()

		err := service.Deletar(ambienteID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), expectedError.Error())
		mockRepo.AssertExpectations(t)
	})
}