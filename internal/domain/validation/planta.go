package validation

import (
	"errors"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/models"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/repository"
)

// internal/domain/validation/plant_validator.go
type PlantValidator struct {
	repo repository.PlantaRepository
}

func (v *PlantValidator) Validate(dto models.Planta) error {
	if dto.Nome == "" {
		return errors.New("nome da planta não pode ser vazio")
	}
	if dto.Especie == "" {
		return errors.New("espécie da planta não pode ser vazia")

	}
	return nil
}
