package clAuth

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lionhart580230/clhttp-framework/clCommon"
	"strings"
)

// 管理auth扩展属性
// 便于对auth包进行自定义数据扩展
// auth将自动处理redis存储与载入

// 设置信息
func (this *AuthInfo) SetItem(_key string, _val interface{}) *AuthInfo {
	if this == nil {
		return this
	}
	this.mLocker.Lock()
	defer this.mLocker.Unlock()

	val := fmt.Sprintf("%v", _val)
	if val != "" {
		if val[0] == '{' || val[0] == '[' || strings.HasPrefix(val, "map[") {
			jsonBytes, _ := json.Marshal(_val)
			val = string(jsonBytes)
		}
	}
	this.ExtraData[_key] = val
	SaveUser(this)
	return this
}

// 批量设置信息
func (this *AuthInfo) SetItems(_datas map[string]interface{}) *AuthInfo {
	if this == nil {
		return this
	}
	this.mLocker.Lock()
	defer this.mLocker.Unlock()

	for _key, _val := range _datas {
		val := fmt.Sprintf("%v", _val)
		if val != "" {
			if val[0] == '{' || val[0] == '[' || strings.HasPrefix(val, "map[") {
				jsonBytes, _ := json.Marshal(_val)
				val = string(jsonBytes)
			}
		}
		this.ExtraData[_key] = val
	}
	SaveUser(this)
	return this
}

// 获取信息 (保留，以便兼容旧版本)
func (this *AuthInfo) GetItem(_key string) string {
	if this == nil {
		return ""
	}
	return this.GetStr(_key)
}

// 获取字符串信息
func (this *AuthInfo) GetStr(_key string) string {
	if this == nil {
		return ""
	}
	this.mLocker.RLock()
	defer this.mLocker.RUnlock()

	val, exists := this.ExtraData[_key]
	if !exists {
		return ""
	}
	return val
}

func (this *AuthInfo) PrintMap() {
	if this == nil {
		return
	}
	this.mLocker.RLock()
	defer this.mLocker.RUnlock()

	for k, v := range this.ExtraData {
		fmt.Printf(">> %v -> %v\n", k, v)
	}
}

// 获取Int64值
func (this *AuthInfo) GetInt64(_key string) int64 {
	if this == nil {
		return 0
	}
	return clCommon.Int64(this.GetStr(_key))
}

// 获取Int32值
func (this *AuthInfo) GetInt32(_key string) int32 {
	if this == nil {
		return 0
	}
	return clCommon.Int32(this.GetStr(_key))
}

// 获取Uint64值
func (this *AuthInfo) GetUint64(_key string) uint64 {
	if this == nil {
		return 0
	}
	return clCommon.Uint64(this.GetStr(_key))
}

// 获取Uint32值
func (this *AuthInfo) GetUint32(_key string) uint32 {
	if this == nil {
		return 0
	}
	return clCommon.Uint32(this.GetStr(_key))
}

// 获取Boolean值
func (this *AuthInfo) GetBool(_key string) bool {
	if this == nil {
		return false
	}
	return clCommon.Bool(this.GetStr(_key))
}

// 获取Float64
func (this *AuthInfo) GetFloat64(_key string) float64 {
	if this == nil {
		return 0
	}
	return clCommon.Float64(this.GetStr(_key))
}

// 获取Float32
func (this *AuthInfo) GetFloat32(_key string) float32 {
	if this == nil {
		return 0
	}
	return clCommon.Float32(this.GetStr(_key))
}

// 将结果根据指定类型分割
func (this *AuthInfo) GetSplitBy(_key string, _ceil string) []string {
	if this == nil {
		return []string{}
	}
	val := this.GetStr(_key)
	if val == "" {
		return []string{}
	}
	return strings.Split(val, _ceil)
}

// 获取指定数据类型
func (this *AuthInfo) GetObject(_key string, _data interface{}) error {
	if this == nil {
		return errors.New("nil pointer")
	}
	jsonStr := this.GetStr(_key)
	if jsonStr == "" {
		return errors.New("not found")
	}
	return json.Unmarshal([]byte(jsonStr), _data)
}

// 获取指定数据类型
func (this *AuthInfo) KeyIsExists(_key string) bool {
	if this == nil {
		return false
	}
	jsonStr := this.GetStr(_key)
	if jsonStr == "" {
		return false
	}
	return true
}
