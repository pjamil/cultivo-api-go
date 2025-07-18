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

type RegistroDiarioController struct {
	registroDiarioService service.RegistroDiarioService
}

func NewRegistroDiarioController(registroDiarioService service.RegistroDiarioService) *RegistroDiarioController {
	return &RegistroDiarioController{registroDiarioService: registroDiarioService}
}

// CreateRegistroDiario godoc
// @Summary      Cria um novo registro diário para um diário de cultivo
// @Description  Adiciona um novo registro diário a um diário de cultivo específico
// @Tags         registros-diario
// @Accept       json
// @Produce      json
// @Param        diario_id      path      int                         true  "ID do Diário de Cultivo"
// @Param        registroDiario  body      dto.CreateRegistroDiarioDTO  true  "Dados para criação do registro diário"
// @Success      201    {object}  dto.RegistroDiarioResponseDTO
// @Failure      400    {object}  map[string]interface{}
// @Failure      404    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Router       /api/v1/diarios-cultivo/{diario_id}/registros [post]
func (c *RegistroDiarioController) Create(ctx *gin.Context) {
	diarioID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		logrus.WithError(err).Error("Erro ao converter ID do diário de cultivo para criar registro")
		utils.RespondWithError(ctx, http.StatusBadRequest, "ID do diário de cultivo inválido", utils.ErrInvalidInput.Error())
		return
	}

	var createDto dto.CreateRegistroDiarioDTO
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
		logrus.WithError(err).Error("Payload da requisição inválido para criar registro diário")
		utils.RespondWithError(ctx, http.StatusBadRequest, "Requisição inválida", utils.ErrInvalidInput.Error())
		return
	}

	if registroCriado, err := c.registroDiarioService.CriarRegistro(uint(diarioID), &createDto); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.WithError(err).Error("Diário de cultivo não encontrado ao criar registro")
			utils.RespondWithError(ctx, http.StatusNotFound, "Diário de cultivo não encontrado", err.Error())
			return
		}
		logrus.WithError(err).Error("Erro ao criar registro diário")
		utils.RespondWithError(ctx, http.StatusInternalServerError, "Erro interno ao criar registro diário", err.Error())
		return
	} else {
		logrus.Debugf("Registro criado: %+v", registroCriado)
		utils.RespondWithJSON(ctx, http.StatusCreated, registroCriado)
	}
}

// ListRegistrosDiario godoc
// @Summary      Lista todos os registros diários de um diário de cultivo com paginação
// @Description  Retorna uma lista paginada de todos os registros diários de um diário de cultivo específico
// @Tags         registros-diario
// @Accept       json
// @Produce      json
// @Param        diario_id  path      int  true  "ID do Diário de Cultivo"
// @Param        page       query     int  false  "Número da página (padrão: 1)"
// @Param        limit      query     int  false  "Limite de itens por página (padrão: 10)"
// @Success      200  {object}  dto.PaginatedResponse{data=[]dto.RegistroDiarioResponseDTO}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /api/v1/diarios-cultivo/{diario_id}/registros [get]
func (c *RegistroDiarioController) List(ctx *gin.Context) {
	diarioID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		logrus.WithError(err).Error("Erro ao converter ID do diário de cultivo para listar registros")
		utils.RespondWithError(ctx, http.StatusBadRequest, "ID do diário de cultivo inválido", utils.ErrInvalidInput.Error())
		return
	}

	var pagination dto.PaginationParams
	if err := ctx.ShouldBindQuery(&pagination); err != nil {
				logrus.WithError(err).Error("Parâmetros de paginação inválidos para listar registros diários")
		utils.RespondWithError(ctx, http.StatusBadRequest, "Requisição inválida", utils.ErrInvalidInput.Error())
		return
	}

	paginatedResponse, err := c.registroDiarioService.ListarRegistrosPorDiarioID(uint(diarioID), pagination.Page, pagination.Limit)
	if err != nil {
		logrus.WithError(err).Error("Erro ao listar registros diários com paginação")
		utils.RespondWithError(ctx, http.StatusInternalServerError, "Erro interno ao listar registros diários", err.Error())
		return
	}

	utils.RespondWithJSON(ctx, http.StatusOK, paginatedResponse)
}