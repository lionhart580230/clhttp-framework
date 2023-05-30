package httpserver

import (
	"github.com/lionhart580230/clUtil/clCrypt"
	"github.com/lionhart580230/clUtil/clHttpClient"
	"github.com/lionhart580230/clUtil/clJson"
	"github.com/lionhart580230/clUtil/clLog"
	"testing"
)

func TestCallHandler(t *testing.T) {

	var aesKey = "5d41402abc4b2a76b9719d911017c592"
	var iv = "5d41222402abc4b222a76b9719d911017c592"
	var data = clJson.CreateBy(clJson.M{"k": "hello world!!"}).ToStr()
	hc := clHttpClient.NewClient("http://127.0.0.1:19999/request")
	hc.SetContentType(clHttpClient.ContentJson)
	hc.AddHeader("Encrypt-Type", "AES")
	hc.AddHeader("Encrypt-iv", iv)
	hc.SetBody(clCrypt.AesCBCEncode(data, aesKey, iv))
	resp, err := hc.Do()
	if err != nil {
		clLog.Error("请求错误: %v", err)
		return
	}
	clLog.Debug("Resp: %+v", resp.Body)
	decodeStr := clCrypt.AesCBCDecode([]byte(resp.Body), []byte(aesKey), []byte(iv))
	clLog.Debug("Resp Decode: %+v", decodeStr)

}
