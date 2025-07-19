package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/entity"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockAmbienteService é um mock para o service.AmbienteService
type MockAmbienteService struct {
	mock.Mock
}

func (m *MockAmbienteService) Criar(ambienteDto *dto.CreateAmbienteDTO) (*entity.Ambiente, error) {
	args := m.Called(ambienteDto)
	return args.Get(0).(*entity.Ambiente), args.Error(1)
}

func (m *MockAmbienteService) ListarTodos(page, limit int) ([]dto.AmbienteResponseDTO, int64, error) {
	args := m.Called(page, limit)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]dto.AmbienteResponseDTO), args.Get(1).(int64), args.Error(2)
}

func (m *MockAmbienteService) BuscarPorID(id uint) (*entity.Ambiente, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.Ambiente), args.Error(1)
}

func (m *MockAmbienteService) Atualizar(id uint, ambienteDto *dto.UpdateAmbienteDTO) (*entity.Ambiente, error) {
	args := m.Called(id, ambienteDto)
	return args.Get(0).(*entity.Ambiente), args.Error(1)
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

		expectedAmbientes := []dto.AmbienteResponseDTO{
			{Nome: "Estufa", Descricao: "Estufa de cultivo", Tipo: "interno"},
			{Nome: "Quarto", Descricao: "Quarto de cultivo", Tipo: "interno"},
		}
		dataBytes, err := json.Marshal(expectedAmbientes)
		assert.NoError(t, err)
		paginatedResponse := &dto.PaginatedResponse{
			Data:  dataBytes,
			Total: int64(len(expectedAmbientes)), // This will be the second return value
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
		err = json.Unmarshal(w.Body.Bytes(), &actualResponse)
		assert.NoError(t, err)
		assert.Equal(t, paginatedResponse.Total, actualResponse.Total)
		assert.Equal(t, paginatedResponse.Page, actualResponse.Page)
		assert.Equal(t, paginatedResponse.Limit, actualResponse.Limit)
		// Convert actualResponse.Data to []dto.AmbienteResponseDTO for comparison
		actualAmbientesBytes, _ := json.Marshal(actualResponse.Data)
		var actualAmbientes []dto.AmbienteResponseDTO

		json.Unmarshal(actualAmbientesBytes, &actualAmbientes)

		assert.Equal(t, expectedAmbientes, actualAmbientes)

		mockService.AssertExpectations(t)
	})

	t.Run("Success - Nenhum Ambiente Encontrado", func(t *testing.T) {
		// Preparação
		mockService := new(MockAmbienteService)
		controller := NewAmbienteController(mockService)

		// The service returns []dto.AmbienteResponseDTO, total, error
		mockService.On("ListarTodos", mock.Anything, mock.Anything).Return(
			[]dto.AmbienteResponseDTO{}, // Empty slice for Data
			int64(0),                    // Total count
			nil,                         // No error
		)

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
		assert.Equal(t, int64(0), actualResponse.Total)
		assert.Equal(t, 1, actualResponse.Page)
		assert.Equal(t, 10, actualResponse.Limit)
		assert.Equal(t, json.RawMessage("[]"), actualResponse.Data) // Data should be an empty JSON array

		mockService.AssertExpectations(t)
	})

	t.Run("Error - Erro Interno do Serviço", func(t *testing.T) {
		// Preparação
		mockService := new(MockAmbienteService)
		controller := NewAmbienteController(mockService)

		mockService.On("ListarTodos", mock.Anything, mock.Anything).Return([]dto.AmbienteResponseDTO{}, int64(0), errors.New("erro no serviço"))

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/ambientes?page=1&limit=10", nil)
		c.Request = req

		controller.Listar(c)

		// Verificação
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Erro interno ao listar ambientes", response["message"])
		assert.Equal(t, "erro no serviço", response["details"])

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
			Nome:           "Estufa Teste",
			Descricao:      "Estufa para testes",
			Tipo:           "interno",
			Comprimento:    10.0,
			Altura:         5.0,
			Largura:        8.0,
			TempoExposicao: 12,
		}
		expectedAmbiente := &entity.Ambiente{
			Nome:           "Estufa Teste",
			Descricao:      "Estufa para testes",
			Tipo:           "interno",
			Comprimento:    10.0,
			Altura:         5.0,
			Largura:        8.0,
			TempoExposicao: 12,
			Orientacao:     "norte",
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

		var actualAmbiente entity.Ambiente
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
			Nome:           "Estufa Teste",
			Descricao:      "Estufa para testes",
			Tipo:           "interno",
			Comprimento:    10.0,
			Altura:         5.0,
			Largura:        8.0,
			TempoExposicao: 12,
		}

		mockService.On("Criar", mock.AnythingOfType("*dto.CreateAmbienteDTO")).Return((*entity.Ambiente)(nil), errors.New("erro interno do serviço"))

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
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Erro interno ao criar ambiente", response["message"])
		assert.Equal(t, "erro interno do serviço", response["details"])

		mockService.AssertExpectations(t)
	})
}

