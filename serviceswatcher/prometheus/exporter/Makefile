install-packege:
	go get github.com/prometheus/client_golang/prometheus
	go get github.com/prometheus/client_golang/prometheus/promauto
	go get github.com/prometheus/client_golang/prometheus/promhttp

grpc:
	protoc -I ../../../api/proto --go_out . --go-grpc_out . \
		--grpc-gateway_out . \
		--openapiv2_out ./api/doc --openapiv2_opt allow_merge=true,merge_file_name=api \
		../../../api/proto/v1/serviceswatcher.proto

	protoc -I ../../../api/proto --go_out . --go-grpc_out . \
		--grpc-gateway_out . \
		--openapiv2_out ./api/doc --openapiv2_opt allow_merge=true,merge_file_name=api \
		../../../api/proto/v1/snmpscan.proto