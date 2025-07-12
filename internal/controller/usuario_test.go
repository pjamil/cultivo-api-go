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

// MockUsuarioService é um mock para o service.UsuarioService
type MockUsuarioService struct {
	mock.Mock
}

func (m *MockUsuarioService) Criar(usuarioDto *dto.UsuarioCreateDTO) (*dto.UsuarioResponseDTO, error) {
	args := m.Called(usuarioDto)
	return args.Get(0).(*dto.UsuarioResponseDTO), args.Error(1)
}

func (m *MockUsuarioService) ListarTodos(page, limit int) (*dto.PaginatedResponse, error) {
	args := m.Called(page, limit)
	return args.Get(0).(*dto.PaginatedResponse), args.Error(1)
}

func (m *MockUsuarioService) BuscarPorID(id uint) (*dto.UsuarioResponseDTO, error) {
	args := m.Called(id)
	return args.Get(0).(*dto.UsuarioResponseDTO), args.Error(1)
}

func (m *MockUsuarioService) Atualizar(id uint, usuarioDto *dto.UsuarioUpdateDTO) (*dto.UsuarioResponseDTO, error) {
	args := m.Called(id, usuarioDto)
	return args.Get(0).(*dto.UsuarioResponseDTO), args.Error(1)
}

func (m *MockUsuarioService) Deletar(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUsuarioService) Login(loginDto *dto.LoginPayload) (string, error) {
	args := m.Called(loginDto)
	return args.Get(0).(string), args.Error(1)
}

func TestUsuarioController_Criar(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		// Preparação
		mockService := new(MockUsuarioService)
		controller := NewUsuarioController(mockService)

		createDTO := &dto.UsuarioCreateDTO{
			Nome:  "Teste",
			Email: "teste@example.com",
			Senha: "password123",
		}
		expectedUsuario := &dto.UsuarioResponseDTO{
			ID:    1,
			Nome:  "Teste",
			Email: "teste@example.com",
		}

		mockService.On("Criar", createDTO).Return(expectedUsuario, nil)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBody, _ := json.Marshal(createDTO)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/usuarios", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req

		controller.Criar(c)

		// Verificação
		assert.Equal(t, http.StatusCreated, w.Code)

		var actualUsuario dto.UsuarioResponseDTO
		err := json.Unmarshal(w.Body.Bytes(), &actualUsuario)
		assert.NoError(t, err)
		assert.Equal(t, expectedUsuario, &actualUsuario)

		mockService.AssertExpectations(t)
	})

	t.Run("Invalid Payload", func(t *testing.T) {
		// Preparação
		mockService := new(MockUsuarioService)
		controller := NewUsuarioController(mockService)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/usuarios", bytes.NewBufferString("invalid-json"))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req

		controller.Criar(c)

		// Verificação
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Email Conflict", func(t *testing.T) {
		// Preparação
		mockService := new(MockUsuarioService)
		controller := NewUsuarioController(mockService)

		createDTO := &dto.UsuarioCreateDTO{
			Nome:  "Teste",
			Email: "teste@example.com",
			Senha: "password123",
		}

		mockService.On("Criar", createDTO).Return((*dto.UsuarioResponseDTO)(nil), errors.New("duplicate entry for key 'email'"))

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBody, _ := json.Marshal(createDTO)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/usuarios", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req

		controller.Criar(c)

		// Verificação
		assert.Equal(t, http.StatusConflict, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "E-mail já cadastrado", response["error"])

		mockService.AssertExpectations(t)
	})

	t.Run("Service Error", func(t *testing.T) {
		// Preparação
		mockService := new(MockUsuarioService)
		controller := NewUsuarioController(mockService)

		createDTO := &dto.UsuarioCreateDTO{
			Nome:  "Teste",
			Email: "teste@example.com",
			Senha: "password123",
		}

		mockService.On("Criar", createDTO).Return((*dto.UsuarioResponseDTO)(nil), errors.New("erro interno do serviço"))

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBody, _ := json.Marshal(createDTO)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/usuarios", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req

		controller.Criar(c)

		// Verificação
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "erro interno do serviço", response["error"])

		mockService.AssertExpectations(t)
	})
}

