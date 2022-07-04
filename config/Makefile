
PROJECT=$(shell basename  $(CURDIR))

.PHONY: bin/$(PROJECT)
bin/$(PROJECT): 
	go build -o ./bin/$(PROJECT) ./cmd/$(PROJECT)


.PHONY: clean
clean: 
	rm ./bin/*

# just a convenience task around `go test`
.PHONY: test
test:
	go test ./...

install: bin/$(PROJECT)
	mv ./bin/* ../bin
image:
	docker build --no-cache -t $(PROJECT):v1 .
