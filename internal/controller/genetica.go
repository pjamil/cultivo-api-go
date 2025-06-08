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
	Nome      string `json:"nome" binding:"required"`
	Descricao string `json:"descricao"`
}

func (ctrl *GeneticaController) Create(c *gin.Context) {
	var dto CreateGeneticaDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	genetica := models.Genetica{
		Nome:      dto.Nome,
		Descricao: dto.Descricao,
	}

	if err := ctrl.service.CreateGenetica(&genetica); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar genética"})
		return
	}

	c.JSON(http.StatusCreated, genetica)
}
