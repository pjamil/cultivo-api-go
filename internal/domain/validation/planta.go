package validation

import (
	"errors"

	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/entity"
	"gitea.paulojamil.dev.br/paulojamil.dev.br/cultivo-api-go/internal/domain/repository"
)

// internal/domain/validation/plant_validator.go
type PlantValidator struct {
	repo repository.PlantaRepositorio
}

func (v *PlantValidator) Validate(dto entity.Planta) error {
	if dto.Nome == "" {
		return errors.New("nome da planta não pode ser vazio")
	}
	if dto.Especie == "" {
		return errors.New("espécie da planta não pode ser vazia")

	}
	return nil
}
