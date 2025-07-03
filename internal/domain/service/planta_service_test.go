package service_test

import (
	"errors"
	"testing"
	"time"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)



func TestPlantaService_Criar(t *testing.T) {
	mockPlantaRepo := new(MockPlantaRepositorio)
	mockGeneticaRepo := new(MockGeneticaRepositorio)
	mockAmbienteRepo := new(MockAmbienteRepositorio)
	mockMeioRepo := new(MockMeioCultivoRepositorio)

	servico := service.NewPlantaService(mockPlantaRepo, mockGeneticaRepo, mockAmbienteRepo, mockMeioRepo, mockPlantaRepo) // Usando mockPlantaRepo para registroDiarioRepositorio

	t.Run("Success", func(t *testing.T) {
		now := time.Now()
		notas := "Algumas notas."
		plantaDto := &dto.CreatePlantaDTO{
			Nome:          "Planta Teste",
			GeneticaID:    1,
			AmbienteID:    1,
			MeioCultivoID: 1,
			Especie:       "Sativa",
			Status:        "ativa",
			DataPlantio:   now,
			Notas:         notas,
			ComecandoDe:   "semente",
			UsuarioID:     1,
		}

		mockPlantaRepo.On("ExistePorNome", plantaDto.Nome).Return(false).Once()
		mockGeneticaRepo.On("BuscarPorID", plantaDto.GeneticaID).Return(&models.Genetica{}, nil).Once()
		mockAmbienteRepo.On("BuscarPorID", plantaDto.AmbienteID).Return(&models.Ambiente{}, nil).Once()
		mockMeioRepo.On("BuscarPorID", plantaDto.MeioCultivoID).Return(&models.MeioCultivo{}, nil).Once()
		mockPlantaRepo.On("Criar", mock.AnythingOfType("*models.Planta")).Run(func(args mock.Arguments) {
			planta := args.Get(0).(*models.Planta)
			planta.ID = 1 // Simula o ID sendo atribuído pelo banco de dados
			planta.DataPlantio = &now
			planta.Notas = &notas
		}).Return(nil).Once()

		response, err := servico.Criar(plantaDto)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, plantaDto.Nome, response.Nome)
		assert.Equal(t, uint(1), response.ID) // Verifica se o ID foi atribuído
		assert.Equal(t, plantaDto.DataPlantio.Format("2006-01-02"), response.DataPlantio.Format("2006-01-02"))
		assert.Equal(t, plantaDto.Notas, response.Notas)
		mockPlantaRepo.AssertExpectations(t)
		mockGeneticaRepo.AssertExpectations(t)
		mockAmbienteRepo.AssertExpectations(t)
		mockMeioRepo.AssertExpectations(t)
	})

	t.Run("Error - Empty Name", func(t *testing.T) {
		plantaDto := &dto.CreatePlantaDTO{
			Nome: "",
		}

		response, err := servico.Criar(plantaDto)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.EqualError(t, err, "o nome da planta não pode estar vazio")
		mockPlantaRepo.AssertNotCalled(t, "ExistePorNome")
	})

	t.Run("Error - Duplicate Name", func(t *testing.T) {
		plantaDto := &dto.CreatePlantaDTO{
			Nome: "Planta Existente",
		}

		mockPlantaRepo.On("ExistePorNome", plantaDto.Nome).Return(true).Once()

		response, err := servico.Criar(plantaDto)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.EqualError(t, err, "uma planta com este nome já existe")
		mockPlantaRepo.AssertExpectations(t)
	})

	t.Run("Error - Genetica Not Found", func(t *testing.T) {
		plantaDto := &dto.CreatePlantaDTO{
			Nome:          "Nova Planta",
			GeneticaID:    999,
			AmbienteID:    1,
			MeioCultivoID: 1,
		}

		mockPlantaRepo.On("ExistePorNome", plantaDto.Nome).Return(false).Once()
		mockGeneticaRepo.On("BuscarPorID", plantaDto.GeneticaID).Return((*models.Genetica)(nil), gorm.ErrRecordNotFound).Once()

		response, err := servico.Criar(plantaDto)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.EqualError(t, err, "genética não encontrada")
		mockPlantaRepo.AssertExpectations(t)
		mockGeneticaRepo.AssertExpectations(t)
	})

	t.Run("Error - Ambiente Not Found", func(t *testing.T) {
		plantaDto := &dto.CreatePlantaDTO{
			Nome:          "Nova Planta",
			GeneticaID:    1,
			AmbienteID:    999,
			MeioCultivoID: 1,
		}

		mockPlantaRepo.On("ExistePorNome", plantaDto.Nome).Return(false).Once()
		mockGeneticaRepo.On("BuscarPorID", plantaDto.GeneticaID).Return(&models.Genetica{}, nil).Once()
		mockAmbienteRepo.On("BuscarPorID", plantaDto.AmbienteID).Return((*models.Ambiente)(nil), gorm.ErrRecordNotFound).Once()

		response, err := servico.Criar(plantaDto)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.EqualError(t, err, "ambiente não encontrado")
		mockPlantaRepo.AssertExpectations(t)
		mockGeneticaRepo.AssertExpectations(t)
		mockAmbienteRepo.AssertExpectations(t)
	})

	t.Run("Error - MeioCultivo Not Found", func(t *testing.T) {
		plantaDto := &dto.CreatePlantaDTO{
			Nome:          "Nova Planta",
			GeneticaID:    1,
			AmbienteID:    1,
			MeioCultivoID: 999,
		}

		mockPlantaRepo.On("ExistePorNome", plantaDto.Nome).Return(false).Once()
		mockGeneticaRepo.On("BuscarPorID", plantaDto.GeneticaID).Return(&models.Genetica{}, nil).Once()
		mockAmbienteRepo.On("BuscarPorID", plantaDto.AmbienteID).Return(&models.Ambiente{}, nil).Once()
		mockMeioRepo.On("BuscarPorID", plantaDto.MeioCultivoID).Return((*models.MeioCultivo)(nil), gorm.ErrRecordNotFound).Once()

		response, err := servico.Criar(plantaDto)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.EqualError(t, err, "meio de cultivo não encontrado")
		mockPlantaRepo.AssertExpectations(t)
		mockGeneticaRepo.AssertExpectations(t)
		mockAmbienteRepo.AssertExpectations(t)
		mockMeioRepo.AssertExpectations(t)
	})

	t.Run("Error - Repository Create Error", func(t *testing.T) {
		now := time.Now()
		notas := "Algumas notas."
		plantaDto := &dto.CreatePlantaDTO{
			Nome:          "Nova Planta",
			GeneticaID:    1,
			AmbienteID:    1,
			MeioCultivoID: 1,
			Especie:       "Sativa",
			Status:        "ativa",
			DataPlantio:   now,
			Notas:         notas,
			ComecandoDe:   "semente",
			UsuarioID:     1,
		}

		mockPlantaRepo.On("ExistePorNome", plantaDto.Nome).Return(false).Once()
		mockGeneticaRepo.On("BuscarPorID", plantaDto.GeneticaID).Return(&models.Genetica{}, nil).Once()
		mockAmbienteRepo.On("BuscarPorID", plantaDto.AmbienteID).Return(&models.Ambiente{}, nil).Once()
		mockMeioRepo.On("BuscarPorID", plantaDto.MeioCultivoID).Return(&models.MeioCultivo{}, nil).Once()
		mockPlantaRepo.On("Criar", mock.AnythingOfType("*models.Planta")).Run(func(args mock.Arguments) {
			planta := args.Get(0).(*models.Planta)
			planta.DataPlantio = &now
			planta.Notas = &notas
		}).Return(errors.New("erro no repositório")).Once()

		response, err := servico.Criar(plantaDto)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.EqualError(t, err, "erro no repositório")
		mockPlantaRepo.AssertExpectations(t)
		mockGeneticaRepo.AssertExpectations(t)
		mockAmbienteRepo.AssertExpectations(t)
		mockMeioRepo.AssertExpectations(t)
	})
}

