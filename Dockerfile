FROM golang:1.17-alpine  AS builder

RUN apk add --update make

WORKDIR /go/src/github.com/ozonva/ova-journey-api/

COPY . /go/src/github.com/ozonva/ova-journey-api/

RUN make deps && make build

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

RUN mkdir swagger
COPY ./swagger ./swagger
COPY ./config/config.yaml ./config/config.yaml
COPY --from=builder /go/src/github.com/ozonva/ova-journey-api/bin/ova-journey-api .

RUN chown root:root ova-journey-api

EXPOSE 8080
EXPOSE 9090
CMD ["./ova-journey-api"]