package KzPack4Go

import (
	"encoding/json"
	"fmt"
)

// #服务配置
type TSrvConfig struct {
	CryptTxt string
	HttpPort string //#后端部署端口(旧)
	ListenAt string //#后端部署端口(新)
	SilentAt string //#后端静默端口(新)
	SocketAt string //#实时通讯中心(新)
	FrontEnd string //#前端部署端口
	JwtToken int64  //#启用请求令牌
	LockKJND int64  //#业务年度封存@作废
	RealKJND int64  //#当前业务年度@作废
	LockYear int64  //#业务年度封存#新的
	RealYear int64  //#当前业务年度#新的

	UseRedis int64 //#启用缓存服务(redis4go)
	UseError int64 //#启用日志收集
	AllError int64 //#收集所有日志
	UseQueue int64 //#启用消息队列(rabbitmq)
	ErrorCnt int64 //#密码容错次数
	CloseErr int64 //#关闭错误详情
	LockFlow int64 //#严格模式流程(谁终审谁撤回)
	WhereAmi int64 //#报告部署信息

	Path2PDF string //#文档转换服务(甬政开发)
	Path2TXT string //#图文转换服务(甬政开发)

	RedisTyp int64  //#缓存服务类型{0,1:单机;2:集群}
	RedisKey string //#缓存服务标识(默认监听端口)
	RedisSrv string //#缓存服务地址
	RedisPwd string //#缓存服务密码
	RedisUsr string //#缓存服务用户

	QueueSrv string //#队列服务地址
	QueuePwd string //#队列服务密码
	QueueUsr string //#队列服务用户

	PrintURL string //#打印服务地址
	EndPoint string //#相对打印服务的后端部署地址

	CloudURL string //#服务部署地址(外网地址)
	Name4App string //#服务部署名称(外网地址)

	Ding4URL string //#钉钉后台服务(甬政开发)
	Ding4APP string //#钉钉项目标识(甬政开发)
	Ding4GOV string //#钉钉官网地址(钉钉官网)
	Ding4Dom string //#钉钉发布域名
	Ding4VER int64  //#钉钉版本区分{1:政务钉钉;2:企业钉钉}
	TryAgain int64  //#重试钉钉发送

	Log4goAt string //#日志存放目录
	UploadAt string //#附件存放目录
	UpFormat int64  //#附件存放格式{0:默认;1:按年;2:按月;3:按周}

	SynzFrom string //#数据同步方向{"->":内到外,"<-":外到内}
	WorkerID int64  //#硬件唯一标识//#https://github.com/yitter/idgenerator

	ThirdURL string //#三方协作地址;etc:内云控
	InnerURL string //#内部协作地址;etc:老内控

	LoginOut int64 //#启用登录超时
	LoginMin int64 //#登录超时分钟

	SecretOn int64 //#开启敏感加密
	ISilence int64 //#静默创建视图
	UsesCron int64 //#启用定时任务
	Delay10s int64 //#消息延时X秒,智慧教育延时10秒,其他默认0秒;

	MaxOpenConns int64 //#最大开放链接
	MaxIdleConns int64 //#最大空闲链接
	MaxLifetimes int64 //#最长保留时间
	MaxIdleTimes int64 //#最长空闲时间

	Path4ExtExcludes string //#禁止附件类型:以逗号隔开
	PathNameExcludes string //#禁止附件名称:以逗号隔开

	UseRedisWhenExpd int64 //#写入修改单据:静默创建缓存
	UpdateOnWhenExpd int64 //#允许修改单据:放开送审锁定

	MakeYearWhenTodo int64 //#待办申请启用年度.{默认:1}
	VoidYearWhenTodo int64 //#待办审核启用年度.{默认:0}

	AccessControlAllowOrigin string
}

