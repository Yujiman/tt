package repository

import (
	"context"
	"errors"

	"git-ffd.kz/fmobile/ferr"
	trmgorm "github.com/avito-tech/go-transaction-manager/gorm"
	"github.com/samber/mo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"{{cookiecutter.package_name}}/internal/model"
)

type Customer interface {
	GetOne(ctx context.Context, params CustomerGetOneParams) (Customer model.Customer, err error)
	Create(ctx context.Context, Customer model.Customer) (model.Customer, error)
	Update(ctx context.Context, customerID uint, params CustomerUpdateParams) (model.Customer, error)
}

type CustomerDB struct {
	db  *gorm.DB
	trx *trmgorm.CtxGetter
}

func NewCustomerDB(db *gorm.DB, trx *trmgorm.CtxGetter) *CustomerDB {
	return &CustomerDB{
		db:  db,
		trx: trx,
	}
}

func (r *CustomerDB) GetOne(ctx context.Context, params CustomerGetOneParams) (customer model.Customer, err error) {
	db := r.trx.DefaultTrOrDB(ctx, r.db)

	q := db.Model(&model.Customer{})

	if CustomerID, ok := params.CustomerID.Get(); ok {
		q = q.Where(`customer_id = ?`, CustomerID)
	}
	if CustomerSession, ok := params.CustomerPhone.Get(); ok {
		q = q.Where(`phone = ?`, CustomerSession)
	}

	err = q.First(&customer).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return customer, ferr.ErrCustomerNotFound.WithErr(err).WithCtx(ctx)
		}

		return customer, ferr.ErrDbUnexpected.WithErr(err).WithCtx(ctx)
	}

	return customer, nil
}

func (r *CustomerDB) Create(ctx context.Context, customer model.Customer) (model.Customer, error) {
	db := r.trx.DefaultTrOrDB(ctx, r.db)

	err := db.Create(&customer).Error
	if err != nil {
		return model.Customer{}, ferr.ErrDbUnexpected.WithErr(err).WithCtx(ctx)
	}

	return customer, nil
}

func (r *CustomerDB) Update(ctx context.Context, customerID uint, params CustomerUpdateParams) (model.Customer, error) {
	db := r.trx.DefaultTrOrDB(ctx, r.db)

	updates := make(map[string]interface{}, 2)

	if fullName, ok := params.FullName.Get(); ok {
		updates["full_name"] = fullName
	}
	if email, ok := params.Email.Get(); ok {
		updates["email"] = email
	}

	var customer model.Customer
	err := db.Model(&customer).
		Where(`customer_id = ?`, customerID).
		Clauses(clause.Returning{}).
		Updates(updates).
		Error
	if err != nil {
		return model.Customer{}, ferr.ErrDbUnexpected.WithErr(err).WithCtx(ctx)
	}

	return customer, nil
}

type CustomerGetOneParams struct {
	CustomerID    mo.Option[uint]
	CustomerPhone mo.Option[string]
}

type CustomerUpdateParams struct {
	FullName mo.Option[string]
	Email    mo.Option[string]
}
