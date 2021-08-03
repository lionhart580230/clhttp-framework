package skylog

import (
	"encoding/json"
	"fmt"
	"github.com/xiaolan580230/clhttp-framework/core/cltime"
	"path"
	"runtime"
)


func init() {
	logObj = &ClLog{
		logLevel: LOG_LEVEL_ALL,
		version:  "nil",
		logType:  LOG_TYPE_CONSOLE,
		sip:      "",
	}
}


//@author xiaolan
//@lastUpdate 2019-08-04
//@comment 创建日志对象
//@param _version 版本号
func New(_version string) error {
	logObj = &ClLog{
		logLevel: LOG_LEVEL_ALL,
		version:  _version,
		logType:  LOG_TYPE_ALI,
		sip:      "",
	}

	return nil
}


/**
	切换日志记录模式
	@param mode uint32 只有指定模式的日志才会被记录下来
 */
func SetLevel(level uint32) {
	logObj.logLevel = level
}


func SetType(_type uint32) {
	logObj.logType = _type
}


/**
	记录一般的日志
	@param fileFormat string  记录的日志格式，类似log.Printf的format
   @param args interface{} 各种参数
 */
func LogInfo(_fileFormat string, _args... interface{}) {

	if (logObj.logLevel & LOG_LEVEL_INFO) > 0 {
		MakeLineHead("", _fileFormat, LOG_LEVEL_INFO, _args...)
	}

}


/**
	记录警告的日志
	@param fileFormat string  记录的日志格式，类似log.Printf的format
	@param args interface{} 各种参数
*/
func LogWarning(_fileFormat string, _args... interface{}) {

	if (logObj.logLevel & LOG_LEVEL_WARNING) > 0 {
		MakeLineHead("", _fileFormat, LOG_LEVEL_WARNING, _args...)
	}
}

/**
	记录错误的日志
	@param fileFormat string  记录的日志格式，类似log.Printf的format
	@param args interface{} 各种参数
*/
func LogErr(_fileFormat string, _args... interface{}) {
	if (logObj.logLevel & LOG_LEVEL_ERR) > 0 {
		MakeLineHead("", _fileFormat, LOG_LEVEL_ERR, _args...)
	}

}


/**
	记录调试的日志
	@param fileFormat string  记录的日志格式，类似log.Printf的format
	@param args interface{} 各种参数
*/
func LogDebug(_fileFormat string, _args... interface{}) {
	if (logObj.logLevel & LOG_LEVEL_DEBUG) > 0 {
		MakeLineHead("", _fileFormat, LOG_LEVEL_DEBUG, _args...)
	}
}

/**
	日志行头
	@param log_type int32	   日志类型
	@param fileFormat string   记录的日志格式，类似log.Printf的format
	@param color uint8         日志打印颜色
	@param args interface{}    各种参数
*/
func MakeLineHead(_title, fileFormat string, level uint8, args... interface{}) string{
	if logObj == nil {
		return ""
	}

	var logContext string
	_, file, line, ok := runtime.Caller(2)
	filenameInfo := ""
	if ok {
		_, filename := path.Split(file)
		filenameInfo = fmt.Sprintf("%v:%d", filename, line)
	}

	switch logObj.logType {
	case LOG_TYPE_CONSOLE:
		linehead := fmt.Sprintf("[%s %v>>%s]%s", cltime.GetDateByFormat(0, "15:04:05"), filenameInfo, logObj.version, fileFormat)
		logContext = fmt.Sprintf(linehead, args...)

		color := COLOR_WHITE
		switch level{
		case LOG_LEVEL_INFO: color = COLOR_WHITE
		case LOG_LEVEL_WARNING: color = COLOR_YELLOW
		case LOG_LEVEL_ERR: color = COLOR_ORINGE
		case LOG_LEVEL_DEBUG: color = COLOR_BLUE
		}
		fmt.Printf("\x1b[0;%dm%v\x1b[0m\n", color, logContext)
	case LOG_TYPE_AWS:
		nowTime, err := GetTimeFormat(0, "01-02 15:04:05")
		if err != nil {
			fmt.Printf("log time format failed: %v", err)
		}
		linehead := fmt.Sprintf("[%s %v>>%s]%s", nowTime, filenameInfo, logObj.version, fileFormat)
		logContext = fmt.Sprintf(linehead, args...)

		loglevel := "UNKNOW"
		switch level{
		case LOG_LEVEL_INFO: loglevel = "INFO"
		case LOG_LEVEL_WARNING: loglevel = "WARNING"
		case LOG_LEVEL_ERR: loglevel = "ERROR"
		case LOG_LEVEL_DEBUG: loglevel = "DEBUG"
		}
		fmt.Printf("(%v-%v)%v\n", logObj.sip, loglevel, logContext)
	case LOG_TYPE_ALI:
		loglevel := "UNKNOW"
		switch level{
		case LOG_LEVEL_INFO: loglevel = "INFO"
		case LOG_LEVEL_WARNING: loglevel = "WARNING"
		case LOG_LEVEL_ERR: loglevel = "ERROR"
		case LOG_LEVEL_DEBUG: loglevel = "DEBUG"
		}
		linehead := fmt.Sprintf("[%v>>%s]%s", filenameInfo, logObj.version, fileFormat)
		logContext = fmt.Sprintf(linehead, args...)
		fmt.Printf("%v <%v> %v\n", cltime.GetDateByFormat(0, "15:04:05"), loglevel, logContext)
	case LOG_TYPE_JSON:
		loglevel := "UNKNOW"
		switch level{
		case LOG_LEVEL_INFO: loglevel = "INFO"
		case LOG_LEVEL_WARNING: loglevel = "WARNING"
		case LOG_LEVEL_ERR: loglevel = "ERROR"
		case LOG_LEVEL_DEBUG: loglevel = "DEBUG"
		}

		nowTime, err := GetTimeFormat(0, "01-02 15:04:05")
		if err != nil {
			fmt.Printf("log time format failed: %v", err)
		}
		var json_str, _ = json.Marshal(jsonLog{
			Time:        nowTime,
			Version:     logObj.version,
			Title:       _title,
			File:        fmt.Sprintf("%v:%v", file, line),
			Level:       loglevel,
			Description: logContext,
		})
		fmt.Println(json_str)
	}

	return logContext
}