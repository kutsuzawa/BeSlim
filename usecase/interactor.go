package usecase

import (
	"time"

	"github.com/kutsuzawa/slim-load-recorder/entity"
	"go.uber.org/zap"
)

// DataAccessor is accessed by Repository
type DataAccessor interface {
	AddLoad(userID string, load entity.Result) error
	GetLoadsByUserID(userID string, start, end time.Time) ([]entity.Load, error)
}

type addLoadFromLine struct {
	outputPort OutputPort
	repository DataAccessor
	logger     *zap.Logger
}

// NewAddLoadFromLine init addLoadFromLine
func NewAddLoadFromLine(outputPort OutputPort, repository DataAccessor, logger *zap.Logger) InputPort {
	return &addLoadFromLine{
		outputPort: outputPort,
		repository: repository,
		logger:     logger,
	}
}

func (u *addLoadFromLine) Handle(request Request) (Response, error) {
	user := entity.User{
		ID: request.UserID,
	}
	load := entity.Result{
		Date:     request.Date,
		Weight:   request.Weight,
		Distance: request.Distance,
	}
	if err := u.repository.AddLoad(user.ID, load); err != nil {
		u.logger.Error("failed to add load data", zap.String("error", err.Error()))
	}

	loads, err := u.repository.GetLoadsByUserID(user.ID, request.StartAt, request.EndAt)
	if err != nil {
		u.logger.Error("failed to get load data", zap.String("error", err.Error()))
	}

	return u.outputPort.Handle(loads)
}
