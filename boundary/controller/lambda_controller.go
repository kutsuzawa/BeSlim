package controller

import (
	"bytes"

	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/kutsuzawa/slim-load-recorder/usecase"
)

// SlimLoadController define method
type SlimLoadController interface {
	Run(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}

type slimLoadController struct {
	inputPort usecase.InputPort
}

// NewSlimLoadController init slimLoadController
func NewSlimLoadController(inputPort usecase.InputPort) SlimLoadController {
	return &slimLoadController{
		inputPort: inputPort,
	}
}

func (sc *slimLoadController) Run(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	bufRequestBody := bytes.NewBufferString(request.Body)
	var slimLoadRequest usecase.Request
	if err := json.NewDecoder(bufRequestBody).Decode(&slimLoadRequest); err != nil {
		return events.APIGatewayProxyResponse{}, nil
	}
	response, err := sc.inputPort.Handle(slimLoadRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{}, nil
	}

	bufResponse := new(bytes.Buffer)
	if err := json.NewEncoder(bufResponse).Encode(response); err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	return events.APIGatewayProxyResponse{Body: bufResponse.String(), StatusCode: http.StatusOK}, nil
}
