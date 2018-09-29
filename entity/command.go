package entity

// Usecase interface define usecase
type Usecase interface {
	GetLoadGraph(Load) error
}
