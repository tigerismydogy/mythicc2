FROM golang:1.21 as builder

WORKDIR /usr/src/app

ARG GOPROXY=proxy.golang.org
ARG GO111MODULE

RUN go env -w GOPROXY=${GOPROXY}
RUN go env -w GO111MODULE=${GO111MODULE}

COPY ["src/", "."]

RUN make build_final

FROM alpine

COPY --from=builder /tiger_server /tiger_server
COPY --from=builder /usr/src/app /usr/src/app

WORKDIR /usr/src/app

HEALTHCHECK --interval=60s --timeout=10s --retries=5 --start-period=20s \
  CMD wget -SqO - http://127.0.0.1:${tiger_SERVER_PORT:-17443}/health || exit 1

CMD ["/bin/sh", "-c", "cp /tiger_server . && ./tiger_server" ]