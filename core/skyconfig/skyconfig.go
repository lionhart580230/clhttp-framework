package skyconfig

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

//@author xiaolan
//@lastUpdate 2019-08-10
//@comment 控制配置的读取
//方式一。 通过文件读取
//方式二。 通过环境变量读取


//@author xiaolan
//@lastUpdate 2019-08-10
//@comment 新建一个配置对象
//@param _filename 配置文件名，如果文件名为空，则默认从环境变量中获取
//@param _autoLoad 自动重载的时间（间隔多久自动重载一次), 如果为0则放弃自动重载机制
func New(_filename string, _autoLoad time.Duration) *Config {
	var config  = Config {
		fileName: _filename,
		config: make(map[string] sectionType),
		autoLoad: _autoLoad,
		stringMap: make(map[string] autoLoadString),
		int64Map: make( map[/*section:key*/ string] autoLoadInt64),
		int32Map: make( map[/*section:key*/ string] autoLoadInt32),
		uint64Map: make( map[/*section:key*/ string] autoLoadUint64),
		uint32Map: make( map[/*section:key*/ string] autoLoadUint32),
		float32Map: make( map[/*section:key*/ string] autoLoadFloat32),
		float64Map: make( map[/*section:key*/ string] autoLoadFloat64),
		boolMap: make( map[/*section:key*/ string] autoLoadBool),
		sectionMap: make(map[/*section*/ string] autoLoadSection),
	}

	loadFile(&config)
	if _autoLoad > 0 {
		go config.autoLoadConfig()
	}
	return &config
}


//@author xiaolan
//@lastUpdate 2019-08-10
//@comment 自动重载机制
func (config *Config) autoLoadConfig() {
	for {
		<-time.After(config.autoLoad * time.Second )
		loadFile(config)
		config.lock.Lock()
		for _, val := range config.stringMap {
			config.GetStr(val.section, val.key, val.def, val.value)
		}
		config.lock.Unlock()
	}
}


//@author xiaolan
//@lastUpdate 2019-08-10
//@comment 开始加载配置
func loadFile(config *Config) {
	defer doWhenErr()

	if config.fileName == "" {
		return
	}

	dat, err := ioutil.ReadFile(config.fileName)
	if err != nil {
		panic(err)
		return
	}

	// 将读入的文件进行换行切割
	confArr := strings.Split(string(dat),"\n")

	// 创建一个全局的配置小节
	section := "global"
	config.config[section] = make(sectionType)

	// 是否注释
	isHelp := false
	// 遍历每一行,提取里面有用的数据
	for _,v := range confArr {
		if len(v) < 2 {
			continue
		}
		v = strings.TrimSpace(v)
		v = strings.TrimPrefix(v, "\n")

		if strings.HasPrefix(v, "#") || strings.HasPrefix(v, "//") {
			continue
		}

		if isHelp {
			if strings.HasPrefix(v, "*/") {
				isHelp = false
			}
			continue
		}

		if strings.HasPrefix(v, "/*") {
			isHelp = true
			continue
		}

		if strings.HasPrefix(v, "[") && strings.HasSuffix(v, "]") {
			// 小节
			section = strings.TrimPrefix(strings.TrimSuffix(v,"]"),"[")
			config.config[section] = make(sectionType)
		} else {
			// 配置项
			arr := strings.SplitN(v, "=", 2)
			if len(arr) != 2 {
				continue
			}
			config.config[section][strings.TrimSpace(arr[0])] = arr[1]
		}
	}
}

//@author xiaolan
//@lastUpdate 2019-08-10
//@comment 错误捕获加异常处理
func doWhenErr() {
	// 发生错误
	if err := recover(); err != nil {
		fmt.Printf("错误: %v\n",err)
	}
}

//@author xiaolan
//@lastUpdate 2019-08-10
//@comment 获取string型的配置, 如果文件获取不到, 尝试从环境变量中获取, 如果都获取不到返回false, 否则返回true
//@param _section 区段名称
//@param _key 配置键名
//@param _def 默认值
//@param _value 接收变量到指针
func (config *Config) GetStr(_section string, _key string, _def string, _value *string) bool {

	var configVal = _def
	var isExists = false
	if len(config.config[_section]) > 0 {
		if len(config.config[_section][_key]) > 0 {
			configVal = config.config[_section][_key]
			isExists = true
		}
	}

	if !isExists {
		// 配置不存在，尝试从环境变量中获取
		envConfig := os.Getenv(strings.ToUpper(_section + "_" + _key))
		if envConfig != "" {
			configVal = envConfig
			isExists = true
		}
	}

	if config.autoLoad > 0 && isExists {
		config.stringMap[_section + ":" + _key] = autoLoadString{
			section: _section,
			key: _key,
			def: _def,
			value: _value,
		}
	}
	*_value = configVal
	return isExists
}

