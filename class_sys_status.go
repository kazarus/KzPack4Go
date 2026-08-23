package KzPack4Go

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/sjson"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type TSYSSTATUS struct {
	WhoBuild int64   `db:"who_build" json:"whoBuild"`
	UserIndx int64   `db:"user_indx" json:"userIndx"`
	StatusID int64   `db:"status_id" json:"statusID"`
	HashCode int64   `db:"hash_code" json:"hashCode"`
	DataLevl int64   `db:"data_levl" json:"dataLevl"`
	DataFrom int64   `db:"data_from" json:"dataFrom"`
	DataType int64   `db:"data_type" json:"dataType"`
	KJNDKJQJ int64   `db:"kjnd_kjqj" json:"kjndKJQJ"`
	DataStat int64   `db:"data_stat" json:"dataStat"`
	DataOrdr int64   `db:"data_ordr" json:"dataOrdr"`
	DataHide int64   `db:"data_hide" json:"dataHide"`
	DataDate int64   `db:"data_date" json:"dataDate"`
	DataTime float64 `db:"data_time" json:"dataTime"`
	DataGUID string  `db:"data_guid" json:"dataGUID"`
	MastName string  `db:"mast_name" json:"mastName"`
	MastGUID string  `db:"mast_guid" json:"mastGUID"`
	CreateBy string  `db:"create_by" json:"createBy"`
	CreateAt float64 `db:"create_at" json:"createAt"`
	UpdateBy string  `db:"update_by" json:"updateBy"`
	UpdateAt float64 `db:"update_at" json:"updateAt"`
	DeleteBy string  `db:"delete_by" json:"deleteBy"`
	DeleteAt float64 `db:"delete_at" json:"deleteAt"`
	QueuesBy int64   `db:"queues_by" json:"queuesBy"`
	EntityID string  `db:"entity_id" json:"entityID"`
	EnrollID string  `db:"enroll_id" json:"enrollID"`
	Name4App string  `db:"name_4app" json:"name4App"`
	HttpPort string  `db:"http_port" json:"httpPort"` //#后端端口
	FrontEnd string  `db:"front_end" json:"frontEnd"` //#前端端口
	EndPoint string  `db:"end_point" json:"endPoint"`
	ClientIP string  `db:"client_ip" json:"clientIP"`
	StatusCD int64   `db:"status_cd" json:"statusCD"`
	SpendDat float64 `db:"spend_dat" json:"spendDat"`
	UserUUID string  `db:"user_uuid" json:"userUUID"`
	RqMethod string  `db:"rq_method" json:"rqMethod"`
	Url4Path string  `db:"url_4path" json:"url4Path"` //#请求链接
	ErrLevel int64   `db:"err_level" json:"errLevel"` //#错误等级
	ErrorMsg string  `db:"error_msg" json:"errorMsg"` //#错误消息
	ReturnID string  `db:"return_id" json:"returnID"` //#处理结果
	BrowerOS string  `db:"brower_os" json:"browerOS"` //#浏览器类型
	BrowerID string  `db:"brower_id" json:"browerID"` //#浏览器版本
	DiskPath string  `db:"disk_path" json:"diskPath"` //#物理路径
	Path4URL string  `db:"path_4url" json:"path4URL"` //#日志路径
	DataBase string  `db:"data_base" json:"dataBase"` //#数据链接
	BkColour string  `db:"bk_colour" json:"bkColour"` //#背景颜色
	AsReally int64   `db:"as_really" json:"asReally"` //#是否生产
	SysCustc int64   `db:"sys_custc" json:"sysCustc"` //#系统客户
	EnvrName string  `db:"envr_name" json:"envrName"` //#环境名称
	YesIknow int64   `db:"yes_iknow" json:"yesIknow"` //#我已知晓
	UserMain int64   `db:"user_main" json:"userMain"` //#用户标识
	ActionID int64   `db:"action_id" json:"actionID"` //#操作类型
	BillUUID string  `db:"bill_uuid" json:"billUUID"` //#单据标识
	RqStream int64   `db:"rq_stream" json:"rqStream"` //#请求表体

	Path4Log string `db:"path_4log,readonly" json:"path4Log"`
	Path4Req string `db:"path_4req,readonly" json:"path4Req"`
	RqSource string `db:"rq_source,readonly" json:"rqSource"`
	FromText string `db:"from_text,readonly" json:"fromText"`
	MsgTitle string `db:"msg_title,readonly" json:"msgTitle"`
	MsgState string `db:"msg_state,readonly" json:"msgState"`
}

