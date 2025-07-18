package entity

// Adicione em um arquivo enums.go
type TipoTarefa string

const (
	Regar        TipoTarefa = "regar"
	Adubar       TipoTarefa = "adubar"
	Podar        TipoTarefa = "podar"
	Transplantar TipoTarefa = "transplantar"
	Monitorar    TipoTarefa = "monitorar"
)

func (t TipoTarefa) Valid() bool {
	switch t {
	case Regar, Adubar, Podar, Transplantar, Monitorar:
		return true
	}
	return false
}
