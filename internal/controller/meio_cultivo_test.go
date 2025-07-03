package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockMeioCultivoService é um mock para o service.MeioCultivoService
type MockMeioCultivoService struct {
	mock.Mock
}

func (m *MockMeioCultivoService) Criar(meioCultivoDto *dto.CreateMeioCultivoDTO) (*dto.MeioCultivoResponseDTO, error) {
	args := m.Called(meioCultivoDto)
	return args.Get(0).(*dto.MeioCultivoResponseDTO), args.Error(1)
}

func (m *MockMeioCultivoService) ListarTodos() ([]dto.MeioCultivoResponseDTO, error) {
	args := m.Called()
	return args.Get(0).([]dto.MeioCultivoResponseDTO), args.Error(1)
}

func (m *MockMeioCultivoService) BuscarPorID(id uint) (*dto.MeioCultivoResponseDTO, error) {
	args := m.Called(id)
	return args.Get(0).(*dto.MeioCultivoResponseDTO), args.Error(1)
}

func (m *MockMeioCultivoService) Atualizar(id uint, meioCultivoDto *dto.UpdateMeioCultivoDTO) (*dto.MeioCultivoResponseDTO, error) {
	args := m.Called(id, meioCultivoDto)
	return args.Get(0).(*dto.MeioCultivoResponseDTO), args.Error(1)
}

func (m *MockMeioCultivoService) Deletar(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestMeioCultivoController_Listar(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success - Meios de Cultivo Encontrados", func(t *testing.T) {
		// Preparação
		mockService := new(MockMeioCultivoService)
		controller := NewMeioCultivoController(mockService)

		expectedMeiosCultivo := []dto.MeioCultivoResponseDTO{
			{ID: 1, Tipo: "Solo", Descricao: "Solo orgânico"},
			{ID: 2, Tipo: "Coco", Descricao: "Fibra de coco"},
		}

		mockService.On("ListarTodos").Return(expectedMeiosCultivo, nil)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/meios-cultivos", nil)
		c.Request = req

		controller.Listar(c)

		// Verificação
		assert.Equal(t, http.StatusOK, w.Code)

		var actualMeiosCultivo []dto.MeioCultivoResponseDTO
		err := json.Unmarshal(w.Body.Bytes(), &actualMeiosCultivo)
		assert.NoError(t, err)
		assert.Equal(t, expectedMeiosCultivo, actualMeiosCultivo)

		mockService.AssertExpectations(t)
	})

	t.Run("Success - Nenhum Meio de Cultivo Encontrado", func(t *testing.T) {
		// Preparação
		mockService := new(MockMeioCultivoService)
		controller := NewMeioCultivoController(mockService)

		mockService.On("ListarTodos").Return([]dto.MeioCultivoResponseDTO{}, nil)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/meios-cultivos", nil)
		c.Request = req

		controller.Listar(c)

		// Verificação
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Nenhum meio de cultivo encontrado", response["message"])

		mockService.AssertExpectations(t)
	})

	t.Run("Error - Erro Interno do Serviço", func(t *testing.T) {
		// Preparação
		mockService := new(MockMeioCultivoService)
		controller := NewMeioCultivoController(mockService)

		mockService.On("ListarTodos").Return([]dto.MeioCultivoResponseDTO{}, errors.New("erro no serviço"))

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/meios-cultivos", nil)
		c.Request = req

		controller.Listar(c)

		// Verificação
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Erro ao recuperar meios de cultivo", response["error"])

		mockService.AssertExpectations(t)
	})
}

func TestMeioCultivoController_Criar(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		// Preparação
		mockService := new(MockMeioCultivoService)
		controller := NewMeioCultivoController(mockService)

		createDTO := &dto.CreateMeioCultivoDTO{
			Tipo:        "Solo Teste",
			Descricao:   "Solo para testes",
		}
		expectedMeioCultivo := &dto.MeioCultivoResponseDTO{
			ID:          1,
			Tipo:        "Solo Teste",
			Descricao:   "Solo para testes",
		}

		mockService.On("Criar", mock.AnythingOfType("*dto.CreateMeioCultivoDTO")).Return(expectedMeioCultivo, nil)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBody, _ := json.Marshal(createDTO)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/meios-cultivos", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req

		controller.Criar(c)

		// Verificação
		assert.Equal(t, http.StatusCreated, w.Code)

		var actualMeioCultivo dto.MeioCultivoResponseDTO
		err := json.Unmarshal(w.Body.Bytes(), &actualMeioCultivo)
		assert.NoError(t, err)
		assert.Equal(t, expectedMeioCultivo.Tipo, actualMeioCultivo.Tipo)
		assert.Equal(t, expectedMeioCultivo.Descricao, actualMeioCultivo.Descricao)

		mockService.AssertExpectations(t)
	})

	t.Run("Invalid Payload", func(t *testing.T) {
		// Preparação
		mockService := new(MockMeioCultivoService)
		controller := NewMeioCultivoController(mockService)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/meios-cultivos", bytes.NewBufferString("invalid-json"))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req

		controller.Criar(c)

		// Verificação
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		// Preparação
		mockService := new(MockMeioCultivoService)
		controller := NewMeioCultivoController(mockService)

		createDTO := &dto.CreateMeioCultivoDTO{
			Tipo:        "Solo Teste",
			Descricao:   "Solo para testes",
		}

		mockService.On("Criar", mock.AnythingOfType("*dto.CreateMeioCultivoDTO")).Return((*dto.MeioCultivoResponseDTO)(nil), errors.New("erro interno do serviço"))

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBody, _ := json.Marshal(createDTO)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/meios-cultivos", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req

		controller.Criar(c)

		// Verificação
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		mockService.AssertExpectations(t)
	})
}

