package controller

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const invalidIDMsg = "ID inválido, deve ser um número inteiro positivo"

type UsuarioController struct {
	servico service.UsuarioService
}

func NewUsuarioController(servico service.UsuarioService) *UsuarioController {
	return &UsuarioController{servico}
}

func (c *UsuarioController) Criar(ctx *gin.Context) {
	var dto dto.UsuarioCreateDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, usuario)
}

func (c *UsuarioController) BuscarPorID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
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

func (c *UsuarioController) Listar(ctx *gin.Context) {
	usuarios, err := c.servico.ListarTodos()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, usuarios)
}

func (c *UsuarioController) Atualizar(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": invalidIDMsg})
		return
	}
	var dto dto.UsuarioUpdateDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = c.servico.Atualizar(uint(id), &dto)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusOK)
}

func (c *UsuarioController) Deletar(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": invalidIDMsg})
		return
	}
	if err := c.servico.Deletar(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}
