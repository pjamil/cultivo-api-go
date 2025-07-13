package test

import (
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"github.com/stretchr/testify/mock"
)

// MockPlantaRepositorio é um mock para a interface PlantaRepositorio.
type MockPlantaRepositorio struct {
	mock.Mock
}

func (m *MockPlantaRepositorio) Criar(planta *models.Planta) error {
	args := m.Called(planta)
	return args.Error(0)
}

func (m *MockPlantaRepositorio) ListarTodos(page, limit int) ([]models.Planta, int64, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]models.Planta), args.Get(1).(int64), args.Error(2)
}

func (m *MockPlantaRepositorio) BuscarPorID(id uint) (*models.Planta, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Planta), args.Error(1)
}

func (m *MockPlantaRepositorio) Atualizar(planta *models.Planta) error {
	args := m.Called(planta)
	return args.Error(0)
}

func (m *MockPlantaRepositorio) Deletar(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPlantaRepositorio) BuscarPorEspecie(especie models.Especie) ([]models.Planta, error) {
	args := m.Called(especie)
	return args.Get(0).([]models.Planta), args.Error(1)
}

func (m *MockPlantaRepositorio) BuscarPorStatus(status string) ([]models.Planta, error) {
	args := m.Called(status)
	return args.Get(0).([]models.Planta), args.Error(1)
}

func (m *MockPlantaRepositorio) ExistePorNome(nome string) bool {
	args := m.Called(nome)
	return args.Bool(0)
}

func (m *MockPlantaRepositorio) CriarRegistroDiario(registro *models.RegistroDiario) error {
	args := m.Called(registro)
	return args.Error(0)
}

// MockAmbienteRepositorio é um mock para a interface AmbienteRepositorio.
type MockAmbienteRepositorio struct {
	mock.Mock
}

func (m *MockAmbienteRepositorio) Criar(ambiente *models.Ambiente) error {
	args := m.Called(ambiente)
	return args.Error(0)
}

func (m *MockAmbienteRepositorio) ListarTodos(page, limit int) ([]models.Ambiente, int64, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]models.Ambiente), args.Get(1).(int64), args.Error(2)
}

func (m *MockAmbienteRepositorio) BuscarPorID(id uint) (*models.Ambiente, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Ambiente), args.Error(1)
}

func (m *MockAmbienteRepositorio) Atualizar(ambiente *models.Ambiente) error {
	args := m.Called(ambiente)
	return args.Error(0)
}

func (m *MockAmbienteRepositorio) Deletar(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockGeneticaRepositorio é um mock para a interface GeneticaRepositorio.
type MockGeneticaRepositorio struct {
	mock.Mock
}

func (m *MockGeneticaRepositorio) Criar(genetica *models.Genetica) error {
	args := m.Called(genetica)
	return args.Error(0)
}

func (m *MockGeneticaRepositorio) ListarTodos(page, limit int) ([]models.Genetica, int64, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]models.Genetica), args.Get(1).(int64), args.Error(2)
}

func (m *MockGeneticaRepositorio) BuscarPorID(id uint) (*models.Genetica, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Genetica), args.Error(1)
}

func (m *MockGeneticaRepositorio) Atualizar(genetica *models.Genetica) error {
	args := m.Called(genetica)
	return args.Error(0)
}

func (m *MockGeneticaRepositorio) Deletar(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockMeioCultivoRepositorio é um mock para a interface MeioCultivoRepositorio.
type MockMeioCultivoRepositorio struct {
	mock.Mock
}

func (m *MockMeioCultivoRepositorio) Criar(meioCultivo *models.MeioCultivo) error {
	args := m.Called(meioCultivo)
	return args.Error(0)
}

func (m *MockMeioCultivoRepositorio) ListarTodos(page, limit int) ([]models.MeioCultivo, int64, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]models.MeioCultivo), args.Get(1).(int64), args.Error(2)
}

func (m *MockMeioCultivoRepositorio) BuscarPorID(id uint) (*models.MeioCultivo, error) {
	args := m.Called(id)
	return args.Get(0).(*models.MeioCultivo), args.Error(1)
}

func (m *MockMeioCultivoRepositorio) Atualizar(meioCultivo *models.MeioCultivo) error {
	args := m.Called(meioCultivo)
	return args.Error(0)
}

func (m *MockMeioCultivoRepositorio) Deletar(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockUsuarioRepositorio é um mock para a interface UsuarioRepositorio.
type MockUsuarioRepositorio struct {
	mock.Mock
}

func (m *MockUsuarioRepositorio) Criar(usuario *models.Usuario) error {
	args := m.Called(usuario)
	return args.Error(0)
}

func (m *MockUsuarioRepositorio) ListarTodos(page, limit int) ([]models.Usuario, int64, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]models.Usuario), args.Get(1).(int64), args.Error(2)
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

func (m *MockUsuarioRepositorio) BuscarPorEmail(email string) (*models.Usuario, error) {
	args := m.Called(email)
	return args.Get(0).(*models.Usuario), args.Error(1)
}