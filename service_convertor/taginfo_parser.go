package service_convertor

import (
	"fmt"
	"github.com/ACking-you/quickstart_project/util"
	"strings"
	"sync"
)

const kTagString = "service"

type TagInfo []string //此tagInfo存储大量以binding(required)类似的语法直接为VO添加tag

type TagInfoMap map[string]TagInfo

var (
	once       sync.Once
	tagInfoMap TagInfoMap
)

func GetTagInfoMap() *TagInfoMap {
	once.Do(func() {
		tagInfoMap = make(map[string]TagInfo)
	})
	return &tagInfoMap
}

func getKey(structName, fieldName string) string {
	return fmt.Sprintf("%s:%s", structName, fieldName)
}

func (t *TagInfoMap) InsertTag(structName, fieldName string, tag TagInfo) {
	key := getKey(structName, fieldName)
	(*t)[key] = tag
}

func (t *TagInfoMap) GetTag(structName, fieldName string) TagInfo {
	key := getKey(structName, fieldName)
	if _, ok := (*t)[key]; ok {
		return (*t)[key]
	}
	return nil
}

func ParseTag(tag string) (tagInfo TagInfo, err error) {
	parts := strings.Split(tag, ";")
	for _, part := range parts {
		part = strings.Trim(part, " \n\t\r")
		var tagName, tagContent string
		err = util.Sscanf(part, "$($)", &tagName, &tagContent)
		if err != nil {
			fmt.Println(err)
		}
		if tagName != "" && tagContent != "" {
			tagInfo = append(tagInfo, fmt.Sprintf("%s:\"%s\"", tagName, tagContent))
		}
	}
	return
}

func ParseTagHandler(structName, fieldName, tag string) {
	tagInfo, err := ParseTag(tag)
	if err != nil {
		fmt.Printf("dao.ParseTagHandler failed in parse %s \n", tag)
	}
	(*GetTagInfoMap())[getKey(structName, fieldName)] = tagInfo
}
