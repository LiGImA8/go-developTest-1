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


## 日志说明（Zap）

项目已接入 `zap` 结构化日志：
- 三个服务启动、连接失败、退出都会输出结构化日志
- user/order/log 的核心请求链路也会输出关键业务日志

可通过环境变量控制日志级别：

```bash
LOG_LEVEL=debug
```

常用级别：`debug` / `info` / `warn` / `error`。

## 后续扩展建议

1. 使用 `proto/minigate.proto` 生成 Go/Python stub，替换手写 gRPC 描述。
2. 在 order-service 增加重试与消息落盘队列，提升消息堆积时可靠性。
3. 在 tests 增加真实接口自动化与 DB 一致性断言。
4. 增加 Prometheus + Grafana 监控，支持压测结果对比。

## 常见问题

### `pointer to interface` 编译报错

如果你在迁移到 `minigatev1` 生成代码时遇到类似报错：

```text
s.userClient.ValidateToken undefined (type *minigatev1.UserServiceClient is pointer to interface, not interface)
```

原因是：`minigatev1.UserServiceClient` 本身是接口，不能写成指向接口的 `*minigatev1.UserServiceClient`。

错误写法：

```go
userClient *minigatev1.UserServiceClient
```

正确写法：

```go
userClient minigatev1.UserServiceClient
```

以及构造时直接传：

```go
internal.NewService(mysql, minigatev1.NewUserServiceClient(userConn), publisher)
```

### `mustEmbedUnimplementedOrderServiceServer` 编译报错

如果你遇到下面错误：

```text
cannot use internal.NewService(...) as minigatev1.OrderServiceServer
(missing method mustEmbedUnimplementedOrderServiceServer)
```

说明当前工程里混用了两套 gRPC 定义：

- `pkg/rpc`（当前仓库的手写 service desc + JSON codec）
- `minigatev1`（proto 生成代码）

二者接口类型不同，不能直接互传。

#### 方案 A（推荐，按当前仓库实现）：统一使用 `pkg/rpc`

`services/order/cmd/main.go` 应保持：

```go
server := grpc.NewServer()
rpc.RegisterOrderServiceServer(server, internal.NewService(mysql, rpc.NewUserServiceClient(userConn), publisher))
```

`services/order/internal/service.go` 的方法签名应是：

```go
func (s *Service) PlaceOrder(ctx context.Context, req *rpc.PlaceOrderRequest) (*rpc.PlaceOrderResponse, error)
```

#### 方案 B（如果你要切到 proto 生成代码）：统一使用 `minigatev1`

则需要：

1. `internal.Service` 嵌入 `minigatev1.UnimplementedOrderServiceServer`
2. `PlaceOrder` 签名改为 `*minigatev1.PlaceOrderRequest/*minigatev1.PlaceOrderResponse`
3. 注册改为 `minigatev1.RegisterOrderServiceServer(...)`

关键点是：**全链路只能选一种接口定义，不要混用。**
