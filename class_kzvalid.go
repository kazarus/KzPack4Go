package KzPack4Go

import (
	"encoding/json"
)

type TKzValid struct {
	ValidVal bool   `json:"validVal"` //#校验结果
	ValidInt int64  `json:"validInt"` //#校验结果
	ValidMsg string `json:"validMsg"` //#校验结果

	FlowStat int64 `json:"flowStat"` //#流程状态
	FlowMain int64 `json:"flowMain"` //#审核流程
	FlowNode int64 `json:"flowNode"` //#审核节点

	EntityID string `json:"entityID"` //#单据标识
	UserUUID string `json:"userUUID"` //#操作用户
	UserName string `json:"userName"` //#操作用户
}

func (this TKzValid) MarshalJSON() ([]byte, error) {

	var result = make(map[string]interface{})

	result["validVal"] = this.ValidVal
	result["validInt"] = this.ValidInt
	result["validMsg"] = this.ValidMsg

	result["flowStat"] = this.FlowStat
	result["flowMain"] = this.FlowMain
	result["flowNode"] = this.FlowNode

	result["entityID"] = this.EntityID
	result["userUUID"] = this.UserUUID
	result["userName"] = this.UserName

	result["metaData"] = "实体类名:TValidValResponse / core.TValidVal; validVal=校验结果;flowStat=流程状态(0:待审,1:通过,-1:退回);validInt与flowStat同值;"

	return json.Marshal(result)
}
