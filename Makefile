include .env

.PHONY: up stop crawl

up:
	docker-compose up -d --remove-orphans

stop:
	docker-compose stop

crawl:
	docker-compose --file=docker-compose-crawler.yml up -d
