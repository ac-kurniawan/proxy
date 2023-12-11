package main

import (
	"context"

	"github.com/ac-kurniawan/proxy/adaptor/datasouce"
	"github.com/ac-kurniawan/proxy/adaptor/gorm"
	"github.com/ac-kurniawan/proxy/adaptor/repository"
	"github.com/ac-kurniawan/proxy/adaptor/util"
	"github.com/ac-kurniawan/proxy/core"
	"github.com/ac-kurniawan/proxy/library"
)

type ProxyApp struct {
	Env        string `mapstructure:"env"`
	AppName    string `mapstructure:"appName"`
	Secret     string `mapstructure:"secret"`
	HttpServer struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"httpServer"`
	Trace struct {
		Enable       bool   `mapstructure:"enable"`
		HostExporter string `mapstructure:"hostExporter"`
		ApiKey       string `mapstructure:"apiKey"`
	} `mapstructure:"trace"`
	SQLite struct {
		RunMigration bool   `mapstructure:"runMigration"`
		FilePath     string `mapstructure:"filePath"`
	} `mapstructure:"sqlite"`
	HttpClient struct {
		Host        string `mapstructure:"host"`
		Port        string `mapstructure:"port"`
		TokenPrefix string `mapstructure:"tokenPrefix"`
	} `mapstructure:"httpClient"`
}

func (t ProxyApp) Init() {
	trace := library.NewAppTrace(context.Background(), t.Trace.Enable, t.Trace.HostExporter, t.Trace.ApiKey, t.AppName, "", t.Env)
	log := library.NewAppLog(false)
	utilCore := util.NewUtil(util.Util{
		AppTrace: trace,
		AppLog:   log,
		Secret:   t.Secret,
	})
	sqlite := library.NewGormSqliteConnection(t.SQLite.FilePath)
	DB := gorm.NewGormDB(gorm.GormDB{
		Gorm:  sqlite,
		Trace: trace,
	}, t.SQLite.RunMigration)
	ds := datasouce.NewDatasource(datasouce.Datasource{
		Host:        t.HttpClient.Host,
		Port:        t.HttpClient.Port,
		TokenPrefix: t.HttpClient.TokenPrefix,
		Trace:       trace,
	})
	repository := repository.NewRepository(repository.Repository{
		IProxyDB:    DB,
		IDataSource: ds,
	})
	service := core.NewProxyService(core.ProxyService{
		Repository: repository,
		Util:       utilCore,
	})
}
