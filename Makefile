.PHONY: setup
setup:
	dep ensure

.PHONY: build
build:
	env GOOS=linux go build -o bin/line-post line-post/main.go

.PHONY: lambda
lambda:
	sls deploy

.PHONY: local
local:
	env APP_ENV=local sam local start-api
