package skyredis

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	rd, err := New("r-j6ckhdr61vu28fn4bzpd.redis.rds.aliyuncs.com:6379","MyInstance001", "")
	if err != nil {
		fmt.Printf(">> connect To redis error: %v\n", err)
		return
	}


	fmt.Printf(">> connect to redis success!!")

	rd.Set("key1", "hello", 600)
}
