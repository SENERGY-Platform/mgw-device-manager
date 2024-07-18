FROM golang:1.22 AS builder

ARG VERSION=dev

COPY . /go/src/app
WORKDIR /go/src/app

RUN CGO_ENABLED=1 GOOS=linux go build -o manager -ldflags="-X 'main.version=$VERSION'" main.go

FROM alpine:3.19

RUN mkdir -p /opt/device-manager /opt/device-manager/data
WORKDIR /opt/device-manager
COPY --from=builder /go/src/app/manager manager

ENTRYPOINT ["./manager"]
