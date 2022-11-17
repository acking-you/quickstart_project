package util

import (
	"fmt"
	"github.com/ACking-you/quickstart_project/common_info"
	"os"
	"path/filepath"
	"strings"
)

var suffixMap = map[string]string{
	"dao":        "dao",
	"service":    "service",
	"vo":         "vo",
	"to":         "to",
	"handler":    "handler",
	"controller": "controller",
}

func EnsureFileDirExist(filedir string) error {
	_, err := os.Stat(filedir)
	if os.IsExist(err) {
		return nil
	}
	if os.IsNotExist(err) {
		return os.MkdirAll(filedir, 0777)
	}
	return err
}

// IsDir 通过是否含有后缀名来判断是否为文件夹
func IsDir(filePath string) bool {
	return filepath.Ext(filePath) == ""
}

func SaveFile(filepath string, content string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	_, err = file.WriteString(content)
	if err != nil {
		return err
	}
	//执行gofmt把代码格式化
	_ = DoGoFmt(filepath)
	if EnableGoTidy {
		_ = DoGoModTidy()
	}
	return nil
}

func SaveAction(savePath, packName string, importsInfo, contents []string) error {
	_, err := os.Stat(savePath)
	helper := &saveHelper{
		savePath:    savePath,
		packName:    packName,
		importsInfo: importsInfo,
		contents:    contents,
	}
	if os.IsNotExist(err) {
		path := savePath
		if !IsDir(path) {
			path = filepath.Dir(path)
		}
		err = EnsureFileDirExist(path)
		if err != nil {
			return err
		}
	}
	if IsDir(savePath) {
		return helper.multiFileSave()
	}
	return helper.singleFileSave()
}

type saveHelper struct {
	savePath, packName    string
	importsInfo, contents []string
}

//根据全局信息获取文件名，并强制使用snack_case
func getFilesName() []string {
	var ret []string
	var objectsName = common_info.GetParserInfo().CurrentObjectsNameByOrder
	if objectsName == nil || len(objectsName) == 0 {
		return nil
	}
	for _, name := range objectsName {
		ret = append(ret, PascalCase2SnackCase(name))
	}
	return ret
}

func importHandler(importName string) string {
	if importName == "" {
		return ""
	}
	return fmt.Sprintf("import \"%s\"\n", importName)
}

func importsHandler(imports []string) string {
	set := Set[string]{}
	var builder strings.Builder
	for _, str := range imports {
		if !set.Contains(str) {
			builder.WriteString(importHandler(str))
			set.Insert(str)
		}
	}
	return builder.String()
}

func packageHandler(packageName string) string {
	if packageName == "" {
		return ""
	}
	return fmt.Sprintf("package %s\n\n", packageName)
}

func contentsHandler(contents []string) string {
	return strings.Join(contents, "\n")
}

func (s *saveHelper) singleFileSave() error {
	//若包名为语义化，则作为文件名
	if filepath.Ext(s.savePath) == "" {
		if v, ok := suffixMap[s.packName]; ok {
			s.savePath = v + ".go"
		}
	}

	return SaveFile(s.savePath,
		packageHandler(s.packName)+importsHandler(s.importsInfo)+contentsHandler(s.contents))
}

func (s *saveHelper) multiFileSave() error {
	filesName := getFilesName()
	if len(filesName) != len(s.contents) {
		panic("len(filesName) != len(Output.contents)")
	}
	var err error
	for i := 0; i < len(filesName); i++ {
		filepath := filepath.Join(s.savePath, filesName[i])
		var builder strings.Builder
		builder.WriteString(packageHandler(s.packName))
		importsString := strings.Split(s.importsInfo[i], ",")
		builder.WriteString(importsHandler(importsString))
		builder.WriteString(s.contents[i])

		//若包名为语义化，则作为文件名的一部分
		if v, ok := suffixMap[s.packName]; ok {
			filepath += "_" + v
		}
		if err = SaveFile(filepath+".go", builder.String()); err != nil {
			fmt.Println(err)
		}
	}
	return err
}