//@author xiaolan
//@lastUpdate 2019-08-10
//@comment 获取float32型的配置, 如果文件获取不到, 尝试从环境变量中获取, 如果都获取不到返回false, 否则返回true
//@param _section 区段名称
//@param _key 配置键名
//@param _def 默认值
//@param _value 接收变量到指针
func (config *Config) GetFloat32(_section string, _key string, _def float32, _value *float32) bool {
	var configVal = _def
	var isExists = false
	if len(config.config[_section]) > 0 {
		if len(config.config[_section][_key]) > 0 {
			fval, err := strconv.ParseFloat(config.config[_section][_key], 32)
			if err == nil {
				configVal = float32(fval)
			}
			isExists = true
		}
	}

	if !isExists {
		// 配置不存在，尝试从环境变量中获取
		envConfig := os.Getenv(strings.ToUpper(_section + "_" + _key))
		if envConfig != "" {
			fval, err := strconv.ParseFloat(envConfig, 32)
			if err == nil {
				configVal = float32(fval)
			}
			isExists = true
		}
	}

	if config.autoLoad > 0 && isExists {
		config.float32Map[_section + ":" + _key] = autoLoadFloat32{
			section: _section,
			key: _key,
			def: _def,
			value: _value,
		}
	}
	*_value = configVal
	return isExists
}


//@author xiaolan
//@lastUpdate 2019-08-10
//@comment 获取float64型的配置, 如果文件获取不到, 尝试从环境变量中获取, 如果都获取不到返回false, 否则返回true
//@param _section 区段名称
//@param _key 配置键名
//@param _def 默认值
//@param _value 接收变量到指针
func (config *Config) GetFloat64(_section string, _key string, _def float64, _value *float64) bool {
	var configVal = _def
	var isExists = false
	if len(config.config[_section]) > 0 {
		if len(config.config[_section][_key]) > 0 {
			fval, err := strconv.ParseFloat(config.config[_section][_key], 64)
			if err == nil {
				configVal = fval
			}
			isExists = true
		}
	}

	if !isExists {
		// 配置不存在，尝试从环境变量中获取
		envConfig := os.Getenv(strings.ToUpper(_section + "_" + _key))
		if envConfig != "" {
			fval, err := strconv.ParseFloat(envConfig, 64)
			if err == nil {
				configVal = fval
			}
			isExists = true
		}
	}

	if config.autoLoad > 0 && isExists {
		config.float64Map[_section + ":" + _key] = autoLoadFloat64{
			section: _section,
			key: _key,
			def: _def,
			value: _value,
		}
	}
	*_value = configVal
	return isExists
}


//@author xiaolan
//@lastUpdate 2019-08-10
//@comment 获取int32型的配置, 如果文件获取不到, 尝试从环境变量中获取, 如果都获取不到返回false, 否则返回true
//@param _section 区段名称
//@param _key 配置键名
//@param _def 默认值
//@param _value 接收变量到指针
func (config *Config) GetInt32(_section string, _key string, _def int32, _value *int32) bool {
	var configVal = _def
	var isExists = false
	if len(config.config[_section]) > 0 {
		if len(config.config[_section][_key]) > 0 {
			ival, err := strconv.ParseInt(config.config[_section][_key], 0, 32)
			if err == nil {
				configVal = int32(ival)
			}
			isExists = true
		}
	}

	if !isExists {
		// 配置不存在，尝试从环境变量中获取
		envConfig := os.Getenv(strings.ToUpper(_section + "_" + _key))
		if envConfig != "" {
			ival, err := strconv.ParseInt(envConfig, 0, 32)
			if err == nil {
				configVal = int32(ival)
			}
			isExists = true
		}
	}

	if config.autoLoad > 0 && isExists {
		config.int32Map[_section + ":" + _key] = autoLoadInt32{
			section: _section,
			key: _key,
			def: _def,
			value: _value,
		}
	}
	*_value = configVal
	return isExists
}


//@author xiaolan
//@lastUpdate 2019-08-10
//@comment 获取uint32型的配置, 如果文件获取不到, 尝试从环境变量中获取, 如果都获取不到返回false, 否则返回true
//@param _section 区段名称
//@param _key 配置键名
//@param _def 默认值
//@param _value 接收变量到指针
func (config *Config) GetUint32(_section string, _key string, _def uint32, _value *uint32) bool {
	var configVal = _def
	var isExists = false
	if len(config.config[_section]) > 0 {
		if len(config.config[_section][_key]) > 0 {
			ival, err := strconv.ParseUint(config.config[_section][_key], 0, 32)
			if err == nil {
				configVal = uint32(ival)
			}
			isExists = true
		}
	}

	if !isExists {
		// 配置不存在，尝试从环境变量中获取
		envConfig := os.Getenv(strings.ToUpper(_section + "_" + _key))
		if envConfig != "" {
			ival, err := strconv.ParseInt(envConfig, 0, 32)
			if err == nil {
				configVal = uint32(ival)
			}
			isExists = true
		}
	}

	if config.autoLoad > 0 && isExists {
		config.uint32Map[_section + ":" + _key] = autoLoadUint32{
			section: _section,
			key: _key,
			def: _def,
			value: _value,
		}
	}
	*_value = configVal
	return isExists
}


