package service

import (
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/repository"
)

type PlantService struct {
	repo repository.PlantaRepository
}

func NewPlantService(repo repository.PlantaRepository) *PlantService {
	return &PlantService{repo: repo}
}

func (s *PlantService) CreatePlant(plant *models.Planta) error {
	return s.repo.Create(plant)
}

func (s *PlantService) GetAllPlants() ([]models.Planta, error) {
	return s.repo.FindAll()
}

func (s *PlantService) GetPlantByID(id uint) (*models.Planta, error) {
	return s.repo.FindByID(id)
}

func (s *PlantService) UpdatePlant(plant *models.Planta) error {
	return s.repo.Update(plant)
}

func (s *PlantService) DeletePlant(id uint) error {
	return s.repo.Delete(id)
}
func (s *PlantService) GetPlantsBySpecies(species models.Especie) ([]models.Planta, error) {
	return s.repo.FindBySpecies(species)
}
func (s *PlantService) GetPlantsByStatus(status string) ([]models.Planta, error) {
	return s.repo.FindByStatus(status)
}
