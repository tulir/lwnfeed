FROM golang:1-alpine AS builder

RUN apk add --no-cache ca-certificates
WORKDIR /build/lwnfeed
COPY . /build/lwnfeed
RUN go build -o /usr/bin/lwnfeed

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/bin/lwnfeed /usr/bin/lwnfeed
VOLUME /data

WORKDIR /data
ENTRYPOINT ["/usr/bin/lwnfeed", "-f", "/data/lwnfeed.cookie.gob"]
CMD ["start", "-c", "/data/lwnfeed.cache.gob"]
