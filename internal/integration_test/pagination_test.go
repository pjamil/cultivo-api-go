package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/utils"
	tu "gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/utils/test_utils"
	"github.com/stretchr/testify/assert"
)

func TestListagemComPaginacaoPadrao(t *testing.T) {
	router := GetTestRouter()
	db := GetTestDB().DB

	// Limpar o banco de dados antes do teste
	db.Exec("DELETE FROM plantas")
	db.Exec("DELETE FROM usuarios")

	// Criar um usuário para associar às plantas
	password := "userpassword"
	hashedPassword, err := utils.HashPassword(password)
	assert.NoError(t, err)
	user := models.Usuario{
		Nome:      "Paginacao User",
		Email:     "paginacao@example.com",
		SenhaHash: hashedPassword,
	}
	err = db.Create(&user).Error
	assert.NoError(t, err)

	// Criar uma genética para associar às plantas
	genetica := models.Genetica{
		Nome:         "Genetica Paginacao",
		TipoGenetica: "Indica",
		TipoEspecie:  "Indica",
		TempoFloracao: 50,
		Origem:       "Teste",
	}
	assert.NoError(t, db.Create(&genetica).Error)

	// Criar um meio de cultivo para associar às plantas
	meioCultivo := models.MeioCultivo{
		Tipo: "solo",
	}
	assert.NoError(t, db.Create(&meioCultivo).Error)

	// Criar um ambiente para associar às plantas
	ambiente := models.Ambiente{
		Nome:        "Ambiente Paginacao",
		Descricao:   "Descrição do ambiente de paginação",
		Tipo:        "interno",
		Comprimento: 10,
		Altura:      10,
		Largura:     10,
	}
	assert.NoError(t, db.Create(&ambiente).Error)

	// Criar 25 plantas para testar a paginação
	for i := 1; i <= 25; i++ {
		planta := models.Planta{
			Nome:        fmt.Sprintf("Planta %d", i),
			Especie:     "Especie Teste",
			DataPlantio: tu.TimePtr(time.Now()),
			Status:      "vegetativo",
			UsuarioID:   user.ID,
			GeneticaID:  genetica.ID, // Usar o ID da genética criada
			MeioCultivoID: meioCultivo.ID, // Usar o ID do meio de cultivo criado
			AmbienteID:  ambiente.ID, // Usar o ID do ambiente criado
		}
		err = db.Create(&planta).Error
		assert.NoError(t, err)
	}

	// Fazer login para obter um token
	loginPayload := dto.LoginPayload{
		Email:    user.Email,
		Password: password,
	}
	jsonPayload, _ := json.Marshal(loginPayload)
	wLogin := httptest.NewRecorder()
	reqLogin, _ := http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(jsonPayload))
	reqLogin.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(wLogin, reqLogin)
	assert.Equal(t, http.StatusOK, wLogin.Code)
	var loginResponse map[string]string
	json.Unmarshal(wLogin.Body.Bytes(), &loginResponse)
	token := loginResponse["token"]

	// Fazer requisição GET para /api/v1/plantas sem parâmetros de paginação
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/plantas", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	// Verificar a resposta
	assert.Equal(t, http.StatusOK, w.Code)

	var paginatedResponse dto.PaginatedResponse
	err = json.Unmarshal(w.Body.Bytes(), &paginatedResponse)
	assert.NoError(t, err)

	assert.Equal(t, int64(25), paginatedResponse.Total) // Total de registros
	assert.Equal(t, 1, paginatedResponse.Page)          // Página padrão
	assert.Equal(t, 10, paginatedResponse.Limit)         // Limite padrão
	assert.Len(t, paginatedResponse.Data.([]interface{}), 10) // Deve retornar 10 itens na primeira página
}

