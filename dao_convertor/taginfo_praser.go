package dao_convertor

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

//CRUD internal flag
const (
	KCreate = iota
	KQuery
	KUpdate
	KDelete
	CRUDLen
)

//CRUD tag str
var tagMapping = map[string]int{"c": KCreate, "r": KQuery, "u": KUpdate, "d": KDelete}

//tag string
const kTagString = "dao"

type TagInfo struct {
	Action    int
	Where     string
	Omit      bool
	rawString string
}

func (t *TagInfo) ParseString(info string) error {
	t.rawString = info
	var err error
	if err = t.parseAction(); err != nil {
		return err
	}
	if err = t.parseFlags(); err != nil {
		return err
	}
	return nil
}

func (t *TagInfo) parseAction() error {
	var action string
	idx := strings.IndexByte(t.rawString, ':')
	if idx == -1 {
		action = t.rawString
		t.rawString = ""
	} else {
		action = t.rawString[0:idx]
		t.rawString = t.rawString[idx+1:]
	}
	action = strings.Trim(action, " \n\t")
	action = strings.ToLower(action)
	if _, ok := tagMapping[action]; ok {
		t.Action = tagMapping[action]
		return nil
	}
	return errors.New("unexpected action in TagInfo.parseAction")
}

func (t *TagInfo) parseFlags() error {
	if t.rawString == "" { //无后续数据可解析
		return nil
	}
	parts := strings.Split(t.rawString, ",")
	if len(parts) == 0 {
		return errors.New("unexpected in TagInfo.parseFlags")
	}
	for _, part := range parts {
		part = strings.Trim(part, " \r\n\t")
		if part == "" {
			continue
		}
		//解析where字段
		if part[0] == '(' && part[len(part)-1] == ')' {
			t.Where = part[1 : len(part)-1]
			continue
		}
		//是否是omit属性
		if part == "omit" {
			t.Omit = true
		}
	}
	return nil
}

type TagInfoMap map[string]*[CRUDLen]TagInfo

var (
	once       sync.Once
	tagInfoMap TagInfoMap
)

func GetTagInfoMap() *TagInfoMap {
	once.Do(func() {
		tagInfoMap = make(map[string]*[CRUDLen]TagInfo)
	})
	return &tagInfoMap
}

func getKey(structName, fieldName string) string {
	return fmt.Sprintf("%s:%s", structName, fieldName)
}

func (t *TagInfoMap) InsertTag(structName, fieldName string, tag TagInfo) {
	key := getKey(structName, fieldName)
	if _, ok := (*t)[key]; ok {
		nums := (*t)[key]
		(*nums)[tag.Action] = tag
		return
	}
	(*t)[key] = &[CRUDLen]TagInfo{}
	nums := (*t)[key]
	(*nums)[tag.Action] = tag
}

func (t *TagInfoMap) GetTag(structName, fieldName string) *[CRUDLen]TagInfo {
	key := getKey(structName, fieldName)
	if _, ok := (*t)[key]; ok {
		return (*t)[key]
	}
	return nil
}

func ParseTag(tag string) (tagNums *[CRUDLen]TagInfo, err error) {
	tagNums = new([CRUDLen]TagInfo)
	parts := strings.Split(tag, ";")
	for _, part := range parts {
		if part != "" {
			t := TagInfo{}
			err = t.ParseString(part)
			if err != nil {
				return nil, err
			}
			tagNums[t.Action] = t
		}
	}
	return
}

func ParseTagHandler(structName, fieldName, tag string) {
	tagInfos, err := ParseTag(tag)
	if err != nil {
		fmt.Printf("dao.ParseTagHandler failed in parse %s \n", tag)
	}
	(*GetTagInfoMap())[getKey(structName, fieldName)] = tagInfos
}

//已被成功解耦，具体实现在util的reflect_util.go中
//func UpdateFromStruct(v interface{}) error {
//	info := reflect.TypeOf(v).Elem()
//
//	//更新该对象需要导入的路径
//	parserCommonInfo := common_info.GetParserInfo()
//	nums := make([]common_info.IColumn, 0)
//	structTypeName := info.Name()
//	packagePath := info.PkgPath()
//	packageName := packagePath[strings.LastIndexByte(packagePath, '/')+1:]
//	(*parserCommonInfo.Dependencies)[structTypeName] = common_info.DependenceInfo{
//		PackagePath: packagePath + ",",
//		packageName: packageName,
//	}
//	for i := 0; i < info.NumField(); i++ {
//		value := info.Field(i)
//		fieldTypeName := value.Type.String()
//		fieldName := value.Name
//		tagName := value.Tag.Get(kTagString)
//		e := Element{
//			EName: fieldName,
//			EType: fieldTypeName,
//			ETag:  tagName,
//		}
//		nums = append(nums, &e)
//		//解析更新tagInfoMap
//		tagInfos, err := ParseTag(tagName)
//		if err != nil {
//			fmt.Printf("dao.ParseTagHandler failed in parse %s \n",tagName)
//		}
//		(*GetTagInfoMap())[getKey(structTypeName, fieldName)] = tagInfos
//	}
//
//	//更新objects
//	(*parserCommonInfo.Objects)[structTypeName] = nums
//
//	return nil
//}
