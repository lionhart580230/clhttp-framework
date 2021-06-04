package clCrypt

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
)

// 对一些加密方式进行封装
// 包括 MD5, Base64, AES等加密手段

//@author xiaolan
//@lastUpdate 2019-08-04
//@comment MD5加密
func Md5(str []byte) string {
	h := md5.New()
	h.Write(str) // 需要加密的字符串为 123456
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}


//@author xiaolan
//@lastUpdate 2021-06-04
//@comment base64加密
//@param str 需要进行base64加密的字符串
func Base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}


//@author xiaolan
//@lastUpdate 2021-06-04
//@comment base64解密
//@param str 需要进行base64解密的字符串
func Base64Decode(str string) []byte {
	res, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return []byte{}
	}
	return res
}



