package common_info

import (
	"strings"
	"sync"
	"unicode"
)

var ProjectPackageName = "" //用于一键生成能够正确的import路径，毕竟一键根据数据库生成则没法利用反射获取信息了

type IColumn interface {
	Name() string
	Type() string
	Tag() string
}

type DependenceInfo struct {
	PackagePath string //以 ',' 分割多个
	PackageName string //每个类只会有一个包名
}

type ObjectInfo map[string][]IColumn

type ObjectDependence map[string]DependenceInfo

func (o *ObjectDependence) GetPackagePath(objectName string) string {
	return (*o)[objectName].PackagePath
}

func (o *ObjectDependence) GetPackageName(objectName string) string {
	return (*o)[objectName].PackageName
}

type ParserInfo struct {
	Objects                   *ObjectInfo
	Dependencies              *ObjectDependence
	CurrentObjectsNameByOrder []string
}

var (
	once  sync.Once
	infos ParserInfo
)

func GetParserInfo() *ParserInfo {
	once.Do(func() {
		if infos.Objects == nil {
			infos.Objects = new(ObjectInfo)
			*infos.Objects = make(map[string][]IColumn)
		}
		if infos.Dependencies == nil {
			infos.Dependencies = new(ObjectDependence)
			*infos.Dependencies = make(map[string]DependenceInfo)
		}
	})
	return &infos
}

// CastObjectInfo copy from any type of implement IColumn
func CastObjectInfo[T IColumn](src map[string][]T) *ObjectInfo {
	ret := make(ObjectInfo)
	for objName, values := range src {
		columns := make([]IColumn, 0)
		for _, column := range values {
			columns = append(columns, column)
		}
		ret[snackCase2PascalCase(objName)] = columns
	}
	return &ret
}

func snackCase2PascalCase(data string) string {
	var ret []rune
	data = strings.ToUpper(data[:1]) + data[1:]
	var flag bool
	for _, ch := range data {
		if ch == '_' {
			flag = true
			continue
		}
		if flag {
			ch = unicode.ToUpper(ch)
			flag = false
		}
		ret = append(ret, ch)
	}
	return string(ret)
}