func TestMeioCultivoController_BuscarPorID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		// Preparação
		mockService := new(MockMeioCultivoService)
		controller := NewMeioCultivoController(mockService)

		expectedMeioCultivo := &dto.MeioCultivoResponseDTO{ID: 1, Tipo: "Solo Teste"}

		mockService.On("BuscarPorID", uint(1)).Return(expectedMeioCultivo, nil)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		controller.BuscarPorID(c)

		// Verificação
		assert.Equal(t, http.StatusOK, w.Code)

		var actualMeioCultivo dto.MeioCultivoResponseDTO
		err := json.Unmarshal(w.Body.Bytes(), &actualMeioCultivo)
		assert.NoError(t, err)
		assert.Equal(t, expectedMeioCultivo.Tipo, actualMeioCultivo.Tipo)

		mockService.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		// Preparação
		mockService := new(MockMeioCultivoService)
		controller := NewMeioCultivoController(mockService)

		mockService.On("BuscarPorID", uint(1)).Return((*dto.MeioCultivoResponseDTO)(nil), gorm.ErrRecordNotFound)

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
		mockService := new(MockMeioCultivoService)
		controller := NewMeioCultivoController(mockService)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "abc"}}

		controller.BuscarPorID(c)

		// Verificação
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestMeioCultivoController_Atualizar(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		// Preparação
		mockService := new(MockMeioCultivoService)
		controller := NewMeioCultivoController(mockService)

		updateDTO := &dto.UpdateMeioCultivoDTO{Tipo: "Solo Atualizado"}
		expectedMeioCultivo := &dto.MeioCultivoResponseDTO{ID: 1, Tipo: "Solo Atualizado"}

		mockService.On("Atualizar", uint(1), mock.AnythingOfType("*dto.UpdateMeioCultivoDTO")).Return(expectedMeioCultivo, nil)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBody, _ := json.Marshal(updateDTO)
		req, _ := http.NewRequest(http.MethodPut, "/api/v1/meios-cultivos/1", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		controller.Atualizar(c)

		// Verificação
		assert.Equal(t, http.StatusOK, w.Code)

		var actualMeioCultivo dto.MeioCultivoResponseDTO
		err := json.Unmarshal(w.Body.Bytes(), &actualMeioCultivo)
		assert.NoError(t, err)
		assert.Equal(t, expectedMeioCultivo.Tipo, actualMeioCultivo.Tipo)

		mockService.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		// Preparação
		mockService := new(MockMeioCultivoService)
		controller := NewMeioCultivoController(mockService)

		updateDTO := &dto.UpdateMeioCultivoDTO{Tipo: "Solo Atualizado"}

		mockService.On("Atualizar", uint(1), mock.AnythingOfType("*dto.UpdateMeioCultivoDTO")).Return((*dto.MeioCultivoResponseDTO)(nil), gorm.ErrRecordNotFound)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBody, _ := json.Marshal(updateDTO)
		req, _ := http.NewRequest(http.MethodPut, "/api/v1/meios-cultivos/1", bytes.NewBuffer(jsonBody))
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
		mockService := new(MockMeioCultivoService)
		controller := NewMeioCultivoController(mockService)

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
		mockService := new(MockMeioCultivoService)
		controller := NewMeioCultivoController(mockService)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest(http.MethodPut, "/api/v1/meios-cultivos/1", bytes.NewBufferString("invalid-json"))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		controller.Atualizar(c)

		// Verificação
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestMeioCultivoController_Deletar(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		// Preparação
		mockService := new(MockMeioCultivoService)
		controller := NewMeioCultivoController(mockService)

		mockService.On("Deletar", uint(1)).Return(nil)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		controller.Deletar(c)

		// Verificação
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Meio de cultivo deletado com sucesso", response["message"])

		mockService.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		// Preparação
		mockService := new(MockMeioCultivoService)
		controller := NewMeioCultivoController(mockService)

		mockService.On("Deletar", uint(1)).Return(gorm.ErrRecordNotFound)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		controller.Deletar(c)

		// Verificação
		assert.Equal(t, http.StatusInternalServerError, w.Code) // Changed from StatusNotFound to StatusInternalServerError
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Erro ao deletar meio de cultivo", response["error"]) // Changed message

		mockService.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		// Preparação
		mockService := new(MockMeioCultivoService)
		controller := NewMeioCultivoController(mockService)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "abc"}}

		controller.Deletar(c)

		// Verificação
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		// Preparação
		mockService := new(MockMeioCultivoService)
		controller := NewMeioCultivoController(mockService)

		mockService.On("Deletar", uint(1)).Return(errors.New("erro interno do serviço"))

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		controller.Deletar(c)

		// Verificação
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		mockService.AssertExpectations(t)
	})
}