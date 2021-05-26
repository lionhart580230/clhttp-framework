package cltime

import (
	"fmt"
	"math"
	"strings"
	"time"
)

/**
  时间类
*/

// 透过时间日期新建一个时间对象
func NewDate(date string) (*clTimer, error) {

	nowtime := getNowTime(date)

	weeks := 0
	switch strings.ToLower(nowtime.Weekday().String()) {
	case "monday":
		weeks = 0
	case "tuesday":
		weeks = 1
	case "wednesday":
		weeks = 2
	case "thursday":
		weeks = 3
	case "friday":
		weeks = 4
	case "saturday":
		weeks = 5
	case "sunday":
		weeks = 6
	}

	var cltimer = clTimer{
		TimeStamp: uint32(nowtime.Unix()),
		Year:      uint32(nowtime.Year()),
		Month:     uint8(nowtime.Month()),
		Days:      uint8(nowtime.Day()),
		Hour:      uint8(nowtime.Hour()),
		Minuter:   uint8(nowtime.Minute()),
		Second:    uint8(nowtime.Second()),
		Week:      uint8(weeks),
	}
	return &cltimer, nil
}

// 透过时间戳新建一个时间对象
func NewTime(timestamp uint32) (*clTimer, error) {

	nowtime := getTargetTime(timestamp)

	weeks := 0
	switch strings.ToLower(nowtime.Weekday().String()) {
	case "monday":
		weeks = 0
	case "tuesday":
		weeks = 1
	case "wednesday":
		weeks = 2
	case "thursday":
		weeks = 3
	case "friday":
		weeks = 4
	case "saturday":
		weeks = 5
	case "sunday":
		weeks = 6
	}

	var cltimer = clTimer{
		TimeStamp: uint32(nowtime.Unix()),
		Year:      uint32(nowtime.Year()),
		Month:     uint8(nowtime.Month()),
		Days:      uint8(nowtime.Day()),
		Hour:      uint8(nowtime.Hour()),
		Minuter:   uint8(nowtime.Minute()),
		Second:    uint8(nowtime.Second()),
		Week:      uint8(weeks),
	}
	return &cltimer, nil
}

// 偏移指定秒数
func (this *clTimer) AfterSecond(sec uint32) (*clTimer, error) {

	if sec == 0 {
		return this, nil
	}

	overTime := int64(this.TimeStamp + sec)
	nowtime := time.Unix(overTime, 0).UTC().Add(8 * time.Hour)
	weeks := 0
	switch strings.ToLower(nowtime.Weekday().String()) {
	case "monday":
		weeks = 0
	case "tuesday":
		weeks = 1
	case "wednesday":
		weeks = 2
	case "thursday":
		weeks = 3
	case "friday":
		weeks = 4
	case "saturday":
		weeks = 5
	case "sunday":
		weeks = 6
	}

	this = &clTimer{
		TimeStamp: uint32(overTime),
		Year:      uint32(nowtime.Year()),
		Month:     uint8(nowtime.Month()),
		Days:      uint8(nowtime.Day()),
		Hour:      uint8(nowtime.Hour()),
		Minuter:   uint8(nowtime.Minute()),
		Second:    uint8(nowtime.Second()),
		Week:      uint8(weeks),
	}
	return this, nil
}

// 获取星期几的文本表示方式
func (this *clTimer) GetWeekStr() string {
	switch this.Week {
	case 0:
		return "星期一"
	case 1:
		return "星期二"
	case 2:
		return "星期三"
	case 3:
		return "星期四"
	case 4:
		return "星期五"
	case 5:
		return "星期六"
	case 6:
		return "星期日"
	default:
		return "NULL"
	}
}

