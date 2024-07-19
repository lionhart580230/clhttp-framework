package clAuth

import (
	"encoding/json"
	"github.com/lionhart580230/clUtil/clCrypt"
	"github.com/lionhart580230/clUtil/clJson"
	"github.com/lionhart580230/clUtil/clTime"
	"sync"
)

var jwtKey string = "003804fc018ad3c4bb4c76d99920b796"
var jwtIv string = "003804fc018ad3c4bb4c76d99920b796"

type JwtAuthResp struct {
	Uid        uint64            `json:"u"`
	ExpireTime uint32            `json:"e"`
	Data       map[string]string `json:"d"`
}

// 设置JWT的验证key和iv
func SetJWTCert(_key string, _iv string) {
	jwtKey = _key
	jwtIv = _iv
}

// 通过jwt来创建验证规则
func CreateAuthByJWT(_jwtStr string) *AuthInfo {
	var jwtJsonStr = clCrypt.AesCBCDecode([]byte(_jwtStr), []byte(jwtKey), []byte(jwtIv))
	if jwtJsonStr == "" {
		return nil
	}
	data := JwtAuthResp{}
	err := json.Unmarshal([]byte(jwtJsonStr), &data)
	if err != nil {
		return nil
	}
	// 过期
	if data.ExpireTime < clTime.GetNowTime() {
		return nil
	}

	return &AuthInfo{
		Uid:        data.Uid,
		Token:      "",
		LastUptime: data.ExpireTime,
		IsLogin:    true,
		ExtraData:  data.Data,
		SessionKey: "",
		mLocker:    sync.RWMutex{},
	}
}

// 创建JWT
func CreateJWT(_uid uint64, _expire uint32, _data map[string]string) string {
	dataStr := clJson.CreateBy(JwtAuthResp{
		Uid:        _uid,
		Data:       _data,
		ExpireTime: clTime.GetNowTime() + _expire,
	}).ToStr()
	return clCrypt.AesCBCEncode(dataStr, jwtKey, jwtIv)
}
