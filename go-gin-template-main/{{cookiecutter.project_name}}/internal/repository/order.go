package repository

import (
	"context"
	"errors"

	"git-ffd.kz/fmobile/ferr"
	trmgorm "github.com/avito-tech/go-transaction-manager/gorm"
	"github.com/samber/mo"
	"gorm.io/gorm"

	"{{cookiecutter.package_name}}/internal/model"
)

type Order interface {
	GetOne(ctx context.Context, params OrderGetOneParams) (order model.Order, err error)
	GetList(ctx context.Context, params OrderGetListParams) (orders []model.Order, err error)
	Create(ctx context.Context, order model.Order) (model.Order, error)
}

type OrderDB struct {
	db  *gorm.DB
	trx *trmgorm.CtxGetter
}

func NewOrderDB(db *gorm.DB, trx *trmgorm.CtxGetter) *OrderDB {
	return &OrderDB{
		db:  db,
		trx: trx,
	}
}

func (r *OrderDB) GetOne(ctx context.Context, params OrderGetOneParams) (order model.Order, err error) {
	db := r.trx.DefaultTrOrDB(ctx, r.db)

	q := db.Model(&model.Order{})

	if orderID, ok := params.OrderID.Get(); ok {
		q = q.Where(`order_id = ?`, orderID)
	}
	if orderSession, ok := params.OrderSession.Get(); ok {
		q = q.Where(`order_session = ?`, orderSession)
	}

	err = q.First(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return order, ferr.ErrOrderNotFound.WithErr(err).WithCtx(ctx)
		}

		return order, ferr.ErrDbUnexpected.WithErr(err).WithCtx(ctx)
	}

	return order, nil
}

func (r *OrderDB) GetList(ctx context.Context, params OrderGetListParams) (orders []model.Order, err error) {
	db := r.trx.DefaultTrOrDB(ctx, r.db)

	q := db.Model(&model.Order{})

	if len(params.OrderIDs) > 0 {
		q = q.Where(`order_id IN (?)`, params.OrderIDs)
	}
	if len(params.OrderSessions) > 0 {
		q = q.Where(`order_session IN (?)`, params.OrderSessions)
	}
	if limit, ok := params.Limit.Get(); ok {
		q = q.Limit(limit)
	}
	if offset, ok := params.Offset.Get(); ok {
		q = q.Offset(offset)
	}

	err = q.Find(&orders).Error
	if err != nil {
		return nil, ferr.ErrDbUnexpected.WithErr(err).WithCtx(ctx)
	}

	return orders, nil
}

func (r *OrderDB) Create(ctx context.Context, order model.Order) (model.Order, error) {
	db := r.trx.DefaultTrOrDB(ctx, r.db)

	err := db.Create(&order).Error
	if err != nil {
		return model.Order{}, ferr.ErrDbUnexpected.WithErr(err).WithCtx(ctx)
	}

	return order, nil
}

type OrderGetOneParams struct {
	OrderID      mo.Option[uint]
	OrderSession mo.Option[string]
}

type OrderGetListParams struct {
	OrderIDs      []uint
	OrderSessions []string
	Limit         mo.Option[int]
	Offset        mo.Option[int]
}
