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
		cli.WithName("go-template-ddd"),
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
	// init logger
	if conf.Logger != nil {
		logger.ResetStandardLoggerWithConfig(*conf.Logger)
	}
	// init memory
	cache.InitMemoryCache(5*time.Minute, time.Minute)
	// init redis
	if conf.Redis != nil {
		if err := cache.InitRedis(ctx, *conf.Redis); err != nil {
			return errors.Wrap(err, "failed to init redis")
		}
	}
	// init mongo
	if conf.Mongo != nil {
		if err := storage.InitMongo(ctx, *conf.Mongo); err != nil {
			return errors.Wrap(err, "failed to init mongo")
		}
	}

	// init mail
	if conf.Email != nil {
		if err := email.Init(*conf.Email); err != nil {
			return errors.Wrap(err, "failed to init mail")
		}
	}

	// init sms
	if conf.Sms != nil {
		if err := sms.Init(conf.Sms.Config, conf.Sms.Message); err != nil {
			return errors.Wrap(err, "failed to init sms")
		}
	}

	// init repository
	if conf.Mongo != nil && conf.Redis != nil {
		repository.Init("hezebin")
	}

	return nil
}
