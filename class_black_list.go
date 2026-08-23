package KzPack4Go

import (
	"sync"
	"time"
)

type TErrorUsr struct {
	UserCode string `json:"userCode"` //#用户编码
	ErrorCnt int64  `json:"errorCnt"` //#错误次数
}

// TBlackList 并发安全黑名单。
//
// 零值可直接使用(无 TTL、无容量上限,等同旧行为)。
// 通过 NewBlackList(MaxSize, TTL) 或直接设置字段启用限制:
//   - MaxSize > 0:达到上限后先清扫过期条目,仍超限则淘汰"最早到期/最早插入"的条目;
//   - TTL > 0:条目过期后 Select 视为不存在,并在下次 Append 时按 TTL/2 间隔整表清扫,
//     避免只设 TTL 时过期条目长期驻留内存。
type TBlackList struct {
	BlackList map[string]TErrorUsr
	ExpireAt  map[string]time.Time // #TTL>0:过期时刻;TTL==0:插入时刻(用于 FIFO 淘汰)
	Lock      sync.RWMutex

	MaxSize int           // #容量上限;<=0 表示不限制
	TTL     time.Duration // #条目有效期;<=0 表示永不过期

	lastClean time.Time // #上次整表清扫时刻(仅 TTL>0 时使用)

	// #内部:时钟函数(测试注入),nil 时使用 time.Now
	now func() time.Time
}

// NewBlackList 构造黑名单;maxSize<=0 不限制容量,ttl<=0 永不过期。
func NewBlackList(maxSize int, ttl time.Duration) *TBlackList {
	return &TBlackList{MaxSize: maxSize, TTL: ttl}
}

func (this *TBlackList) tick() time.Time {
	if this.now != nil {
		return this.now()
	}
	return time.Now()
}

// purgeExpired 整表清扫过期条目;仅 TTL>0 且距上次清扫超过 TTL/2 时执行,摊还 O(n)。
func (this *TBlackList) purgeExpired(now time.Time) {
	if this.TTL <= 0 {
		return
	}
	if this.lastClean.IsZero() || now.Sub(this.lastClean) > this.TTL/2 {
		for key, expire := range this.ExpireAt {
			if now.After(expire) {
				delete(this.BlackList, key)
				delete(this.ExpireAt, key)
			}
		}
		this.lastClean = now
	}
}

// evictOldest 容量超限时淘汰旧条目:先清过期(防御),再淘汰到期/插入时刻最早的条目。
func (this *TBlackList) evictOldest(now time.Time) {
	for key, expire := range this.ExpireAt {
		if now.After(expire) {
			delete(this.BlackList, key)
			delete(this.ExpireAt, key)
		}
	}

	for len(this.BlackList) >= this.MaxSize {
		var oldKey string
		var oldTick time.Time
		first := true
		for key, expire := range this.ExpireAt {
			// #到期时刻最早者优先;相同时按键序决胜,避免依赖 map 随机迭代序
			if first || expire.Before(oldTick) || (expire.Equal(oldTick) && key < oldKey) {
				oldKey, oldTick, first = key, expire, false
			}
		}
		if first {
			break // #空表
		}
		delete(this.BlackList, oldKey)
		delete(this.ExpireAt, oldKey)
	}
}

func (this *TBlackList) Append(EnrollID string, ErrorUsr TErrorUsr) {

	this.Lock.Lock()
	defer this.Lock.Unlock()

	now := this.tick()
	this.purgeExpired(now)

	if this.BlackList == nil {
		this.BlackList = make(map[string]TErrorUsr)
	}
	if this.ExpireAt == nil {
		this.ExpireAt = make(map[string]time.Time)
	}

	_, exists := this.BlackList[EnrollID]

	// #新增键且容量已满:先淘汰旧条目
	if !exists && this.MaxSize > 0 && len(this.BlackList) >= this.MaxSize {
		this.evictOldest(now)
	}

	this.BlackList[EnrollID] = ErrorUsr
	if this.TTL > 0 {
		this.ExpireAt[EnrollID] = now.Add(this.TTL)
	} else {
		this.ExpireAt[EnrollID] = now // #插入时刻,供 FIFO 淘汰
	}
}

func (this *TBlackList) Remove(EnrollID string) {

	this.Lock.Lock()
	defer this.Lock.Unlock()
	if this.BlackList == nil {
		this.BlackList = make(map[string]TErrorUsr)
	}
	delete(this.BlackList, EnrollID)
	if this.ExpireAt != nil {
		delete(this.ExpireAt, EnrollID)
	}
}

func (this *TBlackList) Select(EnrollID string) (TErrorUsr, bool) {

	this.Lock.RLock()
	defer this.Lock.RUnlock()
	if this.BlackList == nil {
		return TErrorUsr{}, false
	}
	ErrorUsr, Valid := this.BlackList[EnrollID]
	if !Valid {
		return ErrorUsr, false
	}
	// #TTL>0 时检查过期:过期视为不存在(惰性清除交由 Append 的整表清扫)
	if this.TTL > 0 {
		if expire, ok := this.ExpireAt[EnrollID]; ok && this.tick().After(expire) {
			return TErrorUsr{}, false
		}
	}
	return ErrorUsr, Valid
}
