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

func (m *MockAmbienteService) ListarTodos(page, limit int) (*dto.PaginatedResponse, error) {
	args := m.Called(page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PaginatedResponse), args.Error(1)
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
		paginatedResponse := &dto.PaginatedResponse{
			Data:  expectedAmbientes,
			Total: int64(len(expectedAmbientes)),
			Page:  1,
			Limit: 10,
		}

		mockService.On("ListarTodos", mock.Anything, mock.Anything).Return(paginatedResponse, nil)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/ambientes?page=1&limit=10", nil)
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

		// Convert actualResponse.Data to []models.Ambiente for comparison
		actualAmbientesBytes, _ := json.Marshal(actualResponse.Data)
		var actualAmbientes []models.Ambiente
		json.Unmarshal(actualAmbientesBytes, &actualAmbientes)

		assert.Equal(t, expectedAmbientes, actualAmbientes)

		mockService.AssertExpectations(t)
	})

	t.Run("Success - Nenhum Ambiente Encontrado", func(t *testing.T) {
		// Preparação
		mockService := new(MockAmbienteService)
		controller := NewAmbienteController(mockService)

		paginatedResponse := &dto.PaginatedResponse{
			Data:  []models.Ambiente{},
			Total: 0,
			Page:  1,
			Limit: 10,
		}

		mockService.On("ListarTodos", mock.Anything, mock.Anything).Return(paginatedResponse, nil)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/ambientes?page=1&limit=10", nil)
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
		assert.Empty(t, actualResponse.Data)

		mockService.AssertExpectations(t)
	})

	t.Run("Error - Erro Interno do Serviço", func(t *testing.T) {
		// Preparação
		mockService := new(MockAmbienteService)
		controller := NewAmbienteController(mockService)

		mockService.On("ListarTodos", mock.Anything, mock.Anything).Return((*dto.PaginatedResponse)(nil), errors.New("erro no serviço"))

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/ambientes?page=1&limit=10", nil)
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
