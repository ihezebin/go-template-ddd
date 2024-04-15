package cmd

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/ihezebin/go-template-ddd/component/cache"
	"github.com/ihezebin/go-template-ddd/component/email"
	"github.com/ihezebin/go-template-ddd/component/pubsub"
	"github.com/ihezebin/go-template-ddd/component/storage"
	"github.com/ihezebin/go-template-ddd/config"
	"github.com/ihezebin/go-template-ddd/cron"
	"github.com/ihezebin/go-template-ddd/domain/repository"
	"github.com/ihezebin/go-template-ddd/domain/service"
	"github.com/ihezebin/go-template-ddd/server"
	"github.com/ihezebin/go-template-ddd/worker"
	"github.com/ihezebin/go-template-ddd/worker/example"
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
			&cli.StringFlag{
				Destination: &configPath,
				Name:        "config", Aliases: []string{"c"},
				Value: "./config/config.toml",
				Usage: "config file path (default find file from pwd and exec dir",
			},
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

			logger.Debugf(ctx, "component init success, config: %+v", *conf)

			return nil
		},
		Action: func(c *cli.Context) error {
			worker.Register(example.NewExampleWorker())
			worker.Run(ctx)
			defer worker.Wait(ctx)

			if err := cron.Run(ctx); err != nil {
				logger.WithError(err).Fatalf(ctx, "cron run error")
			}

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
	if conf.MysqlDsn != "" {
		if err := storage.InitMySQLStorageClient(ctx, conf.MysqlDsn); err != nil {
			return errors.Wrap(err, "init mysql storage client error")
		}
	}
	if conf.MongoDsn != "" {
		if err := storage.InitMongoStorageClient(ctx, conf.MongoDsn); err != nil {
			return errors.Wrap(err, "init mongo storage client error")
		}
	}

	// init cache
	cache.InitMemoryCache(time.Minute*5, time.Minute)
	if conf.Redis != nil {
		if err := cache.InitRedisCache(ctx, conf.Redis.Addr, conf.Redis.Password); err != nil {
			return errors.Wrap(err, "init redis cache client error")
		}
	}

	// init repository
	repository.Init()

	// init pubsub
	if conf.Pulsar != nil {
		if err := pubsub.InitPulsarClient(conf.Pulsar.Url); err != nil {
			return errors.Wrap(err, "init pulsar client error")
		}
	}
	if conf.Kafka != nil {
		if err := pubsub.InitKafkaConn(ctx, conf.Kafka.Address, conf.Kafka.Topic, conf.Kafka.Partition); err != nil {
			return errors.Wrap(err, "init kafka client error")
		}
	}

	// init service
	service.Init()

	// init email
	if conf.Email != nil {
		if err := email.Init(*conf.Email); err != nil {
			return errors.Wrap(err, "init email client error")
		}
	}

	return nil
}
