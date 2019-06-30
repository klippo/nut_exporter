FROM golang:1.12-alpine as builder
LABEL maintainer=klippo@deny.se

WORKDIR /go/src/github.com/klippo/nut_exporter

COPY . .

RUN apk add --no-cache git
RUN go get -d -v ./...
RUN go build -o /go/bin/nut_exporter

# -----

FROM alpine:3.10
COPY --from=builder /go/bin/nut_exporter /bin

RUN echo "http://dl-cdn.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories
RUN apk add --no-cache nut
RUN ln -sf /usr/bin/upsc /bin/upsc

ENTRYPOINT ["/bin/nut_exporter"] 
