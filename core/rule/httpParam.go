package rule

import (
	"strconv"
	"strings"
)

type HttpParam struct {
	values map[string]string
}

//@author xiaolan
//@lastUpdate 2019-08-10
//@comment 添加一个参数到参数列表中
//@param _key 参数的名称
//@param _val 参数的值
func (this *HttpParam) Add(_key, _val string) {
	this.values[_key] = _val
}


func NewHttpParam(_params map[string] string) *HttpParam {
	if _params == nil {
		return &HttpParam{
			values: make(map[string] string),
		}
	}
	return &HttpParam{
		values: _params,
	}
}

//@author xiaolan
//@lastUpdate 2019-08-10
//@comment 获取一个字符串型的参数
//@param _key 要获取的参数名称
//@param _default 如果key不存在, 默认返回什么
func (this *HttpParam) GetStr(_key string, _default string) string {
	val, exists := this.values[_key]
	if !exists {
		return _default
	}
	return val
}


//@author xiaolan
//@lastUpdate 2019-08-10
//@comment 获取uint32类型的参数
//@param _key 要获取的参数名称
//@param _default 如果key不存在, 默认返回什么
func (this *HttpParam) GetUint32(_key string, _default uint32) uint32 {

	val, exists := this.values[_key]
	if !exists {
		return _default
	}

	i, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		return _default
	}
	return uint32(i)
}


//@author xiaolan
//@lastUpdate 2019-08-10
//@comment 获取uint64类型的参数
//@param _key 要获取的参数名称
//@param _default 如果key不存在, 默认返回什么
func (this *HttpParam) GetUint64(_key string, _default uint64) uint64 {

	val, exists := this.values[_key]
	if !exists {
		return _default
	}

	i, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return _default
	}
	return uint64(i)
}


//@author xiaolan
//@lastUpdate 2019-08-10
//@comment 获取int32类型的参数
//@param _key 要获取的参数名称
//@param _default 如果key不存在, 默认返回什么
func (this *HttpParam) GetInt32(_key string, _default int32) int32 {

	val, exists := this.values[_key]
	if !exists {
		return _default
	}

	i, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		return _default
	}
	return int32(i)
}


//@author xiaolan
//@lastUpdate 2019-08-10
//@comment 获取int64类型的参数
//@param _key 要获取的参数名称
//@param _default 如果key不存在, 默认返回什么
func (this *HttpParam) GetInt64(_key string, _default int64) int64 {

	val, exists := this.values[_key]
	if !exists {
		return _default
	}

	i, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return _default
	}
	return int64(i)
}


//@author xiaolan
//@lastUpdate 2019-08-10
//@comment 获取32位浮点数
//@param _key 要获取的参数名称
//@param _default 如果key不存在，默认返回什么
func (this *HttpParam) GetFloat32(_key string, _default float32) float32 {

	val, exists := this.values[_key]
	if !exists {
		return _default
	}

	i, err := strconv.ParseFloat(val, 32)
	if err != nil {
		return _default
	}
	return float32(i)
}


//@author xiaolan
//@lastUpdate 2019-08-10
//@comment 获取64位浮点数
//@param _key 要获取的参数名称
//@param _default 如果key不存在，默认返回什么
func (this *HttpParam) GetFloat64(_key string, _default float64) float64 {

	val, exists := this.values[_key]
	if !exists {
		return _default
	}

	i, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return _default
	}
	return float64(i)
}


//@author xiaolan
//@lastUpdate 2019-08-10
//@comment 获取浮点数类型
//@param _key 要获取的参数名称
//@param _default 如果不存在默认返回什么
func (this *HttpParam) GetBool(_key string, _default bool) bool {

	val, exists := this.values[_key]
	if !exists {
		return _default
	}

	switch strings.ToUpper(val) {
	case "OK", "ON", "YES", "TRUE", "Y", "T":
		return true
	}

	return false
}
