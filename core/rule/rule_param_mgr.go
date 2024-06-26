package rule

import (
	"github.com/lionhart580230/clhttp-framework/clCommon"
	"regexp"
	"strconv"
	"strings"
)

var paramCheckers map[uint32]func( /*_param*/ string /*extra*/, *ParamInfo) string

// 参数检验规则类型定义
const (
	PTYPE_ALL         = 0  // 不检查, 不推荐使用
	PTYPE_SAFE_STR    = 1  // 安全检查, 不包括注入
	PTYPE_TINY_INT    = 2  // 数字, 1-3位
	PTYPE_INT         = 3  // 数字, 1-10位
	PTYPE_LONG        = 4  // 数字, 1-20位
	PTYPE_FLOAT       = 5  // 浮点数（整数也可通过）
	PTYPE_DATE        = 6  // 日期: YYYY-MM-DD
	PTYPE_TIME        = 7  // 时间: HH:MM:SS
	PTYPE_DATETIME    = 8  // 日期时间: YYYY-MM-DD HH:MM:SS
	PTYPE_IP          = 9  // IP: 0.0.0.0 - 255.255.255.255
	PTYPE_MD5         = 10 // MD5: 32位数字+字母 或者 48，64位都可
	PTYPE_ASSERT_NAME = 11 // 支持jpg,jpeg,png,gif,mp4,avi,ogg等
	PTYPE_URL         = 12 // 访问地址
	PTYPE_USERNAME    = 13 // 用户名
	PTYPE_PASSWORD    = 14 // 密码
	PTYPE_EMAIL       = 15 // 邮箱
	PTYPE_NUMBER_LIST = 16 // 数字列表，用半角逗号隔开的数字
	PTYPE_PHONE       = 17 // 手机号码
	PTYPE_VCODE       = 18 // 短信验证码或邮箱验证码，固定6位数字
	PTYPE_ID_CARD     = 19 // 身份证号
	PTYPE_INT_RANGE   = 20 // 整数范围
	PTYPE_STR_RANGE   = 21 // 字符串长度
	PTYPE_V2          = 22 // 2维坐标系
	PTYPE_V3          = 23 // 3维坐标系
	PTYPE_V4          = 24 // 4维坐标系
	PTYPE_REGEX       = 25 // 正则表达式
	PTYPE_STR_IN      = 26 // 正则表达式
)

const (
	PARAM_CHECK_FAIED = "{{CHECK_PARAM_FAILED}}"
)

