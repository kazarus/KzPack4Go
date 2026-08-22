package KzPack4Go

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

import (
	"github.com/bwmarrin/snowflake"
	"github.com/yitter/idgenerator-go/idgen"
	//"github.com/satori/go.uuid"
	"github.com/gofrs/uuid/v5"
)

const CONST_CORE_PATH_TO_目录权限 = 0750
const CONST_CORE_PATH_TO_文件权限 = 0640
const CONST_CORE_PATH_TO_日志权限 = 0664

type KzUnicode string

func (s KzUnicode) MarshalJSON() ([]byte, error) {
	return []byte(strconv.QuoteToASCII(string(s))), nil
}

type TRangeVal struct {
	UseRange int64   `json:"useRange"`
	StartVal float64 `json:"startVal"`
	EndedVal float64 `json:"endedVal"`
}

// 根据不同哈希类型进行哈希：md5、sha1、sha256、sha512
func ToHash(HashText string, HashType string, isHex bool) string {

	var hashInstance hash.Hash // 定义哈希实例

	switch HashType { // 选择哈希类型
	case "md5":
		hashInstance = md5.New()
	case "sha1":
		hashInstance = sha1.New()
	case "sha256":
		hashInstance = sha256.New()
	case "sha512":
		hashInstance = sha512.New()
	default:
		// md4、ripemd160 等类型未实现,直接返回错误提示
		return fmt.Sprintf("unsupported hash type: %s", HashType)
	}
	if isHex {
		arr, _ := hex.DecodeString(HashText) // 十六进制字符串转为十六进制字节数组
		hashInstance.Write(arr)              // 写入哈希实例对象
	} else {
		hashInstance.Write([]byte(HashText)) // 将字符串转换为字节数组，写入哈希对象
	}
	bytes := hashInstance.Sum(nil) // 哈希值追加到参数后面，只获取原始值，不用追加，用nil，返回哈希值字节数组

	return fmt.Sprintf("%x", bytes) // 格式化输出哈希值
}

func ToMsec() float64 {

	result, _ := strconv.ParseFloat(time.Now().Format("20060102150405.000"), 64)
	return result
}

func ToTime() float64 {

	result, _ := strconv.ParseFloat(time.Now().Format("20060102150405"), 64)
	return result
}

func ToHour() int64 {

	result, _ := strconv.ParseInt(time.Now().Format("15"), 10, 64)
	return result
}

func ToDate() int64 {

	result, _ := strconv.ParseInt(time.Now().Format("20060102"), 10, 64)
	return result
}

func ToWeek() int64 {

	Year, Week := time.Now().ISOWeek()
	return int64(Year)*100 + int64(Week)
}

func ToKJND() int64 {

	result, _ := strconv.ParseInt(time.Now().Format("2006"), 10, 64)
	return result
}

func ToKJQJ() int64 {

	result, _ := strconv.ParseInt(time.Now().Format("200601"), 10, 64)
	return result
}

func ToCron(Expire int) string {

	result := time.Now().Add(time.Duration(Expire) * time.Second).Format("05 04 15 02 01")
	result = fmt.Sprintf("%s ?", result)
	return result
}

func NextKJQJ(aValue int64) int64 {

	var KJND int64
	var KJQJ int64

	KJND = int64(aValue) / 100
	KJQJ = int64(aValue) % 100
	KJQJ = KJQJ + 1

	if KJQJ == 13 {
		KJND = KJND + 1
		KJQJ = 1
	}

	return KJND*100 + KJQJ
}

func SubStr(str string, start int, ended int) string {

	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		return ""
	}

	if ended < 0 || ended > length {
		return ""
	}
	return string(rs[start:ended])
}

func TimePlus(aStartDay int) float64 {

	result, _ := strconv.ParseFloat(time.Now().AddDate(0, 0, aStartDay).Format("20060102150405"), 64)
	return result
}

func DatePlus(aStartDay int) int64 {

	result, _ := strconv.ParseInt(time.Now().AddDate(0, 0, aStartDay).Format("20060102"), 10, 64)
	return result
}

