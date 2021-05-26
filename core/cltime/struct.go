package cltime


type clTimer struct {
	TimeStamp uint32		// 时间戳
	Hour uint8				// 小时
	Minuter uint8			// 分钟
	Second uint8			// 秒数
	Year uint32				// 年份
	Month uint8				// 月份
	Days uint8				// 天数
	Week uint8				// 周几，0=周一
}