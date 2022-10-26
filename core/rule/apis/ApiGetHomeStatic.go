package apis

import (
	"github.com/xiaolan580230/clUtil/clJson"
	"github.com/xiaolan580230/clhttp-framework/clResponse"
	"github.com/xiaolan580230/clhttp-framework/core/clAuth"
	"github.com/xiaolan580230/clhttp-framework/core/rule"
)


func ApiGetHomeStatic(_auth *clAuth.AuthInfo, _param *rule.HttpParam, _server *rule.ServerParam) string {

	return clResponse.Success(clJson.M{
		
	})
}