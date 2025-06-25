package controller

import (
	"errors"
	"net/http"
	"strconv"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AmbienteController struct {
	service service.AmbienteService
}

func NewAmbienteController(service service.AmbienteService) *AmbienteController {
	return &AmbienteController{service}
}

// Handler para criar novo ambiente
func (c *AmbienteController) CreateAmbiente(ctx *gin.Context) {
	var dto dto.CreateAmbienteDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ambiente, err := c.service.CreateAmbiente(&dto)
	if err != nil {
		return
	}
	ctx.JSON(http.StatusCreated, ambiente)
}

// Handler para listar todos os ambientes
func (c *AmbienteController) ListAmbientes(ctx *gin.Context) {
	ambientes, err := c.service.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, ambientes)
}

// @Summary Get ambiente by ID
// @Description Get detailed information about a ambiente
// @Tags ambiente
// @Accept json
// @Produce json
// @Param id path int true "Ambiente ID"
// @Success 200 {object} models.Ambiente
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /ambiente/{id} [get]
func (c *AmbienteController) GetAmbienteByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	ambiente, err := c.service.GetAmbienteByID(uint(id))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Genética não encontrada"})
		return
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar genética"})
		return
	}

	ctx.JSON(http.StatusOK, ambiente)
}
