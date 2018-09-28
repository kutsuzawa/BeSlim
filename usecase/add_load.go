package usecase

import (
	"time"

	"github.com/kutsuzawa/slim-load-recorder/entity"
)

// Adapter makes below layer(Adapter layer) to implement two functions
type Adapter interface {
	AddLoad(userID string, load entity.Load) error
	GetLoadsByUserID(userID string, start, end time.Time) ([]entity.Load, error)
}

// LoadInteractor has Adapter Interface
type LoadInteractor struct {
	Adapter Adapter
}

// PostAndGetLoads add load data to db.
// Then, it get load data between start and end.
func (inteactor *LoadInteractor) PostAndGetLoads(userID string, load entity.Load, start, end time.Time) ([]entity.Load, error) {
	if err := inteactor.Adapter.AddLoad(userID, load); err != nil {
		return nil, err
	}
	loads, err := inteactor.Adapter.GetLoadsByUserID(userID, start, end)
	if err != nil {
		return nil, err
	}
	return loads, nil
}
