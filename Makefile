# Use to access to the protoc
# export PATH="$PATH:$(go env GOPATH)/bin"

## help: Show makefile commands
.PHONY: help
help: Makefile
	@echo "---- Project: MaxRazen/crypto-order-manager ----"
	@echo " Usage: make COMMAND"
	@echo
	@echo " Available Commands:"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

## tidy: Ensures fresh go.mod and go.sum
.PHONY: help
tidy:
	go mod tidy
	go mod verify

## grpc-generate: Generates grpc code based on protofiles/*.proto. protoc util is required
.PHONY: grpc-generate
grpc-generate:
	protoc --go_out=internal --go-grpc_out=internal protofiles/ordermanager.proto

## run: Compile and runs ordermanager
.PHONY: run
run:
	go run ./cmd/ordermanager
