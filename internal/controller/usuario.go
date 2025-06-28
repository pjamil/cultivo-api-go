package controller

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/service"
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
// @Success      201      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]string
// @Failure      409      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /usuarios [post]
func (c *UsuarioController) Criar(ctx *gin.Context) {
	var dto dto.UsuarioCreateDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		logrus.WithError(err).Error("Payload da requisição inválido para criar usuário")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	usuario, err := c.servico.Criar(&dto)
	if err != nil {
		// Trate erro de unique constraint (e-mail já cadastrado)
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			ctx.JSON(http.StatusConflict, gin.H{"error": "E-mail já cadastrado"})
			return
		}
		logrus.WithError(err).Error("Erro ao criar usuário")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, usuario)
}

// BuscarPorID godoc
// @Summary      Busca um usuário por ID
// @Description  Retorna os detalhes de um usuário específico
// @Tags         usuario
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID do Usuário"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /usuarios/{id} [get]
func (c *UsuarioController) BuscarPorID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		logrus.WithError(err).Error("ID inválido para buscar usuário por ID")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": invalidIDMsg})
		return
	}
	usuario, err := c.servico.BuscarPorID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	}
	ctx.JSON(http.StatusOK, usuario)
}

// Listar godoc
// @Summary      Lista todos os usuários
// @Description  Retorna uma lista de todos os usuários cadastrados
// @Tags         usuario
// @Produce      json
// @Success      200  {array}   map[string]interface{}
// @Failure      500  {object}  map[string]string
// @Router       /usuarios [get]
func (c *UsuarioController) Listar(ctx *gin.Context) {
	usuarios, err := c.servico.ListarTodos()
	if err != nil {
		logrus.WithError(err).Error("Erro ao listar usuários")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, usuarios)
}

// Atualizar godoc
// @Summary      Atualiza um usuário
// @Description  Atualiza um usuário existente com os dados fornecidos
// @Tags         usuario
// @Accept       json
// @Produce      json
// @Param        id       path      int                 true  "ID do Usuário"
// @Param        usuario  body      dto.UsuarioUpdateDTO  true  "Dados do Usuário para atualização"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /usuarios/{id} [put]
func (c *UsuarioController) Atualizar(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		logrus.WithError(err).Error("ID inválido para atualizar usuário")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": invalidIDMsg})
		return
	}
	var dto dto.UsuarioUpdateDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		logrus.WithError(err).Error("Payload da requisição inválido para atualizar usuário")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = c.servico.Atualizar(uint(id), &dto)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.WithError(err).Error("Usuário não encontrado para atualização")
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	}
	if err != nil {
		logrus.WithError(err).Error("Erro ao atualizar usuário")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusOK)
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
// @Router       /usuarios/{id} [delete]
func (c *UsuarioController) Deletar(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		logrus.WithError(err).Error("ID inválido para deletar usuário")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": invalidIDMsg})
		return
	}
	if err := c.servico.Deletar(uint(id)); err != nil {
		logrus.WithError(err).Error("Erro ao deletar usuário")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}