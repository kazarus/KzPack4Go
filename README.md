# KzPack4Go

Go 通用工具包,整合业务后端日常开发所需的基础能力:**Redis 客户端、JWT 鉴权、AES / 国密 SM4 加解密、哈希与时间/数值工具、ID / UUID 生成、服务与数据库配置实体、日志落盘、钉钉告警、统一响应体、并发安全黑名单**等。

基于 Go Modules 构建(`go 1.21+`),是 GOPATH 时代代码的现代化重构版本。

## 安装

```bash
go get github.com/kazarus/KzPack4Go
```

## 使用

```go
import "github.com/kazarus/KzPack4Go"
```

## 模块结构

| 文件 | 说明 |
|---|---|
| `Redis4Go.go` | Redis 客户端封装,支持单机(`0/1`)与集群(`2`),提供 `NewClient` / `Ping` / `Close` / `Get` / `Set` 等常用操作 |
| `class_jwtauth.go` | JWT 令牌解析,`ParseToken1` / `ParseToken2` 仅接受 HMAC 签名,防算法混淆攻击 |
| `class_yzcrypt.go` | AES-CBC 与国密 SM4 加解密,`YzCrypt.AES` / `YzCrypt.SM4`,并保留兼容旧接口 |
| `class_kzutils.go` | 哈希、时间戳、日期/数值格式化、字符串处理、文件读写、ID / UUID 生成等通用工具 |
| `class_kzerror.go` | 统一响应体 `TEROR` 及 `ToEROR` / `ToTRUE` / `ToPage` 等快捷构造方法 |
| `class_srv_config.go` | 服务配置实体 `TSrvConfig`(Redis/队列/钉钉/上传路径等) |
| `class_uni_config.go` | 数据库连接配置 `TUniConfig`,支持多种数据库 DSN 生成 |
| `class_sys_status.go` | 系统状态/日志实体 `TSYSSTATUS`(日志落盘、钉钉卡片) |
| `class_bot_config.go` | 钉钉机器人告警推送 `TBotConfig` |
| `class_black_list.go` | 并发安全黑名单 `TBlackList` |
| `class_kzvalid.go` | 流程校验响应 `TKzValid` |
| `class_kzdebug.go` | 简易计时调试 `TKzDebug` |
| `class_rmq_config.go` / `class_syn_config.go` | 消息队列 `TRmqConfig` / 数据同步 `TSynConfig` 配置实体 |

## 快速上手

### 1. Redis 客户端

```go
// 单机
r := KzPack4Go.NewClient(0, "127.0.0.1:6379", "", "password")
defer r.Close()

// 集群(地址逗号分隔)
rc := KzPack4Go.NewClient(2, "node1:6379,node2:6379,node3:6379", "", "password")

ctx := context.Background()
r.Set(ctx, "key", "value", 0)
val, _ := r.Get(ctx, "key").Result()
```

> `RedisTyp`: `0` / `1` 为单机模式,`2` 为集群模式。

### 2. JWT 令牌解析

```go
claims, err := KzPack4Go.ParseToken1(tokenString, secret)
if err != nil {
    // 处理错误
}
fmt.Println(claims.EnrollID, claims.UserName)

usr, err := KzPack4Go.ParseToken2(tokenString, secret)
```

> 仅接受 HMAC 签名算法,其他算法(如 `none`、`RS256`)会被拒绝,防止算法混淆攻击。

### 3. 加解密

```go
// AES-CBC
crypted, _ := KzPack4Go.YzCrypt.AES.Encrypt([]byte(key16), []byte(origin))
plain, _ := KzPack4Go.YzCrypt.AES.Decrypt([]byte(key16), crypted)

// SM4(ECB,密钥为 16 字节 Base64)
cipher, _ := KzPack4Go.YzCrypt.SM4.Encrypt(b64Key, "hello")
text, _ := KzPack4Go.YzCrypt.SM4.Decrypt(b64Key, cipher)

// 兼容旧接口
KzPack4Go.EncryptAES(key, data)
KzPack4Go.DecryptAES(key, crypted)
KzPack4Go.EncryptSM4(key, data)
KzPack4Go.DecryptSM4(key, data)
KzPack4Go.SimplyEncryptAES(key, data) // 返回 hex 字符串
KzPack4Go.SimplyDecryptAES(key, hexText)
```

