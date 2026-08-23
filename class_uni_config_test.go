package KzPack4Go

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestUrlBuildMySQL(t *testing.T) {
	uni := TUniConfig{
		UnicType: "MySQL",
		UnicSrvr: "127.0.0.1",
		UnicPort: "3306",
		UnicUser: "root",
		UnicPswd: "secret",
		DataBase: "demo",
	}

	dsn := uni.UrlBuild()
	if !strings.HasPrefix(dsn, "root:secret@tcp(127.0.0.1:3306)/demo") {
		t.Errorf("MySQL DSN 前缀异常: %s", dsn)
	}
	// parseTime=true:让 database/sql 自动解析 DATETIME 为 time.Time
	if !strings.Contains(dsn, "parseTime=true") {
		t.Errorf("MySQL DSN 缺少 parseTime=true: %s", dsn)
	}
}

func TestUrlBuildPostgreSQL(t *testing.T) {
	uni := TUniConfig{
		UnicType: "PostgreSQL",
		UnicSrvr: "127.0.0.1",
		UnicPort: "5432",
		UnicUser: "root",
		UnicPswd: "secret",
		DataBase: "demo",
	}

	dsn := uni.UrlBuild()
	if !strings.Contains(dsn, "host=127.0.0.1 port=5432 user=root password=secret dbname=demo") {
		t.Errorf("PostgreSQL DSN 异常: %s", dsn)
	}
}

// ---------- 其余数据库分支 ----------

func TestUrlBuildSQLServer(t *testing.T) {
	dsn := (TUniConfig{UnicType: "SQL Server", UnicSrvr: "h", UnicPort: "1433", UnicUser: "sa", UnicPswd: "p", DataBase: "db"}).UrlBuild()
	if !strings.Contains(dsn, "server=h;port=1433;user id=sa;password=p;database=db;encrypt=disable") {
		t.Errorf("SQL Server DSN 异常: %s", dsn)
	}
}

func TestUrlBuildOracle(t *testing.T) {
	uni := TUniConfig{UnicType: "Oracle", UnicSrvr: "h", UnicPort: "1521", UnicUser: "u", UnicPswd: "p"}
	dsn := uni.UrlBuild()
	// DataBase 缺省应回退 orcl
	if dsn != "u/p@h:1521/orcl" {
		t.Errorf("Oracle DSN 异常: %s", dsn)
	}
}

func TestUrlBuildSQLite(t *testing.T) {
	dsn := (TUniConfig{UnicType: "SQLite", DataBase: "demo.db"}).UrlBuild()
	if dsn != "demo.db" {
		t.Errorf("SQLite DSN 异常: %s", dsn)
	}
}

func TestUrlBuildDM(t *testing.T) {
	dsn := (TUniConfig{UnicType: "dm", UnicSrvr: "h", UnicPort: "5236", UnicUser: "u", UnicPswd: "p", DataBase: "db"}).UrlBuild()
	if dsn != "dm://u:p@h:5236?schema=db" {
		t.Errorf("dm DSN 异常: %s", dsn)
	}
}

func TestUrlBuildOpenGauss(t *testing.T) {
	dsn := (TUniConfig{UnicType: "opengauss", UnicSrvr: "h", UnicPort: "5432", UnicUser: "u", UnicPswd: "p", DataBase: "db"}).UrlBuild()
	if !strings.Contains(dsn, "target_session_attrs=primary") || !strings.Contains(dsn, "dbname=db") {
		t.Errorf("opengauss DSN 异常: %s", dsn)
	}
}

func TestAdaptive(t *testing.T) {
	valid := []string{"PostgreSQL", "PolorDB", "SQL Server", "Oracle", "MySQL", "SQLite", "dm", "opengauss"}
	for _, typ := range valid {
		if err := (TUniConfig{UnicType: typ}).Adaptive(); err != nil {
			t.Errorf("Adaptive(%s) 应通过: %v", typ, err)
		}
	}
	if err := (TUniConfig{UnicType: "Unknown"}).Adaptive(); err == nil {
		t.Error("Adaptive(Unknown) 应报错")
	}
}

func TestUniConfigClearAttr(t *testing.T) {
	uni := TUniConfig{UnicSrvr: "h", UnicPswd: "p", DataBase: "db", UnicPort: "5432", UnicUser: "u"}
	uni.ClearAttr()
	if uni.UnicSrvr != "" || uni.UnicPswd != "" || uni.DataBase != "" || uni.UnicPort != "" || uni.UnicUser != "" {
		t.Errorf("ClearAttr 后未清空: %+v", uni)
	}
}

func TestUniConfigMarshal(t *testing.T) {
	uni := TUniConfig{WhoBuild: 1, UnicType: "MySQL", UnicUser: "u"}
	data, err := json.Marshal(uni)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}
	for _, k := range []string{"cryptTxt", "whoBuild", "unicIndx", "unicType", "unicSrvr", "unicUser", "unicPswd", "unicPort", "dataBase", "instance", "dbConfig"} {
		if _, ok := m[k]; !ok {
			t.Errorf("UniConfig JSON 缺少键 %s", k)
		}
	}
}

// 凭据输出契约:unicPswd 会被序列化(当前行为,README 要求输出前先 ClearAttr);
// ClearAttr 后凭据/连接字段必须为空。
func TestUniConfigCredentialSanitize(t *testing.T) {
	uni := TUniConfig{UnicPswd: "secret", UnicUser: "u", UnicSrvr: "h", UnicPort: "5432", DataBase: "db"}

	data, err := json.Marshal(uni)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "\"unicPswd\":\"secret\"") {
		t.Error("MarshalJSON 应包含 unicPswd(当前行为,输出前须 ClearAttr)")
	}

	uni.ClearAttr()
	data, err = json.Marshal(uni)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}
	if m["unicPswd"] != "" || m["unicUser"] != "" || m["unicSrvr"] != "" || m["unicPort"] != "" || m["dataBase"] != "" {
		t.Errorf("ClearAttr 后凭据/连接字段应为空: %v", m)
	}
}
