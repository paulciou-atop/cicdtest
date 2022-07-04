FROM golang:1.18.2-bullseye
WORKDIR /build
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN go build -o ./bin/snmpscan ./cmd/snmpscan
# TODO - can't get version in contaier
FROM debian:latest
WORKDIR /app/
COPY --from=0 /build/bin/snmpscan ./
ENTRYPOINT  ["/app/snmpscan"]
