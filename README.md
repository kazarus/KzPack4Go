# KzPack4Go

Go 通用工具包,整合业务后端日常开发所需的基础能力:**Redis 客户端、JWT 鉴权、AES / 国密 SM4 加解密、SM2 密钥配置、哈希与时间/数值工具、ID / UUID 生成、服务与数据库配置实体、日志落盘、钉钉告警、统一响应体、并发安全黑名单**等。

基于 Go Modules 构建(`go 1.27+`),是 GOPATH 时代代码的现代化重构版本。

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
| `Redis4Go.go` | Redis 客户端封装,支持单机(`0/1`)与集群(`2`),提供 `NewClient` / `Ping` / `Close` / `Get` / `Set` / `Keys` 等常用操作 |
| `class_jwtauth.go` | JWT 令牌解析,`ParseToken1` / `ParseToken2` 仅接受 HMAC 签名,防算法混淆攻击 |
| `class_yzcrypt.go` | AES-CBC 与国密 SM4 加解密,`YzCrypt.AES` / `YzCrypt.SM4`,并保留兼容旧接口 |
| `class_enc_config.go` | SM2 密钥配置实体 `TEncConfig`,生成/加载国密密钥分量,`SM2PublicKeyBase64` 导出非压缩公钥 |
| `class_kzutils.go` | 哈希、时间戳、日期/数值格式化、字符串处理、文件读写、ID / UUID 生成等通用工具 |
| `class_kzerror.go` | 统一响应体 `TEROR` 及 `ToEROR` / `ToTRUE` / `ToPage` 等快捷构造方法 |
| `class_srv_config.go` | 服务配置实体 `TSrvConfig`(Redis/队列/钉钉/上传路径等) |
| `class_uni_config.go` | 数据库连接配置 `TUniConfig`,支持多种数据库 DSN 生成 |
| `class_sys_status.go` | 系统状态/日志实体 `TSYSSTATUS`(日志落盘、钉钉卡片) |
| `const_data_from.go` | 请求终端常量(`CONST_DATA_FROM_TO_*`:桌面/钉钉/微信/安卓/苹果) |
| `class_bot_config.go` | 钉钉机器人告警推送 `TBotConfig` |
| `class_black_list.go` | 并发安全黑名单 `TBlackList` |
| `class_kzvalid.go` | 流程校验响应 `TKzValid` |
| `class_kzdebug.go` | 简易计时调试 `TKzDebug` |
| `class_rmq_config.go` / `class_syn_config.go` | 消息队列 `TRmqConfig` / 数据同步 `TSynConfig` 配置实体 |

## 快速上手

### 1. Redis 客户端

```go
// 推荐:NewClientWithOptions,零值字段使用新默认(单机 DB 0 / 读写超时 3s)
r := KzPack4Go.NewClientWithOptions(0, "127.0.0.1:6379", "", "password", KzPack4Go.TRedisOptions{})
defer r.Close()

// 旧接口 NewClient 保留旧默认(单机 DB 11 / 读写超时 1000s),仅用于兼容历史调用点
// r := KzPack4Go.NewClient(0, "127.0.0.1:6379", "redis-user", "password")

// 集群(地址逗号分隔)
rc := KzPack4Go.NewClientWithOptions(2, "node1:6379,node2:6379,node3:6379", "", "password", KzPack4Go.TRedisOptions{})

// 自定义单机库号与读写超时
ro := KzPack4Go.NewClientWithOptions(0, "127.0.0.1:6379", "", "password", KzPack4Go.TRedisOptions{DB: 5, ReadTimeout: 3 * time.Second, WriteTimeout: 3 * time.Second})

ctx := context.Background()
r.Set(ctx, "key", "value", 0)
val, _ := r.Get(ctx, "key").Result()
keys := r.Keys(ctx, "user:*").Val() // 模式匹配
```

> `RedisTyp`: `0` / `1` 为单机模式,`2` 为集群模式;**其它值会在构造期直接 panic(快速失败)**,避免命令方法返回 nil 后引发远程 panic。
>
> 默认值分两套:`NewClientWithOptions` 零值字段使用新默认(单机 DB 0、读写超时 3s),**推荐新代码使用**;`NewClient`(已标 Deprecated)保留旧默认(单机 DB 11、读写超时 1000s)以兼容历史调用点,需要旧行为也可通过 `TRedisOptions{DB: 11, ReadTimeout: 1000 * time.Second, WriteTimeout: 1000 * time.Second}` 显式指定。

