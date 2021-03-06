.PHONY: setup
setup:
	dep ensure

.PHONY: build
build:
	env GOOS=linux go build -o bin/loads cmd/loads/main.go

.PHONY: lambda
lambda:
	make build
	sls deploy

.PHONY: local
local:
	make build
	env APP_ENV=local sam local start-api