// 获取本月开始时间
// @return uint32 当前时间的本月1号0时0分0秒的时间戳
// @return string 当前时间的本月1号0时0分0秒的日期时间格式
func (this *clTimer) GetCurMonth() ( /*timestamp*/ uint32 /*datestr*/, string) {
	dateformat := fmt.Sprintf("%04v-%02v-01 00:00:00", this.Year, this.Month)
	targetTime := getNowTime(dateformat)
	return uint32(targetTime.Unix()), targetTime.Format("2006-01-02 15:04:05")
}

// 获取本周的开始时间
// @return uint32 当前时间的本周星期一的时间戳
// @return string 当前时间的本周星期一的日期时间戳
func (this *clTimer) GetCurWeek() ( /*timestamp*/ uint32 /*datestr*/, string) {
	dateformat := fmt.Sprintf("%04v-%02v-%02v 00:00:00", this.Year, this.Month, this.Days)
	todayTime := getNowTime(dateformat)
	targetTime := todayTime.Add(-time.Duration(this.Week) * 24 * time.Hour)
	return uint32(targetTime.Unix()), targetTime.Format("2006-01-02 15:04:05")
}

// 获取指定跨度的月份时间戳区间
// @param offset int 跨度偏移， 0为当前月份时间周期
// @return uint32 指定时间区间起始时间
// @return uint32 指定时间区间结束时间
func (this *clTimer) GetMonthBetween(offset int) (uint32, uint32) {

	var beginMonth = int(this.Month) + offset
	var beginYear = int(this.Year)
	if beginMonth < 1 || beginMonth > 12 {
		if offset < 0 {
			beginYear = int(this.Year) + beginMonth/12 - 1
			beginMonth = 12 - (int(math.Abs(float64(beginMonth))) % 12)
		} else {
			beginMonth = (int(this.Month) + offset) % 12
			beginYear = int(this.Year) + (beginMonth-1)/12
		}
	}

	var endMonth = beginMonth + 1
	var endYear = beginYear
	if endMonth > 12 {
		endMonth = 1
		endYear++
	}

	begins := GetTimeStamp(fmt.Sprintf("%04v-%02v-01 00:00:00", beginYear, beginMonth))
	ends := GetTimeStamp(fmt.Sprintf("%04v-%02v-01 00:00:00", endYear, endMonth))
	return begins, ends - 1
}

// 跨年
func (this *clTimer) GetYearBetween(offset int) (uint32, uint32) {
	beginYear := int(this.Year) - offset
	begins := GetTimeStamp(fmt.Sprintf("%04v-01-01 00:00:00", beginYear))
	ends := GetTimeStamp(fmt.Sprintf("%04v-01-01 00:00:00", beginYear+1))
	return begins, ends - 1
}

// 获取指定跨度的月份时间戳区间
// @param offset int 跨度偏移， 0为当前月份时间周期
// @return uint64 指定时间区间起始时间
// @return uint64 指定时间区间结束时间
func (this *clTimer) GetMonthBetweenWithMSec(offset int) (uint64, uint64) {

	var beginMonth = int(this.Month) + offset
	var beginYear = int(this.Year)
	if beginMonth < 1 || beginMonth > 12 {
		if offset < 0 {
			beginYear = int(this.Year) + beginMonth/12 - 1
			beginMonth = 12 - (beginMonth % 12)
		} else {
			beginMonth = (int(this.Month) + offset) % 12
			beginYear = int(this.Year) + (beginMonth-1)/12
		}
	}

	var endMonth = beginMonth + 1
	var endYear = beginYear
	if endMonth > 12 {
		endMonth = 1
		endYear++
	}

	begins := GetTimeStampWithMSec(fmt.Sprintf("%04v-%02v-01 00:00:00", beginYear, beginMonth))
	ends := GetTimeStampWithMSec(fmt.Sprintf("%04v-%02v-01 00:00:00", endYear, endMonth))
	return begins, ends - 1000
}

// 获取指定跨度的小时时间戳区间
// @param offset int 跨度偏移， 0为当前小时时间周期
// @return uint32 指定时间区间起始时间
// @return uint32 指定时间区间结束时间
func (this *clTimer) GetHourBetween(offset int) (uint32, uint32) {

	nowTime := uint32(this.TimeStamp)
	beginTime := nowTime - (nowTime % 3600)
	endTime := beginTime + 3599

	return beginTime, endTime
}

