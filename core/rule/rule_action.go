package rule

import (
	"github.com/lionhart580230/clhttp-framework/clResponse"
	"github.com/lionhart580230/clhttp-framework/core/clAuth"
	"net/http"
	"sync"
	"time"
)

// 针对request进行前置动作的池子
var requestActionBefore map[string]func(*RequestBeforeParam) *RequestBeforeReturn

// 针对request进行后置动作的池子
var requestActionEnd map[string]func(*RequestAfterParam) *RequestAfterReturn

var globalActionBefore func(*RequestBeforeParam) *RequestBeforeReturn = nil
var globalActionAfter func(*RequestAfterParam) *RequestAfterReturn = nil

// 登录验证
var globalActionAuthCheck func(*AuthCheckParam) AuthCheckResp

var actionBeforeLocker sync.RWMutex
var actionAfterLocker sync.RWMutex

func init() {
	requestActionBefore = make(map[string]func(*RequestBeforeParam) *RequestBeforeReturn)
	requestActionEnd = make(map[string]func(*RequestAfterParam) *RequestAfterReturn)
}

// 针对指定request的前置动作参数
type RequestBeforeParam struct {
	Request    *http.Request    // 请求指针
	AcName     string           // 请求的acname，可以针对这个进行路由分发
	ServerInfo *ServerParam     // 服务器信息
	UserInfo   *clAuth.AuthInfo // 用户信息
	Param      *HttpParam       // 请求的参数
}

// 针对指定request的前置动作返回值
type RequestBeforeReturn struct {
	ServerInfo *ServerParam     // 服务器信息
	UserInfo   *clAuth.AuthInfo // 用户信息
	Param      *HttpParam       // 请求的参数
	RejectResp string           // 是否拦截
}

// 针对指定request的后置动作参数
type RequestAfterParam struct {
	Request        *http.Request        // 请求指针
	AcName         string               // 请求的acname，可以针对这个进行路由分发
	ServerInfo     *ServerParam         // 服务器信息
	UserInfo       *clAuth.AuthInfo     // 用户信息
	Param          *HttpParam           // 请求的参数
	ResponseText   string               // 返回内容
	ResponseWriter *http.ResponseWriter // 请求写入
}

// 针对指定request的后置动作参数
type AuthCheckParam struct {
	Request    *http.Request // 请求指针
	AcName     string        // 请求的acname，可以针对这个进行路由分发
	ServerInfo *ServerParam  // 服务器信息
	Param      *HttpParam    // 请求的参数
	Uid        uint64        // 用户uid
	Token      string        // 登录密钥
	SessionKey string        // sessionKey (客户端自己生成，生成后立即保存起来，用于识别设备）
}

// 针对指定request的后置动作参数
type AuthCheckResp struct {
	UserInfo     *clAuth.AuthInfo // 用户信息
	ResponseText string           // 返回内容
}

// 针对指定request的后置动作返回值
type RequestAfterReturn struct {
	ResponseText string // 返回内容
}

// 设置request前置条件
func SetRequestBeforeCallback(_request string, _func func(*RequestBeforeParam) *RequestBeforeReturn) {
	actionBeforeLocker.Lock()
	defer actionBeforeLocker.Unlock()

	requestActionBefore[_request] = _func
}

// 设置request后置条件
func SetRequestAfterCallback(_request string, _func func(param *RequestAfterParam) *RequestAfterReturn) {
	actionAfterLocker.Lock()
	defer actionAfterLocker.Unlock()

	requestActionEnd[_request] = _func
}

// 设置全局前置回调，当request没有设置特定的回调的时候才会执行它
func SetGlobalBeforeCallback(_func func(*RequestBeforeParam) *RequestBeforeReturn) {
	globalActionBefore = _func
}

/**
设置全局后置回调，当request没有设置特定回调的时候才会执行它
@param RequestAfterParam 后置回调需要的一些参数
@param RequestAfterReturn 后置回调返回处理过后的数据
*/

func SetGlobalAfterCallback(_func func(param *RequestAfterParam) *RequestAfterReturn) {
	globalActionAfter = _func
}

/**
设置验证登录回调，如果设置为nil（默认），那么就使用系统默认的验证行为进行验证
@param AuthCheckParam 里面包含了一些验证需要的字段
@return AuthCheckResp 返回内容，如果返回的用户对象为nil，那么就是验证失败
*/

func SetAuthCheckCallback(_func func(param *AuthCheckParam) AuthCheckResp) {
	globalActionAuthCheck = _func
}

// 激活前置条件
func DoRequestBefore(_request string, _param *RequestBeforeParam) *RequestBeforeReturn {
	actionBeforeLocker.RLock()
	_callback, exists := requestActionBefore[_request]
	actionBeforeLocker.RUnlock()
	if !exists {
		if globalActionBefore == nil {
			return &RequestBeforeReturn{
				ServerInfo: _param.ServerInfo,
				UserInfo:   _param.UserInfo,
				Param:      _param.Param,
			}
		}
		return globalActionBefore(_param)
	}
	return _callback(_param)
}

// 激活后置条件
func DoRequestAfter(_request string, _param *RequestAfterParam) *RequestAfterReturn {
	actionAfterLocker.RLock()
	_callback, exists := requestActionEnd[_request]
	actionAfterLocker.RUnlock()
	if !exists {
		if globalActionAfter == nil {
			return &RequestAfterReturn{
				ResponseText: _param.ResponseText,
			}
		}
		return globalActionAfter(_param)
	}
	return _callback(_param)
}

func DoAuthCheck(_rq *http.Request, _ac string, _serverParam *ServerParam, _param *HttpParam, _uid uint64, _token, _sessionKey string) (string, *clAuth.AuthInfo) {
	if globalActionAuthCheck != nil {
		resp := globalActionAuthCheck(&AuthCheckParam{
			Request:    _rq,
			AcName:     _ac,
			ServerInfo: _serverParam,
			Param:      _param,
			Uid:        _uid,
			Token:      _token,
		})
		return resp.ResponseText, resp.UserInfo
	}

	if _serverParam.JWT != "" {
		return CheckByJWT(_serverParam.JWT)
	}

	var authInfo *clAuth.AuthInfo
	if _uid > 0 && _token != "" {
		authInfo = clAuth.GetUser(_uid)
	}
	respStr := ""
	if authInfo == nil || !authInfo.IsLogin {
		// 失效了
		respStr = clResponse.NotLogin()
	} else if authInfo.Token != _token {
		// 被其他人顶出, 或者是失效
		if _sessionKey == "" {
			if authInfo.LastUptime > uint32(time.Now().Unix())-600 {
				respStr = clResponse.LogoutByKick()
			} else {
				respStr = clResponse.NotLogin()
			}
		} else {
			if authInfo.SessionKey != _sessionKey {
				respStr = clResponse.LogoutByKick()
			} else {
				respStr = clResponse.NotLogin()
			}
		}

	}
	return respStr, authInfo
}

// 通过JWT进行登录验证
func CheckByJWT(_jwtStr string) (string, *clAuth.AuthInfo) {
	authInfo := clAuth.CreateAuthByJWT(_jwtStr)
	return "", authInfo
}
