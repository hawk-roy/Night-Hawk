# 数据库 ER 说明

## 表关系

- `users` 1 -> N `orders`
- `products` 1 -> 1 `inventory`
- `orders` 1 -> N `order_items`
- `products` 1 -> N `order_items`
- `orders` 1 -> N `payments`

## 为什么要拆 `orders` 和 `order_items`

`orders` 保存订单主信息：用户、订单号、总金额、状态。  
`order_items` 保存商品明细：商品、单价、数量、小计。  
即使当前接口仍然只支持单商品下单，也先按真实订单系统的方式设计，方便后续扩展成多商品订单。

## 为什么 `inventory` 单独建表

`products` 保存商品基础信息。  
`inventory` 保存库存信息，库存会频繁更新，并且在下单时需要事务和行级锁，所以单独拆表更清晰。

## 为什么 `payments` 单独建表

`orders` 表示业务订单，`payments` 表示支付流水。  
一个订单可能对应多次支付尝试，所以支付信息应单独记录。

## 金额字段设计

所有金额字段都使用 `BIGINT`，单位为分，避免小数精度问题。

## 当前进度

Go 服务已经接入 MySQL，`/api/v1/health/db` 可用，`products`、`inventory` 已接入商品列表接口 `/api/v1/products`。  
当前 `POST /api/v1/orders` 也已经迁移到 MySQL，并在事务中完成库存校验、库存扣减、`orders` 写入和 `order_items` 写入。  
Redis 已经接入服务启动流程，并提供 `/api/v1/health/redis` 健康检查接口，用于验证 Redis 连接状态。

## 订单创建流程

`POST /api/v1/orders` 会使用数据库事务：

1. 查询商品和库存
2. 使用 `SELECT ... FOR UPDATE` 锁定库存行
3. 校验商品状态
4. 校验库存是否充足
5. 扣减 `inventory.stock`
6. 写入 `orders`
7. 写入 `order_items`
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

- `Redis` 幂等 key 还未接入
- 支付流程还未实现
- 订单取消、订单列表和订单详情接口还未实现
