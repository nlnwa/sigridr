package api

// Assumes protoc is installed
//go:generate go get github.com/golang/protobuf/proto
//go:generate go get github.com/golang/protobuf/protoc-gen-go
//go:generate protoc -I/usr/include -I. --go_out=plugins=grpc:. agent.proto schema.proto worker.proto
