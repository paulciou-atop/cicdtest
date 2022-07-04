PROJECT=$(shell basename  $(CURDIR))

all: 
	protoc -I ./proto --go_out . --go_opt=paths=source_relative \
   --go-grpc_out . --go-grpc_opt=paths=source_relative \
   --grpc-gateway_out . --grpc-gateway_opt=paths=source_relative \
   --openapiv2_out ./doc --openapiv2_opt allow_merge=true,merge_file_name=api \
   ./proto/v1/**/*.proto


.PHONY: clean
clean: 
	rm -rf ./v1/*