func TestUsuarioController_Listar(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success - Usuarios Encontrados", func(t *testing.T) {
		// Preparação
		mockService := new(MockUsuarioService)
		controller := NewUsuarioController(mockService)

		expectedUsuarios := []models.Usuario{
			{Nome: "Usuario 1", Email: "user1@example.com"},
			{Nome: "Usuario 2", Email: "user2@example.com"},
		}
		paginatedResponse := &dto.PaginatedResponse{
			Data:  expectedUsuarios,
			Total: int64(len(expectedUsuarios)),
			Page:  1,
			Limit: 10,
		}

		mockService.On("ListarTodos", mock.Anything, mock.Anything).Return(paginatedResponse, nil)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/usuarios?page=1&limit=10", nil)
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

		// Convert actualResponse.Data to []models.Usuario for comparison
		actualUsuariosBytes, _ := json.Marshal(actualResponse.Data)
		var actualUsuarios []models.Usuario
		json.Unmarshal(actualUsuariosBytes, &actualUsuarios)

		assert.Equal(t, expectedUsuarios, actualUsuarios)

		mockService.AssertExpectations(t)
	})

	t.Run("Success - Nenhum Usuario Encontrado", func(t *testing.T) {
		// Preparação
		mockService := new(MockUsuarioService)
		controller := NewUsuarioController(mockService)

		paginatedResponse := &dto.PaginatedResponse{
			Data:  []models.Usuario{},
			Total: 0,
			Page:  1,
			Limit: 10,
		}

		mockService.On("ListarTodos", mock.Anything, mock.Anything).Return(paginatedResponse, nil)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/usuarios?page=1&limit=10", nil)
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
		mockService := new(MockUsuarioService)
		controller := NewUsuarioController(mockService)

		mockService.On("ListarTodos", mock.Anything, mock.Anything).Return((*dto.PaginatedResponse)(nil), errors.New("erro interno do serviço"))

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/usuarios?page=1&limit=10", nil)
		c.Request = req

		controller.Listar(c)

		// Verificação
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "erro interno do serviço", response["error"])

		mockService.AssertExpectations(t)
	})
}

func TestUsuarioController_BuscarPorID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		// Preparação
		mockService := new(MockUsuarioService)
		controller := NewUsuarioController(mockService)

																																																																				expectedUsuario := &dto.UsuarioResponseDTO{ID: 1, Nome: "Usuario Teste", Email: "teste@example.com"}

		mockService.On("BuscarPorID", uint(1)).Return(expectedUsuario, nil)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		controller.BuscarPorID(c)

		// Verificação
		assert.Equal(t, http.StatusOK, w.Code)

		var actualUsuario dto.UsuarioResponseDTO
		err := json.Unmarshal(w.Body.Bytes(), &actualUsuario)
		assert.NoError(t, err)
		assert.Equal(t, expectedUsuario, &actualUsuario)

		mockService.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		// Preparação
		mockService := new(MockUsuarioService)
		controller := NewUsuarioController(mockService)

		mockService.On("BuscarPorID", uint(1)).Return((*dto.UsuarioResponseDTO)(nil), errors.New("record not found"))

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
		mockService := new(MockUsuarioService)
		controller := NewUsuarioController(mockService)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "abc"}}

		controller.BuscarPorID(c)

		// Verificação
		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "ID inválido, deve ser um número inteiro positivo", response["error"])
	})
}

