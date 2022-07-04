FROM golang:1.18-alpine	
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go get github.com/lib/pq
COPY . .
RUN go build -o /app/testing
CMD  ["go","run","/app/main.go"]
