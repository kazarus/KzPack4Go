package KzPack4Go

// #同步配置
type TSynConfig struct {
	StartDay int64  `json:"startDay"` //#起始日期
	EndedDay int64  `json:"endedDay"` //#截止日期
	SynField string `json:"synField"` //#同步字段
	FileName string `json:"fileName"` //#输出文件
	FilePath string `json:"filePath"` //#输出路径
	SynToday int64  `json:"synToday"` //#是否同步当天

	ListTabl []TSynTable `json:"listTabl"` //#指标表名
}

type TSynTable struct {
	SynTable string `json:"synTable"` //#同步表名
	SynField string `json:"synField"` //#同步字段
}
