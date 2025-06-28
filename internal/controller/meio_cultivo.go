package controller

import (
	"net/http"
	"strconv"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type MeioCultivoController struct {
	servico service.MeioCultivoService
}

func NewMeioCultivoController(servico service.MeioCultivoService) *MeioCultivoController {
	return &MeioCultivoController{servico}
}

func (c *MeioCultivoController) Criar(ctx *gin.Context) {
	var dto dto.CreateMeioCultivoDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		logrus.WithError(err).Error("Payload da requisição inválido para criar meio de cultivo")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	meioCultivo, err := c.servico.Criar(&dto)
	if err != nil {
		logrus.WithError(err).Error("Erro ao criar meio de cultivo")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar meio de cultivo"})
		return
	}

	ctx.JSON(http.StatusCreated, meioCultivo)
}

func (c *MeioCultivoController) Listar(ctx *gin.Context) {
	meioCultivos, err := c.servico.ListarTodos()
	if err != nil {
		logrus.WithError(err).Error("Erro ao listar meios de cultivo")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao recuperar meios de cultivo"})
		return
	}

	ctx.JSON(http.StatusOK, meioCultivos)
}

func (c *MeioCultivoController) BuscarPorID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		logrus.WithError(err).Error("ID inválido para buscar meio de cultivo por ID")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID Inválido"})
		return
	}
	usuario, err := c.servico.BuscarPorID(uint(id))
	if err != nil {
		logrus.WithError(err).Error("Meio de cultivo não encontrado ao buscar por ID")
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	}
	ctx.JSON(http.StatusOK, usuario)
}
