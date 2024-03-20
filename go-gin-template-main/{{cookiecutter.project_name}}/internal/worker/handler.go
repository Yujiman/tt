package worker

import (
	"fmt"
	"time"

	"git-ffd.kz/pkg/golog"
	"github.com/getsentry/sentry-go"
	"github.com/go-co-op/gocron"

	"{{cookiecutter.package_name}}/internal/config"
	"{{cookiecutter.package_name}}/internal/service"
)

type Handler struct {
	services *service.Services
	Cfg      *config.Config
	logger   golog.ContextLogger
}

func NewHandlerWorker(
	services *service.Services,
	cfg *config.Config,
	logger golog.ContextLogger,
) *Handler {
	return &Handler{
		services: services,
		Cfg:      cfg,
		logger:   logger,
	}
}

// Start запускает воркеров в отдельной горутине
func (h *Handler) Start() error {
	scheduler := gocron.NewScheduler(time.Local)
	scheduler.SingletonModeAll() // максимум будет работать один экземпляр таски
	scheduler.WaitForScheduleAll()
	gocron.SetPanicHandler(recoverer(h.logger))

	if h.Cfg.Worker.RunEchoWorker {
		h.logger.Infow("Register loader worker")
		if _, err := scheduler.Every(5).Minute().StartImmediately().Do(h.echo); err != nil {
			return fmt.Errorf("loader worker: %w", err)
		}
	}

	scheduler.StartAsync()

	return nil
}

func recoverer(logger golog.ContextLogger) func(jobName string, recoverData interface{}) {
	return func(jobName string, recoverData interface{}) {
		if recoverData != nil {
			if sentry.CurrentHub() != nil {
				hub := sentry.CurrentHub().Clone()
				if hub != nil {
					hub.Recover(recoverData)
				}

			}

			logger.Errorw("panic in worker", "panic", recoverData)
		}
	}
}
