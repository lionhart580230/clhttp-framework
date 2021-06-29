package rule_list

import (
	"github.com/xiaolan580230/clhttp-framework/clController/example"
	"github.com/xiaolan580230/clhttp-framework/core/rule"
)

func Init() {

	// 范例: 注册账号接口
	rule.AddRule(rule.Rule{
		Request: "request",
		Name: "api_example",
		Params: []rule.ParamInfo{
			rule.NewParam("user", rule.PTYPE_ALL, true, ""),
		},
		CallBack: example.ApiExample,
	})
}
