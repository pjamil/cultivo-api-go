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
)

// MockGeneticaService é um mock para o service.GeneticaService
type MockGeneticaService struct {
	mock.Mock
}

func (m *MockGeneticaService) Criar(geneticaDto *dto.CreateGeneticaDTO) (*dto.GeneticaResponseDTO, error) {
	args := m.Called(geneticaDto)
	return args.Get(0).(*dto.GeneticaResponseDTO), args.Error(1)
}

func (m *MockGeneticaService) ListarTodas(page, limit int) (*dto.PaginatedResponse, error) {
	args := m.Called(page, limit)
	return args.Get(0).(*dto.PaginatedResponse), args.Error(1)
}

func (m *MockGeneticaService) BuscarPorID(id uint) (*dto.GeneticaResponseDTO, error) {
	args := m.Called(id)
	return args.Get(0).(*dto.GeneticaResponseDTO), args.Error(1)
}

func (m *MockGeneticaService) Atualizar(id uint, geneticaDto *dto.UpdateGeneticaDTO) (*dto.GeneticaResponseDTO, error) {
	args := m.Called(id, geneticaDto)
	return args.Get(0).(*dto.GeneticaResponseDTO), args.Error(1)
}

func (m *MockGeneticaService) Deletar(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestGeneticaController_Listar(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success - Geneticas Encontradas", func(t *testing.T) {
		// Preparação
		mockService := new(MockGeneticaService)
		controller := NewGeneticaController(mockService)

		expectedGeneticas := []entity.Genetica{
			{Nome: "OG Kush", TipoGenetica: "Indica"},
			{Nome: "Sour Diesel", TipoGenetica: "Sativa"},
		}
		dataBytes, err := json.Marshal(expectedGeneticas)
		paginatedResponse := &dto.PaginatedResponse{
			Data:  dataBytes,
			Total: int64(len(expectedGeneticas)),
			Page:  1,
			Limit: 10,
		}

		mockService.On("ListarTodas", mock.Anything, mock.Anything).Return(paginatedResponse, nil)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/geneticas?page=1&limit=10", nil)
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

		// Convert actualResponse.Data to []entity.Genetica for comparison
		actualGeneticasBytes, _ := json.Marshal(actualResponse.Data)
		var actualGeneticas []entity.Genetica
		json.Unmarshal(actualGeneticasBytes, &actualGeneticas)

		assert.Equal(t, expectedGeneticas, actualGeneticas)

		mockService.AssertExpectations(t)
	})

	t.Run("Success - Nenhuma Genetica Encontrada", func(t *testing.T) {
		// Preparação
		mockService := new(MockGeneticaService)
		controller := NewGeneticaController(mockService)

		paginatedResponse := &dto.PaginatedResponse{
			Data:  json.RawMessage("[]"),
			Total: 0,
			Page:  1,
			Limit: 10,
		}

		mockService.On("ListarTodas", mock.Anything, mock.Anything).Return(paginatedResponse, nil)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/geneticas?page=1&limit=10", nil)
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
		assert.Equal(t, json.RawMessage("[]"), actualResponse.Data)

		mockService.AssertExpectations(t)
	})

	t.Run("Error - Erro Interno do Serviço", func(t *testing.T) {
		// Preparação
		mockService := new(MockGeneticaService)
		controller := NewGeneticaController(mockService)

		mockService.On("ListarTodas", mock.Anything, mock.Anything).Return((*dto.PaginatedResponse)(nil), errors.New("erro no serviço"))

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/geneticas?page=1&limit=10", nil)
		c.Request = req

		controller.Listar(c)

		// Verificação
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Erro interno ao listar genéticas", response["message"])
		assert.Equal(t, "erro no serviço", response["details"])

		mockService.AssertExpectations(t)
	})
}

