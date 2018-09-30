package lambda

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	api "github.com/aws/aws-lambda-go/lambda"
	"github.com/kutsuzawa/slim-load-recorder/boundary/controller"
)

type lambda struct {
	ctl controller.SlimLoadController
}

// NewLambda init lambda
func NewLambda(ctl controller.SlimLoadController) controller.SlimLoadReceiver {
	return &lambda{
		ctl: ctl,
	}
}

func (l *lambda) Receive() {
	api.Start(l.handle)
}

func (l *lambda) handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	l.ctl.ParseRequest(request)
	return events.APIGatewayProxyResponse{Body: "ok", StatusCode: http.StatusOK}, nil
}