// 获取指定跨度的周时间戳区间
// @param offset int 跨度偏移， 0为当前周时间周期
// @return uint32 指定时间区间起始时间
// @return uint32 指定时间区间结束时间
func (this *clTimer) GetWeekBetween(offset int) (uint32, uint32) {

	beginTime, _ := NewTime(this.TimeStamp + uint32(offset*7*86400))
	endTime, _ := NewTime(this.TimeStamp + uint32((offset+1)*7*86400))

	begins := GetTimeStamp(fmt.Sprintf("%04v-%02v-%02v 00:00:00", beginTime.Year, beginTime.Month, beginTime.Days))
	if beginTime.Week > 0 {
		begins = begins - uint32(endTime.Week)*86400
	}
	ends := GetTimeStamp(fmt.Sprintf("%04v-%02v-%02v 00:00:00", endTime.Year, endTime.Month, endTime.Days))
	if endTime.Week > 0 {
		ends = ends - uint32(endTime.Week)*86400
	}
	return begins, ends - 1
}

// 获取指定跨度的周时间戳区间
// @param offset int 跨度偏移， 0为当前周时间周期
// @return uint64 指定时间区间起始时间
// @return uint64 指定时间区间结束时间
func (this *clTimer) GetWeekBetweenWithMSec(offset int) (uint64, uint64) {

	beginTime, _ := NewTime(this.TimeStamp + uint32(offset*7*86400))
	endTime, _ := NewTime(this.TimeStamp + uint32((offset+1)*7*86400))

	begins := GetTimeStampWithMSec(fmt.Sprintf("%04v-%02v-%02v 00:00:00", beginTime.Year, beginTime.Month, beginTime.Days))

	fmt.Printf("begins: %v\n", begins)
	if beginTime.Week > 0 {
		begins = begins - uint64(endTime.Week)*86400000
	}
	ends := GetTimeStampWithMSec(fmt.Sprintf("%04v-%02v-%02v 00:00:00", endTime.Year, endTime.Month, endTime.Days))
	if endTime.Week > 0 {
		ends = ends - uint64(endTime.Week)*86400000
	}
	return begins, ends - 1000
}

// 获取指定跨度的天时间戳区间
// @param offset int 跨度偏移， 0为当前周时间周期
// @return uint32 指定时间区间起始时间
// @return uint32 指定时间区间结束时间
func (this *clTimer) GetDayBetween(offset int) (uint32, uint32) {

	begins := GetTimeStamp(fmt.Sprintf("%04v-%02v-%02v 00:00:00", this.Year, this.Month, this.Days))
	btime := begins + uint32(offset*86400)
	return btime, btime + 86400 - 1
}

// 获取指定跨度的天时间戳区间單位為 milisecond
// @param offset int 跨度偏移， 0为当前周时间周期
// @return uint64 指定时间区间起始时间
// @return uint64 指定时间区间结束时间
func (this *clTimer) GetDayBetweenWithMSec(offset int) (uint64, uint64) {

	begins := GetTimeStampWithMSec(fmt.Sprintf("%04v-%02v-%02v 00:00:00", this.Year, this.Month, this.Days))
	btime := begins + uint64(offset*86400000)
	return btime, btime + (86400-1)*1000
}

