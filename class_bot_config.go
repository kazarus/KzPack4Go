package KzPack4Go

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"gopkg.in/resty.v1"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// #错误消息
type TBotConfig struct {
	Ding4Bot int64  `json:"ding4Bot"` //#钉钉平台
	SysCustc int64  `json:"sysCustc"` //#系统客户
	AsReally int64  `json:"asReally"` //#是否生产
	Name4App string `json:"name4App"` //#服务名称
	EnvrName string `json:"envrName"` //#日志环境
	EndPoint string `json:"endPoint"` //#访问网关
	AccessID string `json:"accessID"` //#访问令牌
	SecretID string `json:"secretID"` //#签名密钥
	UploadD2 int64  `json:"uploadD2"` //#上传附件
}

func (this TBotConfig) ToErrorMsg1(ErrorMsg string) {

	if this.EndPoint == "" {
		return
	}

	if ErrorMsg == "" {
		return
	}

	switch this.Ding4Bot {
	case 0:
		{
			var ListText = make([]string, 0)
			ListText = strings.Split(this.EndPoint, ",")

			ClientGo := resty.New()
			for _, ItemPara := range ListText {

				Request1 := ClientGo.R()
				Request1.SetFormData(map[string]string{"params": string(ErrorMsg)})

				TargetTo := fmt.Sprintf("%s/core/tsysstatus/syncfrom?whobuild=3302260000&userindx=1001", ItemPara)

				_, eror := Request1.Post(TargetTo)
				if eror != nil {
					return
				}
			}
		}
	case 1:
		{
			TimeUnix := time.Now().UnixMilli()

			SignText := fmt.Sprintf("%d\n%s", TimeUnix, this.SecretID)

			Hash := hmac.New(sha256.New, []byte(this.SecretID))
			Hash.Write([]byte(SignText))
			ByteData := Hash.Sum(nil)
			SignData := url.QueryEscape(base64.StdEncoding.EncodeToString(ByteData))

			sRequest := fmt.Sprintf("%s?access_token=%s&timestamp=%d&sign=%s", this.EndPoint, this.AccessID, TimeUnix, SignData)

			client := http.Client{}

			requestBody := strings.NewReader(ErrorMsg)
			response, err := client.Post(sRequest, "application/json", requestBody)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			defer response.Body.Close()

			if response.StatusCode == http.StatusOK {
				_, err := io.ReadAll(response.Body)
				if err != nil {
					fmt.Println(err.Error())
				}
			}
		}
	}
}

func (this TBotConfig) ToErrorMsg2(ErrorMsg string) {

	var eror error
	var SysStatus TSYSSTATUS
	eror = json.Unmarshal([]byte(ErrorMsg), &SysStatus)
	if eror != nil {
		return
	}

	if this.EndPoint == "" {
		return
	}

	var TargetTo string
	var FormPara string

	switch this.Ding4Bot {
	case 0:
		{
			var ListText = make([]string, 0)
			ListText = strings.Split(this.EndPoint, ",")

			ClientGo := resty.New()
			for _, ItemPara := range ListText {

				Request1 := ClientGo.R()

				switch this.UploadD2 {
				case -1:
					{
						//#标记-1时, 不上传附件
					}
				default:
					{
						//#上传附件-<
						TargetTo = fmt.Sprintf("%s/core/tpath/uploadd2?whobuild=3302260000&userindx=1001", ItemPara)
						FormPara = "upload[]"

						file, eror := os.Stat(SysStatus.Path4Req)
						if eror != nil {
							fmt.Println(eror.Error())
						} else if file.Size() > 0 {

							FileByte, _ := os.ReadFile(SysStatus.Path4Req)
							Request1.SetFileReader(FormPara, fmt.Sprintf("%s.req", SysStatus.EntityID), bytes.NewReader(FileByte))
							resp, eror := Request1.Post(TargetTo)
							if eror != nil {
								return
							}
							//@fmt.Println(string(resp.Body()))

							JsonText := resp.String()
							OnStatus := gjson.Get(JsonText, "ONSTATUS").Bool()
							//@ErorMemo := gjson.Get(JsonText, "ERORMEMO").String()
							if OnStatus {
								//@JsonText = gjson.Get(resp.String(), "LISTDATA").String()
								SysStatus.RqStream = gjson.Get(JsonText, "LISTDATA.0.sysPathc").Int() //#取出附件标识
							}
						}
						//#上传附件->
					}
				}

				ErrorMsg = SysStatus.ToFeedCard2()
				Request1.SetFormData(map[string]string{"params": string(ErrorMsg)})

				TargetTo = fmt.Sprintf("%s/core/tsysstatus/syncfrom?whobuild=3302260000&userindx=1001", ItemPara)

				_, eror = Request1.Post(TargetTo)
				if eror != nil {
					return
				}
			}
		}
	case 1:
		{
			TimeUnix := time.Now().UnixMilli()

			SignText := fmt.Sprintf("%d\n%s", TimeUnix, this.SecretID)

			Hash := hmac.New(sha256.New, []byte(this.SecretID))
			Hash.Write([]byte(SignText))
			ByteData := Hash.Sum(nil)
			SignData := url.QueryEscape(base64.StdEncoding.EncodeToString(ByteData))

			sRequest := fmt.Sprintf("%s?access_token=%s&timestamp=%d&sign=%s", this.EndPoint, this.AccessID, TimeUnix, SignData)

			client := http.Client{}

			requestBody := strings.NewReader(ErrorMsg)
			response, err := client.Post(sRequest, "application/json", requestBody)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			defer response.Body.Close()

			if response.StatusCode == http.StatusOK {
				_, err := io.ReadAll(response.Body)
				if err != nil {
					fmt.Println(err.Error())
				}
			}
		}
	}
}

//#暂不启用
/*
func (this TBotConfig) ToErrorLog(ErrorMsg string, LoggerFx *logrus.Logger) {

	TimeUnix := time.Now().UnixMilli()

	SignText := fmt.Sprintf("%d\n%s", TimeUnix, this.SecretID)

	Hash := hmac.New(sha256.New, []byte(this.SecretID))
	Hash.Write([]byte(SignText))
	ByteData := Hash.Sum(nil)
	SignData := url.QueryEscape(base64.StdEncoding.EncodeToString(ByteData))

	sRequest := fmt.Sprintf("%s?access_token=%s&timestamp=%d&sign=%s", this.EndPoint, this.AccessID, TimeUnix, SignData)

	client := http.Client{}

	requestBody := strings.NewReader(ErrorMsg)
	response, err := client.Post(sRequest, "application/json", requestBody)
	if err != nil {
		LoggerFx.Error(err.Error())
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		_, err := ioutil.ReadAll(response.Body)
		if err != nil {
			LoggerFx.Error(err.Error())
		}
		//#fmt.Println(string(body))
	}
}
*/