func TestGeneticaController_Criar(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success - Genetica Criada", func(t *testing.T) {
		// Preparação
		mockService := new(MockGeneticaService)
		controller := NewGeneticaController(mockService)

		createDTO := &dto.CreateGeneticaDTO{
			Nome:            "Nova Genetica",
			Descricao:       "Descrição da nova genética",
			TipoGenetica:    "hibrida",
			TipoEspecie:     "cannabis",
			TempoFloracao:   60,
			Origem:          "California",
			Caracteristicas: "Características da nova genética",
		}
		expectedResponse := &dto.GeneticaResponseDTO{
			ID:              3,
			Nome:            "Nova Genetica",
			Descricao:       "Descrição da nova genética",
			TipoGenetica:    "hibrida",
			TipoEspecie:     "cannabis",
			TempoFloracao:   60,
			Origem:          "California",
			Caracteristicas: "Características da nova genética",
		}

		mockService.On("Criar", createDTO).Return(expectedResponse, nil)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBody, _ := json.Marshal(createDTO)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/geneticas", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req

		controller.Criar(c)

		// Verificação
		assert.Equal(t, http.StatusCreated, w.Code)

		var actualResponse dto.GeneticaResponseDTO
		err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, &actualResponse)

		mockService.AssertExpectations(t)
	})

	t.Run("Error - Payload Invalido", func(t *testing.T) {
		// Preparação
		mockService := new(MockGeneticaService)
		controller := NewGeneticaController(mockService)

		// Payload inválido (faltando campos obrigatórios)
		invalidPayload := []byte(`{"nome": "Genetica Invalida"}`)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/geneticas", bytes.NewBuffer(invalidPayload))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req

		controller.Criar(c)

		// Verificação
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Erro de validação", response["message"])
		assert.Contains(t, response["details"].(map[string]interface{})["TipoEspecie"], "Este campo é obrigatório")
		assert.Contains(t, response["details"].(map[string]interface{})["TempoFloracao"], "Este campo é obrigatório")
		assert.Contains(t, response["details"].(map[string]interface{})["Origem"], "Este campo é obrigatório")

		mockService.AssertNotCalled(t, "Criar") // O serviço não deve ser chamado em caso de payload inválido
	})
}

func TestGeneticaController_BuscarPorID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success - Genetica Encontrada", func(t *testing.T) {
		// Preparação
		mockService := new(MockGeneticaService)
		controller := NewGeneticaController(mockService)

		geneticaID := uint(1)
		expectedGenetica := &dto.GeneticaResponseDTO{ID: geneticaID, Nome: "OG Kush", TipoGenetica: "Indica"}

		mockService.On("BuscarPorID", geneticaID).Return(expectedGenetica, nil)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = gin.Params{{Key: "id", Value: "1"}}
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/geneticas/1", nil)
		c.Request = req

		controller.BuscarPorID(c)

		// Verificação
		assert.Equal(t, http.StatusOK, w.Code)

		var actualGenetica dto.GeneticaResponseDTO
		err := json.Unmarshal(w.Body.Bytes(), &actualGenetica)
		assert.NoError(t, err)
		assert.Equal(t, expectedGenetica, &actualGenetica)

		mockService.AssertExpectations(t)
	})

	t.Run("Error - ID Invalido", func(t *testing.T) {
		// Preparação
		mockService := new(MockGeneticaService)
		controller := NewGeneticaController(mockService)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = gin.Params{{Key: "id", Value: "abc"}}
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/geneticas/abc", nil)
		c.Request = req

		controller.BuscarPorID(c)

		// Verificação
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "ID inválido", response["message"])
		assert.Nil(t, response["details"])

		mockService.AssertNotCalled(t, "BuscarPorID") // O serviço não deve ser chamado em caso de ID inválido
	})

	t.Run("Error - Erro Interno do Serviço", func(t *testing.T) {
		// Preparação
		mockService := new(MockGeneticaService)
		controller := NewGeneticaController(mockService)

		geneticaID := uint(1)
		mockService.On("BuscarPorID", geneticaID).Return((*dto.GeneticaResponseDTO)(nil), errors.New("erro interno do serviço"))

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = gin.Params{{Key: "id", Value: "1"}}
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/geneticas/1", nil)
		c.Request = req

		controller.BuscarPorID(c)

		// Verificação
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Erro interno ao buscar genética", response["message"])
		assert.Nil(t, response["details"])

		mockService.AssertExpectations(t)
	})
}

