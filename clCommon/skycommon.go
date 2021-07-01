package clCommon

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/big"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

//@author xiaolan
//@lastUpdate 2019-08-04
//@comment 对一些常用函数进行封装



//@author xiaolan
//@lastUpdate 2019-08-04
//@comment 整数型转IP地址
func Long2IP(ip uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}


//@author xiaolan
//@lastUpdate 2019-08-04
//@comment IP地址转整数型
func IP2Long(ip string) int64 {
	ret := big.NewInt(0)
	ip_in := net.ParseIP(ip).To4()
	if ip_in == nil {
		return 0
	}
	ret.SetBytes(ip_in)
	return ret.Int64()
}


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
//@lastUpdate 2019-08-04
//@comment 字符串转int
func Int(ceil string) int {
	i, err := strconv.Atoi(strings.Trim(ceil, " "))
	if err != nil {
		return 0
	}
	return i
}

//@author xiaolan
//@lastUpdate 2019-08-04
//@comment 字符串转int8
func Int8(ceil string) int8 {
	ib, err := strconv.ParseInt(strings.Trim(ceil, " "), 10, 8)
	if err == nil {
		return int8(ib)
	}
	return int8(0)
}

//@author xiaolan
//@lastUpdate 2019-08-04
//@comment 字符串转int32
func Int32(ceil string) int32 {
	ib, err := strconv.ParseInt(strings.Trim(ceil, " "), 10, 32)
	if err == nil {
		return int32(ib)
	}
	return int32(-1)
}

//@author xiaolan
//@lastUpdate 2019-08-04
//@comment 字符串转int64
func Int64(ceil string) int64 {
	ib, err := strconv.ParseInt(strings.Trim(ceil, " "), 10, 64)
	if err == nil {
		return ib
	}
	return 0
}

//@author xiaolan
//@lastUpdate 2019-08-04
//@comment 字符串转uint8
func Uint8(ceil string) uint8 {
	ib, err := strconv.ParseUint(strings.Trim(ceil, " "), 10, 8)
	if err == nil {
		return uint8(ib)
	}
	return uint8(0)
}

//@author xiaolan
//@lastUpdate 2019-08-04
//@comment 字符串转uint32
func Uint32(ceil string) uint32 {
	ib, err := strconv.ParseUint(strings.Trim(ceil, " "), 10, 32)
	if err == nil {
		return uint32(ib)
	}
	return uint32(0)
}

//@author xiaolan
//@lastUpdate 2019-08-04
//@comment 字符串转uint64
func Uint64(ceil string) uint64 {
	ib, err := strconv.ParseUint(strings.Trim(ceil, " "), 10, 64)
	if err == nil {
		return ib
	}
	return 0
}

//@author xiaolan
//@lastUpdate 2019-08-04
//@comment 16进制转uint8
func HexUnit8(ceil string) uint8 {
	ib, err := strconv.ParseUint(strings.Trim(ceil," "), 16, 8)
	if err == nil {
		return uint8(ib)
	}
	return uint8(0)
}

//@author xiaolan
//@lastUpdate 2019-08-04
//@comment 16进制转uint32
func HexUnit32(ceil string) uint32 {
	ib, err := strconv.ParseUint(strings.Trim(ceil, " "), 16, 32)
	if err == nil {
		return uint32(ib)
	}
	return uint32(0)
}


//@author xiaolan
//@lastUpdate 2019-08-32
//@comment 字符串转float32
func Float32(ceil string) float32 {
	fb, err := strconv.ParseFloat(strings.Trim(ceil, " "), 32)
	if err == nil {
		return float32(fb)
	}
	return float32(0)
}

//@author xiaolan
//@lastUpdate 2019-08-04
//@comment 字符串转float64
func Float64(ceil string) float64 {
	fb, err := strconv.ParseFloat(strings.Trim(ceil, " "), 64)
	if err == nil {
		return fb
	}
	return float64(0)
}

//@author xiaolan
//@lastUpdate 2019-08-04
//@comment 字符串转bool
func Bool(ceil string) bool {
	ceil = strings.ToLower(ceil)
	if ceil == "true" || ceil == "yes" || ceil == "on" || Int32(ceil) > 0 {
		return true
	}
	return false
}



//@author xiaolan
//@lastUpdate 2019-08-05
//@comment 获取指定范围的整数型随机数
func RandInt(_min int64, _max int64) int64 {
	if _max == _min {
		return _min
	}
	rand.Seed(time.Now().UnixNano())
	return (rand.Int63() % (_max - _min)) + _min
}


//@author xiaolan
//@lastUpdate 2019-08-06
//@comment 生成用户token
func GenUserToken(_uid uint64) string {
	return Md5([]byte(fmt.Sprintf("UToken:%v%v%v", _uid, time.Now().Unix(), RandInt(0, 0xFFFFFFFF))))
}


//@author xiaolan
//@lastUpdate 2019-08-06
//@comment 生成用户uid
func GenUserUid() uint64 {
	// 循环500次
	uid := uint64(0)
	for i := 0; i<500; i++ {
		uid = ( uid + uint64(RandInt(0, 0xFFFFFFFFFFFF)) ) % 100000000000
	}
	return uid
}


