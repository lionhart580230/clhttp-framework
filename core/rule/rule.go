package rule

import (
	"clhttp-framework/clCommon"
	"clhttp-framework/core/clAuth"
	"clhttp-framework/core/clCache"
	"clhttp-framework/core/skylog"
	"clhttp-framework/src/skylang"
	"net/url"
	"strings"
	"sync"
)



type ServerParam struct {
	RemoteIP   string
	RequestURI string
	Host       string
	Method     string
	RequestURL string
	UA         string
	UAType     uint32
	Proctol    string
	Port       string
	Language   string
}


//@author xiaolan
//@lastUpdate 2019-08-10
//@comment 路由规则定义

// 路由结构  2019-08-10
type Rule struct {
	Request string	   // 请求的名字
	Name   string      // 规则名称
	Params []ParamInfo // 参数列表
	// 回调函数
	CallBack func(_auth *clAuth.AuthInfo, _param *HttpParam, _server *ServerParam) string
	CacheExpire int		// 缓存秒数, 负数为频率控制, 正数为缓存时间, 0为不缓存
	CacheType int		// 缓存启动的时候才有效 0=全局缓存,1=根据IP缓存,2=根据用户缓存
	Login bool			// 是否登录才可以访问这个接口
	Method string		// 请求方法, 为空则不限制请求方法, POST则为只允许POST请求
}

// 路由列表
var ruleList map[string]Rule
var ruleLocker sync.RWMutex

func init() {
	ruleList = make(map[string]Rule)
}

//@author xiaolan
//@lastUpdate 2019-08-10
//@comment 添加规则
//@param _rule 规则结构体
func AddRule(_rule Rule) {
	ruleLocker.Lock()
	defer ruleLocker.Unlock()

	ruleList[_rule.Name] = _rule
}


// @auth xiaolan
// @lastUpdate 2021-05-26
// @comment 构建缓存key
func BuildCacheKey(_params []string) string {

	var keys = strings.Join(_params, "_")
	return clCommon.Md5([]byte(keys))
}



//@author xiaolan
//@lastUpdate 2019-08-10
//@comment 调用规则
func CallRule(_uri string, _param *HttpParam, _server *ServerParam) string {
	ruleLocker.RLock()
	defer ruleLocker.RUnlock()

	paramStr := _param.GetStr("p", "")
	decodeStr := paramStr

	paramList, err := url.ParseQuery(string(decodeStr))
	if err != nil {
		skylog.LogErr( "HttpRequest 反序列化失败!!错误:%v", err)
		return clCommon.JCode(skylang.MSG_ERR_FAILED_INT, "模块不存在!", nil)
	}

	// 通过AC获取到指定的路由
	acName := paramList.Get("ac")
	ruleinfo, exists := ruleList[_uri + "_" + acName]
	if !exists {
		skylog.LogErr( "AC <%v> 不存在!", acName)
		return clCommon.JCode(skylang.MSG_ERR_FAILED_INT, "模块不存在!", nil)
	}

	var authInfo *clAuth.AuthInfo

	// 需要登录
	if ruleinfo.Login {
		var token = paramList.Get("token")
		if token == "" {
			return clCommon.JCode(skylang.MSG_ERR_FAILED_INT, "请先登录!", nil)
		}

		authInfo = clAuth.GetUser(token)
		if authInfo == nil || !authInfo.IsLogin {
			return clCommon.JCode(skylang.MSG_ERR_FAILED_INT, "请先登录!", nil)
		}
	}

	// 检查参数
	newParam := NewHttpParam(nil)
	paramsKeys := make([]string, 0)
	if ruleinfo.Params != nil {
		for _, pinfo := range ruleinfo.Params {
			value := paramList.Get(pinfo.Name)
			if value == PARAM_CHECK_FAIED || value == "" {
				if pinfo.Static {
					// 严格模式
					return clCommon.JCode(skylang.MSG_ERR_FAILED_INT, "参数:" + pinfo.Name + "不合法!", pinfo.Name)
				} else {
					value = pinfo.Default
				}
			} else {
				if !pinfo.CheckParam(value) {
					if pinfo.Static {
						// 严格模式
						return clCommon.JCode(skylang.MSG_ERR_FAILED_INT, "参数:" + pinfo.Name + "不合法!", pinfo.Name)
					} else {
						value = pinfo.Default
					}
				}
			}
			newParam.Add(pinfo.Name, value)
			paramsKeys = append(paramsKeys, pinfo.Name + "=" + value)
		}
	}

	// 如果回调函数不存在
	if ruleinfo.CallBack == nil {
		skylog.LogErr("AC[%v]回调函数为空!", acName)
		return clCommon.JCode(skylang.MSG_ERR_FAILED_INT, "模块不存在!", nil)
	}

	// 检查是否需要缓存
	var cacheKey = ""
	if ruleinfo.CacheExpire > 0 {
		// 根据用户缓存
		if ruleinfo.CacheType == 2 {
			paramsKeys = append(paramsKeys, "token=" + paramList.Get("token"))
		} else if ruleinfo.CacheType == 1 {
			// 根据IP缓存
			paramsKeys = append(paramsKeys, "ip=" + _server.RemoteIP)
		}
		cacheKey = BuildCacheKey(paramsKeys)
		cacheStr := clCache.GetCache(cacheKey)
		if cacheStr != "" {
			return cacheStr
		}
	}

	// 调用回调函数，并返回结果
	respStr := ruleinfo.CallBack(authInfo, newParam, _server)

	// 检查是否需要缓存
	if ruleinfo.CacheExpire > 0 {
		clCache.UpdateCache(cacheKey, respStr, uint32(ruleinfo.CacheExpire))
	}

	skylog.LogDebug("[ACK][%s] %s", _server.RemoteIP, respStr)

	return respStr
}
