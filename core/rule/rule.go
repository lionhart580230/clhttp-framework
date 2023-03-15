package rule

import (
	"fmt"
	"github.com/xiaolan580230/clUtil/clCrypt"
	"github.com/xiaolan580230/clUtil/clLog"
	"github.com/xiaolan580230/clhttp-framework/clCommon"
	"github.com/xiaolan580230/clhttp-framework/clGlobal"
	"github.com/xiaolan580230/clhttp-framework/clResponse"
	"github.com/xiaolan580230/clhttp-framework/core/clAuth"
	"github.com/xiaolan580230/clhttp-framework/core/clCache"
	"github.com/xiaolan580230/clhttp-framework/src/skylang"
	"net/http"
	"strings"
	"sync"
	"time"
)



type ServerParam struct {
	RemoteIP   string			// 远程IP地址
	RequestURI string			// 请求URI
	UriData   *HttpParam		// uri上的参数列表
	Host       string			// 请求域名
	Method     string			// 请求方法
	Header     http.Header
	RequestURL string			// 请求完整地址
	UA         string			// 目标设备信息
	UAType     uint32			// 目标设备类型
	Proctol    string			// 目标协议
	Port       string			// 端口
	Language   string			// 使用语言信息
	LangType   uint32			// 使用语言信息
	ContentType string			// 提交的方式
	RawData string				// 原始数据
	RawParam *HttpParam			// 原始的参数
	Encrypt bool				// 是否需要加密
	AesKey string				// 加密用的key
	Iv string 					// 加密用的iv
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
	RespContent string  // 返回的结构体内容格式 默认是 text/json
}

// 路由列表
var ruleList map[string] Rule
var ruleLocker sync.RWMutex


// 请求方式
var requestList map[string] string

func init() {
	ruleList = make(map[string] Rule)
	requestList = make(map[string] string)
}


// 添加新的请求方式
func AddRequest(_request string, _acKey string) {
	requestList[_request] = _acKey
}


// 获取请求方式的ackey
func GetRequestAcKey(_request string) string {
	ackey, exists := requestList[_request]
	if !exists {
		return "ac"
	}
	return ackey
}


//@author xiaolan
//@lastUpdate 2019-08-10
//@comment 添加规则
//@param _rule 规则结构体
func AddRule(_rule Rule) {
	ruleLocker.Lock()
	defer ruleLocker.Unlock()

	ruleList[_rule.Request + "_" + _rule.Name] = _rule
}


// @auth xiaolan
// @lastUpdate 2021-05-26
// @comment 构建缓存key
func BuildCacheKey(_params []string) string {

	var keys = strings.Join(_params, "_")
	return clCommon.Md5([]byte(keys))
}



// 删除Api缓存
func DelApiCache(_uri string, _acName string, _uid uint64, _params map[string]string) {
	ruleinfo, exists := ruleList[_uri + "_" + _acName]
	if !exists {
		clLog.Error( "删除缓存失败! AC <%v_%v> 不存在!", _uri, _acName)
		return
	}
	paramsKeys := make([]string, 0)
	if ruleinfo.CacheType == 2 {
		paramsKeys = append(paramsKeys, fmt.Sprintf("uid=%v", _uid))
	}
	if ruleinfo.Params != nil {
		for _, pinfo := range ruleinfo.Params {
			value := _params[pinfo.Name]
			paramsKeys = append(paramsKeys, pinfo.Name + "=" + value)
		}
	}
	cacheKey := BuildCacheKey(paramsKeys)
	clCache.DelCache(cacheKey)
}


// 删除全部Api缓存
func DelApiCacheAll(_uri string, _acName string) {
	clCache.DelCacheContains(_uri + "_" + _acName + "_")
}


