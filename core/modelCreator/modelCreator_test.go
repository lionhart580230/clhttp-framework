package modelCreator

import (
	"testing"
)




func TestCreateModelFile(t *testing.T) {
	CreateAllModelFile("127.0.0.1", "root", "root", "miner_new", "test2")
}
