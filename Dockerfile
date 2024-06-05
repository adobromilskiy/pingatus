FROM golang:1.22-alpine AS builder
WORKDIR /build
COPY . ./
RUN go build -mod=vendor -o pingatus

FROM alpine:latest
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

COPY --from=builder /build/pingatus .
ENTRYPOINT [ "/pingatus" ]