package KzPack4Go

import "sync"

type TErrorUsr struct {
	UserCode string `json:"userCode"` //#用户编码
	ErrorCnt int64  `json:"errorCnt"` //#错误次数
}

type TBlackList struct {
	BlackList map[string]TErrorUsr
	Lock      sync.Mutex
}

func (this *TBlackList) Append(EnrollID string, ErrorUsr TErrorUsr) {

	if this.BlackList == nil {
		this.BlackList = make(map[string]TErrorUsr, 0)
	}

	this.Lock.Lock()
	defer this.Lock.Unlock()
	this.BlackList[EnrollID] = ErrorUsr
}

func (this *TBlackList) Remove(EnrollID string) {

	if this.BlackList == nil {
		this.BlackList = make(map[string]TErrorUsr, 0)
	}

	this.Lock.Lock()
	defer this.Lock.Unlock()
	delete(this.BlackList, EnrollID)
}

func (this *TBlackList) Select(EnrollID string) (TErrorUsr, bool) {

	if this.BlackList == nil {
		this.BlackList = make(map[string]TErrorUsr, 0)
	}

	this.Lock.Lock()
	defer this.Lock.Unlock()
	ErrorUsr, Valid := this.BlackList[EnrollID]
	return ErrorUsr, Valid
}
