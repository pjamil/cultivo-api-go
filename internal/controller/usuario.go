package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/service"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)



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
	usuario, err := c.servico.Criar(&dto)
	if err != nil {
		// Trate erro de unique constraint (e-mail já cadastrado)
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			utils.RespondWithError(ctx, http.StatusConflict, "E-mail já cadastrado", nil)
			return
		}
		logrus.WithError(err).Error("Erro ao criar usuário")
		utils.RespondWithError(ctx, http.StatusInternalServerError, "Erro interno ao criar usuário", err.Error())
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
		utils.RespondWithError(ctx, http.StatusBadRequest, "ID inválido", utils.ErrInvalidInput.Error())
		return
	}
	usuario, err := c.servico.BuscarPorID(uint(id))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.RespondWithError(ctx, http.StatusNotFound, "Usuário não encontrado", utils.ErrNotFound.Error())
		return
	}
	if err != nil {
		logrus.WithError(err).Error("Erro ao buscar usuário por ID")
		utils.RespondWithError(ctx, http.StatusInternalServerError, "Erro interno ao buscar usuário", utils.ErrInternalServer.Error())
		return
	}
	utils.RespondWithJSON(ctx, http.StatusOK, usuario)
}

// Listar godoc
// @Summary      Lista todos os usuários com paginação
// @Description  Retorna uma lista paginada de todos os usuários cadastrados
// @Tags         usuario
// @Produce      json
// @Param        page   query     int  false  "Número da página (padrão: 1)"
// @Param        limit  query     int  false  "Limite de itens por página (padrão: 10)"
// @Success      200  {object}  dto.PaginatedResponse{data=[]dto.UsuarioResponseDTO}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/usuarios [get]
func (c *UsuarioController) Listar(ctx *gin.Context) {
	var pagination dto.PaginationParams
	if err := ctx.ShouldBindQuery(&pagination); err != nil {
		logrus.WithError(err).Error("Parâmetros de paginação inválidos para listar usuários")
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

	paginatedResponse, err := c.servico.ListarTodos(pagination.Page, pagination.Limit)
	if err != nil {
		logrus.WithError(err).Error("Erro ao listar usuários com paginação")
		utils.RespondWithError(ctx, http.StatusInternalServerError, "Erro interno ao listar usuários", utils.ErrInternalServer.Error())
		return
	}

	utils.RespondWithJSON(ctx, http.StatusOK, paginatedResponse)
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
		utils.RespondWithError(ctx, http.StatusBadRequest, "ID inválido", utils.ErrInvalidInput.Error())
		return
	}
	var dto dto.UsuarioUpdateDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		logrus.WithError(err).Error("Payload da requisição inválido para atualizar usuário")
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
	usuarioAtualizado, err := c.servico.Atualizar(uint(id), &dto)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.WithError(err).Error("Usuário não encontrado para atualização")
		utils.RespondWithError(ctx, http.StatusNotFound, "Usuário não encontrado", utils.ErrNotFound.Error())
		return
	}
	if err != nil {
		logrus.WithError(err).Error("Erro ao atualizar usuário")
		utils.RespondWithError(ctx, http.StatusInternalServerError, "Erro interno ao atualizar usuário", utils.ErrInternalServer.Error())
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
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/usuarios/{id} [delete]
func (c *UsuarioController) Deletar(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		logrus.WithError(err).Error("ID inválido para deletar usuário")
		utils.RespondWithError(ctx, http.StatusBadRequest, "ID inválido", utils.ErrInvalidInput.Error())
		return
	}
	if err := c.servico.Deletar(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.WithError(err).Error("Usuário não encontrado para deleção")
			utils.RespondWithError(ctx, http.StatusNotFound, "Usuário não encontrado", utils.ErrNotFound.Error())
			return
		}
		logrus.WithError(err).Error("Erro ao deletar usuário")
		utils.RespondWithError(ctx, http.StatusInternalServerError, "Erro interno ao deletar usuário", utils.ErrInternalServer.Error())
		return
	}
	utils.RespondWithJSON(ctx, http.StatusNoContent, nil)
}

// Login godoc
// @Summary      Autentica um usuário e retorna um token JWT
// @Description  Recebe credenciais de usuário (email e senha) e retorna um token JWT para acesso autenticado.
// @Tags         usuario
// @Accept       json
// @Produce      json
// @Param        credentials  body      dto.LoginPayload  true  "Credenciais do Usuário"
// @Success      200          {object}  map[string]string "token"
// @Failure      400          {object}  map[string]string
// @Failure      401          {object}  map[string]string
// @Failure      500          {object}  map[string]string
// @Router       /api/v1/login [post]
func (c *UsuarioController) Login(ctx *gin.Context) {
	var payload dto.LoginPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		logrus.WithError(err).Error("Payload da requisição de login inválido")
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

	token, err := c.servico.Login(&payload)
	if err != nil {
		if errors.Is(err, utils.ErrInvalidCredentials) {
			utils.RespondWithError(ctx, http.StatusUnauthorized, "Credenciais inválidas", nil)
			return
		}
		logrus.WithError(err).Error("Erro ao tentar login")
		utils.RespondWithError(ctx, http.StatusInternalServerError, "Erro interno ao tentar login", nil)
		return
	}

	utils.RespondWithJSON(ctx, http.StatusOK, gin.H{"token": token})
}