// 获取UTC时间
// @param date string 时间日期格式
// @return *time.Time 返回这个日期格式生成的时间指针
func getNowTime(date string) *time.Time {

	dateTimeArr := strings.Split(date, " ")
	if date == "" || len(dateTimeArr) != 2 {
		date = time.Now().UTC().Add(8 * time.Hour).Format("2006-01-02 15:04:05")
	} else {
		now := time.Now()
		dateArr := strings.Split(dateTimeArr[0], "-")
		timeArr := strings.Split(dateTimeArr[1], ":")

		Year := dateArr[0]
		if Year == "Y" {
			Year = now.Format("2006")
		}
		Month := dateArr[1]
		if Month == "m" {
			Month = now.Format("01")
		}
		Day := dateArr[2]
		if Day == "d" {
			Day = now.Format("02")
		}
		Hour := timeArr[0]
		if Hour == "H" {
			Hour = now.Format("15")
		}
		Minute := timeArr[1]
		if Minute == "i" {
			Minute = now.Format("04")
		}
		Second := timeArr[2]
		if Second == "s" {
			Second = now.Format("05")
		}

		date = fmt.Sprintf("%v-%v-%v %v:%v:%v", Year, Month, Day, Hour, Minute, Second)
	}

	timedate, _ := time.Parse("2006-01-02 15:04:05", date)
	utc := timedate.UTC().Add(-8 * time.Hour)
	//utc := timedate.UTC()
	loc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		fmt.Printf("LoadLocation error: %v\n", err)
		return nil
	}

	utcTime := utc.In(loc)
	return &utcTime
}

func getTargetTime(timestamp uint32) *time.Time {
	timedate := time.Unix(int64(timestamp), 0)
	utc := timedate.UTC().Add(-8 * time.Hour)
	loc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		return nil
	}

	utcTime := utc.In(loc)
	return &utcTime
}

// 获取指定时间日期的时间戳
func GetTimeStamp(date string) uint32 {
	timedate, _ := time.Parse("2006-01-02 15:04:05", date)
	utc := timedate.UTC().Add(-8 * time.Hour)
	//utc := timedate.UTC()
	loc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		return 0
	}

	return uint32(utc.In(loc).Unix())
}

//指定格式获取日期的时间戳
func GetTimeStamp2(date string,format string) uint32 {
	timedate, _ := time.Parse(format, date)
	utc := timedate.UTC().Add(-8 * time.Hour)
	//utc := timedate.UTC()
	loc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		return 0
	}

	return uint32(utc.In(loc).Unix())
}


// 获取指定时间日期的时间戳單位為 milisecond
func GetTimeStampWithMSec(date string) uint64 {
	timedate, _ := time.Parse("2006-01-02 15:04:05", date)
	return uint64(timedate.UnixNano() / int64(time.Millisecond))
}

// 获取指定时间戳的日期格式
func GetDate(timestamp uint32) string {
	utc := time.Unix(int64(timestamp), 0).UTC()
	loc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		return "1970-01-01 00:00:00"
	}

	return utc.In(loc).Format("2006-01-02 15:04:05")
}

// 获取指定时间戳的日期格式
func GetDateByFormat(timestamp uint32, format string) string {
	utc := time.Now().UTC()
	if timestamp > 0 {
		utc = time.Unix(int64(timestamp), 0).UTC()
	}

	loc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		return ""
	}

	return utc.In(loc).Format(format)
}



// 检测某个时间是否是在除夕中
// @param timestamp uint32 需要检测的时间戳
func CheckIsChuxi(timestamp uint32) bool {
	targetTime, _ := NewTime(timestamp)

	btime, etime := GetYearTimeBetween(int32(targetTime.Year))

	if btime <= timestamp && etime > timestamp {
		return true
	}
	return false
}

