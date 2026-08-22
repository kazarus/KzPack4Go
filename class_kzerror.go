package KzPack4Go

import (
	"encoding/json"
)

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
	result.ErrorCnt = ErrorCnt

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
