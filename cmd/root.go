package cmd

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/ihezebin/go-template-ddd/component/cache"
	"github.com/ihezebin/go-template-ddd/component/storage"
	"github.com/ihezebin/go-template-ddd/config"
	"github.com/ihezebin/go-template-ddd/domain/repository"
	"github.com/ihezebin/go-template-ddd/domain/service"
	"github.com/ihezebin/go-template-ddd/server"
	"github.com/ihezebin/oneness/logger"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var (
	configPath string
)

func Run(ctx context.Context) error {

	app := &cli.App{
		Name:    "go-template-ddd",
		Version: "v1.0.0",
		Usage:   "Rapid construction template of Web service based on DDD architecture",
		Authors: []*cli.Author{
			{Name: "hezebin", Email: "ihezebin@qq.com"},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{Destination: &configPath, Name: "config", Aliases: []string{"c"}, Value: "./config/config.toml", Usage: "config file path (default find file from pwd and exec dir"},
		},
		Before: func(c *cli.Context) error {
			if configPath == "" {
				return errors.New("config path is empty")
			}

			conf, err := config.Load(configPath)
			if err != nil {
				return errors.Wrapf(err, "load config error, path: %s", configPath)
			}

			if err = initComponents(ctx, conf); err != nil {
				return errors.Wrap(err, "init components error")
			}

			logger.Infof(ctx, "component init success, config: %+v", *conf)

			return nil
		},
		Action: func(c *cli.Context) error {
			if err := server.Run(ctx, config.GetConfig().Port); err != nil {
				logger.WithError(err).Fatalf(ctx, "server run error, port: %d", config.GetConfig().Port)
			}

			return nil
		},
	}

	return app.Run(os.Args)
}

func initComponents(ctx context.Context, conf *config.Config) error {
	// init logger
	if conf.Logger != nil {
		logger.ResetLoggerWithOptions(
			logger.WithServiceName(conf.ServiceName),
			logger.WithCallerHook(),
			logger.WithTimestampHook(),
			logger.WithLevel(conf.Logger.Level),
			logger.WithLocalFsHook(filepath.Join(conf.Pwd, conf.Logger.Filename)),
		)
	}

	// init storage
	if err := storage.InitMySQLStorageClient(ctx, conf.MysqlDsn); err != nil {
		return errors.Wrap(err, "init mysql storage client error")
	}
	if err := storage.InitMongoStorageClient(ctx, conf.MongoDsn); err != nil {
		return errors.Wrap(err, "init mongo storage client error")
	}

	// init cache
	cache.InitMemoryCache(time.Minute*5, time.Minute)
	if err := cache.InitRedisCache(ctx, conf.Redis.Addr, conf.Redis.Password); err != nil {
		return errors.Wrap(err, "init redis cache client error")
	}

	// init repository
	repository.Init()

	// init service
	service.Init()

	return nil
}
