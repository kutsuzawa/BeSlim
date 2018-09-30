package main

import (
	"log"
	"os"

	"github.com/kutsuzawa/slim-load-recorder/boundary/controller"
	"github.com/kutsuzawa/slim-load-recorder/boundary/presenter"
	"github.com/kutsuzawa/slim-load-recorder/boundary/repository"
	"github.com/kutsuzawa/slim-load-recorder/infrastructure/console"
	"github.com/kutsuzawa/slim-load-recorder/infrastructure/firebase"
	"github.com/kutsuzawa/slim-load-recorder/infrastructure/lambda"
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

	viewer := console.NewConsoleView(os.Stdout, os.Stdin)
	responder := presenter.NewConsolePresenter(viewer)
	requester := interactor.NewAddLoadFromLine(responder, repo, logger)
	ctl := controller.NewSlimLoadController(requester)
	receiver := lambda.NewLambda(ctl)
	receiver.Receive()
}
