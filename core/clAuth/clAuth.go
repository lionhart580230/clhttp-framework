package clAuth

import (
	"encoding/json"
	"fmt"
	"github.com/lionhart580230/clUtil/clCrypt"
	"github.com/lionhart580230/clUtil/clLog"
	"github.com/lionhart580230/clhttp-framework/clGlobal"
	"sync"
	"time"
)

// 验证类。当用户登录后需要将用户信息保存在用户池中

type AuthInfo struct {
	Uid        uint64            // 当前用户Id
	Token      string            // 用户Token数据
	LastUptime uint32            // 最近活跃时间
	IsLogin    bool              // 是否登录
	ExtraData  map[string]string // 附属数据
	SessionKey string            // sessionkey
	mLocker    sync.RWMutex      // 异步锁
}

var mAuthMap map[uint64]*AuthInfo
var mLocker sync.RWMutex
var prefix = "U_INFO_"
var mGetUserByDB func(_uid uint64) *AuthInfo

func init() {
	mAuthMap = make(map[uint64]*AuthInfo)
}

// 设置用户登录缓存前缀
func SetAuthPrefix(_prefix string) {
	prefix = _prefix
}

// 设置通过数据库获取用户结构的回调
func SetGetUserByDB(_func func(_uid uint64) *AuthInfo) {
	mGetUserByDB = _func
}

// 获取新用户
func NewUser(_uid uint64, _token string) *AuthInfo {
	return &AuthInfo{
		Uid:        _uid,
		Token:      _token,
		LastUptime: 0,
		IsLogin:    false,
		ExtraData:  make(map[string]string),
	}
}

// 获取用户缓存Key
func GetUserKey(_uid uint64) string {
	return fmt.Sprintf("%v%v", prefix, _uid)
}

// 添加用户
func AddUser(_auth *AuthInfo) {
	mLocker.Lock()
	defer mLocker.Unlock()

	_auth.LastUptime = uint32(time.Now().Unix())

	if clGlobal.SkyConf.IsCluster {
		SaveUser(_auth)
		return
	}

	mAuthMap[_auth.Uid] = _auth
}

// 保存用户信息到数据库
func SaveUser(_auth *AuthInfo) {
	if _auth.Uid == 0 {
		return
	}

	if !clGlobal.SkyConf.IsCluster {
		return
	}

	redis := clGlobal.GetRedis()
	if redis != nil {
		var userData, err = json.Marshal(_auth)
		if err != nil {
			clLog.Error("序列化用户缓存错误: %v", err)
			return
		}
		redis.SetEx(GetUserKey(_auth.Uid), clCrypt.Base64Encode(userData), 12*3600)
	}
}

// 移除用户
func DelUser(_auth *AuthInfo) {
	mLocker.Lock()
	defer mLocker.Unlock()

	if clGlobal.SkyConf.IsCluster {
		redis := clGlobal.GetRedis()
		if redis != nil {
			redis.Del(GetUserKey(_auth.Uid))
		}
		return
	}
	delete(mAuthMap, _auth.Uid)
}

// 获取用户
func GetUser(_uid uint64) *AuthInfo {
	if _uid == 0 {
		return nil
	}
	var userObj = &AuthInfo{}
	if clGlobal.SkyConf.IsCluster {
		redis := clGlobal.GetRedis()
		if redis != nil {
			var userKey = GetUserKey(_uid)
			var userCache = redis.Get(userKey)
			if userCache != "" {
				err := json.Unmarshal(clCrypt.Base64Decode(userCache), userObj)
				if err != nil {
					clLog.Error("获取反序列化用户缓存错误: %v -> %v", userCache, err)
				} else {
					redis.SetExpire(userKey, 3600)
				}
			}
		}
	} else {
		mLocker.RLock()
		userObj = mAuthMap[_uid]
		mLocker.RUnlock()
	}

	if userObj == nil && mGetUserByDB != nil {
		userObj = mGetUserByDB(_uid)
	}

	if userObj != nil && userObj.Uid != _uid {
		return nil
	}

	return userObj
}

// 设置登录
// 如果设置为登录中状态则 uid必须>0
// 如果没有则自动切换为离线状态
func (this *AuthInfo) SetLogin(_uid uint64, _token string) {
	if this == nil && _uid > 0 && _token != "" {
		this = NewUser(_uid, _token)
	}
	if _uid > 0 && _token != "" {
		this.IsLogin = true
		this.Uid = _uid
		this.Token = _token
		AddUser(this)
	} else if this != nil {
		DelUser(this)
	}
}

// 设置登录
// 如果设置为登录中状态则 uid必须>0
// 如果没有则自动切换为离线状态
func (this *AuthInfo) SetLoginEx(_uid uint64, _token, _sessionKey string) {
	if _uid > 0 && _token != "" {
		this.IsLogin = true
		this.Uid = _uid
		this.Token = _token
		this.SessionKey = _sessionKey
		AddUser(this)
	} else {
		DelUser(this)
	}
}

// 是否登录
func (this *AuthInfo) CheckLogin() bool {
	return this.IsLogin
}
