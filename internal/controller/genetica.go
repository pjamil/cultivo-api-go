package controller

import (
	"errors"
	"net/http"
	"strconv"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/service"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/utils"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type GeneticaController struct {
	servico service.GeneticaService
}

func NewGeneticaController(servico service.GeneticaService) *GeneticaController {
	return &GeneticaController{servico}
}

// Criar godoc
// @Summary      Cria uma nova genética
// @Description  Cria uma nova genética com os dados fornecidos
// @Tags         genetica
// @Accept       json
// @Produce      json
// @Param        genetica  body      dto.CreateGeneticaDTO  true  "Dados da Genética"
// @Success      201      {object}  dto.GeneticaResponseDTO
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /api/v1/geneticas [post]
func (ctrl *GeneticaController) Criar(c *gin.Context) {
	var dto dto.CreateGeneticaDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		logrus.WithError(err).Error("Payload da requisição inválido para criar genética")
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errMsgs := make(map[string]string)
			for _, fe := range ve {
				errMsgs[fe.Field()] = utils.GetErrorMsg(fe)
			}
			utils.RespondWithError(c, http.StatusBadRequest, "Erro de validação", errMsgs)
			return
		}
		utils.RespondWithError(c, http.StatusBadRequest, "Requisição inválida", err.Error())
		return
	}

	geneticaCriada, err := ctrl.servico.Criar(&dto)
	if err != nil {
		logrus.WithError(err).Error("Erro ao criar genética")
		utils.RespondWithError(c, http.StatusInternalServerError, "Erro interno ao criar genética", nil)
		return
	}

	utils.RespondWithJSON(c, http.StatusCreated, geneticaCriada)
}

// Listar lida com requisições GET para retornar todas as genéticas com paginação
// @Summary      Lista todas as genéticas com paginação
// @Description  Retorna uma lista paginada de todas as genéticas cadastradas
// @Tags         genetica
// @Produce      json
// @Param        page   query     int  false  "Número da página (padrão: 1)"
// @Param        limit  query     int  false  "Limite de itens por página (padrão: 10)"
// @Success      200  {object}  dto.PaginatedResponse{data=[]dto.GeneticaResponseDTO}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/geneticas [get]
func (c *GeneticaController) Listar(ctx *gin.Context) {
	var pagination dto.PaginationParams
	if err := ctx.ShouldBindQuery(&pagination); err != nil {
		logrus.WithError(err).Error("Parâmetros de paginação inválidos para listar genéticas")
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

	paginatedResponse, err := c.servico.ListarTodas(pagination.Page, pagination.Limit)
	if err != nil {
		logrus.WithError(err).Error("Erro ao listar genéticas com paginação")
		utils.RespondWithError(ctx, http.StatusInternalServerError, "Erro interno ao listar genéticas", err.Error())
		return
	}

	utils.RespondWithJSON(ctx, http.StatusOK, paginatedResponse)
}

// BuscarPorID godoc
// @Summary Busca uma genética por ID
// @Description Retorna informações detalhadas de uma genética
// @Tags genetica
// @Accept json
// @Produce json
// @Param id path int true "Genetica ID"
// @Success 200 {object} dto.GeneticaResponseDTO
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router       /api/v1/geneticas/{id} [get]
func (c *GeneticaController) BuscarPorID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		logrus.WithError(err).Error("ID inválido para buscar genética por ID")
		utils.RespondWithError(ctx, http.StatusBadRequest, "ID inválido", nil)
		return
	}
	genetics, err := c.servico.BuscarPorID(uint(id))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.WithError(err).Error("Genética não encontrada ao buscar por ID")
		utils.RespondWithError(ctx, http.StatusNotFound, "Genética não encontrada", nil)
		return
	}
	if err != nil {
		logrus.WithError(err).Error("Erro ao buscar genética por ID")
		utils.RespondWithError(ctx, http.StatusInternalServerError, "Erro interno ao buscar genética", nil)
		return
	}
	utils.RespondWithJSON(ctx, http.StatusOK, genetics)
}

// Atualizar godoc
// @Summary      Atualiza uma genética
// @Description  Atualiza uma genética existente com os dados fornecidos
// @Tags         genetica
// @Accept       json
// @Produce      json
// @Param        id        path      int                  true  "ID da Genética"
// @Param        genetica  body      dto.UpdateGeneticaDTO  true  "Dados da Genética para atualização"
// @Success      200       {object}  dto.GeneticaResponseDTO
// @Failure      400       {object}  map[string]string
// @Failure      404       {object}  map[string]string
// @Failure      500       {object}  map[string]string
// @Router       /api/v1/geneticas/{id} [put]
func (c *GeneticaController) Atualizar(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		logrus.WithError(err).Error("ID inválido para atualização de genética")
		utils.RespondWithError(ctx, http.StatusBadRequest, "ID inválido", nil)
		return
	}

	var updateDto dto.UpdateGeneticaDTO
	if err := ctx.ShouldBindJSON(&updateDto); err != nil {
		logrus.WithError(err).Error("Payload da requisição inválido para atualização de genética")
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

	geneticaAtualizada, err := c.servico.Atualizar(uint(id), &updateDto)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.WithError(err).Error("Genética não encontrada para atualização")
		utils.RespondWithError(ctx, http.StatusNotFound, "Genética não encontrada", nil)
		return
	}
	if err != nil {
		logrus.WithError(err).Error("Erro ao atualizar genética")
		utils.RespondWithError(ctx, http.StatusInternalServerError, "Erro interno ao atualizar genética", nil)
		return
	}

	utils.RespondWithJSON(ctx, http.StatusOK, geneticaAtualizada)
}

// Deletar godoc
// @Summary      Deleta uma genética
// @Description  Deleta uma genética existente
// @Tags         genetica
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID da Genética"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/geneticas/{id} [delete]
func (c *GeneticaController) Deletar(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		logrus.WithError(err).Error("ID inválido para deletar genética")
		utils.RespondWithError(ctx, http.StatusBadRequest, "ID inválido", nil)
		return
	}
	if err := c.servico.Deletar(uint(id)); err != nil {
		logrus.WithError(err).Error("Erro ao deletar genética")
		utils.RespondWithError(ctx, http.StatusInternalServerError, "Erro interno ao deletar genética", nil)
		return
	}
	utils.RespondWithJSON(ctx, http.StatusOK, gin.H{"message": "Genética deletada com sucesso"})
}