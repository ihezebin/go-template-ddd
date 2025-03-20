package cmd

import (
	"context"
	"os"
	"path/filepath"
	"time"

	_ "github.com/ihezebin/olympus"
	"github.com/ihezebin/olympus/logger"
	"github.com/ihezebin/olympus/runner"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/ihezebin/go-template-ddd/component/cache"
	"github.com/ihezebin/go-template-ddd/component/emailc"
	"github.com/ihezebin/go-template-ddd/component/oss"
	"github.com/ihezebin/go-template-ddd/component/pubsub"
	"github.com/ihezebin/go-template-ddd/component/storage"
	"github.com/ihezebin/go-template-ddd/config"
	"github.com/ihezebin/go-template-ddd/cron"
	"github.com/ihezebin/go-template-ddd/domain/repository"
	"github.com/ihezebin/go-template-ddd/domain/service"
	"github.com/ihezebin/go-template-ddd/server"
	"github.com/ihezebin/go-template-ddd/worker"
	"github.com/ihezebin/go-template-ddd/worker/example"
)

var (
	configPath string
)

func Run(ctx context.Context) error {

	app := &cli.App{
		Name:    "go-template-ddd",
		Version: "v1.0.1",
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
			logger.Debugf(ctx, "load config: %+v", conf.String())

			if err = initComponents(ctx, conf); err != nil {
				return errors.Wrap(err, "init components error")
			}

			return nil
		},
		Action: func(c *cli.Context) error {
			httpServer, err := server.NewServer(ctx, config.GetConfig())
			if err != nil {
				return errors.Wrap(err, "new http server err")
			}

			tasks := make([]runner.Task, 0)
			tasks = append(tasks, worker.NewWorKeeper(example.NewExampleWorker()))
			tasks = append(tasks, cron.NewCron())
			tasks = append(tasks, httpServer)

			runner.NewRunner(tasks...).Run(ctx)

			return nil
		},
	}

	return app.Run(os.Args)
}

func initComponents(ctx context.Context, conf *config.Config) error {
	// init logger
	if conf.Logger != nil {
		logger.ResetLoggerWithOptions(
			logger.WithLoggerType(logger.LoggerTypeZap),
			logger.WithServiceName(conf.ServiceName),
			logger.WithCaller(),
			logger.WithTimestamp(),
			logger.WithLevel(conf.Logger.Level),
			//logger.WithLocalFsHook(filepath.Join(conf.Pwd, conf.Logger.Filename)),
			logger.WithRotate(logger.RotateConfig{
				Path:               filepath.Join(conf.Pwd, conf.Logger.Filename),
				MaxSizeKB:          1024 * 500, // 500 MB
				MaxAge:             time.Hour * 24 * 7,
				MaxRetainFileCount: 3,
				Compress:           true,
			}),
		)
	}

	// init storage
	if conf.MysqlDsn != "" {
		if err := storage.InitMySQLClient(ctx, conf.MysqlDsn); err != nil {
			return errors.Wrap(err, "init mysql storage client error")
		}
	}
	if conf.MongoDsn != "" {
		if err := storage.InitMongoClient(ctx, conf.MongoDsn); err != nil {
			return errors.Wrap(err, "init mongo storage client error")
		}
	}

	// init oss
	if conf.OSSDsn != "" {
		if err := oss.Init(conf.OSSDsn); err != nil {
			return errors.Wrap(err, "init oss client error")
		}
	}

	// init cache
	cache.InitMemoryCache(time.Minute*5, time.Minute)
	if conf.Redis != nil {
		if err := cache.InitRedisCache(ctx, conf.Redis.Addrs, conf.Redis.Password); err != nil {
			return errors.Wrap(err, "init redis cache client error")
		}
	}

	// init repository
	if conf.MysqlDsn != "" || conf.MongoDsn != "" {
		repository.Init()
	}

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

	// init email
	if conf.Email != nil {
		if err := emailc.Init(*conf.Email); err != nil {
			return errors.Wrap(err, "init email client error")
		}
	}

	// init clickhouse
	if conf.ClickhouseDsn != "" {
		if err := storage.InitClickhouseDatabase(ctx, conf.ClickhouseDsn); err != nil {
			return errors.Wrap(err, "init clickhouse storage database error")
		}
	}

	// init elasticsearch
	if conf.ElasticsearchUrl != "" {
		if err := storage.InitElasticsearchClient(ctx, conf.ElasticsearchUrl); err != nil {
			return errors.Wrap(err, "init elasticsearch client error")
		}
	}

	// init service
	service.Init()

	return nil
}
