FROM golang:alpine

RUN apk add --no-cache --update alpine-sdk protobuf protobuf-dev

COPY . /go/src/github.com/nlnwa/sigridr

RUN cd /go/src/github.com/nlnwa/sigridr \
&& go generate github.com/nlnwa/sigridr/api \
&& go get github.com/golang/dep/cmd/dep \
&& dep ensure -vendor-only \
&& VERSION=$(./scripts/git-version) \
go install -v -ldflags "-w -X github.com/nlnwa/sigridr/version.Version=$(VERSION)" github.com/nlnwa/sigridr/cmd/...
# -w Omit the DWARF symbol table.
# -X Set the value of the string variable in importpath named name to value.


FROM alpine:3.7
LABEL maintainer="nettarkivet@nb.no"

RUN apk add --no-cache --update ca-certificates

COPY --from=0 /go/bin/sigridr /usr/local/bin/

WORKDIR /

ENTRYPOINT ["sigridr"]
CMD ["--help"]

EXPOSE 10000 10001
