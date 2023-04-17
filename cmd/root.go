package cmd

import (
	"context"
	"github.com/ihezebin/go-template-ddd/component/cache"
	"github.com/ihezebin/go-template-ddd/component/email"
	"github.com/ihezebin/go-template-ddd/component/sms"
	"github.com/ihezebin/go-template-ddd/component/storage"
	"github.com/ihezebin/go-template-ddd/config"
	"github.com/ihezebin/go-template-ddd/domain/repository"
	"github.com/ihezebin/go-template-ddd/server"
	"github.com/ihezebin/sdk/cli"
	"github.com/ihezebin/sdk/logger"
	"github.com/pkg/errors"
	"time"
)

func Run() {
	app := cli.NewApp(
		cli.WithName("web-template-ddd"),
		cli.WithVersion("v1.0"),
		cli.WithUsageText("Rapid construction template of Web service based on DDD architecture"),
		cli.WithAuthor("whereabouts.icu"),
	)
	app = app.WithFlagString("config, c", "./config.toml", "config file path (default: ./config/config.json)", false)
	app = app.WithAction(func(v cli.Value) error {
		var (
			err  error
			conf *config.Config
			ctx  = context.Background()
			path = v.String("c")
		)

		// load config
		if conf, err = config.Load(path); err != nil {
			logger.WithError(err).Fatalf("failed to load config path: %s", path)
		}

		// ddd components
		if err = initComponents(ctx, conf); err != nil {
			logger.WithError(err).Fatalf("failed to ddd components")
		}

		if err = server.NewServer(conf.Port).Run(ctx); err != nil {
			logger.WithError(err).Fatalf("failed to ddd components")
		}

		return nil
	})
	_ = app.Run()
}

func initComponents(ctx context.Context, conf *config.Config) error {
	// ddd logger
	logger.ResetStandardLoggerWithConfig(conf.Logger)
	// ddd mongo
	if err := storage.InitMongo(ctx, conf.Mongo); err != nil {
		return errors.Wrap(err, "failed to ddd mongo")
	}
	// ddd redis
	if err := cache.InitRedis(ctx, conf.Redis); err != nil {
		return errors.Wrap(err, "failed to ddd redis")
	}
	// ddd memory
	cache.InitMemoryCache(5*time.Minute, time.Minute)
	// ddd email
	if err := email.Init(conf.Email); err != nil {
		return errors.Wrap(err, "failed to ddd email")
	}
	// ddd sms
	if err := sms.Init(conf.Sms.Config); err != nil {
		return errors.Wrap(err, "failed to ddd sms")
	}

	// ddd repository
	repository.InitTestRepository()

	return nil
}
