package service

import (
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/repository"
)

type PlantService struct {
	repo repository.PlantRepository
}

func NewPlantService(repo repository.PlantRepository) *PlantService {
	return &PlantService{repo: repo}
}

func (s *PlantService) CreatePlant(plant *models.Plant) error {
	return s.repo.Create(plant)
}

func (s *PlantService) GetAllPlants() ([]models.Plant, error) {
	return s.repo.FindAll()
}

func (s *PlantService) GetPlantByID(id uint) (*models.Plant, error) {
	return s.repo.FindByID(id)
}

func (s *PlantService) UpdatePlant(plant *models.Plant) error {
	return s.repo.Update(plant)
}

func (s *PlantService) DeletePlant(id uint) error {
	return s.repo.Delete(id)
}
