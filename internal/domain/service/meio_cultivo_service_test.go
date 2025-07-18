package service_test

import (
	"encoding/json"
	"errors"
	"testing"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/entity"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockMeioCultivoRepositorio struct {
	mock.Mock
}

func (m *MockMeioCultivoRepositorio) Criar(meioCultivo *entity.MeioCultivo) error {
	args := m.Called(meioCultivo)
	return args.Error(0)
}

func (m *MockMeioCultivoRepositorio) ListarTodos(page, limit int) ([]entity.MeioCultivo, int64, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]entity.MeioCultivo), args.Get(1).(int64), args.Error(2)
}

func (m *MockMeioCultivoRepositorio) BuscarPorID(id uint) (*entity.MeioCultivo, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.MeioCultivo), args.Error(1)
}

func (m *MockMeioCultivoRepositorio) Atualizar(meioCultivo *entity.MeioCultivo) error {
	args := m.Called(meioCultivo)
	return args.Error(0)
}

func (m *MockMeioCultivoRepositorio) Deletar(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestMeioCultivoService_Criar(t *testing.T) {
	mockRepo := new(MockMeioCultivoRepositorio)
	service := service.NewMeioCultivoService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		createDTO := &dto.CreateMeioCultivoDTO{Tipo: "Solo", Descricao: "Solo orgânico"}
		expectedResponse := &dto.MeioCultivoResponseDTO{ID: 1, Tipo: "Solo", Descricao: "Solo orgânico"}

		mockRepo.On("Criar", mock.AnythingOfType("*entity.MeioCultivo")).Run(func(args mock.Arguments) {
			meioCultivo := args.Get(0).(*entity.MeioCultivo)
			meioCultivo.ID = 1
		}).Return(nil).Once()

		response, err := service.Criar(createDTO)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, expectedResponse.Tipo, response.Tipo)
		assert.Equal(t, expectedResponse.ID, response.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Repository Error", func(t *testing.T) {
		createDTO := &dto.CreateMeioCultivoDTO{Tipo: "Solo", Descricao: "Solo orgânico"}

		mockRepo.On("Criar", mock.AnythingOfType("*entity.MeioCultivo")).Return(errors.New("erro no repositório")).Once()

		response, err := service.Criar(createDTO)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.EqualError(t, err, "erro no repositório")
		mockRepo.AssertExpectations(t)
	})
}

func TestMeioCultivoService_ListarTodos(t *testing.T) {
	mockRepo := new(MockMeioCultivoRepositorio)
	service := service.NewMeioCultivoService(mockRepo)

	t.Run("Success - Meios de Cultivo Encontrados", func(t *testing.T) {
		// Arrange
		mockMeiosCultivo := []entity.MeioCultivo{
			{Model: gorm.Model{ID: 1}, Tipo: "Solo", Descricao: "Solo Orgânico"},
			{Model: gorm.Model{ID: 2}, Tipo: "Hidroponia", Descricao: "Sistema NFT"},
		}
		expectedTotal := int64(len(mockMeiosCultivo))
		page := 1
		limit := 10

		expectedResponseData := []dto.MeioCultivoResponseDTO{
			{ID: 1, Tipo: "Solo", Descricao: "Solo Orgânico"},
			{ID: 2, Tipo: "Hidroponia", Descricao: "Sistema NFT"},
		}

		mockRepo.On("ListarTodos", page, limit).Return(mockMeiosCultivo, expectedTotal, nil).Once()

		// Act
		responseDTOs, total, err := service.ListarTodos(page, limit)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, responseDTOs)
		assert.Equal(t, expectedTotal, total)
		assert.Equal(t, expectedResponseData, responseDTOs)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success - Nenhum Meio de Cultivo Encontrado", func(t *testing.T) {
		page := 1
		limit := 10
		mockRepo.On("ListarTodos", page, limit).Return([]entity.MeioCultivo{}, int64(0), nil).Once()

		responseDTOs, total, err := service.ListarTodos(page, limit)

		assert.NoError(t, err)
		assert.NotNil(t, responseDTOs)
		assert.Equal(t, int64(0), total)
		assert.Empty(t, responseDTOs)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Repository Error", func(t *testing.T) {
		page := 1
		limit := 10
		mockRepo.On("ListarTodos", page, limit).Return([]entity.MeioCultivo{}, int64(0), errors.New("erro no repositório")).Once()

		responseDTOs, total, err := service.ListarTodos(page, limit)

		assert.Error(t, err)
		assert.Nil(t, responseDTOs)
		assert.Equal(t, int64(0), total)
		assert.EqualError(t, err, "erro no repositório")
		mockRepo.AssertExpectations(t)
	})
}