//@author xiaolan
//@lastUpdate 2019-08-10
//@comment 获取int64型的配置, 如果文件获取不到, 尝试从环境变量中获取, 如果都获取不到返回false, 否则返回true
//@param _section 区段名称
//@param _key 配置键名
//@param _def 默认值
//@param _value 接收变量到指针
func (config *Config) GetInt64(_section string, _key string, _def int64, _value *int64) bool {
	var configVal = _def
	var isExists = false
	if len(config.config[_section]) > 0 {
		if len(config.config[_section][_key]) > 0 {
			ival, err := strconv.ParseInt(config.config[_section][_key], 0, 64)
			if err == nil {
				configVal = ival
			}
			isExists = true
		}
	}

	if !isExists {
		// 配置不存在，尝试从环境变量中获取
		envConfig := os.Getenv(strings.ToUpper(_section + "_" + _key))
		if envConfig != "" {
			ival, err := strconv.ParseInt(envConfig, 0, 64)
			if err == nil {
				configVal = int64(ival)
			}
			isExists = true
		}
	}


	if config.autoLoad > 0 && isExists {
		config.int64Map[_section + ":" + _key] = autoLoadInt64{
			section: _section,
			key: _key,
			def: _def,
			value: _value,
		}
	}
	*_value = configVal
	return isExists
}


//@author xiaolan
//@lastUpdate 2019-08-10
//@comment 获取uint64型的配置, 如果文件获取不到, 尝试从环境变量中获取, 如果都获取不到返回false, 否则返回true
//@param _section 区段名称
//@param _key 配置键名
//@param _def 默认值
//@param _value 接收变量到指针
func (config *Config) GetUint64(_section string, _key string, _def uint64, _value *uint64) bool {
	var configVal = _def
	var isExists = false
	if len(config.config[_section]) > 0 {
		if len(config.config[_section][_key]) > 0 {
			ival, err := strconv.ParseUint(config.config[_section][_key], 0, 64)
			if err == nil {
				configVal = ival
			}
			isExists = true
		}
	}

	if !isExists {
		// 配置不存在，尝试从环境变量中获取
		envConfig := os.Getenv(strings.ToUpper(_section + "_" + _key))
		if envConfig != "" {
			ival, err := strconv.ParseInt(envConfig, 0, 64)
			if err == nil {
				configVal = uint64(ival)
			}
			isExists = true
		}
	}

	if config.autoLoad > 0 && isExists {
		config.uint64Map[_section + ":" + _key] = autoLoadUint64{
			section: _section,
			key: _key,
			def: _def,
			value: _value,
		}
	}
	*_value = configVal
	return isExists
}


//@author xiaolan
//@lastUpdate 2019-08-10
//@comment 获取bool型的配置, 如果文件获取不到, 尝试从环境变量中获取, 如果都获取不到返回false, 否则返回true
//@param _section 区段名称
//@param _key 配置键名
//@param _def 默认值
//@param _value 接收变量到指针
func (config *Config)GetBool(_section string, _key string, _def bool, _value *bool) bool {
	var configVal = _def
	var isExists = false
	if len(config.config[_section]) > 0 {
		if len(config.config[_section][_key]) > 0 {
			bval, err := strconv.ParseBool(config.config[_section][_key])
			if err == nil {
				configVal = bval
			}
			isExists = true
		}
	}

	if !isExists {
		// 配置不存在，尝试从环境变量中获取
		envConfig := os.Getenv(strings.ToUpper(_section + "_" + _key))
		if envConfig != "" {
			bval, err := strconv.ParseBool(envConfig)
			if err == nil {
				configVal = bval
			}
			isExists = true
		}
	}

	if config.autoLoad > 0 && isExists {
		config.boolMap[_section + ":" + _key] = autoLoadBool{
			section: _section,
			key: _key,
			def: _def,
			value: _value,
		}
	}
	*_value = configVal
	return isExists
}


//@author xiaolan
//@lastUpdate 2019-08-10
//@comment 获取完整的区块到map中，此方式暂不支持从环境变量中获取, 如果配置存在返回true, 否则返回false
//@param _section 区块名称
//@param _value 用于存放区块配置列表的变量指针
func (config *Config)GetFullSection(_section string, _value *map[string]string) bool {
	if len(config.config[_section]) > 0 {
		if config.autoLoad > 0 {
			config.sectionMap[_section] = autoLoadSection{
				section: _section,
				value: _value,
			}
		}

		*_value = config.config[_section]
		return true
	}
	return false
}