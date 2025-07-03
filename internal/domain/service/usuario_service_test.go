package service_test

import (
	"errors"
	"testing"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockUsuarioRepositorio struct {
	mock.Mock
}

func (m *MockUsuarioRepositorio) Criar(usuario *models.Usuario) error {
	args := m.Called(usuario)
	return args.Error(0)
}

func (m *MockUsuarioRepositorio) ListarTodos() ([]models.Usuario, error) {
	args := m.Called()
	return args.Get(0).([]models.Usuario), args.Error(1)
}

func (m *MockUsuarioRepositorio) BuscarPorID(id uint) (*models.Usuario, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Usuario), args.Error(1)
}

func (m *MockUsuarioRepositorio) Atualizar(usuario *models.Usuario) error {
	args := m.Called(usuario)
	return args.Error(0)
}

func (m *MockUsuarioRepositorio) Deletar(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUsuarioRepositorio) ExistePorEmail(email string) bool {
	args := m.Called(email)
	return args.Bool(0)
}

func TestUsuarioService_Criar(t *testing.T) {
	mockRepo := new(MockUsuarioRepositorio)
	service := service.NewUsuarioService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		createDTO := &dto.UsuarioCreateDTO{Nome: "Usuario Teste", Email: "teste@example.com", Senha: "password123"}
		expectedUsuario := &models.Usuario{Nome: "Usuario Teste", Email: "teste@example.com"}
		expectedResponse := &dto.UsuarioResponseDTO{ID: 1, Nome: "Usuario Teste", Email: "teste@example.com"}

		mockRepo.On("ExistePorEmail", createDTO.Email).Return(false).Once()
		mockRepo.On("Criar", mock.AnythingOfType("*models.Usuario")).Run(func(args mock.Arguments) {
			usuario := args.Get(0).(*models.Usuario)
			usuario.ID = 1
		}).Return(nil).Once()

		response, err := service.Criar(createDTO)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, expectedResponse.Nome, response.Nome)
		assert.Equal(t, expectedResponse.ID, response.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Email Already Exists", func(t *testing.T) {
		createDTO := &dto.UsuarioCreateDTO{Nome: "Usuario Teste", Email: "teste@example.com", Senha: "password123"}

		mockRepo.On("ExistePorEmail", createDTO.Email).Return(true).Once()

		response, err := service.Criar(createDTO)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.EqualError(t, err, "e-mail já cadastrado")
		mockRepo.AssertExpectations(t)
	})

		t.Run("Error - Repository Error", func(t *testing.T) {
		createDTO := &dto.UsuarioCreateDTO{Nome: "Usuario Teste", Email: "teste@example.com", Senha: "password123"}

		mockRepo.On("ExistePorEmail", createDTO.Email).Return(false).Once()
		mockRepo.On("Criar", mock.AnythingOfType("*models.Usuario")).Return(errors.New("erro no repositório")).Once()

		response, err := service.Criar(createDTO)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.EqualError(t, err, "erro no repositório")
		mockRepo.AssertExpectations(t)
	})
}

func TestUsuarioService_ListarTodos(t *testing.T) {
	mockRepo := new(MockUsuarioRepositorio)
	service := service.NewUsuarioService(mockRepo)

	t.Run("Success - Usuarios Encontrados", func(t *testing.T) {
		expectedUsuarios := []models.Usuario{
			{ID: 1, Nome: "Usuario 1", Email: "user1@example.com"},
			{ID: 2, Nome: "Usuario 2", Email: "user2@example.com"},
		}
		expectedResponse := []dto.UsuarioResponseDTO{
			{ID: 1, Nome: "Usuario 1", Email: "user1@example.com"},
			{ID: 2, Nome: "Usuario 2", Email: "user2@example.com"},
		}

		mockRepo.On("ListarTodos").Return(expectedUsuarios, nil).Once()

		response, err := service.ListarTodos()

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, expectedResponse, response)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success - Nenhum Usuario Encontrado", func(t *testing.T) {
		mockRepo.On("ListarTodos").Return([]models.Usuario{}, nil).Once()

		response, err := service.ListarTodos()

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Empty(t, response)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Repository Error", func(t *testing.T) {
		mockRepo.On("ListarTodos").Return([]models.Usuario{}, errors.New("erro no repositório")).Once()

		response, err := service.ListarTodos()

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.EqualError(t, err, "erro no repositório")
		mockRepo.AssertExpectations(t)
	})
}

