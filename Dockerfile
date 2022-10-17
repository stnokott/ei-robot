FROM alpine:latest
COPY app /
ENTRYPOINT ["/app"]
VOLUME ["/data"]
