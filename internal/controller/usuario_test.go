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

// MockUsuarioService é um mock para o service.UsuarioService
type MockUsuarioService struct {
	mock.Mock
}

func (m *MockUsuarioService) Criar(usuarioDto *dto.UsuarioCreateDTO) (*dto.UsuarioResponseDTO, error) {
	args := m.Called(usuarioDto)
	return args.Get(0).(*dto.UsuarioResponseDTO), args.Error(1)
}

func (m *MockUsuarioService) ListarTodos(page, limit int) ([]dto.UsuarioResponseDTO, int64, error) {
	args := m.Called(page, limit)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]dto.UsuarioResponseDTO), args.Get(1).(int64), args.Error(2)
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
			Nome:         "Teste",
			Email:        "teste@example.com",
			Senha:        "password123",
			Preferencias: json.RawMessage([]byte("null")),
		}
		expectedUsuario := &dto.UsuarioResponseDTO{
			ID:           1,
			Nome:         "Teste",
			Email:        "teste@example.com",
			Preferencias: json.RawMessage([]byte("null")),
		}

		mockService.On("Criar", mock.MatchedBy(func(arg *dto.UsuarioCreateDTO) bool {
			return assert.Equal(t, createDTO.Nome, arg.Nome) &&
				assert.Equal(t, createDTO.Email, arg.Email) &&
				assert.Equal(t, createDTO.Senha, arg.Senha) &&
				bytes.Equal(createDTO.Preferencias, arg.Preferencias)
		})).Return(expectedUsuario, nil)

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
			Nome:         "Teste",
			Email:        "teste@example.com",
			Senha:        "password123",
			Preferencias: json.RawMessage([]byte("null")),
		}

		mockService.On("Criar", mock.MatchedBy(func(arg *dto.UsuarioCreateDTO) bool {
			return assert.Equal(t, createDTO.Nome, arg.Nome) &&
				assert.Equal(t, createDTO.Email, arg.Email) &&
				assert.Equal(t, createDTO.Senha, arg.Senha) &&
				bytes.Equal(createDTO.Preferencias, arg.Preferencias)
		})).Return((*dto.UsuarioResponseDTO)(nil), errors.New("duplicate entry for key 'email'"))

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
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "E-mail já cadastrado", response["message"])
		assert.Nil(t, response["details"])

		mockService.AssertExpectations(t)
	})

	t.Run("Service Error", func(t *testing.T) {
		// Preparação
		mockService := new(MockUsuarioService)
		controller := NewUsuarioController(mockService)

		createDTO := &dto.UsuarioCreateDTO{
			Nome:         "Teste",
			Email:        "teste@example.com",
			Senha:        "password123",
			Preferencias: json.RawMessage([]byte("null")),
		}

		mockService.On("Criar", mock.MatchedBy(func(arg *dto.UsuarioCreateDTO) bool {
			return assert.Equal(t, createDTO.Nome, arg.Nome) &&
				assert.Equal(t, createDTO.Email, arg.Email) &&
				assert.Equal(t, createDTO.Senha, arg.Senha) &&
				bytes.Equal(createDTO.Preferencias, arg.Preferencias)
		})).Return((*dto.UsuarioResponseDTO)(nil), errors.New("erro interno do serviço"))

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
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Erro interno ao criar usuário", response["message"])
		assert.Equal(t, "erro interno do serviço", response["details"])

		mockService.AssertExpectations(t)
	})
}