func TimeToStr(Value float64) string {
	//#20200618105004 -> 2020-06-18 10:50:04
	result := fmt.Sprintf("%.0f", Value)

	if len(result) != 14 {
		return "-1"
	}

	a := result[0:4]
	b := result[4:6]
	c := result[6:8]
	d := result[8:10]
	e := result[10:12]
	f := result[12:14]

	return fmt.Sprintf("%s-%s-%s %s:%s:%s", a, b, c, d, e, f)
}

func DateToStr(Value int64) string {
	//#20200618 -> 2020-06-18
	result := fmt.Sprintf("%d", Value)

	if len(result) != 8 {
		return "-1"
	}

	a := result[0:4]
	b := result[4:6]
	c := result[6:8]

	return fmt.Sprintf("%s-%s-%s", a, b, c)
}

func StrToDate(Value [][]string) int64 {

	if len(Value) != 3 {
		return 0
	}

	a := StrToInt(Value[0][0])
	b := StrToInt(Value[1][0])
	c := StrToInt(Value[2][0])

	return a*10000 + b*100 + c
}

func Format1Time(Value string) string {
	//#20200618 -> 2020-06-18 00:00:00
	result := Value

	if len(result) != 8 {
		return "-1"
	}

	a := result[0:4]
	b := result[4:6]
	c := result[6:8]

	return fmt.Sprintf("%s-%s-%s 00:00:00", a, b, c)
}

func Format2Time(Value string) string {
	//#20200618 -> 2020-06-18 23:59:59
	result := Value

	if len(result) != 8 {
		return "-1"
	}

	a := result[0:4]
	b := result[4:6]
	c := result[6:8]

	return fmt.Sprintf("%s-%s-%s 23:59:59", a, b, c)
}

func Min(x, y int64) int64 {

	if x < y {
		return x
	}

	return y
}

func Max(x, y int64) int64 {

	if x > y {
		return x
	}

	return y
}

func StrToInt(Value string) int64 {

	result, _ := strconv.ParseInt(Value, 10, 64)
	return result
}

func IntToStr(Value int64) string {

	result := strconv.Itoa(int(Value))
	return result
}

func StrToFloat(Value string) float64 {

	result, _ := strconv.ParseFloat(Value, 64)
	return result
}

// #字符串后面没有零
func Float0Text(Value float64) string {

	var NoNumber string
	NoNumber = fmt.Sprintf("%0.4f", Value)
	for strings.HasSuffix(NoNumber, "0") {
		NoNumber = strings.TrimSuffix(NoNumber, "0")
	}
	if strings.HasSuffix(NoNumber, ".") {
		NoNumber = strings.TrimSuffix(NoNumber, ".")
	}

	return NoNumber
}

func Float2Text(Value float64) string {

	result := fmt.Sprintf("%.4f", Value)
	return result
}

func FormatFloat(Value float64, Precision int) float64 {
	if Precision < 0 {
		Precision = 0
	}

	result, _ := strconv.ParseFloat(fmt.Sprintf("%.*f", Precision, Value), 64)
	return result
}

func Format2Float(Value float64) float64 {

	result, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", Value), 64)
	return result
}

func Format4Float(Value float64) float64 {

	result, _ := strconv.ParseFloat(fmt.Sprintf("%.4f", Value), 64)
	return result
}

func Format6Float(Value float64) float64 {

	result, _ := strconv.ParseFloat(fmt.Sprintf("%.6f", Value), 64)
	return result
}

// #interface 转 string to
// #copy from https://blog.csdn.net/xujiamin0022016/article/details/109226605
func AnyToVal(value interface{}) string {

	var key string

	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}

func IfThenI(Value bool, TrueVal, FailVal int64) int64 {

	Result := FailVal
	if Value {
		Result = TrueVal
	}

	return Result
}

func IfThenS(Value bool, TrueVal, FailVal string) string {

	Result := FailVal
	if Value {
		Result = TrueVal
	}

	return Result
}

