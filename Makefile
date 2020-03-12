.PHONY: up stop crawl

up:
  docker-compose --file=docker-compose-es-searchguard.yml up -d --remove-orphans

stop:
  docker-compose stop

crawl:
  docker-compose --file=docker-compose-es-searchguard.yml up -d