### 2. JWT 令牌解析

```go
claims, err := KzPack4Go.ParseToken1(tokenString, secret)
if err != nil {
    // 处理错误
}
fmt.Println(claims.EnrollID, claims.UserName)

usr, err := KzPack4Go.ParseToken2(tokenString, secret)
```

> 仅接受 HMAC 签名算法,其他算法(如 `none`、`RS256`)会被拒绝,防止算法混淆攻击;无效令牌返回错误而非静默零值。

### 3. 加解密

```go
// AES-CBC
crypted, _ := KzPack4Go.YzCrypt.AES.Encrypt([]byte(key16), []byte(origin))
plain, _ := KzPack4Go.YzCrypt.AES.Decrypt([]byte(key16), crypted)

// SM4(ECB,密钥为 16 字节 Base64)
cipher, _ := KzPack4Go.YzCrypt.SM4.Encrypt(b64Key, "hello")
text, _ := KzPack4Go.YzCrypt.SM4.Decrypt(b64Key, cipher)

// 兼容旧接口(旧格式:固定IV=密钥前16字节,⚠️ 仅用于历史数据/旧调用点)
KzPack4Go.EncryptAES(key, data)
KzPack4Go.DecryptAES(key, crypted)
KzPack4Go.EncryptSM4(key, data)
KzPack4Go.DecryptSM4(key, data)
KzPack4Go.SimplyEncryptAES(key, data) // 返回 hex 字符串
KzPack4Go.SimplyDecryptAES(key, hexText)
```

> `YzCrypt.AES` 使用 CBC 模式、PKCS7 填充,每次加密随机生成 16 字节 IV 并**前置**到密文头部(格式 `IV(16B)||密文`),避免固定 IV 导致相同明文前缀产生相同密文前缀;`YzCrypt.SM4` 使用 ECB 模式、密钥须为 16 字节。
>
> ⚠️ **兼容旧接口 `EncryptAES` / `DecryptAES` / `Simply*` 保留旧格式(固定 IV=密钥前 16 字节,不安全),仅用于解密历史数据与旧调用点,新数据请一律使用 `YzCrypt.AES`。两种格式不互通,请勿混用。**

### 4. SM2 密钥配置(TEncConfig)

```go
// 生成 SM2 密钥:D/X/Y 为 base64 编码的定宽 32 字节分量,
// EncryKey 为非压缩公钥(04||X||Y 共 65 字节,base64)
var enc KzPack4Go.TEncConfig
enc.GenerateKey()
KzPack4Go.ToFile("enc.json", enc) // JSON 落盘,自动带 cipherBy/metaData 标注

// 从 JSON 加载(分量完整性由加载方校验:公钥分量须与私钥推导一致)
KzPack4Go.InFile("enc.json", &enc)

// 对外输出前脱敏,清空密钥分量
enc.ClearAttr()
```

> 国密算法基于 `github.com/emmansun/gmsm`。`EncryKey` 交给客户端加密 SM4 会话密钥(SM2-C1C3C2,base64 后放 `X-GM-KEY` 请求头),服务端用 D/X/Y 分量重建私钥完成解密。密钥分量与 SM4 密文格式均与具体实现无关,历史密钥/密文升级后无需迁移。

### 5. 通用工具

```go
// 哈希(支持 md5 / sha1 / sha256 / sha512)
hash, err := KzPack4Go.ToHash("text", "sha256", false)

// 时间戳(紧凑数字格式)
ms   := KzPack4Go.ToMsec() // 20260822150405.000
date := KzPack4Go.ToDate() // 20260822
week := KzPack4Go.ToWeek() // 202634(ISO 年*100+周)

// 数值格式化
fmt.Println(KzPack4Go.Format2Float(3.14159)) // 3.14
fmt.Println(KzPack4Go.Float0Text(1.5000))     // 1.5

// 字符串判断
KzPack4Go.Contains("a,b,c", "b")          // true(逗号分隔串)
KzPack4Go.Contain("b", []string{"a","b"}) // true(切片)

// 三元风格选择
KzPack4Go.IfThenI(ok, 1, 0) // int64
KzPack4Go.IfThenS(ok, "y", "n")

// 文件工具
KzPack4Go.FileExist("cfg.json")          // (bool, error)
KzPack4Go.ToFile("cfg.json", cfg)        // JSON 存档
KzPack4Go.InFile("cfg.json", &cfg)       // 从 JSON 初始化(不存在则创建)
KzPack4Go.FormatFileSize(2048)           // "2.00KB"
```

