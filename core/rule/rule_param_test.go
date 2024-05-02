package rule

import (
	"fmt"
	"github.com/lionhart580230/clUtil/clLog"
	"github.com/lionhart580230/clUtil/clTime"
	"testing"
)

func TestAddRequest(t *testing.T) {

	fmt.Printf("%v", paramCheckers[PTYPE_DATETIME]("2022-02-31 00:00:00", nil))

}

func TestApiGeneralByJson(t *testing.T) {

	jsonStr := `
	[
		{
			"ac": "login",
			"comment": "管理员登录",
			"params": [
				{"name": "username", "strict": true, "type": "PTYPE_USERNAME", "comment": "账号"},
				{"name": "password", "strict": true, "type": "PTYPE_PASSWORD", "comment": "管理员密码"}
			],
			"login": false
		},
		{
			"ac": "getMenu",
			"params": [
			],
			"login": true
		},
		{
			"ac": "getHomeStatic",
			"params": [
			],
			"login": true
		},
		{
			"ac": "getAdminList",
			"params": [
			   {"name": "username", "strict": false, "type": "PTYPE_USERNAME"}
			],
			"login": true
		},
		{
			"ac": "addNewAdmin",
			"params": [
			   {"name": "username", "strict": true, "type": "PTYPE_USERNAME"},
			   {"name": "password", "strict": true, "type": "PTYPE_PASSWORD"},
			   {"name": "ac_type", "strict": true, "type": "PTYPE_INT"}
			],
			"login": true
		}
	]
`

	ApiGeneralByJson(jsonStr, "./apis", "apis", "request", "rulelist")
}

func TestHttpParam_GetTimeStamp(t *testing.T) {
	a := HttpParam{values: map[string]string{
		"time": "2024-01-03 13:22:10",
	}}
	ts1 := a.GetTimeStamp("time")
	ts2 := clTime.GetTimeStamp2("2024-01-03 13:22:10", "2006-01-02 15:04:05")
	if ts1 == ts2 {
		clLog.Info("测试通过!")
	} else {
		clLog.Error("测试不通过! %v != %v", ts1, ts2)
	}
}

func TestHttpParam_GetBeginTime(t *testing.T) {
	a := HttpParam{values: map[string]string{
		"time": "2024-01-03",
	}}
	ts1 := a.GetBeginTime("time")
	ts2 := clTime.GetTimeStamp2("2024-01-03 00:00:00", "2006-01-02 15:04:05")
	if ts1 == ts2 {
		clLog.Info("测试通过!")
	} else {
		clLog.Error("测试不通过! %v != %v", ts1, ts2)
	}
}

func TestHttpParam_GetEndTime(t *testing.T) {
	a := HttpParam{values: map[string]string{
		"time": "2024-01-03",
	}}
	ts1 := a.GetEndTime("time")
	ts2 := clTime.GetTimeStamp2("2024-01-03 23:59:59", "2006-01-02 15:04:05")
	if ts1 == ts2 {
		clLog.Info("测试通过!")
	} else {
		clLog.Error("测试不通过! %v != %v", ts1, ts2)
	}
}
