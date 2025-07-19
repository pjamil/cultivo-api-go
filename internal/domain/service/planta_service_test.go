package service_test

import (
	"errors"
	"testing"
	"time"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/entity"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/service"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/service/test"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestPlantaService_Criar(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockPlantaRepo := new(test.MockPlantaRepositorio)
		mockGeneticaRepo := new(test.MockGeneticaRepositorio)
		mockAmbienteRepo := new(test.MockAmbienteRepositorio)
		mockMeioRepo := new(test.MockMeioCultivoRepositorio)
		mockRegistroDiarioRepo := new(MockRegistroDiarioRepositorio)

		servico := service.NewPlantaService(mockPlantaRepo, mockGeneticaRepo, mockAmbienteRepo, mockMeioRepo, mockRegistroDiarioRepo)

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
		mockGeneticaRepo.On("BuscarPorID", plantaDto.GeneticaID).Return(&entity.Genetica{}, nil).Once()
		mockAmbienteRepo.On("BuscarPorID", plantaDto.AmbienteID).Return(&entity.Ambiente{}, nil).Once()
		mockMeioRepo.On("BuscarPorID", plantaDto.MeioCultivoID).Return(&entity.MeioCultivo{}, nil).Once()
		mockPlantaRepo.On("Criar", mock.AnythingOfType("*entity.Planta")).Run(func(args mock.Arguments) {
			planta := args.Get(0).(*entity.Planta)
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
		assert.Equal(t, plantaDto.Notas, *response.Notas)
		mockPlantaRepo.AssertExpectations(t)
		mockGeneticaRepo.AssertExpectations(t)
		mockAmbienteRepo.AssertExpectations(t)
		mockMeioRepo.AssertExpectations(t)
	})

	t.Run("Error - Empty Name", func(t *testing.T) {
		mockPlantaRepo := new(test.MockPlantaRepositorio) // Keep this as test.MockPlantaRepositorio
		mockGeneticaRepo := new(test.MockGeneticaRepositorio)
		mockAmbienteRepo := new(test.MockAmbienteRepositorio)
		mockMeioRepo := new(test.MockMeioCultivoRepositorio)
		mockRegistroDiarioRepo := new(MockRegistroDiarioRepositorio)

		servico := service.NewPlantaService(mockPlantaRepo, mockGeneticaRepo, mockAmbienteRepo, mockMeioRepo, mockRegistroDiarioRepo)

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
		mockPlantaRepo := new(test.MockPlantaRepositorio) // Keep this as test.MockPlantaRepositorio
		mockGeneticaRepo := new(test.MockGeneticaRepositorio)
		mockAmbienteRepo := new(test.MockAmbienteRepositorio)
		mockMeioRepo := new(test.MockMeioCultivoRepositorio)
		mockRegistroDiarioRepo := new(MockRegistroDiarioRepositorio)

		servico := service.NewPlantaService(mockPlantaRepo, mockGeneticaRepo, mockAmbienteRepo, mockMeioRepo, mockRegistroDiarioRepo)

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
		mockPlantaRepo := new(test.MockPlantaRepositorio) // Keep this as test.MockPlantaRepositorio
		mockGeneticaRepo := new(test.MockGeneticaRepositorio)
		mockAmbienteRepo := new(test.MockAmbienteRepositorio)
		mockMeioRepo := new(test.MockMeioCultivoRepositorio)
		mockRegistroDiarioRepo := new(MockRegistroDiarioRepositorio)

		servico := service.NewPlantaService(mockPlantaRepo, mockGeneticaRepo, mockAmbienteRepo, mockMeioRepo, mockRegistroDiarioRepo)

		plantaDto := &dto.CreatePlantaDTO{
			Nome:          "Nova Planta",
			GeneticaID:    999,
			AmbienteID:    1,
			MeioCultivoID: 1,
		}

		mockPlantaRepo.On("ExistePorNome", plantaDto.Nome).Return(false).Once()
		mockGeneticaRepo.On("BuscarPorID", plantaDto.GeneticaID).Return((*entity.Genetica)(nil), gorm.ErrRecordNotFound).Once()

		response, err := servico.Criar(plantaDto)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.EqualError(t, err, "genética não encontrada")
		mockPlantaRepo.AssertExpectations(t)
		mockGeneticaRepo.AssertExpectations(t)
	})

	t.Run("Error - Ambiente Not Found", func(t *testing.T) {
		mockPlantaRepo := new(test.MockPlantaRepositorio) // Keep this as test.MockPlantaRepositorio
		mockGeneticaRepo := new(test.MockGeneticaRepositorio)
		mockAmbienteRepo := new(test.MockAmbienteRepositorio)
		mockMeioRepo := new(test.MockMeioCultivoRepositorio)
		mockRegistroDiarioRepo := new(MockRegistroDiarioRepositorio)

		servico := service.NewPlantaService(mockPlantaRepo, mockGeneticaRepo, mockAmbienteRepo, mockMeioRepo, mockRegistroDiarioRepo)

		plantaDto := &dto.CreatePlantaDTO{
			Nome:          "Nova Planta",
			GeneticaID:    1,
			AmbienteID:    999,
			MeioCultivoID: 1,
		}

		mockPlantaRepo.On("ExistePorNome", plantaDto.Nome).Return(false).Once()
		mockGeneticaRepo.On("BuscarPorID", plantaDto.GeneticaID).Return(&entity.Genetica{}, nil).Once()
		mockAmbienteRepo.On("BuscarPorID", plantaDto.AmbienteID).Return((*entity.Ambiente)(nil), gorm.ErrRecordNotFound).Once()

		response, err := servico.Criar(plantaDto)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.EqualError(t, err, "ambiente não encontrado")
		mockPlantaRepo.AssertExpectations(t)
		mockGeneticaRepo.AssertExpectations(t)
		mockAmbienteRepo.AssertExpectations(t)
	})

	t.Run("Error - MeioCultivo Not Found", func(t *testing.T) {
		mockPlantaRepo := new(test.MockPlantaRepositorio) // Keep this as test.MockPlantaRepositorio
		mockGeneticaRepo := new(test.MockGeneticaRepositorio)
		mockAmbienteRepo := new(test.MockAmbienteRepositorio)
		mockMeioRepo := new(test.MockMeioCultivoRepositorio)
		mockRegistroDiarioRepo := new(MockRegistroDiarioRepositorio)

		servico := service.NewPlantaService(mockPlantaRepo, mockGeneticaRepo, mockAmbienteRepo, mockMeioRepo, mockRegistroDiarioRepo)

		plantaDto := &dto.CreatePlantaDTO{
			Nome:          "Nova Planta",
			GeneticaID:    1,
			AmbienteID:    1,
			MeioCultivoID: 999,
		}

		mockPlantaRepo.On("ExistePorNome", plantaDto.Nome).Return(false).Once()
		mockGeneticaRepo.On("BuscarPorID", plantaDto.GeneticaID).Return(&entity.Genetica{}, nil).Once()
		mockAmbienteRepo.On("BuscarPorID", plantaDto.AmbienteID).Return(&entity.Ambiente{}, nil).Once()
		mockMeioRepo.On("BuscarPorID", plantaDto.MeioCultivoID).Return((*entity.MeioCultivo)(nil), gorm.ErrRecordNotFound).Once()

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
		mockPlantaRepo := new(test.MockPlantaRepositorio) // Keep this as test.MockPlantaRepositorio
		mockGeneticaRepo := new(test.MockGeneticaRepositorio)
		mockAmbienteRepo := new(test.MockAmbienteRepositorio)
		mockMeioRepo := new(test.MockMeioCultivoRepositorio)
		mockRegistroDiarioRepo := new(MockRegistroDiarioRepositorio)

		servico := service.NewPlantaService(mockPlantaRepo, mockGeneticaRepo, mockAmbienteRepo, mockMeioRepo, mockRegistroDiarioRepo)

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
		mockGeneticaRepo.On("BuscarPorID", plantaDto.GeneticaID).Return(&entity.Genetica{}, nil).Once()
		mockAmbienteRepo.On("BuscarPorID", plantaDto.AmbienteID).Return(&entity.Ambiente{}, nil).Once()
		mockMeioRepo.On("BuscarPorID", plantaDto.MeioCultivoID).Return(&entity.MeioCultivo{}, nil).Once()
		mockPlantaRepo.On("Criar", mock.AnythingOfType("*entity.Planta")).Run(func(args mock.Arguments) {
			planta := args.Get(0).(*entity.Planta)
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

func TestPlantaService_ListarTodas(t *testing.T) {
	t.Run("Success - Plantas Encontradas", func(t *testing.T) {
		mockPlantaRepo := new(test.MockPlantaRepositorio)
		mockGeneticaRepo := new(test.MockGeneticaRepositorio)
		mockAmbienteRepo := new(test.MockAmbienteRepositorio) // Keep this as test.MockAmbienteRepositorio
		mockMeioRepo := new(test.MockMeioCultivoRepositorio)
		mockRegistroDiarioRepo := new(MockRegistroDiarioRepositorio)
		servico := service.NewPlantaService(mockPlantaRepo, mockGeneticaRepo, mockAmbienteRepo, mockMeioRepo, mockRegistroDiarioRepo)

		now := time.Now()
		notas := "Algumas notas."
		mockPlantas := []entity.Planta{
			{Model: gorm.Model{ID: 1}, Nome: "Planta 1", Especie: "Sativa", DataPlantio: &now, Notas: &notas},
			{Model: gorm.Model{ID: 2}, Nome: "Planta 2", Especie: "Indica", DataPlantio: &now, Notas: &notas},
		}
		expectedTotal := int64(len(mockPlantas))
		page := 1
		limit := 10

		expectedResponseData := []dto.PlantaResponseDTO{
			{ID: 1, Nome: "Planta 1", Especie: "Sativa", DataPlantio: utils.TimePtr(now), Notas: utils.StringPtr(notas)},
			{ID: 2, Nome: "Planta 2", Especie: "Indica", DataPlantio: utils.TimePtr(now), Notas: utils.StringPtr(notas)},
		}

		mockPlantaRepo.On("ListarTodos", page, limit).Return(mockPlantas, expectedTotal, nil).Once()

		// Act
		responseDTOs, total, err := servico.ListarTodas(page, limit)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, responseDTOs)
		assert.Equal(t, expectedTotal, total)
		assert.Equal(t, expectedResponseData, responseDTOs)
		mockPlantaRepo.AssertExpectations(t)
	})

	t.Run("Success - Nenhuma Planta Encontrada", func(t *testing.T) {
		mockPlantaRepo := new(test.MockPlantaRepositorio) // Keep this as test.MockPlantaRepositorio
		mockGeneticaRepo := new(test.MockGeneticaRepositorio)
		mockAmbienteRepo := new(test.MockAmbienteRepositorio) // Keep this as test.MockAmbienteRepositorio
		mockMeioRepo := new(test.MockMeioCultivoRepositorio)
		mockRegistroDiarioRepo := new(MockRegistroDiarioRepositorio)
		servico := service.NewPlantaService(mockPlantaRepo, mockGeneticaRepo, mockAmbienteRepo, mockMeioRepo, mockRegistroDiarioRepo)

		page := 1
		limit := 10
		mockPlantaRepo.On("ListarTodos", page, limit).Return([]entity.Planta{}, int64(0), nil).Once()

		responseDTOs, total, err := servico.ListarTodas(page, limit)

		assert.NoError(t, err)
		assert.NotNil(t, responseDTOs)
		assert.Equal(t, int64(0), total)
		assert.Empty(t, responseDTOs)
		mockPlantaRepo.AssertExpectations(t)
	})

	t.Run("Error - Repository Error", func(t *testing.T) {
		mockPlantaRepo := new(test.MockPlantaRepositorio) // Keep this as test.MockPlantaRepositorio
		mockGeneticaRepo := new(test.MockGeneticaRepositorio)
		mockAmbienteRepo := new(test.MockAmbienteRepositorio) // Keep this as test.MockAmbienteRepositorio
		mockMeioRepo := new(test.MockMeioCultivoRepositorio)
		mockRegistroDiarioRepo := new(MockRegistroDiarioRepositorio)
		servico := service.NewPlantaService(mockPlantaRepo, mockGeneticaRepo, mockAmbienteRepo, mockMeioRepo, mockRegistroDiarioRepo)

		page := 1
		limit := 10
		mockPlantaRepo.On("ListarTodos", page, limit).Return([]entity.Planta{}, int64(0), errors.New("erro no repositório")).Once()

		responseDTOs, total, err := servico.ListarTodas(page, limit)

		assert.Error(t, err)
		assert.Nil(t, responseDTOs)
		assert.Equal(t, int64(0), total)
		assert.EqualError(t, err, "erro no repositório")
		mockPlantaRepo.AssertExpectations(t)
	})
}

func TestPlantaService_BuscarPorID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockPlantaRepo := new(test.MockPlantaRepositorio) // Keep this as test.MockPlantaRepositorio
		mockGeneticaRepo := new(test.MockGeneticaRepositorio)
		mockAmbienteRepo := new(test.MockAmbienteRepositorio) // Keep this as test.MockAmbienteRepositorio
		mockMeioRepo := new(test.MockMeioCultivoRepositorio)
		mockRegistroDiarioRepo := new(MockRegistroDiarioRepositorio)
		servico := service.NewPlantaService(mockPlantaRepo, mockGeneticaRepo, mockAmbienteRepo, mockMeioRepo, mockRegistroDiarioRepo)

		plantID := uint(1)
		now := time.Now()
		notas := "Notas de teste"
		expectedPlanta := &entity.Planta{Nome: "Planta Teste", DataPlantio: &now, Notas: &notas}
		expectedPlanta.ID = plantID

		expectedResponse := &dto.PlantaResponseDTO{ID: plantID, Nome: "Planta Teste", DataPlantio: utils.TimePtr(now), Notas: utils.StringPtr(notas)}

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
		mockPlantaRepo := new(test.MockPlantaRepositorio) // Keep this as test.MockPlantaRepositorio
		mockGeneticaRepo := new(test.MockGeneticaRepositorio)
		mockAmbienteRepo := new(test.MockAmbienteRepositorio) // Keep this as test.MockAmbienteRepositorio
		mockMeioRepo := new(test.MockMeioCultivoRepositorio)
		mockRegistroDiarioRepo := new(MockRegistroDiarioRepositorio)
		servico := service.NewPlantaService(mockPlantaRepo, mockGeneticaRepo, mockAmbienteRepo, mockMeioRepo, mockRegistroDiarioRepo)

		plantID := uint(999)

		mockPlantaRepo.On("BuscarPorID", plantID).Return((*entity.Planta)(nil), gorm.ErrRecordNotFound).Once()

		response, err := servico.BuscarPorID(plantID)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, utils.ErrNotFound, err)
		mockPlantaRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		mockPlantaRepo := new(test.MockPlantaRepositorio)     // Keep this as test.MockPlantaRepositorio
		mockGeneticaRepo := new(test.MockGeneticaRepositorio) // Keep this as test.MockGeneticaRepositorio
		mockAmbienteRepo := new(test.MockAmbienteRepositorio)
		mockMeioRepo := new(test.MockMeioCultivoRepositorio)
		mockRegistroDiarioRepo := new(MockRegistroDiarioRepositorio)
		servico := service.NewPlantaService(mockPlantaRepo, mockGeneticaRepo, mockAmbienteRepo, mockMeioRepo, mockRegistroDiarioRepo)

		plantID := uint(1)
		expectedError := errors.New("erro no repositório")

		mockPlantaRepo.On("BuscarPorID", plantID).Return((*entity.Planta)(nil), expectedError).Once()

		response, err := servico.BuscarPorID(plantID)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), expectedError.Error())
		mockPlantaRepo.AssertExpectations(t)
	})
}
