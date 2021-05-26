package example

import (
	"clhttp-framework/clCommon"
	"clhttp-framework/core/clAuth"
	"clhttp-framework/core/rule"
)

func ApiExample(_auth *clAuth.AuthInfo, _param *rule.HttpParam, _server *rule.ServerParam) string {

	return clCommon.JCode(0, "ok", nil)
}