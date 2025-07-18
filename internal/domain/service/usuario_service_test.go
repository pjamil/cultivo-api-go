package service_test

import (
	"bytes"
	"testing"
	"encoding/json" // Adicionado para json.RawMessage

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

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
