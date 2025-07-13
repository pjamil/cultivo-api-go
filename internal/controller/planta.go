package controller

import (
	"errors"
	
	"net/http"
	"strconv"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/service"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)



// PlantaController lida com as requisições HTTP para plantas
type PlantaController struct {
	plantaServico service.PlantaService
}

// NewPlantaController cria um novo PlantController
func NewPlantaController(plantaServico service.PlantaService) *PlantaController {
	return &PlantaController{plantaServico: plantaServico}
}

// CriarPlanta godoc
// @Summary      Cria uma nova planta
// @Description  Adiciona uma nova planta ao sistema de cultivo
// @Tags         plantas
// @Accept       json
// @Produce      json
// @Param        planta  body      dto.CreatePlantaDTO  true  "Dados para criação da planta"
// @Success      201    {object}  dto.PlantaResponseDTO
// @Failure      400    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Router       /api/v1/plantas [post]
func (c *PlantaController) Criar(ctx *gin.Context) {
	var createDto dto.CreatePlantaDTO
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

	if plantaCriada, err := c.plantaServico.Criar(&createDto); err != nil {
		logrus.WithError(err).Error("Erro ao criar planta")
		utils.RespondWithError(ctx, http.StatusInternalServerError, "Erro interno ao criar planta", err.Error())
		return
	} else {
		utils.RespondWithJSON(ctx, http.StatusCreated, plantaCriada)
	}
}

// ListarPlantas godoc
// @Summary      Lista todas as plantas com paginação
// @Description  Retorna uma lista paginada de todas as plantas
// @Tags         plantas
// @Accept       json
// @Produce      json
// @Param        page   query     int  false  "Número da página (padrão: 1)"
// @Param        limit  query     int  false  "Limite de itens por página (padrão: 10)"
// @Success      200  {object}  dto.PaginatedResponse{data=[]dto.PlantaResponseDTO}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /api/v1/plantas [get]
func (c *PlantaController) Listar(ctx *gin.Context) {
	var pagination dto.PaginationParams
	if err := ctx.ShouldBindQuery(&pagination); err != nil {
		logrus.WithError(err).Error("Parâmetros de paginação inválidos para listar plantas")
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

	paginatedResponse, err := c.plantaServico.ListarTodas(pagination.Page, pagination.Limit)
	if err != nil {
		logrus.WithError(err).Error("Erro ao listar plantas com paginação")
		utils.RespondWithError(ctx, http.StatusInternalServerError, "Erro interno ao listar plantas", err.Error())
		return
	}

	utils.RespondWithJSON(ctx, http.StatusOK, paginatedResponse)
}

const (
	ErroIDPlantaInvalido          = "ID da planta inválido"
	ErroPlantaNaoEncontrada       = "Planta não encontrada"
	ErroPayloadRequisicaoInvalido = "Payload da requisição inválido"
	ErroAtualizarPlanta           = "Erro ao atualizar a planta"
	ErroCriarPlanta               = "Erro ao criar a planta"
	ErroRecuperarPlanta           = "Erro ao recuperar a planta"
	ErroDeletarPlanta             = "Erro ao deletar a planta"
	SucessoAtualizarPlanta        = "Planta atualizada com sucesso"
	SucessoCriarPlanta            = "Planta criada com sucesso"
	SucessoRecuperarPlanta        = "Planta recuperada com sucesso"
	SucessoListarPlantas          = "Plantas recuperadas com sucesso"
	ErroListarPlantas             = "Erro ao recuperar as plantas"
)

// BuscarPlantaPorID godoc
// @Summary      Busca uma planta por ID
// @Description  Retorna os detalhes de uma planta específica
// @ID           get-plant-by-id
// @Tags         plantas
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID da Planta"
// @Success      200  {object}  dto.PlantaResponseDTO
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /api/v1/plantas/{id} [get]
func (c *PlantaController) BuscarPorID(ctx *gin.Context) {
	logrus.WithFields(logrus.Fields{
		"id": ctx.Param("id"),
	}).Info("Buscando planta por ID")
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		logrus.WithError(err).Error("Erro ao converter ID da planta")
		utils.RespondWithError(ctx, http.StatusBadRequest, "ID inválido", utils.ErrInvalidInput.Error())
		return
	}

	planta, err := c.plantaServico.BuscarPorID(uint(id))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.WithError(err).Error("Planta não encontrada ao buscar por ID")
		utils.RespondWithError(ctx, http.StatusNotFound, "Planta não encontrada", utils.ErrNotFound.Error())
		return
	}
	utils.RespondWithJSON(ctx, http.StatusOK, planta)
}