func (this TSrvConfig) MarshalJSON() ([]byte, error) {
	var result = make(map[string]interface{})
	result["cryptTxt"] = this.CryptTxt
	result["httpPort"] = this.HttpPort
	result["listenAt"] = this.ListenAt
	result["silentAt"] = this.SilentAt
	result["socketAt"] = this.SocketAt
	result["frontEnd"] = this.FrontEnd
	result["jwtToken"] = this.JwtToken
	result["lockKJND"] = this.LockKJND
	result["realKJND"] = this.RealKJND
	result["lockYear"] = this.LockYear
	result["realYear"] = this.RealYear
	result["useRedis"] = this.UseRedis
	result["useError"] = this.UseError
	result["allError"] = this.AllError
	result["useQueue"] = this.UseQueue
	result["errorCnt"] = this.ErrorCnt
	result["closeErr"] = this.CloseErr
	result["lockFlow"] = this.LockFlow
	result["whereAmi"] = this.WhereAmi
	result["path2PDF"] = this.Path2PDF
	result["path2TXT"] = this.Path2TXT
	result["redisTyp"] = this.RedisTyp
	result["redisKey"] = this.RedisKey
	result["redisSrv"] = this.RedisSrv
	result["redisPwd"] = this.RedisPwd
	result["redisUsr"] = this.RedisUsr
	result["queueSrv"] = this.QueueSrv
	result["queuePwd"] = this.QueuePwd
	result["queueUsr"] = this.QueueUsr
	result["printURL"] = this.PrintURL
	result["endPoint"] = this.EndPoint
	result["cloudURL"] = this.CloudURL
	result["name4App"] = this.Name4App
	result["ding4URL"] = this.Ding4URL
	result["ding4APP"] = this.Ding4APP
	result["ding4GOV"] = this.Ding4GOV
	result["ding4Dom"] = this.Ding4Dom
	result["ding4VER"] = this.Ding4VER
	result["tryAgain"] = this.TryAgain
	result["log4goAt"] = this.Log4goAt
	result["uploadAt"] = this.UploadAt
	result["upFormat"] = this.UpFormat
	result["synzFrom"] = this.SynzFrom
	result["workerID"] = this.WorkerID
	result["thirdURL"] = this.ThirdURL
	result["innerURL"] = this.InnerURL
	result["loginOut"] = this.LoginOut
	result["loginMin"] = this.LoginMin
	result["secretOn"] = this.SecretOn
	result["iSilence"] = this.ISilence
	result["usesCron"] = this.UsesCron
	result["delay10s"] = this.Delay10s

	result["maxOpenConns"] = this.MaxOpenConns
	result["maxIdleConns"] = this.MaxIdleConns
	result["maxLifetimes"] = this.MaxLifetimes
	result["maxIdleTimes"] = this.MaxIdleTimes

	result["path4ExtExcludes"] = this.Path4ExtExcludes
	result["pathNameExcludes"] = this.PathNameExcludes
	result["useRedisWhenExpd"] = this.UseRedisWhenExpd
	result["updateOnWhenExpd"] = this.UpdateOnWhenExpd
	result["makeYearWhenTodo"] = this.MakeYearWhenTodo
	result["voidYearWhenTodo"] = this.VoidYearWhenTodo

	result["accessControlAllowOrigin"] = this.AccessControlAllowOrigin

	return json.Marshal(result)
}

// # 返回路径右侧带/
func (this *TSrvConfig) ToUpFormat(WhoBuild int64) string {

	var result string

	switch this.UpFormat {

	case 1:
		{
			result = fmt.Sprintf("%s/%d/%d/", this.UploadAt, WhoBuild, ToKJND()*100+1)
		}
	case 2:
		{
			result = fmt.Sprintf("%s/%d/%d/", this.UploadAt, WhoBuild, ToKJQJ())
		}
	case 3:
		{
			result = fmt.Sprintf("%s/%d/%d/", this.UploadAt, WhoBuild, ToWeek())
		}
	default:
		{
			result = fmt.Sprintf("%s/%d/", this.UploadAt, WhoBuild)
		}
	}

	return result
}

// # 返回路径右侧带/
func (this *TSrvConfig) ToRealPath(WhoBuild int64) string {

	var result string

	switch this.UpFormat {

	case 1:
		{
			result = fmt.Sprintf("%d/%d/", WhoBuild, ToKJND()*100+1)
		}
	case 2:
		{
			result = fmt.Sprintf("%d/%d/", WhoBuild, ToKJQJ())
		}
	case 3:
		{
			result = fmt.Sprintf("%d/%d/", WhoBuild, ToWeek())
		}
	default:
		{
			result = fmt.Sprintf("%d/", WhoBuild)
		}
	}

	return result
}

