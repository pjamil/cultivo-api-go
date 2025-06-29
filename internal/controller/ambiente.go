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
// @Router       /ambiente [post]
func (c *AmbienteController) Criar(ctx *gin.Context) {
	var dto dto.CreateAmbienteDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		logrus.WithError(err).Error("Payload da requisição inválido para criar ambiente")
		utils.RespondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	ambiente, err := c.servico.Criar(&dto)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"operation": "create_ambiente",
			"error":     err,
		}).Error("Erro ao criar ambiente")
		utils.RespondWithError(ctx, http.StatusInternalServerError, "Erro ao criar ambiente")
		return
	}
	utils.RespondWithJSON(ctx, http.StatusCreated, ambiente)
}

// Listar godoc
// @Summary      Lista todos os ambientes
// @Description  Retorna uma lista de todos os ambientes cadastrados
// @Tags         ambiente
// @Produce      json
// @Success      200  {array}   dto.AmbienteResponseDTO
// @Failure      500  {object}  map[string]string
// @Router       /ambiente [get]
func (c *AmbienteController) Listar(ctx *gin.Context) {
	ambientes, err := c.servico.ListarTodos()
	if err != nil {
		logrus.WithError(err).Error("Erro ao listar ambientes")
		utils.RespondWithError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	if len(ambientes) == 0 {
		utils.RespondWithJSON(ctx, http.StatusOK, gin.H{"message": "Nenhum ambiente encontrado"})
		return
	}

	utils.RespondWithJSON(ctx, http.StatusOK, ambientes)
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
// @Router       /ambiente/{id} [get]
func (c *AmbienteController) BuscarPorID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		logrus.WithError(err).Error("ID inválido para buscar ambiente por ID")
		utils.RespondWithError(ctx, http.StatusBadRequest, "ID inválido")
		return
	}
	ambiente, err := c.servico.BuscarPorID(uint(id))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.WithError(err).Error("Ambiente não encontrado ao buscar por ID")
		utils.RespondWithError(ctx, http.StatusNotFound, "Ambiente não encontrado")
		return
	}
	if err != nil {
		logrus.WithError(err).Error("Erro ao buscar ambiente por ID")
		utils.RespondWithError(ctx, http.StatusInternalServerError, "Erro ao buscar ambiente")
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
// @Router       /ambiente/{id} [put]
func (c *AmbienteController) Atualizar(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		logrus.WithError(err).Error("ID inválido para atualização de ambiente")
		utils.RespondWithError(ctx, http.StatusBadRequest, "ID inválido")
		return
	}

	var updateDto dto.UpdateAmbienteDTO
	if err := ctx.ShouldBindJSON(&updateDto); err != nil {
		logrus.WithError(err).Error("Payload da requisição inválido para atualização de ambiente")
		utils.RespondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ambienteAtualizado, err := c.servico.Atualizar(uint(id), &updateDto)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.WithError(err).Error("Ambiente não encontrado para atualização")
		utils.RespondWithError(ctx, http.StatusNotFound, "Ambiente não encontrado")
		return
	}
	if err != nil {
		logrus.WithError(err).Error("Erro ao atualizar ambiente")
		utils.RespondWithError(ctx, http.StatusInternalServerError, "Erro ao atualizar ambiente")
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
// @Router       /ambientes/{id} [delete]
func (c *AmbienteController) Deletar(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		logrus.WithError(err).Error("ID inválido para deletar ambiente")
		utils.RespondWithError(ctx, http.StatusBadRequest, "ID inválido")
		return
	}
	if err := c.servico.Deletar(uint(id)); err != nil {
		logrus.WithError(err).Error("Erro ao deletar ambiente")
		utils.RespondWithError(ctx, http.StatusInternalServerError, "Erro ao deletar ambiente")
		return
	}
	utils.RespondWithJSON(ctx, http.StatusOK, gin.H{"message": "Ambiente deletado com sucesso"})
}
