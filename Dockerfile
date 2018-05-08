# build stage
FROM golang:1.10.0-alpine AS build-env
ARG NAME
ARG PROJECT
ARG VERSION

RUN apk update && \
    apk upgrade && \
    apk add --no-cache git && \
    apk add --no-cache gcc && \
    apk add --no-cache musl-dev

ADD . /go/src/$PROJECT

# Dep ensure. Uncomment if you don't have a ./vendor folder for go deps.
# RUN cd /go/src/$PROJECT && go get -u github.com/golang/dep/cmd/dep && dep ensure

# Compile .so plugins
RUN cd /go/src/$PROJECT/plugins && go build -buildmode=plugin -o out/github.so github/plugin.go
RUN cd /go/src/$PROJECT/plugins && go build -buildmode=plugin -o out/gitlab.so gitlab/plugin.go
RUN cd /go/src/$PROJECT/plugins && go build -buildmode=plugin -o out/bitbucket.so bitbucket/plugin.go

# Compile project
RUN cd /go/src/$PROJECT && go build -ldflags "-X github.com/italia/developers-italia-backend/version.VERSION=${VERSION}" -o $NAME

# final stage
FROM alpine:3.7
ARG NAME
ARG PROJECT

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

WORKDIR /app
COPY --from=build-env /go/src/$PROJECT/$NAME /app/
COPY --from=build-env /go/src/$PROJECT/domains.yml /app/
COPY --from=build-env /go/src/$PROJECT/plugins/out/ /app/plugins/out/
EXPOSE 8081

# ARG values are not allowed in ENTRYPOINT, pass NAME as ENV variable.
ENV NAME=$NAME
RUN chmod +x ./$NAME

ENTRYPOINT ./$NAME all
