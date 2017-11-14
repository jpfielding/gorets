FROM alpine

RUN apk --update add ca-certificates
ENV CONFIG_URL http://localhost:8888
EXPOSE 8080
EXPOSE 8000

RUN mkdir -p /opt/explorer
COPY . /opt/explorer
RUN chmod +x /opt/explorer/explorer

WORKDIR /opt/explorer

ENTRYPOINT ["/bin/sh", "-c", "./explorer -port 8080 -react /opt/explorer -configService $CONFIG_URL"]
CMD []
