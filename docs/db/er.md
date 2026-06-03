# 数据库 ER 说明

## 表关系

users 1 -> N orders

products 1 -> 1 inventory

orders 1 -> N order_items

products 1 -> N order_items

orders 1 -> N payments

## 为什么要拆 orders 和 order_items

orders 存订单主信息：用户、订单号、总金额、状态。

order_items 存商品明细：商品、单价、数量、小计。

即使当前接口只支持单商品下单，也先按真实订单系统设计，方便后续扩展多商品订单。

## 为什么 inventory 单独建表

products 保存商品基础信息。

inventory 保存库存信息。

库存会被频繁更新，后续会涉及事务和行级锁，因此库存独立建表更清晰。

## 为什么 payments 单独建表

orders 表示业务订单。

payments 表示支付流水。

一个订单可能有多次支付尝试，支付状态和订单状态有关联，但不是同一件事。

## 金额字段设计

所有金额使用 BIGINT，单位为分，避免小数精度问题。

## 当前进度

当前 Go 服务已经接入 MySQL，并且可以通过 `/api/v1/health/db` 验证连接状态。

当前业务 handler 仍然使用内存数据，下一步会把用户注册/登录等业务逐步迁移到 MySQL。

## 本地初始化方式

当前阶段使用 Docker Compose 启动 MySQL 8.0，并通过 `docs/db/schema.sql` 初始化数据库结构。

启动 MySQL：

```bash
docker compose up -d mysql
```

执行 schema：

```powershell
Get-Content docs/db/schema.sql | docker exec -i go-order-service-mysql mysql -uroot -prootpass
```

查看表：

```powershell
docker exec -it go-order-service-mysql mysql -uroot -prootpass -e "USE go_order_service; SHOW TABLES;"
```

## 后续演进

- Go 项目接入 MySQL
- 将内存 users/products/orders 迁移到数据库
- 实现库存扣减事务
- 实现 Redis 幂等 key
