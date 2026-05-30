# API 文档

## 健康检查

### 请求路径

```text
GET /api/v1/health
```

### 成功响应

```json
{
  "service": "go-order-service",
  "status": "ok"
}
```

### curl 验证命令

```powershell
curl.exe http://localhost:8080/api/v1/health
```

## 用户注册

### 请求路径

```text
POST /api/v1/users/register
```

### 请求体

```json
{
  "username": "testuser",
  "password": "123456"
}
```

### 成功响应

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

### 错误响应

```json
{
  "code": 400,
  "message": "invalid request"
}
```

### curl 验证命令

```powershell
curl.exe --% -X POST http://localhost:8080/api/v1/users/register -H "Content-Type: application/json" -d "{\"username\":\"testuser\",\"password\":\"123456\"}"
```

## 用户登录

### 请求路径

```text
POST /api/v1/users/login
```

### 请求体

```json
{
  "username": "testuser",
  "password": "123456"
}
```

### 成功响应

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "xxxxx.yyyyy.zzzzz"
  }
}
```

### 错误响应

缺少 username：

```json
{
  "code": 400,
  "message": "username is required"
}
```

缺少 password：

```json
{
  "code": 400,
  "message": "password is required"
}
```

用户不存在或密码错误：

```json
{
  "code": 401,
  "message": "invalid username or password"
}
```

token 生成失败：

```json
{
  "code": 500,
  "message": "failed to generate token"
}
```

### curl 验证命令

注册用户：

```powershell
curl.exe --% -X POST http://localhost:8080/api/v1/users/register -H "Content-Type: application/json" -d "{\"username\":\"testuser\",\"password\":\"123456\"}"
```

登录成功：

```powershell
curl.exe --% -X POST http://localhost:8080/api/v1/users/login -H "Content-Type: application/json" -d "{\"username\":\"testuser\",\"password\":\"123456\"}"
```

错误密码：

```powershell
curl.exe --% -i -X POST http://localhost:8080/api/v1/users/login -H "Content-Type: application/json" -d "{\"username\":\"testuser\",\"password\":\"wrong\"}"
```

缺少用户名：

```powershell
curl.exe --% -i -X POST http://localhost:8080/api/v1/users/login -H "Content-Type: application/json" -d "{\"password\":\"123456\"}"
```

缺少密码：

```powershell
curl.exe --% -i -X POST http://localhost:8080/api/v1/users/login -H "Content-Type: application/json" -d "{\"username\":\"testuser\"}"
```

## 获取当前用户信息

### 接口说明

获取当前登录用户信息。该接口受 JWT 鉴权中间件保护，请求时必须携带有效 token。

### 请求路径

```text
GET /api/v1/users/me
```

### 请求头

```text
Authorization: Bearer xxxxx.yyyyy.zzzzz
```

### 成功响应

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

### 错误响应

不带 token：

```json
{
  "code": 401,
  "message": "unauthorized",
  "data": null
}
```

错误 token：

```json
{
  "code": 401,
  "message": "unauthorized",
  "data": null
}
```

### curl 验证命令

不带 token：

```powershell
curl.exe -i http://localhost:8080/api/v1/users/me
```

错误 token：

```powershell
curl.exe --% -i http://localhost:8080/api/v1/users/me -H "Authorization: Bearer wrong-token"
```

正确 token：

```powershell
curl.exe --% -i http://localhost:8080/api/v1/users/me -H "Authorization: Bearer 这里替换成真实token"
```
