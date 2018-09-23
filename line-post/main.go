package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kutsuzawa/slim-load-recorder/client"
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
	buf := bytes.NewBufferString(request.Body)
	log.Println(buf.String())
	data := &HealthData{}
	if err := json.NewDecoder(buf).Decode(data); err != nil {
		res := Response{Message: err.Error()}
		return events.APIGatewayProxyResponse{Body: res.Message, StatusCode: http.StatusInternalServerError}, nil
	}
	log.Printf("[DEBUG] weight: %f, distance: %f\n", data.Weight, data.Distance)
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
