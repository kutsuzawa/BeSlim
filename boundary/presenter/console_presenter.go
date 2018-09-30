package presenter

import "github.com/kutsuzawa/slim-load-recorder/interactor"

// ConsoleModel is used to output to console
type ConsoleModel struct {
	UserID   string  `json:"user_id"`
	Date     string  `json:"date"`
	Weight   float64 `json:"weight"`
	Distance float64 `json:"distance"`
}

type consolePresenter struct {
	viewer ConsoleViewer
}

// NewConsolePresenter init consolePresenter
func NewConsolePresenter(viewer ConsoleViewer) interactor.SlimLoadResponder {
	return &consolePresenter{
		viewer: viewer,
	}
}

func (cp *consolePresenter) Present(response []interactor.Response) {
	var models []ConsoleModel
	for _, r := range response {
		model := ConsoleModel{
			UserID:   r.UserID,
			Date:     r.Date,
			Weight:   r.Weight,
			Distance: r.Distance,
		}
		models = append(models, model)
	}
	cp.viewer.View(models)
}
