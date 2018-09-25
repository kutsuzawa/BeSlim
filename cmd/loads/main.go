package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kutsuzawa/slim-load-recorder/adaptor"
	"github.com/kutsuzawa/slim-load-recorder/application"
	"go.uber.org/zap"
)

// Response is format for response from lambda
type Response struct {
	//Data    HealthData `json:"data"`
	Message string `json:"message"`
}

// Receive is LINE request.
// 一時的にこの形
// TODO: 本当のリクエストの形に合わせて宣言
type Receive struct {
	UserID   string  `json:"user_id"`
	Weight   float64 `json:"weight"`
	Distance float64 `json:"distance"`
	Date     string  `json:"date"`
	StartAt  string  `json:"start_at"`
	EndAt    string  `json:"end_at"`
}

type handler struct {
	logger  *zap.Logger
	factory adaptor.Factory
	config  *adaptor.ClientFactoryConfig
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
		zap.Time("start_at", parseTimeStr(rec.StartAt)),
		zap.Time("end_at", parseTimeStr(rec.EndAt)),
	)
	db, err := h.factory.Database(h.config)
	if err != nil {
		res := Response{Message: err.Error()}
		return events.APIGatewayProxyResponse{Body: res.Message, StatusCode: http.StatusInternalServerError}, nil
	}
	if err := db.AddLoad(rec.UserID, rec.Weight, rec.Distance, parseTimeStr(rec.Date)); err != nil {
		res := Response{Message: err.Error()}
		return events.APIGatewayProxyResponse{Body: res.Message, StatusCode: http.StatusInternalServerError}, nil
	}
	results, err := db.GetDataByUserID(rec.UserID, parseTimeStr(rec.StartAt), parseTimeStr(rec.EndAt))
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

func parseRequest(request events.APIGatewayProxyRequest) (Receive, error) {
	buf := bytes.NewBufferString(request.Body)
	rec := Receive{}
	if err := json.NewDecoder(buf).Decode(&rec); err != nil {
		return Receive{}, err
	}
	return rec, nil
}

func parseTimeStr(timeStr string) time.Time {
	t, _ := time.Parse("2006-01-02 15:04:05", timeStr)
	return t
}

func encodeResults(results []application.Load) (string, error) {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(results); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	config := &adaptor.ClientFactoryConfig{
		S3Region: os.Getenv("REGION"),
		S3Bucket: os.Getenv("BUCKET"),
		S3Key:    os.Getenv("KEY"),
	}
	handler := handler{
		logger:  logger,
		factory: adaptor.NewClientFactory(os.Getenv("APP_ENV")),
		config:  config,
	}
	lambda.Start(handler.ServeHTTP)
}
