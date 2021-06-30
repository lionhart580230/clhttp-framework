package clAuth

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/xiaolan580230/clhttp-framework/clCommon"
	"strings"
)

// 管理auth扩展属性
// 便于对auth包进行自定义数据扩展
// auth将自动处理redis存储与载入

// 设置信息
func (this *AuthInfo) SetItem(_key string, _val interface{}) {
	this.mLocker.RLock()
	defer this.mLocker.RLock()

	val := fmt.Sprintf("%v", _val)
	if val[0] == '{' || val[0] == '[' {
		jsonBytes, _ := json.Marshal(_val)
		val = string(jsonBytes)
	}
	this.ExtraData[_key] = val
	SaveUser(this)
}

// 获取信息 (保留，以便兼容旧版本)
func (this *AuthInfo) GetItem(_key string) string {
	this.mLocker.RLock()
	defer this.mLocker.RLock()

	return this.ExtraData[_key]
}


// 获取字符串信息
func (this *AuthInfo) GetStr(_key string) string {
	this.mLocker.RLock()
	defer this.mLocker.RLock()

	return this.ExtraData[_key]
}


// 获取Int64值
func (this *AuthInfo) GetInt64(_key string) int64 {
	this.mLocker.RLock()
	defer this.mLocker.RLock()

	return clCommon.Int64(this.ExtraData[_key])
}


// 获取Int32值
func (this *AuthInfo) GetInt32(_key string) int32 {
	this.mLocker.RLock()
	defer this.mLocker.RLock()

	return clCommon.Int32(this.ExtraData[_key])
}


// 获取Uint64值
func (this *AuthInfo) GetUint64(_key string) uint64 {
	this.mLocker.RLock()
	defer this.mLocker.RLock()

	return clCommon.Uint64(this.ExtraData[_key])
}


// 获取Uint32值
func (this *AuthInfo) GetUint32(_key string) uint32 {
	this.mLocker.RLock()
	defer this.mLocker.RLock()

	return clCommon.Uint32(this.ExtraData[_key])
}


// 获取Boolean值
func (this *AuthInfo) GetBool(_key string) bool {
	this.mLocker.RLock()
	defer this.mLocker.RLock()

	return clCommon.Bool(this.ExtraData[_key])
}


// 将结果根据指定类型分割
func (this *AuthInfo) GetSplitBy(_key string, _ceil string) []string {
	this.mLocker.RLock()
	defer this.mLocker.RLock()

	return strings.Split(this.ExtraData[_key], _ceil)
}


// 获取指定数据类型
func (this *AuthInfo) GetObject(_key string, _data interface{}) error {
	this.mLocker.RLock()
	defer this.mLocker.RLock()
	jsonStr, exists := this.ExtraData[_key]
	if !exists {
		return errors.New("not found")
	}
	return json.Unmarshal([]byte(jsonStr), _data)
}