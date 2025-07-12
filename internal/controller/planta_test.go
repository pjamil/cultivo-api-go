package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockPlantaService é um mock para o service.PlantaService
type MockPlantaService struct {
	mock.Mock
}

func (m *MockPlantaService) Criar(createDto *dto.CreatePlantaDTO) (*dto.PlantaResponseDTO, error) {
	args := m.Called(createDto)
	return args.Get(0).(*dto.PlantaResponseDTO), args.Error(1)
}

func (m *MockPlantaService) ListarTodas(page, limit int) (*dto.PaginatedResponse, error) {
	args := m.Called(page, limit)
	return args.Get(0).(*dto.PaginatedResponse), args.Error(1)
}

func (m *MockPlantaService) BuscarPorID(id uint) (*dto.PlantaResponseDTO, error) {
	args := m.Called(id)
	return args.Get(0).(*dto.PlantaResponseDTO), args.Error(1)
}

func (m *MockPlantaService) Atualizar(id uint, updateDto *dto.UpdatePlantaDTO) (*dto.PlantaResponseDTO, error) {
	args := m.Called(id, updateDto)
	return args.Get(0).(*dto.PlantaResponseDTO), args.Error(1)
}

func (m *MockPlantaService) Deletar(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPlantaService) RegistrarFato(id uint, tipo models.RegistroTipo, titulo string, conteudo string) error {
	args := m.Called(id, tipo, titulo, conteudo)
	return args.Error(0)
}

func (m *MockPlantaService) BuscarPorEspecie(especie models.Especie) ([]models.Planta, error) {
	args := m.Called(especie)
	return args.Get(0).([]models.Planta), args.Error(1)
}

func (m *MockPlantaService) BuscarPorStatus(status string) ([]models.Planta, error) {
	args := m.Called(status)
	return args.Get(0).([]models.Planta), args.Error(1)
}

func TestPlantaController_Listar(t *testing.T) {
	// Configuração do Gin em modo de teste
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		// Preparação
		mockService := new(MockPlantaService)
		controller := NewPlantaController(mockService)

		expectedPlantas := []models.Planta{
			{Nome: "Sativa"},
			{Nome: "Indica"},
		}
		paginatedResponse := &dto.PaginatedResponse{
			Data:  expectedPlantas,
			Total: int64(len(expectedPlantas)),
			Page:  1,
			Limit: 10,
		}

		mockService.On("ListarTodas", mock.Anything, mock.Anything).Return(paginatedResponse, nil)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/plantas?page=1&limit=10", nil)
		c.Request = req

		controller.Listar(c)

		// Verificação
		assert.Equal(t, http.StatusOK, w.Code)

		var actualResponse dto.PaginatedResponse
		err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
		assert.NoError(t, err)
		assert.Equal(t, paginatedResponse.Total, actualResponse.Total)
		assert.Equal(t, paginatedResponse.Page, actualResponse.Page)
		assert.Equal(t, paginatedResponse.Limit, actualResponse.Limit)

		// Convert actualResponse.Data to []models.Planta for comparison
		actualPlantasBytes, _ := json.Marshal(actualResponse.Data)
		var actualPlantas []models.Planta
		json.Unmarshal(actualPlantasBytes, &actualPlantas)

		assert.Equal(t, expectedPlantas, actualPlantas)

		mockService.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		// Preparação
		mockService := new(MockPlantaService)
		controller := NewPlantaController(mockService)

		mockService.On("ListarTodas", mock.Anything, mock.Anything).Return((*dto.PaginatedResponse)(nil), errors.New("erro no serviço"))

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/plantas?page=1&limit=10", nil)
		c.Request = req

		controller.Listar(c)

		// Verificação
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Erro ao listar plantas", response["error"])

		mockService.AssertExpectations(t)
	})
}

