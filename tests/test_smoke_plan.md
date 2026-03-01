# pytest + allure 测试计划（骨架）

1. 登录成功与失败（密码错误/用户不存在）
2. 下单成功，并校验 orders + logs 表一致性
3. token 过期场景
4. 参数异常（空 item_name、quantity<=0）
5. 服务异常与消息堆积（通过停止 log-service 观察）

> 当前仓库先搭建运行骨架，后续可基于 `proto/minigate.proto` 生成 Python gRPC stub 完成自动化脚本。
