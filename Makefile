#!/usr/bin/make -f

.ONESHELL:
.SHELL := /usr/bin/bash

PROJECTNAME := $(shell basename "$$(pwd)")
PROJECTPATH := $(shell pwd)
COMPOSE_SERVICE_NAME = "gocallbacksvc"
GOFLAGS :=

help:
	echo "Usage: make [options] [arguments]\n"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

go-build: ## Compiles packages and dependencies. Builds a binary for the callback-service under bin/. *Accepts GOFLAGS and LDFLAGS
	@[ -d bin ] || mkdir bin
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go build $(GOFLAGS) -o bin/ "$(PROJECTPATH)/cmd/..."

go-run: ## Starts callback-service project. *Accepts GOFLAGS and LDFLAGS
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go run $(GOFLAGS) "$(PROJECTPATH)/cmd/callback-service/main.go" $(LDFLAGS)

client-start: ## Starts client-service project. *Accepts GOFLAGS and LDFLAGS
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go run $(GOFLAGS) "$(PROJECTPATH)/cmd/client-service/main.go" $(LDFLAGS)

go-doc: ## Generates static docs
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) godoc -http=localhost:6060

go-vendor: ## Updates vendor dependencies
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go mod vendor && go mod tidy

test: ## Runs the tests
	docker start postgres_test || docker compose up postgres_test -d
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go test ./...

test-integration: ## Runs the tests
	docker start postgres_test || docker compose up postgres_test -d
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go test ./... -tags=integration

docker-build: ## Builds the project binary inside a docker image
	@docker build -t $(PROJECTNAME) .

docker-run:	## Runs the previosly build docker image
	@docker run $(PROJECTNAME) -p 9090:9090

compose-build: ## Builds the compose image of the callback-service
	docker compose --file "$(PROJECTPATH)/docker-compose.yaml" build $(COMPOSE_SERVICE_NAME)

compose-mac: ## Starts docker-compose project for a Mac OS (Darwin arch). Accepts LDFLAGS
	docker compose up -d zipkin postgres gocallbackclient
	TRACE_URL=$$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' zipkin) ;\
	docker compose --file "$(PROJECTPATH)/docker-compose.yaml" run --service-ports -e ENV_ARGS="--database-host=docker.for.mac.localhost --trace-url=http://$$TRACE_URL:9411/api/v2/spans $(LDFLAGS)" $(COMPOSE_SERVICE_NAME) -p $(COMPOSE_SERVICE_NAME)

compose-linux: ## Starts docker-compose project for a Linux OS
	docker compose up -d zipkin postgres gocallbackclient
	docker compose --file "$(PROJECTPATH)/docker-compose.yaml" --verbose up $(COMPOSE_SERVICE_NAME)

compose-client: ## Starts docker-compose callback client
	docker compose up gocallbackclient
