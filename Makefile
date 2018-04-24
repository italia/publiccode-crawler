include .env

.PHONY: run run-all run-version build lint test up stop prod-up prod-stop

default: build

run:
	go run main.go

run-all:
	go run main.go all

run-github:
	go run main.go github

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
	docker-compose pull --parallel
	docker-compose up -d --remove-orphans

stop:
	docker-compose stop

prod-up:
	docker-compose --file=docker-compose-prod.yml pull --parallel
	docker-compose --file=docker-compose-prod.yml up -d --remove-orphans

prod-stop:
	docker-compose --file=docker-compose-prod.yml stop
