package main

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kutsuzawa/slim-load-recorder/client"
	"go.uber.org/zap"
)

// Response is format for response from lambda
type Response struct {
	//Data    HealthData `json:"data"`
	Message string `json:"message"`
}

// HealthData has weight and distance
type HealthData struct {
	Weight   float64 `json:"weight"`
	Distance float64 `json:"distance"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	// parse post data
	buf := bytes.NewBufferString(request.Body)
	data := &HealthData{}
	if err := json.NewDecoder(buf).Decode(data); err != nil {
		res := Response{Message: err.Error()}
		return events.APIGatewayProxyResponse{Body: res.Message, StatusCode: http.StatusInternalServerError}, nil
	}
	logger.Debug("post data",
		zap.Float64("weight", data.Weight),
		zap.Float64("distance", data.Distance),
	)

	// insert weight and distance data to DB
	db, err := client.NewDatabase()
	if err != nil {
		res := Response{Message: err.Error()}
		return events.APIGatewayProxyResponse{Body: res.Message, StatusCode: http.StatusInternalServerError}, nil
	}
	if err := db.AddUser("test_user", data.Weight, data.Distance); err != nil {
		res := Response{Message: err.Error()}
		return events.APIGatewayProxyResponse{Body: res.Message, StatusCode: http.StatusInternalServerError}, nil
	}
	res := Response{Message: "hello lambda"}
	return events.APIGatewayProxyResponse{Body: res.Message, StatusCode: http.StatusOK}, nil
}

func main() {
	lambda.Start(Handler)
}