// #初始化默认属性
func (this *TSrvConfig) ToDefault() {
	//#默认配置
	this.UploadAt = IfThenS(this.UploadAt == "", "upload", this.UploadAt)
	this.Log4goAt = IfThenS(this.Log4goAt == "", "log4go", this.Log4goAt)
	this.WorkerID = IfThenI(this.WorkerID == 0, 1, this.WorkerID)
	this.SynzFrom = IfThenS(this.SynzFrom == "", "NTOW", this.SynzFrom)
	this.LoginMin = IfThenI(this.LoginMin == 0, 9, this.LoginMin)

	this.CloudURL = IfThenS(this.CloudURL == "", this.EndPoint, this.CloudURL)
	this.ListenAt = IfThenS(this.ListenAt == "", this.HttpPort, this.ListenAt)
	this.SilentAt = IfThenS(this.SilentAt == "", this.HttpPort, this.SilentAt)
	this.RedisKey = IfThenS(this.RedisKey == "", this.HttpPort, this.RedisKey)

	this.MaxOpenConns = IfThenI(this.MaxOpenConns == 0, 500, this.MaxOpenConns)
	this.MaxIdleConns = IfThenI(this.MaxIdleConns == 0, 200, this.MaxIdleConns)
	this.MaxLifetimes = IfThenI(this.MaxLifetimes == 0, 60, this.MaxLifetimes)
	//this.MaxIdleTimes = IfThenI(this.MaxIdleTimes == 0, 10, this.MaxIdleTimes)

	//#配置升级.默认赋值
	this.Path4ExtExcludes = IfThenS(this.Path4ExtExcludes == "", ".html,.js,.ts,.exe", this.Path4ExtExcludes)
	this.AccessControlAllowOrigin = IfThenS(this.AccessControlAllowOrigin == "", "*", this.AccessControlAllowOrigin)

	//#20250110
	this.LockYear = this.LockKJND
	this.RealYear = this.RealKJND
	this.MakeYearWhenTodo = IfThenI(this.MakeYearWhenTodo == 0, 1, this.MakeYearWhenTodo) //#申请用户必传年度.{默认:1,配置文件可覆盖}
	this.VoidYearWhenTodo = IfThenI(this.VoidYearWhenTodo == 0, 0, this.VoidYearWhenTodo) //#审核用户不传年度.{默认:0,配置文件可覆盖}
}

func (this *TSrvConfig) ClearAttr() {

	this.EndPoint = ""
	this.CloudURL = ""
	this.PrintURL = ""

	this.FrontEnd = ""

	this.HttpPort = ""
	this.ListenAt = ""

	this.Ding4APP = ""
	this.Ding4Dom = ""
	this.Ding4GOV = ""
	this.Ding4URL = ""

	this.Path4ExtExcludes = ""
	this.RealYear = 0
	this.RedisKey = ""
	this.RedisPwd = ""
	this.RedisSrv = ""
	this.RedisUsr = ""
	this.QueueSrv = ""
	this.QueuePwd = ""
	this.QueueUsr = ""
	this.SocketAt = ""
	this.SilentAt = ""

	this.UploadAt = ""
	this.Log4goAt = ""
}

func (this *TSrvConfig) ToRabbitMQ(queueName string) string {

	return fmt.Sprintf("tech.cxcloud.%s.%s", this.HttpPort, queueName)
}

func (this *TSrvConfig) GetMakeYearWhenTodo(SourceND int64) int64 {

	var result int64

	switch this.MakeYearWhenTodo {
	case 1:
		{
			result = SourceND
		}
	default:
		{
			result = -1
		}
	}

	return result
}

func (this *TSrvConfig) GetVoidYearWhenTodo(SourceND int64) int64 {

	var result int64

	switch this.VoidYearWhenTodo {
	case 1:
		{
			result = SourceND
		}
	default:
		{
			result = -1
		}
	}

	return result
}
