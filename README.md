# go-order-service

一个使用 Go + Gin 实现的电商订单后端服务。

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
- [x] 创建订单接口雏形 `/api/v1/orders`
- [x] JWT 保护订单创建接口
- [ ] 数据库接入
- [ ] 库存扣减事务
- [ ] Redis 幂等 key

## 启动方式

```bash
go run ./cmd/server
```

## 接口测试小工具

启动服务后，可以用项目内置的小工具测试接口，避免手写复杂的 PowerShell curl 命令。

```bash
go run ./cmd/apitest health
go run ./cmd/apitest products
go run ./cmd/apitest register JulieJaps 112233
go run ./cmd/apitest login JulieJaps 112233
go run ./cmd/apitest me
go run ./cmd/apitest me-wrong
go run ./cmd/apitest orders
go run ./cmd/apitest orders 1 2
```

说明：

- `login` 成功后会把 JWT token 保存到 `.night-hawk-token`
- `.night-hawk-token` 已加入 `.gitignore`，不要提交
- `me` 会自动读取 `.night-hawk-token` 并访问受保护接口 `/api/v1/users/me`
- `orders` 会自动读取 `.night-hawk-token` 并访问受保护接口 `/api/v1/orders`

## 用户注册接口

接口路径：

```text
POST /api/v1/users/register
```

请求 JSON 示例：

```json
{
  "username": "testuser",
  "password": "123456"
}
```

curl 示例：

```bash
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"123456"}'
```

成功响应示例：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "user_id": 1,
    "username": "testuser"
  }
}
```
