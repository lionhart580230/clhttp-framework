package clAuth

import (
	"encoding/json"
	"github.com/xiaolan580230/clhttp-framework/clCrypt"
	"github.com/xiaolan580230/clhttp-framework/clGlobal"
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

	if clGlobal.SkyConf.IsCluster {
		SaveUser(_auth)
		return
	}

	mAuthMap[_auth.Token] = _auth
}


// 保存用户信息到数据库
func SaveUser(_auth *AuthInfo) {
	if !clGlobal.SkyConf.IsCluster {
		return
	}

	redis := clGlobal.GetRedis()
	if redis != nil {
		var userData, err = json.Marshal(_auth)
		if err != nil {
			return
		}
		redis.Set(_auth.Token, clCrypt.Base64Encode(string(userData)), 3600)
	}
}


// 移除用户
func DelUser( _auth *AuthInfo) {
	mLocker.Lock()
	defer mLocker.Unlock()

	if clGlobal.SkyConf.IsCluster {
		redis := clGlobal.GetRedis()
		if redis != nil {
			redis.Del(_auth.Token)
		}
		return
	}


	delete(mAuthMap, _auth.Token )
}


// 获取用户
func GetUser( _token string ) *AuthInfo {
	mLocker.RLock()
	defer mLocker.RUnlock()

	if clGlobal.SkyConf.IsCluster {
		redis := clGlobal.GetRedis()
		if redis != nil {
			var userCache = redis.Get(_token)
			if userCache == "" {
				return nil
			}

			var userObj AuthInfo
			err := json.Unmarshal(clCrypt.Base64Decode(userCache), &userObj)
			if err != nil {
				return nil
			}

			return &userObj
		}
		return nil
	}

	return mAuthMap[_token]
}


// 设置信息
func (this *AuthInfo) SetItem(_key string, _val string) {
	mLocker.Lock()
	defer mLocker.Unlock()

	this.extraData[_key] = _val

	SaveUser(this)
}


// 获取信息
func (this *AuthInfo) GetItem(_key string) string {
	mLocker.RLock()
	defer mLocker.RUnlock()

	return this.extraData[_key]
}


// 设置登录
// 如果设置为登录中状态则 uid必须>0
// 如果没有则自动切换为离线状态
func (this *AuthInfo) SetLogin(_uid uint64, _token string) {
	if _uid > 0 && _token != "" {
		this.IsLogin = true
		this.Uid = _uid
		this.Token = _token
		AddUser(this)
	} else {
		DelUser(this)
	}
}


// 是否登录
func (this *AuthInfo) CheckLogin() bool {
	return this.IsLogin
}