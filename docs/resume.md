# 简历项目描述

## 精简版

**go-order-service | Go 订单后端服务**

基于 Go + Gin + MySQL + Redis 实现的模拟订单后端服务，覆盖用户注册登录、JWT 鉴权、商品查询、订单创建、库存扣减、Redis 幂等和模拟支付状态流转。
使用 MySQL 事务和 `SELECT ... FOR UPDATE` 保证下单时库存扣减、订单写入和订单明细写入的一致性。
使用 Redis `Idempotency-Key` 防止重复下单，并通过 Docker Compose、单元测试、MySQL 集成测试和 GitHub Actions CI 提升项目可交付性。

## 详细版

项目名称：go-order-service  
项目类型：Go 后端作品集项目  
技术栈：Go、Gin、MySQL、Redis、JWT、Docker Compose、GitHub Actions、database/sql、go-redis/v9

项目描述：

该项目是一个模拟电商订单核心链路的 Go 后端服务，覆盖用户注册登录、JWT 鉴权、商品查询、订单创建、库存扣减、订单幂等、支付状态流转、统一响应、请求日志、Docker Compose 和自动化测试。

主要工作：

- 设计 `users`、`products`、`inventory`、`orders`、`order_items`、`payments` 等核心表结构。
- 基于 Gin 实现 RESTful API，包括注册登录、商品查询、创建订单、模拟支付等接口。
- 使用 MySQL 事务和 `SELECT ... FOR UPDATE` 实现下单库存扣减，保证库存、订单和订单明细的一致性。
- 使用 Redis `SET NX EX` 实现 `Idempotency-Key`，防止重复点击或网络重试导致重复下单。
- 实现模拟支付状态流转，支持 `PENDING_PAYMENT -> PAID / PAYMENT_FAILED`，支付失败时恢复库存。
- 封装统一响应结构和请求日志中间件，支持 `X-Request-ID` 追踪请求。
- 使用 Docker Compose 编排 MySQL、Redis 和 Go 服务，支持一键启动。
- 编写单元测试和 MySQL Repository 集成测试，并接入 GitHub Actions CI 自动验证核心链路。

## 可量化亮点

- 设计并实现 6 张核心业务表，覆盖用户、商品、库存、订单、订单明细和支付流水。
- 实现订单创建事务，保证库存扣减、订单写入、订单明细写入在同一事务中完成。
- 使用 Redis 幂等 key 拦截重复下单请求，避免库存被重复扣减。
- 补充单元测试和 MySQL 集成测试，覆盖 JWT、鉴权中间件、请求日志、Redis 幂等、订单事务和支付事务。
- 接入 GitHub Actions CI，在 push 和 pull request 时自动执行单元测试和 MySQL 集成测试。

## 可选简历标题

- Go 订单后端服务
- Go + Gin 电商订单系统
- 基于 Go 的订单链路后端项目

## 注意事项

简历中不要写：

- 高并发订单系统
- 分布式微服务
- 真实支付系统
- 秒杀系统
- 生产级电商平台

除非后续确实完成压测、消息队列、监控、限流等能力，否则保持描述诚实、克制，更容易让面试官信任项目完成度。
