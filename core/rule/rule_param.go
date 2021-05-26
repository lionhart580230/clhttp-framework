package rule


type ParamInfo struct {
	Name string					// 参数名称
	ParamType uint32			// 参数类型
	Static bool					// 是否严格检查，如果为true, 参数检测失败后拒绝处理这个请求，否则使用默认值继续处理
	Default string				// 如果参数检测不通过则使用这个默认值进行处理
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
	}
}


