package rule_list

import (
	"clhttp-framework/clController/example"
	"clhttp-framework/core/rule"
)

func Init() {

	// 范例: 注册账号接口
	rule.AddRule(rule.Rule{
		Request: "request",
		Name: "api_example",
		Params: []rule.ParamInfo{
			rule.NewParam("user", rule.PTYPE_SAFE_STR, true, ""),
		},
		CallBack: example.ApiExample,
	})
}
