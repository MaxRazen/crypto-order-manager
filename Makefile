# Use to access to the protoc
# export PATH="$PATH:$(go env GOPATH)/bin"

## help: Show makefile commands
help: Makefile
	@echo "---- Project: MaxRazen/crypto-order-manager ----"
	@echo " Usage: make COMMAND"
	@echo
	@echo " Available Commands:"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

## tidy: Ensures fresh go.mod and go.sum
tidy:
	go mod tidy
	go mod verify

## grpc-generate: Generates grpc code based on protofiles/*.proto. protoc util is required
grpc-generate:
	protoc --go_out=internal --go-grpc_out=internal protofiles/ordermanager.proto

## run: Compile and runs ordermanager
run:
	go run ./cmd/ordermanager

## test: Runs tests across the project with no cache
test:
	go test -count=1 ./...
