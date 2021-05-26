package skycrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

// 服务端解密程序
// 使用 AES 加密方式

type SkyPacket struct {
	Iv    string `json:"iv"`
	Value string `json:"value"`
}

//@author xiaolan
//@lastUpdate 2019-09-29
//@comment 解密程序
func Decode(_buffer []byte, _key []byte) []byte {
	dData := base64Decode([]byte(_buffer))
	var Packet SkyPacket

	bufferResp := make([]byte, 0)
	// 解析失败
	err := json.Unmarshal(dData, &Packet)
	if err != nil {
		return nil
	}

	dIv := base64Decode([]byte(Packet.Iv))
	if len(dIv) != 16 {
		return nil
	}
	dValue := base64Decode([]byte(Packet.Value))
	if len(dValue) == 0 {
		return nil
	}

	_cipher, err := aes.NewCipher([]byte(_key))
	if err != nil {
		return nil
	}

	blockMode := cipher.NewCBCDecrypter(_cipher, []byte(dIv))

	origData := make([]byte, len(dValue))
	blockMode.CryptBlocks(origData, dValue)
	// 解析
	bufferResp = PKCS5UnPadding(origData)

	return bufferResp
}

//@author xiaolan
//@lastUpdate 2019-09-29
//@comment 加密程序
func Encode(_buffer []byte, _key []byte) []byte {

	_cipher, err := aes.NewCipher([]byte(_key))
	if err != nil {
		return nil
	}

	iv := RandomBlock()
	value := _buffer

	// 生成CBC加密对象
	blockMode := cipher.NewCBCEncrypter(_cipher, iv)

	// 填充字节
	dValue := PKCS5Padding(value)
	origData := make([]byte, len(dValue))

	// 加密
	blockMode.CryptBlocks(origData, dValue)

	// 组装结构
	var skyPacket = SkyPacket{
		Iv:    string(base64Encode(iv)),
		Value: string(base64Encode(origData)),
	}

	// 生成json字串
	finallyData, err := json.Marshal(skyPacket)
	if err != nil {
		fmt.Printf(">> 生成json字符串失败! 错误: %v\n", err)
		return nil
	}

	// base64加密
	return []byte(base64Encode(finallyData))
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func PKCS5Padding(_origData []byte) []byte {
	padnum := 16 - len(_origData)%16
	pad := bytes.Repeat([]byte{byte(padnum)}, padnum)
	return append(_origData, pad...)
}

// Base64解密
func base64Decode(data []byte) []byte {

	lenOfData := len(data)
	if lenOfData%4 > 0 {
		data = append(data, bytes.Repeat([]byte("="), 4-(lenOfData%4))...)
	}
	srcByte, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return []byte{}
	}

	return srcByte
}

// Base64加密
func base64Encode(data []byte) []byte {

	srcByte := base64.StdEncoding.EncodeToString(data)

	return []byte(srcByte)
}

// 生成16位随机字符串
func RandomBlock() []byte {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ/"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 16; i++ {
		result = append(result, bytes[r.Int31n(int32(len(str)))])
	}
	return result
}
