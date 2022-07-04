FROM golang:1.18.2-bullseye
WORKDIR /build
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .

# Build services
RUN go build -o ./bin/config ./config/cmd/config
RUN go build -o ./bin/serviceswatcher ./serviceswatcher/cmd/serviceswatcher
RUN go build -o ./bin/snmpscan ./snmpscan/cmd/snmpscan
RUN go build -o ./bin/atopudpscan ./atopudpscan/cmd/atopudpscan
RUN go build -o ./bin/inventory ./inventory/cmd/inventory
RUN go build -o ./bin/scanservice ./scanservice/cmd/scanservice
# TODO - can't get version in contaier

# config
FROM debian:latest AS config
WORKDIR /app/
COPY --from=0 /build/bin/config ./
ENTRYPOINT  ["/app/config"]

# serviceswatcher
FROM debian:latest AS serviceswatcher
WORKDIR /app/
COPY --from=0 /build/bin/serviceswatcher ./
ENTRYPOINT  ["/app/serviceswatcher"]

# snmpscan
FROM debian:latest AS snmpscan
WORKDIR /app/
COPY --from=0 /build/bin/snmpscan ./
ENTRYPOINT  ["/app/snmpscan"]

# atopudpscan
FROM debian:latest AS atopudpscan
WORKDIR /app/
COPY --from=0 /build/bin/atopudpscan ./
ENTRYPOINT  ["/app/atopudpscan"]

# inventory
FROM debian:latest AS inventory
WORKDIR /app/
COPY --from=0 /build/bin/inventory ./
ENTRYPOINT  ["/app/inventory"]

# scan service
FROM debian:latest AS scanservice
WORKDIR /app/
COPY --from=0 /build/bin/scanservice ./
ENTRYPOINT  ["/app/scanservice"]