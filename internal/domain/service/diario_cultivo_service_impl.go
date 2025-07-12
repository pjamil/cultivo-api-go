package service

import (
	"errors"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/repository"
)

type diarioCultivoService struct {
	repo            repository.DiarioCultivoRepository
	plantaRepo      repository.PlantaRepositorio
	ambientesRepo   repository.AmbienteRepositorio
}

func NewDiarioCultivoService(
	repo repository.DiarioCultivoRepository,
	plantaRepo repository.PlantaRepositorio,
	ambientesRepo repository.AmbienteRepositorio,
) *diarioCultivoService {
	return &diarioCultivoService{
		repo:            repo,
		plantaRepo:      plantaRepo,
		ambientesRepo:   ambientesRepo,
	}
}

func (s *diarioCultivoService) Create(dto *dto.CreateDiarioCultivoDTO) (*dto.DiarioCultivoResponseDTO, error) {
	diarioCultivo := models.DiarioCultivo{
		Nome:        dto.Nome,
		DataInicio:  dto.DataInicio,
		DataFim:     dto.DataFim,
		Ativo:       true,
		UsuarioID:   dto.UsuarioID,
		Privacidade: dto.Privacidade,
		Tags:        dto.Tags,
	}

	if dto.Ativo != nil {
		diarioCultivo.Ativo = *dto.Ativo
	}

	if err := s.repo.Create(&diarioCultivo); err != nil {
		return nil, err
	}

	// Adicionar plantas
	if len(dto.PlantasIDs) > 0 {
		var plantas []*models.Planta
		for _, id := range dto.PlantasIDs {
			planta, err := s.plantaRepo.BuscarPorID(id)
			if err != nil {
				return nil, errors.New("Planta não encontrada")
			}
			plantas = append(plantas, planta)
		}
		if err := s.repo.AddPlantas(diarioCultivo.ID, plantas); err != nil {
			return nil, err
		}
	}

	// Adicionar ambientes
	if len(dto.AmbientesIDs) > 0 {
		var ambientes []*models.Ambiente
		for _, id := range dto.AmbientesIDs {
			ambiente, err := s.ambientesRepo.BuscarPorID(id)
			if err != nil {
				return nil, errors.New("Ambiente não encontrado")
			}
			ambientes = append(ambientes, ambiente)
		}
		if err := s.repo.AddAmbientes(diarioCultivo.ID, ambientes); err != nil {
			return nil, err
		}
	}

	return s.GetByID(diarioCultivo.ID)
}

func (s *diarioCultivoService) GetByID(id uint) (*dto.DiarioCultivoResponseDTO, error) {
	diarioCultivo, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return &dto.DiarioCultivoResponseDTO{
		ID:          diarioCultivo.ID,
		Nome:        diarioCultivo.Nome,
		DataInicio:  diarioCultivo.DataInicio,
		DataFim:     diarioCultivo.DataFim,
		Ativo:       diarioCultivo.Ativo,
		UsuarioID:   diarioCultivo.UsuarioID,
		Plantas:     diarioCultivo.Plantas,
		Ambientes:   diarioCultivo.Ambientes,
		Privacidade: diarioCultivo.Privacidade,
		Tags:        diarioCultivo.Tags,
		CreatedAt:   diarioCultivo.CreatedAt,
		UpdatedAt:   diarioCultivo.UpdatedAt,
	}, nil
}

func (s *diarioCultivoService) GetAll(page, limit int) (*dto.PaginatedResponse, error) {
	diariosCultivo, total, err := s.repo.GetAll(page, limit)
	if err != nil {
		return nil, err
	}

	var diarioCultivoDTOs []interface{}
	for _, dc := range diariosCultivo {
		diarioCultivoDTOs = append(diarioCultivoDTOs, dto.DiarioCultivoResponseDTO{
			ID:          dc.ID,
			Nome:        dc.Nome,
			DataInicio:  dc.DataInicio,
			DataFim:     dc.DataFim,
			Ativo:       dc.Ativo,
			UsuarioID:   dc.UsuarioID,
			Plantas:     dc.Plantas,
			Ambientes:   dc.Ambientes,
			Privacidade: dc.Privacidade,
			Tags:        dc.Tags,
			CreatedAt:   dc.CreatedAt,
			UpdatedAt:   dc.UpdatedAt,
		})
	}

	return &dto.PaginatedResponse{
		Data:  diarioCultivoDTOs,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

func (s *diarioCultivoService) Update(id uint, dto *dto.UpdateDiarioCultivoDTO) (*dto.DiarioCultivoResponseDTO, error) {
	diarioCultivo, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if dto.Nome != "" {
		diarioCultivo.Nome = dto.Nome
	}
	if dto.DataInicio != nil {
		diarioCultivo.DataInicio = *dto.DataInicio
	}
	if dto.DataFim != nil {
		diarioCultivo.DataFim = dto.DataFim
	}
	if dto.Ativo != nil {
		diarioCultivo.Ativo = *dto.Ativo
	}
	if dto.Privacidade != "" {
		diarioCultivo.Privacidade = dto.Privacidade
	}
	if dto.Tags != "" {
		diarioCultivo.Tags = dto.Tags
	}

	// Atualizar relacionamentos de plantas
	if dto.PlantasIDs != nil {
		// Remover todas as plantas existentes e adicionar as novas
		var existingPlantasPointers []*models.Planta
		for i := range diarioCultivo.Plantas {
			existingPlantasPointers = append(existingPlantasPointers, &diarioCultivo.Plantas[i])
		}
		if err := s.repo.RemovePlantas(diarioCultivo.ID, existingPlantasPointers); err != nil {
			return nil, err
		}
		var novasPlantas []*models.Planta
		for _, pID := range dto.PlantasIDs {
			planta, err := s.plantaRepo.BuscarPorID(pID)
			if err != nil {
				return nil, errors.New("Planta não encontrada")
			}
			novasPlantas = append(novasPlantas, planta)
		}
		if err := s.repo.AddPlantas(diarioCultivo.ID, novasPlantas); err != nil {
			return nil, err
		}
	}

	// Atualizar relacionamentos de ambientes
	if dto.AmbientesIDs != nil {
		// Remover todos os ambientes existentes e adicionar os novos
		var existingAmbientesPointers []*models.Ambiente
		for i := range diarioCultivo.Ambientes {
			existingAmbientesPointers = append(existingAmbientesPointers, &diarioCultivo.Ambientes[i])
		}
		if err := s.repo.RemoveAmbientes(diarioCultivo.ID, existingAmbientesPointers); err != nil {
			return nil, err
		}
		var novosAmbientes []*models.Ambiente
		for _, aID := range dto.AmbientesIDs {
			ambiente, err := s.ambientesRepo.BuscarPorID(aID)
			if err != nil {
				return nil, errors.New("Ambiente não encontrado")
			}
			novosAmbientes = append(novosAmbientes, ambiente)
		}
		if err := s.repo.AddAmbientes(diarioCultivo.ID, novosAmbientes); err != nil {
			return nil, err
		}
	}

	if err := s.repo.Update(diarioCultivo); err != nil {
		return nil, err
	}

	return s.GetByID(diarioCultivo.ID)
}

func (s *diarioCultivoService) Delete(id uint) error {
	return s.repo.Delete(id)
}