package controller

import (
	"net/http"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/service"

	"github.com/gin-gonic/gin"
)

type GeneticaController struct {
	service service.GeneticaService
}

func NewGeneticaController(service service.GeneticaService) *GeneticaController {
	return &GeneticaController{service}
}

type CreateGeneticaDTO struct {
	Nome            string          `json:"nome" binding:"required"`
	Descricao       string          `json:"descricao"`
	TipoGenetica    string          `json:"tipoGenetica" binding:"required"`
	TipoEspecie     string          `json:"tipoEspecie" binding:"required"`
	TempoFloracao   int             `json:"tempoFloracao" binding:"required"`
	Origem          string          `json:"origem" binding:"required"`
	Caracteristicas string          `json:"caracteristicas"`
	Plantas         []models.Planta `json:"plantas,omitempty"`
}

func (ctrl *GeneticaController) Create(c *gin.Context) {
	var dto CreateGeneticaDTO
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
