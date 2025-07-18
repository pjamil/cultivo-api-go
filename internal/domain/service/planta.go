package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/entity"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/repository"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/utils"
	"gorm.io/gorm"
)

// PlantaService define os métodos que um serviço de planta deve implementar.
// Fornece uma interface para gerenciar plantas no sistema.
type PlantaService interface {
	BuscarPorID(id uint) (*dto.PlantaResponseDTO, error)
	Criar(plantaDto *dto.CreatePlantaDTO) (*dto.PlantaResponseDTO, error)
	ListarTodas(page, limit int) (*dto.PaginatedResponse, error)
	Atualizar(id uint, plantaDto *dto.UpdatePlantaDTO) (*dto.PlantaResponseDTO, error)
	Deletar(id uint) error
	BuscarPorEspecie(especie entity.Especie) ([]entity.Planta, error)
	BuscarPorStatus(status string) ([]entity.Planta, error)
	RegistrarFato(plantaID uint, tipo entity.RegistroTipo, titulo string, conteudo string) error
}

// PlantaService fornece métodos para gerenciar plantas no sistema.
// Ele interage com o PlantaRepository para realizar operações CRUD em plantas.
type plantaService struct {
	repositorio               repository.PlantaRepositorio
	geneticaRepositorio       repository.GeneticaRepositorio
	ambienteRepositorio       repository.AmbienteRepositorio
	meioRepositorio           repository.MeioCultivoRepositorio
	registroDiarioRepositorio repository.RegistroDiarioRepositorio
}

func NewPlantaService(
	repositorio repository.PlantaRepositorio,
	geneticaRepositorio repository.GeneticaRepositorio,
	ambienteRepositorio repository.AmbienteRepositorio,
	meioRepositorio repository.MeioCultivoRepositorio,
	registroDiarioRepositorio repository.RegistroDiarioRepositorio,
) PlantaService {
	return &plantaService{
		repositorio:               repositorio,
		geneticaRepositorio:       geneticaRepositorio,
		ambienteRepositorio:       ambienteRepositorio,
		meioRepositorio:           meioRepositorio,
		registroDiarioRepositorio: registroDiarioRepositorio,
	}
}

func (s *plantaService) Criar(plantaDto *dto.CreatePlantaDTO) (*dto.PlantaResponseDTO, error) {
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

	planta := entity.Planta{
		Nome:          plantaDto.Nome,
		ComecandoDe:   plantaDto.ComecandoDe,
		Especie:       entity.Especie(plantaDto.Especie),
		Status:        entity.PlantaStatus(plantaDto.Status),
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
	return &dto.PlantaResponseDTO{
		ID:            planta.ID,
		Nome:          planta.Nome,
		ComecandoDe:   planta.ComecandoDe,
		Especie:       string(planta.Especie),
		DataPlantio:   utils.TimePtr(utils.DereferenceTimePtr(planta.DataPlantio)),
		DataColheita:  planta.DataColheita,
		Status:        string(planta.Status),
		Notas:         utils.StringPtr(utils.DereferenceStringPtr(planta.Notas)),
		GeneticaID:    planta.GeneticaID,
		MeioCultivoID: planta.MeioCultivoID,
		AmbienteID:    planta.AmbienteID,
		UsuarioID:     planta.UsuarioID,
		CreatedAt:     planta.CreatedAt,
		UpdatedAt:     planta.UpdatedAt,
	}, nil
}

func (s *plantaService) BuscarPorID(id uint) (*dto.PlantaResponseDTO, error) {
	planta, err := s.repositorio.BuscarPorID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("falha ao obter planta com ID %d: %w", id, err)
	}
	return &dto.PlantaResponseDTO{
		ID:            planta.ID,
		Nome:          planta.Nome,
		ComecandoDe:   planta.ComecandoDe,
		Especie:       string(planta.Especie),
		DataPlantio:   utils.TimePtr(utils.DereferenceTimePtr(planta.DataPlantio)),
		DataColheita:  planta.DataColheita,
		Status:        string(planta.Status),
		Notas:         utils.StringPtr(utils.DereferenceStringPtr(planta.Notas)),
		GeneticaID:    planta.GeneticaID,
		MeioCultivoID: planta.MeioCultivoID,
		AmbienteID:    planta.AmbienteID,
		UsuarioID:     planta.UsuarioID,
		CreatedAt:     planta.CreatedAt,
		UpdatedAt:     planta.UpdatedAt,
	}, nil
}

