# build stage
FROM golang:1.10.0-alpine AS build-env
ARG NAME
ARG PROJECT
ARG VERSION

ADD . /go/src/$PROJECT

RUN cd /go/src/$PROJECT && go build -ldflags "-X main.Version=${VERSION}" -o $NAME

# final stage
FROM alpine:3.7
ARG NAME
ARG PROJECT

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

WORKDIR /app
COPY --from=build-env /go/src/$PROJECT/$NAME /app/
COPY --from=build-env /go/src/$PROJECT/hosting.yml /app/
EXPOSE 8081

# ARG values are not allowed in ENTRYPOINT, pass NAME as ENV variable.
ENV NAME=$NAME
ENTRYPOINT ./$NAME