> `YzCrypt.AES` 使用 CBC 模式、PKCS7 填充;`YzCrypt.SM4` 使用 ECB 模式、密钥须为 16 字节。

### 4. 通用工具

```go
// 哈希(支持 md5 / sha1 / sha256 / sha512)
hash := KzPack4Go.ToHash("text", "sha256", false)

// 时间戳
ms   := KzPack4Go.ToMsec() // 20260822150405.000
date := KzPack4Go.ToDate() // 20260822

// 数值格式化
fmt.Println(KzPack4Go.Format2Float(3.14159)) // 3.14
fmt.Println(KzPack4Go.Float0Text(1.5000))     // 1.5

// ID / UUID 生成
uuid := KzPack4Go.ToUUID()   // {xxxxxxxx-...}
id16 := KzPack4Go.ToNext()   // 16 位雪花 ID
id19 := KzPack4Go.ToNext19() // 19 位雪花 ID(超出 2^53)

// 字符串判断
KzPack4Go.Contains("a,b,c", "b") // true
KzPack4Go.Contain("b", []string{"a", "b"})
```

> `ToHash` 仅支持 `md5` / `sha1` / `sha256` / `sha512`,其它类型返回错误提示字符串。

### 5. 统一响应体

```go
// 成功返回数据
KzPack4Go.ToTRUE(listData)

// 分页返回
KzPack4Go.ToPage(listData, len, page, size, total)

// 错误返回
KzPack4Go.ToEROR("发生错误")
```

### 6. 配置实体

```go
// 服务配置
var cfg KzPack4Go.TSrvConfig
cfg.ToDefault()   // 填充默认值
cfg.ClearAttr()   // 对外输出前脱敏,清空口令等敏感字段

// 数据库配置
var uni KzPack4Go.TUniConfig
uni.UnicType = "PostgreSQL"
dsn := uni.UrlBuild() // 生成数据库连接串
uni.ClearAttr()       // 脱敏
```

> `TUniConfig.UrlBuild` 支持 `PostgreSQL` / `PolorDB` / `SQL Server` / `Oracle` / `MySQL` / `SQLite` / `dm`(达梦) / `opengauss`。

### 7. 钉钉告警与日志

```go
var bot KzPack4Go.TBotConfig
bot.EndPoint = "https://oapi.dingtalk.com/robot/send"
bot.AccessID = "access_token"
bot.SecretID = "secret"
bot.ToErrorMsg1(jsonStr) // 推送错误消息(钉钉机器人)

// 日志落盘 + 生成钉钉卡片
var st KzPack4Go.TSYSSTATUS
st.Initialize()
st.InTokenUsr(usr)
st.ToErrorLog(usr, cfg, bot, logger)
card := st.ToFeedCard1()
```

## 注意事项

- `TSrvConfig` / `TUniConfig` 序列化时包含明文口令字段,对外输出前请先调用 `ClearAttr()` 脱敏。
- `ToHash` 仅支持 `md5` / `sha1` / `sha256` / `sha512`,其它类型返回错误提示。
- `SM4` 加解密密钥须为 16 字节(Base64 编码);`AES` 密钥长度须符合 AES 规范(16/24/32 字节)。
- 依赖第三方库:`github.com/redis/go-redis/v9`、`github.com/golang-jwt/jwt/v5`、`github.com/tjfoc/gmsm`、`github.com/sirupsen/logrus`、`github.com/go-resty/resty`、`github.com/bwmarrin/snowflake`、`github.com/yitter/idgenerator-go`、`github.com/gofrs/uuid/v5` 等。
