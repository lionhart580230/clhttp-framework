package clAuth

import (
	"fmt"
	"github.com/xiaolan580230/clhttp-framework/core/cljson"
	"reflect"
	"testing"
)

type TestObj struct {
	a string
	b string
}

func TestGetUser(t *testing.T) {
	_val := cljson.A{"a", "sss"}
	_type := reflect.TypeOf(_val)

	fmt.Printf("%+v 类型: %v\n", _val, _type.Align())
}