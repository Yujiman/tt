package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"git-ffd.kz/fmobile/ferr"
	"git-ffd.kz/pkg/goauth/contrib/ginauth"
	"git-ffd.kz/pkg/golog"
	"git-ffd.kz/pkg/golog/contrib/ginlog"
	"git-ffd.kz/pkg/requestid"
	"git-ffd.kz/pkg/requestid/contrib/ginid"
	"github.com/Depado/ginprom"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/ulule/limiter/v3"
	ginLimiter "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"

	"{{cookiecutter.package_name}}/internal/common/middleware"
	"{{cookiecutter.package_name}}/pkg/middlewares"

	_ "{{cookiecutter.package_name}}/docs"
	"{{cookiecutter.package_name}}/internal/config"
	v1 "{{cookiecutter.package_name}}/internal/delivery/http/v1"
	"{{cookiecutter.package_name}}/internal/service"
)

type Handler struct {
	logger         golog.ContextLogger
	services       *service.Services
	baseUrl        string
	authMiddleware *ginauth.AuthMiddleware
	healthcheckFn  func() error
}

func NewHandlerDelivery(
	logger golog.ContextLogger,
	services *service.Services,
	authMiddleware *ginauth.AuthMiddleware,
	healthcheckFn func() error,
	baseUrl string,
) *Handler {
	return &Handler{
		logger:         logger,
		services:       services,
		baseUrl:        baseUrl,
		authMiddleware: authMiddleware,
		healthcheckFn:  healthcheckFn,
	}
}

func (h *Handler) Init(cfg *config.Config) (*gin.Engine, error) {
	if !cfg.Service.Environment.IsLocal() {
		gin.SetMode(gin.ReleaseMode)
	}

	app := gin.New()
	prom := ginprom.New(
		ginprom.Engine(app),
		ginprom.Path("/actuator/prometheus"),
		ginprom.Subsystem(""),
		ginprom.Namespace(""),
	)

	app.Use(
		middlewares.Cors(),
		ginid.NewMiddleware(requestid.WithHandler(func(ctx context.Context, requestID string) context.Context {
			ctx, _ = h.logger.FromContext(ctx, "request_id", requestID)
			return ctx
		})),
		ginlog.New(
			h.logger,
			ginlog.WithLogRequestBody(cfg.Service.Environment.IsLocal()),
			ginlog.WithLogResponseBody(cfg.Service.Environment.IsLocal()),
		),
		middlewares.Recovery(middleware.GinRecoveryFn),
		func(c *gin.Context) {
			ctx := sentry.SetHubOnContext(c.Request.Context(), sentry.CurrentHub().Clone())
			c.Request = c.Request.WithContext(ctx)
			c.Next()
		},
		sentrygin.New(sentrygin.Options{
			Repanic:         true,
			WaitForDelivery: true,
			Timeout:         time.Second,
		}),
		prom.Instrument(),
		h.authMiddleware.SetCurrentUser(),
	)

	if err := h.applyLimiter(cfg.Service.Limiter, app); err != nil {
		return nil, fmt.Errorf("can't apply limiter: %w", err)
	}

	app.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]string{"message": "pong"})
	})
	app.GET("/readiness", func(c *gin.Context) {
		if err := h.healthcheckFn(); err != nil {
			c.JSON(http.StatusServiceUnavailable, map[string]string{"message": err.Error()})
			c.Error(err)
		} else {
			c.JSON(http.StatusOK, map[string]string{"message": "ok"})
		}
	})
	app.GET("/liveness", func(c *gin.Context) {
		if err := h.healthcheckFn(); err != nil {
			c.JSON(http.StatusServiceUnavailable, map[string]string{"message": err.Error()})
			c.Error(err)
		} else {
			c.JSON(http.StatusOK, map[string]string{"message": "ok"})
		}
	})

	h.initAPI(app)

	return app, nil
}

func (h *Handler) applyLimiter(limiterSettings string, router *gin.Engine) error {
	if limiterSettings != "" {
		rate, err := limiter.NewRateFromFormatted(limiterSettings)
		if err != nil {
			return err
		}

		limit := limiter.New(
			memory.NewStore(),
			rate,
			limiter.WithTrustForwardHeader(true),
		)

		router.Use(
			ginLimiter.NewMiddleware(limit, ginLimiter.WithLimitReachedHandler(func(c *gin.Context) {
				c.Error(ferr.ErrLimitReached)
				c.AbortWithStatus(http.StatusTooManyRequests)
			})),
		)
	}

	return nil
}

func (h *Handler) initAPI(router *gin.Engine) {
	baseUrl := router.Group(h.baseUrl)

	router.GET(h.baseUrl+"-docs/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	handlerV1 := v1.NewHandler(h.services, h.authMiddleware)
	api := baseUrl.Group("/api")
	{
		handlerV1.Init(api)
	}
}
