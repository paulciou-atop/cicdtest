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
RUN go build -o ./security ./main.go


FROM alpine:latest
WORKDIR /security
COPY --from=0 /app ./
RUN apk add step-cli 
EXPOSE 8080
ENTRYPOINT ["./security"]

