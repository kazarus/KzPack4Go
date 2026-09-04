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
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/yitter/idgenerator-go/idgen"
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

// 根据不同哈希类型进行哈希:md5、sha1、sha256、sha512
func ToHash(HashText string, HashType string, isHex bool) (string, error) {

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
		// md4、ripemd160 等类型未实现,返回错误
		return "", fmt.Errorf("unsupported hash type: %s", HashType)
	}
	if isHex {
		arr, eror := hex.DecodeString(HashText) // 十六进制字符串转为十六进制字节数组
		if eror != nil {
			return "", fmt.Errorf("hex decode error: %v", eror)
		}
		hashInstance.Write(arr) // 写入哈希实例对象
	} else {
		hashInstance.Write([]byte(HashText)) // 将字符串转换为字节数组，写入哈希对象
	}
	bytes := hashInstance.Sum(nil) // 哈希值追加到参数后面，只获取原始值，不用追加，用nil，返回哈希值字节数组

	return fmt.Sprintf("%x", bytes), nil // 格式化输出哈希值
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

func NextKJQJ(Value int64) int64 {

	var KJND int64
	var KJQJ int64

	KJND = int64(Value) / 100
	KJQJ = int64(Value) % 100
	KJQJ = KJQJ + 1

	if KJQJ == 13 {
		KJND = KJND + 1
		KJQJ = 1
	}

	return KJND*100 + KJQJ
}

func SubStr(value string, start int, ended int) string {

	rs := []rune(value)
	length := len(rs)

	if start < 0 || start > length {
		return ""
	}

	if ended < 0 || ended > length {
		return ""
	}

	// #ended 小于 start 时,直接切片会触发 slice bounds out of range panic,提前拦截
	if ended < start {
		return ""
	}
	return string(rs[start:ended])
}

func TimePlus(StartDay int) float64 {

	result, _ := strconv.ParseFloat(time.Now().AddDate(0, 0, StartDay).Format("20060102150405"), 64)
	return result
}

func DatePlus(StartDay int) int64 {

	result, _ := strconv.ParseInt(time.Now().AddDate(0, 0, StartDay).Format("20060102"), 10, 64)
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
	NoNumber = strings.TrimSuffix(NoNumber, ".")

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

	switch v := value.(type) {
	case float64:
		key = strconv.FormatFloat(v, 'f', -1, 64)
	case float32:
		key = strconv.FormatFloat(float64(v), 'f', -1, 64)
	case int:
		key = strconv.Itoa(v)
	case uint:
		key = strconv.Itoa(int(v))
	case int8:
		key = strconv.Itoa(int(v))
	case uint8:
		key = strconv.Itoa(int(v))
	case int16:
		key = strconv.Itoa(int(v))
	case uint16:
		key = strconv.Itoa(int(v))
	case int32:
		key = strconv.Itoa(int(v))
	case uint32:
		key = strconv.Itoa(int(v))
	case int64:
		key = strconv.FormatInt(v, 10)
	case uint64:
		key = strconv.FormatUint(v, 10)
	case string:
		key = v
	case []byte:
		key = string(v)
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}

// #ToEscape interface 转字符串并对引号/反斜杠/控制字符做 Go 转义(去掉外层引号);
// #与 AnyToVal 共享转换逻辑
func ToEscape(value interface{}) string {

	// #strconv.Quote 保证首尾字符必为外层引号、内部引号必被转义,精确剥掉首尾各一个字符;
	// #不能使用 strings.Trim(...,"\"")——它会连转义引号 \" 的尾部引号一起剥掉
	var quoted = strconv.Quote(AnyToVal(value))

	return quoted[1 : len(quoted)-1]
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

func Contain(Val string, list []string) bool {

	for _, One := range list {
		if One == Val {
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

func GetEnvDef(Key string, Value string) string {

	var result string
	result = os.Getenv(Value)
	if result == "" {
		return Value
	}

	return result
}

// #判断是否文件
func IsFile(FilePath string) bool {

	fi, e := os.Stat(FilePath)
	if e != nil {
		return false
	}

	return !fi.IsDir()
}

// #配置存档
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

// #从文件初始化配置
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

// #判断文件存在
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

// #格式化文件大小,更加易读
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

// #生成长度16位 未超 9007199254740992(2的53次方)
func ToNext() int64 {
	return idgen.NextId()
}

// #应用所在目录
func AtRoot() string {

	FilePath, eror := os.Executable()
	if eror != nil {
		return ""
	}

	FilePath, _ = filepath.Split(FilePath)

	return FilePath
}

// #FatalWhenError 检查错误:出错时直接退出进程(log.Fatalln / os.Exit)。
// #⚠️ 危险:会终止整个应用且无法被 recover 拦截,仅建议用于启动期的必要性检查
// #(如配置加载、核心依赖初始化失败)。严禁在业务逻辑、HTTP 处理或并发任务中使用;
// #业务错误请自行记录日志并返回,由调用方决定处理方式。
func FatalWhenError(err error, msg string) {

	if err != nil {
		log.Fatalln(msg, err)
	}
}

// #初始化ID生成器;WorkerID须在[0,63]区间(默认WorkerIdBitLength=6)
func NewIdGenerator(WorkerID int64) {

	var options = idgen.NewIdGeneratorOptions(uint16(WorkerID))
	idgen.SetIdGenerator(options)
}
