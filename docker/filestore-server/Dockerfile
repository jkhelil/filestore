FROM alpine:3.11

COPY filestore-server /
RUN chmod +x /filestore-server
VOLUME /store
EXPOSE 8080
ENTRYPOINT ["/filestore-server"]
CMD ["--store-dir", "/store"]
