package controller

import (
	"bytes"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gin-gonic/gin/json"
	"github.com/kutsuzawa/slim-load-recorder/interactor"
)

// SlimLoadController define method for parsing request
type SlimLoadController interface {
	ParseRequest(request events.APIGatewayProxyRequest)
}

// SlimLoadReceiver (DIP)
type SlimLoadReceiver interface {
	Receive()
}

type slimLoadController struct {
	interactor interactor.SlimLoadRequester
}

// NewSlimLoadController init slimLoadController
func NewSlimLoadController(interactor interactor.SlimLoadRequester) SlimLoadController {
	return &slimLoadController{
		interactor: interactor,
	}
}

func (sc *slimLoadController) ParseRequest(request events.APIGatewayProxyRequest) {
	buf := bytes.NewBufferString(request.Body)
	var slimLoadRequest interactor.Request
	if err := json.NewDecoder(buf).Decode(&slimLoadRequest); err != nil {
		return
	}
	sc.interactor.Handle(slimLoadRequest)
}
