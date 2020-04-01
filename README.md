# Mirai HTTP Center

mirai-api-http 插件的前置代理。

## Why

这个项目作为 [mirai-api-http](https://github.com/mamoe/mirai-api-http) 插件的增强应用，主要解决以下两个问题：

- mirai-api-http 插件不支持通过 http 上报数据，对于函数计算等 serverless 应用场景不友好。
- mirai-api-http 插件认证过程较为繁琐，不便于调用。

## TODO

- [ ] 支持接收 websocket 数据并通过 http 转发。
- [ ] 支持接收 http 请求并转发至 mirai-api-http 插件。
- [ ] 全部功能可配置。
- [ ] Docker 部署。