func (s *plantaService) ListarTodas(page, limit int) (*dto.PaginatedResponse, error) {
	plantas, total, err := s.repositorio.ListarTodos(page, limit)
	if err != nil {
		return nil, err
	}

	responseDTOs := make([]dto.PlantaResponseDTO, 0, len(plantas))
	for _, planta := range plantas {
		responseDTOs = append(responseDTOs, dto.PlantaResponseDTO{
			ID:            planta.ID,
			Nome:          planta.Nome,
			ComecandoDe:   planta.ComecandoDe,
			Especie:       string(planta.Especie),
			DataPlantio:   utils.TimePtr(utils.DereferenceTimePtr(planta.DataPlantio)),
			DataColheita:  planta.DataColheita,
			Status:        string(planta.Status),
			Notas:         utils.StringPtr(utils.DereferenceStringPtr(planta.Notas)),
			GeneticaID:    planta.GeneticaID,
			MeioCultivoID: planta.MeioCultivoID,
			AmbienteID:    planta.AmbienteID,
			UsuarioID:     planta.UsuarioID,
			CreatedAt:     planta.CreatedAt,
			UpdatedAt:     planta.UpdatedAt,
		})
	}

		dataBytes, err := json.Marshal(responseDTOs)
		if err != nil {
			return nil, fmt.Errorf("falha ao serializar plantas: %w", err)
		}

	return &dto.PaginatedResponse{
		Data:  dataBytes,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

func (s *plantaService) Atualizar(id uint, plantaDto *dto.UpdatePlantaDTO) (*dto.PlantaResponseDTO, error) {
	plantaExistente, err := s.repositorio.BuscarPorID(id)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar planta com ID %d: %w", id, err)
	}

	// Atualiza os campos da planta existente com os dados do DTO
	if plantaDto.Nome != "" {
		plantaExistente.Nome = plantaDto.Nome
	}
	if plantaDto.ComecandoDe != "" {
		plantaExistente.ComecandoDe = plantaDto.ComecandoDe
	}
	if plantaDto.Especie != "" {
		plantaExistente.Especie = entity.Especie(plantaDto.Especie)
	}
	if !plantaDto.DataPlantio.IsZero() {
		plantaExistente.DataPlantio = &plantaDto.DataPlantio
	}
	if !plantaDto.DataColheita.IsZero() {
		plantaExistente.DataColheita = &plantaDto.DataColheita
	}
	if plantaDto.Status != "" {
		plantaExistente.Status = entity.PlantaStatus(plantaDto.Status)
	}
	// O campo EstagioCrescimento não existe em entity.Planta, então não pode ser atribuído diretamente.
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
				return nil, errors.New("genética não encontrada")
			}
			return nil, fmt.Errorf("falha ao buscar genética com ID %d: %w", plantaDto.GeneticaID, err)
		}
	}
	if plantaDto.AmbienteID != 0 {
		if _, err := s.ambienteRepositorio.BuscarPorID(plantaDto.AmbienteID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("ambiente não encontrado")
			}
			return nil, fmt.Errorf("falha ao buscar ambiente com ID %d: %w", plantaDto.AmbienteID, err)
		}
	}
	if plantaDto.MeioCultivoID != 0 {
		if _, err := s.meioRepositorio.BuscarPorID(plantaDto.MeioCultivoID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("meio de cultivo não encontrado")
			}
			return nil, fmt.Errorf("falha ao buscar meio de cultivo com ID %d: %w", plantaDto.MeioCultivoID, err)
		}
	}

	if err := s.repositorio.Atualizar(plantaExistente); err != nil {
		return nil, err
	}

	return &dto.PlantaResponseDTO{
		ID:            plantaExistente.ID,
		Nome:          plantaExistente.Nome,
		ComecandoDe:   plantaExistente.ComecandoDe,
		Especie:       string(plantaExistente.Especie),
		DataPlantio:   utils.TimePtr(utils.DereferenceTimePtr(plantaExistente.DataPlantio)),
		DataColheita:  plantaExistente.DataColheita,
		Status:        string(plantaExistente.Status),
		Notas:         utils.StringPtr(utils.DereferenceStringPtr(plantaExistente.Notas)),
		GeneticaID:    plantaExistente.GeneticaID,
		MeioCultivoID: plantaExistente.MeioCultivoID,
		AmbienteID:    plantaExistente.AmbienteID,
		UsuarioID:     plantaExistente.UsuarioID,
		CreatedAt:     plantaExistente.CreatedAt,
		UpdatedAt:     plantaExistente.UpdatedAt,
	}, nil
}

func (s *plantaService) Deletar(id uint) error {
	if _, err := s.repositorio.BuscarPorID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("planta não encontrada")
		}
	}
	return s.repositorio.Deletar(id)
}

func (s *plantaService) BuscarPorEspecie(especie entity.Especie) ([]entity.Planta, error) {
	return s.repositorio.BuscarPorEspecie(especie)
}
func (s *plantaService) BuscarPorStatus(status string) ([]entity.Planta, error) {
	return s.repositorio.BuscarPorStatus(status)
}

func (s *plantaService) RegistrarFato(plantaID uint, tipo entity.RegistroTipo, titulo string, conteudo string) error {
	// Verificar se a planta existe
	_, err := s.repositorio.BuscarPorID(plantaID)
	if err != nil {
		return fmt.Errorf("planta com ID %d não encontrada: %w", plantaID, err)
	}

	// Criar o novo registro diário
	registro := &entity.RegistroDiario{
		Data:           time.Now(),
		Tipo:           tipo,
		Titulo:         titulo,
		Conteudo:       conteudo,
		ReferenciaID:   &plantaID,
		ReferenciaTipo: entity.String("planta"),
	}

	// Salvar o registro diário usando o repositório
	if err := s.registroDiarioRepositorio.Create(registro); err != nil {
		return fmt.Errorf("falha ao registrar fato para a planta %d: %w", plantaID, err)
	}

	return nil
}
