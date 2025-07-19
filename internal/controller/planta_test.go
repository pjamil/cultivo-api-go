package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/entity"
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

func (m *MockPlantaService) ListarTodas(page, limit int) ([]dto.PlantaResponseDTO, int64, error) {
	args := m.Called(page, limit)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]dto.PlantaResponseDTO), args.Get(1).(int64), args.Error(2)
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

func (m *MockPlantaService) RegistrarFato(id uint, tipo entity.RegistroTipo, titulo string, conteudo string) error {
	args := m.Called(id, tipo, titulo, conteudo)
	return args.Error(0)
}

func (m *MockPlantaService) BuscarPorEspecie(especie entity.Especie) ([]entity.Planta, error) {
	args := m.Called(especie)
	return args.Get(0).([]entity.Planta), args.Error(1)
}

func (m *MockPlantaService) BuscarPorStatus(status string) ([]entity.Planta, error) {
	args := m.Called(status)
	return args.Get(0).([]entity.Planta), args.Error(1)
}

func TestPlantaController_Listar(t *testing.T) {
	// Configuração do Gin em modo de teste
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		// Preparação
		mockService := new(MockPlantaService)
		controller := NewPlantaController(mockService)

		expectedPlantas := []dto.PlantaResponseDTO{
			{ID: 1, Nome: "Sativa"},
			{ID: 2, Nome: "Indica"},
		}
		dataBytes, err := json.Marshal(expectedPlantas)
		paginatedResponse := &dto.PaginatedResponse{
			Data:  dataBytes,
			Total: int64(len(expectedPlantas)),
			// The ListarTodas mock expects a []dto.PlantaResponseDTO and an int64 for total.
			// The PaginatedResponse struct is used for the controller's response, not the service's return.
			// So, we need to pass the actual slice and total count to the mock.
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
		err = json.Unmarshal(w.Body.Bytes(), &actualResponse)
		assert.NoError(t, err)
		assert.Equal(t, paginatedResponse.Total, actualResponse.Total)
		assert.Equal(t, paginatedResponse.Page, actualResponse.Page)
		assert.Equal(t, paginatedResponse.Limit, actualResponse.Limit)

		// Convert actualResponse.Data to []entity.Planta for comparison
		actualPlantasBytes, _ := json.Marshal(actualResponse.Data)
		var actualPlantas []dto.PlantaResponseDTO
		json.Unmarshal(actualPlantasBytes, &actualPlantas)

		assert.Equal(t, expectedPlantas, actualPlantas)

		mockService.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		// Preparação
		mockService := new(MockPlantaService)
		controller := NewPlantaController(mockService)

		mockService.On("ListarTodas", mock.Anything, mock.Anything).Return([]dto.PlantaResponseDTO{}, int64(0), errors.New("erro no serviço"))

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/plantas?page=1&limit=10", nil)
		c.Request = req

		controller.Listar(c)

		// Verificação
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Erro interno ao listar plantas", response["message"])
		assert.Equal(t, "erro no serviço", response["details"])

		mockService.AssertExpectations(t)
	})
}

func TestPlantaController_BuscarPorID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		// Preparação
		mockService := new(MockPlantaService)
		controller := NewPlantaController(mockService)

		fixedTime, _ := time.Parse("2006-01-02", "2023-01-01")
		notasContent := "Algumas notas."

		expectedPlanta := &dto.PlantaResponseDTO{
			ID:            1,
			Nome:          "Sativa",
			ComecandoDe:   "semente",
			Especie:       "Sativa",
			DataPlantio:   &fixedTime,
			Status:        "vegetativo",
			Notas:         &notasContent,
			GeneticaID:    1,
			MeioCultivoID: 1,
			AmbienteID:    1,
			UsuarioID:     1,
		}

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
		assert.Equal(t, expectedPlanta.ID, actualPlanta.ID)
		assert.Equal(t, expectedPlanta.Nome, actualPlanta.Nome)
		assert.Equal(t, expectedPlanta.ComecandoDe, actualPlanta.ComecandoDe)
		assert.Equal(t, expectedPlanta.Especie, actualPlanta.Especie)
		assert.Equal(t, expectedPlanta.DataPlantio.Format("2006-01-02"), actualPlanta.DataPlantio.Format("2006-01-02"))
		assert.Equal(t, expectedPlanta.Status, actualPlanta.Status)
		assert.Equal(t, *expectedPlanta.Notas, *actualPlanta.Notas)
		assert.Equal(t, expectedPlanta.GeneticaID, actualPlanta.GeneticaID)
		assert.Equal(t, expectedPlanta.MeioCultivoID, actualPlanta.MeioCultivoID)
		assert.Equal(t, expectedPlanta.AmbienteID, actualPlanta.AmbienteID)
		assert.Equal(t, expectedPlanta.UsuarioID, actualPlanta.UsuarioID)

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
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Planta não encontrada", response["message"])
		assert.Equal(t, "recurso não encontrado", response["details"])

		mockService.AssertExpectations(t)
	})
}

