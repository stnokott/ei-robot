FROM alpine:3.18.0
COPY app /
ENTRYPOINT ["/app"]
VOLUME ["/var/lib/data"]