func TestGeneticaController_Atualizar(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success - Genetica Atualizada", func(t *testing.T) {
		// Preparação
		mockService := new(MockGeneticaService)
		controller := NewGeneticaController(mockService)

		geneticaID := uint(1)
		updateDTO := &dto.UpdateGeneticaDTO{Nome: "OG Kush Atualizada", TipoGenetica: "Indica"}
		expectedResponse := &dto.GeneticaResponseDTO{ID: geneticaID, Nome: "OG Kush Atualizada", TipoGenetica: "Indica"}

		mockService.On("Atualizar", geneticaID, updateDTO).Return(expectedResponse, nil)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = gin.Params{{Key: "id", Value: "1"}}
		jsonBody, _ := json.Marshal(updateDTO)
		req, _ := http.NewRequest(http.MethodPut, "/api/v1/geneticas/1", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req

		controller.Atualizar(c)

		// Verificação
		assert.Equal(t, http.StatusOK, w.Code)

		var actualResponse dto.GeneticaResponseDTO
		err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, &actualResponse)

		mockService.AssertExpectations(t)
	})

	t.Run("Error - ID Invalido", func(t *testing.T) {
		// Preparação
		mockService := new(MockGeneticaService)
		controller := NewGeneticaController(mockService)

		updateDTO := &dto.UpdateGeneticaDTO{Nome: "OG Kush Atualizada", TipoGenetica: "Indica"}

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = gin.Params{{Key: "id", Value: "abc"}}
		jsonBody, _ := json.Marshal(updateDTO)
		req, _ := http.NewRequest(http.MethodPut, "/api/v1/geneticas/abc", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req

		controller.Atualizar(c)

		// Verificação
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "ID inválido", response["message"])
		assert.Nil(t, response["details"])

		mockService.AssertNotCalled(t, "Atualizar") // O serviço não deve ser chamado em caso de ID inválido
	})

	t.Run("Error - Erro Interno do Serviço", func(t *testing.T) {
		// Preparação
		mockService := new(MockGeneticaService)
		controller := NewGeneticaController(mockService)

		geneticaID := uint(1)
		updateDTO := &dto.UpdateGeneticaDTO{Nome: "OG Kush Atualizada", TipoGenetica: "Indica"}
		mockService.On("Atualizar", geneticaID, updateDTO).Return((*dto.GeneticaResponseDTO)(nil), errors.New("erro interno do serviço"))

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = gin.Params{{Key: "id", Value: "1"}}
		jsonBody, _ := json.Marshal(updateDTO)
		req, _ := http.NewRequest(http.MethodPut, "/api/v1/geneticas/1", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req

		controller.Atualizar(c)

		// Verificação
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Erro interno ao atualizar genética", response["message"])
		assert.Nil(t, response["details"])

		mockService.AssertExpectations(t)
	})
}

func TestGeneticaController_Deletar(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success - Genetica Deletada", func(t *testing.T) {
		// Preparação
		mockService := new(MockGeneticaService)
		controller := NewGeneticaController(mockService)

		geneticaID := uint(1)
		mockService.On("Deletar", geneticaID).Return(nil)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = gin.Params{{Key: "id", Value: "1"}}
		req, _ := http.NewRequest(http.MethodDelete, "/api/v1/geneticas/1", nil)
		c.Request = req

		controller.Deletar(c)

		// Verificação
		assert.Equal(t, http.StatusOK, w.Code) // Changed from StatusNoContent to StatusOK

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Genética deletada com sucesso", response["message"])

		mockService.AssertExpectations(t)
	})

	t.Run("Error - ID Invalido", func(t *testing.T) {
		// Preparação
		mockService := new(MockGeneticaService)
		controller := NewGeneticaController(mockService)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = gin.Params{{Key: "id", Value: "abc"}}
		req, _ := http.NewRequest(http.MethodDelete, "/api/v1/geneticas/abc", nil)
		c.Request = req

		controller.Deletar(c)

		// Verificação
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "ID inválido", response["message"])
		assert.Nil(t, response["details"])

		mockService.AssertNotCalled(t, "Deletar") // O serviço não deve ser chamado em caso de ID inválido
	})

	t.Run("Error - Erro Interno do Serviço", func(t *testing.T) {
		// Preparação
		mockService := new(MockGeneticaService)
		controller := NewGeneticaController(mockService)

		geneticaID := uint(1)
		mockService.On("Deletar", geneticaID).Return(errors.New("erro interno do serviço"))

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = gin.Params{{Key: "id", Value: "1"}}
		req, _ := http.NewRequest(http.MethodDelete, "/api/v1/geneticas/1", nil)
		c.Request = req

		controller.Deletar(c)

		// Verificação
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Erro interno ao deletar genética", response["message"])
		assert.Nil(t, response["details"])

		mockService.AssertExpectations(t)
	})
}
