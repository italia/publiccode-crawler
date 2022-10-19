FROM golang:1.18 as build

WORKDIR /src
COPY . .
RUN go build -ldflags "-s -w -X 'github.com/italia/developers-italia-backend/internal.VERSION=$(git describe --abbrev=0 --tags)' -X 'github.com/italia/developers-italia-backend/internal.BuildTime=$(date)'"

FROM alpine:3

COPY --from=build /src/developers-italia-backend /usr/local/bin/developers-italia-backend
CMD ["developers-italia-backend", "crawl"]
