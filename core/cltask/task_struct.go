package cltask

import (
	"github.com/xiaolan580230/clhttp-framework/core/skyredis"
	"sync"
)

// 计划任务结构
type TaskInfo struct {
	Tag string				// 任务Tag，便于做SetNX
	Name string				// 回调任务名称
	Callback func ()		// 回调处理
	Delay uint32			// 每个计划区间的间隔多久
	Repeattime uint32		// 在区间内每多少秒执行一次
	StartTime uint32		// 每天起始秒数
	Duration uint32			// 持续多少秒
	Lasttime uint32			// 上次执行时间
	Week uint32				// 是否限制周几执行, 位运算
	running bool			// 是否正在运行中
}

const (
	WEEK_MONDAY = 1
	WEEK_TUESDAY = 2
	WEEK_WEDNESDAY = 4
	WEEK_THURDAY = 8
	WEEK_FRIDAY = 16
	WEEK_SATURDAY = 32
	WEEK_SUNDAY = 64

	WEEK_WORKING_DAY = WEEK_MONDAY | WEEK_TUESDAY | WEEK_WEDNESDAY | WEEK_THURDAY | WEEK_FRIDAY
	WEEK_ALL = WEEK_WORKING_DAY | WEEK_SATURDAY | WEEK_SUNDAY
)

// 任务池
type TaskPool struct{
	pool []*TaskInfo
	locker sync.RWMutex
	redisPtr *skyredis.RedisObject
	recoverCallback func(string, string)
}