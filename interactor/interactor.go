package interactor

import (
	"time"

	"github.com/kutsuzawa/slim-load-recorder/entity"
	"go.uber.org/zap"
)

// DataAccessor is accessed by Repository
type DataAccessor interface {
	AddLoad(userID string, load entity.Load) error
	GetLoadsByUserID(userID string, start, end time.Time) ([]entity.Load, error)
}

type addLoadFromLine struct {
	responder  SlimLoadResponder
	repository DataAccessor
	logger     *zap.Logger
}

// NewAddLoadFromLine init addLoadFromLine
func NewAddLoadFromLine(responder SlimLoadResponder, repository DataAccessor, logger *zap.Logger) SlimLoadRequester {
	return &addLoadFromLine{
		responder:  responder,
		repository: repository,
		logger:     logger,
	}
}

func (u *addLoadFromLine) Handle(request Request) {
	t, err := request.date()
	if err != nil {
		u.logger.Error("error occurred when parsing date", zap.String("error", err.Error()))
	}
	load := entity.Load{
		Date:     t,
		Weight:   request.Weight,
		Distance: request.Distance,
	}
	if err := u.repository.AddLoad(request.UserID, load); err != nil {
		u.logger.Error("failed to add load data", zap.String("error", err.Error()))
	}

	startAt, err := request.startAt()
	if err != nil {
		u.logger.Error("error occurred when parsing startAt", zap.String("error", err.Error()))
	}
	endAt, err := request.endAt()
	if err != nil {
		u.logger.Error("error occurred when parsing endAt", zap.String("error", err.Error()))
	}
	loads, err := u.repository.GetLoadsByUserID(request.UserID, startAt, endAt)
	if err != nil {
		u.logger.Error("failed to get load data", zap.String("error", err.Error()))
	}

	var responses []Response
	for _, l := range loads {
		response := Response{
			UserID:   request.UserID,
			Weight:   l.Weight,
			Distance: l.Distance,
			Date:     l.Date.String(),
		}
		responses = append(responses, response)
	}
	u.responder.Present(responses)
}
