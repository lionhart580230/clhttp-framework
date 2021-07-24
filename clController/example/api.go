package example

import (
	"github.com/xiaolan580230/clhttp-framework/clResponse"
	"github.com/xiaolan580230/clhttp-framework/core/clAuth"
	"github.com/xiaolan580230/clhttp-framework/core/cljson"
	"github.com/xiaolan580230/clhttp-framework/core/rule"
)

func ApiExample(_auth *clAuth.AuthInfo, _param *rule.HttpParam, _server *rule.ServerParam) string {

	return clResponse.Success(_server)
}


// 上传范例
func ApiUploadExample(_auth *clAuth.AuthInfo, _param *rule.HttpParam, _server *rule.ServerParam) string {

	return clResponse.Success(cljson.M{
		"filename": _param.GetStr("filename", ""),
		"fileExt": _param.GetStr("fileExt", ""),
		"localPath": _param.GetStr("localPath", ""),
	})
}