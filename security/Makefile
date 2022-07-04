grpc : 
	protoc -I ../api/proto --proto_path=./vender --go_out . --go-grpc_out . \
		--grpc-gateway_out . \
		--openapiv2_out ./api/doc --openapiv2_opt allow_merge=true,merge_file_name=api \
		../api/proto/v1/security.proto
	


	