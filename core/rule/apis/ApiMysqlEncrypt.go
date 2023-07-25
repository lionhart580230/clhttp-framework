package apis

import (
	"github.com/lionhart580230/clhttp-framework/clGlobal"
	"github.com/lionhart580230/clhttp-framework/clResponse"
	"github.com/lionhart580230/clhttp-framework/core/clAuth"
	"github.com/lionhart580230/clhttp-framework/core/rule"
	"strings"
)

// 数据库加密
func ApiMysqlEncrypt(_auth *clAuth.AuthInfo, _param *rule.HttpParam, _server *rule.ServerParam) string {
	_mysqlConfs := strings.Split(_param.GetStr("p", ""), "$$")

	// 格式检查
	for i := 0; i < len(_mysqlConfs); i++ {
		mysqlItems := strings.Split(_mysqlConfs[i], "|")
		if len(mysqlItems) != 4 {
			return clResponse.Failed(3, "请求参数错误", nil)
		}
	}

	retStr := clGlobal.EncryptMysql(_param.GetStr("p", ""))
	return clResponse.Success(retStr)
}
