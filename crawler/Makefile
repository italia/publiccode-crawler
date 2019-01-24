include ../.env

.PHONY: run run-all run-version build lint test

default: build

run-all:
	go run -race main.go all

run-version:
	go run -race main.go version

build:
	docker build -t italia/${NAME}-crawler:${VERSION} \
	    --build-arg VERSION=${VERSION} \
	    .

lint:
	gometalinter --install
	gometalinter --exclude=vendor ./...

test:
	go test -race "${PROJECT}"/...
