package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response struct {
	Message string `json:"message"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 200}, nil
	res := Response{Message: "hello lambda"}
	return events.APIGatewayProxyResponse{Body: res.Message, StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