func TestPlantaController_Criar(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		// Preparação
		mockService := new(MockPlantaService)
		controller := NewPlantaController(mockService)

		createDTO := &dto.CreatePlantaDTO{
			Nome:          "Planta Teste",
			ComecandoDe:   "semente",
			Especie:       "Sativa",
			DataPlantio:   time.Now(),
			Status:        "vegetativo",
			Notas:         "Algumas notas.",
			GeneticaID:    1,
			MeioCultivoID: 1,
			AmbienteID:    1,
			UsuarioID:     1,
		}
		expectedPlanta := &dto.PlantaResponseDTO{
			ID:            1,
			Nome:          "Planta Teste",
			ComecandoDe:   "semente",
			Especie:       "Sativa",
			DataPlantio:   &createDTO.DataPlantio,
			Status:        "vegetativo",
			Notas:         &createDTO.Notas,
			GeneticaID:    1,
			MeioCultivoID: 1,
			AmbienteID:    1,
			UsuarioID:     1,
		}

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
		assert.Equal(t, expectedPlanta.ID, actualPlanta.ID)
		assert.Equal(t, expectedPlanta.Nome, actualPlanta.Nome)
		assert.Equal(t, expectedPlanta.ComecandoDe, actualPlanta.ComecandoDe)
		assert.Equal(t, expectedPlanta.Especie, actualPlanta.Especie)
		assert.Equal(t, expectedPlanta.DataPlantio.Format("2006-01-02"), actualPlanta.DataPlantio.Format("2006-01-02"))
		assert.Equal(t, expectedPlanta.Status, actualPlanta.Status)
		assert.Equal(t, *expectedPlanta.Notas, *actualPlanta.Notas)
		assert.Equal(t, expectedPlanta.GeneticaID, actualPlanta.GeneticaID)
		assert.Equal(t, expectedPlanta.MeioCultivoID, actualPlanta.MeioCultivoID)
		assert.Equal(t, expectedPlanta.AmbienteID, actualPlanta.AmbienteID)
		assert.Equal(t, expectedPlanta.UsuarioID, actualPlanta.UsuarioID)

		mockService.AssertExpectations(t)
	})

	t.Run("Error - Invalid Payload", func(t *testing.T) {
		// Preparação
		mockService := new(MockPlantaService)
		controller := NewPlantaController(mockService)

		// Payload inválido (faltando campos obrigatórios)
		invalidPayload := []byte(`{"nome": "Planta Invalida"}`)

		// Execução
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/plantas", bytes.NewBuffer(invalidPayload))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req

		controller.Criar(c)

		// Verificação
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Erro de validação", response["message"])
		assert.Contains(t, response["details"].(map[string]interface{})["ComecandoDe"], "Este campo é obrigatório")
		assert.Contains(t, response["details"].(map[string]interface{})["Especie"], "Este campo é obrigatório")
		assert.Contains(t, response["details"].(map[string]interface{})["DataPlantio"], "Este campo é obrigatório")
		assert.Contains(t, response["details"].(map[string]interface{})["Status"], "Este campo é obrigatório")
		assert.Contains(t, response["details"].(map[string]interface{})["GeneticaID"], "Este campo é obrigatório")
		assert.Contains(t, response["details"].(map[string]interface{})["MeioCultivoID"], "Este campo é obrigatório")
		assert.Contains(t, response["details"].(map[string]interface{})["AmbienteID"], "Este campo é obrigatório")
		assert.Contains(t, response["details"].(map[string]interface{})["UsuarioID"], "Este campo é obrigatório")

		mockService.AssertNotCalled(t, "Criar") // O serviço não deve ser chamado em caso de payload inválido
	})

	t.Run("Error - Service Error", func(t *testing.T) {
		// Preparação
		mockService := new(MockPlantaService)
		controller := NewPlantaController(mockService)

		createDTO := &dto.CreatePlantaDTO{
			Nome:          "Planta Teste",
			ComecandoDe:   "semente",
			Especie:       "Sativa",
			DataPlantio:   time.Now(),
			Status:        "vegetativo",
			Notas:         "Algumas notas.",
			GeneticaID:    1,
			MeioCultivoID: 1,
			AmbienteID:    1,
			UsuarioID:     1,
		}

		mockService.On("Criar", mock.AnythingOfType("*dto.CreatePlantaDTO")).Return((*dto.PlantaResponseDTO)(nil), errors.New("erro interno do serviço"))

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
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Erro interno ao criar planta", response["message"])
		assert.Equal(t, "erro interno do serviço", response["details"])

		mockService.AssertExpectations(t)
	})
}

