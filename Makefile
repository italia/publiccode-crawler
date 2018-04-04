include .env

.PHONY: run run-all run-version build lint test up stop prod-up prod-stop

default: build

run:
	go run main.go

run-all:
	go run main.go all

run-version:
	go run main.go version

build:
	docker build -t italia/${NAME}:${VERSION} \
	    --build-arg NAME=${NAME} \
	    --build-arg PROJECT=${PROJECT} \
	    --build-arg VERSION=${VERSION} \
	    ./

lint:
	gometalinter --install
	gometalinter --exclude=vendor ./...

test:
	go test -race "${PROJECT}"/...

up:
	docker-compose up -d

stop:
	docker-compose stop

prod-up:
	docker-compose --file=docker-compose-prod.yml up -d

prod-stop:
	docker-compose --file=docker-compose-prod.yml stop
