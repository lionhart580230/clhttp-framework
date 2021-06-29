package cltask

import (
	"fmt"
	"github.com/xiaolan580230/clhttp-framework/core/cltime"
	"runtime"
	"strings"
)

// 添加一个计划任务到系统
// 注: 可通过设置starttime和duration来设置定时任务
//    如 starttime=3600, duration=0, delay=86400 则任务到每天的01:00则会执行一次
//    如 starttime=0, duration=86400, delay=0, repeattime=60 则任务每60秒就会执行一次
//    如 starttime=8*3600, duration=3600, delay=86400, repeattime=60 则任务会在每天的08:00到09:00时间段内每60秒执行一次
// @param name string 计划任务名称
// @param callback func() 计划回调
// @param delay uint32 计划区间间隔, 从starttime开始计算
// @param repeattime uint32 计划任务重复时间
// @param starttime uint32 起始时间
// @param duration uint32 持续时间
func NewTask(tag string, name string, callback func(), delay uint32, starttime uint32, duration uint32, repeattime uint32, run bool) *TaskInfo {

	return NewTaskByWeek(tag, name, callback, delay, starttime, duration, repeattime, WEEK_ALL, run)
}


// 添加一个计划任务到系统
// 注: 可通过设置starttime和duration来设置定时任务
//    如 starttime=3600, duration=0, delay=86400 则任务到每天的01:00则会执行一次
//    如 starttime=0, duration=86400, delay=0, repeattime=60 则任务每60秒就会执行一次
//    如 starttime=8*3600, duration=3600, delay=86400, repeattime=60 则任务会在每天的08:00到09:00时间段内每60秒执行一次
// @param name string 计划任务名称
// @param callback func() 计划回调
// @param delay uint32 计划区间间隔, 从starttime开始计算
// @param repeattime uint32 计划任务重复时间
// @param starttime uint32 起始时间
// @param duration uint32 持续时间
func NewTaskByWeek(tag string, name string, callback func(), delay uint32, starttime uint32, duration uint32, repeattime uint32, week uint32, run bool) *TaskInfo {
	task := TaskInfo{
		Tag: tag,
		Name: name,
		Callback: callback,
		Delay: delay,
		StartTime: starttime,
		Duration: duration,
		Repeattime: repeattime,
		Lasttime: 0,
		Week: week,
	}

	if !run {
		nowTime, _ := cltime.NewDate("")
		task.Lasttime = nowTime.TimeStamp
	}

	return &task
}


// 添加一个每天固定几点执行的任务
// @param name string 任务名称
// @param callback func() 指定的任务回调
// @param hour uint32 每天的几点开始执行
// @param run bool 程序重启的时候是否执行
func NewTaskPerdayHour(tag string, name string, callback func(), hour uint32, run bool) *TaskInfo {
	return NewTask(tag, name, callback, 86400, hour * 3600, 3600, 3600, run)
}

// 添加一个每天固定几秒执行的任务
// @param name string 任务名称
// @param callback func() 指定的任务回调
// @param hour uint32 每天的几点开始执行
// @param run bool 程序重启的时候是否执行
func NewTaskPerdaySec(tag string, name string, callback func(), sec uint32, run bool) *TaskInfo {
	return NewTask(tag, name, callback, 86400, sec, 3600, 3600, run)
}

// 每天整点触发
func NewTaskPerHour(tag string, name string, callback func(), run bool) *TaskInfo {
	return NewTask(tag, name, callback, 3600, 0, 86400, 3600, run)
}

// 每天整10分钟触发
func NewTaskPerTenMinute(tag string, name string, callback func(), run bool) *TaskInfo {
	return NewTask(tag, name, callback, 600, 0, 86400, 600, run)
}

// 每天整分钟触发
func NewTaskPerMinute(tag string, name string, callback func(), run bool) *TaskInfo {
	return NewTask(tag, name, callback, 60, 0, 86400, 60, run)
}

// 每周几几秒执行
func NewTaskPerWeekSec(tag string, name string, week uint32, sec uint32, callback func(), run bool) *TaskInfo {
	return NewTaskByWeek(tag, name, callback, 86400, sec, 86400, 86400, week, run)
}


// 添加一个每多少秒执行一次的任务
// @param name string 任务名称
// @param callback func() 指定的任务回调
// @param sec uint32 多少秒执行一次
// @param run bool 程序重启的时候是否执行
func NewTaskPerSec(tag string, name string, callback func(), sec uint32, run bool) *TaskInfo {
	return NewTask(tag, name, callback, 0, 0, 86400, sec, run)
}


// 启动任务执行判断
func (this *TaskInfo) Run(nowTime uint32, dayTime uint32, week uint32, _recoverCallback func(string, string)) {

	if _recoverCallback != nil {
		defer func() {
			if r := recover(); r != nil {

				_panicStr := strings.Builder{}
				_panicStr.WriteString(fmt.Sprintf("异常: %v\n", r))

				// 最多显示10层
				for i :=0; i < 10; i++ {
					ptr, file, line, ok := runtime.Caller(i)
					if !ok {
						break
					}

					_panicStr.WriteString(fmt.Sprintf("%v\n%v:%v\n", runtime.FuncForPC(ptr).Name(), file, line))

				}

				_recoverCallback(this.Tag, _panicStr.String())
			}
		}()
	}


	if this.Week & week == 0 {
		return
	}

	if this.Lasttime > 0 {
		if this.StartTime > dayTime || this.StartTime + this.Duration < dayTime {
			return
		}

		if this.Lasttime + this.Repeattime > nowTime  {
			if this.Delay > 0 {
				if  (this.StartTime + dayTime) % this.Delay != 0 {
					return
				}
			} else {
				return
			}
		}
	}

	if redisPtr == nil {
		if this.running {
			this.Lasttime = nowTime
			return
		}
	}

	this.Lasttime = nowTime
	if this.Callback != nil {

		if redisPtr  != nil {
			if !redisPtr.SetNx("NX_" + this.Tag, "1", 600) {
				return
			}
		} else {
			this.running = true
		}

		this.Callback()

		if redisPtr  != nil {
			redisPtr.Del("NX_" + this.Tag)
		} else {
			this.running = false
		}
	}
}