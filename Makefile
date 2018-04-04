RUBBERNECKER_DIR=$(shell pwd)

DOCKER         = $(shell command -v docker)
DOCKER_COMPOSE = $(shell command -v docker-compose)

GO     := $(if $(shell command -v go),go,docker run --rm -v $(RUBBERNECKER_DIR):$(RUBBERNECKER_DIR) -w $(RUBBERNECKER_DIR) golang:1.9 go)
GOLINT := $(if $(shell command -v golint),golint,$(GO) get -u github.com/golang/lint/golint && golint)
NPM    := $(if $(shell command -v npm),npm,docker run --rm -v $(RUBBERNECKER_DIR):$(RUBBERNECKER_DIR) -w $(RUBBERNECKER_DIR) node:carbon-alpine npm)
DEP    := $(if $(shell command -v dep),dep,docker run --rm -v $(RUBBERNECKER_DIR):$(RUBBERNECKER_DIR) -w $(RUBBERNECKER_DIR) golang:1.9 dep)

build: scripts styles compile

check_evn:
	$(if $(DOCKER),,$(error "docker not found in PATH"))

compile:
	$(GO) build -o bin/rubbernecker

dependencies:
	$(NPM) install
	$(DEP) ensure

lint:
	$(GO) fmt ./...
	$(GO) vet ./...
	$(GOLINT) . pkg/...
	$(NPM) run tslint

scripts:
	$(NPM) run tsc

styles:
	$(NPM) run sass

test:
ifndef TRAVIS
	$(if $(DOCKER_COMPOSE),,$(error "docker-compose not found in PATH"))
	docker-compose run rubbernecker go test -v ./...
else
	go test -v ./...
endif
