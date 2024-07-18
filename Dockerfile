FROM alpine:3.20 AS builder

ARG VERSION=dev

RUN apk add --no-cache --update go gcc g++

WORKDIR /app
ENV GOPATH /app
COPY . /app

RUN CGO_ENABLED=1 GOOS=linux go build -o manager -ldflags="-X 'main.version=$VERSION'" main.go

FROM alpine:3.20

RUN mkdir -p /opt/device-manager /opt/device-manager/data
WORKDIR /opt/device-manager
COPY --from=builder /app/manager manager

ENTRYPOINT ["./manager"]