func TestPlantaController_Atualizar(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		// Preparação
		mockService := new(MockPlantaService)
		controller := NewPlantaController(mockService)

		updateDTO := &dto.UpdatePlantaDTO{
			Nome:          "Sativa Nova",
			ComecandoDe:   "semente",
			Especie:       "Sativa",
			DataPlantio:   time.Now(),
			Status:        "vegetativo",
			Notas:         "Algumas notas atualizadas.",
			GeneticaID:    1,
			MeioCultivoID: 1,
			AmbienteID:    1,
			UsuarioID:     1,
		}
		expectedPlanta := &dto.PlantaResponseDTO{
			ID:            1,
			Nome:          "Sativa Nova",
			ComecandoDe:   "semente",
			Especie:       "Sativa",
			DataPlantio:   &updateDTO.DataPlantio,
			Status:        "vegetativo",
			Notas:         &updateDTO.Notas,
			GeneticaID:    1,
			MeioCultivoID: 1,
			AmbienteID:    1,
			UsuarioID:     1,
		}

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
		assert.Equal(t, expectedPlanta.ID, actualPlanta.ID)
		assert.Equal(t, expectedPlanta.Nome, actualPlanta.Nome)
		assert.Equal(t, expectedPlanta.ComecandoDe, actualPlanta.ComecandoDe)
		assert.Equal(t, expectedPlanta.Especie, actualPlanta.Especie)
		assert.Equal(t, expectedPlanta.DataPlantio.Format("2006-01-02"), actualPlanta.DataPlantio.Format("2006-01-02"))
		assert.Equal(t, expectedPlanta.Status, actualPlanta.Status)
		assert.Equal(t, *expectedPlanta.Notas, *actualPlanta.Notas)
		assert.Equal(t, expectedPlanta.GeneticaID, actualPlanta.GeneticaID)
		assert.Equal(t, expectedPlanta.MeioCultivoID, actualPlanta.MeioCultivoID)
		assert.Equal(t, expectedPlanta.AmbienteID, actualPlanta.AmbienteID)
		assert.Equal(t, expectedPlanta.UsuarioID, actualPlanta.UsuarioID)

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
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Planta não encontrada", response["message"])
		assert.Equal(t, "recurso não encontrado", response["details"])

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
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Requisição inválida", response["message"])
		assert.Equal(t, "entrada inválida", response["details"])
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
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Planta não encontrada", response["message"])
		assert.Equal(t, "recurso não encontrado", response["details"])

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
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "ID inválido", response["message"])
		assert.Equal(t, "entrada inválida", response["details"])
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
			Tipo:     "observacao",
			Titulo:   "Primeira Observação",
			Conteudo: "A planta está crescendo bem.",
		}

		mockService.On("RegistrarFato", plantID, entity.RegistroTipo(fatoDTO.Tipo), fatoDTO.Titulo, fatoDTO.Conteudo).Return(nil)

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
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "ID inválido", response["message"])
		assert.Equal(t, "entrada inválida", response["details"])
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
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Requisição inválida", response["message"])
		assert.Equal(t, "entrada inválida", response["details"])
	})

	t.Run("Plant Not Found", func(t *testing.T) {
		// Preparação
		mockService := new(MockPlantaService)
		controller := NewPlantaController(mockService)

		plantID := uint(1)
		fatoDTO := &dto.RegistrarFatoDTO{
			Tipo:     "observacao",
			Titulo:   "Primeira Observação",
			Conteudo: "A planta está crescendo bem.",
		}

		mockService.On("RegistrarFato", plantID, entity.RegistroTipo(fatoDTO.Tipo), fatoDTO.Titulo, fatoDTO.Conteudo).Return(gorm.ErrRecordNotFound)

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
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Planta não encontrada", response["message"])
		assert.Equal(t, "recurso não encontrado", response["details"])

		mockService.AssertExpectations(t)
	})

	t.Run("Service Error", func(t *testing.T) {
		// Preparação
		mockService := new(MockPlantaService)
		controller := NewPlantaController(mockService)

		plantID := uint(1)
		fatoDTO := &dto.RegistrarFatoDTO{
			Tipo:     "observacao",
			Titulo:   "Primeira Observação",
			Conteudo: "A planta está crescendo bem.",
		}

		mockService.On("RegistrarFato", plantID, entity.RegistroTipo(fatoDTO.Tipo), fatoDTO.Titulo, fatoDTO.Conteudo).Return(errors.New("erro interno do serviço"))

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
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Erro interno ao registrar fato", response["message"])
		assert.Equal(t, "erro interno do serviço", response["details"])

		mockService.AssertExpectations(t)
	})
}