func TestPlantaController_BuscarPorID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		// Preparação
		mockService := new(MockPlantaService)
		controller := NewPlantaController(mockService)

		expectedPlanta := &dto.PlantaResponseDTO{ID: 1, Nome: "Sativa"}

		mockService.On("BuscarPorID", uint(1)).Return(expectedPlanta, nil)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		controller.BuscarPorID(c)

		// Verificação
		assert.Equal(t, http.StatusOK, w.Code)

		var actualPlanta dto.PlantaResponseDTO
		err := json.Unmarshal(w.Body.Bytes(), &actualPlanta)
		assert.NoError(t, err)
		assert.Equal(t, *expectedPlanta, actualPlanta)

		mockService.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		// Preparação
		mockService := new(MockPlantaService)
		controller := NewPlantaController(mockService)

		mockService.On("BuscarPorID", uint(1)).Return((*dto.PlantaResponseDTO)(nil), gorm.ErrRecordNotFound)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		controller.BuscarPorID(c)

		// Verificação
		assert.Equal(t, http.StatusNotFound, w.Code)

		mockService.AssertExpectations(t)
	})
}

func TestPlantaController_Criar(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		// Preparação
		mockService := new(MockPlantaService)
		controller := NewPlantaController(mockService)

		createDTO := &dto.CreatePlantaDTO{Nome: "Sativa", Especie: "Sativa"}
		expectedPlanta := &dto.PlantaResponseDTO{ID: 1, Nome: "Sativa"}

		mockService.On("Criar", mock.MatchedBy(func(arg *dto.CreatePlantaDTO) bool {
			return arg.Nome == createDTO.Nome && arg.Especie == createDTO.Especie
		})).Return(expectedPlanta, nil)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBody, _ := json.Marshal(createDTO)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/plantas", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req

		controller.Criar(c)

		// Verificação
		assert.Equal(t, http.StatusCreated, w.Code)

		var actualPlanta dto.PlantaResponseDTO
		err := json.Unmarshal(w.Body.Bytes(), &actualPlanta)
		assert.NoError(t, err)
		assert.Equal(t, *expectedPlanta, actualPlanta)

		mockService.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		// Preparação
		mockService := new(MockPlantaService)
		controller := NewPlantaController(mockService)

		createDTO := &dto.CreatePlantaDTO{Nome: "Sativa", Especie: "Sativa"}

		mockService.On("Criar", mock.MatchedBy(func(arg *dto.CreatePlantaDTO) bool {
			return arg.Nome == createDTO.Nome && arg.Especie == createDTO.Especie
		})).Return((*dto.PlantaResponseDTO)(nil), errors.New("erro no serviço"))

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBody, _ := json.Marshal(createDTO)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/plantas", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req

		controller.Criar(c)

		// Verificação
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		mockService.AssertExpectations(t)
	})
}

