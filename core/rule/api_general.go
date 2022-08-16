package rule

import (
	"fmt"
	"github.com/xiaolan580230/clUtil/clFile"
	"github.com/xiaolan580230/clUtil/clLog"
	"io/ioutil"
	"strings"
)


const ApiTemp = `package %v

import (
	"github.com/xiaolan580230/clUtil/clJson"
	"github.com/xiaolan580230/clUtil/clLog"
	"github.com/xiaolan580230/clhttp-framework/clResponse"
	"github.com/xiaolan580230/clhttp-framework/core/clAuth"
	"github.com/xiaolan580230/clhttp-framework/core/rule"
	"strings"
)


func %v(_auth *clAuth.AuthInfo, _param *rule.HttpParam, _server *rule.ServerParam) string {

	return clResponse.Success(clJson.M{
		
	})
}`


// api文件生成器
//@param _pathName 要生成在哪个目录
func ApiGeneral(_pathName string, _package, _request string) {

	for _, val := range ruleList {
		if val.Request != _request {
			continue
		}

		apiName := "Api" + strings.ToUpper(string(val.Name[0])) + val.Name[1:]
		fileName := _pathName + "/" + apiName + ".go"
		writeBuffer := fmt.Sprintf(ApiTemp, _package, apiName)
		if clFile.FileIsExists(fileName) {
			clLog.Info("文件: %v 已经存在, 跳过!", fileName)
			continue
		}

		if err := ioutil.WriteFile(fileName, []byte(writeBuffer), 0666); err != nil {
			clLog.Error("生成Api文件: %v 失败: %v", fileName, err)
		} else {
			clLog.Info("生成文件: %v 成功!", fileName)
		}
	}
}