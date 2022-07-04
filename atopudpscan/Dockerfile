FROM alpine:latest
RUN apk add --no-cache go


RUN GOARCH=amd64 \
    GOOS=linux \
    CGO_ENABLED=0 	
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY ./ ./
RUN go build -o ./cmd/atopudpscan/atopudpscan ./cmd/atopudpscan/main.go


FROM alpine:latest
WORKDIR /atopUdp
COPY --from=0 /app ./
EXPOSE 8083
EXPOSE 8093
EXPOSE 55954/udp
ENTRYPOINT ["./cmd/atopudpscan/atopudpscan"]