//@author xiaolan
//@lastUpdate 2019-08-10
//@comment 调用规则
func CallRule(rq *http.Request, rw *http.ResponseWriter, _uri string, _param *HttpParam, _server *ServerParam) (string, string) {
	ruleLocker.RLock()
	defer ruleLocker.RUnlock()

	var acKey = GetRequestAcKey(_uri)

	// 通过AC获取到指定的路由
	acName := _param.GetStr(acKey, "")
	ruleinfo, exists := ruleList[_uri + "_" + acName]
	if !exists {
		if clGlobal.SkyConf.DebugRouter {
			clLog.Error( "AC <%v_%v> 不存在! IP: %v", _uri, acName, _server.RemoteIP)
			clLog.Debug("%+v", ruleList)
		}
		respStr := clResponse.JCode(skylang.MSG_ERR_FAILED_INT, "模块不存在!", nil)
		if _server.Encrypt {
			respStr = clCrypt.AesCBCEncode(respStr, _server.AesKey, _server.Iv)
		}
		return respStr, "text/json"
	}

	if ruleinfo.RespContent == "" {
		ruleinfo.RespContent = "text/json"
	}

	var authInfo *clAuth.AuthInfo
	var token = _param.GetStr("token", "")
	var uid  = _param.GetUint64("uid", 0)
	var sessionKey = _param.GetStr("session_key", "")
	if uid > 0 && token != "" {
		authInfo = clAuth.GetUser(uid)
	}

	paramsKeys := make([]string, 0)

	// 需要登录
	if ruleinfo.Login {
		respStr := ""
		if authInfo == nil || !authInfo.IsLogin {
			// 失效了
			respStr = clResponse.NotLogin()
		} else if authInfo.Token != token {

			// 被其他人顶出, 或者是失效
			if sessionKey == "" {
				if authInfo.LastUptime > uint32(time.Now().Unix()) - 600 {
					respStr = clResponse.LogoutByKick()
				} else {
					respStr = clResponse.NotLogin()
				}
			} else {
				if authInfo.SessionKey != sessionKey {
					respStr = clResponse.LogoutByKick()
				} else {
					respStr = clResponse.NotLogin()
				}
			}


		}
		if respStr != "" {
			if _server.Encrypt {
				respStr = clCrypt.AesCBCEncode(respStr, _server.AesKey, _server.Iv)
			}
			return respStr, ruleinfo.RespContent
		}

		//if authInfo == nil || authInfo.Token != token || !authInfo.IsLogin {
		//	if clGlobal.SkyConf.DebugRouter {
		//		clLog.Debug("Uid: %v TOKEN: %v 登录状态无效!", uid, token)
		//	}
		//	clLog.Info("用户: [%v] %v 登录状态失效!", uid, token)
		//
		//
		//
		//
		//}
	} else {
		authInfo = clAuth.NewUser(0, "")
	}

	// 检查参数
	newParam := NewHttpParam(nil)
	if ruleinfo.Params != nil {
		for _, pinfo := range ruleinfo.Params {
			value := _param.GetStr(pinfo.Name, "")
			if value == PARAM_CHECK_FAIED || value == "" {
				if pinfo.Static {
					// 严格模式
					return clResponse.JCode(skylang.MSG_ERR_FAILED_INT, "参数:" + pinfo.Name + "不合法!", pinfo.Name), ruleinfo.RespContent
				} else {
					value = pinfo.Default
				}
			} else {
				if !pinfo.CheckParam(value) {
					if pinfo.Static {
						// 严格模式
						return clResponse.JCode(skylang.MSG_ERR_FAILED_INT, "参数:" + pinfo.Name + "不合法!", pinfo.Name), ruleinfo.RespContent
					} else {
						value = pinfo.Default
					}
				}
			}
			newParam.Add(pinfo.Name, value)
			paramsKeys = append(paramsKeys, pinfo.Name + "=" + value)
		}
	} else {
		// 如果路由配置上参数列表为nil，那么就不过滤参数，所有参数都接收
		for key, val := range _param.values {
			newParam.Add(key, val)
		}
	}

	// 如果回调函数不存在
	if ruleinfo.CallBack == nil {
		if ruleinfo.RespContent != "" {
			return ruleinfo.RespContent, ruleinfo.RespContent
		}

		if clGlobal.SkyConf.DebugRouter {
			clLog.Error("AC[%v]回调函数为空!", acName)
		}
		respStr := clResponse.JCode(skylang.MSG_ERR_FAILED_INT, "模块不存在!", nil)
		if _server.Encrypt {
			respStr = clCrypt.AesCBCEncode(respStr, _server.AesKey, _server.Iv)
		}
		return respStr, "text/json"
	}

	// 检查是否需要缓存
	var cacheKey = ""
	if ruleinfo.CacheExpire > 0 {
		// 根据用户缓存
		if ruleinfo.CacheType == 2 {
			paramsKeys = append(paramsKeys, fmt.Sprintf("uid=%v", authInfo.Uid))
		} else if ruleinfo.CacheType == 1 {
			// 根据IP缓存
			paramsKeys = append(paramsKeys, "ip=" + _server.RemoteIP)
		}
		cacheKey = _uri + "_" + acName + "_" + BuildCacheKey(paramsKeys)
		if _server.Encrypt {	// 如果是加密的话，需要带上Iv
			cacheKey += _server.Iv
		}
		cacheStr := clCache.GetCache(cacheKey)
		if cacheStr != "" {
			return cacheStr, ruleinfo.RespContent
		}
	}

	// 调用前置函数，并返回结果
	var beforeParam = DoRequestBefore(_uri, &RequestBeforeParam{
		Request:    rq,
		AcName:     acName,
		ServerInfo: _server,
		UserInfo:   authInfo,
		Param:      _param,
	})

	nowTime := time.Now()
	respStr := ruleinfo.CallBack(beforeParam.UserInfo, beforeParam.Param, beforeParam.ServerInfo)
	diffTime := time.Since(nowTime).Seconds()
	if diffTime > 5 {
		clLog.Error("接口:%v.%v 处理耗时(%0.2fs)过长!", _uri, acName, diffTime)
	}

	afterResp := DoRequestAfter(_uri, &RequestAfterParam{
		Request:    rq,
		AcName:     acName,
		ServerInfo: _server,
		UserInfo:   authInfo,
		Param:      _param,
		ResponseText:   respStr,
		ResponseWriter: rw,
	})

	respStr = afterResp.ResponseText
	// 需要加密
	if _server.Encrypt {
		respStr = clCrypt.AesCBCEncode(respStr, _server.AesKey, _server.Iv)
	}

	// 检查是否需要缓存
	if ruleinfo.CacheExpire > 0 {
		clCache.UpdateCacheSimple(cacheKey, respStr, uint32(ruleinfo.CacheExpire))
	}

	if clGlobal.SkyConf.DebugRouter {
		clLog.Debug("[%s][%s] REQUEST: %s / RESPONSE: %s", acName, _server.RemoteIP, strings.Join(paramsKeys, "&"), respStr)
	}

	return respStr, ruleinfo.RespContent
}
