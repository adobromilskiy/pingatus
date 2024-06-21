# PINGATUS

Pingatus is a simple service to monitor your HTTP endpoints and notify you when it's down.

## Features

- Monitor TCP endpoints
- Notify via email
- Frontend dashboard more dynamic

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