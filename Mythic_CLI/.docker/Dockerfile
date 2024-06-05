#FROM itsafeaturetiger/tiger_go_base:latest
FROM golang:1.21-alpine
WORKDIR /usr/src/app

ARG GOPROXY=proxy.golang.org
ARG GO111MODULE

RUN go env -w GOPROXY=${GOPROXY}
RUN go env -w GO111MODULE=${GO111MODULE}

COPY ["src/", "."]

RUN apk add --no-cache make

RUN make build_all

FROM alpine

COPY --from=0 /usr/src/app/tiger-cli_linux /tiger-cli_linux
COPY --from=0 /usr/src/app/tiger-cli_macos /tiger-cli_macos