> `ToHash` 仅支持 `md5` / `sha1` / `sha256` / `sha512`,其它类型返回空串与错误;`isHex=true` 时输入非法 hex 同样返回错误。

### 6. ID / UUID 生成

```go
// UUID(带花括号,多用于 EntityID)
uuid := KzPack4Go.ToUUID()        // {xxxxxxxx-...}
lite := KzPack4Go.ToLite(uuid)    // 1234...9012(简写便于对比)

// 16 位雪花 ID(未超 2^53,前端安全)
id := KzPack4Go.ToNext()

// 初始化 ID 生成器(WorkerID 须在 [0,63] 区间,默认 WorkerIdBitLength=6)
KzPack4Go.NewIdGenerator(5)
id = KzPack4Go.ToNext()
```

### 7. 统一响应体

```go
// 成功返回数据
KzPack4Go.ToTRUE(listData)

// 分页返回
KzPack4Go.ToPage(listData, len, page, size, total)

// 错误返回
KzPack4Go.ToEROR("发生错误")
KzPack4Go.ToErrorCnt("密码错误", "PWD", 3) // 携带错误次数
```

### 8. 配置实体

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

### 9. 钉钉告警与日志

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

### 10. 并发安全黑名单

```go
// 零值可直接使用(无 TTL、无容量上限)
var bl KzPack4Go.TBlackList
bl.Append("u001", KzPack4Go.TErrorUsr{UserCode: "u001", ErrorCnt: 1})
usr, ok := bl.Select("u001")
bl.Remove("u001")

// 长跑服务建议限制内存:容量上限 + 有效期(登录失败计数等场景)
bl2 := KzPack4Go.NewBlackList(10000, 10*time.Minute)
bl2.Append("u002", KzPack4Go.TErrorUsr{UserCode: "u002", ErrorCnt: 3})
```

> `MaxSize>0` 时达到上限按"最早到期/最早插入"淘汰旧条目;`TTL>0` 时条目过期后 `Select` 视为不存在,并在追加时按 `TTL/2` 间隔整表清扫,防止过期条目长期驻留。

## 注意事项

- `TSrvConfig` / `TUniConfig` 序列化时包含明文口令字段,对外输出前请先调用 `ClearAttr()` 脱敏。
- `FatalWhenError(err, msg)` 出错时**直接退出进程**(`log.Fatalln`),仅建议用于应用启动期的必要性检查,业务代码请自行处理错误。
- `ToHash` 仅支持 `md5` / `sha1` / `sha256` / `sha512`,其它类型与非法 hex 输入返回空串与错误。
- `SM4` 加解密密钥须为 16 字节(Base64 编码);`AES` 密钥长度须符合 AES 规范(16/24/32 字节)。
- `YzCrypt.AES` 密文为 `随机IV(16B)||密文` 格式;旧接口 `EncryptAES`/`DecryptAES`/`Simply*` 为旧格式(固定 IV=密钥前 16 字节),仅用于兼容历史数据,新数据请使用 `YzCrypt.AES`,两种格式不互通。
- `NewClient` 对非法 `RedisTyp` 构造期快速失败(panic)。默认值:`NewClientWithOptions` 零值字段使用新默认(单机 DB 0、读写超时 3s);`NewClient`(Deprecated)保留旧默认(单机 DB 11、读写超时 1000s)兼容历史调用点,新代码请勿依赖。
- 依赖第三方库:`github.com/redis/go-redis/v9`、`github.com/golang-jwt/jwt/v5`、`github.com/emmansun/gmsm`(国密 SM2/SM4)、`github.com/sirupsen/logrus`、`github.com/go-resty/resty/v2`、`github.com/yitter/idgenerator-go`、`github.com/gofrs/uuid/v5` 等。
