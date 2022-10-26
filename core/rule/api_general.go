package rule

import (
	"encoding/json"
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

%v
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
		writeBuffer := fmt.Sprintf(ApiTemp, _package, "", apiName)
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

type GeneralParams struct {
	Name string `json:"name"`
	Strict bool `json:"strict"`
	Type string `json:"type"`
	Comment string `json:"comment"`
}
type GeneralRouterInfo struct {
	Ac string `json:"ac"`
	Comment string `json:"comment"`
	Params []GeneralParams `json:"params"`
	Login bool `json:"login"`
}

// api文件根据json
//@param _pathName 要生成在哪个目录
func ApiGeneralByJson(_jsonStr, _pathName string, _package, _request string, _file string) {

	var data = make([]GeneralRouterInfo, 0)
	err := json.Unmarshal([]byte(_jsonStr), &data)
	if err != nil {
		clLog.Info("生成API失败! 解析目标JSON出错: %v", err)
		return
	}

	var tempFunc = `func Init() {

	// 添加请求类型
	rule.AddRequest("` + _request + `", "ac")

// 接口列表
%v

}`

	var tempItem = `
	%v
	rule.AddRule(rule.Rule{
		Request:     "%v",
		Name:        "%v",
		Params:      []rule.ParamInfo{` +"\n" + `%v		},
		CallBack:    %v,
		CacheExpire: 0,
		CacheType:   0,
		Login:       %v,
	})
	` + "\n\n"
	var tempParam = `			rule.NewParam("%v", rule.%v, %v, ""),` + "\n"

	// 生成路由表缓冲区
	var ruleListContent = strings.Builder{}

	ruleListContent.WriteString(`package rulelist

import "github.com/xiaolan580230/clhttp-framework/core/rule"

`)

	clFile.CreateDirIFNotExists(_pathName)

	var ruleBodyBuilder = strings.Builder{}

	for _, val := range data {

		apiName := "Api" + strings.ToUpper(string(val.Ac[0])) + val.Ac[1:]
		fileName := _pathName + "/" + apiName + ".go"
		comment := ""
		if val.Comment != "" {
			comment = "// " + val.Comment
		}
		for _, p := range val.Params {
			comment += "\n// " + p.Name + " " + p.Comment
		}
		writeBuffer := fmt.Sprintf(ApiTemp, _package, comment, apiName)
		if clFile.FileIsExists(fileName) {
			clLog.Info("文件: %v 已经存在, 跳过!", fileName)
			continue
		}

		var paramsBuilder = strings.Builder{}
		for _, p := range val.Params {
			paramsBuilder.WriteString(fmt.Sprintf(tempParam, p.Name, p.Type, p.Strict))
		}

		comment = ""
		if val.Comment != "" {
			comment = "// " + val.Comment
		}

		ruleBodyBuilder.WriteString(fmt.Sprintf(tempItem, comment, _request, val.Ac, paramsBuilder.String(), _package + "." + apiName, val.Login))

		if err := ioutil.WriteFile(fileName, []byte(writeBuffer), 0666); err != nil {
			clLog.Error("生成Api文件: %v 失败: %v", fileName, err)
		} else {
			clLog.Info("生成文件: %v 成功!", fileName)
		}
	}

	ruleListContent.WriteString(fmt.Sprintf(tempFunc, ruleBodyBuilder.String()))
	clFile.CreateDirIFNotExists("rulelist")
	clFile.AppendFile("rulelist/" + _file + ".go", ruleListContent.String())
	clLog.Info("文件内容:\n\n %+v\n", ruleListContent.String())
}