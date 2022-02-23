package rule_list

import (
	"github.com/xiaolan580230/clhttp-framework/clController/example"
	"github.com/xiaolan580230/clhttp-framework/core/rule"
)

func Init() {

	// 添加请求类型
	rule.AddRequest("request", "api")

	// 范例: 注册账号接口
	rule.AddRule(rule.Rule{
		Request: "request",
		Name: "api_example",
		Params: []rule.ParamInfo{
			rule.NewParam("user", rule.PTYPE_ALL, true, ""),
		},
		Login: true,
		CallBack: example.ApiExample,
		CacheExpire: 180,
	})



	// 范例: 注册账号接口
	rule.AddRule(rule.Rule{
		Request: "upload",
		Name: "UploadFile",
		Params: []rule.ParamInfo{
			rule.NewParam("filename", rule.PTYPE_ALL, true, ""),
			rule.NewParam("fileExt", rule.PTYPE_ALL, true, ""),
			rule.NewParam("localPath", rule.PTYPE_ALL, true, ""),
		},
		CallBack: example.ApiUploadExample,
	})
}
