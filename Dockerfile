FROM golang:1.25@sha256:cbff9d1a9041b316010f2da6b701b6c0d597718cb90928c85eb597334a0d23d4 AS build

WORKDIR /src
COPY . .
RUN go build -ldflags "-s -w -X 'github.com/italia/publiccode-crawler/v4/internal.VERSION=$(git describe --abbrev=0 --tags)' -X 'github.com/italia/publiccode-crawler/v4/internal.BuildTime=$(date)'"

FROM alpine:3@sha256:28bd5fe8b56d1bd048e5babf5b10710ebe0bae67db86916198a6eec434943f8b

RUN apk add --no-cache gcompat \
    && mkdir -p /var/crawler/data

COPY --from=build /src/publiccode-crawler /usr/local/bin/publiccode-crawler
CMD ["publiccode-crawler", "crawl"]
