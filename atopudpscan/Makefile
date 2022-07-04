

CGO_CPPFLAGS ?= ${CPPFLAGS}
export CGO_CPPFLAGS
CGO_CFLAGS ?= ${CFLAGS}
export CGO_CFLAGS
CGO_LDFLAGS ?= $(filter -g -L% -l% -O%,${LDFLAGS})
export CGO_LDFLAGS

app = atopudpscan
grpc_prot=8083
http_prot=8093
udp_port=55954

PROJECT=$(shell basename  $(CURDIR))

.PHONY: clean  test  bin/$(PROJECT)

bin/$(PROJECT): 
	go build -o ./bin/$(PROJECT) ./cmd/${PROJECT}/main.go

build :
	go build -o ${app} cmd/${app}/main.go

run: build
	${app} run

clean:
	@go clean -x
	rm -rf ./bin/* $(PROJECT)

test: 	
	@go build -o ${app}.test cmd/${app}/main.go
	./${app}.test run

image:
	@docker build -t ${app}:v1 .

docker-run: image
	@docker run -itd --name ${app} --net=host -p ${grpc_prot}:${grpc_prot} -p ${http_prot}:${http_prot} -p ${udp_port}:${udp_port}/udp ${app} run

install: bin/$(PROJECT)
	mv ./bin/* ../bin


grpc: 
	protoc -I ../api/proto --proto_path=./vendor --proto_path=.  --go_out . --go-grpc_out . \
		--grpc-gateway_out=allow_delete_body=true:. \
		--openapiv2_out ./api/doc --openapiv2_opt allow_merge=true,merge_file_name=api,allow_delete_body=true \
		 ../api/proto/v1/atopudpscan/atopudpscan.proto 

#	protoc -I ../api/proto --proto_path=./vender --proto_path=.   --go_out=plugins=grpc:.  ../api/proto/v1/db_store.proto
#	protoc -I ../api/proto --proto_path=./vender --go_out . --go-grpc_out . \
		--grpc-gateway_out . \
		--openapiv2_out ./api/doc --openapiv2_opt allow_merge=true,merge_file_name=api \
		../api/proto/v1/serviceswatcher.proto



help:
	@echo Usage:
	@echo "make  <target>"
	@echo Targets:	
	@echo bin/$(PROJECT): build bin/$(PROJECT)
	@echo build:		build ${app}
	@echo run:  		run ${app} service 
	@echo image:	docker build ${app}
	@echo docker-run:	docker run ${app}
	@echo test:	 	run test ${app} service
	@echo clean:	 	clean test data
	@echo grpc:		build *.proto file