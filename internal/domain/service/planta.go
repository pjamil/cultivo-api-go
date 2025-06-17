package service

import (
	"errors"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/repository"
)

type PlantaService struct {
	repo repository.PlantaRepository
}

func NewPlantService(repo repository.PlantaRepository) *PlantaService {
	return &PlantaService{repo: repo}
}

func (s *PlantaService) CreatePlanta(planta *models.Planta) error {
	if planta.Nome == "" {
		return errors.New("plant name cannot be empty")
	}
	return s.repo.Create(planta)
}

func (s *PlantaService) GetAllPlants() ([]models.Planta, error) {
	return s.repo.FindAll()
}

func (s *PlantaService) GetPlantByID(id uint) (*models.Planta, error) {
	return s.repo.FindByID(id)
}

func (s *PlantaService) UpdatePlant(plant *models.Planta) error {
	return s.repo.Update(plant)
}

func (s *PlantaService) DeletePlant(id uint) error {
	return s.repo.Delete(id)
}
func (s *PlantaService) GetPlantsBySpecies(species models.Especie) ([]models.Planta, error) {
	return s.repo.FindBySpecies(species)
}
func (s *PlantaService) GetPlantsByStatus(status string) ([]models.Planta, error) {
	return s.repo.FindByStatus(status)
}
