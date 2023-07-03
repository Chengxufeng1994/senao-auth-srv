package service

import (
	"context"

	"senao-auth-srv/model"
	"senao-auth-srv/repository"
)

type AccountService interface {
	CreateAccount(ctx context.Context, account *model.Account) error
	GetAccounts(ctx context.Context) ([]*model.Account, error)
	GetAccountsByUsername(ctx context.Context, username string) (*model.Account, error)
	UpdateAccount(ctx context.Context, account *model.Account) (*model.Account, error)

	GetAccountRetryById(ctx context.Context, id string) (string, error)
	UpdateAccountRetryById(ctx context.Context, id string) error
}

func NewAccountServiceImpl(accountRepo repository.AccountRepo) AccountService {
	return &AccountServiceImpl{
		AccountRepo: accountRepo,
	}
}

type AccountServiceImpl struct {
	AccountRepo repository.AccountRepo
}

func (svc *AccountServiceImpl) CreateAccount(ctx context.Context, account *model.Account) error {
	return svc.AccountRepo.CreateAccount(ctx, account)
}

func (svc *AccountServiceImpl) GetAccounts(ctx context.Context) ([]*model.Account, error) {
	return svc.AccountRepo.GetAccounts(ctx)
}

func (svc *AccountServiceImpl) GetAccountsByUsername(ctx context.Context, username string) (*model.Account, error) {
	return svc.AccountRepo.GetAccountsByUsername(ctx, username)
}

func (svc *AccountServiceImpl) UpdateAccount(ctx context.Context, account *model.Account) (*model.Account, error) {
	return svc.AccountRepo.UpdateAccount(ctx, account)
}

func (svc *AccountServiceImpl) GetAccountRetryById(ctx context.Context, id string) (string, error) {
	return svc.AccountRepo.GetAccountRetryById(ctx, id)
}

func (svc *AccountServiceImpl) UpdateAccountRetryById(ctx context.Context, id string) error {
	return svc.AccountRepo.UpdateAccountRetryById(ctx, id)
}
