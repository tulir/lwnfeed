FROM golang:1-alpine AS builder

RUN apk add --no-cache ca-certificates
WORKDIR /build/lwnfeed
COPY . /build/lwnfeed
ENV CGO_ENABLED=0
RUN go build -ldflags "-X 'main.BuildTime=`date -u +'%Y-%m-%dT%H:%M:%S+00:00'`'" -o /usr/bin/lwnfeed

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/bin/lwnfeed /usr/bin/lwnfeed
VOLUME /data

WORKDIR /data
ENTRYPOINT ["/usr/bin/lwnfeed", "-f", "/data/lwnfeed.cookie.gob"]
CMD ["start", "-l", ":8080", "-c", "/data/lwnfeed.cache.gob"]
