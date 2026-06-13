# 数据库 ER 说明

## 表关系

- `users` 1 -> N `orders`
- `products` 1 -> 1 `inventory`
- `orders` 1 -> N `order_items`
- `products` 1 -> N `order_items`
- `orders` 1 -> N `payments`

## 为什么要拆 `orders` 和 `order_items`

`orders` 保存订单主信息，例如用户、订单号、总金额和状态。  
`order_items` 保存商品明细，例如商品、单价、数量和小计。  

即使当前接口仍然是单商品下单，也按真实电商订单的方式设计，方便后续扩展成多商品订单。

## 为什么要单独建 `inventory`

`products` 保存商品基础信息。  
`inventory` 保存库存信息，库存会频繁更新，而且下单时需要事务和行级锁，所以拆表更清晰，也更便于控制并发。

## 为什么要单独建 `payments`

`orders` 表示业务订单，`payments` 表示支付流水。  
一个订单可能对应一次或多次支付尝试，所以支付信息应单独记录。

当前项目里的模拟支付接口 `POST /api/v1/payments/mock` 会在数据库里完成下面的事情：

1. `SELECT orders ... FOR UPDATE` 锁定订单
2. 校验订单是否属于当前用户
3. 校验订单状态是否为 `PENDING_PAYMENT`
4. 写入 `payments`
5. `SUCCESS` 时更新 `orders.status = PAID`
6. `FAILED` 时更新 `orders.status = PAYMENT_FAILED`
7. `FAILED` 时把该订单扣减过的库存加回 `inventory.stock`
8. 提交事务

## 金额字段设计

所有金额字段都使用 `BIGINT`，单位为分，避免小数精度问题。

## 订单状态

当前 `orders.status` 支持以下状态：

- `PENDING_PAYMENT`：订单已创建，等待支付
- `PAID`：支付成功
- `PAYMENT_FAILED`：支付失败

## 当前进度

Go 服务已经接入 MySQL，`/api/v1/health/db` 可用，`products` 和 `inventory` 已接入商品列表接口 `/api/v1/products`。  
`POST /api/v1/orders` 已接入 Redis 幂等 key，并在事务中完成库存校验、库存扣减、`orders` 写入和 `order_items` 写入。  
`POST /api/v1/payments/mock` 已接入 JWT 鉴权，并在事务中完成支付记录写入、订单状态流转和失败库存回补。

Redis 也已经接入服务启动流程，并提供 `/api/v1/health/redis` 健康检查接口，用于验证 Redis 连接状态。  
`POST /api/v1/orders` 会先用 `userID + Idempotency-Key` 组成 Redis key，重复请求会返回 `409 duplicate request`，不会重复扣减库存。

## 订单创建流程

`POST /api/v1/orders` 会先处理 Redis 幂等 key，再使用数据库事务：

1. 读取 `userID + Idempotency-Key`
2. 使用 Redis `SET NX EX` 创建幂等占位
3. 查询商品和库存
4. 使用 `SELECT ... FOR UPDATE` 锁定库存行
5. 校验商品状态
6. 校验库存是否充足
7. 扣减 `inventory.stock`
8. 写入 `orders`
9. 写入 `order_items`
10. 提交事务
11. 成功后把 Redis key 标记为 `SUCCESS:{order_no}`

## 支付状态流转流程

`POST /api/v1/payments/mock` 使用数据库事务完成支付状态流转：

1. 锁定订单
2. 校验订单属于当前用户
3. 校验订单状态为 `PENDING_PAYMENT`
4. 写入 `payments`
5. `SUCCESS` 时更新 `orders.status = PAID`
6. `FAILED` 时更新 `orders.status = PAYMENT_FAILED`
7. `FAILED` 时按 `order_items` 回补 `inventory.stock`
8. 提交事务

## 本地初始化方式

当前阶段使用 Docker Compose 启动 MySQL 8.0，并通过 `docs/db/schema.sql` 初始化数据库结构。

```powershell
docker compose up -d mysql
Get-Content docs/db/schema.sql | docker exec -i go-order-service-mysql mysql -uroot -prootpass
```

查看表：

```powershell
docker exec -it go-order-service-mysql mysql -uroot -prootpass -e "USE go_order_service; SHOW TABLES;"
```

## 后续演进

- 订单列表和订单详情接口仍可继续扩展
- 取消订单接口仍可继续扩展
- 目前 Redis 幂等 key 只覆盖订单创建，后续可以扩展到更多需要防重的写接口
- 支付流程当前是模拟流转，不接入第三方支付渠道

