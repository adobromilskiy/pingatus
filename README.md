# PINGATUS

[![build](https://github.com/adobromilskiy/pingatus/actions/workflows/ci.yml/badge.svg)](https://github.com/adobromilskiy/pingatus/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/adobromilskiy/pingatus)](https://goreportcard.com/report/github.com/adobromilskiy/pingatus)

Pingatus is a simple service to monitor your HTTP/ICMP endpoints and notify you when it's down/up.

Setup **config.yaml** via environment variables `PINGATUS_CONFIG_PATH`.

Example of **config.yaml**:

```yaml
mongouri: mongodb://localhost:27017/pingatus?timeoutMS=5000

debug: true

webapi:
  listenaddr: :8080
  assetsdir: ./assets

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

notifier:
  type: telegram
  tgtoken: <telegram-bot-token>
  tgchatid: <telegram-chat-id>
```