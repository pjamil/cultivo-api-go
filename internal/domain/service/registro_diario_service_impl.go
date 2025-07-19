package service

import (
	"errors"
	"fmt"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/dto"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/entity"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/repository"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/utils"
	"gorm.io/gorm"
)

type registroDiarioService struct {
	registroRepo repository.RegistroDiarioRepository
	diarioRepo   repository.DiarioCultivoRepositorio
}

// NewRegistroDiarioService cria uma nova instância do serviço de RegistroDiario.
func NewRegistroDiarioService(registroRepo repository.RegistroDiarioRepository, diarioRepo repository.DiarioCultivoRepositorio) RegistroDiarioService {
	return &registroDiarioService{
		registroRepo: registroRepo,
		diarioRepo:   diarioRepo,
	}
}

// CriarRegistro cria um novo registro no diário de um cultivo específico.
func (s *registroDiarioService) CriarRegistro(diarioID uint, registroDTO *dto.CreateRegistroDiarioDTO) (*dto.RegistroDiarioResponseDTO, error) {
	// 1. Validar se o diário de cultivo existe
	_, err := s.diarioRepo.GetByID(diarioID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound
		}
		return nil, fmt.Errorf("falha ao verificar diário de cultivo: %w", err)
	}

	// 2. Mapear DTO para o modelo
	registro := &entity.RegistroDiario{
		Titulo:          registroDTO.Titulo,
		Conteudo:        registroDTO.Conteudo,
		Data:            registroDTO.Data,
		Tipo:            registroDTO.Tipo,
		DiarioCultivoID: diarioID,
	}

	// 3. Chamar o repositório para criar
	if err := s.registroRepo.Criar(registro); err != nil {
		return nil, fmt.Errorf("falha ao criar registro no diário: %w", err)
	}

	// 4. Mapear modelo para DTO de resposta
	return &dto.RegistroDiarioResponseDTO{
		ID:       registro.ID,
		Titulo:   registro.Titulo,
		Conteudo: registro.Conteudo,
		Data:     registro.Data,
		Tipo:     registro.Tipo,
	}, nil
}

// ListarRegistrosPorDiarioID lista todos os registros de um diário específico com paginação.
func (s *registroDiarioService) ListarRegistrosPorDiarioID(diarioID uint, page, limit int) ([]dto.RegistroDiarioResponseDTO, int64, error) {
	registros, total, err := s.registroRepo.ListarPorDiarioCultivoID(diarioID, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("falha ao listar registros do diário: %w", err)
	}

	var responseDTOs []dto.RegistroDiarioResponseDTO
	for _, r := range registros {
		responseDTOs = append(responseDTOs, dto.RegistroDiarioResponseDTO{ID: r.ID, Titulo: r.Titulo, Conteudo: r.Conteudo, Data: r.Data, Tipo: r.Tipo})
	}

	return responseDTOs, total, nil
}
