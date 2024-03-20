package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git-ffd.kz/fmobile/events/goevents/event_storage"
	"git-ffd.kz/pkg/goauth"
	"git-ffd.kz/pkg/goauth/contrib/ginauth"
	"git-ffd.kz/pkg/golog"
	watermilllog "git-ffd.kz/pkg/golog/contrib/watermill"
	"git-ffd.kz/pkg/gosentry"
	"git-ffd.kz/pkg/gowatermill/producer"
	"git-ffd.kz/pkg/gowatermill/subscriber"
	trmgorm "github.com/avito-tech/go-transaction-manager/gorm"
	"github.com/avito-tech/go-transaction-manager/trm/manager"
	"{{cookiecutter.package_name}}/docs"
	"{{cookiecutter.package_name}}/internal/common/middleware"
	"{{cookiecutter.package_name}}/internal/config"
	eventDelivery "{{cookiecutter.package_name}}/internal/delivery/event"
	httpDelivery "{{cookiecutter.package_name}}/internal/delivery/http"
	"{{cookiecutter.package_name}}/internal/repository"
	"{{cookiecutter.package_name}}/internal/server"
	"{{cookiecutter.package_name}}/internal/service"
	"{{cookiecutter.package_name}}/internal/worker"
	"{{cookiecutter.package_name}}/pkg/db/migrate"
	"{{cookiecutter.package_name}}/pkg/db/mongodb"
	"{{cookiecutter.package_name}}/pkg/db/postgresql"
	"{{cookiecutter.package_name}}/pkg/openapi"
)

func Run(cfg *config.Config) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	loggerConfig := golog.Config{
		Mode:              golog.ProductionMode,
		Level:             golog.InfoLevel,
		AppName:           cfg.Service.AppName,
		DisableStacktrace: true,
	}

	if cfg.Service.Environment.IsLocal() {
		loggerConfig.Mode = golog.DevelopmentMode
		loggerConfig.Level = golog.DebugLevel
	}

	logger, err := golog.NewZapLogger(loggerConfig)
	if err != nil {
		panic(err)
	}

	if cfg.Service.Domain == "127.0.0.1" || cfg.Service.Domain == "localhost" {
		docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%s", cfg.Service.Port)
		docs.SwaggerInfo.Schemes = []string{"http", "https"}
	} else {
		docs.SwaggerInfo.Host = cfg.Service.Domain
		docs.SwaggerInfo.Schemes = []string{"https", "http"}
		err = openapi.NewOpenApiClient(
			cfg.Service.AppName,
			cfg.Service.OpenapiEndpoint,
			docs.SwaggerInfo.ReadDoc(),
		).Send(context.Background()).Error()
		if err != nil {
			logger.Warnw("send openapi", err)
		}
	}

	connection, db, err := postgresql.NewDB(cfg.Database.WriteDSNs, cfg.Database.ReadDSNs, cfg.Service.AppName, cfg.Service.Environment.IsLocal(), logger)
	if err != nil {
		logger.Fatalw(err.Error())
	}

	if err := migrate.Up(connection, cfg.Database.Schema, "migrations"); err != nil {
		logger.Fatalw(err.Error())
	}

	if err := gosentry.SentryInit(cfg.Sentry.DSN, cfg.Service.Environment.String()); err != nil {
		logger.Fatalw(err.Error())
	}

	mongoDb, err := mongodb.NewClient(cfg.Database.MongoDSN)
	if err != nil {
		logger.Fatalw(err.Error())
	}

	var (
		gowatermillProducer   producer.Producer
		gowatermillSubscriber subscriber.Subscriber
	)

	if cfg.Nats.DSN != "" {
		natsCon, err := nats.Connect(cfg.Nats.DSN, nats.DrainTimeout(time.Second*5))
		if err != nil {
			logger.Fatalw("error connect to NATS", "err", err.Error())
		}

		jetStream, err := natsCon.JetStream()
		if err != nil {
			logger.Fatalw("error connect to JetStream", "err", err.Error())
		}

		gowatermillProducer = producer.NewJetStreamProducer(
			jetStream,
			producer.SaveLostMessages(event_storage.NewMongoStorage(mongoDb, logger), logger),
			logger,
		)
		gowatermillSubscriber = subscriber.NewWatermillSubscriber(
			natsCon,
			jetStream,
			cfg.Service.AppName,
			watermilllog.NewWatermillAdapter(logger),
		)
	} else {
		logger.Infow("NATS DSN is not set! Subscriber and producer not using!")
		gowatermillProducer = producer.NewNoOpProducer()
		gowatermillSubscriber = subscriber.NewNoOpSubscriber()
	}

	repos, err := repository.NewRepositories(db, trmgorm.DefaultCtxGetter, cfg, logger)
	if err != nil {
		logger.Fatalw("repository.NewRepositories", "err", err)
	}

	trxManager := manager.Must(trmgorm.NewDefaultFactory(db))
	services := service.NewServices(service.Deps{
		Repos:              repos,
		TransactionManager: trxManager,
		Cgf:                cfg,
		Logger:             logger,
		Producer:           gowatermillProducer,
	})

	eventConsumer, err := eventDelivery.NewHandler(
		gowatermillSubscriber,
		services,
		gowatermillProducer,
		logger,
	)
	if err != nil {
		logger.Fatalw(err.Error())
	}

	go func() {
		if err := eventConsumer.Run(ctx); err != nil {
			logger.Fatalw(err.Error())
		}
	}()

	healthCheckFn := func() error {
		if err := connection.Ping(); err != nil {
			return fmt.Errorf("database is not responding: %w", err)
		}

		return nil
	}

	authClient := goauth.NewAuthorizationClient(cfg.Authorization.BaseURL, logger)
	authManager := goauth.NewHttpAuthManager(authClient)
	authMiddleware := ginauth.NewAuthMiddleware(authManager, middleware.GinRecoveryFn, logger)

	handlerDelivery := httpDelivery.NewHandlerDelivery(
		logger,
		services,
		authMiddleware,
		healthCheckFn,
		"{{cookiecutter.project_name | slugify | lower}}",
	)
	workerHandler := worker.NewHandlerWorker(services, cfg, logger)

	if err := workerHandler.Start(); err != nil {
		logger.Fatalw(err.Error())
	}

	// HTTP Server
	srv, err := server.NewServer(cfg, handlerDelivery)
	if err != nil {
		logger.Fatalw(err.Error())
	}

	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalw("🔥 Server stopped due error", "err", err.Error())
		} else {
			logger.Infow("✅ Server shutdown successfully")
		}
	}()

	logger.Infow(fmt.Sprintf("🚀 Starting server at http://0.0.0.0:%s", cfg.Service.Port))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	shutdownCtx, shutdownCtxCancel := context.WithTimeout(context.Background(), time.Second*30)
	defer shutdownCtxCancel()

	isShutdownErrors := false

	if err = srv.Shutdown(shutdownCtx); err != nil {
		logger.Errorw(err.Error())
		isShutdownErrors = true
	}

	cancel()

	if err = connection.Close(); err != nil {
		logger.Errorw(err.Error())
		isShutdownErrors = true
	}
	if err = mongoDb.Disconnect(shutdownCtx); err != nil {
		logger.Errorw(err.Error())
		isShutdownErrors = true
	}

	if isShutdownErrors {
		logger.Warnw("Server closed, but not all resources closed properly!")
	} else {
		logger.Infow("✅ Server shutdown successfully")
	}
}
