build: compile

compile:
	go build -o bin/rubbernecker

lint:
	go fmt ./...
	go vet ./...
	golint . pkg/...

test:
	go run github.com/onsi/ginkgo/v2/ginkgo -r
