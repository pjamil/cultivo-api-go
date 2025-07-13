package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/service"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/service/test"
)

func TestUsuarioService_Criar(t *testing.T) {
	mockRepo := new(test.MockUsuarioRepositorio)
	s := service.NewUsuarioService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("ExistePorEmail", "test@example.com").Return(false).Once()
		mockRepo.On("Criar", mock.Anything).Return(nil).Once()

		_, err := s.Criar(nil) // Passar um DTO real aqui
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}