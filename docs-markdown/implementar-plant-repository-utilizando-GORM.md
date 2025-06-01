# Implementação do PlantRepository utilizando GORM

Vou implementar o `PlantRepository` seguindo a interface definida anteriormente, utilizando GORM como ORM para interação com o banco de dados. Esta implementação será concreta e será usada pela camada de serviço.

## Implementação Completa do PlantRepository

`internal/infrastructure/database/plant_repository.go`:

```go
package database

import (
 "errors"

 "gorm.io/gorm"

 "github.com/yourusername/plant-cultivation-api/internal/domain/models"
)

// PlantRepository implementa a interface repository.PlantRepository
type PlantRepository struct {
 db *gorm.DB
}

// NewPlantRepository cria uma nova instância do PlantRepository
func NewPlantRepository(db *gorm.DB) *PlantRepository {
 return &PlantRepository{db: db}
}

// Create insere uma nova planta no banco de dados
func (r *PlantRepository) Create(plant *models.Plant) error {
 if plant == nil {
  return errors.New("plant cannot be nil")
 }

 result := r.db.Create(plant)
 if result.Error != nil {
  return result.Error
 }

 return nil
}

// FindAll retorna todas as plantas cadastradas
func (r *PlantRepository) FindAll() ([]models.Plant, error) {
 var plants []models.Plant

 result := r.db.Find(&plants)
 if result.Error != nil {
  return nil, result.Error
 }

 return plants, nil
}

// FindByID busca uma planta pelo ID
func (r *PlantRepository) FindByID(id uint) (*models.Plant, error) {
 var plant models.Plant

 result := r.db.First(&plant, id)
 if result.Error != nil {
  if errors.Is(result.Error, gorm.ErrRecordNotFound) {
   return nil, nil // Retorna nil sem erro quando não encontrado
  }
  return nil, result.Error
 }

 return &plant, nil
}

// Update atualiza uma planta existente
func (r *PlantRepository) Update(plant *models.Plant) error {
 if plant == nil {
  return errors.New("plant cannot be nil")
 }

 // Verifica se a planta existe
 existingPlant, err := r.FindByID(plant.ID)
 if err != nil {
  return err
 }
 if existingPlant == nil {
  return errors.New("plant not found")
 }

 result := r.db.Save(plant)
 if result.Error != nil {
  return result.Error
 }

 return nil
}

// Delete remove uma planta pelo ID
func (r *PlantRepository) Delete(id uint) error {
 result := r.db.Delete(&models.Plant{}, id)
 if result.Error != nil {
  return result.Error
 }

 // Verifica se algum registro foi realmente deletado
 if result.RowsAffected == 0 {
  return errors.New("no plant found with the given ID")
 }

 return nil
}

// FindBySpecies retorna plantas por espécie (método adicional)
func (r *PlantRepository) FindBySpecies(species models.PlantSpecies) ([]models.Plant, error) {
 var plants []models.Plant

 result := r.db.Where("species = ?", species).Find(&plants)
 if result.Error != nil {
  return nil, result.Error
 }

 return plants, nil
}

// FindByStatus retorna plantas por status (método adicional)
func (r *PlantRepository) FindByStatus(status string) ([]models.Plant, error) {
 var plants []models.Plant

 result := r.db.Where("status = ?", status).Find(&plants)
 if result.Error != nil {
  return nil, result.Error
 }

 return plants, nil
}
```

## Atualização da Interface

`internal/domain/repository/plant_repository.go` (atualizado):

```go
package repository

import "github.com/yourusername/plant-cultivation-api/internal/domain/models"

type PlantRepository interface {
 Create(plant *models.Plant) error
 FindAll() ([]models.Plant, error)
 FindByID(id uint) (*models.Plant, error)
 Update(plant *models.Plant) error
 Delete(id uint) error
 FindBySpecies(species models.PlantSpecies) ([]models.Plant, error) // Novo método
 FindByStatus(status string) ([]models.Plant, error)               // Novo método
}
```

## Explicação da Implementação

1. **Estrutura do Repositório**:
   - O `PlantRepository` contém uma instância do GORM DB para operações no banco
   - Todos os métodos seguem a interface definida

2. **Métodos CRUD**:
   - `Create`: Insere uma nova planta
   - `FindAll`: Retorna todas as plantas
   - `FindByID`: Busca uma planta específica pelo ID
   - `Update`: Atualiza uma planta existente
   - `Delete`: Remove uma planta pelo ID

3. **Métodos Adicionais**:
   - `FindBySpecies`: Filtra plantas por espécie
   - `FindByStatus`: Filtra plantas por status

4. **Tratamento de Erros**:
   - Verificação de nil pointers
   - Tratamento específico para "registro não encontrado"
   - Retorno de erro quando nenhum registro é afetado

5. **Boas Práticas**:
   - Validações básicas nos métodos
   - Uso adequado do GORM
   - Separação clara de responsabilidades

## Como Usar o Repositório

Aqui está um exemplo de como o repositório seria usado na camada de serviço:

```go
// Inicialização
db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
plantRepo := database.NewPlantRepository(db)

// Criar uma planta
newPlant := &models.Plant{
    Name:         "Tomate Cereja",
    Species:      models.SpeciesTomato,
    PlantingDate: "2023-05-15",
    Status:       "growing",
    GrowthStage:  "vegetative",
}
err := plantRepo.Create(newPlant)

// Buscar todas as plantas
plants, err := plantRepo.FindAll()

// Buscar por ID
plant, err := plantRepo.FindByID(1)

// Atualizar
plant.Name = "Tomate Cereja Orgânico"
err = plantRepo.Update(plant)

// Deletar
err = plantRepo.Delete(1)

// Buscar por espécie
tomatoes, err := plantRepo.FindBySpecies(models.SpeciesTomato)
```

## Próximas Melhorias Possíveis

1. **Paginação**: Adicionar suporte a paginação nos métodos FindAll
2. **Filtros Avançados**: Implementar um sistema de query builder para filtros complexos
3. **Cache**: Adicionar camada de cache para consultas frequentes
4. **Transações**: Implementar métodos que usam transações para operações complexas
5. **Logging**: Adicionar logging detalhado para operações do repositório

Esta implementação fornece uma base sólida para a persistência de dados no domínio de cultivo de plantas, seguindo os princípios do padrão Repository e permitindo fácil substituição da implementação concreta se necessário.