// AtualizarPlanta godoc
// @Summary      Atualiza uma planta
// @Description  Atualiza uma planta existente
// @Tags         plantas
// @Accept       json
// @Produce      json
// @Param        id      path      int            true  "ID da Planta"
// @Param        planta  body      dto.UpdatePlantaDTO  true  "Dados para atualização da planta"
// @Success      200    {object}  dto.PlantaResponseDTO
// @Failure      400    {object}  map[string]interface{}
// @Failure      404    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Router       /api/v1/plantas/{id} [put]
func (c *PlantaController) Atualizar(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		logrus.WithError(err).Error("Erro ao converter ID da planta para atualização")
		utils.RespondWithError(ctx, http.StatusBadRequest, "ID inválido", utils.ErrInvalidInput.Error())
		return
	}
	var updateDto dto.UpdatePlantaDTO
	if err := ctx.ShouldBindJSON(&updateDto); err != nil {
		logrus.WithError(err).Error("Payload da requisição inválido para atualização de planta")
		utils.RespondWithError(ctx, http.StatusBadRequest, "Requisição inválida", utils.ErrInvalidInput.Error())
		return
	}

	// A chamada ao serviço de atualização agora retorna o DTO de resposta
	plantaAtualizada, err := c.plantaServico.Atualizar(uint(id), &updateDto)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.WithError(err).Error("Planta não encontrada para atualização")
		utils.RespondWithError(ctx, http.StatusNotFound, "Planta não encontrada", utils.ErrNotFound.Error())
		return
	}
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"operation": "update_plant",
			"plant_id":  id,
			"error":     err,
		}).Error(ErroAtualizarPlanta)
		utils.RespondWithError(ctx, http.StatusInternalServerError, "Erro interno ao atualizar planta", err.Error())
		return
	}

	utils.RespondWithJSON(ctx, http.StatusOK, plantaAtualizada)
}

// DeletarPlanta godoc
// @Summary      Deleta uma planta
// @Description  Deleta uma planta existente
// @Tags         plantas
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID da Planta"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /api/v1/plantas/{id} [delete]
func (c *PlantaController) Deletar(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		logrus.WithError(err).Error("Erro ao converter ID da planta para deleção")
		utils.RespondWithError(ctx, http.StatusBadRequest, "ID inválido", utils.ErrInvalidInput.Error())
		return
	}
	if err := c.plantaServico.Deletar(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.WithError(err).Error("Planta não encontrada para deleção")
			utils.RespondWithError(ctx, http.StatusNotFound, "Planta não encontrada", utils.ErrNotFound.Error())
			return
		}
		logrus.WithError(err).Error("Erro ao deletar planta")
		utils.RespondWithError(ctx, http.StatusInternalServerError, "Erro interno ao deletar planta", err.Error())
		return
	}
	logrus.WithFields(logrus.Fields{
		"operation": "delete_plant",
		"plant_id":  id,
	}).Info("Planta deletada com sucesso")
	utils.RespondWithJSON(ctx, http.StatusOK, gin.H{"message": "Planta deletada com sucesso"})
}

// RegistrarFato godoc
// @Summary      Registra um fato na vida da planta
// @Description  Registra um evento ou observação importante sobre a planta
// @Tags         plantas
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID da Planta"
// @Param        fato  body      dto.RegistrarFatoDTO  true  "Dados do fato a ser registrado"
// @Success      200    {object}  map[string]interface{}
// @Failure      400    {object}  map[string]interface{}
// @Failure      404    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Router       /api/v1/plantas/{id}/registrar-fato [post]
func (c *PlantaController) RegistrarFato(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		logrus.WithError(err).Error("Erro ao converter ID da planta para registrar fato")
		utils.RespondWithError(ctx, http.StatusBadRequest, "ID inválido", utils.ErrInvalidInput.Error())
		return
	}

	var fatoDto dto.RegistrarFatoDTO
	if err := ctx.ShouldBindJSON(&fatoDto); err != nil {
		logrus.WithError(err).Error("Payload da requisição inválido para registrar fato")
		utils.RespondWithError(ctx, http.StatusBadRequest, "Requisição inválida", utils.ErrInvalidInput.Error())
		return
	}

	if err := c.plantaServico.RegistrarFato(uint(id), models.RegistroTipo(fatoDto.Tipo), fatoDto.Titulo, fatoDto.Conteudo); err != nil {
		logrus.WithError(err).Error("Erro ao registrar fato para a planta")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.RespondWithError(ctx, http.StatusNotFound, "Planta não encontrada", utils.ErrNotFound.Error())
			return
		}
		utils.RespondWithError(ctx, http.StatusInternalServerError, "Erro interno ao registrar fato", err.Error())
		return
	}

	utils.RespondWithJSON(ctx, http.StatusOK, gin.H{"message": "Fato registrado com sucesso"})
}
