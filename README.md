# Sigríðr

## Compiling protocol buffers with gRPC

### Prerequisites for golang

```
go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
go get -u google.golang.org/grpc
```

On Fedora the `google.protobuf.*` definitions can be installed with: 
```
sudo dnf install protobuf-devel
```

### Example invocation of the protobuf compiler

```
protoc -I. -I/usr/include schema.proto --go_out=plugins=grpc:.
```