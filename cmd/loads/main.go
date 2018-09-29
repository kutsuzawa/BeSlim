package main

import (
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kutsuzawa/slim-load-recorder/adapter/repository"
	"github.com/kutsuzawa/slim-load-recorder/driver/firebase"
	"github.com/kutsuzawa/slim-load-recorder/driver/handler"
	"github.com/kutsuzawa/slim-load-recorder/interactor"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	db, err := firebase.NewFirebase(os.Getenv("APP_ENV"), os.Getenv("REGION"), os.Getenv("BUCKET"), os.Getenv("KEY"))
	if err != nil {
		log.Fatal(err)
	}
	repository := &repository.Repository{
		Driver: db,
	}
	usecase := &interactor.Interactor{
		Repository: repository,
	}

	handler := &handler.Handle{
		Logger:  logger,
		Usecase: usecase,
	}
	lambda.Start(handler.ServeHTTP)
}
