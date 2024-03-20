package v1

import (
	"git-ffd.kz/pkg/goerr"
	"github.com/gin-gonic/gin"

	"{{cookiecutter.package_name}}/internal/common/middleware"
	"{{cookiecutter.package_name}}/internal/schema"
)

func (h *Handler) initOrderHandler(v1 *gin.RouterGroup) {
	v1.GET("orders/item", middleware.GinErrorHandle(h.GetOrderByID))
	v1.POST("orders", middleware.GinErrorHandle(h.CreateOrder))
}

// GetOrderByID
// WhoAmi godoc
// @Summary Получение заказа по ID
// @Param data query schema.OrderGetByIDRequest true "data"
// @Produce json
// @Success 200 {object} schema.Response[schema.OrderResponse]
// @tags orders
// @Router /api/v1/orders/item [get]
func (h *Handler) GetOrderByID(c *gin.Context) (err error) {
	ctx := c.Request.Context()

	var data schema.OrderGetByIDRequest
	if err := c.Bind(&data); err != nil {
		return goerr.Wrap(err).WithCtx(ctx)
	}

	if err := data.Validate(); err != nil {
		return goerr.Wrap(err).WithCtx(ctx)
	}

	order, err := h.services.Order.GetById(ctx, data)
	if err != nil {
		return err
	}

	return schema.Respond(order, c)
}

// CreateOrder
// WhoAmi godoc
// @Summary Создание заказа
// @Accept json
// @Produce json
// @Param data body schema.OrderCreateRequest true "OrderCreateRequest Data"
// @Success 200 {object} schema.Response[schema.OrderResponse]
// @Failure 400 {object} schema.Response[schema.Empty]
// @tags orders
// @Router /api/v1/orders [post]
func (h *Handler) CreateOrder(c *gin.Context) error {
	var data schema.OrderCreateRequest

	if err := c.BindJSON(&data); err != nil {
		return err
	}

	if err := data.Validate(); err != nil {
		return goerr.Wrap(err)
	}

	createdOrder, err := h.services.Order.Create(c.Request.Context(), data)
	if err != nil {
		return err
	}

	return schema.Respond(createdOrder, c)
}
