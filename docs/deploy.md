## Deploy architecture

**Deploy Architecture: containers**

Docker-compose

- elasticsearch:khezen/elasticsearch:6.2.2
  Service elasticsearch: contains all the data saved by the crawler and used from the website developers.italia.it
  Elasticsearch is a distributed, RESTful search and analytics engine.

- kibana: khezen/kibana:6.2.2
  Kibana is the GUI for elasticsearch. Not exposed and not used by the crawler.

- prometheus: quay.io/prometheus/prometheus:v2.2.1
  Prometheus offers a monitoring system and time series database. Used in order to save and review the metrics for the crawler service.

- proxy: containous/traefik:experimental
  Træfik is the reverse proxy that serves all the containers over the different ports.

- crawler: italia/developers-italia-backend:0.0.1
  The crawler that will visit a repository, validate, download and save the publiccode.yml file.
