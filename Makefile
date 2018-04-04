RUBBERNECKER_DIR=$(shell pwd)

DOCKER         = $(shell command -v docker)
DOCKER_COMPOSE = $(shell command -v docker-compose)

GO     := $(if $(shell command -v go),go,docker run --rm -v $(RUBBERNECKER_DIR):$(RUBBERNECKER_DIR) -w $(RUBBERNECKER_DIR) golang:1.9 go)
GOLINT := $(if $(shell command -v golint),golint,$(GO) get -u github.com/golang/lint/golint && golint)
NPM    := $(if $(shell command -v npm),npm,docker run --rm -v $(RUBBERNECKER_DIR):$(RUBBERNECKER_DIR) -w $(RUBBERNECKER_DIR) node:carbon-alpine npm)

build: scripts styles compile

check_evn:
	$(if $(DOCKER),,$(error "docker not found in PATH"))

compile: check_evn
	$(GO) build -o bin/rubbernecker

dependencies:
	$(NPM) install

lint: check_evn
	$(GO) fmt ./...
	$(GO) vet ./...
	$(GOLINT) . pkg/...
	$(NPM) run tslint

scripts: check_evn
	$(NPM) run tsc

styles: check_evn
	$(NPM) run sass

test:
ifndef TRAVIS
	$(if $(DOCKER_COMPOSE),,$(error "docker-compose not found in PATH"))
	docker-compose run rubbernecker go test -v ./...
else
	go test -v ./...
endif
