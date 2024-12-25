# PINGATUS

![Static Badge](https://img.shields.io/badge/Go-1.23.4-blue)
[![build](https://github.com/adobromilskiy/pingatus/actions/workflows/ci.yml/badge.svg)](https://github.com/adobromilskiy/pingatus/actions/workflows/ci.yml)
[![Coverage Status](https://coveralls.io/repos/github/adobromilskiy/pingatus/badge.svg?branch=main)](https://coveralls.io/github/adobromilskiy/pingatus?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/adobromilskiy/pingatus)](https://goreportcard.com/report/github.com/adobromilskiy/pingatus)

Pingatus is a simple service to monitor your HTTP/ICMP endpoints and notify you when it's down/up.

![Pingatus dashboard](.github/assets/example.png)

## Getting started

- Create a **config.yaml** file.
- Setup **config.yaml** via environment variables `PINGATUS_CONFIG_PATH`.
- Run the docker container or build a binary.

Run via docker:

```sh
docker run -d \
    --name pingatus \
    -p 8080:8080 \
    -v $(pwd)/mycfg.yaml:/config.yaml \
    -e PINGATUS_CONFIG_PATH=/config.yaml \
    ghcr.io/adobromilskiy/pingatus:latest
```

Run via binary:

```sh
go install github.com/adobromilskiy/pingatus/cmd/pingatus@latest
export PINGATUS_CONFIG_PATH=/path/to/config.yaml
pingatus
```

Example of **config.yaml**:

```yaml
dbdsn: sqlite://pingatus.db # you can use sqlite://:memory:

listenaddr: :8080

logger:
  json: true # true - json format, false - plain text (default false)
  level: debug # debug, info, warn, error (default info)

notifier:
  type: telegram
  tgtoken: <telegram-bot-token>
  tgchatid: <telegram-chat-id>

endpoints:
  - name: server1
    type: http
    address: https://yourdomain.com
    status: 200
    timeout: 3s
    interval: 1m

  - name: server2
    type: icmp
    address: 8.8.8.8
    packetcount: 5
    interval: 1m
```

### ICMP Pinger Timeout Behavior

The ICMP pinger expects a response for sent packets within 3 seconds. If no response is received during this period, the packet is considered lost.

## Dependencies for `make` (**optional**):

```
// make sec
$ go install golang.org/x/vuln/cmd/govulncheck@latest
$ go install github.com/zricethezav/gitleaks/v8@latest

// make fmt
$ go install mvdan.cc/gofumpt@latest

// make vet
$ go install honnef.co/go/tools/cmd/staticcheck@latest
$ go install golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow@latest

// make lint
$ go install golang.org/x/tools/cmd/deadcode@latest
$ curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
| sh -s -- -b $(go env GOPATH)/bin
```
