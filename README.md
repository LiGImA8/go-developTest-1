# MiniGate

基于 Go + gRPC + ZeroMQ + MySQL 的微服务骨架项目，包含：
- `user-service`：登录、token 鉴权
- `order-service`：同步下单（gRPC 调 user-service）
- `log-service`：异步日志消费（ZeroMQ）

## 调用链路

`登录 -> 下单 -> 鉴权 -> 订单入库 -> ZeroMQ 发送日志 -> 日志入库`

## 项目结构

- `services/user`：用户服务
- `services/order`：订单服务
- `services/log`：日志服务
- `pkg/rpc`：gRPC 手写 service desc + JSON codec
- `pkg/logger`：ZeroMQ 事件发布
- `sql/init.sql`：数据库初始化脚本
- `tests/`：pytest + allure 测试骨架
- `scripts/jmeter_plan.md`：JMeter 压测方案

## 快速启动

```bash
docker compose up --build
```

服务端口：
- user-service: `50051`
- order-service: `50052`
- log-service(ZMQ): `5557`
- mysql: `3306`

## 后续扩展建议

1. 使用 `proto/minigate.proto` 生成 Go/Python stub，替换手写 gRPC 描述。
2. 在 order-service 增加重试与消息落盘队列，提升消息堆积时可靠性。
3. 在 tests 增加真实接口自动化与 DB 一致性断言。
4. 增加 Prometheus + Grafana 监控，支持压测结果对比。