func TestPlantaService_BuscarPorID(t *testing.T) {
	mockPlantaRepo := new(MockPlantaRepositorio)
	mockGeneticaRepo := new(MockGeneticaRepositorio)
	mockAmbienteRepo := new(MockAmbienteRepositorio)
	mockMeioRepo := new(MockMeioCultivoRepositorio)

	servico := service.NewPlantaService(mockPlantaRepo, mockGeneticaRepo, mockAmbienteRepo, mockMeioRepo, mockPlantaRepo)

	t.Run("Success", func(t *testing.T) {
		plantID := uint(1)
		now := time.Now()
		notas := "Notas de teste"
		// Removed ID from struct literal as it's part of gorm.Model and might cause issues
		expectedPlanta := &models.Planta{Nome: "Planta Teste", DataPlantio: &now, Notas: &notas}
		expectedPlanta.ID = plantID // Set ID after initialization

		expectedResponse := &dto.PlantaResponseDTO{ID: plantID, Nome: "Planta Teste", DataPlantio: now, Notas: notas}

		mockPlantaRepo.On("BuscarPorID", plantID).Return(expectedPlanta, nil).Once()

		response, err := servico.BuscarPorID(plantID)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, expectedResponse.ID, response.ID)
		assert.Equal(t, expectedResponse.Nome, response.Nome)
		assert.Equal(t, expectedResponse.DataPlantio.Format("2006-01-02"), response.DataPlantio.Format("2006-01-02"))
		assert.Equal(t, expectedResponse.Notas, response.Notas)
		mockPlantaRepo.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		plantID := uint(999)

		mockPlantaRepo.On("BuscarPorID", plantID).Return((*models.Planta)(nil), gorm.ErrRecordNotFound).Once()

		response, err := servico.BuscarPorID(plantID)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		mockPlantaRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		plantID := uint(1)
		expectedError := errors.New("erro no repositório")

		mockPlantaRepo.On("BuscarPorID", plantID).Return((*models.Planta)(nil), expectedError).Once()

		response, err := servico.BuscarPorID(plantID)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), expectedError.Error())
		mockPlantaRepo.AssertExpectations(t)
	})
}