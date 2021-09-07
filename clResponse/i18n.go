package clResponse

import (
	"fmt"
	"github.com/xiaolan580230/clUtil/clLog"
	"github.com/xiaolan580230/clhttp-framework/clCommon"
	"strings"
)

var mI18NMap I18NMap
func init() {
	mI18NMap = make(I18NMap)
}
// 导入i18n数据
// data必须是通过\n分割的数据，格式为 id,字符串
func ImportI18N(_langType uint32, _data string) {

	mI18NMap[_langType] = make(map[uint32] string)
	data := strings.Split(_data, "\n")
	for _, val := range data {
		items := strings.Split(strings.Trim(val, "\n\r"), "=")
		if len(items) < 2 {
			continue
		}
		mI18NMap[_langType][ clCommon.Uint32(items[0])] = strings.Join(items[1:], "=")
	}
	clLog.Debug("加载:%v语言包成功! %+v", _langType, mI18NMap[_langType])
}


// 生成字符串
func GenStr(_langType uint32, _id uint32, _param ...interface{}) string {

	strMap, exists := mI18NMap[_langType]
	if !exists {
		return ""
	}

	strFmt, exists := strMap[_id]
	if !exists {
		return ""
	}

	if _param == nil || len(_param) == 0 {
		return strFmt
	}
	return fmt.Sprintf(strFmt, _param ...)
}