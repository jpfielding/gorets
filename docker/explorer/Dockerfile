FROM alpine

RUN apk --update add ca-certificates
EXPOSE 8080

RUN mkdir -p /opt/explorer
COPY . /opt/explorer
RUN chmod +x /opt/explorer/explorer

WORKDIR /opt/explorer

ENTRYPOINT ["/bin/sh", "-c","./explorer -port 8080 -config ./config.json -react /opt/explorer"]
CMD []