// 获取指定年份的除夕时间区间
// @param year int32 年份
func GetYearTimeBetween(year int32) (/*btime*/uint32, /*etime*/uint32) {
	beginTime := uint32(0)
	switch year {
	case 2018:
		chuxiBegin, _ := NewDate("2018-02-15 00:00:00")
		beginTime = chuxiBegin.TimeStamp
	case 2019:
		chuxiBegin, _ := NewDate("2019-02-04 00:00:00")
		beginTime = chuxiBegin.TimeStamp
	case 2020:
		chuxiBegin, _ := NewDate("2020-01-24 00:00:00")
		beginTime = chuxiBegin.TimeStamp
	case 2021:
		chuxiBegin, _ := NewDate("2021-02-11 00:00:00")
		beginTime = chuxiBegin.TimeStamp
	case 2022:
		chuxiBegin, _ := NewDate("2022-01-31 00:00:00")
		beginTime = chuxiBegin.TimeStamp
	case 2023:
		chuxiBegin, _ := NewDate("2023-01-21 00:00:00")
		beginTime = chuxiBegin.TimeStamp
	case 2024:
		chuxiBegin, _ := NewDate("2024-02-09 00:00:00")
		beginTime = chuxiBegin.TimeStamp
	case 2025:
		chuxiBegin, _ := NewDate("2025-01-28 00:00:00")
		beginTime = chuxiBegin.TimeStamp
	case 2026:
		chuxiBegin, _ := NewDate("2026-02-16 00:00:00")
		beginTime = chuxiBegin.TimeStamp
	case 2027:
		chuxiBegin, _ := NewDate("2027-02-05 00:00:00")
		beginTime = chuxiBegin.TimeStamp
	case 2028:
		chuxiBegin, _ := NewDate("2028-01-25 00:00:00")
		beginTime = chuxiBegin.TimeStamp
	case 2029:
		chuxiBegin, _ := NewDate("2029-02-12 00:00:00")
		beginTime = chuxiBegin.TimeStamp
	case 2030:
		chuxiBegin, _ := NewDate("2030-02-02 00:00:00")
		beginTime = chuxiBegin.TimeStamp
	case 2031:
		chuxiBegin, _ := NewDate("2031-01-22 00:00:00")
		beginTime = chuxiBegin.TimeStamp
	case 2032:
		chuxiBegin, _ := NewDate("2032-02-10 00:00:00")
		beginTime = chuxiBegin.TimeStamp
	case 2033:
		chuxiBegin, _ := NewDate("2033-01-30 00:00:00")
		beginTime = chuxiBegin.TimeStamp
	case 2034:
		chuxiBegin, _ := NewDate("2034-02-18 00:00:00")
		beginTime = chuxiBegin.TimeStamp
	case 2035:
		chuxiBegin, _ := NewDate("2035-02-07 00:00:00")
		beginTime = chuxiBegin.TimeStamp
	case 2036:
		chuxiBegin, _ := NewDate("2036-01-27 00:00:00")
		beginTime = chuxiBegin.TimeStamp
	case 2037:
		chuxiBegin, _ := NewDate("2037-02-14 00:00:00")
		beginTime = chuxiBegin.TimeStamp
	case 2038:
		chuxiBegin, _ := NewDate("2038-02-03 00:00:00")
		beginTime = chuxiBegin.TimeStamp
	}
	return beginTime, beginTime+7*86400
}


// 间隔多少天的除夕
func BetweenChuxiDays(timestamp uint32, target uint32) uint32 {

	// 过滤垃圾参数
	if timestamp > target {
		return 0
	}
	oldtime, _ := NewTime(timestamp)
	newtime, _ := NewTime(target)

	if oldtime.Year == newtime.Year {
		btime, etime := GetYearTimeBetween(int32(oldtime.Year))

		if oldtime.TimeStamp < btime && newtime.TimeStamp < btime {
			// 都在除夕之前
			return 0
		} else if oldtime.TimeStamp > etime && newtime.TimeStamp > etime {
			// 都在除夕之后
			return 0
		}
		// 一前一后
		return 7
	}


	yearBetween := newtime.Year - oldtime.Year
	betweenDays := (yearBetween-1) * 7

	obtime, _ := GetYearTimeBetween(int32(oldtime.Year))

	if oldtime.TimeStamp < obtime {
		// 旧时间在除夕之前 + 7天
		betweenDays += 7
	}

	_, netime := GetYearTimeBetween(int32(newtime.Year))
	if newtime.TimeStamp > netime {
		betweenDays += 7
	}

	return betweenDays
}