func TestMeioCultivoService_BuscarPorID(t *testing.T) {
	mockRepo := new(MockMeioCultivoRepositorio)
	service := service.NewMeioCultivoService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		meioCultivoID := uint(1)
		expectedMeioCultivo := &entity.MeioCultivo{Tipo: "Solo"}
		expectedMeioCultivo.ID = meioCultivoID
		expectedResponse := &dto.MeioCultivoResponseDTO{ID: meioCultivoID, Tipo: "Solo"}

		mockRepo.On("BuscarPorID", meioCultivoID).Return(expectedMeioCultivo, nil).Once()

		response, err := service.BuscarPorID(meioCultivoID)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, expectedResponse.ID, response.ID)
		assert.Equal(t, expectedResponse.Tipo, response.Tipo)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		meioCultivoID := uint(999)

		mockRepo.On("BuscarPorID", meioCultivoID).Return((*entity.MeioCultivo)(nil), gorm.ErrRecordNotFound).Once()

		response, err := service.BuscarPorID(meioCultivoID)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, service.ErrNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		meioCultivoID := uint(1)
		expectedError := errors.New("erro no repositório")

		mockRepo.On("BuscarPorID", meioCultivoID).Return((*entity.MeioCultivo)(nil), expectedError).Once()

		response, err := service.BuscarPorID(meioCultivoID)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), expectedError.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestMeioCultivoService_Atualizar(t *testing.T) {
	mockRepo := new(MockMeioCultivoRepositorio)
	service := service.NewMeioCultivoService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		meioCultivoID := uint(1)
		updateDTO := &dto.UpdateMeioCultivoDTO{Tipo: "Solo Atualizado"}
		existingMeioCultivo := &entity.MeioCultivo{Tipo: "Solo Antigo"}
		existingMeioCultivo.ID = meioCultivoID
		expectedResponse := &dto.MeioCultivoResponseDTO{ID: meioCultivoID, Tipo: "Solo Atualizado"}

		mockRepo.On("BuscarPorID", meioCultivoID).Return(existingMeioCultivo, nil).Once()
		mockRepo.On("Atualizar", mock.AnythingOfType("*entity.MeioCultivo")).Run(func(args mock.Arguments) {
			meioCultivo := args.Get(0).(*entity.MeioCultivo)
			meioCultivo.Tipo = updateDTO.Tipo
		}).Return(nil).Once()

		response, err := service.Atualizar(meioCultivoID, updateDTO)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, expectedResponse.Tipo, response.Tipo)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		meioCultivoID := uint(999)
		updateDTO := &dto.UpdateMeioCultivoDTO{Tipo: "Solo Atualizado"}

		mockRepo.On("BuscarPorID", meioCultivoID).Return((*entity.MeioCultivo)(nil), gorm.ErrRecordNotFound).Once()

		response, err := service.Atualizar(meioCultivoID, updateDTO)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, service.ErrNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error on Update", func(t *testing.T) {
		meioCultivoID := uint(1)
		updateDTO := &dto.UpdateMeioCultivoDTO{Tipo: "Solo Atualizado"}
		existingMeioCultivo := &entity.MeioCultivo{Tipo: "Solo Antigo"}
		existingMeioCultivo.ID = meioCultivoID
		expectedError := errors.New("erro no repositório")

		mockRepo.On("BuscarPorID", meioCultivoID).Return(existingMeioCultivo, nil).Once()
		mockRepo.On("Atualizar", mock.AnythingOfType("*entity.MeioCultivo")).Return(expectedError).Once()

		response, err := service.Atualizar(meioCultivoID, updateDTO)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), expectedError.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestMeioCultivoService_Deletar(t *testing.T) {
	mockRepo := new(MockMeioCultivoRepositorio)
	service := service.NewMeioCultivoService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		meioCultivoID := uint(1)
		mockRepo.On("Deletar", meioCultivoID).Return(nil).Once()

		err := service.Deletar(meioCultivoID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		meioCultivoID := uint(999)
		mockRepo.On("Deletar", meioCultivoID).Return(gorm.ErrRecordNotFound).Once()

		err := service.Deletar(meioCultivoID)

		assert.Error(t, err)
		assert.Equal(t, service.ErrNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		meioCultivoID := uint(1)
		expectedError := errors.New("erro no repositório")
		mockRepo.On("Deletar", meioCultivoID).Return(expectedError).Once()

		err := service.Deletar(meioCultivoID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), expectedError.Error())
		mockRepo.AssertExpectations(t)
	})
}
