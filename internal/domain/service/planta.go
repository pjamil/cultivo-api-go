package service

import (
	"errors"
	"fmt"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/repository"
	"gorm.io/gorm"
)

// PlantaService define os métodos que um serviço de planta deve implementar.
// Fornece uma interface para gerenciar plantas no sistema.
type PlantaService interface {
	BuscarPorID(id uint) (*models.Planta, error)
	Criar(planta *models.Planta) error
	ListarTodas() ([]models.Planta, error)
	Atualizar(planta *models.Planta) error
	Deletar(id uint) error
	BuscarPorEspecie(especie models.Especie) ([]models.Planta, error)
	BuscarPorStatus(status string) ([]models.Planta, error)
}

// PlantaService fornece métodos para gerenciar plantas no sistema.
// Ele interage com o PlantaRepository para realizar operações CRUD em plantas.
type plantaService struct {
	repositorio           repository.PlantaRepositorio
	geneticaRepositorio   repository.GeneticaRepositorio
	ambienteRepositorio repository.AmbienteRepositorio
	meioRepositorio     repository.MeioCultivoRepositorio
}

func NewPlantaService(
	repositorio repository.PlantaRepositorio,
	geneticaRepositorio repository.GeneticaRepositorio,
	ambienteRepositorio repository.AmbienteRepositorio,
	meioRepositorio repository.MeioCultivoRepositorio,
) PlantaService {
	return &plantaService{
		repositorio:           repositorio,
		geneticaRepositorio:   geneticaRepositorio,
		ambienteRepositorio: ambienteRepositorio,
		meioRepositorio:     meioRepositorio,
	}
}

func (s *plantaService) Criar(planta *models.Planta) error {
	if planta.Nome == "" {
		return errors.New("o nome da planta não pode estar vazio")
	}
	if exists := s.repositorio.ExistePorNome(planta.Nome); exists {
		return errors.New("uma planta com este nome já existe")
	}
	// Validação de entidades relacionadas
	if _, err := s.geneticaRepositorio.BuscarPorID(planta.GeneticaID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("genética não encontrada")
		}
		return fmt.Errorf("falha ao buscar genética com ID %d: %w", planta.GeneticaID, err)
	}
	if _, err := s.ambienteRepositorio.BuscarPorID(planta.AmbienteID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("ambiente não encontrado")
		}
		return fmt.Errorf("falha ao buscar ambiente com ID %d: %w", planta.AmbienteID, err)
	}
	if _, err := s.meioRepositorio.BuscarPorID(planta.MeioCultivoID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("meio de cultivo não encontrado")
		}
		return fmt.Errorf("falha ao buscar meio de cultivo com ID %d: %w", planta.MeioCultivoID, err)
	}
	return s.repositorio.Criar(planta)
}

func (s *plantaService) BuscarPorID(id uint) (*models.Planta, error) {
	planta, err := s.repositorio.BuscarPorID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("falha ao obter planta com ID %d: %w", id, err)
	}
	return planta, nil
}

func (s *plantaService) ListarTodas() ([]models.Planta, error) {
	return s.repositorio.ListarTodos()
}

func (s *plantaService) Atualizar(planta *models.Planta) error {
	if _, err := s.repositorio.BuscarPorID(planta.ID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}
	}
	// Validação de entidades relacionadas
	if _, err := s.geneticaRepositorio.BuscarPorID(planta.GeneticaID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("genética não encontrada")
		}
	}
	if _, err := s.ambienteRepositorio.BuscarPorID(planta.AmbienteID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("ambiente não encontrado")
		}
	}
	if _, err := s.meioRepositorio.BuscarPorID(planta.MeioCultivoID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("meio de cultivo não encontrado")
		}
	}

	return s.repositorio.Atualizar(planta)
}

func (s *plantaService) Deletar(id uint) error {
	if _, err := s.repositorio.BuscarPorID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("planta não encontrada")
		}
	}
	return s.repositorio.Deletar(id)
}

func (s *plantaService) BuscarPorEspecie(especie models.Especie) ([]models.Planta, error) {
	return s.repositorio.BuscarPorEspecie(especie)
}
func (s *plantaService) BuscarPorStatus(status string) ([]models.Planta, error) {
	return s.repositorio.BuscarPorStatus(status)
}
