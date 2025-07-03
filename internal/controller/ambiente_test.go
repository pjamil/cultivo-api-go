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

// MockAmbienteService é um mock para o service.AmbienteService
type MockAmbienteService struct {
	mock.Mock
}

func (m *MockAmbienteService) Criar(ambienteDto *dto.CreateAmbienteDTO) (*models.Ambiente, error) {
	args := m.Called(ambienteDto)
	return args.Get(0).(*models.Ambiente), args.Error(1)
}

func (m *MockAmbienteService) ListarTodos() ([]models.Ambiente, error) {
	args := m.Called()
	return args.Get(0).([]models.Ambiente), args.Error(1)
}

func (m *MockAmbienteService) BuscarPorID(id uint) (*models.Ambiente, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Ambiente), args.Error(1)
}

func (m *MockAmbienteService) Atualizar(id uint, ambienteDto *dto.UpdateAmbienteDTO) (*models.Ambiente, error) {
	args := m.Called(id, ambienteDto)
	return args.Get(0).(*models.Ambiente), args.Error(1)
}

func (m *MockAmbienteService) Deletar(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestAmbienteController_Listar(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success - Ambientes Encontrados", func(t *testing.T) {
		// Preparação
		mockService := new(MockAmbienteService)
		controller := NewAmbienteController(mockService)

		expectedAmbientes := []models.Ambiente{
			{Nome: "Estufa", Descricao: "Estufa de cultivo"},
			{Nome: "Quarto", Descricao: "Quarto de cultivo"},
		}

		mockService.On("ListarTodos").Return(expectedAmbientes, nil)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/ambientes", nil)
		c.Request = req

		controller.Listar(c)

		// Verificação
		assert.Equal(t, http.StatusOK, w.Code)

		var actualAmbientes []models.Ambiente
		err := json.Unmarshal(w.Body.Bytes(), &actualAmbientes)
		assert.NoError(t, err)
		assert.Equal(t, expectedAmbientes, actualAmbientes)

		mockService.AssertExpectations(t)
	})

	t.Run("Success - Nenhum Ambiente Encontrado", func(t *testing.T) {
		// Preparação
		mockService := new(MockAmbienteService)
		controller := NewAmbienteController(mockService)

		mockService.On("ListarTodos").Return([]models.Ambiente{}, nil)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/ambientes", nil)
		c.Request = req

		controller.Listar(c)

		// Verificação
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Nenhum ambiente encontrado", response["message"])

		mockService.AssertExpectations(t)
	})

	t.Run("Error - Erro Interno do Serviço", func(t *testing.T) {
		// Preparação
		mockService := new(MockAmbienteService)
		controller := NewAmbienteController(mockService)

		mockService.On("ListarTodos").Return([]models.Ambiente{}, errors.New("erro no serviço"))

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/ambientes", nil)
		c.Request = req

		controller.Listar(c)

		// Verificação
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "erro no serviço", response["error"])

		mockService.AssertExpectations(t)
	})
}

func TestAmbienteController_Criar(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		// Preparação
		mockService := new(MockAmbienteService)
		controller := NewAmbienteController(mockService)

		createDTO := &dto.CreateAmbienteDTO{
			Nome:        "Estufa Teste",
			Descricao:   "Estufa para testes",
			Tipo:        "interno",
			Comprimento: 10.0,
			Altura:      5.0,
			Largura:     8.0,
			TempoExposicao: 12,
		}
		expectedAmbiente := &models.Ambiente{
			Nome:        "Estufa Teste",
			Descricao:   "Estufa para testes",
			Tipo:        "interno",
			Comprimento: 10.0,
			Altura:      5.0,
			Largura:     8.0,
			TempoExposicao: 12,
			Orientacao:  "norte",
		}

		mockService.On("Criar", mock.AnythingOfType("*dto.CreateAmbienteDTO")).Return(expectedAmbiente, nil)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBody, _ := json.Marshal(createDTO)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/ambientes", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req

		controller.Criar(c)

		// Verificação
		assert.Equal(t, http.StatusCreated, w.Code)

		var actualAmbiente models.Ambiente
		err := json.Unmarshal(w.Body.Bytes(), &actualAmbiente)
		assert.NoError(t, err)
		assert.Equal(t, expectedAmbiente.Nome, actualAmbiente.Nome)
		assert.Equal(t, expectedAmbiente.Tipo, actualAmbiente.Tipo)

		mockService.AssertExpectations(t)
	})

	t.Run("Invalid Payload", func(t *testing.T) {
		// Preparação
		mockService := new(MockAmbienteService)
		controller := NewAmbienteController(mockService)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/ambientes", bytes.NewBufferString("invalid-json"))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req

		controller.Criar(c)

		// Verificação
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		// Preparação
		mockService := new(MockAmbienteService)
		controller := NewAmbienteController(mockService)

		createDTO := &dto.CreateAmbienteDTO{
			Nome:        "Estufa Teste",
			Descricao:   "Estufa para testes",
			Tipo:        "interno",
			Comprimento: 10.0,
			Altura:      5.0,
			Largura:     8.0,
			TempoExposicao: 12,
		}

		mockService.On("Criar", mock.AnythingOfType("*dto.CreateAmbienteDTO")).Return((*models.Ambiente)(nil), errors.New("erro interno do serviço"))

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBody, _ := json.Marshal(createDTO)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/ambientes", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req

		controller.Criar(c)

		// Verificação
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		mockService.AssertExpectations(t)
	})
}

