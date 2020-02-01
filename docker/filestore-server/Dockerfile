FROM alpine:3.11

COPY filestore-server /
RUN chmod +x /filestore-server
VOLUME /store
ENTRYPOINT ["/filestore-server","--store-dir", "/store"]

EXPOSE 8080