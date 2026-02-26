# JMeter 压测设计（初版）

- 线程组：20/50/100 并发，Ramp-up 30s，持续 5 分钟
- 接口：登录 + 下单（串联）
- 观测指标：QPS、P95/P99、错误率
- MySQL 观测：`SHOW ENGINE INNODB STATUS`、慢查询日志、锁等待
- 异常注入：无效 token、参数异常、停掉 log-service 观察消息堆积
