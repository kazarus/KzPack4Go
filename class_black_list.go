package KzPack4Go

import "sync"

type TErrorUsr struct {
	UserCode string `json:"userCode"` //#用户编码
	ErrorCnt int64  `json:"errorCnt"` //#错误次数
}

type TBlackList struct {
	BlackList map[string]TErrorUsr
	Lock      sync.RWMutex
}

func (this *TBlackList) Append(EnrollID string, ErrorUsr TErrorUsr) {

	this.Lock.Lock()
	defer this.Lock.Unlock()
	if this.BlackList == nil {
		this.BlackList = make(map[string]TErrorUsr)
	}
	this.BlackList[EnrollID] = ErrorUsr
}

func (this *TBlackList) Remove(EnrollID string) {

	this.Lock.Lock()
	defer this.Lock.Unlock()
	if this.BlackList == nil {
		this.BlackList = make(map[string]TErrorUsr)
	}
	delete(this.BlackList, EnrollID)
}

func (this *TBlackList) Select(EnrollID string) (TErrorUsr, bool) {

	this.Lock.RLock()
	defer this.Lock.RUnlock()
	if this.BlackList == nil {
		return TErrorUsr{}, false
	}
	ErrorUsr, Valid := this.BlackList[EnrollID]
	return ErrorUsr, Valid
}
