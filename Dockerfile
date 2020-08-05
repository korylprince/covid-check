FROM golang:1.14-alpine as builder

ARG VERSION

RUN apk add --no-cache git

RUN GO111MODULE=on go get github.com/korylprince/fileenv@v1.1.0

RUN git clone --branch "$VERSION" --single-branch --depth 1 \
    https://github.com/korylprince/covid-check.git  /go/src/github.com/korylprince/covid-check

RUN cd /go/src/github.com/korylprince/covid-check && \
    go install -mod=vendor github.com/korylprince/covid-check


FROM alpine:3.12

RUN apk add --no-cache ca-certificates

COPY --from=builder /go/bin/fileenv /
COPY --from=builder /go/bin/covid-check /

CMD ["/fileenv", "/covid-check"]
