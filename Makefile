RUBBERNECKER_DIR=$(shell pwd)

DOCKER         = $(shell command -v docker)
DOCKER_COMPOSE = $(shell command -v docker-compose)

GO     := $(if $(shell command -v go),go,docker run --rm -v $(RUBBERNECKER_DIR):$(RUBBERNECKER_DIR) -w $(RUBBERNECKER_DIR) golang:1.9 go)
GOLINT := $(if $(shell command -v golint),golint,$(GO) get -u github.com/golang/lint/golint && golint)
SASS   := $(if $(shell command -v sass),sass,docker run --rm -v $(RUBBERNECKER_DIR):$(RUBBERNECKER_DIR) -w $(RUBBERNECKER_DIR) ubuntudesign/sass sass)
TSC    := $(if $(shell command -v tsc),tsc,docker run --rm -v $(RUBBERNECKER_DIR):$(RUBBERNECKER_DIR) -w $(RUBBERNECKER_DIR) sandrokeil/typescript tsc)

build: scripts styles compile

check_evn:
	$(if $(DOCKER),,$(error "docker not found in PATH"))

compile: check_evn
	$(GO) build -o bin/rubbernecker

lint: check_evn
	$(GO) fmt ./...
	$(GO) vet ./...
	$(GOLINT) . pkg/...

scripts: check_evn
	$(TSC) build/ts/app.ts --outFile dist/app.js --removeComments

styles: check_evn
	$(SASS) build/scss/app.scss > dist/app.css --style compressed --no-cache

test:
ifndef TRAVIS
	$(if $(DOCKER_COMPOSE),,$(error "docker-compose not found in PATH"))
	docker-compose run rubbernecker go test -v ./...
else
	go test -v ./...
endif
