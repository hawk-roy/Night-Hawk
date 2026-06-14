# API 文档

## 通用响应格式

项目内的核心接口统一返回以下结构：

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

- `code = 0` 表示成功
- `message` 是提示信息
- `data` 是业务数据，失败时通常为 `null`
- 当前实现里，常见错误码会直接沿用 HTTP 状态码语义，例如 `400 / 401 / 404 / 409 / 500`

## 通用响应头

### `X-Request-ID`

- 客户端可以传入 `X-Request-ID`
- 如果客户端没有传，服务端会生成一个
- 服务端会把同一个 `X-Request-ID` 写回响应头
- 请求日志会打印同一个 `request_id`

示例：

```text
X-Request-ID: test-request-0611-001
```

---

## 健康检查

### GET /api/v1/health

#### 成功响应

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "service": "go-order-service",
    "status": "ok"
  }
}
```

#### curl

```powershell
curl.exe http://localhost:8500/api/v1/health
```

---

## 数据库健康检查

### GET /api/v1/health/db

用于检查 Go 服务与 MySQL 的连接状态。

#### 成功响应

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "database": "mysql",
    "status": "ok"
  }
}
```

#### 数据库不可用响应

```json
{
  "code": 500,
  "message": "database unavailable",
  "data": null
}
```

#### curl

```powershell
curl.exe http://localhost:8500/api/v1/health/db
```

---

## Redis 健康检查

### GET /api/v1/health/redis

用于检查 Go 服务与 Redis 的连接状态。

#### 成功响应

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "cache": "redis",
    "status": "ok"
  }
}
```

#### Redis 不可用响应

```json
{
  "code": 500,
  "message": "redis unavailable",
  "data": null
}
```

#### curl

```powershell
curl.exe http://localhost:8500/api/v1/health/redis
```

---

## 用户注册

### POST /api/v1/users/register

#### 请求体

```json
{
  "username": "testuser",
  "password": "123456"
}
```

#### 成功响应

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "username": "testuser"
  }
}
```

#### 错误响应

```json
{
  "code": 400,
  "message": "invalid request",
  "data": null
}
```

```json
{
  "code": 400,
  "message": "username is required",
  "data": null
}
```

```json
{
  "code": 400,
  "message": "password is required",
  "data": null
}
```

#### curl

```powershell
curl.exe -X POST http://localhost:8500/api/v1/users/register `
  -H "Content-Type: application/json" `
  -d '{"username":"testuser","password":"123456"}'
```

---

## 用户登录

### POST /api/v1/users/login

#### 请求体

```json
{
  "username": "testuser",
  "password": "123456"
}
```

#### 成功响应

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "xxxxx.yyyyy.zzzzz"
  }
}
```

#### 常见错误

```json
{
  "code": 400,
  "message": "invalid request",
  "data": null
}
```

```json
{
  "code": 401,
  "message": "invalid username or password",
  "data": null
}
```

#### curl

```powershell
curl.exe -X POST http://localhost:8500/api/v1/users/login `
  -H "Content-Type: application/json" `
  -d '{"username":"testuser","password":"123456"}'
```

登录成功后，`cmd/apitest login` 会把 token 保存到 `.night-hawk-token`。

---

## 当前用户

### GET /api/v1/users/me

该接口受 JWT 鉴权保护，请求头需要携带：

```text
Authorization: Bearer xxxxx.yyyyy.zzzzz
```

#### 成功响应

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "user_id": 1,
    "username": "testuser"
  }
}
```

#### 未授权响应

```json
{
  "code": 401,
  "message": "unauthorized",
  "data": null
}
```

#### curl

```powershell
curl.exe -H "Authorization: Bearer xxxxx.yyyyy.zzzzz" http://localhost:8500/api/v1/users/me
```

---

## 商品列表

### GET /api/v1/products

直接从 MySQL 的 `products` 和 `inventory` 查询商品与库存，不需要 JWT。

#### 成功响应

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "Go Backend Course",
      "description": "A practical Go backend course",
      "price": 19900,
      "stock": 100
    }
  ]
}
```

#### curl

```powershell
curl.exe http://localhost:8500/api/v1/products
```

---

## 创建订单

### POST /api/v1/orders

创建订单接口需要 JWT 鉴权，并且必须携带 `Idempotency-Key`。

处理流程：

1. 查询商品与库存
2. 使用 `SELECT ... FOR UPDATE` 锁定库存行
3. 校验商品状态和库存数量
4. 扣减 `inventory.stock`
5. 写入 `orders`
6. 写入 `order_items`
7. 提交事务

#### 请求头

```text
Authorization: Bearer xxxxx.yyyyy.zzzzz
Idempotency-Key: <unique-request-key>
```

#### 请求体

```json
{
  "product_id": 1,
  "quantity": 2
}
```

