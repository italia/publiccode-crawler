## Jekyll

Jekyll package create the 4 yaml files used by Jekyll plugin in developers.italia.it
in order to generate the pages and data structures of OSS catalog.

It uses ElasticSearch data store as source of truth.

The result files are saved into: `jekyll/generated/` (mounted as volume in docker-compose stack).

**File generation**

The four yaml files are:

- amministrazioni.yml
- softwares.yml
- software-riuso.yml
- software-open-source.yml

_amministrazioni.yml_

It contains all the Public Administrations with: name, url (website) and codice iPA.

_softwares.yml_

It contains all the software that the crawler scraped, validated and saved into ElasticSearch.
The structure is similar to publiccode data structure with some other fields like vitality and vitality score.

_software-riuso.yml_

It contains all the software in _softwares.yml_ with a not empty codiceIPA.

_software-open-source.yml_

It contains all the software in _softwares.yml_ with an empty codiceIPA.
