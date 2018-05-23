FROM alpine:3.6

RUN apk add --no-cache ca-certificates

ADD tencentcloud-cloud-controller-manager /bin/

CMD ["/bin/tencentcloud-cloud-controller-manager"]