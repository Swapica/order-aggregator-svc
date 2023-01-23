FROM golang:1.18-alpine as buildbase

RUN apk add git build-base

WORKDIR /go/src/github.com/Swapica/order-aggregator-svc
COPY vendor .
COPY . .

RUN GOOS=linux go build  -o /usr/local/bin/order-aggregator-svc /go/src/github.com/Swapica/order-aggregator-svc


FROM alpine:3.9

COPY --from=buildbase /usr/local/bin/order-aggregator-svc /usr/local/bin/order-aggregator-svc
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["order-aggregator-svc"]
