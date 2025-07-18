package repository

import (
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/entity"
)

type UsuarioRepositorio interface {
	Criar(usuario *entity.Usuario) error
	ListarTodos(page, limit int) ([]entity.Usuario, int64, error)
	BuscarPorID(id uint) (*entity.Usuario, error)
	Atualizar(usuario *entity.Usuario) error
	Deletar(id uint) error
	BuscarPorEmail(email string) (*entity.Usuario, error)
	ExistePorEmail(email string) bool
}