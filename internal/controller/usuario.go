package controller

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/service"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const invalidIDMsg = "ID inválido, deve ser um número inteiro positivo"

type UsuarioController struct {
	servico service.UsuarioService
}

func NewUsuarioController(servico service.UsuarioService) *UsuarioController {
	return &UsuarioController{servico}
}

// Criar godoc
// @Summary      Cria um novo usuário
// @Description  Cria um novo usuário com os dados fornecidos
// @Tags         usuario
// @Accept       json
// @Produce      json
// @Param        usuario  body      dto.UsuarioCreateDTO  true  "Dados do Usuário"
// @Success      201      {object}  dto.UsuarioResponseDTO
// @Failure      400      {object}  map[string]string
// @Failure      409      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /api/v1/usuarios [post]
func (c *UsuarioController) Criar(ctx *gin.Context) {
	var dto dto.UsuarioCreateDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		logrus.WithError(err).Error("Payload da requisição inválido para criar usuário")
		utils.RespondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	usuario, err := c.servico.Criar(&dto)
	if err != nil {
		// Trate erro de unique constraint (e-mail já cadastrado)
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			utils.RespondWithError(ctx, http.StatusConflict, "E-mail já cadastrado")
			return
		}
		logrus.WithError(err).Error("Erro ao criar usuário")
		utils.RespondWithError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(ctx, http.StatusCreated, usuario)
}

// BuscarPorID godoc
// @Summary      Busca um usuário por ID
// @Description  Retorna os detalhes de um usuário específico
// @Tags         usuario
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID do Usuário"
// @Success      200  {object}  dto.UsuarioResponseDTO
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/usuarios/{id} [get]
func (c *UsuarioController) BuscarPorID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		logrus.WithError(err).Error("ID inválido para buscar usuário por ID")
		utils.RespondWithError(ctx, http.StatusBadRequest, invalidIDMsg)
		return
	}
	usuario, err := c.servico.BuscarPorID(uint(id))
	if err != nil {
		utils.RespondWithError(ctx, http.StatusNotFound, "Usuário não encontrado")
		return
	}
	utils.RespondWithJSON(ctx, http.StatusOK, usuario)
}

// Listar godoc
// @Summary      Lista todos os usuários
// @Description  Retorna uma lista de todos os usuários cadastrados
// @Tags         usuario
// @Produce      json
// @Success      200  {array}   dto.UsuarioResponseDTO
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/usuarios [get]
func (c *UsuarioController) Listar(ctx *gin.Context) {
	usuarios, err := c.servico.ListarTodos()
	if err != nil {
		logrus.WithError(err).Error("Erro ao listar usuários")
		utils.RespondWithError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	if len(usuarios) == 0 {
		utils.RespondWithJSON(ctx, http.StatusOK, gin.H{"message": "Nenhum usuário encontrado"})
		return
	}

	utils.RespondWithJSON(ctx, http.StatusOK, usuarios)
}

// Atualizar godoc
// @Summary      Atualiza um usuário
// @Description  Atualiza um usuário existente com os dados fornecidos
// @Tags         usuario
// @Accept       json
// @Produce      json
// @Param        id       path      int                 true  "ID do Usuário"
// @Param        usuario  body      dto.UsuarioUpdateDTO  true  "Dados do Usuário para atualização"
// @Success      200      {object}  dto.UsuarioResponseDTO
// @Failure      400      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /api/v1/usuarios/{id} [put]
func (c *UsuarioController) Atualizar(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		logrus.WithError(err).Error("ID inválido para atualizar usuário")
		utils.RespondWithError(ctx, http.StatusBadRequest, invalidIDMsg)
		return
	}
	var dto dto.UsuarioUpdateDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		logrus.WithError(err).Error("Payload da requisição inválido para atualizar usuário")
		utils.RespondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	usuarioAtualizado, err := c.servico.Atualizar(uint(id), &dto)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.WithError(err).Error("Usuário não encontrado para atualização")
		utils.RespondWithError(ctx, http.StatusNotFound, "Usuário não encontrado")
		return
	}
	if err != nil {
		logrus.WithError(err).Error("Erro ao atualizar usuário")
		utils.RespondWithError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(ctx, http.StatusOK, usuarioAtualizado)
}

// Deletar godoc
// @Summary      Deleta um usuário
// @Description  Deleta um usuário existente
// @Tags         usuario
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID do Usuário"
// @Success      204  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/usuarios/{id} [delete]
func (c *UsuarioController) Deletar(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		logrus.WithError(err).Error("ID inválido para deletar usuário")
		utils.RespondWithError(ctx, http.StatusBadRequest, invalidIDMsg)
		return
	}
	if err := c.servico.Deletar(uint(id)); err != nil {
		logrus.WithError(err).Error("Erro ao deletar usuário")
		utils.RespondWithError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(ctx, http.StatusNoContent, nil)
}
