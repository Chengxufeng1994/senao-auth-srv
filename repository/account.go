package repository

import (
	"context"

	"go.uber.org/fx"

	"senao-auth-srv/db"
	"senao-auth-srv/model"
)

type AccountRepo interface {
	CreateAccount(ctx context.Context, account *model.Account) error
	GetAccounts(ctx context.Context) ([]*model.Account, error)
	GetAccountsByUsername(ctx context.Context, username string) (*model.Account, error)
	UpdateAccount(ctx context.Context, account *model.Account) (*model.Account, error)

	GetAccountRetryById(ctx context.Context, id string) (string, error)
	UpdateAccountRetryById(ctx context.Context, id string) error
}

var AccountRepoModule = fx.Options(
	fx.Provide(
		NewAccountRepoImpl,
	),
)

func NewAccountRepoImpl(store *db.DbImpl) AccountRepo {
	return &AccountRepoImpl{
		store: store,
	}
}

type AccountRepoImpl struct {
	store *db.DbImpl
}

func (repo *AccountRepoImpl) CreateAccount(ctx context.Context, account *model.Account) error {
	return repo.store.CreateAccount(ctx, account)
}

func (repo *AccountRepoImpl) GetAccounts(ctx context.Context) ([]*model.Account, error) {
	return repo.store.GetAccounts(ctx)
}

func (repo *AccountRepoImpl) GetAccountsByUsername(ctx context.Context, username string) (*model.Account, error) {
	return repo.store.GetAccountsByUsername(ctx, username)
}

func (repo *AccountRepoImpl) UpdateAccount(ctx context.Context, account *model.Account) (*model.Account, error) {
	return repo.store.UpdateAccount(ctx, account)
}

func (repo *AccountRepoImpl) GetAccountRetryById(ctx context.Context, id string) (string, error) {
	return repo.store.GetAccountRetryById(ctx, id)
}

func (repo *AccountRepoImpl) UpdateAccountRetryById(ctx context.Context, id string) error {
	return repo.store.SetAccountRetryById(ctx, id)
}
