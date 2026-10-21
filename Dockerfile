FROM golang:1.26@sha256:3745f1aa11822297b79113f5e84ce2bffeb039734571c39855809cfe9e5d7be6 AS build

WORKDIR /src
COPY . .
RUN go build -trimpath -ldflags "-s -w -buildid= -X 'github.com/italia/publiccode-crawler/v4/internal.VERSION=$(git describe --abbrev=0 --tags)' -X 'github.com/italia/publiccode-crawler/v4/internal.BuildTime=$(git log -1 --format=%cd --date=format-local:%Y-%m-%dT%H:%M:%SZ)'"

FROM alpine:3@sha256:294b683cb724975bec92580e1e685676bd4b50bda910ddb8c51d4cabeaec77e6

RUN apk add --no-cache gcompat \
    && mkdir -p /var/crawler/data

COPY --from=build /src/publiccode-crawler /usr/local/bin/publiccode-crawler
CMD ["publiccode-crawler", "crawl"]
