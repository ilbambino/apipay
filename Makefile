SHA := $(shell git rev-parse --short HEAD)

all: build

test:
	go get golang.org/x/tools/cmd/cover
	go get github.com/stretchr/testify/assert
	go test -race -timeout 3m  ./...


build:
	go build -ldflags "-X main.gitHash=$(SHA)"

clean:
	rm apipay
	go clean -testcache
