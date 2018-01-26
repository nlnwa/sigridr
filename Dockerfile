FROM golang:alpine

RUN apk add --no-cache --update alpine-sdk protobuf protobuf-dev

COPY . /go/src/github.com/nlnwa/sigridr
RUN cd /go/src/github.com/nlnwa/sigridr && make release-binary


FROM alpine:3.7
LABEL maintainer="nettarkivet@nb.no"

RUN apk add --no-cache --update ca-certificates

COPY --from=0 /go/bin/sigridr /usr/local/bin/

WORKDIR /

ENTRYPOINT ["sigridr"]
CMD ["--help"]

EXPOSE 10000 10001
