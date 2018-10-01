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

// DataAccessor is accessed by Repository
type DataAccessor interface {
	AddLoad(userID string, load entity.Load) error
	GetLoadsByUserID(userID string, start, end time.Time) ([]entity.Load, error)
}

type Receiver interface {
	GetRequest() (Request, error)
}

type Request struct {
	UserID   string    `json:"user_id"`
	Weight   float64   `json:"weight"`
	Distance float64   `json:"distance"`
	Date     time.Time `json:"date"`
	StartAt  time.Time `json:"start_at"`
	EndAt    time.Time `json:"end_at"`
}

type Sender interface {
	Post([]entity.Load) error
}

// GraphGeneration has DataAccessor Interface
type GraphGeneration struct {
	repository DataAccessor
	receiver   Receiver
	sender     Sender
}

// Run add load data to db.
// Then, it get load data between start and end.
func (gg *GraphGeneration) Run() error {
	rec, err := gg.receiver.GetRequest()
	if err != nil {
		return err
	}
	load := entity.Load{
		Weight:   rec.Weight,
		Distance: rec.Distance,
		Date:     rec.Date,
	}
	if err := gg.repository.AddLoad(rec.UserID, load); err != nil {
		return err
	}
	loads, err := gg.repository.GetLoadsByUserID(rec.UserID, rec.StartAt, rec.EndAt)
	if err != nil {
		return err
	}
	if err := gg.sender.Post(loads); err != nil {
		return err
	}
	return nil
}
