package controller_convertor

import (
	"fmt"
	"github.com/ACking-you/quickstart_project/common_info"
	"github.com/ACking-you/quickstart_project/controller_convertor/template"
	"github.com/ACking-you/quickstart_project/util"
	"path/filepath"
	"strings"
)

type Struct2Controller struct {
	config  *Config
	objects *common_info.ObjectInfo
	err     error
}

func NewStruct2Controller(config *Config) *Struct2Controller {
	return &Struct2Controller{
		config:  config,
		objects: common_info.GetParserInfo().Objects,
		err:     nil,
	}
}

func (s *Struct2Controller) Config(config *Config) *Struct2Controller {
	s.config = config
	return s
}

func (s *Struct2Controller) AutoMigrate(t ...interface{}) *Struct2Controller {
	for _, p := range t {
		s.err = util.UpdateFromStruct(p, "", nil) //本层的解析并不提供tag的拓展作用
	}
	return s
}

func (s *Struct2Controller) generateGinExample() string {
	return template.KGinExampleTemplate
}

func (s *Struct2Controller) generateResponse() string {
	return template.KResponseTemplate
}

func (s *Struct2Controller) generateCodeByAction(rawString, structName string) string {
	var content strings.Builder
	h := util.StrHandleByChain{Str: rawString}
	if s.config.enableCreate {
		content.WriteString(h.ReplaceAll("${action}", "Create").ReplaceAll("${class}", structName).Str)
	}
	if s.config.enableQuery {
		h.Str = rawString
		content.WriteString(h.ReplaceAll("${action}", "Query").ReplaceAll("${class}", structName).Str)
	}
	if s.config.enableUpdate {
		h.Str = rawString
		content.WriteString(h.ReplaceAll("${action}", "Update").ReplaceAll("${class}", structName).Str)
	}
	if s.config.enableDelete {
		h.Str = rawString
		content.WriteString(h.ReplaceAll("${action}", "Delete").ReplaceAll("${class}", structName).Str)
	}
	return content.String()
}

func (s *Struct2Controller) generateBind(structName string) string {
	if s.config.enableVOBind {
		return strings.ReplaceAll(template.KHandlerBind, "${class}", structName)
	}
	return ""
}

func (s *Struct2Controller) generateContentByStructName(structName string) string {
	var handlerCode string
	var bindCode string
	var defineCode string
	var controllerCode string
	if s.config.enableResponse {
		handlerCode = s.generateCodeByAction(template.KResponseHandlerByAction, structName)
	} else {
		handlerCode = s.generateCodeByAction(template.KNoResponseHandlerByAction, structName)
	}
	if s.config.enableVOBind {
		bindCode = s.generateBind(structName)
	}
	defineCode = strings.ReplaceAll(template.KControllerDefine, "${class}", structName)
	controllerCode = s.generateCodeByAction(template.KControllerAction, structName)

	return handlerCode + bindCode + defineCode + controllerCode
}

func (s *Struct2Controller) Run() error {

	imports := make([]string, 0)
	contents := make([]string, 0)
	currentNames := make([]string, 0)
	for structName := range *s.objects {
		importx := s.config.ginPackageName + ","
		p := (*common_info.GetParserInfo().Dependencies)[structName].PackagePath
		name := (*common_info.GetParserInfo().Dependencies)[structName].PackageName
		if s.config.enableVOBind {
			importx += strings.Replace(p, name, "service/vo", 1) + ","
		}
		if s.config.enableResponse {
			importx += strings.Replace(p, name, "controller/r", 1) + ","
		}
		imports = append(imports, importx)
		contents = append(contents, s.generateContentByStructName(structName))
		currentNames = append(currentNames, structName)
	}

	common_info.GetParserInfo().CurrentObjectsNameByOrder = currentNames
	if s.config.enableDebug {
		for _, c := range contents {
			fmt.Println(c)
		}
	}
	//是否生成gin的默认example
	basePath := s.config.savePath
	if !util.IsDir(s.config.savePath) {
		basePath = filepath.Dir(s.config.savePath)
	}
	if s.config.enableGinExample {
		//注意import数组在单文件会把数组所有的内容作为一个import，而多文件时单个数组元素会以 ',' 分割每个文件的import
		s.err = util.SaveAction(filepath.Join(basePath, "gin_example", "example.go"), "gin_example",
			[]string{s.config.ginPackageName, "net/http"}, []string{s.generateGinExample()})
		if s.err != nil {
			return s.err
		}
	}
	if s.config.enableResponse {
		s.err = util.SaveAction(filepath.Join(basePath, "r", "response.go"), "r",
			[]string{s.config.ginPackageName, "net/http"}, []string{s.generateResponse()})
		if s.err != nil {
			return s.err
		}
	}
	//生成
	s.err = util.SaveAction(basePath, KPackageName, imports, contents)
	if s.err != nil {
		return s.err
	}
	return nil
}
