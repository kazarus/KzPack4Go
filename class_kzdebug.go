package KzPack4Go

import (
	"fmt"
)

type TKzDebug struct {
	ThisTick float64
	LastTick float64
	TotlTick float64
	RunDebug bool
}

func (this *TKzDebug) StartEd(aRunDebug bool) {
	this.ThisTick = 0
	this.LastTick = 0
	this.TotlTick = 0

	this.RunDebug = aRunDebug
}

func (this *TKzDebug) TickFmt(aTickMsg string) string {

	var Result = ""
	var DateTime = ToMsec()

	switch this.ThisTick {
	case 0:
		{
			this.ThisTick = DateTime
			this.LastTick = DateTime
			this.TotlTick = DateTime
			Result = fmt.Sprintf("|%15s |%15.3f | %15.3f | %15.3f |", aTickMsg, this.ThisTick, 0.0, 0.0)
		}
	default:
		{
			this.ThisTick = DateTime
			Result = fmt.Sprintf("|%15s |%15.3f | %15.3f | %15.3f |", aTickMsg, this.ThisTick, this.ThisTick-this.LastTick, this.ThisTick-this.TotlTick)
			this.LastTick = DateTime
		}
	}

	switch this.RunDebug {
	case true:
		{
			fmt.Println(Result)
		}
	}

	return Result
}
