package skycrypt

import (
	"fmt"
	"testing"
)

var cryptKey = "0e67f635e4bbeb74b7609def8dab6e6e"

func TestDecode(t *testing.T) {

	var cryptStr = "eyJpdiI6Ik1rdFZSMHA1WkRZNFoyRkRaVk16WXc9PSIsInZhbHVlIjoiSnhFRmVoN1dPYUxaYnNFRExVNEwxWVdRNE0xOVpWTE1aak5DOEd6TWl5SHB0MWJCUjVMYkxTT0NlMksvSm4wb3BkSU5uZ3JLbmNPTThMd0pDWGZid3hUNFJDRkpSSktDRG41elpEQjVnWkk9In0="
	//Encode([]byte("hello world"), []byte(skyCommon.Md5([]byte("11asdasd23123"))))
	fmt.Printf(">> Value: (%v)\n", string(Decode([]byte(cryptStr), []byte(cryptKey))))
}

func TestRandomBlock(t *testing.T) {

	fmt.Printf(">> CryptValue: (%v)\n", string(Encode([]byte("hello world"), []byte(cryptKey))))
}

func TestEncode(t *testing.T) {
	//fmt.Printf(">> CruptValue: (%v)\n", string(Encode([]byte("ac=k8sHeathCheck"), []byte("0e67f635e4bbeb74b7609def8dab6e6e"))))
	fmt.Printf(">> CruptValue: (%v)\n", string(Encode([]byte("ac=captcha&width=200&height=80"), []byte("0e67f635e4bbeb74b7609def8dab6e6e"))))
	fmt.Printf(">> CruptValue: (%v)\n", string(Encode([]byte("width=200&height=80&captcha_id=SsvuvJDBQ7BfykolyZN3"), []byte("0e67f635e4bbeb74b7609def8dab6e6e"))))
	fmt.Printf(">> CruptValue: (%v)\n", string(Encode([]byte("ac=sendPhoneVCode&phone_code=%2B86&phone=15982418951&login=0&captcha_id=oqi4g4HG3YOlW03uKl4S&captcha_input=887033"), []byte("0e67f635e4bbeb74b7609def8dab6e6e"))))
}

func TestDecode2(t *testing.T) {
	fmt.Printf(">> DecodeValue: (%v)\n", string(Decode([]byte("eyJpdiI6InFiUGNNNEhMY2FwQXpkMHUiLCJ2YWx1ZSI6ImdtVHBLcHZ5bVpvTWdQUGlqcXRpL3JjRC9xNmhtNWhiQmVaSXZEM2tJb295Q2VqckN3N3k1dm1wQ1lWYmpxQWQifQ=="), []byte("0e67f635e4bbeb74b7609def8dab6e6e"))))
}
