.PHONY: run build lint test

default: build

NAME=developers-italia-backend
PROJECT?=github.com/italia/developers-italia-backend
VERSION?=0.0.1

run:
	go run main.go version.go

build: lint
	go build -ldflags "-X main.Version=${VERSION}" -o ${NAME} "${PROJECT}"

lint:
	gometalinter --install
	gometalinter --exclude=vendor --exclude=middleware ./...

test:
	go test -race "${PROJECT}"/...
