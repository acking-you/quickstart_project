package dao_convertor

import (
	"fmt"
	"github.com/ACking-you/quickstart_project/common_info"
	"github.com/ACking-you/quickstart_project/util"
	"path/filepath"
	"strings"
)

func parseDBInit(dsn, packageName, ormPackageName, sqlPackageName string) string {
	return fmt.Sprintf(`
package %s

import (
	"%s"
	"%s"
)

var db *gorm.DB

func InitDB() (err error) {
	dsn := "%s"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return
}
`, packageName, sqlPackageName, ormPackageName, dsn)
}

func parseDAONewStruct(structName string) string {
	public := structName
	private := strings.ToLower(structName[0:1]) + structName[1:]
	return fmt.Sprintf(`

type %sDAO struct {
}

var %sDAO = %sDAO{}

func New%sDAO() *%sDAO {
	return &%sDAO
}

`, public, private, public, public, public, private)
}

func getStringByParams(enableQuota bool, params []string) string {
	if params == nil || len(params) == 0 {
		return ""
	}
	var builder strings.Builder
	for idx, param := range params {
		if enableQuota {
			builder.WriteString("\"")
		}
		builder.WriteString(param)
		if enableQuota {
			builder.WriteString("\"")
		}
		if idx == len(params)-1 {
			break
		}
		builder.WriteByte(',')
	}
	return builder.String()
}

func getOmitString(params []string) string {
	if params == nil || len(params) == 0 {
		return ""
	}
	return fmt.Sprintf(".Omit(%s)", getStringByParams(true, params))
}

func getWhereString(isOr bool, wheres []string, params []string) string {
	if wheres == nil || len(wheres) == 0 || params == nil || len(params) == 0 {
		return ""
	}
	var whereBuilder strings.Builder
	connStr := "AND"
	if isOr {
		connStr = "OR"
	}
	for idx, where := range wheres {
		whereBuilder.WriteByte(' ')
		whereBuilder.WriteString(where)
		if idx == len(wheres)-1 {
			break
		}
		whereBuilder.WriteByte(' ')
		whereBuilder.WriteString(connStr)
	}
	return fmt.Sprintf(".Where(\"%s\",%s)", whereBuilder.String(), getStringByParams(false, params))
}

func parseCreate(structName string, omits []string) string {
	omit := getOmitString(omits)
	fullName := common_info.GetParserInfo().Dependencies.GetPackageName(structName) + "." + structName
	return fmt.Sprintf(`

func (%sDAO) Create%s(info *%s) error {
	return db%s.Create(info).Error
}

func (%sDAO) CreateAll%s(infos *[]%s) error {
	return db%s.Create(infos).Error
}

`, structName, structName, fullName, omit,
		structName, structName, fullName, omit)
}

func parseQuery(structName string, omits []string) string {

	omitsString := getOmitString(omits)

	fullName := common_info.GetParserInfo().Dependencies.GetPackageName(structName) + "." + structName
	return fmt.Sprintf(`

func (%sDAO) Query%s(cond, result *%s) error {
	return db.Model(cond)%s.First(result).Error
}

func (%sDAO) QueryAll%s(cond, results *[]%s) error {
	return db.Model(cond)%s.Find(results).Error
}

`,
		structName, structName, fullName, omitsString,
		structName, structName, fullName, omitsString)
}

func parseUpdate(structName string, omits []string) string {
	omitsString := getOmitString(omits)

	fullName := common_info.GetParserInfo().Dependencies.GetPackageName(structName) + "." + structName
	return fmt.Sprintf(`
	
func (%sDAO) Update%s(old, new *%s) error {
	return db.Model(old)%s.Updates(*new).Error
}

`, structName, structName, fullName, omitsString)
}

func parseDelete(structName string, isOr bool, wheres, whereParams []string) string {
	whereString := getWhereString(isOr, wheres, whereParams)
	fullName := common_info.GetParserInfo().Dependencies.GetPackageName(structName) + "." + structName
	return fmt.Sprintf(`

func (%sDAO) Delete%s(cond *%s) error {
	return db%s.Delete(cond).Error
}

`, structName, structName, fullName, whereString)
}

