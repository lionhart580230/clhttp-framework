package model_real_list

import (
	"errors"
	"github.com/xiaolan580230/clUtil/clLog"
	"github.com/xiaolan580230/clUtil/clReflect"
	"github.com/xiaolan580230/clhttp-framework/clGlobal"
)

const TableName = "real_list"
const PK = "id"

type Model struct {
	Id uint32 `db:"id" primary:"true"`
	ParentId string `db:"parent_id"`
	MenuName string `db:"menu_name"`
	Url string `db:"url"`

	changeFiled map[string] interface{}
}


func New() *Model {
	return &Model{
		changeFiled: make(map[string] interface{}),
	}
}


// 保存
func (this *Model) SaveToDB(_pkVal uint32) error {
	DB := clGlobal.GetMysql()
	if DB == nil {
		return errors.New("数据库连线丢失")
	}
	_, err := DB.NewBuilder().Table(TableName).Where("`%v` = '%v'", PK, _pkVal).Save(this.changeFiled)
	return err
}

// 保存
func (this *Model) AddToDB() error {
	DB := clGlobal.GetMysql()
	if DB == nil {
		return errors.New("数据库连线丢失")
	}
	_, err := DB.NewBuilder().Table(TableName).Add(this.changeFiled)
	return err
}

// 通过id找数据
func (this *Model) FindById(_id uint32) (*Model, error) {
	DB := clGlobal.GetMysql()
	if DB == nil {
		return nil, errors.New("数据库连线丢失")
	}
	var model = New()
	err := DB.NewBuilder().Table(TableName).Where("id = %v", _id).FindOne(model)
	return model, err
}

// 通过id找数据
func (this *Model) GetListByLimit(_page int32, _count int32) (*Model, error) {
	DB := clGlobal.GetMysql()
	if DB == nil {
		return nil, errors.New("数据库连线丢失")
	}
	var model = New()
	err := DB.NewBuilder().Table(TableName).Order("id ASC").Page(_page, _count).FindOne(model)
	return model, err
}


// 设置各自的值
func (this *Model)SetId(_id uint32) *Model {
	this.Id = _id
	this.changeFiled["id"] = _id
	return this
}

// 设置各自的值
func (this *Model)SetParent(_parent string) *Model {
	this.ParentId = _parent
	this.changeFiled["parent_id"] = _parent
	return this
}

// 设置menuname
func (this *Model)SetMenuName(_menuName string) *Model {
	this.MenuName = _menuName
	this.changeFiled["menu_name"] = _menuName
	return this
}

// 设置menuname
func (this *Model)SetUrl(_url string) *Model {
	this.Url = _url
	this.changeFiled["url"] = _url
	return this
}

// 填充整个数据
func (this *Model)FillData(_model Model) *Model {
	this = &_model

	this.changeFiled = clReflect.StructToMap(_model)
	clLog.Debug("changeFiled: %+v", this.changeFiled)
	return this
}