func TestUsuarioController_Listar(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success - Usuarios Encontrados", func(t *testing.T) {
		// Preparação
		mockService := new(MockUsuarioService)
		controller := NewUsuarioController(mockService)

		expectedUsuarios := []dto.UsuarioResponseDTO{
			{ID: 1, Nome: "Usuario 1", Email: "user1@example.com", Preferencias: json.RawMessage("null")},
			{ID: 2, Nome: "Usuario 2", Email: "user2@example.com", Preferencias: json.RawMessage("null")},
		}
		dataBytes, err := json.Marshal(expectedUsuarios)
		paginatedResponse := &dto.PaginatedResponse{
			Data:  dataBytes,
			Total: int64(len(expectedUsuarios)),
			Page:  1,
			Limit: 10,
		}
		// The mock now returns the slice of DTOs and the total count directly
		mockService.On("ListarTodos", mock.Anything, mock.Anything).Return(expectedUsuarios, int64(len(expectedUsuarios)), nil)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/usuarios?page=1&limit=10", nil)
		c.Request = req

		controller.Listar(c)

		// Verificação
		assert.Equal(t, http.StatusOK, w.Code)

		var actualPaginatedResponse dto.PaginatedResponse
		err = json.Unmarshal(w.Body.Bytes(), &actualPaginatedResponse)
		assert.NoError(t, err)
		assert.Equal(t, paginatedResponse.Total, actualPaginatedResponse.Total)
		assert.Equal(t, paginatedResponse.Page, actualPaginatedResponse.Page)
		assert.Equal(t, paginatedResponse.Limit, actualPaginatedResponse.Limit)

		// Unmarshal the Data field into a slice of entity.Usuario
		var actualUsers []dto.UsuarioResponseDTO
		err = json.Unmarshal(actualPaginatedResponse.Data, &actualUsers) // This will unmarshal into the DTO slice
		assert.NoError(t, err)

		assert.Equal(t, expectedUsuarios, actualUsers)

		mockService.AssertExpectations(t)
	})

	t.Run("Success - Nenhum Usuario Encontrado", func(t *testing.T) {
		// Preparação
		mockService := new(MockUsuarioService)
		controller := NewUsuarioController(mockService)

		// The mock now returns an empty slice of DTOs and a total count of 0
		mockService.On("ListarTodos", mock.Anything, mock.Anything).Return(
			[]dto.UsuarioResponseDTO{}, // Empty slice for Data
			int64(0),                   // Total count
			nil)                        // No error

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
		assert.Equal(t, int64(0), actualResponse.Total)
		assert.Equal(t, 1, actualResponse.Page)
		assert.Equal(t, 10, actualResponse.Limit)
		assert.Equal(t, json.RawMessage("[]"), actualResponse.Data)

		mockService.AssertExpectations(t)
	})

	t.Run("Error - Erro Interno do Serviço", func(t *testing.T) {
		// Preparação
		mockService := new(MockUsuarioService)
		controller := NewUsuarioController(mockService)

		mockService.On("ListarTodos", mock.Anything, mock.Anything).Return([]dto.UsuarioResponseDTO{}, int64(0), errors.New("erro interno do serviço"))

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/usuarios?page=1&limit=10", nil)
		c.Request = req

		controller.Listar(c)

		// Verificação
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Erro interno ao listar usuários", response["message"])
		assert.Equal(t, "erro interno do servidor", response["details"])
	})
}

func TestUsuarioController_BuscarPorID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		// Preparação
		mockService := new(MockUsuarioService)
		controller := NewUsuarioController(mockService)

		expectedUsuario := &dto.UsuarioResponseDTO{ID: 1, Nome: "Usuario Teste", Email: "teste@example.com", Preferencias: json.RawMessage("null")}

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
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Usuário não encontrado", response["message"])
		assert.Equal(t, "recurso não encontrado", response["details"])

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
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "ID inválido", response["message"])
		assert.Equal(t, "entrada inválida", response["details"])
	})
}

func TestUsuarioController_Atualizar(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		// Preparação
		mockService := new(MockUsuarioService)
		controller := NewUsuarioController(mockService)

		updateDTO := &dto.UsuarioUpdateDTO{
			Nome:         "Usuario Atualizado",
			Preferencias: json.RawMessage([]byte("null")),
		}
		expectedUsuario := &dto.UsuarioResponseDTO{
			ID:           1,
			Nome:         "Usuario Atualizado",
			Email:        "teste@example.com",
			Preferencias: json.RawMessage([]byte("null")),
		}

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

		updateDTO := &dto.UsuarioUpdateDTO{Nome: "Usuario Atualizado", Preferencias: json.RawMessage([]byte("null"))}

		mockService.On("Atualizar", uint(1), mock.AnythingOfType("*dto.UsuarioUpdateDTO")).Return((*dto.UsuarioResponseDTO)(nil), gorm.ErrRecordNotFound)

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
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "ID inválido", response["message"])
		assert.Equal(t, "entrada inválida", response["details"])

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
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Requisição inválida", response["message"])
		assert.Contains(t, response["details"], "invalid character")
	})

	t.Run("Service Error", func(t *testing.T) {
		// Preparação
		mockService := new(MockUsuarioService)
		controller := NewUsuarioController(mockService)

		updateDTO := &dto.UsuarioUpdateDTO{Nome: "Usuario Atualizado", Preferencias: json.RawMessage([]byte("null"))}

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
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Erro interno ao atualizar usuário", response["message"])
		assert.Equal(t, "erro interno do servidor", response["details"])

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

		// Para StatusNoContent (204), o corpo da resposta deve ser vazio.
		assert.Empty(t, w.Body.String())

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
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Usuário não encontrado", response["message"])
		assert.Equal(t, gorm.ErrRecordNotFound.Error(), response["details"])

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
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "ID inválido", response["message"])
		assert.Equal(t, "entrada inválida", response["details"])

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
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Erro interno ao deletar usuário", response["message"])
		assert.Equal(t, "erro interno do servidor", response["details"])

		mockService.AssertExpectations(t)
	})
}
