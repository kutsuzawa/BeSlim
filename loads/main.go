package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

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

// Recieve is LINE request.
// 一時的にこの形
// TODO: 本当のリクエストの形に合わせて宣言
type Recieve struct {
	UserID   string  `json:"user_id"`
	Weight   float64 `json:"weight"`
	Distance float64 `json:"distance"`
	Date     string  `json:"date"`
	Start_at string  `json:"start_at"`
	End_at   string  `json:"end_at"`
}

type handler struct {
	logger *zap.Logger
}

func (h *handler) ServeHTTP(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	rec, err := parseRequest(request)
	if err != nil {
		h.logger.Error("parse error", zap.String("error", err.Error()))
		res := Response{Message: err.Error()}
		return events.APIGatewayProxyResponse{Body: res.Message, StatusCode: http.StatusInternalServerError}, nil
	}
	h.logger.Info("post data",
		zap.String("user_id", rec.UserID),
		zap.Float64("weight", rec.Weight),
		zap.Float64("distance", rec.Distance),
		zap.Time("date", parseTimeStr(rec.Date)),
		zap.Time("start_at", parseTimeStr(rec.Start_at)),
		zap.Time("end_at", parseTimeStr(rec.End_at)),
	)

	db, err := client.NewDatabase()
	if err != nil {
		res := Response{Message: err.Error()}
		return events.APIGatewayProxyResponse{Body: res.Message, StatusCode: http.StatusInternalServerError}, nil
	}
	if err := db.AddLoad(rec.UserID, rec.Weight, rec.Distance, parseTimeStr(rec.Date)); err != nil {
		res := Response{Message: err.Error()}
		return events.APIGatewayProxyResponse{Body: res.Message, StatusCode: http.StatusInternalServerError}, nil
	}
	results, err := db.GetDataByUserID(rec.UserID, parseTimeStr(rec.Start_at), parseTimeStr(rec.End_at))
	if err != nil {
		res := Response{Message: err.Error()}
		return events.APIGatewayProxyResponse{Body: res.Message, StatusCode: http.StatusInternalServerError}, nil
	}
	message, err := encodeResults(results)
	if err != nil {
		res := Response{Message: err.Error()}
		return events.APIGatewayProxyResponse{Body: res.Message, StatusCode: http.StatusInternalServerError}, nil
	}
	res := Response{Message: message}
	return events.APIGatewayProxyResponse{Body: res.Message, StatusCode: http.StatusOK}, nil
}

func parseRequest(request events.APIGatewayProxyRequest) (Recieve, error) {
	buf := bytes.NewBufferString(request.Body)
	rec := Recieve{}
	if err := json.NewDecoder(buf).Decode(&rec); err != nil {
		return Recieve{}, err
	}
	return rec, nil
}

func parseTimeStr(timeStr string) time.Time {
	t, _ := time.Parse("2006-01-02 15:04:05", timeStr)
	return t
}

func encodeResults(results []client.Load) (string, error) {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(results); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	handler := handler{logger}
	lambda.Start(handler.ServeHTTP)
}
