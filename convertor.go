package quickstart

import (
	"errors"
	"fmt"
	"github.com/ACking-you/quickstart_project/controller_convertor"
	"github.com/ACking-you/quickstart_project/dao_convertor"
	"github.com/ACking-you/quickstart_project/model_convertor"
	"github.com/ACking-you/quickstart_project/service_convertor"
	"github.com/ACking-you/quickstart_project/util"
	"path/filepath"
)

func Quickstart(config *Config) error {
	if config.enableGoTidy {
		util.EnableGoTidy = true
	}
	modelConfig := model_convertor.DefaultConfig(config.username, config.password, config.host, config.port, config.database)
	modelConfig.EnableDebug(config.enableDebug).
		EnableJsonTag(config.enableJsonTag).
		RealTableNameMethod(config.realTableNameMethod).
		Table(config.table).
		TableFilterHook(config.tableFilterHook).
		TagKey(config.tagKey).
		SavePath(filepath.Join(config.basePath, "models"))
	err := model_convertor.NewTable2Struct(modelConfig).Run()
	if err != nil {
		return errors.New(fmt.Sprintf("err in model: %s", err.Error()))
	}

	daoConfig := dao_convertor.DefaultConfig(config.username, config.password, config.host, config.port, config.database)
	daoConfig.SqlPackageName(config.sqlPackageName).
		OrmPackageName(config.ormPackageName).
		EnableDebug(config.enableDebug).
		SavePath(filepath.Join(config.basePath, "dao"))
	err = dao_convertor.NewStruct2DAO(daoConfig).Run()
	if err != nil {
		return errors.New(fmt.Sprintf("err in dao: %s", err.Error()))
	}

	serviceConfig := service_convertor.DefaultConfig().
		EnableDebug(config.enableDebug).
		SavePath(filepath.Join(config.basePath, "service"))
	err = service_convertor.NewStruct2Service(serviceConfig).Run()
	if err != nil {
		return errors.New(fmt.Sprintf("err in service: %s", err.Error()))
	}

	controllerConfig := controller_convertor.DefaultCConfig().
		GinPackageName(config.ginPackageName).
		EnableDebug(config.enableDebug).
		SavePath(filepath.Join(config.basePath, "controller"))
	err = controller_convertor.NewStruct2Controller(controllerConfig).Run()
	if err != nil {
		return errors.New(fmt.Sprintf("err in controller: %s", err.Error()))
	}
	return nil
}
