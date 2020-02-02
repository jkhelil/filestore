FROM alpine:3.11

COPY filestore-client /
RUN chmod +x /filestore-client

ENTRYPOINT ["/filestore-client"] 
CMD ["--help"]
