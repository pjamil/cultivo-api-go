package service

import (
	"errors"

	"context"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/repository"
	"gorm.io/gorm"
)

// PlantService defines the methods that a plant service should implement.
// It provides an interface for managing plants in the system.
type PlantService interface {
	GetPlantaById(ctx context.Context, id uint) (*models.Planta, error)
	CreatePlanta(ctx context.Context, planta *models.Planta) error
	GetAllPlants(ctx context.Context) ([]models.Planta, error)
	UpdatePlant(ctx context.Context, plant *models.Planta) error
	DeletePlant(ctx context.Context, id uint) error
	GetPlantsBySpecies(ctx context.Context, species models.Especie) ([]models.Planta, error)
	GetPlantsByStatus(ctx context.Context, status string) ([]models.Planta, error)
}

// PlantaService provides methods to manage plants in the system.
// It interacts with the PlantaRepository to perform CRUD operations on plants.
type PlantaService struct {
	repo         repository.PlantaRepository
	geneticaRepo repository.GeneticaRepository
	ambienteRepo repository.AmbienteRepository
	meioRepo     repository.MeioCultivoRepository
}

func NewPlantService(
	repo repository.PlantaRepository,
	geneticaRepo repository.GeneticaRepository,
	ambienteRepo repository.AmbienteRepository,
	meioRepo repository.MeioCultivoRepository,
) *PlantaService {
	return &PlantaService{
		repo:         repo,
		geneticaRepo: geneticaRepo,
		ambienteRepo: ambienteRepo,
		meioRepo:     meioRepo,
	}
}

func (s *PlantaService) CreatePlanta(ctx context.Context, planta *models.Planta) error {
	if planta.Nome == "" {
		return errors.New("plant name cannot be empty")
	}
	if exists := s.repo.ExistsByName(planta.Nome); exists {
		return errors.New("plant with this name already exists")
	}
	// Validação de entidades relacionadas
	if _, err := s.geneticaRepo.FindByID(ctx, planta.GeneticaID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("genética não encontrada")
		}
		return err
	}
	if _, err := s.ambienteRepo.FindByID(ctx, planta.AmbienteID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("ambiente não encontrado")
		}
		return err
	}
	if _, err := s.meioRepo.FindByID(planta.MeioCultivoID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("meio de cultivo não encontrado")
		}
		return err
	}
	return s.repo.Create(planta)
}

func (s *PlantaService) GetPlantaById(ctx context.Context, id uint) (*models.Planta, error) {
	planta, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return planta, nil
}

func (s *PlantaService) GetAllPlants() ([]models.Planta, error) {
	return s.repo.FindAll()
}

func (s *PlantaService) UpdatePlant(ctx context.Context, plant *models.Planta) error {
	if _, err := s.repo.FindByID(plant.ID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Planta não encontrada")
		}
	}
	// Validação de entidades relacionadas
	if _, err := s.geneticaRepo.FindByID(ctx, plant.GeneticaID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("genética não encontrada")
		}
	}
	if _, err := s.ambienteRepo.FindByID(ctx, plant.AmbienteID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("ambiente não encontrado")
		}
	}
	if _, err := s.meioRepo.FindByID(plant.MeioCultivoID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("meio de cultivo não encontrado")
		}
	}

	return s.repo.Update(plant)
}

func (s *PlantaService) DeletePlant(id uint) error {
	if _, err := s.repo.FindByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Planta não encontrada")
		}
	}
	return s.repo.Delete(id)
}

func (s *PlantaService) GetPlantsBySpecies(species models.Especie) ([]models.Planta, error) {
	return s.repo.FindBySpecies(species)
}
func (s *PlantaService) GetPlantsByStatus(status string) ([]models.Planta, error) {
	return s.repo.FindByStatus(status)
}
