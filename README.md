# go-order-service

一个使用 Go + Gin 实现的电商订单后端服务。

## 当前进度

- [x] 初始化 Go 项目
- [x] 接入 Gin
- [x] 实现健康检查接口 `/api/v1/health`

## 启动方式

```bash
go run ./cmd/server

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