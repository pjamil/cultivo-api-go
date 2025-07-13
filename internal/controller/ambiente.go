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

type AmbienteController struct {
	servico service.AmbienteService
}

func NewAmbienteController(servico service.AmbienteService) *AmbienteController {
	return &AmbienteController{servico}
}

// Criar godoc
// @Summary      Cria um novo ambiente
// @Description  Cria um novo ambiente com os dados fornecidos
// @Tags         ambiente
// @Accept       json
// @Produce      json
// @Param        ambiente  body      dto.CreateAmbienteDTO  true  "Dados do Ambiente"
// @Success      201      {object}  dto.AmbienteResponseDTO
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /api/v1/ambientes [post]
func (c *AmbienteController) Criar(ctx *gin.Context) {
	var dto dto.CreateAmbienteDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		logrus.WithError(err).Error("Payload da requisição inválido para criar ambiente")
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errMsgs := make(map[string]string)
			for _, fe := range ve {
				errMsgs[fe.Field()] = utils.GetErrorMsg(fe)
			}
			utils.RespondWithError(ctx, http.StatusBadRequest, "Erro de validação", errMsgs)
			return
		}
		utils.RespondWithError(ctx, http.StatusBadRequest, "Requisição inválida", err.Error())
		return
	}
	ambiente, err := c.servico.Criar(&dto)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"operation": "create_ambiente",
			"error":     err,
		}).Error("Erro ao criar ambiente")
		utils.RespondWithError(ctx, http.StatusInternalServerError, "Erro interno ao criar ambiente", nil)
		return
	}
	utils.RespondWithJSON(ctx, http.StatusCreated, ambiente)
}

// Listar godoc
// @Summary      Lista todos os ambientes com paginação
// @Description  Retorna uma lista paginada de todos os ambientes cadastrados
// @Tags         ambiente
// @Produce      json
// @Param        page   query     int  false  "Número da página (padrão: 1)"
// @Param        limit  query     int  false  "Limite de itens por página (padrão: 10)"
// @Success      200  {object}  dto.PaginatedResponse{data=[]dto.AmbienteResponseDTO}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/ambientes [get]
func (c *AmbienteController) Listar(ctx *gin.Context) {
	var pagination dto.PaginationParams
	if err := ctx.ShouldBindQuery(&pagination); err != nil {
		logrus.WithError(err).Error("Parâmetros de paginação inválidos para listar ambientes")
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errMsgs := make(map[string]string)
			for _, fe := range ve {
				errMsgs[fe.Field()] = utils.GetErrorMsg(fe)
			}
			utils.RespondWithError(ctx, http.StatusBadRequest, "Erro de validação", errMsgs)
			return
		}
		utils.RespondWithError(ctx, http.StatusBadRequest, "Requisição inválida", err.Error())
		return
	}

	paginatedResponse, err := c.servico.ListarTodos(pagination.Page, pagination.Limit)
	if err != nil {
		logrus.WithError(err).Error("Erro ao listar ambientes com paginação")
		utils.RespondWithError(ctx, http.StatusInternalServerError, "Erro interno ao listar ambientes", err.Error())
		return
	}

	utils.RespondWithJSON(ctx, http.StatusOK, paginatedResponse)
}

// BuscarPorID godoc
// @Summary      Busca um ambiente por ID
// @Description  Retorna os detalhes de um ambiente específico
// @Tags         ambiente
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID do Ambiente"
// @Success      200  {object}  dto.AmbienteResponseDTO
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/ambientes/{id} [get]
func (c *AmbienteController) BuscarPorID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		logrus.WithError(err).Error("ID inválido para buscar ambiente por ID")
		utils.RespondWithError(ctx, http.StatusBadRequest, "ID inválido", nil)
		return
	}
	ambiente, err := c.servico.BuscarPorID(uint(id))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.WithError(err).Error("Ambiente não encontrado ao buscar por ID")
		utils.RespondWithError(ctx, http.StatusNotFound, "Ambiente não encontrado", nil)
		return
	}
	if err != nil {
		logrus.WithError(err).Error("Erro ao buscar ambiente por ID")
		utils.RespondWithError(ctx, http.StatusInternalServerError, "Erro interno ao buscar ambiente", nil)
		return
	}

	utils.RespondWithJSON(ctx, http.StatusOK, ambiente)
}

// Atualizar godoc
// @Summary      Atualiza um ambiente
// @Description  Atualiza um ambiente existente com os dados fornecidos
// @Tags         ambiente
// @Accept       json
// @Produce      json
// @Param        id        path      int                  true  "ID do Ambiente"
// @Param        ambiente  body      dto.UpdateAmbienteDTO  true  "Dados do Ambiente para atualização"
// @Success      200       {object}  dto.AmbienteResponseDTO
// @Failure      400       {object}  map[string]string
// @Failure      404       {object}  map[string]string
// @Failure      500       {object}  map[string]string
// @Router       /api/v1/ambientes/{id} [put]
func (c *AmbienteController) Atualizar(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		logrus.WithError(err).Error("ID inválido para atualização de ambiente")
		utils.RespondWithError(ctx, http.StatusBadRequest, "ID inválido", nil)
		return
	}

	var updateDto dto.UpdateAmbienteDTO
	if err := ctx.ShouldBindJSON(&updateDto); err != nil {
		logrus.WithError(err).Error("Payload da requisição inválido para atualização de ambiente")
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errMsgs := make(map[string]string)
			for _, fe := range ve {
				errMsgs[fe.Field()] = utils.GetErrorMsg(fe)
			}
			utils.RespondWithError(ctx, http.StatusBadRequest, "Erro de validação", errMsgs)
			return
		}
		utils.RespondWithError(ctx, http.StatusBadRequest, "Requisição inválida", err.Error())
		return
	}

	ambienteAtualizado, err := c.servico.Atualizar(uint(id), &updateDto)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.WithError(err).Error("Ambiente não encontrado para atualização")
		utils.RespondWithError(ctx, http.StatusNotFound, "Ambiente não encontrado", nil)
		return
	}
	if err != nil {
		logrus.WithError(err).Error("Erro ao atualizar ambiente")
		utils.RespondWithError(ctx, http.StatusInternalServerError, "Erro interno ao atualizar ambiente", nil)
		return
	}

	utils.RespondWithJSON(ctx, http.StatusOK, ambienteAtualizado)
}

// Deletar godoc
// @Summary      Deleta um ambiente
// @Description  Deleta um ambiente existente
// @Tags         ambiente
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID do Ambiente"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/ambientes/{id} [delete]
func (c *AmbienteController) Deletar(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		logrus.WithError(err).Error("ID inválido para deletar ambiente")
		utils.RespondWithError(ctx, http.StatusBadRequest, "ID inválido", nil)
		return
	}
	if err := c.servico.Deletar(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.WithError(err).Error("Ambiente não encontrado para deleção")
			utils.RespondWithError(ctx, http.StatusNotFound, "Ambiente não encontrado", nil)
			return
		}
		logrus.WithError(err).Error("Erro ao deletar ambiente")
		utils.RespondWithError(ctx, http.StatusInternalServerError, "Erro interno ao deletar ambiente", nil)
		return
	}
	utils.RespondWithJSON(ctx, http.StatusOK, gin.H{"message": "Ambiente deletado com sucesso"})
}