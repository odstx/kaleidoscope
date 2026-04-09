# TODO

## Features

- [ ] 应用注册表 - API 让微应用注册自己（name, js path, version）
- [ ] 应用白名单 - 代理层增加白名单，只允许通过验证的 appname
- [ ] 应用间通信 - 定义 postMessage 协议，微应用和主应用通信
- [ ] 应用配置 - 微应用获取配置（API地址等）

## DevOps

- [ ] 健康检查 - 后端代理增加对下游服务的健康检查
- [ ] 重试/熔断 - 下游服务不可用时的降级策略

## Frontend

- [ ] 共享 React - 多个微应用共用一份 React，减少 bundle 体积

## Observability

- [ ] 代理指标 - Prometheus 增加代理相关的 metrics（请求延迟、错误率等）
