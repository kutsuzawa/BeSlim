package repository

import (
	"time"

	"github.com/kutsuzawa/slim-load-recorder/entity"
)

// Driver is the interface that wraps methods for operating db
type Driver interface {
	Add(userID string, load entity.Load) error
	Search(userID string, start, end time.Time) ([]entity.Load, error)
}

// Repository has a Driver interface
type Repository struct {
	Driver Driver
}

// AddLoad add load structure to db
func (repo *Repository) AddLoad(userID string, load entity.Load) error {
	if err := repo.Driver.Add(userID, load); err != nil {
		return err
	}
	return nil
}

// GetLoadsByUserID search load data.
func (repo *Repository) GetLoadsByUserID(userID string, start, end time.Time) ([]entity.Load, error) {
	loads, err := repo.Driver.Search(userID, start, end)
	if err != nil {
		return nil, err
	}
	return loads, nil
}
