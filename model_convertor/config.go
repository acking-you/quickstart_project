package model_convertor

import "fmt"

type Config struct {
	dsn                 string                      //数据库连接字符串
	savePath            string                      //如果指定了文件，则为单文件，如果只指定了目录，则生成多文件
	table               string                      //如果指定了表，则只针对该表生成
	realTableNameMethod string                      //若该字段被赋值，则以该字段名生成获取真实表明的方法
	packageName         string                      // 生成struct的包名(默认为空的话, 则取名为: package models)
	tagKey              string                      // tag字段的key值,默认是gorm
	enableJsonTag       bool                        // 是否添加json的tag, 默认不添加
	enableDebug         bool                        //是否打印结构体的解析结果
	tableFilterHook     func(tableName string) bool //用于过滤不需要生成的表结构
}

func DefaultConfig(username, password, host string, port int, database string) *Config {
	return &Config{
		dsn:                 fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", username, password, host, port, database),
		savePath:            "./models/",
		realTableNameMethod: "TableName",
		packageName:         "models",
		tagKey:              "gorm",
		enableJsonTag:       false,
		enableDebug:         false,
		tableFilterHook:     nil,
	}
}

func (s *Config) Dsn(t string) *Config {
	s.dsn = t
	return s
}

func (s *Config) SavePath(t string) *Config {
	s.savePath = t
	return s
}

func (s *Config) Table(t string) *Config {
	s.table = t
	return s
}

func (s *Config) RealTableNameMethod(t string) *Config {
	s.realTableNameMethod = t
	return s
}

func (s *Config) PackageName(t string) *Config {
	s.packageName = t
	return s
}

func (s *Config) TagKey(t string) *Config {
	s.tagKey = t
	return s
}

func (s *Config) EnableJsonTag(t bool) *Config {
	s.enableJsonTag = t
	return s
}

func (s *Config) EnableDebug(t bool) *Config {
	s.enableDebug = t
	return s
}

func (s *Config) TableFilterHook(t func(string) bool) *Config {
	s.tableFilterHook = t
	return s
}
