package rule

import "fmt"

type ParamInfo struct {
	Name string					// 参数名称
	ParamType uint32			// 参数类型
	Static bool					// 是否严格检查，如果为true, 参数检测失败后拒绝处理这个请求，否则使用默认值继续处理
	Default string				// 如果参数检测不通过则使用这个默认值进行处理
	Extra []string				// 扩展参数
}


//@author xiaolan
//@lastUpdate 2019-08-10
//@comment 获取一个新的参数配置
//@param _name 参数名称
//@param _type 参数类型
//@param _static 是否严格检查
//@param _default 默认方式
func NewParam(_name string, _type uint32, _static bool, _default string) ParamInfo {
	return ParamInfo {
		Name: _name,
		ParamType: _type,
		Static: _static,
		Default: _default,
		Extra: []string{},
	}
}


//@author xiaolan
//@lastUpdate 2022-09-15
//@comment 创建一个int类型的参数配置, 使得这个整数必须处于指定数值区间下
//@param _name 参数名称
//@param _static 是否严格检查
//@param _default 默认值（如果非严格检查情况下，出现不符合检查规则的参数，则用默认值代替）
//@param _min 最小整数
//@param _max 最大整数
func NewIntParamRange(_name string, _static bool, _default string, _min, _max int) ParamInfo {
	return ParamInfo{
		Name: _name,
		ParamType: PTYPE_INT_RANGE,
		Static: _static,
		Default: _default,
		Extra: []string{
			fmt.Sprintf("%v", _min),
			fmt.Sprintf("%v", _max),
		},
	}
}


//@author xiaolan
//@lastUpdate 2022-09-15
//@comment 创建一个字符串类型的参数配置, 使得这个字符串长度必须在指定的数值范围内
//@param _name 参数名称
//@param _static 是否严格检查
//@param _default 默认值（如果非严格检查情况下，出现不符合检查规则的参数，则用默认值代替）
//@param _min 最小长度
//@param _max 最大长度
func NewStrParamRange(_name string, _static bool, _default string, _min, _max int) ParamInfo {
	return ParamInfo{
		Name: _name,
		ParamType: PTYPE_STR_RANGE,
		Static: _static,
		Default: _default,
		Extra: []string{
			fmt.Sprintf("%v", _min),
			fmt.Sprintf("%v", _max),
		},
	}
}