func TestAmbienteController_BuscarPorID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		// Preparação
		mockService := new(MockAmbienteService)
		controller := NewAmbienteController(mockService)

		expectedAmbiente := &entity.Ambiente{Nome: "Estufa Teste"}

		mockService.On("BuscarPorID", uint(1)).Return(expectedAmbiente, nil)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		controller.BuscarPorID(c)

		// Verificação
		assert.Equal(t, http.StatusOK, w.Code)

		var actualAmbiente entity.Ambiente
		err := json.Unmarshal(w.Body.Bytes(), &actualAmbiente)
		assert.NoError(t, err)
		assert.Equal(t, expectedAmbiente.Nome, actualAmbiente.Nome)

		mockService.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		// Preparação
		mockService := new(MockAmbienteService)
		controller := NewAmbienteController(mockService)

		mockService.On("BuscarPorID", uint(1)).Return((*entity.Ambiente)(nil), gorm.ErrRecordNotFound)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		controller.BuscarPorID(c)

		// Verificação
		assert.Equal(t, http.StatusNotFound, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Ambiente não encontrado", response["message"])
		assert.Equal(t, "recurso não encontrado", response["details"])

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
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "ID inválido", response["message"])
		assert.Equal(t, "entrada inválida", response["details"])
	})

	t.Run("Service Error", func(t *testing.T) {
		// Preparação
		mockService := new(MockAmbienteService)
		controller := NewAmbienteController(mockService)

		mockService.On("BuscarPorID", uint(1)).Return((*entity.Ambiente)(nil), errors.New("erro interno do serviço"))

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		controller.BuscarPorID(c)

		// Verificação
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Erro interno ao buscar ambiente", response["message"])
		assert.Equal(t, "erro interno do serviço", response["details"])

		mockService.AssertExpectations(t)
	})

}
func TestAmbienteController_Atualizar(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		// Preparação
		mockService := new(MockAmbienteService)
		controller := NewAmbienteController(mockService)

		updateDTO := &dto.UpdateAmbienteDTO{Nome: "Estufa Atualizada"}
		expectedAmbiente := &entity.Ambiente{Nome: "Estufa Atualizada"}

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

		var actualAmbiente entity.Ambiente
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

		mockService.On("Atualizar", uint(1), mock.AnythingOfType("*dto.UpdateAmbienteDTO")).Return((*entity.Ambiente)(nil), gorm.ErrRecordNotFound)

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
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "ID inválido", response["message"])
		assert.Equal(t, "entrada inválida", response["details"])
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
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Requisição inválida", response["message"])
		assert.Equal(t, "invalid character 'i' looking for beginning of value", response["details"])
	})

	t.Run("Service Error", func(t *testing.T) {
		// Preparação
		mockService := new(MockAmbienteService)
		controller := NewAmbienteController(mockService)

		updateDTO := &dto.UpdateAmbienteDTO{Nome: "Estufa Atualizada"}

		mockService.On("Atualizar", uint(1), mock.AnythingOfType("*dto.UpdateAmbienteDTO")).Return((*entity.Ambiente)(nil), errors.New("erro interno do serviço"))

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
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Erro interno ao atualizar ambiente", response["message"])
		assert.Equal(t, "erro interno do serviço", response["details"])

		mockService.AssertExpectations(t)
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
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "ID inválido", response["message"])
		assert.Equal(t, "entrada inválida", response["details"])
	})

	t.Run("Service Error", func(t *testing.T) {
		// Preparação
		mockService := new(MockAmbienteService)
		controller := NewAmbienteController(mockService)

		mockService.On("Deletar", uint(1)).Return(errors.New("erro interno do serviço"))

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		controller.Deletar(c)

		// Verificação
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Erro interno ao deletar ambiente", response["message"])
		assert.Equal(t, "erro interno do serviço", response["details"])

		mockService.AssertExpectations(t)
	})
}
