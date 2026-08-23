package KzPack4Go

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

// 这些测试锁定 TSrvConfig / TSYSSTATUS 的 JSON 输出(键集合与值),
// 防止未来增删结构体字段或手写 MarshalJSON 时输出无声变化。
// 注意:json.Marshal 对 map 按键排序、对结构体按字段序,键顺序不同是正常的。

func TestTSrvConfigJSONStable(t *testing.T) {
	cfg := TSrvConfig{
		CryptTxt: "ct", HttpPort: "8080", ListenAt: "8081", SilentAt: "8082", SocketAt: "8083",
		FrontEnd: "3000", JwtToken: 1, LockKJND: 2025, RealKJND: 2026, LockYear: 2025, RealYear: 2026,
		UseRedis: 1, UseError: 1, AllError: 0, UseQueue: 1, ErrorCnt: 5, CloseErr: 1, LockFlow: 1, WhereAmi: 1,
		Path2PDF: "pdf", Path2TXT: "txt",
		RedisTyp: 1, RedisKey: "rk", RedisSrv: "127.0.0.1:6379", RedisPwd: "pwd", RedisUsr: "usr",
		QueueSrv: "amqp://q", QueuePwd: "qpwd", QueueUsr: "qusr",
		PrintURL: "print", EndPoint: "http://ep", CloudURL: "http://cloud", Name4App: "app",
		Ding4URL: "ding", Ding4APP: "dapp", Ding4GOV: "dgov", Ding4Dom: "ddom", Ding4VER: 2, TryAgain: 3,
		Log4goAt: "log4go", UploadAt: "upload", UpFormat: 2,
		SynzFrom: "NTOW", WorkerID: 7,
		ThirdURL: "third", InnerURL: "inner",
		LoginOut: 1, LoginMin: 9,
		SecretOn: 1, ISilence: 0, UsesCron: 1, Delay10s: 10,
		MaxOpenConns: 500, MaxIdleConns: 200, MaxLifetimes: 60, MaxIdleTimes: 10,
		Path4ExtExcludes: ".exe", PathNameExcludes: "evil",
		UseRedisWhenExpd: 1, UpdateOnWhenExpd: 1,
		MakeYearWhenTodo: 1, VoidYearWhenTodo: 0,
		AccessControlAllowOrigin: "*",
	}

	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]interface{}
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}

	want := map[string]interface{}{
		"cryptTxt": "ct", "httpPort": "8080", "listenAt": "8081", "silentAt": "8082", "socketAt": "8083",
		"frontEnd": "3000", "jwtToken": float64(1), "lockKJND": float64(2025), "realKJND": float64(2026),
		"lockYear": float64(2025), "realYear": float64(2026),
		"useRedis": float64(1), "useError": float64(1), "allError": float64(0), "useQueue": float64(1),
		"errorCnt": float64(5), "closeErr": float64(1), "lockFlow": float64(1), "whereAmi": float64(1),
		"path2PDF": "pdf", "path2TXT": "txt",
		"redisTyp": float64(1), "redisKey": "rk", "redisSrv": "127.0.0.1:6379", "redisPwd": "pwd", "redisUsr": "usr",
		"queueSrv": "amqp://q", "queuePwd": "qpwd", "queueUsr": "qusr",
		"printURL": "print", "endPoint": "http://ep", "cloudURL": "http://cloud", "name4App": "app",
		"ding4URL": "ding", "ding4APP": "dapp", "ding4GOV": "dgov", "ding4Dom": "ddom",
		"ding4VER": float64(2), "tryAgain": float64(3),
		"log4goAt": "log4go", "uploadAt": "upload", "upFormat": float64(2),
		"synzFrom": "NTOW", "workerID": float64(7),
		"thirdURL": "third", "innerURL": "inner",
		"loginOut": float64(1), "loginMin": float64(9),
		"secretOn": float64(1), "iSilence": float64(0), "usesCron": float64(1), "delay10s": float64(10),
		"maxOpenConns": float64(500), "maxIdleConns": float64(200), "maxLifetimes": float64(60), "maxIdleTimes": float64(10),
		"path4ExtExcludes": ".exe", "pathNameExcludes": "evil",
		"useRedisWhenExpd": float64(1), "updateOnWhenExpd": float64(1),
		"makeYearWhenTodo": float64(1), "voidYearWhenTodo": float64(0),
		"accessControlAllowOrigin": "*",
	}

	if len(got) != len(want) {
		t.Errorf("TSrvConfig JSON 键数 = %d; want %d; got keys: %v", len(got), len(want), keysOf(got))
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("TSrvConfig JSON 输出与预期不符\n got: %v\nwant: %v", got, want)
	}
}

