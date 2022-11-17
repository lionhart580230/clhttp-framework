package rule

import (
	"github.com/xiaolan580230/clhttp-framework/core/clAuth"
	"net/http"
	"sync"
)

// 针对request进行前置动作的池子
var requestActionBefore map[string] func (*RequestBeforeParam) *RequestBeforeReturn
// 针对request进行后置动作的池子
var requestActionEnd map[string] func(*RequestAfterParam) *RequestAfterReturn

var actionBeforeLocker sync.RWMutex
var actionAfterLocker sync.RWMutex

func init() {
	requestActionBefore = make(map[string] func (*RequestBeforeParam) *RequestBeforeReturn)
	requestActionEnd = make(map[string] func(*RequestAfterParam) *RequestAfterReturn)
}

// 针对指定request的前置动作参数
type RequestBeforeParam struct {
	Request *http.Request			// 请求指针
	AcName string					// 请求的acname，可以针对这个进行路由分发
	ServerInfo *ServerParam			// 服务器信息
	UserInfo *clAuth.AuthInfo		// 用户信息
	Param *HttpParam				// 请求的参数
}
// 针对指定request的前置动作返回值
type RequestBeforeReturn struct {
	ServerInfo *ServerParam			// 服务器信息
	UserInfo *clAuth.AuthInfo		// 用户信息
	Param *HttpParam				// 请求的参数
}

// 针对指定request的后置动作参数
type RequestAfterParam struct {
	Request *http.Request			// 请求指针
	AcName string					// 请求的acname，可以针对这个进行路由分发
	ServerInfo *ServerParam			// 服务器信息
	UserInfo *clAuth.AuthInfo		// 用户信息
	Param *HttpParam				// 请求的参数
	ResponseText string					// 返回内容
	ResponseWriter *http.ResponseWriter	// 请求写入
}

// 针对指定request的后置动作返回值
type RequestAfterReturn struct {
	ResponseText string					// 返回内容
}


// 设置前置条件
func SetRequestBeforeCallback(_request string, _func func(*RequestBeforeParam) *RequestBeforeReturn) {
	actionBeforeLocker.Lock()
	defer actionBeforeLocker.Unlock()

	requestActionBefore[_request] = _func
}

// 设置后置条件
func SetRequestAfterCallback(_request string, _func func(param *RequestAfterParam) *RequestAfterReturn) {
	actionAfterLocker.Lock()
	defer actionAfterLocker.Unlock()

	requestActionEnd[_request] = _func
}


// 激活前置条件
func DoRequestBefore(_request string, _param *RequestBeforeParam) *RequestBeforeReturn {
	actionBeforeLocker.RLock()
	_callback, exists := requestActionBefore[_request]
	actionBeforeLocker.RUnlock()
	if !exists {
		return &RequestBeforeReturn{
			ServerInfo: _param.ServerInfo,
			UserInfo:   _param.UserInfo,
			Param:      _param.Param,
		}
	}
	return _callback(_param)
}


// 激活后置条件
func DoRequestAfter(_request string, _param *RequestAfterParam) *RequestAfterReturn {
	actionAfterLocker.RLock()
	_callback, exists := requestActionEnd[_request]
	actionAfterLocker.RUnlock()
	if !exists {
		return &RequestAfterReturn{
			ResponseText: _param.ResponseText,
		}
	}
	return _callback(_param)
}