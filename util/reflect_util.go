package util

import (
	"bufio"
	"github.com/ACking-you/quickstart_project/common_info"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
)

type Element struct {
	EName string
	EType string
	ETag  string
}

func (e *Element) Name() string {
	return e.EName
}
func (e *Element) Type() string {
	return e.EType
}
func (e *Element) Tag() string {
	return e.ETag
}

type ParseTagHandler func(structName, fieldName, tagName string)

func UpdateFromStruct(v interface{}, tagKey string, handler ParseTagHandler) error {
	info := reflect.TypeOf(v).Elem()

	//更新该对象需要导入的路径
	parserCommonInfo := common_info.GetParserInfo()
	nums := make([]common_info.IColumn, 0)
	structTypeName := info.Name()
	packagePath := info.PkgPath()
	packageName := packagePath[strings.LastIndexByte(packagePath, '/')+1:]
	(*parserCommonInfo.Dependencies)[structTypeName] = common_info.DependenceInfo{
		PackagePath: packagePath + ",",
		PackageName: packageName,
	}
	//根据该类型字段更新对应元信息
	for i := 0; i < info.NumField(); i++ {
		value := info.Field(i)
		fieldTypeName := value.Type.String()
		fieldName := value.Name
		tagName := value.Tag.Get(tagKey)
		e := Element{
			EName: fieldName,
			EType: fieldTypeName,
			ETag:  tagName,
		}
		nums = append(nums, &e)
		//解析tag
		if tagName != "" && handler != nil {
			handler(structTypeName, fieldName, tagName)
		}
	}

	//更新objects
	(*parserCommonInfo.Objects)[structTypeName] = nums

	return nil
}

const KChain = `

func (Output *${class}) ${argPublic}(t ${type}) *${class} {
	Output.${argPrivate} = t
	return Output
}

`

// GenerateMethodByChainStyle 根据struct字段自动生成对应的链式访问方法，注意生成链式访问后，原本public的字段会自动变成private
func GenerateMethodByChainStyle(t interface{}, fromPath string, toPath string) error {
	err := UpdateFromStruct(t, "", nil)
	if err != nil {
		return err
	}
	fi, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer fi.Close()

	rawBytes, err := ioutil.ReadAll(fi)
	if err != nil {
		return err
	}
	preContent := string(rawBytes)
	var content strings.Builder
	for objectName, columns := range *common_info.GetParserInfo().Objects {
		helpPre := StrHandleByChain{preContent}
		for _, arg := range columns {
			publicName := arg.Name()
			privateName := arg.Name()
			if IsUpper(privateName[0]) {
				privateName = strings.ToLower(privateName[0:1]) + privateName[1:]
			}
			if !IsUpper(publicName[0]) {
				publicName = strings.ToUpper(publicName[0:1]) + publicName[1:]
			}
			h := StrHandleByChain{Str: KChain}
			h.ReplaceAll("${class}", objectName).ReplaceAll("${argPublic}", publicName).ReplaceAll("${argPrivate}", privateName).ReplaceAll("${type}", arg.Type())
			content.WriteString(h.Str)
			helpPre.ReplaceAll(arg.Name(), privateName)
		}
		preContent = helpPre.Str
		//只解析一个同时也应该只有一个
		break
	}

	fo, err := os.Create(toPath)
	defer fo.Close()
	outBuf := bufio.NewWriter(fo)
	_, err = outBuf.WriteString(preContent)
	if err != nil {
		return err
	}
	_, err = outBuf.WriteString(content.String())
	if err != nil {
		return err
	}
	err = outBuf.Flush()
	if err != nil {
		return err
	}
	_ = DoGoFmt(toPath)
	return nil
}
