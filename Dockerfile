FROM golang:1.22 AS build-stage

WORKDIR /build

ENV GOPROXY https://goproxy.cn,direct

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /build/easeway ./cmd

FROM debian:12 AS build-release-stage

RUN apt-get update && apt-get install -y ca-certificates
WORKDIR /app
COPY --from=build-stage /build/easeway easeway
EXPOSE 8055
ENTRYPOINT ["/app/easeway", "-f", "/conf/config.yaml"]
