package KzPack4Go

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"log"
	"path/filepath"

	"github.com/bwmarrin/snowflake"
	"github.com/idgenerator-go/idgen"
	"github.com/kardianos/osext"
	//"github.com/satori/go.uuid"
	"github.com/gofrs/uuid/v5"
)

var eror error

type TErrorMsg struct {
	RowIndex int64  `json:"rowIndex"`
	ErrorKey string `json:"errorKey"`
	ErrorMsg string `json:"errorMsg"`
	ErrorCnt int64  `json:"errorCnt"` //#错误次数
}

func (this TErrorMsg) MarshalJSON() ([]byte, error) {
	var result = make(map[string]interface{})

	result["rowIndex"] = this.RowIndex
	result["errorKey"] = this.ErrorKey
	result["errorMsg"] = this.ErrorMsg
	result["errorCnt"] = this.ErrorCnt

	result["metaData"] = "实体类名: IErrorMsgResponse / eror.TErrorMsg"

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	return json.Marshal(result)
}

type TEROR struct {
	OnStatus bool
	ErorCode string
	ErorMemo string
	DataSize int
	ErrLevel int
	DataDate int64
	DataTime float64

	RealPage int64
	PageSize int64
	TotalCnt int64

	ErrorCnt int64
	ErrorKey string
	CryptTxt string

	ListData interface{}
}

func (this TEROR) MarshalJSON() ([]byte, error) {

	var result = make(map[string]interface{})

	result["ONSTATUS"] = this.OnStatus
	result["ERORCODE"] = this.ErorCode
	result["ERORMEMO"] = this.ErorMemo
	result["DATASIZE"] = this.DataSize
	result["ERRLEVEL"] = this.ErrLevel
	result["DATADATE"] = this.DataDate
	result["DATATIME"] = this.DataTime

	result["REALPAGE"] = this.RealPage
	result["PAGESIZE"] = this.PageSize
	result["TOTALCNT"] = this.TotalCnt

	result["ERRORCNT"] = this.ErrorCnt
	result["ERRORKEY"] = this.ErrorKey
	result["CRYPTTXT"] = this.CryptTxt

	if this.ListData == nil {
		this.ListData = make([]interface{}, 0)
	}
	result["LISTDATA"] = this.ListData

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	return json.Marshal(result)
}

func ToEROR(ErorMemo string) TEROR {

	var result TEROR

	result.DataSize = 0
	result.OnStatus = false
	result.ErorCode = "EROR"
	result.ErorMemo = ErorMemo

	result.DataDate = ToDate()
	result.DataTime = ToTime()

	return result
}

func ToMEMO(ErorMemo string) TEROR {

	var result TEROR

	result.DataSize = 0
	result.OnStatus = true
	result.ErorCode = "TRUE"
	result.ErorMemo = ErorMemo

	result.DataDate = ToDate()
	result.DataTime = ToTime()

	return result
}

func ToTRUE(ListData interface{}) TEROR {

	var result TEROR

	result.DataSize = 0
	result.OnStatus = true
	result.ErorCode = "TRUE"
	result.ErorMemo = "TRUE"
	result.ErrLevel = 0

	result.DataDate = ToDate()
	result.DataTime = ToTime()

	result.ListData = ListData

	return result
}

func ToSIZE(ListData interface{}, DataSize int) TEROR {

	var result TEROR

	result.OnStatus = true
	result.ErorCode = "TRUE"
	result.ErorMemo = "TRUE"
	result.ErrLevel = 0

	result.DataDate = ToDate()
	result.DataTime = ToTime()

	result.ListData = ListData
	result.DataSize = DataSize

	return result
}

func ToPage(ListData interface{}, DataSize int, RealPage int64, PageSize int64, TotalCnt int64) TEROR {

	var result TEROR

	result.OnStatus = true
	result.ErorCode = "TRUE"
	result.ErorMemo = "TRUE"
	result.ErrLevel = 0

	result.DataDate = ToDate()
	result.DataTime = ToTime()

	result.ListData = ListData
	result.DataSize = DataSize

	result.RealPage = RealPage
	result.PageSize = PageSize
	result.TotalCnt = TotalCnt

	return result
}

