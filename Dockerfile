FROM golang:1.18 as build

WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags "-X github.com/italia/developers-italia-backend/crawler/version.VERSION=$(shell git describe --abbrev=0 --tags)"

FROM alpine:3

COPY --from=build /src/developers-italia-backend /usr/local/bin/developers-italia-backend
CMD ["developers-italia-backend", "crawl"]
