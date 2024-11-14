This directory contains the protocol buffers used for gRPC protocols

## Developer set up

1. Install Protocol Buffers (protoc)

- For macOS (using Homebrew):

```bash
brew install protobuf
```

- For Ubuntu/Debian:

```bash
sudo apt-get install -y protobuf-compiler
```

2. Install gRPC and Protocol Buffers Go Plugin

Next, install the Go-specific protoc-gen-go and protoc-gen-go-grpc plugins:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Make sure to add the GOPATH/bin directory to your PATH so that protoc can find the plugins:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

3. Generate Go Code from `.proto` Files

Now use the protoc compiler to generate the Go code for the service and messages.

```bash
protoc --go_out=../matching-service --go-grpc_out=../matching-service ./questionmatching.proto
protoc --go_out=../question-service --go-grpc_out=../question-service ./questionmatching.proto
```

This command will generate two files:

- `questionmatching.pb.go`: Contains the message structures for requests and responses.
- `questionmatching_grpc.pb.go`: Contains the interface and client/server code for the service.
