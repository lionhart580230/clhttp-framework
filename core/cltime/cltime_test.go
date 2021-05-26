package cltime

import (
	"testing"
	"fmt"
)

func TestNew(t *testing.T) {

	nowtime, _ := NewDate("")

	fmt.Printf(">> 年份: %v\n", nowtime.Year)
	fmt.Printf(">> 月份: %v\n", nowtime.Month)
	fmt.Printf(">> 日期: %v\n", nowtime.Days)
	fmt.Printf(">> 小时: %v\n", nowtime.Hour)
	fmt.Printf(">> 分钟: %v\n", nowtime.Minuter)
	fmt.Printf(">> 秒数: %v\n", nowtime.Second)
	fmt.Printf(">> 星期: %v(%v)\n", nowtime.GetWeekStr(), nowtime.Week)
	fmt.Printf(">> 时间戳: %v\n", nowtime.TimeStamp)

	// 测试获取本月数据
	curmonthtime, curmonthdate := nowtime.GetCurMonth()
	fmt.Printf(">> 本月时间戳: %v %v\n", curmonthtime, curmonthdate )

	// 测试获取本周数据
	curweektime, curweekdate := nowtime.GetCurWeek()
	fmt.Printf(">> 本周时间戳: %v %v\n", curweektime, curweekdate )


	// 测试上周时间区间
	lastweekbegin, lastweekend := nowtime.GetWeekBetween(-1)
	fmt.Printf(">> 上周时间戳: %v - %v\n", lastweekbegin, lastweekend)
	fmt.Printf(">> 上周日期: %v - %v\n", GetDate(lastweekbegin), GetDate(lastweekend))


	// 测试上月时间区间
	lastmonthbegin, lastmonthend := nowtime.GetMonthBetween(-4)
	fmt.Printf(">> 上三月时间戳: %v - %v\n", lastmonthbegin, lastmonthend)
	fmt.Printf(">> 上三月日期: %v - %v\n", GetDate(lastmonthbegin), GetDate(lastmonthend))

	// 测试这个小时时间
	curHourBegin, curHourEnd := nowtime.GetHourBetween(0)
	fmt.Printf(">> 小时时间戳: %v - %v\n", curHourBegin, curHourEnd)
	fmt.Printf(">> 小时日期: %v - %v\n", GetDate(curHourBegin), GetDate(curHourEnd))

}