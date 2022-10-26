package rulelist

import (
	"github.com/xiaolan580230/clhttp-framework/core/rule"
	"github.com/xiaolan580230/clhttp-framework/core/rule/apis"
)

func Init() {

	// 添加请求类型
	rule.AddRequest("request", "ac")

// 接口列表

	// 管理员登录
	rule.AddRule(rule.Rule{
		Request:     "request",
		Name:        "login",
		Params:      []rule.ParamInfo{
			rule.NewParam("username", rule.PTYPE_USERNAME, true, ""),
			rule.NewParam("password", rule.PTYPE_PASSWORD, true, ""),
		},
		CallBack:    apis.ApiLogin,
		CacheExpire: 0,
		CacheType:   0,
		Login:       false,
	})
	


	
	rule.AddRule(rule.Rule{
		Request:     "request",
		Name:        "getMenu",
		Params:      []rule.ParamInfo{
		},
		CallBack:    apis.ApiGetMenu,
		CacheExpire: 0,
		CacheType:   0,
		Login:       true,
	})
	


	
	rule.AddRule(rule.Rule{
		Request:     "request",
		Name:        "getHomeStatic",
		Params:      []rule.ParamInfo{
		},
		CallBack:    apis.ApiGetHomeStatic,
		CacheExpire: 0,
		CacheType:   0,
		Login:       true,
	})
	


	
	rule.AddRule(rule.Rule{
		Request:     "request",
		Name:        "getAdminList",
		Params:      []rule.ParamInfo{
			rule.NewParam("username", rule.PTYPE_USERNAME, false, ""),
		},
		CallBack:    apis.ApiGetAdminList,
		CacheExpire: 0,
		CacheType:   0,
		Login:       true,
	})
	


	
	rule.AddRule(rule.Rule{
		Request:     "request",
		Name:        "addNewAdmin",
		Params:      []rule.ParamInfo{
			rule.NewParam("username", rule.PTYPE_USERNAME, true, ""),
			rule.NewParam("password", rule.PTYPE_PASSWORD, true, ""),
			rule.NewParam("ac_type", rule.PTYPE_INT, true, ""),
		},
		CallBack:    apis.ApiAddNewAdmin,
		CacheExpire: 0,
		CacheType:   0,
		Login:       true,
	})
	



}