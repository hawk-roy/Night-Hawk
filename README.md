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
- [ ] 商品列表接口
- [ ] 创建订单接口

## 启动方式

```bash
go run ./cmd/server
```

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