func TestUsuarioController_Atualizar(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		// Preparação
		mockService := new(MockUsuarioService)
		controller := NewUsuarioController(mockService)

		updateDTO := &dto.UsuarioUpdateDTO{Nome: "Usuario Atualizado"}
		expectedUsuario := &dto.UsuarioResponseDTO{Nome: "Usuario Atualizado", Email: "teste@example.com"}
		expectedUsuario.ID = 1

		mockService.On("Atualizar", uint(1), updateDTO).Return(expectedUsuario, nil)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBody, _ := json.Marshal(updateDTO)
		req, _ := http.NewRequest(http.MethodPut, "/api/v1/usuarios/1", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		controller.Atualizar(c)

		// Verificação
		assert.Equal(t, http.StatusOK, w.Code)

		var actualUsuario dto.UsuarioResponseDTO
		err := json.Unmarshal(w.Body.Bytes(), &actualUsuario)
		assert.NoError(t, err)
		assert.Equal(t, expectedUsuario, &actualUsuario)

		mockService.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		// Preparação
		mockService := new(MockUsuarioService)
		controller := NewUsuarioController(mockService)

		updateDTO := &dto.UsuarioUpdateDTO{Nome: "Usuario Atualizado"}

		mockService.On("Atualizar", uint(1), updateDTO).Return((*dto.UsuarioResponseDTO)(nil), gorm.ErrRecordNotFound)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBody, _ := json.Marshal(updateDTO)
		req, _ := http.NewRequest(http.MethodPut, "/api/v1/usuarios/1", bytes.NewBuffer(jsonBody))
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
		mockService := new(MockUsuarioService)
		controller := NewUsuarioController(mockService)

		updateDTO := &dto.UsuarioUpdateDTO{Nome: "Usuario Atualizado"}

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBody, _ := json.Marshal(updateDTO)
		req, _ := http.NewRequest(http.MethodPut, "/api/v1/usuarios/abc", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "abc"}}

		controller.Atualizar(c)

		// Verificação
		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "ID inválido, deve ser um número inteiro positivo", response["error"])

		mockService.AssertNotCalled(t, "Atualizar")
	})

	t.Run("Invalid Payload", func(t *testing.T) {
		// Preparação
		mockService := new(MockUsuarioService)
		controller := NewUsuarioController(mockService)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest(http.MethodPut, "/api/v1/usuarios/1", bytes.NewBufferString("invalid-json"))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		controller.Atualizar(c)

		// Verificação
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		// Preparação
		mockService := new(MockUsuarioService)
		controller := NewUsuarioController(mockService)

		updateDTO := &dto.UsuarioUpdateDTO{Nome: "Usuario Atualizado"}

		mockService.On("Atualizar", uint(1), updateDTO).Return((*dto.UsuarioResponseDTO)(nil), errors.New("erro interno do serviço"))

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBody, _ := json.Marshal(updateDTO)
		req, _ := http.NewRequest(http.MethodPut, "/api/v1/usuarios/1", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		controller.Atualizar(c)

		// Verificação
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "erro interno do serviço", response["error"])

		mockService.AssertExpectations(t)
	})
}

func TestUsuarioController_Deletar(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		// Preparação
		mockService := new(MockUsuarioService)
		controller := NewUsuarioController(mockService)

		usuarioID := uint(1)
		mockService.On("Deletar", usuarioID).Return(nil)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = gin.Params{{Key: "id", Value: "1"}}
		req, _ := http.NewRequest(http.MethodDelete, "/api/v1/usuarios/1", nil)
		c.Request = req

		controller.Deletar(c)

		// Verificação
		assert.Equal(t, http.StatusNoContent, w.Code)

		// Não deve haver corpo de resposta para StatusNoContent, mas o Gin pode adicionar um vazio.
		// Se houver um corpo, ele deve ser vazio ou conter uma mensagem de sucesso.
		if w.Body.Len() > 0 {
			var response map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Usuário deletado com sucesso", response["message"])
		}

		mockService.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		// Preparação
		mockService := new(MockUsuarioService)
		controller := NewUsuarioController(mockService)

		usuarioID := uint(999)
		mockService.On("Deletar", usuarioID).Return(gorm.ErrRecordNotFound)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = gin.Params{{Key: "id", Value: "999"}}
		req, _ := http.NewRequest(http.MethodDelete, "/api/v1/usuarios/999", nil)
		c.Request = req

		controller.Deletar(c)

		// Verificação
		assert.Equal(t, http.StatusNotFound, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "registro não encontrado", response["error"])

		mockService.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		// Preparação
		mockService := new(MockUsuarioService)
		controller := NewUsuarioController(mockService)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = gin.Params{{Key: "id", Value: "abc"}}
		req, _ := http.NewRequest(http.MethodDelete, "/api/v1/usuarios/abc", nil)
		c.Request = req

		controller.Deletar(c)

		// Verificação
		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "ID inválido, deve ser um número inteiro positivo", response["error"])

		mockService.AssertNotCalled(t, "Deletar")
	})

	t.Run("Service Error", func(t *testing.T) {
		// Preparação
		mockService := new(MockUsuarioService)
		controller := NewUsuarioController(mockService)

		usuarioID := uint(1)
		mockService.On("Deletar", usuarioID).Return(errors.New("erro interno do serviço"))

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = gin.Params{{Key: "id", Value: "1"}}
		req, _ := http.NewRequest(http.MethodDelete, "/api/v1/usuarios/1", nil)
		c.Request = req

		controller.Deletar(c)

		// Verificação
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "erro interno do serviço", response["error"])

		mockService.AssertExpectations(t)
	})
}