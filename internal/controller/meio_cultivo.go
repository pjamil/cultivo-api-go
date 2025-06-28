package controller

import (
	"errors"
	"net/http"
	"strconv"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
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

// Atualizar godoc
// @Summary      Atualiza um meio de cultivo
// @Description  Atualiza um meio de cultivo existente com os dados fornecidos
// @Tags         meio_cultivo
// @Accept       json
// @Produce      json
// @Param        id        path      int                     true  "ID do Meio de Cultivo"
// @Param        meioCultivo  body      dto.UpdateMeioCultivoDTO  true  "Dados do Meio de Cultivo para atualização"
// @Success      200       {object}  models.MeioCultivo
// @Failure      400       {object}  map[string]string
// @Failure      404       {object}  map[string]string
// @Failure      500       {object}  map[string]string
// @Router       /meios-cultivos/{id} [put]
func (c *MeioCultivoController) Atualizar(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		logrus.WithError(err).Error("ID inválido para atualização de meio de cultivo")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var updateDto dto.UpdateMeioCultivoDTO
	if err := ctx.ShouldBindJSON(&updateDto); err != nil {
		logrus.WithError(err).Error("Payload da requisição inválido para atualização de meio de cultivo")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	meioCultivoAtualizado, err := c.servico.Atualizar(uint(id), &updateDto)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.WithError(err).Error("Meio de cultivo não encontrado para atualização")
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Meio de cultivo não encontrado"})
		return
	}
	if err != nil {
		logrus.WithError(err).Error("Erro ao atualizar meio de cultivo")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar meio de cultivo"})
		return
	}

	ctx.JSON(http.StatusOK, meioCultivoAtualizado)
}
