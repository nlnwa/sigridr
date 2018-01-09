# Sigríðr

A Twitter API client.

### Usage

First time with consumer key and consumer secret (after which the access token is stored in a config file `~/.sigridr.yaml`):
```
./sigridr -k <consumer key> -s <consumer secret> search from:nasjonalbibl
```

Access token provided as environment variable:
```
ACCESS_TOKEN=<access token> ./sigridr search from:nasjonalbibl
```

With filters (no replies and no retweets):
```
./sigridr search from:nasjonalbibl -- -filter:replies -filter:retweets
```

## Development

### Build
```
go build -o sigridr
```

### Compiling protocol buffers (with gRPC)

#### Prerequisites (golang)

```
go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
go get -u google.golang.org/grpc
```

On Fedora the `google.protobuf.*` definitions can be installed with: 
```
sudo dnf install protobuf-devel
```

#### Compile stubs
Example invocation of the protobuf compiler to generate go stubs:

```
protoc -I. -I/usr/include *.proto --go_out=plugins=grpc:.
```
