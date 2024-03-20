package v1

import (
	"git-ffd.kz/pkg/goauth/contrib/ginauth"
	"github.com/gin-gonic/gin"

	"{{cookiecutter.package_name}}/internal/service"
)

type Handler struct {
	services *service.Services
	auth     *ginauth.AuthMiddleware
}

func NewHandler(services *service.Services, auth *ginauth.AuthMiddleware) *Handler {
	return &Handler{
		services: services,
		auth:     auth,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initOrderHandler(v1)
	}
}
