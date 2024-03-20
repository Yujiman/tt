package service

import (
	"context"
	"errors"

	"git-ffd.kz/fmobile/ferr"
	"git-ffd.kz/pkg/goerr"
	"git-ffd.kz/pkg/golog"
	"github.com/avito-tech/go-transaction-manager/trm"
	"github.com/google/uuid"
	"github.com/samber/mo"

	"{{cookiecutter.package_name}}/internal/model"
	"{{cookiecutter.package_name}}/internal/repository"
	"{{cookiecutter.package_name}}/internal/schema"
)

type Order interface {
	GetById(ctx context.Context, req schema.OrderGetByIDRequest) (schema.OrderResponse, error)
	Create(ctx context.Context, order schema.OrderCreateRequest) (schema.OrderResponse, error)
}

type OrderService struct {
	orderRepo    repository.Order
	customerRepo repository.Customer
	logger       golog.ContextLogger
	transaction  trm.Manager
}

func NewOrderService(
	orderRepo repository.Order,
	customerRepo repository.Customer,
	logger golog.ContextLogger,
	transaction trm.Manager,
) *OrderService {
	return &OrderService{
		orderRepo:    orderRepo,
		customerRepo: customerRepo,
		logger:       logger,
		transaction:  transaction,
	}
}

func (s *OrderService) GetById(ctx context.Context, req schema.OrderGetByIDRequest) (schema.OrderResponse, error) {
	ctx, _ = s.logger.FromContext(ctx, "OrderService", "service", "GetById", "method")

	order, err := s.orderRepo.GetOne(ctx, repository.OrderGetOneParams{
		OrderID:      mo.Some(req.OrderID),
		OrderSession: mo.None[string](),
	})
	if err != nil {
		return schema.OrderResponse{}, err
	}

	return schema.NewOrderResponseFromModel(order), nil
}

func (s *OrderService) Create(ctx context.Context, order schema.OrderCreateRequest) (schema.OrderResponse, error) {
	ctx, logger := s.logger.FromContext(ctx, "OrderService", "service", "Create", "method")

	var (
		createdOrder model.Order
		err          error
	)

	// Здесь все выполняется внутри транзакции в базе данных. Если вернется ошибка - транзакция будет откачена
	err = s.transaction.Do(ctx, func(ctx context.Context) error {
		customer, err := s.customerRepo.GetOne(ctx, repository.CustomerGetOneParams{
			CustomerID:    mo.None[uint](),
			CustomerPhone: mo.Some(order.Customer.Phone),
		})
		if errors.Is(err, ferr.ErrCustomerNotFound) {
			logger.Infow("customer not found, create new one", "phone", order.Customer.Phone)

			customer, err = s.customerRepo.Create(ctx, model.Customer{
				CustomerID: 0,
				FullName:   order.Customer.FullName,
				Phone:      order.Customer.Phone,
				Email:      order.Customer.Email,
			})
		}
		if err != nil {
			return goerr.Wrap(err).WithCtx(ctx)
		}

		createdOrder, err = s.orderRepo.Create(ctx, model.Order{
			TimestampMixin: model.TimestampMixin{},
			DeleteMixin:    model.DeleteMixin{},
			OrderID:        0,
			OrderSession:   uuid.NewString(),
			ChannelID:      order.ChannelID,
			BuyerID:        customer.CustomerID,
			OrderComment:   order.OrderComment,
		})
		if err != nil {
			return goerr.Wrap(err).WithCtx(ctx)
		}

		return nil
	})
	if err != nil {
		return schema.OrderResponse{}, goerr.Wrap(err).WithCtx(ctx)
	}

	return schema.NewOrderResponseFromModel(createdOrder), nil
}
