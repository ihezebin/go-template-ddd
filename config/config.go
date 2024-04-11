package config

import (
	"os"

	"github.com/ihezebin/oneness/config"
	"github.com/ihezebin/oneness/logger"
	"github.com/pkg/errors"
)

type Config struct {
	ServiceName string        `json:"service_name" mapstructure:"service_name"`
	Port        uint          `json:"port" mapstructure:"port"`
	MongoDsn    string        `json:"mongo_dsn" mapstructure:"mongo_dsn"`
	MysqlDsn    string        `json:"mysql_dsn" mapstructure:"mysql_dsn"`
	Logger      *LoggerConfig `json:"logger" mapstructure:"logger"`
	Redis       *RedisConfig  `json:"redis" mapstructure:"redis"`
	Pwd         string        `json:"-" mapstructure:"-"`
}

type RedisConfig struct {
	Addr     string `json:"addr" mapstructure:"addr"`
	Password string `json:"password" mapstructure:"password"`
}

type LoggerConfig struct {
	Level    logger.Level `json:"level" mapstructure:"level"`
	Filename string       `json:"filename" mapstructure:"filename"`
}

var gConfig *Config

func GetConfig() *Config {
	if gConfig == nil {
		gConfig = &Config{}
	}
	return gConfig
}

func Load(path string) (*Config, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, errors.Wrap(err, "get pwd error")
	}

	if err = config.NewWithFilePath(path).Load(&gConfig); err != nil {
		return nil, errors.Wrap(err, "load config error")
	}

	gConfig.Pwd = pwd

	return gConfig, nil
}
