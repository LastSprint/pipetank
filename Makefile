init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

lint:
	go tool -modfile=golangci-lint.mod golangci-lint run

fmt:
	go tool -modfile=golangci-lint.mod golangci-lint fmt