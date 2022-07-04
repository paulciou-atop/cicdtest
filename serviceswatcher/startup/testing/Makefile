grpc:
	protoc -I ../../../api/proto --go_out . --go-grpc_out . \
		--grpc-gateway_out . \
		--openapiv2_out ./api/doc --openapiv2_opt allow_merge=true,merge_file_name=api \
		../../../api/proto/v1/serviceswatcher.proto

	protoc -I ../../../api/proto --go_out . --go-grpc_out . \
		--grpc-gateway_out . \
		--openapiv2_out ./api/doc --openapiv2_opt allow_merge=true,merge_file_name=api \
		../../../api/proto/v1/snmpscan.proto

setupgrpc:
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
	
	