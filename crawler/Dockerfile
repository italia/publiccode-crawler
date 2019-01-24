# build stage
FROM golang:1.11.0-alpine AS build-env
ARG VERSION

RUN apk update && \
    apk upgrade && \
    apk add git gcc musl-dev && \
    rm -rf /var/cache/apk/*

ADD . /crawler

# Compile project
RUN cd /crawler && \
    go build -ldflags "-X github.com/italia/developers-italia-backend/crawler/version.VERSION=${VERSION}" -o bin/crawler && \
    chmod +x /crawler/bin/crawler

# final stage
FROM alpine:3.7

RUN apk update && \
    apk upgrade && \
    apk add ca-certificates && \
    apk add git && \
    rm -rf /var/cache/apk/*

WORKDIR /app
COPY --from=build-env /crawler/bin/crawler /app/
COPY ipa /app/ipa
COPY jekyll/tags.yml /app/jekyll/tags.yml

EXPOSE 8081

#ENTRYPOINT /crawler/bin/crawler one https://github.com/italia/developers.italia.it
#ENTRYPOINT /crawler/bin/crawler one https://bitbucket.org/marco-capobussi/publiccode-example
#ENTRYPOINT /crawler/bin/crawler one https://gitlab.com/inkscape/inkscape
#ENTRYPOINT /crawler/bin/crawler one https://gitlab.marconirovereto.it/Cristian.Tonelli/ProvaInfo
#ENTRYPOINT /crawler/bin/crawler one https://github.com/arpa-puglia/infoaria
ENTRYPOINT /crawler/bin/crawler crawl whitelistPA.yml whitelistGeneric.yml
