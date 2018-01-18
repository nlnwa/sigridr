FROM golang:1.9.2-alpine

RUN apk add --no-cache --update alpine-sdk

COPY . /go/src/github.com/nlnwa/sigridr
RUN cd /go/src/github.com/nlnwa/sigridr && make release-binary

FROM alpine:3.4

COPY --from=0 /go/bin/sigridr /go/bin/sigridrctl /usr/local/bin/

WORKDIR /

ENTRYPOINT ["sigridr"]

CMD ["version"]

EXPOSE 10000 10001
