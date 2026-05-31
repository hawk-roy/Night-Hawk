# JWT 三段式结构理解

JWT 最终长这样：

```txt
xxxxxx.yyyyyy.zzzzzz
```

它由三段组成：

```txt
header.payload.signature
```

生成过程可以理解为：

```txt
headerEncoded = Base64URL(header)
payloadEncoded = Base64URL(payload)

signatureRaw = HMACSHA256(
  headerEncoded + "." + payloadEncoded,
  serverSecret   
)

signatureEncoded = Base64URL(signatureRaw)

token = headerEncoded + "." + payloadEncoded + "." + signatureEncoded
```

其中：

- `xxxxxx` 是第一段，也就是 `header` 编码后的结果。
- `yyyyyy` 是第二段，也就是 `payload` 编码后的结果。
- `zzzzzz` 是第三段，也就是签名结果再编码后的结果。

`header` 通常保存 token 类型和签名算法，例如：

```json
{
  "alg": "HS256",
  "typ": "JWT"
}
```

`payload` 保存用户信息或业务声明，例如：

```json
{
  "user_id": 1,
  "username": "JulieJaps",
  "iat": 1710000000,
  "exp": 1710086400
}
```

第三段 `signature` 的作用是防篡改。

服务器会把前两段：

```txt
xxxxxx.yyyyyy
```

再加上服务器自己的密钥 `serverSecret`，通过签名算法计算出签名：

```txt
HMACSHA256(xxxxxx.yyyyyy, serverSecret)
```

然后把签名结果编码成第三段 `zzzzzz`。

所以 JWT 的核心不是“加密”，而是“编码 + 签名”：

- 前两段只是 Base64URL 编码，别人可以解码看到内容。
- 第三段是服务器用密钥计算出来的签名，用来证明前两段没有被修改。
- 如果别人修改了 payload，比如把 `username` 改成 `admin`，服务器重新计算签名时会发现第三段对不上，于是 token 验证失败，返回 `401`。

一句话总结：

> JWT 前两段负责携带信息，第三段负责证明前两段没有被改过。
