package service

import (
	"errors"
	"fmt"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/repository"
	"gorm.io/gorm"
)

// PlantaService define os métodos que um serviço de planta deve implementar.
// Fornece uma interface para gerenciar plantas no sistema.
type PlantaService interface {
	BuscarPorID(id uint) (*models.Planta, error)
	Criar(plantaDto *dto.CreatePlantaDTO) (*models.Planta, error)
	ListarTodas() ([]models.Planta, error)
	Atualizar(id uint, plantaDto *dto.UpdatePlantaDTO) error
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

func (s *plantaService) Criar(plantaDto *dto.CreatePlantaDTO) (*models.Planta, error) {
	if plantaDto.Nome == "" {
		return nil, errors.New("o nome da planta não pode estar vazio")
	}
	if exists := s.repositorio.ExistePorNome(plantaDto.Nome); exists {
		return nil, errors.New("uma planta com este nome já existe")
	}
	// Validação de entidades relacionadas
	if _, err := s.geneticaRepositorio.BuscarPorID(plantaDto.GeneticaID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("genética não encontrada")
		}
		return nil, fmt.Errorf("falha ao buscar genética com ID %d: %w", plantaDto.GeneticaID, err)
	}
	if _, err := s.ambienteRepositorio.BuscarPorID(plantaDto.AmbienteID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("ambiente não encontrado")
		}
		return nil, fmt.Errorf("falha ao buscar ambiente com ID %d: %w", plantaDto.AmbienteID, err)
	}
	if _, err := s.meioRepositorio.BuscarPorID(plantaDto.MeioCultivoID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("meio de cultivo não encontrado")
		}
		return nil, fmt.Errorf("falha ao buscar meio de cultivo com ID %d: %w", plantaDto.MeioCultivoID, err)
	}

	planta := models.Planta{
		Nome:          plantaDto.Nome,
		ComecandoDe:   plantaDto.ComecandoDe,
		Especie:       models.Especie(plantaDto.Especie),
		Status:        models.PlantaStatus(plantaDto.Status),
		GeneticaID:    plantaDto.GeneticaID,
		MeioCultivoID: plantaDto.MeioCultivoID,
		AmbienteID:    plantaDto.AmbienteID,
		UsuarioID:     plantaDto.UsuarioID,
	}

	if !plantaDto.DataPlantio.IsZero() {
		planta.DataPlantio = &plantaDto.DataPlantio
	}
	if plantaDto.Notas != "" {
		planta.Notas = &plantaDto.Notas
	}

	if err := s.repositorio.Criar(&planta); err != nil {
		return nil, err
	}
	return &planta, nil
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

func (s *plantaService) Atualizar(id uint, plantaDto *dto.UpdatePlantaDTO) error {
	plantaExistente, err := s.repositorio.BuscarPorID(id)
	if err != nil {
		return fmt.Errorf("falha ao buscar planta com ID %d: %w", id, err)
	}

	// Atualiza os campos da planta existente com os dados do DTO
	if plantaDto.Nome != "" {
		plantaExistente.Nome = plantaDto.Nome
	}
	if plantaDto.ComecandoDe != "" {
		plantaExistente.ComecandoDe = plantaDto.ComecandoDe
	}
	if plantaDto.Especie != "" {
		plantaExistente.Especie = models.Especie(plantaDto.Especie)
	}
	if !plantaDto.DataPlantio.IsZero() {
		plantaExistente.DataPlantio = &plantaDto.DataPlantio
	}
	if !plantaDto.DataColheita.IsZero() {
		plantaExistente.DataColheita = &plantaDto.DataColheita
	}
	if plantaDto.Status != "" {
		plantaExistente.Status = models.PlantaStatus(plantaDto.Status)
	}
	// O campo EstagioCrescimento não existe em models.Planta, então não pode ser atribuído diretamente.
	// if plantaDto.EstagioCrescimento != "" {
	// 	plantaExistente.EstagioCrescimento = plantaDto.EstagioCrescimento
	// }
	if plantaDto.Notas != "" {
		plantaExistente.Notas = &plantaDto.Notas
	}
	if plantaDto.GeneticaID != 0 {
		plantaExistente.GeneticaID = plantaDto.GeneticaID
	}
	if plantaDto.MeioCultivoID != 0 {
		plantaExistente.MeioCultivoID = plantaDto.MeioCultivoID
	}
	if plantaDto.AmbienteID != 0 {
		plantaExistente.AmbienteID = plantaDto.AmbienteID
	}
	if plantaDto.UsuarioID != 0 {
		plantaExistente.UsuarioID = plantaDto.UsuarioID
	}

	// Validação de entidades relacionadas (se os IDs foram alterados)
	if plantaDto.GeneticaID != 0 {
		if _, err := s.geneticaRepositorio.BuscarPorID(plantaDto.GeneticaID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("genética não encontrada")
			}
			return fmt.Errorf("falha ao buscar genética com ID %d: %w", plantaDto.GeneticaID, err)
		}
	}
	if plantaDto.AmbienteID != 0 {
		if _, err := s.ambienteRepositorio.BuscarPorID(plantaDto.AmbienteID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("ambiente não encontrado")
			}
			return fmt.Errorf("falha ao buscar ambiente com ID %d: %w", plantaDto.AmbienteID, err)
		}
	}
	if plantaDto.MeioCultivoID != 0 {
		if _, err := s.meioRepositorio.BuscarPorID(plantaDto.MeioCultivoID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("meio de cultivo não encontrado")
			}
			return fmt.Errorf("falha ao buscar meio de cultivo com ID %d: %w", plantaDto.MeioCultivoID, err)
		}
	}

	return s.repositorio.Atualizar(plantaExistente)
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
