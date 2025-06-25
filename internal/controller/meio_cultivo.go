package controller

import (
	"net/http"
	"strconv"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/service"

	"github.com/gin-gonic/gin"
)

type MeioCultivoController struct {
	service service.MeioCultivoService
}

func NewMeioCultivoController(service service.MeioCultivoService) *MeioCultivoController {
	return &MeioCultivoController{service}
}

func (c *MeioCultivoController) Create(ctx *gin.Context) {
	var dto dto.CreateMeioCultivoDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	meioCultivo, err := c.service.CreateMeioCultivo(&dto)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar meio de cultivo"})
		return
	}

	ctx.JSON(http.StatusCreated, meioCultivo)
}

func (c *MeioCultivoController) GetAll(ctx *gin.Context) {
	meioCultivos, err := c.service.GetAllMeioCultivos()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao recuperar meios de cultivo"})
		return
	}

	ctx.JSON(http.StatusOK, meioCultivos)
}

func (c *MeioCultivoController) GetByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID Inválido"})
		return
	}
	usuario, err := c.service.GetByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	}
	ctx.JSON(http.StatusOK, usuario)
}
