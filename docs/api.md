# API 文档

## 通用响应格式

项目中的接口通常返回以下结构：

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

---

## 健康检查

### GET /api/v1/health

#### 成功响应

```json
{
  "service": "go-order-service",
  "status": "ok"
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

该接口直接从 MySQL 的 `products` 和 `inventory` 表查询商品与库存，不需要 JWT。

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

创建订单接口需要 JWT 鉴权。接口会在 MySQL 事务中完成：

1. 查询商品和库存
2. 使用 `SELECT ... FOR UPDATE` 锁定库存行
3. 校验商品状态和库存数量
4. 扣减 `inventory.stock`
5. 写入 `orders`
6. 写入 `order_items`
7. 提交事务

#### 请求头

```text
Authorization: Bearer xxxxx.yyyyy.zzzzz
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

#### 错误响应

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
| status | 订单状态，当前固定为 `PENDING_PAYMENT` |
| created_at | 订单创建时间 |

---

## 请求示例说明

- `cmd/apitest db` 用于检查数据库连接
- `cmd/apitest products` 用于检查商品列表
- `cmd/apitest register/login/me` 用于检查用户注册、登录和 JWT 鉴权
- `cmd/apitest orders` 用于检查订单创建与库存扣减事务
