package controller

import (
	"errors"
	"net/http"
	"strconv"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/service"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type GeneticaController struct {
	service service.GeneticaService
}

func NewGeneticaController(service service.GeneticaService) *GeneticaController {
	return &GeneticaController{service}
}

func (ctrl *GeneticaController) Create(c *gin.Context) {
	var dto dto.CreateGeneticaDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	genetica := models.Genetica{
		Nome:            dto.Nome,
		Descricao:       dto.Descricao,
		TipoGenetica:    dto.TipoGenetica,
		TipoEspecie:     dto.TipoEspecie,
		TempoFloracao:   dto.TempoFloracao,
		Origem:          dto.Origem,
		Caracteristicas: dto.Caracteristicas,
		Plantas:         dto.Plantas,
	}

	if err := ctrl.service.CreateGenetica(&genetica); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar genética"})
		return
	}

	c.JSON(http.StatusCreated, genetica)
}

// GetAll handles GET requests to retrieve all geneticas
func (c *GeneticaController) GetAll(ctx *gin.Context) {
	geneticas, err := c.service.GetAll()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, geneticas)
}

// @Summary Get genetica by ID
// @Description Get detailed information about a genetics strain
// @Tags genetica
// @Accept json
// @Produce json
// @Param id path int true "Genetica ID"
// @Success 200 {object} models.Genetica
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /genetica/{id} [get]
func (c *GeneticaController) GetGeneticaByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	genetics, err := c.service.GetGeneticaByID(uint(id))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Genética não encontrada"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar genética"})
		return
	}
	ctx.JSON(http.StatusOK, genetics)
}
