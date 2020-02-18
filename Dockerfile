FROM golang:1.12 as golangBuilder

COPY . /go/src/github.com/jpfielding/gorets

WORKDIR /go/src/github.com/jpfielding/gorets

RUN GO111MODULE="on" make build-explorer-svc

# ==============

FROM node:8.9 as nodeBuilder

ARG HEADERS_CONFIG={}
ARG CONFIG_ENV=docker
ARG CONFIG_DOCKER="module.exports = { \
  staticAssetPath: '.', \
  api: '', \
};"

RUN echo ${CONFIG_ENV} "\n" ${CONFIG_DOCKER} "\n" ${HEADERS_CONFIG}

COPY . /gorets

WORKDIR /gorets/

RUN echo ${CONFIG_DOCKER} >> web/explorer/client/config/docker.js
RUN make build-explorer-client
RUN echo ${HEADERS_CONFIG} >> bin/explorer/config.json

# ==============

FROM alpine

RUN apk --update add ca-certificates
EXPOSE 8080

RUN mkdir -p /opt/explorer
COPY --from=golangBuilder /go/src/github.com/jpfielding/gorets/bin/explorer/ /opt/explorer/
COPY --from=nodeBuilder /gorets/bin/explorer/ /opt/explorer/
RUN chmod +x /opt/explorer/explorer

WORKDIR /opt/explorer

ENTRYPOINT ["/bin/sh", "-c","./explorer -port 8080 -config ./config.json -react /opt/explorer"]
CMD []
