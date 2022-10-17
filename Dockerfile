FROM alpine:latest
COPY app /
ENTRYPOINT ["/app"]
VOLUME ["/var/lib/data"]