func TestAmbienteController_BuscarPorID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		// Preparação
		mockService := new(MockAmbienteService)
		controller := NewAmbienteController(mockService)

		expectedAmbiente := &models.Ambiente{Nome: "Estufa Teste"}

		mockService.On("BuscarPorID", uint(1)).Return(expectedAmbiente, nil)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		controller.BuscarPorID(c)

		// Verificação
		assert.Equal(t, http.StatusOK, w.Code)

		var actualAmbiente models.Ambiente
		err := json.Unmarshal(w.Body.Bytes(), &actualAmbiente)
		assert.NoError(t, err)
		assert.Equal(t, expectedAmbiente.Nome, actualAmbiente.Nome)

		mockService.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		// Preparação
		mockService := new(MockAmbienteService)
		controller := NewAmbienteController(mockService)

		mockService.On("BuscarPorID", uint(1)).Return((*models.Ambiente)(nil), gorm.ErrRecordNotFound)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		controller.BuscarPorID(c)

		// Verificação
		assert.Equal(t, http.StatusNotFound, w.Code)

		mockService.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		// Preparação
		mockService := new(MockAmbienteService)
		controller := NewAmbienteController(mockService)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "abc"}}

		controller.BuscarPorID(c)

		// Verificação
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestAmbienteController_Atualizar(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		// Preparação
		mockService := new(MockAmbienteService)
		controller := NewAmbienteController(mockService)

		updateDTO := &dto.UpdateAmbienteDTO{Nome: "Estufa Atualizada"}
				expectedAmbiente := &models.Ambiente{Nome: "Estufa Atualizada"}

		mockService.On("Atualizar", uint(1), mock.AnythingOfType("*dto.UpdateAmbienteDTO")).Return(expectedAmbiente, nil)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBody, _ := json.Marshal(updateDTO)
		req, _ := http.NewRequest(http.MethodPut, "/api/v1/ambientes/1", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		controller.Atualizar(c)

		// Verificação
		assert.Equal(t, http.StatusOK, w.Code)

		var actualAmbiente models.Ambiente
		err := json.Unmarshal(w.Body.Bytes(), &actualAmbiente)
		assert.NoError(t, err)
		assert.Equal(t, expectedAmbiente.Nome, actualAmbiente.Nome)

		mockService.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		// Preparação
		mockService := new(MockAmbienteService)
		controller := NewAmbienteController(mockService)

		updateDTO := &dto.UpdateAmbienteDTO{Nome: "Estufa Atualizada"}

		mockService.On("Atualizar", uint(1), mock.AnythingOfType("*dto.UpdateAmbienteDTO")).Return((*models.Ambiente)(nil), gorm.ErrRecordNotFound)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBody, _ := json.Marshal(updateDTO)
		req, _ := http.NewRequest(http.MethodPut, "/api/v1/ambientes/1", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		controller.Atualizar(c)

		// Verificação
		assert.Equal(t, http.StatusNotFound, w.Code)

		mockService.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		// Preparação
		mockService := new(MockAmbienteService)
		controller := NewAmbienteController(mockService)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "abc"}}

		controller.Atualizar(c)

		// Verificação
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Invalid Payload", func(t *testing.T) {
		// Preparação
		mockService := new(MockAmbienteService)
		controller := NewAmbienteController(mockService)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest(http.MethodPut, "/api/v1/ambientes/1", bytes.NewBufferString("invalid-json"))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		controller.Atualizar(c)

		// Verificação
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestAmbienteController_Deletar(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		// Preparação
		mockService := new(MockAmbienteService)
		controller := NewAmbienteController(mockService)

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
		mockService := new(MockAmbienteService)
		controller := NewAmbienteController(mockService)

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
		mockService := new(MockAmbienteService)
		controller := NewAmbienteController(mockService)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "abc"}}

		controller.Deletar(c)

		// Verificação
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
