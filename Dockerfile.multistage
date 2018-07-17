FROM golang:1.10

ADD . /go/src/github.com/tencentcloud/tencentcloud-cloud-controller-manager

WORKDIR /go/src/github.com/tencentcloud/tencentcloud-cloud-controller-manager

RUN go build --ldflags '-linkmode external -extldflags "-static"' -v -o /go/src/bin/tencentcloud-cloud-controller-manager


FROM alpine:3.6

RUN apk add --no-cache ca-certificates

COPY --from=0 /go/src/bin/tencentcloud-cloud-controller-manager /bin/tencentcloud-cloud-controller-manager

CMD ["/bin/tencentcloud-cloud-controller-manager"]