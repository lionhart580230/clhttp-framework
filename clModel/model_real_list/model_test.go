package model_real_list

import "testing"

func TestNew(t *testing.T) {
	var model = New()
	model.FillData(Model{
		Id:          1,
		ParentId:    "dasda",
		MenuName:    "asdasd",
		Url:         "asdasd",
	})
}