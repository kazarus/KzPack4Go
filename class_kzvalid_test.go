package KzPack4Go

import (
	"encoding/json"
	"testing"
)

func TestTKzValidMarshal(t *testing.T) {
	v := TKzValid{
		ValidVal: true,
		ValidInt: 1,
		ValidMsg: "ok",
		FlowStat: 1,
		FlowMain: 2,
		FlowNode: 3,
		EntityID: "e1",
		UserUUID: "u1",
		UserName: "张三",
	}

	data, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	want := map[string]interface{}{
		"validVal": true,
		"validInt": float64(1),
		"validMsg": "ok",
		"flowStat": float64(1),
		"flowMain": float64(2),
		"flowNode": float64(3),
		"entityID": "e1",
		"userUUID": "u1",
		"userName": "张三",
	}
	for k, v := range want {
		if m[k] != v {
			t.Errorf("TKzValid JSON 键 %s = %v; want %v", k, m[k], v)
		}
	}
	if _, ok := m["metaData"]; !ok {
		t.Error("TKzValid JSON 应包含 metaData(当前行为)")
	}
}
