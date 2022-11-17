package service_convertor

import (
	"fmt"
	"github.com/ACking-you/quickstart_project/common_info"
	"github.com/ACking-you/quickstart_project/util"
	"path/filepath"
	"strings"
)

type Struct2Service struct {
	config     *Config
	objectInfo *common_info.ObjectInfo
	err        error
}

func NewStruct2Service(config *Config) *Struct2Service {
	return &Struct2Service{
		config:     config,
		objectInfo: common_info.GetParserInfo().Objects,
		err:        nil,
	}
}

func (s *Struct2Service) AutoMigrate(t ...interface{}) *Struct2Service {
	if s.err != nil {
		return s
	}
	for _, v := range t {
		err := util.UpdateFromStruct(v, kTagString, ParseTagHandler)
		if err != nil {
			s.err = err
		}
	}
	return s
}

func (s *Struct2Service) getTypeInfoStrings(structName string, columns []common_info.IColumn) string {
	var ret strings.Builder
	for _, column := range columns {
		tags := GetTagInfoMap().GetTag(structName, column.Name())
		//若未有tag，则该字段不进入VO
		if tags == nil {
			continue
		}
		concat := strings.Join(tags, " ")
		//若本身未打json/form tag，则根据配置决定是否自动补上
		if !strings.Contains(concat, "json") && s.config.enableAutoJsonTag {
			concat += fmt.Sprintf(" %s:\"%s\"", "json", util.PascalCase2SnackCase(column.Name()))
		}
		if !strings.Contains(concat, "form") && s.config.enableAutoJsonTag {
			concat += fmt.Sprintf(" %s:\"%s\"", "form", util.PascalCase2SnackCase(column.Name()))
		}
		//字段名 类型 tag
		ret.WriteString(fmt.Sprintf("%s %s `%s`\n", column.Name(), column.Type(), concat))
	}
	return ret.String()
}

func (s *Struct2Service) parseVOString(structName string, columns []common_info.IColumn) string {
	return fmt.Sprintf(`

type %sVO struct {
	%s
}

`, structName, s.getTypeInfoStrings(structName, columns))
}

func (s *Struct2Service) parseTOString(structName string) string {
	return fmt.Sprintf(`

type %sTO struct {

}

`, structName)
}

func (s *Struct2Service) parseServiceString(structName string) string {
	return fmt.Sprintf(`

func %sService(vo *vo.%sVO) *%sServiceHelper {
	return &%sServiceHelper{vo: vo}
}

type %sServiceHelper struct {
	vo *vo.%sVO
	to *to.%sTO
}

func (u *%sServiceHelper) Do%s() (*to.%sTO, error) {

	//将打包好的to返回
	return u.to, nil
}

`, structName, structName, structName, structName, structName,
		structName, structName, structName, s.config.defaultMethodName, structName)
}

func (s *Struct2Service) Run() error {
	if s.err != nil {
		return s.err
	}
	mapLen := len(*s.objectInfo)
	importsVO := make([]string, mapLen)
	importsTO := make([]string, mapLen)
	importsService := make([]string, mapLen)
	contentsVO := make([]string, mapLen)
	contentsTO := make([]string, mapLen)
	contentsService := make([]string, mapLen)
	currentObjects := make([]string, 0)
	idx := 0
	for objectName, column := range *s.objectInfo {
		importBase := (*common_info.GetParserInfo().Dependencies)[objectName].PackagePath
		if i := strings.LastIndexByte(importBase, '/'); i != -1 {
			importBase = importBase[:i]
		}
		//处理VO和TO
		contentsVO[idx] = s.parseVOString(objectName, column)
		contentsTO[idx] = s.parseTOString(objectName)
		if strings.Contains(contentsVO[idx], "time.Time") {
			importsVO[idx] += "time,"
		}
		if strings.Contains(contentsTO[idx], "time.Time") {
			importsTO[idx] += "time,"
		}
		//处理service
		contentsService[idx] = s.parseServiceString(objectName)
		importsService[idx] = fmt.Sprintf("%s/service/vo,%s/service/to", importBase, importBase)
		if s.config.enableDebug {
			fmt.Println(contentsVO[idx], contentsTO[idx], contentsService[idx])
		}
		currentObjects = append(currentObjects, objectName)
		idx++
	}
	common_info.GetParserInfo().CurrentObjectsNameByOrder = currentObjects

	basePath := s.config.savePath
	if !util.IsDir(basePath) {
		basePath = filepath.Dir(s.config.savePath)
	}
	voSavePath := filepath.Join(basePath, "vo")
	toSavePath := filepath.Join(basePath, "to")
	//根据config配置vo或者to是否单文件存储
	if s.config.enableVOFileSingle {
		voSavePath = filepath.Join(voSavePath, "vo_entity.go")
	}
	if s.config.enableTOFileSingle {
		toSavePath = filepath.Join(toSavePath, "to_entity.go")
	}
	//保存service
	_ = util.SaveAction(basePath, kTagString, importsService, contentsService)

	//保存vo和to
	_ = util.SaveAction(voSavePath, "vo", importsVO, contentsVO)
	_ = util.SaveAction(toSavePath, "to", importsVO, contentsTO)
	return s.err
}
