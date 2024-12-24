FROM golang:1.23-alpine AS builder
WORKDIR /build
COPY . ./
ENV GOBIN="/build/bin"
RUN go install -mod=vendor ./cmd/...

FROM alpine:latest
RUN apk update && apk add ca-certificates tzdata && rm -rf /var/cache/apk/*
ENV TZ="Europe/Kyiv"
COPY --from=builder /build/bin/pingatus .
ENTRYPOINT [ "/pingatus" ]
