package interactor

import (
	"time"

	"github.com/kutsuzawa/slim-load-recorder/entity"
)

// Adapter makes below layer(Adapter layer) to implement two functions
//type Adapter interface {
//	AddLoad(userID string, load entity.Load) error
//	GetLoadsByUserID(userID string, start, end time.Time) ([]entity.Load, error)
//}

// DataAccesser is accessed by Repository
type DataAccesser interface {
	AddLoad(userID string, load entity.Load) error
	GetLoadsByUserID(userID string, start, end time.Time) ([]entity.Load, error)
}

// Interactor has DataAccesser Interface
type Interactor struct {
	Repository DataAccesser
}

// LoadGenerate add load data to db.
// Then, it get load data between start and end.
func (inteactor *Interactor) LoadGenerate(userID string, load entity.Load, start, end time.Time) ([]entity.Load, error) {
	if err := inteactor.Repository.AddLoad(userID, load); err != nil {
		return nil, err
	}
	loads, err := inteactor.Repository.GetLoadsByUserID(userID, start, end)
	if err != nil {
		return nil, err
	}
	return loads, nil
}
