package clAuth

import (
	"encoding/json"
	"github.com/xiaolan580230/clhttp-framework/clCrypt"
	"github.com/xiaolan580230/clhttp-framework/clGlobal"
	"github.com/xiaolan580230/clhttp-framework/core/skylog"
	"github.com/xiaolan580230/clhttp-framework/core/skyredis"
	"sync"
	"time"
)

// 验证类。当用户登录后需要将用户信息保存在用户池中

type AuthInfo struct {
	Uid uint64					// 当前用户Id
	Token string				// 用户Token数据
	LastUptime uint32			// 最近活跃时间
	IsLogin bool				// 是否登录
	ExtraData map[string]string	// 附属数据
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
		ExtraData:  make(map[string] string),
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


// 加载所有用户数据
func LoadUsers() {
	mLocker.Lock()
	defer mLocker.Unlock()

	mAuthMap = make(map[string] *AuthInfo)

	redis := clGlobal.GetRedis()

	keys := redis.GetKeys("U_INFO_*")
	for _, ukey := range keys {
		var data AuthInfo
		var jsonStr = clCrypt.Base64Decode(redis.Get(ukey))
		var err = json.Unmarshal(jsonStr, &data)
		if err != nil {
			skylog.LogErr("加载: %v 用户数据失败! 错误: %v -> %v", err, string(jsonStr))
			continue
		}
		skylog.LogInfo("成功加载用户: %v -> %v", data.Token, data.Uid)
		mAuthMap[ data.Token ] = &data
	}
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
		redis.Set("U_INFO_" + _auth.Token, clCrypt.Base64Encode(string(userData)), 3600)
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
	delete(mAuthMap, _auth.Token)
}


// 获取用户
func GetUser( _token string ) *AuthInfo {

	if clGlobal.SkyConf.IsCluster {
		redis := clGlobal.GetRedis()
		if redis != nil {
			var userCache = redis.Get("U_INFO_" + _token)
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

	mLocker.RLock()
	defer mLocker.RUnlock()
	return mAuthMap[_token]
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

		// 清理内存
		ClearUserData(_uid, _token)
	} else {
		DelUser(this)
	}
}


// 清理内存 (如果内存中或redis中存在这个uid的数据，则自动清理)
func ClearUserData(_uid uint64, _ignToken string) {
	mLocker.Lock()
	defer mLocker.Unlock()

	var redis *skyredis.RedisObject
	if clGlobal.SkyConf.IsCluster {
		redis = clGlobal.GetRedis()
	}
	for _, val := range mAuthMap {
		if val.Uid == _uid && val.Token != _ignToken {
			delete(mAuthMap, val.Token)
			if redis != nil {
				redis.Del("U_INFO_" + val.Token)
			}
		}
	}
}


// 是否登录
func (this *AuthInfo) CheckLogin() bool {
	return this.IsLogin
}