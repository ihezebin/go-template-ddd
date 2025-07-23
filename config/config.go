package config

import (
	"encoding/json"
	"os"

	"github.com/ihezebin/olympus/config"
	"github.com/ihezebin/olympus/email"
	"github.com/ihezebin/olympus/logger"
	"github.com/ihezebin/olympus/sms/tencent"
	"github.com/pkg/errors"
)

type Config struct {
	ServiceName      string          `json:"service_name" mapstructure:"service_name"`
	Port             uint            `json:"port" mapstructure:"port"`
	MongoDsn         string          `json:"mongo_dsn" mapstructure:"mongo_dsn"`
	MysqlDsn         string          `json:"mysql_dsn" mapstructure:"mysql_dsn"`
	ClickhouseDsn    string          `json:"clickhouse_dsn" mapstructure:"clickhouse_dsn"`
	OSSDsn           string          `json:"oss_dsn" mapstructure:"oss_dsn"`
	ElasticsearchUrl string          `json:"elasticsearch_url" mapstructure:"elasticsearch_url"`
	Pwd              string          `json:"-" mapstructure:"-"`
	Logger           *LoggerConfig   `json:"logger" mapstructure:"logger"`
	Redis            *RedisConfig    `json:"redis" mapstructure:"redis"`
	Email            *email.Config   `json:"email" mapstructure:"email"`
	Sms              *tencent.Config `json:"sms" mapstructure:"sms"`
	Pulsar           *PulsarConfig   `json:"pulsar" mapstructure:"pulsar"`
	Kafka            *KafkaConfig    `json:"kafka" mapstructure:"kafka"`
	Remote           *RemoteConfig   `json:"remote" mapstructure:"remote"`
	Wx               *WxConfig       `json:"wx" mapstructure:"wx"`
}

type PulsarConfig struct {
	Url          string `json:"url" mapstructure:"url"`
	Topic        string `json:"topic" mapstructure:"topic"`
	Subscription string `json:"subscription" mapstructure:"subscription"`
}

type KafkaConfig struct {
	Address   string `json:"address" mapstructure:"address"`
	Topic     string `json:"topic" mapstructure:"topic"`
	Partition int    `json:"partition" mapstructure:"partition"`
}

type RedisConfig struct {
	Addrs    []string `json:"addrs" mapstructure:"addrs"`
	Password string   `json:"password" mapstructure:"password"`
}

type LoggerConfig struct {
	Level    logger.Level `json:"level" mapstructure:"level"`
	Filename string       `json:"filename" mapstructure:"filename"`
}

type RemoteConfig struct {
	UserCenterHost string `json:"user_center_host" mapstructure:"user_center_host"`
}

type WxConfig struct {
	Host   string `json:"host" mapstructure:"host"`
	AppId  string `json:"app_id" mapstructure:"app_id"`
	Secret string `json:"secret" mapstructure:"secret"`
}

var gConfig *Config = &Config{}

func (c *Config) String() string {
	data, _ := json.Marshal(c)
	return string(data)
}

func GetConfig() *Config {
	return gConfig
}

func Load(path string) (*Config, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, errors.Wrap(err, "get pwd error")
	}

	if err = config.NewWithFilePath(path).Load(gConfig); err != nil {
		return nil, errors.Wrap(err, "load config error")
	}

	gConfig.Pwd = pwd

	return gConfig, nil
}