func TestUsuarioService_BuscarPorID(t *testing.T) {
	mockRepo := new(MockUsuarioRepositorio)
	service := service.NewUsuarioService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		usuarioID := uint(1)
		expectedUsuario := &models.Usuario{ID: usuarioID, Nome: "Usuario Teste", Email: "teste@example.com"}
		expectedResponse := &dto.UsuarioResponseDTO{ID: usuarioID, Nome: "Usuario Teste", Email: "teste@example.com"}

		mockRepo.On("BuscarPorID", usuarioID).Return(expectedUsuario, nil).Once()

		response, err := service.BuscarPorID(usuarioID)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, expectedResponse.ID, response.ID)
		assert.Equal(t, expectedResponse.Nome, response.Nome)
		assert.Equal(t, expectedResponse.Email, response.Email)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		usuarioID := uint(999)

		mockRepo.On("BuscarPorID", usuarioID).Return((*models.Usuario)(nil), gorm.ErrRecordNotFound).Once()

		response, err := service.BuscarPorID(usuarioID)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		usuarioID := uint(1)
		expectedError := errors.New("erro no repositório")

		mockRepo.On("BuscarPorID", usuarioID).Return((*models.Usuario)(nil), expectedError).Once()

		response, err := service.BuscarPorID(usuarioID)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), expectedError.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestUsuarioService_Atualizar(t *testing.T) {
	mockRepo := new(MockUsuarioRepositorio)
	service := service.NewUsuarioService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		usuarioID := uint(1)
		updateDTO := &dto.UsuarioUpdateDTO{Nome: "Usuario Atualizado", Email: "atualizado@example.com"}
		existingUsuario := &models.Usuario{ID: usuarioID, Nome: "Usuario Antigo", Email: "antigo@example.com"}
		expectedResponse := &dto.UsuarioResponseDTO{ID: usuarioID, Nome: "Usuario Atualizado", Email: "atualizado@example.com"}

		mockRepo.On("BuscarPorID", usuarioID).Return(existingUsuario, nil).Once()
		mockRepo.On("ExistePorEmail", updateDTO.Email).Return(false).Once()
		mockRepo.On("Atualizar", mock.AnythingOfType("*models.Usuario")).Run(func(args mock.Arguments) {
			usuario := args.Get(0).(*models.Usuario)
			usuario.Nome = updateDTO.Nome
			usuario.Email = updateDTO.Email
		}).Return(nil).Once()

		response, err := service.Atualizar(usuarioID, updateDTO)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, expectedResponse.Nome, response.Nome)
		assert.Equal(t, expectedResponse.Email, response.Email)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		usuarioID := uint(999)
		updateDTO := &dto.UsuarioUpdateDTO{Nome: "Usuario Atualizado"}

		mockRepo.On("BuscarPorID", usuarioID).Return((*models.Usuario)(nil), gorm.ErrRecordNotFound).Once()

		response, err := service.Atualizar(usuarioID, updateDTO)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Email Already Exists", func(t *testing.T) {
		usuarioID := uint(1)
		updateDTO := &dto.UsuarioUpdateDTO{Nome: "Usuario Atualizado", Email: "existente@example.com"}
		existingUsuario := &models.Usuario{ID: usuarioID, Nome: "Usuario Antigo", Email: "antigo@example.com"}

		mockRepo.On("BuscarPorID", usuarioID).Return(existingUsuario, nil).Once()
		mockRepo.On("ExistePorEmail", updateDTO.Email).Return(true).Once()

		response, err := service.Atualizar(usuarioID, updateDTO)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.EqualError(t, err, "e-mail já cadastrado")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error on Update", func(t *testing.T) {
		usuarioID := uint(1)
		updateDTO := &dto.UsuarioUpdateDTO{Nome: "Usuario Atualizado", Email: "atualizado@example.com"}
		existingUsuario := &models.Usuario{ID: usuarioID, Nome: "Usuario Antigo", Email: "antigo@example.com"}
		expectedError := errors.New("erro no repositório")

		mockRepo.On("BuscarPorID", usuarioID).Return(existingUsuario, nil).Once()
		mockRepo.On("ExistePorEmail", updateDTO.Email).Return(false).Once()
		mockRepo.On("Atualizar", mock.AnythingOfType("*models.Usuario")).Return(expectedError).Once()

		response, err := service.Atualizar(usuarioID, updateDTO)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), expectedError.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestUsuarioService_Deletar(t *testing.T) {
	mockRepo := new(MockUsuarioRepositorio)
	service := NewUsuarioService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		usuarioID := uint(1)
		mockRepo.On("Deletar", usuarioID).Return(nil).Once()

		err := service.Deletar(usuarioID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		usuarioID := uint(999)
		mockRepo.On("Deletar", usuarioID).Return(gorm.ErrRecordNotFound).Once()

		err := service.Deletar(usuarioID)

		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		usuarioID := uint(1)
		expectedError := errors.New("erro no repositório")
		mockRepo.On("Deletar", usuarioID).Return(expectedError).Once()

		err := service.Deletar(usuarioID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), expectedError.Error())
		mockRepo.AssertExpectations(t)
	})
}