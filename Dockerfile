FROM golang:latest as builder

WORKDIR /go/src
COPY . .

ARG VERSION
ARG COMMIT_ID

RUN set -ex && \
    go mod download && \
    GO_VERSION=$(go version) && \
    BUILD_TIME=$(date) && \
    CGO_ENABLED=0 go build -ldflags "-X 'main.version=${VERSION}' -X 'main.goVersion=${GO_VERSION}' -X 'main.buildTime=${BUILD_TIME}' -X 'main.commitID=${COMMIT_ID}'"

FROM alpine:latest

WORKDIR /usr/local/bin/mirai-http-center/
COPY --from=builder /go/src/mirai-http-center .
COPY config.json .

RUN set -ex && \
    apk --no-cache add ca-certificates curl && \
    chmod +x /usr/local/bin/mirai-http-center/mirai-http-center

HEALTHCHECK --interval=30s --timeout=3s --retries=3 \
  CMD curl -fs -H 'authorization: Bearer _healthcheck' http://localhost/_healthcheck || exit 1

ENV PATH /usr/local/bin/mirai-http-center:$PATH

CMD ["mirai-http-center"]
