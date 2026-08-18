package KzPack4Go

import (
	"encoding/json"
	"errors"
	"fmt"
)

// # Copy From DbaUniEngine
const CONST_PRIVIDER_NAME_TO_POSTGR = "PostgreSQL"
const CONST_PRIVIDER_NAME_TO_SQLSRV = "SQL Server"
const CONST_PRIVIDER_NAME_TO_ORACLE = "Oracle"
const CONST_PRIVIDER_NAME_TO_ACCESS = "Access"
const CONST_PRIVIDER_NAME_TO_SQLITE = "SQLite"
const CONST_PRIVIDER_NAME_TO_MYSQLN = "MySQL"
const CONST_PRIVIDER_NAME_TO_KINGES = "kb"
const CONST_PRIVIDER_NAME_TO_DAMENG = "dm"
const CONST_PRIVIDER_NAME_TO_OPENGS = "opengauss"
const CONST_PRIVIDER_NAME_TO_POLODB = "PolorDB"

// #数据配置
type TUniConfig struct {
	CryptTxt string
	WhoBuild int64
	UnicIndx int64
	UnicType string //#数据库驱动
	UnicSrvr string //#数据库地址
	UnicUser string //#数据库用户
	UnicPswd string //#数据库密码
	UnicPort string //#数据库端口
	DataBase string //#数据库名称
	Instance string //#数据库实例
	DbConfig string //#数据库链接
}

func (this TUniConfig) MarshalJSON() ([]byte, error) {
	var result = make(map[string]interface{})

	result["cryptTxt"] = this.CryptTxt
	result["whoBuild"] = this.WhoBuild
	result["unicIndx"] = this.UnicIndx

	result["unicType"] = this.UnicType
	result["unicSrvr"] = this.UnicSrvr
	result["unicUser"] = this.UnicUser
	result["unicPswd"] = this.UnicPswd
	result["unicPort"] = this.UnicPort
	result["dataBase"] = this.DataBase
	result["instance"] = this.Instance
	result["dbConfig"] = this.DbConfig

	return json.Marshal(result)
}

func (this TUniConfig) Adaptive() error {

	if (this.UnicType != "PostgreSQL") && (this.UnicType != "SQL Server") && (this.UnicType != "Oracle") && (this.UnicType != "MySQL") && (this.UnicType != "dm") && (this.UnicType != "opengauss") && (this.UnicType != "PolorDB") {
		return errors.New("error: unsupport database; value from config.uni.json->unictype should be 'PostgreSQL' or 'SQL Server' or 'Oracle' or 'MySQL' or 'dm' or 'opengauss' ")
	}

	return nil
}

func (this TUniConfig) UrlBuild() string {

	var DbSource string

	switch this.UnicType {
	case "PostgreSQL", "PolorDB":
		{
			DbSource = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable application_name=publish10001", this.UnicSrvr, this.UnicPort, this.UnicUser, this.UnicPswd, this.DataBase)
		}
	case "SQL Server":
		{
			DbSource = fmt.Sprintf("server=%s;port=%s;user id=%s;password=%s;database=%s;encrypt=disable;Pooling=True;Max Pool Size=1000", this.UnicSrvr, this.UnicPort, this.UnicUser, this.UnicPswd, this.DataBase) //#
		}
	case "Oracle":
		{
			this.DataBase = IfThenS(this.DataBase == "", "orcl", this.DataBase)
			this.Instance = this.DataBase
			DbSource = fmt.Sprintf("%s/%s@%s:%s/%s", this.UnicUser, this.UnicPswd, this.UnicSrvr, this.UnicPort, this.DataBase)
		}
	case "dm":
		{
			//#`"dm://user:pswd@srvr:port?appName=yzcloud&connectTimeout=30000&schema=database"`
			DbSource = fmt.Sprintf("dm://%s:%s@%s:%s?schema=%s", this.UnicUser, this.UnicPswd, this.UnicSrvr, this.UnicPort, this.DataBase)
		}
	case "MySQL":
		{
			//#DbSource = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local", UniConfig.UnicUser, UniConfig.UnicPswd, UniConfig.UnicSrvr, UniConfig.UnicPort, UniConfig.DataBase)
			DbSource = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", this.UnicUser, this.UnicPswd, this.UnicSrvr, this.UnicPort, this.DataBase)
		}
	case "SQLite":
		{
			DbSource = fmt.Sprintf("%s", this.DataBase)
		}
	case "opengauss":
		{
			DbSource = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable target_session_attrs=primary connect_timeout=1", this.UnicSrvr, this.UnicPort, this.UnicUser, this.UnicPswd, this.DataBase)
		}
	}

	return DbSource
}

func (this *TUniConfig) ClearAttr() {

	this.UnicSrvr = ""
	this.UnicPswd = ""
	this.DataBase = ""
	this.UnicPort = ""
	this.UnicUser = ""
}
