all: export GOPATH=${PWD}/../../../..
all: format
	@mkdir -p bin
	@echo "--> Running go build ${VERSION}"
	@go build -ldflags "-s -w" -v -i -o bin/go-kipper github.com/untoldwind/go-kipper/cli

format: export GOPATH=${PWD}/../../../..
format:
	@echo "--> Running go fmt"
	@go fmt ./...
