package main

import (
	"github.com/xiaolan580230/clUtil/clLog"
	"github.com/xiaolan580230/clhttp-framework/clGlobal"
	"github.com/xiaolan580230/clhttp-framework/core/clAuth"
	"github.com/xiaolan580230/clhttp-framework/core/httpserver"
	"github.com/xiaolan580230/clhttp-framework/core/modelCreator"
	"github.com/xiaolan580230/clhttp-framework/core/rule"
	"github.com/xiaolan580230/clhttp-framework/src/rule_list"
)

// HTTP服务默认使用端口号
const HTTPServerPort = 19999


func main() {

	clGlobal.Init("cl.conf")

	rule_list.Init()

	clAuth.SetAuthPrefix("U_INFO")

	httpserver.SetAESKey("5d41402abc4b2a76b9719d911017c592")
	// 关闭上传功能
	httpserver.SetEnableUploadFile(false)
	// 关闭上传调试页
	httpserver.SetEnableUploadTest(false)

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

	// 根据路由配置表生成api文档
	rule.ApiGeneral("./apis", "apis", "/request")

	// 根据数据库中的配置生成模型
	modelCreator.CreateAllModelFile("127.0.0.1", "root", "root", "testdb", "testModel")

	httpserver.StartServer(HTTPServerPort)
}