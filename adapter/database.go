package adapter

import (
	"time"

	"github.com/kutsuzawa/slim-load-recorder/entity"
)

// DatabaseDriver is the interface that wraps methods for operating db
type DatabaseDriver interface {
	Add(userID string, load entity.Load) error
	Search(userID string, start, end time.Time) ([]entity.Load, error)
}

// Adapt has DatabaseDriver interface
type Adapt struct {
	DatabaseDriver DatabaseDriver
}

// AddLoad add load struct to db
func (adapt *Adapt) AddLoad(userID string, load entity.Load) error {
	if err := adapt.DatabaseDriver.Add(userID, load); err != nil {
		return err
	}
	return nil
}

// GetLoadsByUserID search load data.
func (adapt *Adapt) GetLoadsByUserID(userID string, start, end time.Time) ([]entity.Load, error) {
	loads, err := adapt.DatabaseDriver.Search(userID, start, end)
	if err != nil {
		return nil, err
	}
	return loads, nil
}
