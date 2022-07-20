FROM golang:1.18.3-buster AS builder

WORKDIR /opt

COPY suites suites
COPY internal internal
COPY go.mod go.mod
COPY go.sum go.sum
COPY Makefile Makefile

RUN make tests-binary

FROM alpine:3.16

WORKDIR /opt

RUN addgroup -g 1000 testgroup && \
    adduser -D testuser -u 1000 -G testgroup && \
    apk add \
        --no-cache \
        ginkgo=2.1.4-r1

USER testuser:testgroup

COPY --from=builder /opt/merkle.test merkle.test
COPY docker/tests/merkle.yml merkle.yml

CMD [ "ginkgo", "-v", "-p", "-race", "merkle.test" ]
