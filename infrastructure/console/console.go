package console

import (
	"fmt"
	"io"

	"github.com/kutsuzawa/slim-load-recorder/boundary/presenter"
)

type consoleView struct {
	outStream io.Writer
	errStream io.Writer
}

// NewConsoleView init console view
func NewConsoleView(outStream, errStream io.Writer) presenter.ConsoleViewer {
	return &consoleView{
		outStream: outStream,
		errStream: errStream,
	}
}

func (cv *consoleView) View(models []presenter.ConsoleModel) {
	for _, m := range models {
		fmt.Fprint(cv.outStream, "==============================\n")
		fmt.Fprintf(cv.outStream, "user_id: %s\n", m.UserID)
		fmt.Fprintf(cv.outStream, "date: %s\n", m.Date)
		fmt.Fprintf(cv.outStream, "weight: %f\n", m.Weight)
		fmt.Fprintf(cv.outStream, "distance: %f\n", m.Distance)
	}
}
