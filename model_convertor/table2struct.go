package model_convertor

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/ACking-you/quickstart_project/common_info"
	"github.com/ACking-you/quickstart_project/util"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

//mapping  mysql type to golang types
var typeForMysqlToGo = map[string]string{
	"int":                "int",
	"integer":            "int",
	"tinyint":            "int",
	"smallint":           "int",
	"mediumint":          "int",
	"bigint":             "int",
	"int unsigned":       "int",
	"integer unsigned":   "int",
	"tinyint unsigned":   "int",
	"smallint unsigned":  "int",
	"mediumint unsigned": "int",
	"bigint unsigned":    "int",
	"bit":                "int",
	"bool":               "bool",
	"enum":               "string",
	"set":                "string",
	"varchar":            "string",
	"char":               "string",
	"tinytext":           "string",
	"mediumtext":         "string",
	"text":               "string",
	"longtext":           "string",
	"blob":               "string",
	"tinyblob":           "string",
	"mediumblob":         "string",
	"longblob":           "string",
	"date":               "time.Time",
	"datetime":           "time.Time",
	"timestamp":          "time.Time",
	"time":               "time.Time",
	"float":              "float64",
	"double":             "float64",
	"decimal":            "float64",
	"binary":             "string",
	"varbinary":          "string",
}

type Table2Struct struct {
	db     *sql.DB
	config *Config
	err    error //内部错误
}

func NewTable2Struct(config *Config) *Table2Struct {
	return &Table2Struct{
		config: config,
	}
}

func (t *Table2Struct) Run() error {
	// 连接mysql, 获取db对象
	t.dialMysql()
	if t.err != nil {
		return t.err
	}

	// 获取表和字段的schema
	tableColumns, err := t.getColumns()
	if err != nil {
		return err
	}

	// 包名
	var packageName string
	if t.config.packageName == "" {
		packageName = "models"
	} else {
		packageName = t.config.packageName
	}

	// 组装struct
	structsContent := make([]string, len(tableColumns))
	importsContent := make([]string, len(tableColumns))
	currentObjectsName := make([]string, 0)

	var idx = 0
	for tableRealName, item := range tableColumns {
		tableName := tableRealName
		//filter hook
		if t.config.tableFilterHook != nil && t.config.tableFilterHook(tableName) {
			continue
		}
		var structContent string
		switch len(tableName) {
		case 0:
		case 1:
			tableName = strings.ToUpper(tableName[0:1])
		default:
			tableName = util.SnackCase2PascalCase(tableName)
		}
		depth := 1
		structContent += "type " + tableName + " struct {\n"
		for _, v := range item {
			// 字段注释
			var columnComment string
			if v.ColumnComment != "" {
				columnComment = fmt.Sprintf(" // %s", v.ColumnComment)
			}
			structContent += fmt.Sprintf("%s%s %s %s%s\n",
				tab(depth), v.ColumnName, v.ColumnType, v.ColumnTag, columnComment)
		}
		structContent += tab(depth-1) + "}\n\n"

		// 添加 method 获取真实表名
		if t.config.realTableNameMethod != "" {
			structContent += fmt.Sprintf("func (*%s) %s() string {\n",
				tableName, t.config.realTableNameMethod)
			structContent += fmt.Sprintf("%sreturn \"%s\"\n",
				tab(depth), tableRealName)
			structContent += "}\n\n"
		}
		//debug打印
		if t.config.enableDebug {
			fmt.Println(structContent)
		}
		importContent := ""
		// 如果有引入 time.Time, 则需要引入 time 包
		if strings.Contains(structContent, "time.Time") {
			importContent = "time"
		}
		structsContent[idx] = structContent
		importsContent[idx] = importContent
		currentObjectsName = append(currentObjectsName, tableRealName)
		//更新依赖信息
		(*common_info.GetParserInfo().Dependencies)[util.SnackCase2PascalCase(tableName)] = common_info.DependenceInfo{
			PackagePath: common_info.ProjectPackageName + "/" + packageName,
			PackageName: packageName,
		}
		idx++
	}
	common_info.GetParserInfo().CurrentObjectsNameByOrder = currentObjectsName
	// 写入文件struct
	err = util.SaveAction(t.config.savePath, packageName, importsContent, structsContent)
	if err != nil {
		return err
	}
	return nil
}

func (t *Table2Struct) dialMysql() {
	if t.db == nil {
		if t.config.dsn == "" {
			t.err = errors.New("dsn数据库配置缺失")
			return
		}
		t.db, t.err = sql.Open("mysql", t.config.dsn)
	}
	return
}

func tab(depth int) string {
	return strings.Repeat("\t", depth)
}