// #允许指定所有字段
func ToDATA(OnStatus bool, ErorCode string, ErorMemo string, DataSize int, ListData interface{}) TEROR {

	var result TEROR

	result.OnStatus = OnStatus
	result.ErorCode = ErorCode
	result.ErorMemo = ErorMemo

	result.DataDate = ToDate()
	result.DataTime = ToTime()

	result.DataSize = DataSize
	result.ListData = ListData

	return result
}

// #允许指定消息级别,ErrLevel<0时, 不予消息提示
func ToErrorKey(ErorMemo string, ErrorKey string) TEROR {

	var result TEROR

	result.DataSize = 0
	result.OnStatus = false
	result.ErorCode = "EROR"
	result.ErorMemo = ErorMemo
	result.ErrLevel = 0
	result.ErrorKey = ErrorKey

	result.DataDate = ToDate()
	result.DataTime = ToTime()

	return result
}

// #允许指定消息级别,ErrLevel<0时, 携带出错次数
func ToErrorCnt(ErorMemo string, ErrorKey string, ErrorCnt int64) TEROR {

	var result TEROR

	result.DataSize = 0
	result.OnStatus = false
	result.ErorCode = "EROR"
	result.ErorMemo = ErorMemo
	result.ErrLevel = 0
	result.ErrorKey = ErrorKey

	result.DataDate = ToDate()
	result.DataTime = ToTime()

	return result
}

// #允许指定消息级别,ErrLevel<0时, 不予消息提示
func ToErrLevel(ErorMemo string, ErrLevel int, ErrorCnt int64, ErrorKey string) TEROR {

	var result TEROR

	result.DataSize = 0
	result.OnStatus = false
	result.ErorCode = "EROR"
	result.ErorMemo = ErorMemo
	result.ErrLevel = ErrLevel
	result.ErrorCnt = ErrorCnt
	result.ErrorKey = ErrorKey

	result.DataDate = ToDate()
	result.DataTime = ToTime()

	return result
}

// #允许指定所有字段
func ToCryptTxt(OnStatus bool, ErorCode string, ErorMemo string, DataSize int, CryptTxt string, ListData interface{}) TEROR {

	var result TEROR

	result.OnStatus = OnStatus
	result.ErorCode = ErorCode
	result.ErorMemo = ErorMemo

	result.DataDate = ToDate()
	result.DataTime = ToTime()

	result.DataSize = DataSize
	result.ListData = ListData
	result.CryptTxt = CryptTxt

	return result
}

// #生成唯一标识,多用于EntityID
func ToUUID() string {

	/*
		uuid := uuid.NewV4()
		return fmt.Sprintf("{%s}", uuid)
	*/

	uuid, _ := uuid.NewV7()
	return fmt.Sprintf("{%s}", uuid)
}

// #生成唯一标识的简写形式,便于对比
func ToLite(EntityID string) string {

	if EntityID == "" {
		return ""
	}
	return fmt.Sprintf("%s...%s", EntityID[1:5], EntityID[len(EntityID)-5:len(EntityID)-1])
}

// #生成长度19位 超出 9007199254740992(2的53次方)
func ToNext19() int64 {

	var Snowflake *snowflake.Node
	Snowflake, eror = snowflake.NewNode(1)
	if eror != nil {
		return 0
	}

	return Snowflake.Generate().Int64()
}

// #生成长度16位 未超 9007199254740992(2的53次方)
func ToNext() int64 {

	var result int64
	result = idgen.NextId()
	//if result%2 == 1 {
	//	result = idgen.NextId()
	//}
	//if result%2 == 1 {
	//	result = idgen.NextId()
	//}

	return result
}

func AtRoot() string {

	FilePath, eror := osext.Executable()
	if eror != nil {
		return ""
	}

	FilePath, _ = filepath.Split(FilePath)

	return FilePath
}

func ChEror(err error, msg string) {

	if err != nil {
		log.Fatalln(msg, err)
	}
}
