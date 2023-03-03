FROM alpine:3.17.2
COPY app /
ENTRYPOINT ["/app"]
VOLUME ["/var/lib/data"]
