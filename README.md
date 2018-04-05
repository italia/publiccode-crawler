# developers-italia-backend

[![CircleCI](https://circleci.com/gh/italia/developers-italia-backend/tree/master.svg?style=svg)](https://circleci.com/gh/italia/developers-italia-backend/tree/master)

Backend &amp; crawler for the OSS catalog of Developers Italia

**This document is a Work in progress**

## How to contribute

### Dependencies

* [Go](https://golang.org/)
* [dep](https://github.com/golang/dep)
* [Docker](https://www.docker.com/)
* [Docker-compose](https://docs.docker.com/compose/)

### Starting steps

* rename .env.example to .env and fill the variables with your values

  Default elastic user and password are `elastic`

  Default kibana user and password are `kibana`

  Basic Auth token is generated with `echo -n "user:password" | openssl base64 -base64`

* start the Docker stack: `docker-compose up -d`
* execute the application: `make run all`

### Copyright

```
Copyright (c) the respective contributors, as shown by the AUTHORS file.

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published
by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
```
