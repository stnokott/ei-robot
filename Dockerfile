FROM alpine:3.17.3
COPY app /
ENTRYPOINT ["/app"]
VOLUME ["/var/lib/data"]
