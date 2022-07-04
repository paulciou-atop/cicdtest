app=scanservice

PROJECT=$(shell basename  $(CURDIR))

.PHONY: clean bin/$(PROJECT)

bin/$(PROJECT): 
	go build -o ./bin/$(PROJECT) ./cmd/${PROJECT}/main.go

build:
	go build -o ${app} cmd/${app}/main.go

run: build
	${app} run

clean:
	@go clean -x
	rm -rf ./bin/* $(PROJECT)

install: bin/$(PROJECT)
	mv ./bin/* ../bin

image:
	@docker build -t ${app}:v1 .