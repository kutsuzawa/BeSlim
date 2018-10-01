package main

import (
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kutsuzawa/slim-load-recorder/boundary/controller"
	"github.com/kutsuzawa/slim-load-recorder/boundary/presenter"
	"github.com/kutsuzawa/slim-load-recorder/boundary/repository"
	"github.com/kutsuzawa/slim-load-recorder/infrastructure/firebase"
	"github.com/kutsuzawa/slim-load-recorder/usecase"
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

	outputPort := presenter.NewLambdaPresenter()
	inputPort := usecase.NewAddLoadFromLine(outputPort, repo, logger)
	ctl := controller.NewSlimLoadController(inputPort)
	lambda.Start(ctl.Run)
}