#### 成功响应

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "order_no": "ORD1780894800552424800",
    "user_id": 3,
    "username": "orderuser_0608",
    "product_id": 1,
    "product_name": "Go Backend Course",
    "unit_price": 19900,
    "quantity": 2,
    "total_amount": 39800,
    "status": "PENDING_PAYMENT",
    "created_at": "2026-06-08T13:00:00.5524248+08:00"
  }
}
```

#### 常见错误

```json
{
  "code": 401,
  "message": "unauthorized",
  "data": null
}
```

```json
{
  "code": 400,
  "message": "idempotency key is required",
  "data": null
}
```

```json
{
  "code": 409,
  "message": "duplicate request",
  "data": null
}
```

```json
{
  "code": 500,
  "message": "idempotency service unavailable",
  "data": null
}
```

```json
{
  "code": 400,
  "message": "quantity must be greater than 0",
  "data": null
}
```

```json
{
  "code": 400,
  "message": "product_id must be greater than 0",
  "data": null
}
```

```json
{
  "code": 404,
  "message": "product not found",
  "data": null
}
```

```json
{
  "code": 400,
  "message": "insufficient stock",
  "data": null
}
```

#### curl

```powershell
curl.exe -X POST http://localhost:8500/api/v1/orders `
  -H "Content-Type: application/json" `
  -H "Authorization: Bearer xxxxx.yyyyy.zzzzz" `
  -H "Idempotency-Key: order-test-0610-001" `
  -d '{"product_id":1,"quantity":2}'
```

#### apitest

```powershell
go run ./cmd/apitest orders
go run ./cmd/apitest orders 1 2
```

#### 字段说明

| 字段 | 说明 |
| --- | --- |
| id | 订单 ID |
| order_no | 订单号 |
| user_id | 下单用户 ID |
| username | 下单用户名，仅用于返回展示 |
| product_id | 商品 ID |
| product_name | 商品名称 |
| unit_price | 商品单价，单位为分 |
| quantity | 购买数量 |
| total_amount | 订单总金额，单位为分 |
| status | 订单状态，目前固定为 `PENDING_PAYMENT` |
| created_at | 订单创建时间 |

---

## 请求示例说明

- `cmd/apitest db` 用于检查数据库连接
- `cmd/apitest redis` 用于检查 Redis 连接
- `cmd/apitest products` 用于检查商品列表
- `cmd/apitest register` 用于检查用户注册
- `cmd/apitest login` 用于检查用户登录并保存 JWT token
- `cmd/apitest me` 用于检查 JWT 鉴权后的当前用户接口
- `cmd/apitest orders` 用于检查订单创建、Redis 幂等和库存扣减事务
- `cmd/apitest orders` 会依次验证 `400 / 401 / 200 / 409`
- `cmd/apitest orders` 每次最多成功创建 1 个订单，重复请求会返回 `409`

---

## 请求日志

服务端已经接入请求日志中间件，会记录：

- `request_id`
- `method`
- `path`
- `status`
- `latency`
- `client_ip`

如果请求头没有传 `X-Request-ID`，服务端会生成一个并写回响应头。

---

## 支付状态变更

### POST /api/v1/payments/mock

用于模拟订单支付状态变更。这个接口不接入第三方支付，也不做真实回调验签，主要方便本地联调和测试。

#### 请求说明

- 需要 JWT 鉴权
- 请求体必须包含 `order_id` 和 `result`
- `result` 只允许 `SUCCESS` 或 `FAILED`

#### 请求头

```text
Authorization: Bearer <token>
Content-Type: application/json
```

#### 请求体

```json
{
  "order_id": 1,
  "result": "SUCCESS"
}
```

#### 成功响应

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "payment": {
      "id": 1,
      "payment_no": "PAY...",
      "order_id": 1,
      "amount": 39800,
      "status": "SUCCESS"
    },
    "order_status": "PAID"
  }
}
```

#### 失败响应

```json
{
  "code": 400,
  "message": "invalid payment result",
  "data": null
}
```

```json
{
  "code": 401,
  "message": "unauthorized",
  "data": null
}
```

```json
{
  "code": 404,
  "message": "order not found",
  "data": null
}
```

```json
{
  "code": 409,
  "message": "order is not pending payment",
  "data": null
}
```

```json
{
  "code": 500,
  "message": "internal server error",
  "data": null
}
```

#### curl

```powershell
curl.exe -X POST http://localhost:8500/api/v1/payments/mock `
  -H "Content-Type: application/json" `
  -H "Authorization: Bearer xxxxx.yyyyy.zzzzz" `
  -d '{"order_id":1,"result":"SUCCESS"}'
```

#### apitest

```powershell
go run ./cmd/apitest payments
go run ./cmd/apitest payments 1 2
```
