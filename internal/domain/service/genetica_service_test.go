package service_test

import (
	"errors"
	"testing"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockGeneticaRepositorio struct {
	mock.Mock
}

func (m *MockGeneticaRepositorio) Criar(genetica *models.Genetica) error {
	args := m.Called(genetica)
	return args.Error(0)
}

func (m *MockGeneticaRepositorio) ListarTodos() ([]models.Genetica, error) {
	args := m.Called()
	return args.Get(0).([]models.Genetica), args.Error(1)
}

func (m *MockGeneticaRepositorio) BuscarPorID(id uint) (*models.Genetica, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Genetica), args.Error(1)
}

func (m *MockGeneticaRepositorio) Atualizar(genetica *models.Genetica) error {
	args := m.Called(genetica)
	return args.Error(0)
}

func (m *MockGeneticaRepositorio) Deletar(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestGeneticaService_Criar(t *testing.T) {
	mockRepo := new(MockGeneticaRepositorio)
	service := service.NewGeneticaService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		createDTO := &dto.CreateGeneticaDTO{Nome: "Genetica Teste", TipoGenetica: "indica", TipoEspecie: "sativa", TempoFloracao: 8, Origem: "Brasil"}
		expectedResponse := &dto.GeneticaResponseDTO{ID: 1, Nome: "Genetica Teste", TipoGenetica: "indica", TipoEspecie: "sativa", TempoFloracao: 8, Origem: "Brasil"}

		mockRepo.On("Criar", mock.AnythingOfType("*models.Genetica")).Run(func(args mock.Arguments) {
			genetica := args.Get(0).(*models.Genetica)
			genetica.ID = 1
		}).Return(nil).Once()

		response, err := service.Criar(createDTO)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, expectedResponse.Nome, response.Nome)
		assert.Equal(t, expectedResponse.ID, response.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Repository Error", func(t *testing.T) {
		createDTO := &dto.CreateGeneticaDTO{Nome: "Genetica Teste", TipoGenetica: "indica", TipoEspecie: "sativa", TempoFloracao: 8, Origem: "Brasil"}

		mockRepo.On("Criar", mock.AnythingOfType("*models.Genetica")).Return(errors.New("erro no repositório")).Once()

		response, err := service.Criar(createDTO)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.EqualError(t, err, "erro no repositório")
		mockRepo.AssertExpectations(t)
	})
}

func TestGeneticaService_ListarTodas(t *testing.T) {
	mockRepo := new(MockGeneticaRepositorio)
	service := service.NewGeneticaService(mockRepo)

	t.Run("Success - Geneticas Encontradas", func(t *testing.T) {
		expectedGeneticas := []models.Genetica{
			{Nome: "Genetica 1"},
			{Nome: "Genetica 2"},
		}
		expectedGeneticas[0].ID = 1
		expectedGeneticas[1].ID = 2
		expectedResponse := []dto.GeneticaResponseDTO{
			{ID: 1, Nome: "Genetica 1"},
			{ID: 2, Nome: "Genetica 2"},
		}

		mockRepo.On("ListarTodos").Return(expectedGeneticas, nil).Once()

		response, err := service.ListarTodas()

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, expectedResponse, response)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success - Nenhuma Genetica Encontrada", func(t *testing.T) {
		mockRepo.On("ListarTodos").Return([]models.Genetica{}, nil).Once()

		response, err := service.ListarTodas()

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Empty(t, response)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Repository Error", func(t *testing.T) {
		mockRepo.On("ListarTodos").Return([]models.Genetica{}, errors.New("erro no repositório")).Once()

		response, err := service.ListarTodas()

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.EqualError(t, err, "erro no repositório")
		mockRepo.AssertExpectations(t)
	})
}

func TestGeneticaService_BuscarPorID(t *testing.T) {
	mockRepo := new(MockGeneticaRepositorio)
	service := service.NewGeneticaService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		geneticaID := uint(1)
		expectedGenetica := &models.Genetica{Nome: "Genetica Teste"}
		expectedGenetica.ID = geneticaID
		expectedResponse := &dto.GeneticaResponseDTO{ID: geneticaID, Nome: "Genetica Teste"}

		mockRepo.On("BuscarPorID", geneticaID).Return(expectedGenetica, nil).Once()

		response, err := service.BuscarPorID(geneticaID)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, expectedResponse.ID, response.ID)
		assert.Equal(t, expectedResponse.Nome, response.Nome)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		geneticaID := uint(999)

		mockRepo.On("BuscarPorID", geneticaID).Return((*models.Genetica)(nil), gorm.ErrRecordNotFound).Once()

		response, err := service.BuscarPorID(geneticaID)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		geneticaID := uint(1)
		expectedError := errors.New("erro no repositório")

		mockRepo.On("BuscarPorID", geneticaID).Return((*models.Genetica)(nil), expectedError).Once()

		response, err := service.BuscarPorID(geneticaID)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), expectedError.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestGeneticaService_Atualizar(t *testing.T) {
	mockRepo := new(MockGeneticaRepositorio)
	service := service.NewGeneticaService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		geneticaID := uint(1)
		updateDTO := &dto.UpdateGeneticaDTO{Nome: "Genetica Atualizada"}
		existingGenetica := &models.Genetica{Nome: "Genetica Antiga"}
		existingGenetica.ID = geneticaID
		expectedResponse := &dto.GeneticaResponseDTO{ID: geneticaID, Nome: "Genetica Atualizada"}

		mockRepo.On("BuscarPorID", geneticaID).Return(existingGenetica, nil).Once()
		mockRepo.On("Atualizar", mock.AnythingOfType("*models.Genetica")).Run(func(args mock.Arguments) {
			genetica := args.Get(0).(*models.Genetica)
			genetica.Nome = updateDTO.Nome
		}).Return(nil).Once()

		response, err := service.Atualizar(geneticaID, updateDTO)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, expectedResponse.Nome, response.Nome)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		geneticaID := uint(999)
		updateDTO := &dto.UpdateGeneticaDTO{Nome: "Genetica Atualizada"}

		mockRepo.On("BuscarPorID", geneticaID).Return((*models.Genetica)(nil), gorm.ErrRecordNotFound).Once()

		response, err := service.Atualizar(geneticaID, updateDTO)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error on Update", func(t *testing.T) {
		geneticaID := uint(1)
		updateDTO := &dto.UpdateGeneticaDTO{Nome: "Genetica Atualizada"}
		existingGenetica := &models.Genetica{Nome: "Genetica Antiga"}
		existingGenetica.ID = geneticaID
		expectedError := errors.New("erro no repositório")

		mockRepo.On("BuscarPorID", geneticaID).Return(existingGenetica, nil).Once()
		mockRepo.On("Atualizar", mock.AnythingOfType("*models.Genetica")).Return(expectedError).Once()

		response, err := service.Atualizar(geneticaID, updateDTO)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), expectedError.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestGeneticaService_Deletar(t *testing.T) {
	mockRepo := new(MockGeneticaRepositorio)
	service := service.NewGeneticaService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		geneticaID := uint(1)
		mockRepo.On("Deletar", geneticaID).Return(nil).Once()

		err := service.Deletar(geneticaID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		geneticaID := uint(999)
		mockRepo.On("Deletar", geneticaID).Return(gorm.ErrRecordNotFound).Once()

		err := service.Deletar(geneticaID)

		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		geneticaID := uint(1)
		expectedError := errors.New("erro no repositório")
		mockRepo.On("Deletar", geneticaID).Return(expectedError).Once()

		err := service.Deletar(geneticaID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), expectedError.Error())
		mockRepo.AssertExpectations(t)
	})
}