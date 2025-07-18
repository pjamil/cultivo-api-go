package service_test

import (
	"bytes"
	"testing"
	"encoding/json" // Adicionado para json.RawMessage

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto" // Adicionado para dto.UsuarioCreateDTO
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/entity" // Adicionado para entity.Usuario
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/service"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/service/test"
)

func TestUsuarioService_Criar(t *testing.T) {
	mockRepo := new(test.MockUsuarioRepositorio)
	s := service.NewUsuarioService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		createDTO := &dto.UsuarioCreateDTO{
			Nome:        "Test User",
			Email:       "test@example.com",
			Senha:       "password",
			Preferencias: json.RawMessage([]byte("null")),
		}
		expectedResponse := &dto.UsuarioResponseDTO{
			ID:           1,
			Nome:         "Test User",
			Email:        "test@example.com",
			Preferencias: json.RawMessage([]byte("null")),
		}

		mockRepo.On("Criar", mock.MatchedBy(func(arg *entity.Usuario) bool {
			return assert.Equal(t, createDTO.Nome, arg.Nome) &&
				   assert.Equal(t, createDTO.Email, arg.Email) &&
				   bytes.Equal(createDTO.Preferencias, arg.Preferencias)
		})).Return(nil).Once().Run(func(args mock.Arguments) {
			usuario := args.Get(0).(*entity.Usuario)
			usuario.ID = 1 // Set the ID before returning
		}).Return(nil)

		response, err := s.Criar(createDTO)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, expectedResponse, response)
		mockRepo.AssertExpectations(t)
	})
}

func TestUsuarioService_ListarTodos(t *testing.T) {
	mockRepo := new(test.MockUsuarioRepositorio)
	s := service.NewUsuarioService(mockRepo)

	t.Run("Success - Usuarios Encontrados", func(t *testing.T) {
		// Arrange
		mockUsuarios := []entity.Usuario{
			{Model: gorm.Model{ID: 1}, Nome: "Usuario 1", Email: "user1@example.com", Preferencias: json.RawMessage("null")},
			{Model: gorm.Model{ID: 2}, Nome: "Usuario 2", Email: "user2@example.com", Preferencias: json.RawMessage("null")},
		}
		expectedTotal := int64(len(mockUsuarios))
		page := 1
		limit := 10

		expectedResponseDTOs := []dto.UsuarioResponseDTO{
			{ID: 1, Nome: "Usuario 1", Email: "user1@example.com", Preferencias: json.RawMessage("null")},
			{ID: 2, Nome: "Usuario 2", Email: "user2@example.com", Preferencias: json.RawMessage("null")},
		}

		mockRepo.On("ListarTodos", page, limit).Return(mockUsuarios, expectedTotal, nil).Once()

		// Act
		responseDTOs, total, err := s.ListarTodos(page, limit)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, responseDTOs)
		assert.Equal(t, expectedTotal, total)
		assert.Equal(t, expectedResponseDTOs, responseDTOs)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success - Nenhum Usuario Encontrado", func(t *testing.T) {
		page := 1
		limit := 10
		mockRepo.On("ListarTodos", page, limit).Return([]entity.Usuario{}, int64(0), nil).Once()

		responseDTOs, total, err := s.ListarTodos(page, limit)

		assert.NoError(t, err)
		assert.NotNil(t, responseDTOs)
		assert.Equal(t, int64(0), total)
		assert.Empty(t, responseDTOs)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Repository Error", func(t *testing.T) {
		page := 1
		limit := 10
		mockRepo.On("ListarTodos", page, limit).Return([]entity.Usuario{}, int64(0), errors.New("erro no repositório")).Once()

		responseDTOs, total, err := s.ListarTodos(page, limit)

		assert.Error(t, err)
		assert.Nil(t, responseDTOs)
		assert.Equal(t, int64(0), total)
		assert.EqualError(t, err, "erro no repositório")
		mockRepo.AssertExpectations(t)
	})
}