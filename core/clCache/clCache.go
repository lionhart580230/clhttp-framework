package clCache

import (
	"fmt"
	"github.com/xiaolan580230/clUtil/clCrypt"
	"github.com/xiaolan580230/clUtil/clJson"
	"github.com/xiaolan580230/clhttp-framework/clGlobal"
	"reflect"
	"strings"
	"sync"
	"time"
)

// 缓存管理器
var mMemoryCache map[string] clCache
var mLocker sync.RWMutex


type clCache struct {
	Data string `json:"d"`
	Expire uint32 `json:"e"`
}


func init() {
	mMemoryCache = make(map[string] clCache)
}


// 简单类型写入缓存
func UpdateCacheSimple(_key string, _obj interface{}, _expire uint32) {
	mLocker.Lock()
	defer mLocker.Unlock()

	var jsonStr = fmt.Sprintf("%v", _obj)
	var data = clCrypt.Base64Encode( []byte( jsonStr ) )
	var cacheObj = clCache{
		Data:   data,
		Expire: uint32(time.Now().Unix()) + _expire,
	}

	if clGlobal.SkyConf.IsCluster {
		redis := clGlobal.GetRedis()
		if redis != nil {
			redis.Set(_key, data, int32(_expire))
		}
	} else {
		mMemoryCache[_key] = cacheObj
	}
}



// 写入缓存
func UpdateCache(_key string, _obj interface{}, _expire uint32) {
	mLocker.Lock()
	defer mLocker.Unlock()

	var jsonStr string
	if IsSimpleType(_obj) {
		jsonStr = fmt.Sprintf("%v", _obj)
	} else {
		jsonStr = clJson.CreateBy(_obj).ToStr()
	}

	var data = clCrypt.Base64Encode( []byte( jsonStr ) )
	var cacheObj = clCache{
		Data:   data,
		Expire: uint32(time.Now().Unix()) + _expire,
	}

	if clGlobal.SkyConf.IsCluster {
		redis := clGlobal.GetRedis()
		if redis != nil {
			redis.Set(_key, data, int32(_expire))
		}
	} else {
		mMemoryCache[_key] = cacheObj
	}
}



// 是否是简单类型
func IsSimpleType(_val interface{}) bool {
	switch reflect.TypeOf(_val).Kind() {
	case reflect.Int:
		fallthrough
	case reflect.Float32:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		fallthrough
	case reflect.Float64:
		fallthrough
	case reflect.String:
		fallthrough
	case reflect.Bool:
		return true
	default:
		return false
	}
}




// 获取缓存
func GetCache(_key string) string {
	mLocker.RLock()
	defer mLocker.RUnlock()

	if clGlobal.SkyConf.IsCluster {
		redis := clGlobal.GetRedis()
		if redis == nil {
			return ""
		}

		val := redis.Get(_key)
		if val == "" {
			return ""
		}
		return string(clCrypt.Base64Decode(val))
	} else {
		val, exist := mMemoryCache[_key]
		if !exist {
			return ""
		}
		if val.Expire < uint32(time.Now().Unix()) {
			return ""
		}
		return string(clCrypt.Base64Decode(val.Data))
	}
}



// 删除缓存
func DelCache(_key string){
	mLocker.Lock()
	defer mLocker.Lock()

	if clGlobal.SkyConf.IsCluster {
		redis := clGlobal.GetRedis()
		if redis != nil {
			redis.Del(_key)
		}
	} else {
		delete(mMemoryCache, _key)
	}
}


// 批量删除缓存
func DelCacheContains(_key string) {
	mLocker.Lock()
	defer mLocker.Lock()

	if clGlobal.SkyConf.IsCluster {
		redis := clGlobal.GetRedis()
		if redis != nil {
			keyList := redis.GetKeys(fmt.Sprintf("*%v*", _key))
			for _, key := range keyList {
				redis.Del(key)
			}
		}
	} else {
		for key, _ := range mMemoryCache {
			if strings.Contains(key, _key) {
				delete(mMemoryCache, key)
			}
		}

	}
}


// 移除接口缓存
func DelApiCache() {

}