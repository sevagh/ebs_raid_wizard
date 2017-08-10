VERSION := 0.3.0

build:
	@go build -ldflags "-X main.VERSION=$(VERSION)" .

install:
	@go install .

deps:
	@go get -u github.com/golang/dep
	@dep ensure

test: build
	@go fmt . 
	@go vet .
	@go test -v ./...

lint:
	@gofmt -w .
