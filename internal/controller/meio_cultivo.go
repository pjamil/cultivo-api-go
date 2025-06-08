package controller

import (
	"net/http"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/service"

	"github.com/gin-gonic/gin"
)

type MeioCultivoController struct {
	service service.MeioCultivoService
}

func NewMeioCultivoController(service service.MeioCultivoService) *MeioCultivoController {
	return &MeioCultivoController{service}
}

type CreateMeioCultivoDTO struct {
	Tipo      string `json:"tipo" binding:"required"`
	Descricao string `json:"descricao"`
}

func (ctrl *MeioCultivoController) Create(c *gin.Context) {
	var dto CreateMeioCultivoDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	meioCultivo := models.MeioCultivo{
		Tipo:      dto.Tipo,
		Descricao: dto.Descricao,
	}

	if err := ctrl.service.CreateMeioCultivo(&meioCultivo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar meio de cultivo"})
		return
	}

	c.JSON(http.StatusCreated, meioCultivo)
}
