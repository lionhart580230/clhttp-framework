package main

import (
	"github.com/xiaolan580230/clhttp-framework/clGlobal"
	"github.com/xiaolan580230/clhttp-framework/core/httpserver"
	"github.com/xiaolan580230/clhttp-framework/core/skylog"
	"github.com/xiaolan580230/clhttp-framework/src/rule_list"
)

// HTTP服务默认使用端口号
const HTTPServerPort = 19999


func main() {

	clGlobal.Init("cl.conf")

	rule_list.Init()

	skylog.LogInfo( "正在启动服务，端口: %v", HTTPServerPort)
	httpserver.StartServer(HTTPServerPort)
}


