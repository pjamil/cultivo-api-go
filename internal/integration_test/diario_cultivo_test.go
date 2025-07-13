package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	
	test_utils "gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/utils/test_utils"
)

func TestDiarioCultivoCRUD(t *testing.T) {
	router := GetTestRouter()
	db := GetTestDB().DB

	// Limpar o banco de dados antes de cada teste
	LimparBancoDeDados(db)

	// Criar um usuário para associar ao diário de cultivo
	usuario := models.Usuario{
		Nome:      "Teste Diario",
		Email:     fmt.Sprintf("teste_diario_%d@example.com", time.Now().UnixNano()),
		SenhaHash: "senha_hash",
		Preferencias: json.RawMessage("{}"),
	}
	assert.NoError(t, db.Create(&usuario).Error)

	// Criar um ambiente para associar ao diário de cultivo
	ambiente := models.Ambiente{
		Nome:        "Ambiente Teste Diario",
		Descricao:   "Descrição do ambiente",
		Tipo:        "interno",
		Comprimento: 10,
		Altura:      10,
		Largura:     10,
	}
	assert.NoError(t, db.Create(&ambiente).Error)

	// Criar uma genética para associar à planta
	genetica := models.Genetica{
		Nome:         "Genetica Teste Diario",
		TipoGenetica: "Hibrida",
		TipoEspecie:  "Sativa",
		TempoFloracao: 60,
		Origem:       "Teste",
	}
	assert.NoError(t, db.Create(&genetica).Error)

	// Criar um meio de cultivo para associar à planta
	meioCultivo := models.MeioCultivo{
		Tipo: "solo",
	}
	assert.NoError(t, db.Create(&meioCultivo).Error)

	// Criar uma planta para associar ao diário de cultivo
	planta := models.Planta{
		Nome:        "Planta Teste Diario",
		Especie:     "Especie Teste",
		DataPlantio: test_utils.TimePtr(time.Now()),
		Status:      "vegetativo",
		UsuarioID:   usuario.ID,
		GeneticaID:  genetica.ID, // Usar o ID da genética criada
		MeioCultivoID: meioCultivo.ID, // Usar o ID do meio de cultivo criado
		AmbienteID:  ambiente.ID, // Usar o ID do ambiente criado
	}
	assert.NoError(t, db.Create(&planta).Error)

	// 1. Teste de Criação (POST /api/v1/diarios-cultivo)
	createDto := dto.CreateDiarioCultivoDTO{
		Nome:        "Meu Primeiro Cultivo",
		DataInicio:  time.Now(),
		UsuarioID:   usuario.ID,
		PlantasIDs:  []uint{planta.ID},
		AmbientesIDs: []uint{ambiente.ID},
		Privacidade: "privado",
		Tags:        "teste,primeiro",
	}
	jsonBody, _ := json.Marshal(createDto)
	req, _ := http.NewRequest("POST", "/api/v1/diarios-cultivo", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Fazer login para obter um token
	loginPayload := dto.LoginPayload{
		Email:    usuario.Email,
		Password: "senha_hash", // Use a senha real do usuário criado
	}
	jsonLoginPayload, _ := json.Marshal(loginPayload)
	wLogin := httptest.NewRecorder()
	reqLogin, _ := http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(jsonLoginPayload))
	reqLogin.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(wLogin, reqLogin)
	assert.Equal(t, http.StatusOK, wLogin.Code)
	var loginResponse map[string]string
	json.Unmarshal(wLogin.Body.Bytes(), &loginResponse)
	token := loginResponse["token"]

	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var responseDiario dto.DiarioCultivoResponseDTO
	json.Unmarshal(w.Body.Bytes(), &responseDiario)
	assert.Equal(t, createDto.Nome, responseDiario.Nome)
	assert.Equal(t, createDto.UsuarioID, responseDiario.UsuarioID)
	assert.True(t, responseDiario.ID > 0)

	// Salvar o ID do diário de cultivo criado para testes futuros
	diarioCultivoID := responseDiario.ID

	// 2. Teste de Listagem (GET /api/v1/diarios-cultivo)
	req, _ = http.NewRequest("GET", "/api/v1/diarios-cultivo", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var paginatedResponse dto.PaginatedResponse
	json.Unmarshal(w.Body.Bytes(), &paginatedResponse)
	t.Logf("paginatedResponse.Data: %+v", paginatedResponse.Data)
	assert.True(t, paginatedResponse.Total >= 1)
	var listedDiarios []dto.DiarioCultivoResponseDTO
	jsonBytes, _ := json.Marshal(paginatedResponse.Data)
	json.Unmarshal(jsonBytes, &listedDiarios)
	t.Logf("listedDiarios: %+v", listedDiarios)
	assert.True(t, len(listedDiarios) >= 1) // Deve retornar pelo menos 1 item

	// 3. Teste de Busca por ID (GET /api/v1/diarios-cultivo/{id})
	req, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/diarios-cultivo/%d", diarioCultivoID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var fetchedDiario dto.DiarioCultivoResponseDTO
	json.Unmarshal(w.Body.Bytes(), &fetchedDiario)
	assert.Equal(t, diarioCultivoID, fetchedDiario.ID)
	assert.Equal(t, createDto.Nome, fetchedDiario.Nome)

	// 4. Teste de Atualização (PUT /api/v1/diarios-cultivo/{id})
	updateDto := dto.UpdateDiarioCultivoDTO{
		Nome:        "Cultivo Atualizado",
		Privacidade: "publico",
	}
	jsonBody, _ = json.Marshal(updateDto)
	req, _ = http.NewRequest("PUT", fmt.Sprintf("/api/v1/diarios-cultivo/%d", diarioCultivoID), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var updatedDiario dto.DiarioCultivoResponseDTO
	json.Unmarshal(w.Body.Bytes(), &updatedDiario)
	assert.Equal(t, updateDto.Nome, updatedDiario.Nome)
	assert.Equal(t, updateDto.Privacidade, updatedDiario.Privacidade)

	// 5. Teste de Deleção (DELETE /api/v1/diarios-cultivo/{id})
	req, _ = http.NewRequest("DELETE", fmt.Sprintf("/api/v1/diarios-cultivo/%d", diarioCultivoID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var deleteResponse map[string]string
	json.Unmarshal(w.Body.Bytes(), &deleteResponse)
	assert.Equal(t, "Diário de cultivo deletado com sucesso", deleteResponse["message"])

	// Verificar se o diário de cultivo foi realmente deletado
	req, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/diarios-cultivo/%d", diarioCultivoID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDiarioCultivoCreateInvalidInput(t *testing.T) {
	router := GetTestRouter()
	db := GetTestDB().DB
	LimparBancoDeDados(db)

	// Teste de criação com nome vazio
	createDto := dto.CreateDiarioCultivoDTO{
		Nome:       "",
		DataInicio: time.Now(),
		UsuarioID:  1, // ID de usuário inválido para este teste
	}
	jsonBody, _ := json.Marshal(createDto)
	req, _ := http.NewRequest("POST", "/api/v1/diarios-cultivo", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var errorResponse map[string]string
	json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.Contains(t, errorResponse, "Nome")
}

func TestDiarioCultivoUpdateNotFound(t *testing.T) {
	router := GetTestRouter()
	db := GetTestDB().DB
	LimparBancoDeDados(db)

	// Criar um usuário para associar ao diário de cultivo
	usuario := models.Usuario{
		Nome:      "Teste Update",
		Email:     fmt.Sprintf("teste_update_%d@example.com", time.Now().UnixNano()),
		SenhaHash: "senha_hash",
		Preferencias: json.RawMessage("{}"),
	}
	assert.NoError(t, db.Create(&usuario).Error)

	// Fazer login para obter um token
	loginPayload := dto.LoginPayload{
		Email:    usuario.Email,
		Password: "senha_hash",
	}
	jsonLoginPayload, _ := json.Marshal(loginPayload)
	wLogin := httptest.NewRecorder()
	reqLogin, _ := http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(jsonLoginPayload))
	reqLogin.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(wLogin, reqLogin)
	assert.Equal(t, http.StatusOK, wLogin.Code)
	var loginResponse map[string]string
	json.Unmarshal(wLogin.Body.Bytes(), &loginResponse)
	token := loginResponse["token"]

	updateDto := dto.UpdateDiarioCultivoDTO{
		Nome: "Nome Inexistente",
	}
	jsonBody, _ := json.Marshal(updateDto)
	req, _ := http.NewRequest("PUT", "/api/v1/diarios-cultivo/9999", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDiarioCultivoDeleteNotFound(t *testing.T) {
	router := GetTestRouter()
	db := GetTestDB().DB
	LimparBancoDeDados(db)

	// Criar um usuário para associar ao diário de cultivo
	usuario := models.Usuario{
		Nome:      "Teste Delete",
		Email:     fmt.Sprintf("teste_delete_%d@example.com", time.Now().UnixNano()),
		SenhaHash: "senha_hash",
		Preferencias: json.RawMessage("{}"),
	}
	assert.NoError(t, db.Create(&usuario).Error)

	// Fazer login para obter um token
	loginPayload := dto.LoginPayload{
		Email:    usuario.Email,
		Password: "senha_hash",
	}
	jsonLoginPayload, _ := json.Marshal(loginPayload)
	wLogin := httptest.NewRecorder()
	reqLogin, _ := http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(jsonLoginPayload))
	reqLogin.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(wLogin, reqLogin)
	assert.Equal(t, http.StatusOK, wLogin.Code)
	var loginResponse map[string]string
	json.Unmarshal(wLogin.Body.Bytes(), &loginResponse)
	token := loginResponse["token"]

	req, _ := http.NewRequest("DELETE", "/api/v1/diarios-cultivo/9999", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
