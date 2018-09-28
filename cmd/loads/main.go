package main

import (
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kutsuzawa/slim-load-recorder/adapter"
	"github.com/kutsuzawa/slim-load-recorder/driver"
	"github.com/kutsuzawa/slim-load-recorder/usecase"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	db, err := driver.NewFirebase(os.Getenv("APP_ENV"), os.Getenv("REGION"), os.Getenv("BUCKET"), os.Getenv("KEY"))
	if err != nil {
		log.Fatal(err)
	}
	repository := &adapter.Adapt{
		DatabaseDriver: db,
	}
	usecase := &usecase.LoadInteractor{
		Adapter: repository,
	}

	handler := &driver.Handle{
		Logger:  logger,
		Usecase: usecase,
	}
	lambda.Start(handler.ServeHTTP)
}
