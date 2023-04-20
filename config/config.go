package config

import (
	"github.com/ihezebin/sdk/config"
	"github.com/ihezebin/sdk/emailc"
	"github.com/ihezebin/sdk/logger"
	"github.com/ihezebin/sdk/model/mongoc"
	"github.com/ihezebin/sdk/model/redisc"
	smsc "github.com/ihezebin/sdk/smsc/tencent"
)

type Config struct {
	Name   string         `mapstructure:"name"`
	Port   int            `mapstructure:"port"`
	Logger *logger.Config `mapstructure:"logger"`
	Mongo  *mongoc.Config `mapstructure:"mongo"`
	Redis  *redisc.Config `mapstructure:"redis"`
	Email  *emailc.Config `mapstructure:"email"`
	Sms    *Sms           `mapstructure:"sms"`
}

type Sms struct {
	Config  smsc.Config  `mapstructure:"config"`
	Message smsc.Message `mapstructure:"message"`
}

var gConfig Config

func GetConfig() Config {
	return gConfig
}

func Load(path string) (*Config, error) {
	if err := config.NewWithFilePath(path).Load(&gConfig); err != nil {
		return nil, err
	}
	return &gConfig, nil
}
