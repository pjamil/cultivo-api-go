package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/service"
)

// DiarioCultivoHandler encapsula as dependências do handler de diário de cultivo.
type DiarioCultivoHandler struct {
	service service.DiarioCultivoService
}

// NewDiarioCultivoHandler cria uma nova instância do handler.
func NewDiarioCultivoHandler(s service.DiarioCultivoService) *DiarioCultivoHandler {
	return &DiarioCultivoHandler{
		service: s,
	}
}

// CreateDiario lida com a requisição de criação de um novo diário de cultivo.
func (h *DiarioCultivoHandler) Create(c *gin.Context) {
	var input dto.CreateDiarioCultivoDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "corpo da requisição inválido"})
		return
	}

	diario, err := h.service.CreateDiario(input)
	if err != nil {
		// A camada de serviço já logou o erro, então aqui apenas retornamos a resposta HTTP.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao criar o diário"})
		return
	}

	c.JSON(http.StatusCreated, diario)
}

// GetDiarioByID lida com a busca de um diário por ID.
func (h *DiarioCultivoHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	diario, err := h.service.GetDiarioByID(uint(id))
	if err != nil {
		switch err {
		case service.ErrDiarioNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case service.ErrForbidden:
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao buscar o diário"})
		}
		return
	}

	c.JSON(http.StatusOK, diario)
}

// GetAllDiarios lida com a busca de todos os diários.
func (h *DiarioCultivoHandler) GetAllDiarios(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	paginatedResponse, err := h.service.GetAllDiarios(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao buscar os diários"})
		return
	}

	c.JSON(http.StatusOK, paginatedResponse)
}

// UpdateDiario lida com a atualização de um diário existente.
func (h *DiarioCultivoHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var input dto.UpdateDiarioCultivoDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "corpo da requisição inválido"})
		return
	}

	diario, err := h.service.UpdateDiario(uint(id), input)
	if err != nil {
		switch err {
		case service.ErrDiarioNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case service.ErrForbidden:
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao atualizar o diário"})
		}
		return
	}

	c.JSON(http.StatusOK, diario)
}

// DeleteDiario lida com a exclusão de um diário.
func (h *DiarioCultivoHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	err = h.service.DeleteDiario(uint(id))
	if err != nil {
		switch err {
		case service.ErrDiarioNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case service.ErrForbidden:
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao deletar o diário"})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// --- Funções Auxiliares de Mapeamento ---
// As funções de mapeamento toDiarioCultivoResponseDTO e toDiarioCultivoResponseDTOList
// foram removidas pois o serviço já retorna o DTO de resposta diretamente.