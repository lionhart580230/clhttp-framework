package main

import (
	"clhttp-framework/clGlobal"
	"clhttp-framework/core/httpserver"
	"clhttp-framework/core/skylog"
	"clhttp-framework/src/rule_list"
)

// HTTP服务默认使用端口号
const HTTPServerPort = 80


func main() {

	clGlobal.Init("cl.conf")
	rule_list.Init()

	skylog.LogInfo( "正在启动服务，端口: %v", HTTPServerPort)
	httpserver.StartServer(HTTPServerPort)
}


