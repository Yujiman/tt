package repository

import (
	"fmt"

	"git-ffd.kz/pkg/golog"
	trmgorm "github.com/avito-tech/go-transaction-manager/gorm"
	"gorm.io/gorm"

	"{{cookiecutter.package_name}}/internal/config"
)

type Repositories struct {
	Order    Order
	Customer Customer
	OneC     OneC
}

func NewRepositories(
	db *gorm.DB,
	trx *trmgorm.CtxGetter,
	cfg *config.Config,
	logger golog.ContextLogger,
) (*Repositories, error) {
	orderRepo := NewOrderDB(db, trx)
	customerRepo := NewCustomerDB(db, trx)
	oneCRepo, err := NewOneCClient(
		cfg.Integration.ExampleIntegrationBaseUrl,
		cfg.Integration.ExampleIntegrationToken,
		logger,
	)
	if err != nil {
		return nil, fmt.Errorf("NewOneCClient: %w", err)
	}

	return &Repositories{
		Order:    orderRepo,
		Customer: customerRepo,
		OneC:     oneCRepo,
	}, nil
}
