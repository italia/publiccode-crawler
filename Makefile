include .env

.PHONY: run build lint test up stop up-prod stop-prod

default: build

run:
	go run main.go version.go

build:
	docker build -t italia/${NAME}:${VERSION} --build-arg NAME=${NAME} --build-arg PROJECT=${PROJECT} .

lint:
	gometalinter --install
	gometalinter --exclude=vendor --exclude=middleware ./...

test:
	go test -race "${PROJECT}"/...

up:
	docker-compose up -d

stop:
	docker-compose stop

up-prod:
	docker-compose --file=docker-compose-prod.yml up -d

stop-prod:
	docker-compose --file=docker-compose-prod.yml stop
