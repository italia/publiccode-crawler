FROM golang:1.18 as build

WORKDIR /src
COPY . .
RUN go build -ldflags "-s -w -X 'github.com/italia/publiccode-crawler/v3/internal.VERSION=$(git describe --abbrev=0 --tags)' -X 'github.com/italia/publiccode-crawler/v3/internal.BuildTime=$(date)'"

FROM alpine:3

RUN apk add --no-cache gcompat \
    && mkdir -p /var/crawler/data

COPY --from=build /src/publiccode-crawler /usr/local/bin/publiccode-crawler
CMD ["publiccode-crawler", "crawl"]