func TestPlantaController_Atualizar(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		// Preparação
		mockService := new(MockPlantaService)
		controller := NewPlantaController(mockService)

		updateDTO := &dto.UpdatePlantaDTO{Nome: "Sativa Nova"}
		expectedPlanta := &dto.PlantaResponseDTO{ID: 1, Nome: "Sativa Nova"}

		mockService.On("Atualizar", uint(1), mock.MatchedBy(func(arg *dto.UpdatePlantaDTO) bool {
			return arg.Nome == updateDTO.Nome
		})).Return(expectedPlanta, nil)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBody, _ := json.Marshal(updateDTO)
		req, _ := http.NewRequest(http.MethodPut, "/api/v1/plantas/1", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		controller.Atualizar(c)

		// Verificação
		assert.Equal(t, http.StatusOK, w.Code)

		var actualPlanta dto.PlantaResponseDTO
		err := json.Unmarshal(w.Body.Bytes(), &actualPlanta)
		assert.NoError(t, err)
		assert.Equal(t, *expectedPlanta, actualPlanta)

		mockService.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		// Preparação
		mockService := new(MockPlantaService)
		controller := NewPlantaController(mockService)

		updateDTO := &dto.UpdatePlantaDTO{Nome: "Sativa Nova"}

		mockService.On("Atualizar", uint(1), mock.MatchedBy(func(arg *dto.UpdatePlantaDTO) bool {
			return arg.Nome == updateDTO.Nome
		})).Return((*dto.PlantaResponseDTO)(nil), gorm.ErrRecordNotFound)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBody, _ := json.Marshal(updateDTO)
		req, _ := http.NewRequest(http.MethodPut, "/api/v1/plantas/1", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		controller.Atualizar(c)

		// Verificação
		assert.Equal(t, http.StatusNotFound, w.Code)

		mockService.AssertExpectations(t)
	})

	t.Run("Invalid Payload", func(t *testing.T) {
		// Preparação
		mockService := new(MockPlantaService)
		controller := NewPlantaController(mockService)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest(http.MethodPut, "/api/v1/plantas/1", bytes.NewBufferString("invalid-json"))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		controller.Atualizar(c)

		// Verificação
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestPlantaController_Deletar(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		// Preparação
		mockService := new(MockPlantaService)
		controller := NewPlantaController(mockService)

		mockService.On("Deletar", uint(1)).Return(nil)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		controller.Deletar(c)

		// Verificação
		assert.Equal(t, http.StatusOK, w.Code)

		mockService.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		// Preparação
		mockService := new(MockPlantaService)
		controller := NewPlantaController(mockService)

		mockService.On("Deletar", uint(1)).Return(gorm.ErrRecordNotFound)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		controller.Deletar(c)

		// Verificação
		assert.Equal(t, http.StatusNotFound, w.Code)

		mockService.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		// Preparação
		mockService := new(MockPlantaService)
		controller := NewPlantaController(mockService)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "abc"}}

		controller.Deletar(c)

		// Verificação
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestPlantaController_RegistrarFato(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		// Preparação
		mockService := new(MockPlantaService)
		controller := NewPlantaController(mockService)

		plantID := uint(1)
		fatoDTO := &dto.RegistrarFatoDTO{
			Tipo:    "observacao",
			Titulo:  "Primeira Observação",
			Conteudo: "A planta está crescendo bem.",
		}

		mockService.On("RegistrarFato", plantID, models.RegistroTipo(fatoDTO.Tipo), fatoDTO.Titulo, fatoDTO.Conteudo).Return(nil)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBody, _ := json.Marshal(fatoDTO)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/plantas/1/registrar-fato", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		controller.RegistrarFato(c)

		// Verificação
		assert.Equal(t, http.StatusOK, w.Code)

		mockService.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		// Preparação
		mockService := new(MockPlantaService)
		controller := NewPlantaController(mockService)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "abc"}}

		controller.RegistrarFato(c)

		// Verificação
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Invalid Payload", func(t *testing.T) {
		// Preparação
		mockService := new(MockPlantaService)
		controller := NewPlantaController(mockService)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/plantas/1/registrar-fato", bytes.NewBufferString("invalid-json"))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		controller.RegistrarFato(c)

		// Verificação
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Plant Not Found", func(t *testing.T) {
		// Preparação
		mockService := new(MockPlantaService)
		controller := NewPlantaController(mockService)

		plantID := uint(1)
		fatoDTO := &dto.RegistrarFatoDTO{
			Tipo:    "observacao",
			Titulo:  "Primeira Observação",
			Conteudo: "A planta está crescendo bem.",
		}

		mockService.On("RegistrarFato", plantID, models.RegistroTipo(fatoDTO.Tipo), fatoDTO.Titulo, fatoDTO.Conteudo).Return(gorm.ErrRecordNotFound)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBody, _ := json.Marshal(fatoDTO)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/plantas/1/registrar-fato", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		controller.RegistrarFato(c)

		// Verificação
		assert.Equal(t, http.StatusNotFound, w.Code)

		mockService.AssertExpectations(t)
	})

	t.Run("Service Error", func(t *testing.T) {
		// Preparação
		mockService := new(MockPlantaService)
		controller := NewPlantaController(mockService)

		plantID := uint(1)
		fatoDTO := &dto.RegistrarFatoDTO{
			Tipo:    "observacao",
			Titulo:  "Primeira Observação",
			Conteudo: "A planta está crescendo bem.",
		}

		mockService.On("RegistrarFato", plantID, models.RegistroTipo(fatoDTO.Tipo), fatoDTO.Titulo, fatoDTO.Conteudo).Return(errors.New("erro interno do serviço"))

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBody, _ := json.Marshal(fatoDTO)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/plantas/1/registrar-fato", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		controller.RegistrarFato(c)

		// Verificação
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		mockService.AssertExpectations(t)
	})
}