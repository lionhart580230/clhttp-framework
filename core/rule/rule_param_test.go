package rule

import (
	"fmt"
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