//检查对象的每个属性中tagOperation(对应CURD操作)操作是否被Omit，返回被Omit的属性名(强制使用snack_case)
func getOmitNumsString(tagOperation int, structName string, columns []common_info.IColumn) []string {
	numsString := make([]string, 0)
	for _, column := range columns {
		tagInfos := GetTagInfoMap().GetTag(structName, column.Name())
		if tagInfos == nil {
			continue
		}
		tagInfo := tagInfos[tagOperation]
		if tagInfo.Omit {
			numsString = append(numsString, util.PascalCase2SnackCase(column.Name()))
		}
	}
	return numsString
}

func getWhereAndParamNumsString(tagOperation int, structName string, columns []common_info.IColumn) (wheres []string, params []string) {
	for _, column := range columns {
		tagInfos := GetTagInfoMap().GetTag(structName, column.Name())
		if tagInfos == nil {
			continue
		}
		tagInfo := tagInfos[tagOperation]
		if tagInfo.Where != "" {
			wheres = append(wheres, tagInfo.Where)
			params = append(params, "cond."+column.Name())
		}
	}
	return
}

type Struct2DAO struct {
	config     *Config
	objectInfo *common_info.ObjectInfo
	err        error
}

func NewStruct2DAO(config *Config) *Struct2DAO {
	return &Struct2DAO{
		config:     config,
		objectInfo: common_info.GetParserInfo().Objects,
		err:        nil,
	}
}

func (s *Struct2DAO) Config(config Config) *Struct2DAO {
	s.config = &config
	return s
}

func (s *Struct2DAO) Error() error {
	return s.err
}

func (s *Struct2DAO) AutoMigrate(t ...interface{}) *Struct2DAO {
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

func (s *Struct2DAO) Run() error {
	if s.err != nil {
		return s.err
	}
	initContent := parseDBInit(s.config.dsn, s.config.packageName, s.config.ormPackageName, s.config.sqlPackageName)
	//只采用文件夹不支持单文件存储,文件夹不存在时自动创建
	if !util.IsDir(s.config.savePath) {
		s.config.savePath = filepath.Dir(s.config.savePath)
	}
	err := util.EnsureFileDirExist(s.config.savePath)
	if err != nil {
		s.err = err
		return s.err
	}
	err = util.SaveFile(filepath.Join(s.config.savePath, "init.go"), initContent)
	if err != nil {
		s.err = err
		return s.err
	}

	//根据CRUD标记确定是否生成对应内容
	objNums := len(*s.objectInfo)

	imports := make([]string, objNums)
	contents := make([]string, objNums)
	currentObjectsName := make([]string, 0)

	idx := 0
	for key, value := range *s.objectInfo {
		var content strings.Builder
		var importx string //import以 ‘,’ 分割多个
		content.WriteString(parseDAONewStruct(key))
		if s.config.enableCreate {
			content.WriteString(parseCreate(key, getOmitNumsString(KCreate, key, value)))
		}
		if s.config.enableQuery {
			content.WriteString(parseQuery(key, getOmitNumsString(KQuery, key, value)))
		}
		if s.config.enableUpdate {
			content.WriteString(parseUpdate(key, getOmitNumsString(KUpdate, key, value)))
		}
		//delete 允许where语句，故需要额外解析得到where字符串
		if s.config.enableDelete {
			wheres, params := getWhereAndParamNumsString(KDelete, key, value)
			content.WriteString(parseDelete(key, s.config.enableOr, wheres, params))
		}
		contents[idx] = content.String()
		importx += common_info.GetParserInfo().Dependencies.GetPackagePath(key)
		if strings.Contains(contents[idx], "time.Time") {
			importx += "time,"
		}
		imports[idx] = importx
		currentObjectsName = append(currentObjectsName, key)
		idx++
	}
	common_info.GetParserInfo().CurrentObjectsNameByOrder = currentObjectsName
	err = util.SaveAction(s.config.savePath, s.config.packageName, imports, contents)
	if err != nil {
		s.err = err
		return s.err
	}
	return s.err
}
