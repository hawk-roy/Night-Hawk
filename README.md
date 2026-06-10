# go-order-service

一个使用 Go + Gin 实现的电商订单后端服务。

## 当前进度

- [x] 初始 Go 项目
- [x] 接入 Gin
- [x] 实现健康检查接口 `/api/v1/health`
- [x] 用户注册接口
- [x] 注册接口参数校验
- [x] 注册接口错误返回
- [x] 用户登录接口
- [x] 登录成功返回 JWT token
- [x] JWT 鉴权中间件
- [x] 受保护接口 `/api/v1/users/me`
- [x] 商品列表接口 `/api/v1/products`
- [x] 订单创建接口
- [x] JWT 保护订单创建接口
- [x] 数据库表结构设计
- [x] Docker Compose 启动 MySQL
- [x] `schema.sql` 初始化
- [x] 用户注册/登录迁移到 MySQL
- [x] bcrypt 密码 hash 存储
- [x] 用户数据服务重启后仍可登录
- [x] 商品列表迁移到 MySQL
- [x] 商品 `seed.sql` 初始化数据
- [x] Go 服务接入 MySQL
- [x] `/api/v1/health/db`
- [x] 订单创建迁移到 MySQL
- [x] 库存扣减事务
- [x] Redis 接入 Docker Compose
- [x] Go 服务接入 Redis
- [x] Redis 健康检查接口 `/api/v1/health/redis`
- [x] Redis 幂等 key
- [x] 订单创建防重复提交

## 启动方式

```powershell
go run ./cmd/server
```

## 本地 MySQL 开发环境

### 1. 准备环境变量

```powershell
Copy-Item .env.example .env
```

### 2. 启动 MySQL

```powershell
docker compose up -d mysql
```

### 3. 初始化数据库

```powershell
Get-Content docs/db/schema.sql | docker exec -i go-order-service-mysql mysql -uroot -prootpass
```

### 4. 查看表

```powershell
docker exec -it go-order-service-mysql mysql -uroot -prootpass -e "USE go_order_service; SHOW TABLES;"
```

### 5. 验证 Go 服务连接 MySQL

```powershell
go run ./cmd/server
curl.exe http://localhost:8500/api/v1/health/db
```

## 启动 Redis

```powershell
docker compose up -d redis
docker exec -it go-order-service-redis redis-cli ping
```

期望返回：

```txt
PONG
```

## 验证 Redis 健康检查

```powershell
curl.exe http://localhost:8500/api/v1/health/redis
```

期望返回：

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

## 接口测试小工具

启动服务后，可以直接用项目内置的小工具验证接口：

```powershell
go run ./cmd/apitest health
go run ./cmd/apitest db
go run ./cmd/apitest redis
go run ./cmd/apitest products
go run ./cmd/apitest register JulieJaps 112233
go run ./cmd/apitest login JulieJaps 112233
go run ./cmd/apitest me
go run ./cmd/apitest orders
go run ./cmd/apitest orders 1 2
```

说明：

- `login` 成功后会把 JWT token 保存到 `.night-hawk-token`
- `.night-hawk-token` 已加入 `.gitignore`，不要提交
- `me` 会自动读取 `.night-hawk-token` 并访问受保护接口 `/api/v1/users/me`
- `orders` 会自动读取 `.night-hawk-token` 并访问受保护接口 `/api/v1/orders`
- `db` 用来验证数据库连接，不需要 JWT token
- `redis` 用来验证 Redis 连接，不需要 JWT token

## 订单创建

订单创建接口已经迁移到 MySQL，并接入 Redis 幂等 key。

- 接口会在事务中完成商品校验、库存校验、库存扣减、`orders` 写入和 `order_items` 写入
- 成功后会返回 `order_no`
- 当前订单状态固定为 `PENDING_PAYMENT`
- `POST /api/v1/orders` 必须携带 `Idempotency-Key` header
- 服务端使用 `userID + Idempotency-Key` 组成 Redis key
- 第一次请求会创建订单并扣减库存
- 重复请求会返回 `409 duplicate request`
- 重复请求不会再次扣减库存，也不会重复写入 `orders` 和 `order_items`

更多接口细节请看 [docs/api.md](docs/api.md)。
