package cltask

import (
	"fmt"
	"testing"
)

func TestNewPool(t *testing.T) {

	pool := NewPool(nil)
	pool.AddNew(NewTaskPerSec("aa", "aa", func(){
		fmt.Printf("执行..\n")
		var a map[string] string
		a["1"] = "a"
	}, 10, true))

	// 添加异常捕获回调机制,
	// tag 为任务的tag
	// panic 为异常捕获信息
	pool.SetRecoverCallback(func(tag string, panic string){
		fmt.Printf("异常捕获:[%v]\n%v\n", tag, panic)
	})

	pool.Start()
}