func TestListagemComPaginacaoCustomizada(t *testing.T) {
	router := GetTestRouter()
	db := GetTestDB().DB

	// Limpar o banco de dados antes do teste
	db.Exec("DELETE FROM plantas")
	db.Exec("DELETE FROM usuarios")

	// Criar um usuário para associar às plantas
	password := "userpassword"
	hashedPassword, err := utils.HashPassword(password)
	assert.NoError(t, err)
	user := models.Usuario{
		Nome:      "Paginacao User 2",
		Email:     "paginacao2@example.com",
		SenhaHash: hashedPassword,
	}
	err = db.Create(&user).Error
	assert.NoError(t, err)

	// Criar uma genética para associar às plantas
	genetica := models.Genetica{
		Nome:         "Genetica Paginacao Custom",
		TipoGenetica: "Sativa",
		TipoEspecie:  "Sativa",
		TempoFloracao: 70,
		Origem:       "Teste Custom",
	}
	assert.NoError(t, db.Create(&genetica).Error)

	// Criar um meio de cultivo para associar às plantas
	meioCultivo := models.MeioCultivo{
		Tipo: "coco",
	}
	assert.NoError(t, db.Create(&meioCultivo).Error)

	// Criar um ambiente para associar às plantas
	ambiente := models.Ambiente{
		Nome:        "Ambiente Paginacao Custom",
		Descricao:   "Descrição do ambiente de paginação customizada",
		Tipo:        "externo",
		Comprimento: 20,
		Altura:      15,
		Largura:     12,
	}
	assert.NoError(t, db.Create(&ambiente).Error)

	// Criar 25 plantas para testar a paginação
	for i := 1; i <= 25; i++ {
		planta := models.Planta{
			Nome:        fmt.Sprintf("Planta Custom %d", i),
			Especie:     "Especie Teste Custom",
			DataPlantio: tu.TimePtr(time.Now()),
			Status:      "vegetativo",
			UsuarioID:   user.ID,
			GeneticaID:  genetica.ID, // Usar o ID da genética criada
			MeioCultivoID: meioCultivo.ID, // Usar o ID do meio de cultivo criado
			AmbienteID:  ambiente.ID, // Usar o ID do ambiente criado
		}
		err = db.Create(&planta).Error
		assert.NoError(t, err)
	}

	// Fazer login para obter um token
	loginPayload := dto.LoginPayload{
		Email:    user.Email,
		Password: password,
	}
	jsonPayload, _ := json.Marshal(loginPayload)
	wLogin := httptest.NewRecorder()
	reqLogin, _ := http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(jsonPayload))
	reqLogin.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(wLogin, reqLogin)
	assert.Equal(t, http.StatusOK, wLogin.Code)
	var loginResponse map[string]string
	json.Unmarshal(wLogin.Body.Bytes(), &loginResponse)
	token := loginResponse["token"]

	// Cenário 1: Página 2, Limite 5
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/plantas?page=2&limit=5", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var paginatedResponse dto.PaginatedResponse
	err = json.Unmarshal(w.Body.Bytes(), &paginatedResponse)
	assert.NoError(t, err)

	assert.Equal(t, int64(25), paginatedResponse.Total)
	assert.Equal(t, 2, paginatedResponse.Page)
	assert.Equal(t, 5, paginatedResponse.Limit)
	assert.Len(t, paginatedResponse.Data.([]interface{}), 5)

	// Cenário 2: Última página (Página 5, Limite 5)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/plantas?page=5&limit=5", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	err = json.Unmarshal(w.Body.Bytes(), &paginatedResponse)
	assert.NoError(t, err)

	assert.Equal(t, int64(25), paginatedResponse.Total)
	assert.Equal(t, 5, paginatedResponse.Page)
	assert.Equal(t, 5, paginatedResponse.Limit)
	assert.Len(t, paginatedResponse.Data.([]interface{}), 5)

	// Cenário 3: Página que não existe (Página 10, Limite 5) - deve retornar vazio
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/plantas?page=10&limit=5", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	err = json.Unmarshal(w.Body.Bytes(), &paginatedResponse)
	assert.NoError(t, err)

	assert.Equal(t, int64(25), paginatedResponse.Total)
	assert.Equal(t, 10, paginatedResponse.Page)
	assert.Equal(t, 5, paginatedResponse.Limit)
	assert.Len(t, paginatedResponse.Data.([]interface{}), 0)
}

func TestListagemComPaginacaoVazia(t *testing.T) {
	router := GetTestRouter()
	db := GetTestDB().DB

	// Limpar o banco de dados para garantir que não há registros
	db.Exec("DELETE FROM plantas")
	db.Exec("DELETE FROM usuarios")

	// Criar um usuário para associar às plantas (mesmo que não haja plantas)
	password := "userpassword"
	hashedPassword, err := utils.HashPassword(password)
	assert.NoError(t, err)
	user := models.Usuario{
		Nome:      "Empty Pagination User",
		Email:     "empty@example.com",
		SenhaHash: hashedPassword,
	}
	err = db.Create(&user).Error
	assert.NoError(t, err)

	// Fazer login para obter um token
	loginPayload := dto.LoginPayload{
		Email:    user.Email,
		Password: password,
	}
	jsonPayload, _ := json.Marshal(loginPayload)
	wLogin := httptest.NewRecorder()
	reqLogin, _ := http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(jsonPayload))
	reqLogin.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(wLogin, reqLogin)
	assert.Equal(t, http.StatusOK, wLogin.Code)
	var loginResponse map[string]string
	json.Unmarshal(wLogin.Body.Bytes(), &loginResponse)
	token := loginResponse["token"]

	// Fazer requisição GET para /api/v1/plantas
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/plantas", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	// Verificar a resposta
	assert.Equal(t, http.StatusOK, w.Code)

	var paginatedResponse dto.PaginatedResponse
	err = json.Unmarshal(w.Body.Bytes(), &paginatedResponse)
	assert.NoError(t, err)

	assert.Equal(t, int64(0), paginatedResponse.Total) // Total deve ser 0
	assert.Equal(t, 1, paginatedResponse.Page)          // Página padrão
	assert.Equal(t, 10, paginatedResponse.Limit)         // Limite padrão
	assert.Len(t, paginatedResponse.Data.([]interface{}), 0) // Data deve ser um array vazio
}
