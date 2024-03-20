package service

import (
	"git-ffd.kz/pkg/golog"
	"git-ffd.kz/pkg/gowatermill/producer"
	"github.com/avito-tech/go-transaction-manager/trm"

	"{{cookiecutter.package_name}}/internal/config"
	"{{cookiecutter.package_name}}/internal/repository"
)

type Services struct {
	Order Order
}

type Deps struct {
	Repos              *repository.Repositories
	Cgf                *config.Config
	TransactionManager trm.Manager
	Logger             golog.ContextLogger
	Producer           producer.Producer
}

func NewServices(deps Deps) *Services {
	orderService := NewOrderService(deps.Repos.Order, deps.Repos.Customer, deps.Logger, deps.TransactionManager)

	return &Services{
		Order: orderService,
	}
}
