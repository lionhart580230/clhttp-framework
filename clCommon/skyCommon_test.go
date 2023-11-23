package clCommon

import (
	"fmt"
	"testing"
)

func TestGenUserPassword(t *testing.T) {
	fmt.Printf("%v\n", GenUserPassword("admin001", "qweqwe001", false))
	fmt.Printf("%v\n", GenUserPassword("admin002", "qweqwe002", false))
	fmt.Printf("%v\n", GenUserPassword("admin003", "qweqwe003", false))

}
