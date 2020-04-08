# Mirai HTTP Center

[![GitHub release (latest by date)](https://img.shields.io/github/v/release/tarocch1/mirai-http-center)](https://github.com/Tarocch1/mirai-http-center/releases)
[![Docker Image Version (latest semver)](https://img.shields.io/docker/v/tarocch1/mirai-http-center?label=docker%20hub)](https://hub.docker.com/r/tarocch1/mirai-http-center)
[![Docker Image Size (latest semver)](https://img.shields.io/docker/image-size/tarocch1/mirai-http-center)](https://hub.docker.com/r/tarocch1/mirai-http-center)
[![GitHub All Releases](https://img.shields.io/github/downloads/tarocch1/mirai-http-center/total)](https://github.com/Tarocch1/mirai-http-center/releases)
[![Docker Pulls](https://img.shields.io/docker/pulls/tarocch1/mirai-http-center)](https://hub.docker.com/r/tarocch1/mirai-http-center)
[![GitHub](https://img.shields.io/github/license/tarocch1/mirai-http-center)](https://github.com/Tarocch1/mirai-http-center/blob/master/LICENSE)
![Release Workflow](https://github.com/Tarocch1/mirai-http-center/workflows/Release%20Workflow/badge.svg)
![Docker Workflow](https://github.com/Tarocch1/mirai-http-center/workflows/Docker%20Workflow/badge.svg)

mirai-api-http 插件的前置代理。

## Feature

这个项目作为 [mirai-api-http](https://github.com/mamoe/mirai-api-http) 插件的增强应用，主要有以下功能：

- 通过 websocket 监听 mirai-api-http 插件的消息，并通过 http post 方式上报给指定地址。
- 可以使用 schema 对收到的消息进行筛选，并给不同的 schema 指定不同的上报地址。
- 转发 http 请求到 mirai-api-http 插件，转发时自动携带 sessionKey，无需通过 mirai-api-http 插件繁琐的认证流程。
- 转发请求使用 Bearer 认证。
- 单个二进制文件，方便部署，同时支持 docker 部署。

## Config

``` jsonc
{
  "log": {
    "level": "info" // 日志等级，包含 trace、debug、info、warn、error、fatal 和 panic，默认为 info
  },
  "mirai": {
    "authKey": "1234567890", // mirai-api-http 的 authKey
    "apiBaseURL": "http://127.0.0.1:8080", // mirai-api-http 的 api 地址
    "wsBaseURL": "ws://127.0.0.1:8080", // mirai-api-http 的 websocket 地址
    "qq": 1234567890 // bot qq 号
  },
  "schemas": {
    "/all": [ // 用于筛选 /all 频道（接收事件和消息）消息的 shcema 列表，至少应存在一个，否则不会有消息上报
      {
        "name": "default", // schema 名称
        "schema": {}, // schema 内容
        "postURL": "http://127.0.0.1" // 匹配到 chema 时的数据上报地址
      }
    ],
    "/command": [ // 用于筛选 /command 频道（接收指令）消息的 shcema 列表，至少应存在一个，否则不会有消息上报
      {
        "name": "default", // schema 名称
        "schema": {}, // schema 内容
        "postURL": "http://127.0.0.1" // 匹配到 chema 时的数据上报地址
      }
    ]
  },
  "http": {
    "host": "0.0.0.0:80", // 监听 http 请求的地址
    "authKey": "1234567890" // 用于 http 请求认证的 key
  }
}
```

## Usage

### Docker

``` bash
docker run -d -v /path/to/config.json:/usr/local/bin/mirai-http-center/config.json tarocch1/mirai-http-center
```

### Binary

```
mirai-http-center
```
