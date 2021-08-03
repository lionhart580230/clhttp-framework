package clCache

import (
	"fmt"
	"reflect"
	"testing"
)

func TestUpdateCache(t *testing.T) {

	var i = 1000

	switch reflect.TypeOf(i).Kind() {
		case reflect.Int:
			fallthrough
		case reflect.Float32:
			fallthrough
		case reflect.Int32:
			fallthrough
		case reflect.Int64:
			fallthrough
		case reflect.Uint32:
			fallthrough
		case reflect.Uint64:
			fallthrough
		case reflect.Float64:
			fallthrough
		case reflect.String:
			fallthrough
		case reflect.Bool:

	}

	fmt.Printf("i的类型: %v\n", reflect.TypeOf(i))

}
