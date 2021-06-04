package clCrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)


// 微信相关操作
func pkcs7Unpad(b []byte, blocksize int) ([]byte, error) {
	if blocksize <= 0 {
		return nil, errors.New("尺寸非法")
	}
	if b == nil || len(b) == 0 {
		return nil, errors.New("尺寸非法")
	}
	if len(b)%blocksize != 0 {
		return nil, errors.New("尺寸非法")
	}
	c := b[len(b)-1]
	n := int(c)
	if n == 0 || n > len(b) {
		return nil, errors.New("尺寸非法")
	}
	for i := 0; i < n; i++ {
		if b[len(b)-n+i] != c {
			return nil, errors.New("尺寸非法")
		}
	}
	return b[:len(b)-n], nil
}



func base64Decode(data []byte) []byte {
	lenOfData := len(data)
	if lenOfData%4 > 0{
		data = append(data, bytes.Repeat([]byte("="), 4-(lenOfData%4))...)
	}
	srcByte, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil{
		return []byte{}
	}

	return srcByte
}


// AES解密数据
func AesDecode(_buffer []byte, _key []byte, _iv []byte) string {

	_buffer = bytes.ReplaceAll(_buffer, []byte(" "), []byte("+"))
	_iv = bytes.ReplaceAll(_iv, []byte(" "), []byte("+"))

	dData := base64Decode(_buffer)
	bufferResp := make([]byte, 0)


	dIv := base64Decode( _iv )
	dKey := base64Decode( _key )
	_cipher, err := aes.NewCipher( dKey )
	if err != nil {
		return ""
	}

	blockMode := cipher.NewCBCDecrypter(_cipher, dIv)

	origData := make([]byte, len(dData))
	blockMode.CryptBlocks(origData, dData)

	// 解析
	bufferResp, err = pkcs7Unpad(origData, aes.BlockSize)
	if len(bufferResp) == 0 {
		return ""
	}

	return string(bufferResp)
}