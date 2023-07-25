package rule_list

import (
	"github.com/lionhart580230/clhttp-framework/clController/example"
	"github.com/lionhart580230/clhttp-framework/core/rule"
	"github.com/lionhart580230/clhttp-framework/core/rule/apis"
)

func Init() {

	// 添加请求类型
	rule.AddRequest("request", "api")

	// 范例: 注册账号接口
	rule.AddRule(rule.Rule{
		Request: "request",
		Name:    "api_example",
		Params: []rule.ParamInfo{
			// 参数名为id，它必须是整数，并且值范围必须在1到10之间
			rule.NewIntParamRange("id", true, "1", 1, 10),
			// 参数名为name, 它必须是字符串，并且这个字符串的长度必须为2到5之间
			rule.NewStrParamRange("name", true, "", 2, 5),
		},
		Login:       true,
		CallBack:    example.ApiExample,
		CacheExpire: 180,
	})

	// 范例: 注册账号接口
	rule.AddRule(rule.Rule{
		Request: "upload",
		Name:    "UploadFile",
		Params: []rule.ParamInfo{
			rule.NewParam("filename", rule.PTYPE_ALL, true, ""),
			rule.NewParam("fileExt", rule.PTYPE_ALL, true, ""),
			rule.NewParam("localPath", rule.PTYPE_ALL, true, ""),
		},
		CallBack: example.ApiUploadExample,
	})

}

func InitSuperAPI() {

	// 添加请求类型
	rule.AddRequest("sys", "ac")

	rule.AddRule(rule.Rule{
		Request: "sys",
		Name:    "mysql_encrypt",
		Params: []rule.ParamInfo{
			rule.NewParam("p", rule.PTYPE_ALL, true, ""),
		},
		CallBack: apis.ApiMysqlEncrypt,
		Login:    false,
	})
}
