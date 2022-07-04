CGO_CPPFLAGS ?= ${CPPFLAGS}
export CGO_CPPFLAGS
CGO_CFLAGS ?= ${CFLAGS}
export CGO_CFLAGS
CGO_LDFLAGS ?= $(filter -g -L% -l% -O%,${LDFLAGS})
export CGO_LDFLAGS


PROJECT=$(shell basename  $(CURDIR))

.PHONY: bin/$(PROJECT)
bin/$(PROJECT): 
	go build -o ./bin/messaging ./


.PHONY: clean
clean: 
	rm ./bin/*

# just a convenience task around `go test`
.PHONY: test
test:
	go test ../test

install: bin/$(PROJECT)
	mv ./bin/* ../bin

setuposx:
	brew install protobuf
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
	export PATH="$PATH:$(go env GOPATH)/bin"
	go get github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor@v2.9.0
	go install  github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gatewa
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2


setupwsl:
	sudo apt install protobuf
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
	export PATH="$PATH:$(go env GOPATH)/bin"