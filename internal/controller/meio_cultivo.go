package controller

import (
	"errors"
	"net/http"
	"strconv"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/service"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/utils"

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

// Criar godoc
// @Summary      Cria um novo meio de cultivo
// @Description  Cria um novo meio de cultivo com os dados fornecidos
// @Tags         meio_cultivo
// @Accept       json
// @Produce      json
// @Param        meioCultivo  body      dto.CreateMeioCultivoDTO  true  "Dados do Meio de Cultivo"
// @Success      201       {object}  dto.MeioCultivoResponseDTO
// @Failure      400       {object}  map[string]string
// @Failure      500       {object}  map[string]string
// @Router       /api/v1/meios-cultivos [post]
func (c *MeioCultivoController) Criar(ctx *gin.Context) {
	var dto dto.CreateMeioCultivoDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		logrus.WithError(err).Error("Payload da requisição inválido para criar meio de cultivo")
		utils.RespondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	meioCultivo, err := c.servico.Criar(&dto)
	if err != nil {
		logrus.WithError(err).Error("Erro ao criar meio de cultivo")
		utils.RespondWithError(ctx, http.StatusInternalServerError, "Erro ao criar meio de cultivo")
		return
	}

	utils.RespondWithJSON(ctx, http.StatusCreated, meioCultivo)
}

// Listar godoc
// @Summary      Lista todos os meios de cultivo
// @Description  Retorna uma lista de todos os meios de cultivo cadastrados
// @Tags         meio_cultivo
// @Produce      json
// @Success      200  {array}   dto.MeioCultivoResponseDTO
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/meios-cultivos [get]
func (c *MeioCultivoController) Listar(ctx *gin.Context) {
	meioCultivos, err := c.servico.ListarTodos()
	if err != nil {
		logrus.WithError(err).Error("Erro ao listar meios de cultivo")
		utils.RespondWithError(ctx, http.StatusInternalServerError, "Erro ao recuperar meios de cultivo")
		return
	}

	if len(meioCultivos) == 0 {
		utils.RespondWithJSON(ctx, http.StatusOK, gin.H{"message": "Nenhum meio de cultivo encontrado"})
		return
	}

	utils.RespondWithJSON(ctx, http.StatusOK, meioCultivos)
}

// BuscarPorID godoc
// @Summary      Busca um meio de cultivo por ID
// @Description  Retorna os detalhes de um meio de cultivo específico
// @Tags         meio_cultivo
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID do Meio de Cultivo"
// @Success      200  {object}  dto.MeioCultivoResponseDTO
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/meios-cultivos/{id} [get]
func (c *MeioCultivoController) BuscarPorID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		logrus.WithError(err).Error("ID inválido para buscar meio de cultivo por ID")
		utils.RespondWithError(ctx, http.StatusBadRequest, "ID Inválido")
		return
	}
	meioCultivo, err := c.servico.BuscarPorID(uint(id))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.WithError(err).Error("Meio de cultivo não encontrado ao buscar por ID")
		utils.RespondWithError(ctx, http.StatusNotFound, "Meio de cultivo não encontrado")
		return
	}
	utils.RespondWithJSON(ctx, http.StatusOK, meioCultivo)
}

// Atualizar godoc
// @Summary      Atualiza um meio de cultivo
// @Description  Atualiza um meio de cultivo existente com os dados fornecidos
// @Tags         meio_cultivo
// @Accept       json
// @Produce      json
// @Param        id        path      int                     true  "ID do Meio de Cultivo"
// @Param        meioCultivo  body      dto.UpdateMeioCultivoDTO  true  "Dados do Meio de Cultivo para atualização"
// @Success      200       {object}  dto.MeioCultivoResponseDTO
// @Failure      400       {object}  map[string]string
// @Failure      404       {object}  map[string]string
// @Failure      500       {object}  map[string]string
// @Router       /api/v1/meios-cultivos/{id} [put]
func (c *MeioCultivoController) Atualizar(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		logrus.WithError(err).Error("ID inválido para atualização de meio de cultivo")
		utils.RespondWithError(ctx, http.StatusBadRequest, "ID inválido")
		return
	}

	var updateDto dto.UpdateMeioCultivoDTO
	if err := ctx.ShouldBindJSON(&updateDto); err != nil {
		logrus.WithError(err).Error("Payload da requisição inválido para atualização de meio de cultivo")
		utils.RespondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	meioCultivoAtualizado, err := c.servico.Atualizar(uint(id), &updateDto)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.WithError(err).Error("Meio de cultivo não encontrado para atualização")
		utils.RespondWithError(ctx, http.StatusNotFound, "Meio de cultivo não encontrado")
		return
	}
	if err != nil {
		logrus.WithError(err).Error("Erro ao atualizar meio de cultivo")
		utils.RespondWithError(ctx, http.StatusInternalServerError, "Erro ao atualizar meio de cultivo")
		return
	}

	utils.RespondWithJSON(ctx, http.StatusOK, meioCultivoAtualizado)
}

// Deletar godoc
// @Summary      Deleta um meio de cultivo
// @Description  Deleta um meio de cultivo existente
// @Tags         meio_cultivo
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID do Meio de Cultivo"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/meios-cultivos/{id} [delete]
func (c *MeioCultivoController) Deletar(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		logrus.WithError(err).Error("ID inválido para deletar meio de cultivo")
		utils.RespondWithError(ctx, http.StatusBadRequest, "ID inválido")
		return
	}
	if err := c.servico.Deletar(uint(id)); err != nil {
		logrus.WithError(err).Error("Erro ao deletar meio de cultivo")
		utils.RespondWithError(ctx, http.StatusInternalServerError, "Erro ao deletar meio de cultivo")
		return
	}
	utils.RespondWithJSON(ctx, http.StatusOK, gin.H{"message": "Meio de cultivo deletado com sucesso"})
}