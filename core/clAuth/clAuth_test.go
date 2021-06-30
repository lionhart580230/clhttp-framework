package clAuth

import (
	"fmt"
	"reflect"
	"testing"
)

type TestObj struct {
	a string
	b string
}

func TestGetUser(t *testing.T) {
	_val := "1111"
	_type := reflect.TypeOf(_val)

	fmt.Printf("%+v 类型: %v\n", _val, _type.Align())
}