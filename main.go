package main

import (
	"github.com/xiaolan580230/clUtil/clLog"
	"github.com/xiaolan580230/clhttp-framework/clGlobal"
	"github.com/xiaolan580230/clhttp-framework/core/clAuth"
	"github.com/xiaolan580230/clhttp-framework/core/httpserver"
	"github.com/xiaolan580230/clhttp-framework/src/rule_list"
)

// HTTP服务默认使用端口号
const HTTPServerPort = 19999


func main() {

	clGlobal.Init("cl.conf")

	rule_list.Init()

	clLog.Info( "正在启动服务，端口: %v", HTTPServerPort)
	clAuth.SetGetUserByDB(func(_uid uint64) *clAuth.AuthInfo {
		return &clAuth.AuthInfo{
			Uid:        1,
			Token:      "1000",
			LastUptime: 0,
			IsLogin:    true,
			ExtraData:  nil,
		}
	})
	httpserver.SetUploadFileSizeLimit(1024 * 1024 * 300)
	httpserver.StartServer(HTTPServerPort)
}