func (this TSYSSTATUS) MarshalJSON() ([]byte, error) {
	var result = make(map[string]interface{})
	result["whoBuild"] = this.WhoBuild
	result["userIndx"] = this.UserIndx
	result["statusID"] = this.StatusID
	result["hashCode"] = this.HashCode
	result["dataLevl"] = this.DataLevl
	result["dataFrom"] = this.DataFrom
	result["dataType"] = this.DataType
	result["kjndKJQJ"] = this.KJNDKJQJ
	result["dataStat"] = this.DataStat
	result["dataOrdr"] = this.DataOrdr
	result["dataHide"] = this.DataHide
	result["dataDate"] = this.DataDate
	result["dataTime"] = this.DataTime
	result["dataGUID"] = this.DataGUID

	result["mastName"] = this.MastName
	result["mastGUID"] = this.MastGUID
	result["createBy"] = this.CreateBy
	result["createAt"] = this.CreateAt
	result["updateBy"] = this.UpdateBy
	result["updateAt"] = this.UpdateAt
	result["deleteBy"] = this.DeleteBy
	result["deleteAt"] = this.DeleteAt
	result["queuesBy"] = this.QueuesBy
	result["entityID"] = this.EntityID
	result["enrollID"] = this.EnrollID
	result["name4App"] = this.Name4App
	result["endPoint"] = this.EndPoint
	result["httpPort"] = this.HttpPort
	result["frontEnd"] = this.FrontEnd
	result["clientIP"] = this.ClientIP
	result["statusCD"] = this.StatusCD
	result["spendDat"] = this.SpendDat
	result["userUUID"] = this.UserUUID
	result["rqMethod"] = this.RqMethod
	result["url4Path"] = this.Url4Path
	result["errLevel"] = this.ErrLevel
	result["errorMsg"] = this.ErrorMsg
	result["returnID"] = this.ReturnID
	result["browerOS"] = this.BrowerOS
	result["browerID"] = this.BrowerID
	result["diskPath"] = this.DiskPath
	result["path4URL"] = this.Path4URL
	result["dataBase"] = this.DataBase
	result["bkColour"] = this.BkColour
	result["asReally"] = this.AsReally
	result["sysCustc"] = this.SysCustc
	result["envrName"] = this.EnvrName
	result["yesIknow"] = this.YesIknow
	result["userMain"] = this.UserMain
	result["actionID"] = this.ActionID
	result["billUUID"] = this.BillUUID
	result["rqStream"] = this.RqStream

	result["path4Log"] = this.Path4Log
	result["path4Req"] = this.Path4Req
	result["rqSource"] = this.RqSource
	result["fromText"] = this.FromText
	result["msgTitle"] = this.MsgTitle
	result["msgState"] = this.MsgState

	result["metaData"] = "实体类名: core.TSysStatus"

	return json.Marshal(result)
}

func (this *TSYSSTATUS) Initialize() {

	this.HashCode = 0
	this.DataStat = 1

	this.DataHide = 0
	this.DataFrom = 0

	this.KJNDKJQJ = -1
	this.DataGUID = ToUUID()
	this.DataDate = ToDate()
	this.DataTime = ToTime()
	this.EntityID = ToUUID()
}

