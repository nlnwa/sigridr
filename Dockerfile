FROM golang:1.9.2-alpine

RUN apk add --no-cache --update alpine-sdk

COPY . /go/src/github.com/nlnwa/sigridr
RUN cd /go/src/github.com/nlnwa/sigridr && make release-binary

FROM alpine:3.4

RUN apk add --no-cache --update ca-certificates

COPY --from=0 /go/bin/sigridr /usr/local/bin/

WORKDIR /

ENTRYPOINT ["sigridr"]

EXPOSE 10000 10001
