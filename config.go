package quickstart

import "github.com/ACking-you/quickstart_project/common_info"

type Config struct {
	//数据库连接
	username            string
	password            string
	host                string
	database            string
	port                int
	ormPackageName      string                      //gorm包名
	sqlPackageName      string                      //sql包名
	ginPackageName      string                      //gin包名
	basePath            string                      //用于生成文件的基本路径，需要为文件夹
	table               string                      //如果指定了表，则只针对该表生成
	realTableNameMethod string                      //若该字段被赋值，则以该字段名生成获取真实表明的方法
	tagKey              string                      // tag字段的key值,默认是gorm
	enableGoTidy        bool                        //是否自动go get
	enableJsonTag       bool                        // 是否添加json的tag, 默认不添加
	enableDebug         bool                        //是否打印结构体的解析结果
	tableFilterHook     func(tableName string) bool //用于过滤不需要生成的表结构
}

func DefaultConfig(projectName, username, password, host string, port int, database string) *Config {
	common_info.ProjectPackageName = projectName
	return &Config{
		username:            username,
		password:            password,
		host:                host,
		port:                port,
		ormPackageName:      "gorm.io/gorm",
		sqlPackageName:      "gorm.io/driver/mysql",
		ginPackageName:      "github.com/gin-gonic/gin",
		database:            database,
		basePath:            "./", //默认为当前项目根目录
		realTableNameMethod: "TableName",
		tagKey:              "gorm",
		enableGoTidy:        false,
		enableJsonTag:       false,
		enableDebug:         false,
		tableFilterHook:     nil,
	}
}

func (Output *Config) Username(t string) *Config {
	Output.username = t
	return Output
}

func (Output *Config) Password(t string) *Config {
	Output.password = t
	return Output
}

func (Output *Config) Host(t string) *Config {
	Output.host = t
	return Output
}

func (Output *Config) Database(t string) *Config {
	Output.database = t
	return Output
}

func (Output *Config) Port(t int) *Config {
	Output.port = t
	return Output
}

func (Output *Config) OrmPackageName(t string) *Config {
	Output.ormPackageName = t
	return Output
}

func (Output *Config) SqlPackageName(t string) *Config {
	Output.sqlPackageName = t
	return Output
}

func (Output *Config) GinPackageName(t string) *Config {
	Output.ginPackageName = t
	return Output
}

func (Output *Config) BasePath(t string) *Config {
	Output.basePath = t
	return Output
}

func (Output *Config) Table(t string) *Config {
	Output.table = t
	return Output
}

func (Output *Config) RealTableNameMethod(t string) *Config {
	Output.realTableNameMethod = t
	return Output
}

func (Output *Config) TagKey(t string) *Config {
	Output.tagKey = t
	return Output
}

func (Output *Config) EnableGoTidy(t bool) *Config {
	Output.enableGoTidy = t
	return Output
}

func (Output *Config) EnableJsonTag(t bool) *Config {
	Output.enableJsonTag = t
	return Output
}

func (Output *Config) EnableDebug(t bool) *Config {
	Output.enableDebug = t
	return Output
}

func (Output *Config) TableFilterHook(t func(string) bool) *Config {
	Output.tableFilterHook = t
	return Output
}
