package config

import (
	"github.com/whereabouts/sdk/config"
	"github.com/whereabouts/sdk/db/mongoc"
	"github.com/whereabouts/sdk/db/redisc"
	"github.com/whereabouts/sdk/emailc"
	"github.com/whereabouts/sdk/logger"
	smsc "github.com/whereabouts/sdk/smsc/tencent"
)

type Config struct {
	Name   string        `mapstructure:"name"`
	Port   int           `mapstructure:"port"`
	Logger logger.Config `mapstructure:"logger"`
	Mongo  mongoc.Config `mapstructure:"mongo"`
	Redis  redisc.Config `mapstructure:"redis"`
	Email  emailc.Config `mapstructure:"email"`
	Sms    Sms           `mapstructure:"sms"`
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
	if err := config.LoadWithFilePath(path, &gConfig); err != nil {
		return nil, err
	}
	return &gConfig, nil
}
