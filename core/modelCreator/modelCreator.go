package modelCreator

import (
	"fmt"
	"github.com/lionhart580230/clUtil/clCommon"
	"github.com/lionhart580230/clUtil/clFile"
	"github.com/lionhart580230/clUtil/clLog"
	"github.com/lionhart580230/clUtil/clMysql"
	"strings"
)

// 模型文件生成器
const modelTemplatePackage = `package %v

// 导入文件
import (
	"errors"
	"github.com/lionhart580230/clUtil/clLog"
	"github.com/lionhart580230/clUtil/clReflect"
	"github.com/lionhart580230/clhttp-framework/clGlobal"
)
`

// 模型区
const modelTemplateModel = `
// 表名
const TableName = "%v"
// 主键名
const PK = "%v"

type Model struct{
%v
	
	changeField map[string] interface{}
}`

// 函数区
const modelTemplateFunction = `
// 生成模型
func New() *Model {
	return &Model{
		changeField: make(map[string] interface{}),
	}
}


// 保存
func (this *Model) SaveToDB(_pkVal %v) error {
	DB := clGlobal.GetMysql()
	if DB == nil {
		return errors.New("数据库连线丢失")
	}
	_, err := DB.NewBuilder().Table(TableName).Where("%%v = '%%v'", PK, _pkVal).Save(this.changeField)
	return err
}

// 保存
func (this *Model) AddToDB() error {
	DB := clGlobal.GetMysql()
	if DB == nil {
		return errors.New("数据库连线丢失")
	}
	_, err := DB.NewBuilder().Table(TableName).Add(this.changeField)
	return err
}


// 生成自定义函数
%v
`

// @auth ciaolan
// @param _dbHost 数据库连线
// @param _dbUser 数据库用户
// @param _dbPass 数据库密码
// @param _dbName 数据库名称
// @param _path 数据库路径
func CreateAllModelFile(_dbHost, _dbUser, _dbPass, _dbName, _path string) {
	DB, err := clMysql.NewDB(_dbHost, _dbUser, _dbPass, _dbName)
	if err != nil {
		clLog.Error("连接数据库错误: %v", err)
		return
	}

	tableList, err := DB.GetTables("")
	if err != nil {
		clLog.Error("获取数据表列表错误: %v", err)
		return
	}

	clFile.CreateDirIFNotExists(_path)

	for _, table := range tableList {
		cols, err := DB.GetTableColumns(table)
		if err != nil {
			clLog.Error("获取字段列表错误: %v", err)
			continue
		}

		CreateModelFile(_path, table, cols)
	}
}

// 生成模型文件
func CreateModelFile(_path string, _tableName string, _columns []clMysql.DBColumnInfo) {
	var fileContent = strings.Builder{}
	// 生成模块名字
	var modelName = "model" + clCommon.UnderlineToUppercase(true, _tableName)

	clFile.CreateDirIFNotExists(_path + "/" + modelName)

	// 生成模型头部
	fileContent.WriteString(fmt.Sprintf(modelTemplatePackage, modelName))

	var pkString = ""
	var pkType = ""
	// 结构体
	var structBody = strings.Builder{}
	// 函数体
	var setFunctions = strings.Builder{}

	for _, val := range _columns {

		var fieldType = "string"
		if strings.HasPrefix(val.Type, "int") || strings.HasPrefix(val.Type, "tinyint") {
			if strings.Contains(val.Type, "unsigned") {
				fieldType = "uint32"
			} else {
				fieldType = "int32"
			}
		} else if strings.HasPrefix(val.Type, "bigint") {
			if strings.Contains(val.Type, "unsigned") {
				fieldType = "uint64"
			} else {
				fieldType = "int64"
			}
		} else if strings.HasPrefix(val.Type, "float") || strings.HasPrefix(val.Type, "decimal") {
			fieldType = "float64"
		}
		if val.KeyType == "PRI" {
			pkString = val.Field
			pkType = fieldType
		}

		fieldName := clCommon.UnderlineToUppercase(true, val.Field)
		// 写入结构体
		structBody.WriteString(fmt.Sprintf("\t%v %v `db:\"%v\"`\n", fieldName, fieldType, val.Field))

		setFunctions.WriteString(fmt.Sprintf(`
func (this *Model)Set%v(_val %v) *Model {
	this.%v = _val
	this.changeField["%v"] = _val
	return this
}
`, fieldName, fieldType, fieldName, val.Field))

	}

	fileContent.WriteString(fmt.Sprintf(modelTemplateModel, _tableName, pkString, structBody.String()))
	fileContent.WriteString(fmt.Sprintf(modelTemplateFunction, pkType, setFunctions.String()))

	clFile.AppendFile(_path+"/"+modelName+"/"+modelName+".go", fileContent.String())

	fmt.Printf(">> 模型文件生成完毕!!\n\n%v\n", fileContent.String())
}
