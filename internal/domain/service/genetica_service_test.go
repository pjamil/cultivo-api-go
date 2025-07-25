package service_test

import (
	"errors"
	"testing"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/entity"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/service"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockGeneticaRepositorio struct {
	mock.Mock
}

func (m *MockGeneticaRepositorio) Criar(genetica *entity.Genetica) error {
	args := m.Called(genetica)
	return args.Error(0)
}

func (m *MockGeneticaRepositorio) ListarTodos(page, limit int) ([]entity.Genetica, int64, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]entity.Genetica), args.Get(1).(int64), args.Error(2)
}

func (m *MockGeneticaRepositorio) BuscarPorID(id uint) (*entity.Genetica, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.Genetica), args.Error(1)
}

func (m *MockGeneticaRepositorio) Atualizar(genetica *entity.Genetica) error {
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

		mockRepo.On("Criar", mock.AnythingOfType("*entity.Genetica")).Run(func(args mock.Arguments) {
			genetica := args.Get(0).(*entity.Genetica)
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

		mockRepo.On("Criar", mock.AnythingOfType("*entity.Genetica")).Return(errors.New("erro no repositório")).Once()

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
		// Arrange
		mockGeneticas := []entity.Genetica{
			{Model: gorm.Model{ID: 1}, Nome: "Genetica 1", TipoGenetica: "indica", TipoEspecie: "sativa", Origem: "Brasil"},
			{Model: gorm.Model{ID: 2}, Nome: "Genetica 2", TipoGenetica: "sativa", TipoEspecie: "indica", Origem: "Afeganistão"},
		}
		expectedTotal := int64(len(mockGeneticas))
		page := 1
		limit := 10

		expectedResponseData := []dto.GeneticaResponseDTO{
			{ID: 1, Nome: "Genetica 1", TipoGenetica: "indica", TipoEspecie: "sativa", Origem: "Brasil"},
			{ID: 2, Nome: "Genetica 2", TipoGenetica: "sativa", TipoEspecie: "indica", Origem: "Afeganistão"},
		}

		mockRepo.On("ListarTodos", page, limit).Return(mockGeneticas, expectedTotal, nil).Once()

		// Act
		responseDTOs, total, err := service.ListarTodas(page, limit)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, responseDTOs)
		assert.Equal(t, expectedTotal, total)
		assert.Equal(t, expectedResponseData, responseDTOs)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success - Nenhuma Genetica Encontrada", func(t *testing.T) {
		page := 1
		limit := 10
		mockRepo.On("ListarTodos", page, limit).Return([]entity.Genetica{}, int64(0), nil).Once()

		responseDTOs, total, err := service.ListarTodas(page, limit)

		assert.NoError(t, err)
		assert.NotNil(t, responseDTOs)
		assert.Equal(t, int64(0), total)
		assert.Empty(t, responseDTOs)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Repository Error", func(t *testing.T) {
		page := 1
		limit := 10
		mockRepo.On("ListarTodos", page, limit).Return([]entity.Genetica{}, int64(0), errors.New("erro no repositório")).Once()

		responseDTOs, total, err := service.ListarTodas(page, limit)

		assert.Error(t, err)
		assert.Nil(t, responseDTOs)
		assert.Equal(t, int64(0), total)
		assert.EqualError(t, err, "erro no repositório")
		mockRepo.AssertExpectations(t)
	})
}

func TestGeneticaService_BuscarPorID(t *testing.T) {
	mockRepo := new(MockGeneticaRepositorio)
	service := service.NewGeneticaService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		geneticaID := uint(1)
		expectedGenetica := &entity.Genetica{Nome: "Genetica Teste"}
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

		mockRepo.On("BuscarPorID", geneticaID).Return((*entity.Genetica)(nil), gorm.ErrRecordNotFound).Once()

		response, err := service.BuscarPorID(geneticaID)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, utils.ErrNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		geneticaID := uint(1)
		expectedError := errors.New("erro no repositório")

		mockRepo.On("BuscarPorID", geneticaID).Return((*entity.Genetica)(nil), expectedError).Once()

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
		existingGenetica := &entity.Genetica{Nome: "Genetica Antiga"}
		existingGenetica.ID = geneticaID
		expectedResponse := &dto.GeneticaResponseDTO{ID: geneticaID, Nome: "Genetica Atualizada"}

		mockRepo.On("BuscarPorID", geneticaID).Return(existingGenetica, nil).Once()
		mockRepo.On("Atualizar", mock.AnythingOfType("*entity.Genetica")).Run(func(args mock.Arguments) {
			genetica := args.Get(0).(*entity.Genetica)
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

		mockRepo.On("BuscarPorID", geneticaID).Return((*entity.Genetica)(nil), gorm.ErrRecordNotFound).Once()

		response, err := service.Atualizar(geneticaID, updateDTO)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, utils.ErrNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error on Update", func(t *testing.T) {
		geneticaID := uint(1)
		updateDTO := &dto.UpdateGeneticaDTO{Nome: "Genetica Atualizada"}
		existingGenetica := &entity.Genetica{Nome: "Genetica Antiga"}
		existingGenetica.ID = geneticaID
		expectedError := errors.New("erro no repositório")

		mockRepo.On("BuscarPorID", geneticaID).Return(existingGenetica, nil).Once()
		mockRepo.On("Atualizar", mock.AnythingOfType("*entity.Genetica")).Return(expectedError).Once()

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
		assert.Equal(t, utils.ErrNotFound, err)
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
