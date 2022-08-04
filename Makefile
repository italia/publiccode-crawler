.PHONY: build test

default: build

build:
	go build -ldflags "-X github.com/italia/developers-italia-backend/crawler/version.VERSION=$(shell git describe --abbrev=0 --tags)" -o bin/crawler

test:
	go test -race ./...
