package main

import (
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kutsuzawa/slim-load-recorder/boundary/repository"
	"github.com/kutsuzawa/slim-load-recorder/infrastructure/firebase"
	"github.com/kutsuzawa/slim-load-recorder/infrastructure/handler"
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
	repo := &repository.Repository{
		Driver: db,
	}
	usecase := &interactor.Interactor{
		Repository: repo,
	}

	h := &handler.Handle{
		Logger:  logger,
		Usecase: usecase,
	}
	lambda.Start(h.ServeHTTP)
}
