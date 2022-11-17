package dao_convertor

import "fmt"

type Config struct {
	dsn            string //数据库连接字符串
	savePath       string //如果指定了文件，则为单文件，如果只指定了目录，则生成多文件
	packageName    string // 生成struct的包名(默认为空的话, 则取名为: package model)
	ormPackageName string //本dao层用到的orm框架包名
	sqlPackageName string //本dao层用到的sql驱动包名
	//是否生成CRUD
	enableCreate bool
	enableQuery  bool
	enableUpdate bool
	enableDelete bool
	//where语句的连接方式（and/or）
	enableOr    bool
	enableDebug bool
}

func DefaultConfig(username, password, host string, port int, database string) *Config {
	return &Config{
		dsn:            fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", username, password, host, port, database),
		savePath:       "./dao/",
		packageName:    "dao",
		ormPackageName: "gorm.io/gorm",
		sqlPackageName: "gorm.io/driver/mysql",
		enableDebug:    false,
		enableOr:       false,
		enableCreate:   true,
		enableQuery:    true,
		enableUpdate:   true,
		enableDelete:   true,
	}
}

func (Output *Config) Dsn(t string) *Config {
	Output.dsn = t
	return Output
}

func (Output *Config) SavePath(t string) *Config {
	Output.savePath = t
	return Output
}

func (Output *Config) PackageName(t string) *Config {
	Output.packageName = t
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

func (Output *Config) EnableCreate(t bool) *Config {
	Output.enableCreate = t
	return Output
}

func (Output *Config) EnableQuery(t bool) *Config {
	Output.enableQuery = t
	return Output
}

func (Output *Config) EnableUpdate(t bool) *Config {
	Output.enableUpdate = t
	return Output
}

func (Output *Config) EnableDelete(t bool) *Config {
	Output.enableDelete = t
	return Output
}

func (Output *Config) EnableOr(t bool) *Config {
	Output.enableOr = t
	return Output
}

func (Output *Config) EnableDebug(t bool) *Config {
	Output.enableDebug = t
	return Output
}
