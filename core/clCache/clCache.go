package clCache

import (
	"encoding/json"
	"fmt"
	"github.com/xiaolan580230/clhttp-framework/clCrypt"
	"github.com/xiaolan580230/clhttp-framework/clGlobal"
	"github.com/xiaolan580230/clhttp-framework/core/cljson"
	"reflect"
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
	var data = clCrypt.Base64Encode( jsonStr )
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
		jsonStr = cljson.CreateBy(_obj).ToStr()
	}

	var data = clCrypt.Base64Encode( jsonStr )
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
		var cacheObj clCache
		err := json.Unmarshal([]byte(val), &cacheObj)
		if err != nil {
			return ""
		}
		return string(clCrypt.Base64Decode(cacheObj.Data))
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
