package skyredis

import (
	"testing"
)

func TestNew(t *testing.T) {
	//rd, _ := New("localhost:6379","MyInstance001", "")
	//if err != nil {
	//	fmt.Printf(">> connect To redis error: %v\n", err)
	//	return
	//}

	clrd := &RedisObject{
		myredis:   nil,
		prefix:    "",
		isCluster: false,
	}
	//fmt.Printf(">> connect to redis success!!")

	clrd.Set("key1", "hello", 600)
}
