package clCache

import (
	"github.com/xiaolan580230/clhttp-framework/clCommon"
	"github.com/xiaolan580230/clhttp-framework/core/cljson"
	"sync"
	"time"
)

// 缓存管理器

// 缓存类型 0=内存缓存,1=文件缓存,2=redis缓存

const CacheTypeMemory = 0
const CacheTypeFile = 1
const CacheTypeRedis = 2

var mCacheType = CacheTypeMemory
var mMemoryCache map[string] clCache
var mLocker sync.RWMutex


type clCache struct {
	Data string `json:"d"`
	Expire uint32 `json:"e"`
}


func init() {
	mMemoryCache = make(map[string] clCache)
}

// 设置缓存类型
func SetCacheType(_type int) {
	mCacheType = _type
}


// 写入缓存
func UpdateCache(_key string, _obj interface{}, _expire uint32) {
	mLocker.Lock()
	defer mLocker.Unlock()

	var data = clCommon.Base64Encode( cljson.CreateBy(_obj).ToStr() )
	var cacheObj = clCache{
		Data:   data,
		Expire: uint32(time.Now().Unix()) + _expire,
	}


	if mCacheType == 0 {
		mMemoryCache[_key] = cacheObj
	} else if mCacheType == 1 {
		// TODO 写入到文件
	} else if mCacheType == 2 {
		// TODO 写入到redis
	}
}


// 获取缓存
func GetCache(_key string) string {
	mLocker.RLock()
	defer mLocker.RUnlock()

	val, exist := mMemoryCache[_key]
	if !exist {
		return ""
	}
	if val.Expire < uint32(time.Now().Unix()) {
		return ""
	}
	return string(clCommon.Base64Decode(val.Data))
}
