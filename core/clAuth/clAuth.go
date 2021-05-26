package clAuth

import (
	"sync"
	"time"
)

// 验证类。当用户登录后需要将用户信息保存在用户池中

type AuthInfo struct {
	Uid uint64					// 当前用户Id
	Token string				// 用户Token数据
	LastUptime uint32			// 最近活跃时间
	IsLogin bool				// 是否登录
	extraData map[string]string	// 附属数据
	mLocker sync.RWMutex		// 异步锁
}

var mAuthMap map[string] *AuthInfo
var mLocker sync.RWMutex


func init() {
	mAuthMap = make(map[string] *AuthInfo)
}


// 获取新用户
func NewUser(_uid uint64, _token string ) *AuthInfo {
	return &AuthInfo{
		Uid:        _uid,
		Token:      _token,
		LastUptime: 0,
		IsLogin:    false,
		extraData:  make(map[string] string),
	}
}


// 添加用户
func AddUser( _auth *AuthInfo) {
	mLocker.Lock()
	defer mLocker.Unlock()

	_auth.LastUptime = uint32(time.Now().Unix())

	mAuthMap[_auth.Token] = _auth
}


// 移除用户
func DelUser( _auth *AuthInfo) {
	mLocker.Lock()
	defer mLocker.Unlock()

	delete(mAuthMap, _auth.Token )
}


// 获取用户
func GetUser( _token string ) *AuthInfo {
	mLocker.RLock()
	defer mLocker.RUnlock()

	return mAuthMap[_token]
}


// 设置信息
func (this *AuthInfo) SetItem(_key string, _val string) {
	mLocker.Lock()
	defer mLocker.Unlock()

	this.extraData[_key] = _val
}


// 获取信息
func (this *AuthInfo) GetItem(_key string) string {
	mLocker.RLock()
	defer mLocker.RUnlock()

	return this.extraData[_key]
}


// 设置登录
func (this *AuthInfo) SetLogin(_login bool) {
	this.IsLogin = _login
}


// 是否登录
func (this *AuthInfo) CheckLogin() bool {
	return this.IsLogin
}