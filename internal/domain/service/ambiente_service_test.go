package service_test

import (
	"encoding/json"
	"errors"
	"testing"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/entity"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/service"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/service/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestAmbienteService_Criar(t *testing.T) {
	mockRepo := new(test.MockAmbienteRepositorio)
	service := service.NewAmbienteService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		createDTO := &dto.CreateAmbienteDTO{Nome: "Ambiente Teste", Tipo: "interno", Comprimento: 10, Altura: 5, Largura: 8, TempoExposicao: 12}
		expectedResponse := &dto.AmbienteResponseDTO{ID: 1, Nome: "Ambiente Teste", Tipo: "interno", Comprimento: 10, Altura: 5, Largura: 8, TempoExposicao: 12, Orientacao: "Norte"}

		mockRepo.On("Criar", mock.AnythingOfType("*entity.Ambiente")).Run(func(args mock.Arguments) {
			ambiente := args.Get(0).(*entity.Ambiente)
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
		createDTO := &dto.CreateAmbienteDTO{Nome: "Ambiente Teste", Tipo: "interno", Comprimento: 10, Altura: 5, Largura: 8, TempoExposicao: 12}

		mockRepo.On("Criar", mock.AnythingOfType("*entity.Ambiente")).Return(errors.New("erro no repositório")).Once()

		response, err := service.Criar(createDTO)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.EqualError(t, err, "falha ao criar ambiente com nome Ambiente Teste: erro no repositório")
		mockRepo.AssertExpectations(t)
	})
}

func TestAmbienteService_ListarTodos(t *testing.T) {
	mockRepo := new(test.MockAmbienteRepositorio)
	service := service.NewAmbienteService(mockRepo)

	t.Run("Success - Ambientes Encontrados", func(t *testing.T) {
		// Arrange
		mockAmbientes := []entity.Ambiente{
			{Model: gorm.Model{ID: 1}, Nome: "Ambiente 1", Tipo: "interno"},
			{Model: gorm.Model{ID: 2}, Nome: "Ambiente 2", Tipo: "externo"},
		}
		expectedTotal := int64(len(mockAmbientes))
		page := 1
		limit := 10

		expectedResponseData := []dto.AmbienteResponseDTO{
			{ID: 1, Nome: "Ambiente 1", Tipo: "interno"},
			{ID: 2, Nome: "Ambiente 2", Tipo: "externo"},
		}

		mockRepo.On("ListarTodos", page, limit).Return(mockAmbientes, expectedTotal, nil).Once()

		// Act
		response, err := service.ListarTodos(page, limit)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, expectedTotal, response.Total)
		assert.Equal(t, page, response.Page)
		assert.Equal(t, limit, response.Limit)
		var actualResponseData []dto.AmbienteResponseDTO
		err = json.Unmarshal(response.Data, &actualResponseData)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponseData, actualResponseData)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success - Nenhum Ambiente Encontrado", func(t *testing.T) {
		page := 1
		limit := 10
		mockRepo.On("ListarTodos", page, limit).Return([]entity.Ambiente{}, int64(0), nil).Once()

		response, err := service.ListarTodos(page, limit)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, int64(0), response.Total)
		assert.Equal(t, json.RawMessage("[]"), response.Data)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Repository Error", func(t *testing.T) {
		page := 1
		limit := 10
		mockRepo.On("ListarTodos", page, limit).Return([]entity.Ambiente{}, int64(0), errors.New("erro no repositório")).Once()

		response, err := service.ListarTodos(page, limit)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.EqualError(t, err, "erro no repositório")
		mockRepo.AssertExpectations(t)
	})
}

func TestAmbienteService_BuscarPorID(t *testing.T) {
	mockRepo := new(test.MockAmbienteRepositorio)
	service := service.NewAmbienteService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		ambienteID := uint(1)
		expectedAmbiente := &entity.Ambiente{Nome: "Ambiente Teste", Tipo: "interno"}
		expectedAmbiente.ID = ambienteID
		expectedResponse := &dto.AmbienteResponseDTO{ID: ambienteID, Nome: "Ambiente Teste", Tipo: "interno"}

		mockRepo.On("BuscarPorID", ambienteID).Return(expectedAmbiente, nil).Once()

		response, err := service.BuscarPorID(ambienteID)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, expectedResponse.ID, response.ID)
		assert.Equal(t, expectedResponse.Nome, response.Nome)
		assert.Equal(t, expectedResponse.Tipo, response.Tipo)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		ambienteID := uint(999)

		mockRepo.On("BuscarPorID", ambienteID).Return((*entity.Ambiente)(nil), gorm.ErrRecordNotFound).Once()

		response, err := service.BuscarPorID(ambienteID)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		ambienteID := uint(1)
		expectedError := errors.New("erro no repositório")

		mockRepo.On("BuscarPorID", ambienteID).Return((*entity.Ambiente)(nil), expectedError).Once()

		response, err := service.BuscarPorID(ambienteID)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), expectedError.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestAmbienteService_Atualizar(t *testing.T) {
	mockRepo := new(test.MockAmbienteRepositorio)
	service := service.NewAmbienteService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		ambienteID := uint(1)
		updateDTO := &dto.UpdateAmbienteDTO{Nome: "Ambiente Atualizado"}
		existingAmbiente := &entity.Ambiente{Nome: "Ambiente Antigo", Tipo: "interno"}
		existingAmbiente.ID = ambienteID
		expectedResponse := &dto.AmbienteResponseDTO{ID: ambienteID, Nome: "Ambiente Atualizado", Tipo: "interno"}

		mockRepo.On("BuscarPorID", ambienteID).Return(existingAmbiente, nil).Once()
		mockRepo.On("Atualizar", mock.AnythingOfType("*entity.Ambiente")).Run(func(args mock.Arguments) {
			ambiente := args.Get(0).(*entity.Ambiente)
			ambiente.Nome = updateDTO.Nome
		}).Return(nil).Once()

		response, err := service.Atualizar(ambienteID, updateDTO)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, expectedResponse.Nome, response.Nome)
		assert.Equal(t, expectedResponse.Tipo, response.Tipo)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		ambienteID := uint(999)
		updateDTO := &dto.UpdateAmbienteDTO{Nome: "Ambiente Atualizado"}

		mockRepo.On("BuscarPorID", ambienteID).Return((*entity.Ambiente)(nil), gorm.ErrRecordNotFound).Once()

		response, err := service.Atualizar(ambienteID, updateDTO)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error on Update", func(t *testing.T) {
		ambienteID := uint(1)
		updateDTO := &dto.UpdateAmbienteDTO{Nome: "Ambiente Atualizado"}
		existingAmbiente := &entity.Ambiente{Nome: "Ambiente Antigo", Tipo: "interno"}
		existingAmbiente.ID = ambienteID
		expectedError := errors.New("erro no repositório")

		mockRepo.On("BuscarPorID", ambienteID).Return(existingAmbiente, nil).Once()
		mockRepo.On("Atualizar", mock.AnythingOfType("*entity.Ambiente")).Return(expectedError).Once()

		response, err := service.Atualizar(ambienteID, updateDTO)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), expectedError.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestAmbienteService_Deletar(t *testing.T) {
	mockRepo := new(test.MockAmbienteRepositorio)
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