func (this *TSYSSTATUS) InTokenUsr(TokenUsr TTokenUsr) {

	this.MastName = TokenUsr.UserName
	this.EnrollID = TokenUsr.EnrollID
	this.MastGUID = TokenUsr.EntityID

	this.CreateBy = TokenUsr.EntityID
	this.CreateAt = ToTime()

	this.UpdateBy = TokenUsr.EntityID
	this.UpdateAt = ToTime()

	this.DataFrom = TokenUsr.DataFrom
}

func (this TSYSSTATUS) ToMarkDown0(msgTitle string, msgState string) string {

	r := make(map[string]interface{})
	r["msgtype"] = "markdown"
	r["markdown"] = map[string]string{
		"title": "日志提醒",
		"text": fmt.Sprintf("#### %s\n > 服务名称: %s\n\n > 后端端口: %s\n\n > 前端端口: %s\n\n > 数据链接: %s\n\n > 日志提醒: %s;\n\n",
			msgTitle,
			this.Name4App,
			this.HttpPort,
			this.FrontEnd,
			this.DataBase,
			msgState,
		),
	}

	result, _ := json.Marshal(r)
	return string(result)
}

func (this TSYSSTATUS) ToMarkDown1() string {

	r := make(map[string]interface{})
	r["msgtype"] = "markdown"
	r["markdown"] = map[string]string{
		"title": "日志提醒",
		"text": fmt.Sprintf("#### 服务重启\n > 服务名称: %s\n\n > 后端端口: %s\n\n > 前端端口: %s\n\n > 数据链接: %s\n\n > 日志提醒: 数据链接变更频繁请相关人员留意;\n\n",
			this.Name4App,
			this.HttpPort,
			this.FrontEnd,
			this.DataBase,
		),
	}

	result, _ := json.Marshal(r)
	return string(result)
}

func (this TSYSSTATUS) ToMarkDown2() string {

	r := make(map[string]interface{})
	r["msgtype"] = "markdown"
	r["markdown"] = map[string]string{
		"title": "日志提醒",
		"text": fmt.Sprintf("#### 服务暂停\n > 服务名称: %s\n\n > 后端端口: %s\n\n > 前端端口: %s\n\n > 数据链接: %s\n\n > 日志提醒: 后端日常运维,服务暂停;\n\n",
			this.Name4App,
			this.HttpPort,
			this.FrontEnd,
			this.DataBase,
		),
	}

	result, _ := json.Marshal(r)
	return string(result)
}

func (this TSYSSTATUS) ToFeedCard1() string {

	var JsonText string
	JsonText, _ = sjson.Set(JsonText, "msgtype", "feedCard")
	JsonText, _ = sjson.Set(JsonText, "feedCard.links.-1", map[string]interface{}{"messageURL": this.Path4URL, "picURL": this.BkColour, "title": fmt.Sprintf("日志标识: %s", ToLite(this.EntityID))})
	JsonText, _ = sjson.Set(JsonText, "feedCard.links.-1", map[string]interface{}{"messageURL": this.Path4URL, "title": fmt.Sprintf("服务标识: %s", this.Name4App)})
	JsonText, _ = sjson.Set(JsonText, "feedCard.links.-1", map[string]interface{}{"messageURL": this.Path4URL, "title": fmt.Sprintf("请求链接: %s", this.Url4Path)})
	JsonText, _ = sjson.Set(JsonText, "feedCard.links.-1", map[string]interface{}{"messageURL": this.Path4URL, "title": fmt.Sprintf("请求地址: %s", this.ClientIP)})
	JsonText, _ = sjson.Set(JsonText, "feedCard.links.-1", map[string]interface{}{"messageURL": this.Path4URL, "title": fmt.Sprintf("后端端口: %s(%s)", this.EndPoint, this.HttpPort)})
	JsonText, _ = sjson.Set(JsonText, "feedCard.links.-1", map[string]interface{}{"messageURL": this.Path4URL, "title": fmt.Sprintf("前端端口: %s", this.FrontEnd)})
	JsonText, _ = sjson.Set(JsonText, "feedCard.links.-1", map[string]interface{}{"messageURL": this.Path4URL, "title": fmt.Sprintf("会话用户: %s(%s)", this.MastName, this.EnrollID)})
	JsonText, _ = sjson.Set(JsonText, "feedCard.links.-1", map[string]interface{}{"messageURL": this.Path4URL, "title": fmt.Sprintf("数据连接: %s", this.DataBase)})
	JsonText, _ = sjson.Set(JsonText, "feedCard.links.-1", map[string]interface{}{"messageURL": this.Path4URL, "title": fmt.Sprintf("请求终端: %s", this.FromText)})
	JsonText, _ = sjson.Set(JsonText, "feedCard.links.-1", map[string]interface{}{"messageURL": this.Path4URL, "title": fmt.Sprintf("响应文本: %s", this.ErrorMsg)})
	JsonText, _ = sjson.Set(JsonText, "feedCard.links.-1", map[string]interface{}{"messageURL": this.Path4URL, "title": fmt.Sprintf("响应时间: %s", time.Now().Format("2006-01-02 15:04:05"))})

	return JsonText
}

