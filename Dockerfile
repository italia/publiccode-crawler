FROM golang:1.26@sha256:3745f1aa11822297b79113f5e84ce2bffeb039734571c39855809cfe9e5d7be6 AS build

WORKDIR /src
COPY . .
RUN go build -trimpath -ldflags "-s -w -buildid= -X 'github.com/italia/publiccode-crawler/v4/internal.VERSION=$(git describe --abbrev=0 --tags)' -X 'github.com/italia/publiccode-crawler/v4/internal.BuildTime=$(git log -1 --format=%cd --date=format-local:%Y-%m-%dT%H:%M:%SZ)'"

FROM alpine:3@sha256:f59abe656b692d901e184be3445e084786602108599d2a51d229c0f5a48c71bf

RUN apk add --no-cache gcompat \
    && mkdir -p /var/crawler/data

COPY --from=build /src/publiccode-crawler /usr/local/bin/publiccode-crawler
CMD ["publiccode-crawler", "crawl"]