func TestTSYSSTATUSJSONStable(t *testing.T) {
	st := TSYSSTATUS{
		WhoBuild: 1, UserIndx: 2, StatusID: 3, HashCode: 4, DataLevl: 5, DataFrom: 6, DataType: 7,
		KJNDKJQJ: 8, DataStat: 9, DataOrdr: 10, DataHide: 11, DataDate: 20260822, DataTime: 20260822.5,
		DataGUID: "g", MastName: "m", MastGUID: "mg", CreateBy: "cb", CreateAt: 1, UpdateBy: "ub",
		UpdateAt: 2, DeleteBy: "db", DeleteAt: 3, QueuesBy: 12, EntityID: "e", EnrollID: "en",
		Name4App: "app", HttpPort: "8080", FrontEnd: "3000", EndPoint: "ep", ClientIP: "ip",
		StatusCD: 13, SpendDat: 14, UserUUID: "uu", RqMethod: "POST", Url4Path: "/x", ErrLevel: 0,
		ErrorMsg: "err", ReturnID: "ok", BrowerOS: "chrome", BrowerID: "120",
		DiskPath: "/d", Path4URL: "/p", DataBase: "db",
		BkColour: "c", AsReally: 1, SysCustc: 2, EnvrName: "env", YesIknow: 3, UserMain: 4,
		ActionID: 5, BillUUID: "b", RqStream: 6, Path4Log: "l", Path4Req: "r", RqSource: "s",
		FromText: "f", MsgTitle: "t", MsgState: "st",
	}

	data, err := json.Marshal(st)
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]interface{}
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}

	// 当前 MarshalJSON 手写 map 的完整键集合(59 个,含补上的 browerOS/browerID/metaData)
	wantKeys := []string{
		"whoBuild", "userIndx", "statusID", "hashCode", "dataLevl", "dataFrom", "dataType",
		"kjndKJQJ", "dataStat", "dataOrdr", "dataHide", "dataDate", "dataTime", "dataGUID",
		"mastName", "mastGUID", "createBy", "createAt", "updateBy", "updateAt", "deleteBy",
		"deleteAt", "queuesBy", "entityID", "enrollID", "name4App", "endPoint", "httpPort",
		"frontEnd", "clientIP", "statusCD", "spendDat", "userUUID", "rqMethod", "url4Path",
		"errLevel", "errorMsg", "returnID", "browerOS", "browerID", "diskPath", "path4URL",
		"dataBase", "bkColour", "asReally", "sysCustc", "envrName", "yesIknow", "userMain",
		"actionID", "billUUID", "rqStream", "path4Log", "path4Req", "rqSource", "fromText",
		"msgTitle", "msgState", "metaData",
	}

	if len(got) != len(wantKeys) {
		t.Errorf("TSYSSTATUS JSON 键数 = %d; want %d; got keys: %v", len(got), len(wantKeys), keysOf(got))
	}
	for _, k := range wantKeys {
		if _, ok := got[k]; !ok {
			t.Errorf("TSYSSTATUS JSON 缺少键 %q", k)
		}
	}
	// 抽查值:浏览器字段与 metaData 已纳入输出
	if got["whoBuild"] != float64(1) || got["name4App"] != "app" || got["errorMsg"] != "err" {
		t.Errorf("TSYSSTATUS JSON 值异常: %v", got)
	}
	if got["browerOS"] != "chrome" || got["browerID"] != "120" {
		t.Errorf("TSYSSTATUS JSON 浏览器字段异常: %v", got)
	}
	if got["metaData"] != "实体类名: core.TSysStatus" {
		t.Errorf("TSYSSTATUS JSON metaData = %v", got["metaData"])
	}
}

func keysOf(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// 凭据输出契约:MarshalJSON 会序列化 redisPwd/queuePwd(README 要求对外输出前先 ClearAttr 脱敏);
// ClearAttr 后这些字段必须为空,防止凭据外泄。测试锁定该行为,防止未来增删凭据字段时无声变化。
func TestTSrvConfigCredentialSanitize(t *testing.T) {
	cfg := TSrvConfig{RedisPwd: "rp", QueuePwd: "qp", RedisSrv: "srv", RedisUsr: "usr", QueueSrv: "q", QueueUsr: "qu"}

	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "\"redisPwd\":\"rp\"") || !strings.Contains(string(data), "\"queuePwd\":\"qp\"") {
		t.Error("MarshalJSON 应包含 redisPwd/queuePwd(当前行为,输出前须 ClearAttr)")
	}

	cfg.ClearAttr()
	data, err = json.Marshal(cfg)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}
	for _, k := range []string{"redisPwd", "redisSrv", "redisUsr", "queuePwd", "queueSrv", "queueUsr"} {
		if m[k] != "" {
			t.Errorf("ClearAttr 后凭据字段 %s 应为空: %v", k, m[k])
		}
	}
}