func (this TSYSSTATUS) ToFeedCard2() string {

	var JsonByte []byte

	JsonByte, eror := json.Marshal(this)
	if eror != nil {
		return ""
	}

	return string(JsonByte)
}

func (this *TSYSSTATUS) ToErrorLog(TokenUsr TTokenUsr, SrvConfig TSrvConfig, BotConfig TBotConfig, LoggerFx *logrus.Logger) error {

	FilePath, eror := os.Executable()
	if eror != nil {
		return eror
	}

	TokenUsr.WhoBuild = IfThenI(TokenUsr.WhoBuild == 0, 1000000000, TokenUsr.WhoBuild)

	FilePath, _ = filepath.Split(FilePath)
	eror = os.MkdirAll(FilePath+fmt.Sprintf("/%s/%d", SrvConfig.Log4goAt, TokenUsr.WhoBuild), os.ModePerm)
	if eror != nil {
		return eror
	}

	RootPath := fmt.Sprintf("%s%s/%d/", FilePath, SrvConfig.Log4goAt, TokenUsr.WhoBuild)
	RealPath := fmt.Sprintf("%s%s.log", RootPath, this.EntityID)

	this.DiskPath = RealPath
	this.Path4Log = RealPath
	this.Path4URL = fmt.Sprintf("%s/%s/%d/%s.log", SrvConfig.CloudURL, SrvConfig.Log4goAt, TokenUsr.WhoBuild, this.EntityID)
	//@this.BkColour = fmt.Sprintf("%s/%s/%d/bkcolour.jpg", SrvConfig.CloudURL, SrvConfig.UploadAt, TokenUsr.WhoBuild)

	switch this.DataFrom {
	case CONST_DATA_FROM_TO_桌面:
		{
			this.FromText = "桌面终端"
		}
	case CONST_DATA_FROM_TO_钉钉:
		{
			this.FromText = "移动终端"
		}
	case CONST_DATA_FROM_TO_微信:
		{
			this.FromText = "微信终端"
		}
	case CONST_DATA_FROM_TO_安卓:
		{
			this.FromText = "安卓终端"
		}
	case CONST_DATA_FROM_TO_苹果:
		{
			this.FromText = "苹果终端"
		}
	}

	//#保存日志
	eror = ToFile(RealPath, this)
	if eror != nil {
		return eror
	}

	//#保存请求
	RealPath = strings.ReplaceAll(RealPath, ".log", ".req")
	this.Path4Req = RealPath
	eror = os.WriteFile(RealPath, []byte(this.RqSource), CONST_CORE_PATH_TO_文件权限)
	if eror != nil {
		return eror
	}

	return nil
}
