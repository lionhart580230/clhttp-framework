package skyconfig

import (
	"sync"
	"time"
)

type sectionType map[string] string


// 配置结构 2019-08-10
type Config struct {
	config map[string] sectionType
	fileName string
	autoLoad time.Duration
	lock sync.RWMutex
	stringMap map[/*section:key*/ string] autoLoadString
	int64Map map[/*section:key*/ string] autoLoadInt64
	int32Map map[/*section:key*/ string] autoLoadInt32
	uint64Map map[/*section:key*/ string] autoLoadUint64
	uint32Map map[/*section:key*/ string] autoLoadUint32
	float32Map map[/*section:key*/ string] autoLoadFloat32
	float64Map map[/*section:key*/ string] autoLoadFloat64
	boolMap map[/*section:key*/ string] autoLoadBool
	sectionMap map[/*section*/ string] autoLoadSection
	ArrMap map[/*section*/ string] autoLoadArr
}

// (自动重载) 字符串存储结构
type autoLoadString struct {
	section string
	key string
	def string
	value *string
}

// (自动重载) int64数据结构
type autoLoadInt64 struct {
	section string
	key string
	def int64
	value *int64
}

// (自动重载) int32数据结构
type autoLoadInt32 struct {
	section string
	key string
	def int32
	value *int32
}

// (自动重载) uint64数据结构
type autoLoadUint64 struct {
	section string
	key string
	def uint64
	value *uint64
}

// (自动重载) uint32数据结构
type autoLoadUint32 struct {
	section string
	key string
	def uint32
	value *uint32
}

// (自动重载) bool数据结构
type autoLoadBool struct {
	section string
	key string
	def bool
	value *bool
}

// (自动重载) float32数据结构
type autoLoadFloat32 struct {
	section string
	key string
	def float32
	value *float32
}

// (自动重载) float64数据结构
type autoLoadFloat64 struct {
	section string
	key string
	def float64
	value *float64
}

// (自动重载) map[string]string 数据结构
type autoLoadSection struct {
	section string
	value *map[string] string
}

// (自动重载) []string 数据结构
type autoLoadArr struct {
	section string
	value *[] string
}