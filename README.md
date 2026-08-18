# KzPack4Go

Go 通用工具包:封装 Redis、JWT、AES/SM4 加解密、哈希、时间/数值处理、服务与数据库配置实体、钉钉告警、错误响应体、日志落盘、黑名单、ID/UUID 生成等常用组件。

## 使用

```go
import "github.com/kazarus/KzPack4Go"
```

## 模块结构

| 文件 | 说明 |
|---|---|
| `Redis4Go.go` | Redis 客户端封装(单机/集群,`NewClient` / `Close`) |
| `class_jwtauth.go` | JWT 令牌解析(`ParseToken1` / `ParseToken2`,仅接受 HMAC 签名) |
| `class_yzcrypt.go` | AES-CBC 与国密 SM4 加解密(`YzCrypt.AES` / `YzCrypt.SM4`) |
| `class_kzutils.go` | 哈希、时间戳、数值格式化、文件读写、通用工具函数 |
| `class_kzerror.go` | 统一响应体 `TEROR`、UUID / 雪花 ID 生成、错误工具 |
| `class_srv_config.go` | 服务配置实体 `TSrvConfig`(Redis/队列/钉钉/上传路径等) |
| `class_uni_config.go` | 数据库连接配置 `TUniConfig`(多数据库 DSN 生成) |
| `class_sys_status.go` | 系统状态/日志实体 `TSYSSTATUS`(日志落盘、钉钉卡片) |
| `class_bot_config.go` | 钉钉机器人告警推送 |
| `class_black_list.go` | 并发安全黑名单 `TBlackList` |
| `class_kzvalid.go` | 流程校验响应 `TKzValid` |
| `class_kzdebug.go` | 简易计时调试 `TKzDebug` |
| `class_rmq_config.go` / `class_syn_config.go` | 消息队列 / 数据同步配置实体 |

## 注意事项

- `TSrvConfig` / `TUniConfig` 序列化时包含明文口令字段,对外输出前请先调用 `ClearAttr()` 脱敏。
- `ToHash` 仅支持 `md5` / `sha1` / `sha256` / `sha512`,其它类型返回错误提示。
- 本包为 GOPATH 时代代码的现代化版本,现使用 Go Modules(`go 1.21+`)。
