package model_convertor

import (
	"fmt"
	"github.com/ACking-you/quickstart_project/common_info"
	"github.com/ACking-you/quickstart_project/util"
)

type column struct {
	ColumnName    string
	ColumnType    string
	Nullable      string
	TableName     string
	ColumnComment string
	ColumnTag     string
}

func (c *column) Name() string {
	return c.ColumnName
}

func (c *column) Type() string {
	return c.ColumnType
}

func (c *column) Tag() string {
	return c.ColumnTag
}

// 根据数据库内容获取结构体中列的映射
func (t *Table2Struct) getColumns() (tableColumns map[string][]*column, err error) {
	//在解析完成后，同时更新有效的解析信息用于共享
	tableColumns = make(map[string][]*column)

	// sql
	sqlStr := `SELECT COLUMN_NAME,DATA_TYPE,IS_NULLABLE,TABLE_NAME,COLUMN_COMMENT
		FROM information_schema.COLUMNS 
		WHERE table_schema = DATABASE()`

	// 是否指定了具体的table
	if t.config.table != "" {
		sqlStr += fmt.Sprintf(" AND TABLE_NAME = '%s'", t.config.table)
	}
	// sql排序
	sqlStr += " ORDER BY TABLE_NAME asc, ORDINAL_POSITION asc"

	rows, err := t.db.Query(sqlStr)
	if err != nil {
		fmt.Println("Error reading table information: ", err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		col := column{}
		err = rows.Scan(&col.ColumnName, &col.ColumnType, &col.Nullable, &col.TableName, &col.ColumnComment)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		col.ColumnTag = col.ColumnName //用于给gorm映射的tag名称，需要和数据库field对应
		col.ColumnName = util.SnackCase2PascalCase(col.ColumnName)
		col.ColumnType = typeForMysqlToGo[col.ColumnType]
		// 是否需要将json_tag转换成snack_case
		jsonTag := col.ColumnTag
		jsonTag = util.PascalCase2SnackCase(jsonTag)
		if t.config.enableJsonTag {
			col.ColumnTag = fmt.Sprintf("`%s:\"%s\" json:\"%s\"`", t.config.tagKey, col.ColumnTag, jsonTag)
		} else {
			col.ColumnTag = fmt.Sprintf("`%s:\"%s\"`", t.config.tagKey, col.ColumnTag)
		}
		if _, ok := tableColumns[col.TableName]; !ok {
			tableColumns[col.TableName] = []*column{}
		}
		tableColumns[col.TableName] = append(tableColumns[col.TableName], &col)
	}

	//这里赋值object数据方便后面的数据解析
	common_info.GetParserInfo().Objects = common_info.CastObjectInfo[*column](tableColumns)
	return
}
