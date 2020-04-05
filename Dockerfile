FROM golang:latest as builder

WORKDIR /go/src
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build

FROM alpine:latest

WORKDIR /usr/local/bin/mirai-http-center/
COPY --from=builder /go/src/mirai-http-center .
COPY config.json .

RUN set -ex && \
    apk --no-cache add ca-certificates && \
    chmod +x /usr/local/bin/mirai-http-center/mirai-http-center

ENV PATH /usr/local/bin/mirai-http-center:$PATH

CMD ["mirai-http-center"]
