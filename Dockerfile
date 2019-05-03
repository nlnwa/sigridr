FROM golang:alpine

ENV GO111MODULE=on

RUN apk add --no-cache --update alpine-sdk protobuf protobuf-dev
WORKDIR /go/src/github.com/nlnwa/sigridr

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN VERSION=$(./scripts/git-version) \
&& go generate ./api \
&& go install -v -ldflags "-w -X github.com/nlnwa/sigridr/version.Version=${VERSION}" ./cmd/...
# -w Omit the DWARF symbol table.
# -X Set the value of the string variable in importpath named name to value.


FROM alpine:latest

LABEL maintainer="marius.beck@nb.no"

RUN apk add --no-cache --update ca-certificates

COPY --from=0 /go/bin/sigridr /usr/local/bin/

WORKDIR /

ENTRYPOINT ["sigridr"]
CMD ["--help"]

EXPOSE 10000 10001
