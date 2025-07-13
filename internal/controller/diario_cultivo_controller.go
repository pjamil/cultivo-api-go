package controller

import (
	"errors"
	"net/http"
	"strconv"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/service"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type DiarioCultivoController struct {
	diarioCultivoServico service.DiarioCultivoService
}

func NewDiarioCultivoController(diarioCultivoServico service.DiarioCultivoService) *DiarioCultivoController {
	return &DiarioCultivoController{diarioCultivoServico: diarioCultivoServico}
}

// CreateDiarioCultivo godoc
// @Summary      Cria um novo diário de cultivo
// @Description  Adiciona um novo diário de cultivo ao sistema
// @Tags         diarios-cultivo
// @Accept       json
// @Produce      json
// @Param        diarioCultivo  body      dto.CreateDiarioCultivoDTO  true  "Dados para criação do diário de cultivo"
// @Success      201    {object}  dto.DiarioCultivoResponseDTO
// @Failure      400    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Router       /api/v1/diarios-cultivo [post]
func (c *DiarioCultivoController) Create(ctx *gin.Context) {
	var createDto dto.CreateDiarioCultivoDTO
	if err := ctx.ShouldBindJSON(&createDto); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errMsgs := make(map[string]string)
			for _, fe := range ve {
				errMsgs[fe.Field()] = utils.GetErrorMsg(fe)
			}
			utils.RespondWithError(ctx, http.StatusBadRequest, "Erro de validação", errMsgs)
			return
		}
		utils.RespondWithError(ctx, http.StatusBadRequest, "Requisição inválida", utils.ErrInvalidInput.Error())
		return
	}

	if diarioCultivoCriado, err := c.diarioCultivoServico.Create(&createDto); err != nil {
		logrus.WithError(err).Error("Erro ao criar diário de cultivo")
		utils.RespondWithError(ctx, http.StatusInternalServerError, "Erro interno ao criar diário de cultivo", err.Error())
		return
	} else {
		utils.RespondWithJSON(ctx, http.StatusCreated, diarioCultivoCriado)
	}
}

// ListDiariosCultivo godoc
// @Summary      Lista todos os diários de cultivo com paginação
// @Description  Retorna uma lista paginada de todos os diários de cultivo
// @Tags         diarios-cultivo
// @Accept       json
// @Produce      json
// @Param        page   query     int  false  "Número da página (padrão: 1)"
// @Param        limit  query     int  false  "Limite de itens por página (padrão: 10)"
// @Success      200  {object}  dto.PaginatedResponse{data=[]dto.DiarioCultivoResponseDTO}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /api/v1/diarios-cultivo [get]
func (c *DiarioCultivoController) List(ctx *gin.Context) {
	var pagination dto.PaginationParams
	if err := ctx.ShouldBindQuery(&pagination); err != nil {
		logrus.WithError(err).Error("Parâmetros de paginação inválidos para listar diários de cultivo")
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errMsgs := make(map[string]string)
			for _, fe := range ve {
				errMsgs[fe.Field()] = utils.GetErrorMsg(fe)
			}
			utils.RespondWithError(ctx, http.StatusBadRequest, "Erro de validação", errMsgs)
			return
		}
		utils.RespondWithError(ctx, http.StatusBadRequest, "Requisição inválida", utils.ErrInvalidInput.Error())
		return
	}

	paginatedResponse, err := c.diarioCultivoServico.GetAll(pagination.Page, pagination.Limit)
	if err != nil {
		logrus.WithError(err).Error("Erro ao listar diários de cultivo com paginação")
		utils.RespondWithError(ctx, http.StatusInternalServerError, "Erro interno ao listar diários de cultivo", err.Error())
		return
	}

	utils.RespondWithJSON(ctx, http.StatusOK, paginatedResponse)
}

