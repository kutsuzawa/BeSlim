package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kutsuzawa/slim-load-recorder/client"
)

// Response is format for response from lambda
type Response struct {
	Message string `json:"message"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	db, err := client.NewDatabase()
	if err != nil {
		res := Response{Message: err.Error()}
		return events.APIGatewayProxyResponse{Body: res.Message, StatusCode: 500}, nil
	}
	db.AddUser("test_group", "test_user")
	res := Response{Message: "hello lambda"}
	return events.APIGatewayProxyResponse{Body: res.Message, StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
