package apis

import (
	"github.com/lionhart580230/clUtil/clJson"
	"github.com/lionhart580230/clhttp-framework/clResponse"
	"github.com/lionhart580230/clhttp-framework/core/clAuth"
	"github.com/lionhart580230/clhttp-framework/core/rule"
)

func ApiGetMenu(_auth *clAuth.AuthInfo, _param *rule.HttpParam, _server *rule.ServerParam) string {

	return clResponse.Success(clJson.M{})
}
