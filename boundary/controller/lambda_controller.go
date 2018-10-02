package controller

import (
	"bytes"

	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/kutsuzawa/slim-load-recorder/usecase"
	"time"
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
	preRequest := struct {
		UserID   string  `json:"user_id"`
		Weight   float64 `json:"weight"`
		Distance float64 `json:"distance"`
		Date     string  `json:"date"`
		StartAt  string  `json:"start_at"`
		EndAt    string  `json:"end_at"`
	}{}
	bufRequestBody := bytes.NewBufferString(request.Body)
	if err := json.NewDecoder(bufRequestBody).Decode(&preRequest); err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusInternalServerError}, nil
	}
	date, err := sc.parseStrToTime(preRequest.Date)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusInternalServerError}, nil
	}
	startAt, err := sc.parseStrToTime(preRequest.StartAt)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusInternalServerError}, nil
	}
	endAt, err := sc.parseStrToTime(preRequest.EndAt)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusInternalServerError}, nil
	}
	slimLoadRequest := usecase.Request{
		UserID:   preRequest.UserID,
		Weight:   preRequest.Weight,
		Distance: preRequest.Distance,
		Date:     date,
		StartAt:  startAt,
		EndAt:    endAt,
	}
	response, err := sc.inputPort.Handle(slimLoadRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusInternalServerError}, nil
	}

	bufResponse := new(bytes.Buffer)
	if err := json.NewEncoder(bufResponse).Encode(response); err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	return events.APIGatewayProxyResponse{Body: bufResponse.String(), StatusCode: http.StatusOK}, nil
}

func (sc *slimLoadController) parseStrToTime(str string) (time.Time, error) {
	t, err := time.Parse("2006-01-02 15:04:05", str)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}
