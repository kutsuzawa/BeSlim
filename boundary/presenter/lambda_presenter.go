package presenter

import (
	"github.com/kutsuzawa/slim-load-recorder/entity"
	"github.com/kutsuzawa/slim-load-recorder/usecase"
)

type lambdaPresenter struct {
}

// NewLambdaPresenter init lambdaPresenter
func NewLambdaPresenter() usecase.OutputPort {
	return &lambdaPresenter{}
}

func (lp *lambdaPresenter) Handle(loads []entity.Load) (usecase.Response, error) {
	var response usecase.Response
	for _, l := range loads {
		r := struct {
			UserID   string  `json:"user_id"`
			Weight   float64 `json:"weight"`
			Distance float64 `json:"distance"`
			Date     string  `json:"date"`
		}{UserID: l.User.ID, Weight: l.Result.Weight, Distance: l.Result.Distance, Date: l.Result.Date.String()}

		response = append(response, r)
	}
	return response, nil
}