func init() {
	paramCheckers = make(map[uint32]func(string, *ParamInfo) string)

	// 无视检查条件
	paramCheckers[PTYPE_ALL] = func(_param string, _extra *ParamInfo) string {
		return _param
	}

	// 安全字符串
	paramCheckers[PTYPE_SAFE_STR] = func(_param string, _extra *ParamInfo) string {
		match, err := regexp.Match(`[\;\'\"\<\>]`, []byte(_param))
		if err != nil || match {
			return PARAM_CHECK_FAIED
		}
		return _param
	}

	// 短数字
	paramCheckers[PTYPE_TINY_INT] = func(_param string, _extra *ParamInfo) string {
		match, err := regexp.Match(`^(\-)?[0-9]{1,3}$`, []byte(_param))
		if err != nil || !match {
			return PARAM_CHECK_FAIED
		}

		return _param
	}

	// 中数字
	paramCheckers[PTYPE_INT] = func(_param string, _extra *ParamInfo) string {
		match, err := regexp.Match(`^(\-)?[0-9]{1,10}$`, []byte(_param))
		if err != nil || !match {
			return PARAM_CHECK_FAIED
		}
		return _param
	}

	// 长数字
	paramCheckers[PTYPE_LONG] = func(_param string, _extra *ParamInfo) string {
		match, err := regexp.Match(`^(\-)?[0-9]{1,20}$`, []byte(_param))
		if err != nil || !match {
			return PARAM_CHECK_FAIED
		}
		return _param
	}

	// 浮点数
	paramCheckers[PTYPE_FLOAT] = func(_param string, _extra *ParamInfo) string {
		match, err := regexp.Match(`^(\-)?[0-9]{1,20}(\.[0-9]{1,10})?$`, []byte(_param))
		if err != nil || !match {
			return PARAM_CHECK_FAIED
		}
		return _param
	}

	// 日期
	paramCheckers[PTYPE_DATE] = func(_param string, _extra *ParamInfo) string {
		match, err := regexp.Match(`^[0-9]{4}\-[01][0-9]\-([012][0-9]|[3][01])$`, []byte(_param))
		if err != nil || !match {
			return PARAM_CHECK_FAIED
		}
		return _param
	}

	// 时间
	paramCheckers[PTYPE_TIME] = func(_param string, _extra *ParamInfo) string {
		match, err := regexp.Match(`^[0-2][0-9]\:[0-5][0-9]\:[0-5][0-9]$`, []byte(_param))
		if err != nil || !match {
			return PARAM_CHECK_FAIED
		}
		return _param
	}

	// 日期时间
	paramCheckers[PTYPE_DATETIME] = func(_param string, _extra *ParamInfo) string {
		match, err := regexp.Match(`^[0-9]{4}\-[01][0-9]\-([012][0-9]|[3][01])\s[0-2][0-9]\:[0-5][0-9]\:[0-5][0-9]$`, []byte(_param))
		if err != nil || !match {
			return PARAM_CHECK_FAIED
		}
		return _param
	}

	// IP地址
	paramCheckers[PTYPE_IP] = func(_param string, _extra *ParamInfo) string {
		match, err := regexp.Match(`^[0-9]{1,3}(\.[0-9]{1,3}){3}$`, []byte(_param))
		if err != nil || !match {
			return PARAM_CHECK_FAIED
		}
		return _param
	}

	// MD5
	paramCheckers[PTYPE_MD5] = func(_param string, _extra *ParamInfo) string {
		match, err := regexp.Match(`^([0-9a-zA-Z]{32})|([0-9a-zA-Z]{48})|([0-9a-zA-Z]{64})$`, []byte(_param))
		if err != nil || !match {
			return PARAM_CHECK_FAIED
		}
		return _param
	}

	// 资源文件名
	paramCheckers[PTYPE_ASSERT_NAME] = func(_param string, _extra *ParamInfo) string {
		match, err := regexp.Match(`^[0-9a-zA-Z\_\.]+\.(jpg|jpeg|png|gif|mp4|avi|ogg|txt)$`, []byte(_param))
		if err != nil || !match {
			return PARAM_CHECK_FAIED
		}
		return _param
	}

	// URL
	paramCheckers[PTYPE_URL] = func(_param string, _extra *ParamInfo) string {
		match, err := regexp.Match(`^([hH][tT]{2}[pP]:\/\/|[hH][tT]{2}[pP][sS]:\/\/)([a-zA-Z0-9\-\_\p{Han}\%\&\?\#\@\:\.\/\=])+[0-9a-zA-Z\#]+$`, []byte(_param))
		if err != nil || !match {
			return PARAM_CHECK_FAIED
		}
		return _param
	}

	// 邮箱
	paramCheckers[PTYPE_EMAIL] = func(_param string, _extra *ParamInfo) string {
		match, err := regexp.Match(`^([a-zA-Z0-9_\.\-])+\@(([a-zA-Z0-9\-])+\.)+([a-zA-Z0-9]{2,4})+$`, []byte(_param))
		if err != nil || !match {
			return PARAM_CHECK_FAIED
		}
		return _param
	}

	// 用户名
	paramCheckers[PTYPE_USERNAME] = func(_param string, _extra *ParamInfo) string {
		match, err := regexp.Match(`^([0-9a-zA-Z]{5,16})|(([a-zA-Z0-9_\.\-])+\@(([a-zA-Z0-9\-])+\.)+([a-zA-Z0-9]{2,4})+)$`, []byte(_param))
		if err != nil || !match {
			return PARAM_CHECK_FAIED
		}
		return _param
	}

	// 密码
	paramCheckers[PTYPE_PASSWORD] = func(_param string, _extra *ParamInfo) string {
		match, err := regexp.Match(`^[0-9a-zA-Z\.\,\!\@\#\$\%\^\&\*]{6,20}$`, []byte(_param))
		if err != nil || !match {
			return PARAM_CHECK_FAIED
		}
		return _param
	}

	// 数字列表
	paramCheckers[PTYPE_NUMBER_LIST] = func(_param string, _extra *ParamInfo) string {
		match, err := regexp.Match(`^[0-9]{1,20}(\,[0-9]{1,20})*$`, []byte(_param))
		if err != nil || !match {
			return PARAM_CHECK_FAIED
		}
		return _param
	}

	// 手机号码
	paramCheckers[PTYPE_PHONE] = func(_param string, _extra *ParamInfo) string {
		match, err := regexp.Match(`^(13|14|15|16|17|18|19)[0-9]{9}$`, []byte(_param))
		if err != nil || !match {
			return PARAM_CHECK_FAIED
		}
		return _param
	}

	// 短信或邮箱验证码
	paramCheckers[PTYPE_VCODE] = func(_param string, _extra *ParamInfo) string {
		match, err := regexp.Match(`^[0-9]{6}$`, []byte(_param))
		if err != nil || !match {
			return PARAM_CHECK_FAIED
		}
		return _param
	}

	// 身份证
	paramCheckers[PTYPE_ID_CARD] = func(_param string, _extra *ParamInfo) string {
		_IDRe18 := `/^([1-6][1-9]|50)\d{4}(18|19|20)\d{2}((0[1-9])|10|11|12)(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$/`
		_IDre15 := `/^([1-6][1-9]|50)\d{4}\d{2}((0[1-9])|10|11|12)(([0-2][1-9])|10|20|30|31)\d{3}$/`

		match18, err18 := regexp.Match(_IDRe18, []byte(_param))
		match15, err15 := regexp.Match(_IDre15, []byte(_param))
		if err18 != nil {
		}
		// 校验身份证：
		if (err18 == nil && match18) || (err15 == nil && match15) {
			return _param
		}
		return PARAM_CHECK_FAIED
	}

	// 范围内的整数
	paramCheckers[PTYPE_INT_RANGE] = func(_param string, _extra *ParamInfo) string {
		ib, err := strconv.ParseInt(strings.Trim(_param, " "), 10, 32)
		if err != nil {
			return PARAM_CHECK_FAIED
		}
		if _extra == nil {
			return _param
		}

		if len(_extra.Extra) != 2 {
			return _param
		}

		_min := clCommon.Int(_extra.Extra[0])
		_max := clCommon.Int(_extra.Extra[1])
		if ib >= int64(_min) && ib <= int64(_max) {
			return _param
		}
		return PARAM_CHECK_FAIED
	}

	// 长度范围内的字符串
	paramCheckers[PTYPE_STR_RANGE] = func(_param string, _extra *ParamInfo) string {
		if _extra == nil {
			return _param
		}

		if len(_extra.Extra) != 2 {
			return _param
		}

		_min := clCommon.Int(_extra.Extra[0])
		_max := clCommon.Int(_extra.Extra[1])
		if len(_param) >= _min && len(_param) <= _max {
			return _param
		}
		return PARAM_CHECK_FAIED
	}

	// 长度范围内的字符串
	paramCheckers[PTYPE_V2] = func(_param string, _extra *ParamInfo) string {

		params := strings.Split(_param, ",")
		if len(params) != 2 {
			return PARAM_CHECK_FAIED
		}

		for _, val := range params {
			match, err := regexp.Match(`^(\-)?[0-9]{1,20}(\.[0-9]{1,10})?$`, []byte(val))
			if err != nil || !match {
				return PARAM_CHECK_FAIED
			}
		}

		return _param
	}

	// 长度范围内的字符串
	paramCheckers[PTYPE_V3] = func(_param string, _extra *ParamInfo) string {

		params := strings.Split(_param, ",")
		if len(params) != 3 {
			return PARAM_CHECK_FAIED
		}

		for _, val := range params {
			match, err := regexp.Match(`^(\-)?[0-9]{1,20}(\.[0-9]{1,10})?$`, []byte(val))
			if err != nil || !match {
				return PARAM_CHECK_FAIED
			}
		}

		return _param
	}

	// 四维坐标系
	paramCheckers[PTYPE_V4] = func(_param string, _extra *ParamInfo) string {

		params := strings.Split(_param, ",")
		if len(params) != 4 {
			return PARAM_CHECK_FAIED
		}

		for _, val := range params {
			match, err := regexp.Match(`^(\-)?[0-9]{1,20}(\.[0-9]{1,10})?$`, []byte(val))
			if err != nil || !match {
				return PARAM_CHECK_FAIED
			}
		}

		return _param
	}

	// 指定正则
	paramCheckers[PTYPE_REGEX] = func(_param string, _extra *ParamInfo) string {
		match, err := regexp.Match(_extra.Extra[0], []byte(_param))
		if err != nil || !match {
			return PARAM_CHECK_FAIED
		}
		return _param
	}

	// 指定字符串区间, 多个字符串用逗号分割
	paramCheckers[PTYPE_STR_IN] = func(_param string, _extra *ParamInfo) string {
		checkStrList := strings.Split(_extra.Extra[0], "&$&")
		caseMatch := _extra.Extra[1] == "true"
		for _, str := range checkStrList {
			if !caseMatch {
				// 大小写不敏感
				if strings.ToUpper(str) == strings.ToUpper(_param) {
					return _param
				}
			} else {
				// 大小写敏感
				if str == _param {
					return _param
				}
			}
		}
		// 没有匹配成功
		return PARAM_CHECK_FAIED
	}

}

// @author xiaolan
// @lastUpdate 2019-08-04
// @comment 验证参数合法性
// @param _param 参数值
func (this *ParamInfo) CheckParam(_param string) bool {
	checkFunc, exists := paramCheckers[this.ParamType]
	if !exists {
		return false
	}
	return checkFunc(_param, this) != PARAM_CHECK_FAIED
}
