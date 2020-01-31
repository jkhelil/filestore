FROM alpine:3.11

COPY webcrawler-client /
RUN chmod +x /webcrawler-client

ENTRYPOINT ["/webcrawler-client", "-page"] 
