package clResponse


type SkyResp struct {
	Code uint32 `json:"code"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"`
}


type I18NMap map[uint32] map[uint32]string