// GetDiarioCultivoByID godoc
// @Summary      Busca um diário de cultivo por ID
// @Description  Retorna os detalhes de um diário de cultivo específico
// @ID           get-diario-cultivo-by-id
// @Tags         diarios-cultivo
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID do Diário de Cultivo"
// @Success      200  {object}  dto.DiarioCultivoResponseDTO
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /api/v1/diarios-cultivo/{id} [get]
func (c *DiarioCultivoController) GetByID(ctx *gin.Context) {
	logrus.WithFields(logrus.Fields{
		"id": ctx.Param("id"),
	}).Info("Buscando diário de cultivo por ID")
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		logrus.WithError(err).Error("Erro ao converter ID do diário de cultivo")
		utils.RespondWithError(ctx, http.StatusBadRequest, "ID inválido", utils.ErrInvalidInput.Error())
		return
	}

	diarioCultivo, err := c.diarioCultivoServico.GetByID(uint(id))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.WithError(err).Error("Diário de cultivo não encontrado ao buscar por ID")
		utils.RespondWithError(ctx, http.StatusNotFound, "Diário de cultivo não encontrado", utils.ErrNotFound.Error())
		return
	}
	utils.RespondWithJSON(ctx, http.StatusOK, diarioCultivo)
}

// UpdateDiarioCultivo godoc
// @Summary      Atualiza um diário de cultivo
// @Description  Atualiza um diário de cultivo existente
// @Tags         diarios-cultivo
// @Accept       json
// @Produce      json
// @Param        id      path      int            true  "ID do Diário de Cultivo"
// @Param        diarioCultivo  body      dto.UpdateDiarioCultivoDTO  true  "Dados para atualização do diário de cultivo"
// @Success      200    {object}  dto.DiarioCultivoResponseDTO
// @Failure      400    {object}  map[string]interface{}
// @Failure      404    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Router       /api/v1/diarios-cultivo/{id} [put]
func (c *DiarioCultivoController) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		logrus.WithError(err).Error("Erro ao converter ID do diário de cultivo para atualização")
		utils.RespondWithError(ctx, http.StatusBadRequest, "ID inválido", utils.ErrInvalidInput.Error())
		return
	}
	var updateDto dto.UpdateDiarioCultivoDTO
	if err := ctx.ShouldBindJSON(&updateDto); err != nil {
		logrus.WithError(err).Error("Payload da requisição inválido para atualização de diário de cultivo")
		utils.RespondWithError(ctx, http.StatusBadRequest, "Requisição inválida", utils.ErrInvalidInput.Error())
		return
	}

	diarioCultivoAtualizado, err := c.diarioCultivoServico.Update(uint(id), &updateDto)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.WithError(err).Error("Diário de cultivo não encontrado para atualização")
		utils.RespondWithError(ctx, http.StatusNotFound, "Diário de cultivo não encontrado", utils.ErrNotFound.Error())
		return
	}
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"operation": "update_diario_cultivo",
			"diario_cultivo_id":  id,
			"error":     err,
		}).Error("Erro ao atualizar diário de cultivo")
		utils.RespondWithError(ctx, http.StatusInternalServerError, "Erro interno ao atualizar diário de cultivo", err.Error())
		return
	}

	utils.RespondWithJSON(ctx, http.StatusOK, diarioCultivoAtualizado)
}

// DeleteDiarioCultivo godoc
// @Summary      Deleta um diário de cultivo
// @Description  Deleta um diário de cultivo existente
// @Tags         diarios-cultivo
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID do Diário de Cultivo"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /api/v1/diarios-cultivo/{id} [delete]
func (c *DiarioCultivoController) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		logrus.WithError(err).Error("Erro ao converter ID do diário de cultivo para deleção")
		utils.RespondWithError(ctx, http.StatusBadRequest, "ID inválido", utils.ErrInvalidInput.Error())
		return
	}
	if err := c.diarioCultivoServico.Delete(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.WithError(err).Error("Diário de cultivo não encontrado para deleção")
			utils.RespondWithError(ctx, http.StatusNotFound, "Diário de cultivo não encontrado", utils.ErrNotFound.Error())
			return
		}
		logrus.WithError(err).Error("Erro ao deletar diário de cultivo")
		utils.RespondWithError(ctx, http.StatusInternalServerError, "Erro interno ao deletar diário de cultivo", err.Error())
		return
	}
	logrus.WithFields(logrus.Fields{
		"operation": "delete_diario_cultivo",
		"diario_cultivo_id":  id,
	}).Info("Diário de cultivo deletado com sucesso")
	utils.RespondWithJSON(ctx, http.StatusOK, gin.H{"message": "Diário de cultivo deletado com sucesso"})
}
