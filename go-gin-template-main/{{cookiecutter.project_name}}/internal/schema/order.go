package schema

import (
	"time"

	"git-ffd.kz/fmobile/ferr"
	"git-ffd.kz/pkg/goerr"

	"{{cookiecutter.package_name}}/internal/model"
)

type OrderGetByIDRequest struct {
	OrderID uint `json:"order_id" form:"order_id"`
}

func (r OrderGetByIDRequest) Validate() error {
	if r.OrderID == 0 {
		return ferr.ErrOrderIdRequired.WithStack()
	}

	return nil
}

type OrderCustomer struct {
	FullName string `json:"full_name"`
	Phone    string `json:"phone" validate:"required"`
	Email    string `json:"email"`
}

func (c OrderCustomer) Validate() error {
	if c.Phone == "" {
		return ferr.ErrPhoneRequired.WithStack()
	}

	if len(c.Phone) != 11 {
		return ferr.ErrPhoneWrongFormat.WithStack()
	}

	return nil
}

type OrderCreateRequest struct {
	ChannelID    uint          `json:"channel_id"`
	Customer     OrderCustomer `json:"customer"`
	OrderComment string        `json:"order_comment,omitempty"`
}

func (o OrderCreateRequest) Validate() error {
	if o.ChannelID == 0 {
		return ferr.ErrOrderChannelIdRequired.WithStack()
	}

	if err := o.Customer.Validate(); err != nil {
		return goerr.Wrap(err)
	}

	return nil
}

type OrderResponse struct {
	CreatedAt    time.Time `json:"created_at"`
	OrderID      uint      `json:"order_id"`
	OrderSession string    `json:"order_session"`
	ChannelID    uint      `json:"channel_id"`
	BuyerID      uint      `json:"buyer_id"`
	OrderComment string    `json:"order_comment"`
}

func NewOrderResponseFromModel(order model.Order) OrderResponse {
	return OrderResponse{
		CreatedAt:    order.CreatedAt,
		OrderID:      order.OrderID,
		OrderSession: order.OrderSession,
		ChannelID:    order.ChannelID,
		BuyerID:      order.BuyerID,
		OrderComment: order.OrderComment,
	}
}
