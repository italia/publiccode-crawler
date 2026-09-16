FROM golang:1.27@sha256:f44f6e88636cfb311f9ebace870ded69d943f227bb3cb27d32ffd84ea18c43ea AS build

WORKDIR /src
COPY . .
RUN go build -trimpath -ldflags "-s -w -buildid= -X 'github.com/italia/publiccode-crawler/v4/internal.VERSION=$(git describe --abbrev=0 --tags)' -X 'github.com/italia/publiccode-crawler/v4/internal.BuildTime=$(git log -1 --format=%cd --date=format-local:%Y-%m-%dT%H:%M:%SZ)'"

FROM alpine:3@sha256:28bd5fe8b56d1bd048e5babf5b10710ebe0bae67db86916198a6eec434943f8b

RUN apk add --no-cache gcompat \
    && mkdir -p /var/crawler/data

COPY --from=build /src/publiccode-crawler /usr/local/bin/publiccode-crawler
CMD ["publiccode-crawler", "crawl"]
