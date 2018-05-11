# Populate the ElasticSearch index

In this directory there is a php script that generate a handfull of documents and puts them in elasticsearch.

## Populate the Elasticsearch index

* copy `config.inc.dist` to `config.inc`
```
$ cp config.inc.dist config.inc
```
* edit `config.inc` and set the actual Elasticsearch connection parameters.
* connect to php container and move to the php directory
```
$ docker exec -it developers-italia-backend_php /bin/bash
$ cd /var/www/php/
```
* install the Elasticsearch php library
```
$ composer install
```
* run the script to insert documents
```
$ php insertDocument.php
```