func IfThenF(Value bool, TrueVal, FailVal float64) float64 {

	Result := FailVal
	if Value {
		Result = TrueVal
	}

	return Result
}

func IfThenB(Value bool, TrueVal, FailVal bool) bool {

	Result := FailVal
	if Value {
		Result = TrueVal
	}

	return Result
}

func Contain(a string, list []string) bool {

	for _, b := range list {
		if b == a {
			return true
		}
	}

	return false
}

func Contains(AllString string, SubString string) bool {

	var ListText = make([]string, 0)
	ListText = strings.Split(AllString, ",")

	for _, ItemText := range ListText {

		if ItemText == SubString {
			return true
		}
	}

	return false
}

func FileExist(FilePath string) (bool, error) {

	_, eror := os.Stat(FilePath)
	if eror == nil {
		return true, nil
	}
	if os.IsNotExist(eror) {
		return false, nil
	}

	return false, eror
}

func IsFile(FilePath string) bool {

	fi, e := os.Stat(FilePath)
	if e != nil {
		return false
	}

	return !fi.IsDir()
}

func ToFile(FilePath string, v interface{}) error {

	var eror error
	var source []byte

	/*
		if !IsFile(aFilePath) {
			return errors.New("Error: The FilePath is not a FileName;")
		}
	*/

	source, eror = json.MarshalIndent(v, "", "\t")
	if eror != nil {
		return eror
	}

	eror = os.WriteFile(FilePath, source, CONST_CORE_PATH_TO_文件权限)
	if eror != nil {
		return eror
	}

	return nil
}

func InFile(FilePath string, v interface{}) error {

	var eror error
	var didExist bool

	/*
		if !IsFile(aFilePath) {
			return errors.New("Error: The FilePath is not a FileName;")
		}
	*/

	didExist, _ = FileExist(FilePath)
	if !didExist {
		eror = ToFile(FilePath, v)
		if eror != nil {
			return eror
		}
		return nil
	}

	source, eror := os.ReadFile(FilePath)
	if eror != nil {
		return eror
	}

	eror = json.Unmarshal(source, &v)
	if eror != nil {
		return eror
	}

	return nil
}

func FormatFileSize(FileSize int64) (size string) {
	if FileSize < 1024 {
		//return strconv.FormatInt(fileSize, 10) + "B"
		return fmt.Sprintf("%.2fB", float64(FileSize)/float64(1))
	} else if FileSize < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", float64(FileSize)/float64(1024))
	} else if FileSize < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", float64(FileSize)/float64(1024*1024))
	} else if FileSize < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(FileSize)/float64(1024*1024*1024))
	} else if FileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fTB", float64(FileSize)/float64(1024*1024*1024*1024))
	} else { //if fileSize < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.2fEB", float64(FileSize)/float64(1024*1024*1024*1024*1024))
	}
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
	if len(EntityID) < 6 {
		return EntityID
	}
	return fmt.Sprintf("%s...%s", EntityID[1:5], EntityID[len(EntityID)-5:len(EntityID)-1])
}

var (
	snowflakeOnce sync.Once
	snowflakeNode *snowflake.Node
	snowflakeEror error
)

// #生成长度19位 超出 9007199254740992(2的53次方)
func ToNext19() int64 {

	snowflakeOnce.Do(func() {
		snowflakeNode, snowflakeEror = snowflake.NewNode(1)
	})
	if snowflakeEror != nil {
		return 0
	}

	return snowflakeNode.Generate().Int64()
}

// #生成长度16位 未超 9007199254740992(2的53次方)
func ToNext() int64 {

	var result int64
	result = idgen.NextId()

	return result
}

func AtRoot() string {

	FilePath, eror := os.Executable()
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

// #初始化ID生成器;WorkerID须在[0,63]区间(默认WorkerIdBitLength=6)
func NewIdGenerator(WorkerID int64) {

	var options = idgen.NewIdGeneratorOptions(uint16(WorkerID))
	idgen.SetIdGenerator(options)
}
