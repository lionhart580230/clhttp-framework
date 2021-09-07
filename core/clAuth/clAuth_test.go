package clAuth

import (
	"fmt"
	"github.com/xiaolan580230/clUtil/clJson"
	"reflect"
	"testing"
)

type TestObj struct {
	a string
	b string
}

func TestGetUser(t *testing.T) {
	_val := clJson.A{"a", "sss"}
	_type := reflect.TypeOf(_val)

	fmt.Printf("%+v 类型: %v\n", _val, _type.Align())
}