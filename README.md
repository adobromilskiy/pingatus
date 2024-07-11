# PINGATUS

[![build](https://github.com/adobromilskiy/pingatus/actions/workflows/ci.yml/badge.svg)](https://github.com/adobromilskiy/pingatus/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/adobromilskiy/pingatus)](https://goreportcard.com/report/github.com/adobromilskiy/pingatus)
[![Coverage Status](https://coveralls.io/repos/github/adobromilskiy/pingatus/badge.svg?branch=main&kill_cache=1)](https://coveralls.io/github/adobromilskiy/pingatus?branch=main)

Pingatus is a simple service to monitor your HTTP endpoints and notify you when it's down.

## Roadmap

- Monitor TCP endpoints
- Notify via email

Setup **config.yaml** via environment variables `PINGATUS_CONFIG_PATH`.

Example of **config.yaml**:

```yaml
mongouri: mongodb://localhost:27017/pingatus?timeoutMS=5000

debug: true

webapi:
  listenaddr: :8080
  assetsdir: ./assets

httppoint:
  - name: server1
    url: https://yourdomain.com
    status: 200
    timeout: 3s
    interval: 1m

notifier:
  type: telegram
  tgtoken: <telegram-bot-token>
  tgchatid: <telegram-chat-id>
```