FROM golang:1.18.2-bullseye
WORKDIR /build
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN go build -o ./bin/config ./cmd/config
# TODO - can't get version in contaier
FROM debian:latest
WORKDIR /app/
COPY --from=0 /build/bin/config ./
ENTRYPOINT  ["/app/config"]