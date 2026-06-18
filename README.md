# go-order-service

一个使用 Go + Gin 实现的订单后端服务，模拟电商核心链路，覆盖用户认证、商品查询、订单创建、库存扣减、Redis 幂等、支付状态流转、统一响应、请求日志和 Docker Compose 一键启动。

## 项目亮点

- JWT 用户认证：注册、登录、受保护接口访问
- MySQL 持久化：用户、商品、库存、订单、支付流水落库
- 订单库存事务：创建订单时使用 MySQL 事务和 `SELECT ... FOR UPDATE` 扣减库存
- Redis 幂等 key：使用 `Idempotency-Key` 防止重复下单
- 支付状态流转：支持 `PENDING_PAYMENT -> PAID / PAYMENT_FAILED`
- 统一响应结构：所有接口统一返回 `code/message/data`
- 请求日志：支持 `X-Request-ID`，记录 `method/path/status/latency/client_ip`
- Docker Compose：一键启动 MySQL、Redis 和 Go 服务
- 单元测试：覆盖 response、JWT、AuthMiddleware、RequestLogger、Redis 幂等逻辑
- MySQL 仓储集成测试：覆盖订单创建事务和支付状态流转事务

## 技术栈

- Go
- Gin
- MySQL 8.0
- Redis 7
- JWT
- Docker Compose
- `database/sql`
- `go-redis/v9`
- `miniredis`

## 架构设计

详见：[系统架构说明](docs/architecture.md)

其中包括：

- 系统架构图
- 订单创建流程图
- 支付状态流转图

## 当前进度

- [x] 初始化 Go 项目
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
- [x] `seed.sql` 初始化数据
- [x] Go 服务接入 MySQL
- [x] `/api/v1/health/db`
- [x] 订单创建迁移到 MySQL
- [x] 库存扣减事务
- [x] Redis 接入 Docker Compose
- [x] Go 服务接入 Redis
- [x] Redis 健康检查 `/api/v1/health/redis`
- [x] Redis 幂等 key
- [x] 订单创建防重复提交
- [x] 统一响应结构
- [x] 统一错误响应
- [x] 请求日志中间件
- [x] `X-Request-ID`
- [x] 支付状态流转

## 启动方式

```powershell
go run ./cmd/server
```

## Docker Compose 一键启动

启动 MySQL、Redis 和 Go 服务：

```powershell
docker compose up -d --build
```

查看容器：

```powershell
docker ps
```

期望看到：

```txt
go-order-service-mysql
go-order-service-redis
go-order-service-app
```

验证服务：

```powershell
curl.exe http://localhost:9000/api/v1/health
curl.exe http://localhost:9000/api/v1/health/db
curl.exe http://localhost:9000/api/v1/health/redis
curl.exe http://localhost:9000/api/v1/products
```

Docker Compose 模式访问地址：

```txt
http://localhost:9000
```

查看 app 日志：

```powershell
docker compose logs app
```

停止服务：

```powershell
docker compose down
```

注意：不要随便执行 `docker compose down -v`，它会删除 MySQL 数据卷。

### 容器内连接地址

Go app 容器中连接 MySQL 使用：

```txt
MYSQL_HOST=mysql
```

连接 Redis 使用：

```txt
REDIS_ADDR=redis:6379
```

不要在 app 容器中使用 `127.0.0.1` 连接 MySQL 或 Redis，因为容器内的 `127.0.0.1` 指向 app 容器自身。

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
curl.exe http://localhost:9000/api/v1/health/db
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

## 接口测试小工具

启动服务后，可以直接运行：

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
go run ./cmd/apitest payments
go run ./cmd/apitest payments 1 2
```

`cmd/apitest` 会自动尝试 `http://localhost:8080` 和 `http://localhost:9000`。如果你想手动指定地址，可以加 `-base`，例如：

```powershell
go run ./cmd/apitest -base http://localhost:9000 health
```

说明：
- `login` 成功后会把 JWT token 保存到 `.night-hawk-token`
- `.night-hawk-token` 已加入 `.gitignore`，不要提交
- `me` 会自动读取 `.night-hawk-token` 并访问受保护接口 `/api/v1/users/me`
- `orders` 会自动读取 `.night-hawk-token` 并访问受保护接口 `/api/v1/orders`
- `payments` 会自动读取 `.night-hawk-token` 并访问受保护接口 `/api/v1/payments/mock`
- `db` 用于验证数据库连接，不需要 JWT token
- `redis` 用于验证 Redis 连接，不需要 JWT token

## 核心接口

- `GET /api/v1/health`
- `GET /api/v1/health/db`
- `GET /api/v1/health/redis`
- `POST /api/v1/users/register`
- `POST /api/v1/users/login`
- `GET /api/v1/users/me`
- `GET /api/v1/products`
- `POST /api/v1/orders`
- `POST /api/v1/payments/mock`

## 测试

默认测试不依赖 MySQL / Redis 外部服务：

```powershell
go test ./...
```

当前测试覆盖：

- 统一响应结构 `response.Success / response.SuccessWithStatus / response.Error`
- JWT token 生成与解析
- JWT `AuthMiddleware` 鉴权行为
- 请求日志中间件和 `X-Request-ID`
- Redis `Idempotency-Key` 首次请求、重复请求、成功标记、失败释放、用户隔离
- MySQL 仓储集成测试，默认跳过

### MySQL 集成测试

Repository 集成测试会连接真实 MySQL，默认不运行。  
需要先启动 MySQL：

```powershell
docker compose up -d mysql
```

然后执行：

```powershell
$env:RUN_INTEGRATION_TESTS="1"
go test ./internal/repository -v
Remove-Item Env:RUN_INTEGRATION_TESTS
```

当前集成测试覆盖：

- OrderRepository 创建订单事务
- 创建订单时 `inventory.stock` 扣减
- 库存不足时事务回滚，不写 `orders` / `order_items`
- PaymentRepository 支付成功状态流转
- PaymentRepository 重复支付校验
- PaymentRepository 支付失败库存回补

## 接口回归验证

```powershell
go run ./cmd/apitest health
go run ./cmd/apitest db
go run ./cmd/apitest redis
go run ./cmd/apitest products
go run ./cmd/apitest register JulieJaps 112233
go run ./cmd/apitest login JulieJaps 112233
go run ./cmd/apitest me
go run ./cmd/apitest orders
go run ./cmd/apitest payments
```

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

## 支付状态流转

当前项目已支持模拟支付状态流转：
- 订单创建后状态为 `PENDING_PAYMENT`
- `POST /api/v1/payments/mock` 携带 `result=SUCCESS` 时，订单状态变为 `PAID`，并写入 `payments`
- `POST /api/v1/payments/mock` 携带 `result=FAILED` 时，订单状态变为 `PAYMENT_FAILED`，并恢复本次订单扣减的库存
- 同一订单重复支付会返回 `409`
- 支付接口需要 JWT 鉴权

## 当前边界

- 当前支付为模拟支付，不接入真实第三方支付
- 当前订单创建只支持单商品下单
- 当前重复 Idempotency-Key 返回 `409`，不返回第一次请求的完整响应
- 当前尚未接入消息队列
- 当前尚未做压测和 pprof 性能分析
