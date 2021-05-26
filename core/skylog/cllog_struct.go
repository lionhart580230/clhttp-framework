package skylog


const (
	LOG_LEVEL_INFO = 1
	LOG_LEVEL_WARNING = 2
	LOG_LEVEL_ERR = 4
	LOG_LEVEL_DEBUG = 8
	LOG_LEVEL_ALL = LOG_LEVEL_INFO | LOG_LEVEL_WARNING | LOG_LEVEL_ERR | LOG_LEVEL_DEBUG
)

const (
	COLOR_ORINGE = uint8(iota+91)
	COLOR_GREEN
	COLOR_YELLOW
	COLOR_PURPLE
	COLOR_MAGENTA
	COLOR_BLUE
	COLOR_WHITE
)


const (
	LOTTERY_DEBUG_VERSION = "v1.0.1"
)

const (
	LOG_TYPE_CONSOLE = 0
	LOG_TYPE_AWS = 1
	LOG_TYPE_ALI = 2
	LOG_TYPE_JSON = 3
)


type ClLog struct {
	logLevel uint32
	version string
	logType uint32
	sip string
}


type jsonLog struct {
	Time string `json:"time"`			// 时间
	Version string `json:"version"`		// 版本号
	Title string `json:"title"`			// 标题
	File string `json:"file"`			// 所在文件
	Level string `json:"level"`			// 日志级别
	Description string `json:"detail"`	// 日志详情
}


var logObj *ClLog