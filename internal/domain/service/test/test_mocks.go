package test

import (
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/entity"
	"github.com/stretchr/testify/mock"
)

// MockPlantaRepositorio é um mock para a interface PlantaRepositorio.
type MockPlantaRepositorio struct {
	mock.Mock
}

func (m *MockPlantaRepositorio) Criar(planta *entity.Planta) error {
	args := m.Called(planta)
	return args.Error(0)
}

func (m *MockPlantaRepositorio) ListarTodos(page, limit int) ([]entity.Planta, int64, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]entity.Planta), args.Get(1).(int64), args.Error(2)
}

func (m *MockPlantaRepositorio) BuscarPorID(id uint) (*entity.Planta, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.Planta), args.Error(1)
}

func (m *MockPlantaRepositorio) Atualizar(planta *entity.Planta) error {
	args := m.Called(planta)
	return args.Error(0)
}

func (m *MockPlantaRepositorio) Deletar(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPlantaRepositorio) BuscarPorEspecie(especie entity.Especie) ([]entity.Planta, error) {
	args := m.Called(especie)
	return args.Get(0).([]entity.Planta), args.Error(1)
}

func (m *MockPlantaRepositorio) BuscarPorStatus(status string) ([]entity.Planta, error) {
	args := m.Called(status)
	return args.Get(0).([]entity.Planta), args.Error(1)
}

func (m *MockPlantaRepositorio) ExistePorNome(nome string) bool {
	args := m.Called(nome)
	return args.Bool(0)
}

func (m *MockPlantaRepositorio) CriarRegistroDiario(registro *entity.RegistroDiario) error {
	args := m.Called(registro)
	return args.Error(0)
}

// MockAmbienteRepositorio é um mock para a interface AmbienteRepositorio.
type MockAmbienteRepositorio struct {
	mock.Mock
}

func (m *MockAmbienteRepositorio) Criar(ambiente *entity.Ambiente) error {
	args := m.Called(ambiente)
	return args.Error(0)
}

func (m *MockAmbienteRepositorio) ListarTodos(page, limit int) ([]entity.Ambiente, int64, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]entity.Ambiente), args.Get(1).(int64), args.Error(2)
}

func (m *MockAmbienteRepositorio) BuscarPorID(id uint) (*entity.Ambiente, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.Ambiente), args.Error(1)
}

func (m *MockAmbienteRepositorio) Atualizar(ambiente *entity.Ambiente) error {
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

func (m *MockGeneticaRepositorio) Criar(genetica *entity.Genetica) error {
	args := m.Called(genetica)
	return args.Error(0)
}

func (m *MockGeneticaRepositorio) ListarTodos(page, limit int) ([]entity.Genetica, int64, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]entity.Genetica), args.Get(1).(int64), args.Error(2)
}

func (m *MockGeneticaRepositorio) BuscarPorID(id uint) (*entity.Genetica, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.Genetica), args.Error(1)
}

func (m *MockGeneticaRepositorio) Atualizar(genetica *entity.Genetica) error {
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

func (m *MockMeioCultivoRepositorio) Criar(meioCultivo *entity.MeioCultivo) error {
	args := m.Called(meioCultivo)
	return args.Error(0)
}

func (m *MockMeioCultivoRepositorio) ListarTodos(page, limit int) ([]entity.MeioCultivo, int64, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]entity.MeioCultivo), args.Get(1).(int64), args.Error(2)
}

func (m *MockMeioCultivoRepositorio) BuscarPorID(id uint) (*entity.MeioCultivo, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.MeioCultivo), args.Error(1)
}

func (m *MockMeioCultivoRepositorio) Atualizar(meioCultivo *entity.MeioCultivo) error {
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

func (m *MockUsuarioRepositorio) Criar(usuario *entity.Usuario) error {
	args := m.Called(usuario)
	return args.Error(0)
}

func (m *MockUsuarioRepositorio) ListarTodos(page, limit int) ([]entity.Usuario, int64, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]entity.Usuario), args.Get(1).(int64), args.Error(2)
}

func (m *MockUsuarioRepositorio) BuscarPorID(id uint) (*entity.Usuario, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.Usuario), args.Error(1)
}

func (m *MockUsuarioRepositorio) Atualizar(usuario *entity.Usuario) error {
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

func (m *MockUsuarioRepositorio) BuscarPorEmail(email string) (*entity.Usuario, error) {
	args := m.Called(email)
	return args.Get(0).(*entity.Usuario), args.Error(1)
}

// MockRegistroDiarioRepositorio is a mock for the RegistroDiarioRepositorio interface.
type MockRegistroDiarioRepositorio struct {
	mock.Mock
}

func (m *MockRegistroDiarioRepositorio) Create(registroDiario *entity.RegistroDiario) error {
	args := m.Called(registroDiario)
	return args.Error(0)
}

func (m *MockRegistroDiarioRepositorio) GetByID(id uint) (*entity.RegistroDiario, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.RegistroDiario), args.Error(1)
}

func (m *MockRegistroDiarioRepositorio) GetAll(page, limit int) ([]entity.RegistroDiario, int64, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]entity.RegistroDiario), args.Get(1).(int64), args.Error(2)
}

func (m *MockRegistroDiarioRepositorio) Update(registroDiario *entity.RegistroDiario) error {
	args := m.Called(registroDiario)
	return args.Error(0)
}

func (m *MockRegistroDiarioRepositorio) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRegistroDiarioRepositorio) ListarPorDiarioCultivoID(diarioCultivoID uint, page, limit int) ([]entity.RegistroDiario, int64, error) {
	args := m.Called(diarioCultivoID, page, limit)
	return args.Get(0).([]entity.RegistroDiario), args.Get(1).(int64), args.Error(2)
}