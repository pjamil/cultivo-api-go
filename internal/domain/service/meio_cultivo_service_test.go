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

type MockMeioCultivoRepositorio struct {
	mock.Mock
}

func (m *MockMeioCultivoRepositorio) Criar(meioCultivo *models.MeioCultivo) error {
	args := m.Called(meioCultivo)
	return args.Error(0)
}

func (m *MockMeioCultivoRepositorio) ListarTodos() ([]models.MeioCultivo, error) {
	args := m.Called()
	return args.Get(0).([]models.MeioCultivo), args.Error(1)
}

func (m *MockMeioCultivoRepositorio) BuscarPorID(id uint) (*models.MeioCultivo, error) {
	args := m.Called(id)
	return args.Get(0).(*models.MeioCultivo), args.Error(1)
}

func (m *MockMeioCultivoRepositorio) Atualizar(meioCultivo *models.MeioCultivo) error {
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
		createDTO := &dto.CreateMeioCultivoDTO{Nome: "Solo", Descricao: "Solo orgânico"}
		expectedMeioCultivo := &models.MeioCultivo{Nome: "Solo", Descricao: "Solo orgânico"}
		expectedResponse := &dto.MeioCultivoResponseDTO{ID: 1, Nome: "Solo", Descricao: "Solo orgânico"}

		mockRepo.On("Criar", mock.AnythingOfType("*models.MeioCultivo")).Run(func(args mock.Arguments) {
			meioCultivo := args.Get(0).(*models.MeioCultivo)
			meioCultivo.ID = 1
		}).Return(nil).Once()

		response, err := service.Criar(createDTO)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, expectedResponse.Nome, response.Nome)
		assert.Equal(t, expectedResponse.ID, response.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Repository Error", func(t *testing.T) {
		createDTO := &dto.CreateMeioCultivoDTO{Nome: "Solo", Descricao: "Solo orgânico"}

		mockRepo.On("Criar", mock.AnythingOfType("*models.MeioCultivo")).Return(errors.New("erro no repositório")).Once()

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
		expectedMeiosCultivo := []models.MeioCultivo{
			{ID: 1, Nome: "Solo"},
			{ID: 2, Nome: "Hidroponia"},
		}
		expectedResponse := []dto.MeioCultivoResponseDTO{
			{ID: 1, Nome: "Solo"},
			{ID: 2, Nome: "Hidroponia"},
		}

		mockRepo.On("ListarTodos").Return(expectedMeiosCultivo, nil).Once()

		response, err := service.ListarTodos()

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, expectedResponse, response)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success - Nenhum Meio de Cultivo Encontrado", func(t *testing.T) {
		mockRepo.On("ListarTodos").Return([]models.MeioCultivo{}, nil).Once()

		response, err := service.ListarTodos()

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Empty(t, response)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Repository Error", func(t *testing.T) {
		mockRepo.On("ListarTodos").Return([]models.MeioCultivo{}, errors.New("erro no repositório")).Once()

		response, err := service.ListarTodos()

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.EqualError(t, err, "erro no repositório")
		mockRepo.AssertExpectations(t)
	})
}

func TestMeioCultivoService_BuscarPorID(t *testing.T) {
	mockRepo := new(MockMeioCultivoRepositorio)
	service := service.NewMeioCultivoService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		meioCultivoID := uint(1)
		expectedMeioCultivo := &models.MeioCultivo{ID: meioCultivoID, Nome: "Solo"}
		expectedResponse := &dto.MeioCultivoResponseDTO{ID: meioCultivoID, Tipo: "Solo"}

		mockRepo.On("BuscarPorID", meioCultivoID).Return(expectedMeioCultivo, nil).Once()

		response, err := service.BuscarPorID(meioCultivoID)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, expectedResponse.ID, response.ID)
		assert.Equal(t, expectedResponse.Nome, response.Nome)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		meioCultivoID := uint(999)

		mockRepo.On("BuscarPorID", meioCultivoID).Return((*models.MeioCultivo)(nil), gorm.ErrRecordNotFound).Once()

		response, err := service.BuscarPorID(meioCultivoID)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		meioCultivoID := uint(1)
		expectedError := errors.New("erro no repositório")

		mockRepo.On("BuscarPorID", meioCultivoID).Return((*models.MeioCultivo)(nil), expectedError).Once()

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
		existingMeioCultivo := &models.MeioCultivo{ID: meioCultivoID, Tipo: "Solo Antigo"}
		expectedResponse := &dto.MeioCultivoResponseDTO{ID: meioCultivoID, Tipo: "Solo Atualizado"}

		mockRepo.On("BuscarPorID", meioCultivoID).Return(existingMeioCultivo, nil).Once()
		mockRepo.On("Atualizar", mock.AnythingOfType("*models.MeioCultivo")).Run(func(args mock.Arguments) {
			meioCultivo := args.Get(0).(*models.MeioCultivo)
									meioCultivo.Tipo = updateDTO.Tipo
		}).Return(nil).Once()

		response, err := service.Atualizar(meioCultivoID, updateDTO)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, expectedResponse.Nome, response.Nome)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		meioCultivoID := uint(999)
		updateDTO := &dto.UpdateMeioCultivoDTO{Nome: "Solo Atualizado"}

		mockRepo.On("BuscarPorID", meioCultivoID).Return((*models.MeioCultivo)(nil), gorm.ErrRecordNotFound).Once()

		response, err := service.Atualizar(meioCultivoID, updateDTO)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error on Update", func(t *testing.T) {
		meioCultivoID := uint(1)
		updateDTO := &dto.UpdateMeioCultivoDTO{Nome: "Solo Atualizado"}
		existingMeioCultivo := &models.MeioCultivo{ID: meioCultivoID, Tipo: "Solo Antigo"}
		expectedError := errors.New("erro no repositório")

		mockRepo.On("BuscarPorID", meioCultivoID).Return(existingMeioCultivo, nil).Once()
		mockRepo.On("Atualizar", mock.AnythingOfType("*models.MeioCultivo")).Return(expectedError).Once()

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
		assert.Equal(t, gorm.ErrRecordNotFound, err)
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
