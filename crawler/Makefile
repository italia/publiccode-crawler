.PHONY: build lint test

default: build

build:
	go build -ldflags "-X github.com/italia/developers-italia-backend/crawler/version.VERSION=$(shell git describe --abbrev=0 --tags)" -o bin/crawler
	chmod +x bin/crawler

lint:
	gometalinter --install
	gometalinter --exclude=vendor ./...

test:
	go test -race ./...
