package cmd

import (
	"context"
	"github.com/pkg/errors"
	"github.com/whereabouts/sdk/cli"
	"github.com/whereabouts/sdk/logger"
	"github.com/whereabouts/web-template-ddd/component/cache"
	"github.com/whereabouts/web-template-ddd/component/email"
	"github.com/whereabouts/web-template-ddd/component/sms"
	"github.com/whereabouts/web-template-ddd/component/storage"
	"github.com/whereabouts/web-template-ddd/config"
	"github.com/whereabouts/web-template-ddd/domain/repository"
	"github.com/whereabouts/web-template-ddd/server"
	"time"
)

func Run() {
	app := cli.NewApp(
		cli.WithName("web-template-ddd"),
		cli.WithVersion("v1.0"),
		cli.WithUsageText("Rapid construction template of Web service based on DDD architecture"),
		cli.WithAuthor("whereabouts.icu"),
	)
	app = app.WithFlagString("config, c", "./config.json", "config file path (default: ./config/config.json)", false)
	app = app.WithAction(func(v cli.Value) error {
		var (
			err  error
			conf *config.Config
			ctx  = context.Background()
			path = v.String("c")
		)

		// load config
		if conf, err = config.Load(path); err != nil {
			logger.Fatalf("failed to load config path: %s", path)
		}

		// init components
		if err = initComponents(ctx, conf); err != nil {
			logger.Fatalf("failed to init components: %v", err)
		}

		if err = server.NewServer(conf.Port).Run(ctx); err != nil {
			logger.Fatalf("failed to init components: %v", err)
		}

		return nil
	})
	_ = app.Run()
}

func initComponents(ctx context.Context, conf *config.Config) error {
	// init logger
	logger.ResetStandardLoggerWithConfig(conf.Logger)
	// init mongo
	if err := storage.InitMongo(ctx, conf.Mongo); err != nil {
		return errors.Wrap(err, "failed to init mongo")
	}
	// init redis
	if err := cache.InitRedis(ctx, conf.Redis); err != nil {
		return errors.Wrap(err, "failed to init redis")
	}
	// init memory
	cache.InitMemoryCache(5*time.Minute, time.Minute)
	// init email
	if err := email.Init(conf.Email); err != nil {
		return errors.Wrap(err, "failed to init email")
	}
	// init sms
	if err := sms.Init(conf.Sms.Config); err != nil {
		return errors.Wrap(err, "failed to init sms")
	}

	// init repository
	repository.InitTestRepository()

	return nil
}
