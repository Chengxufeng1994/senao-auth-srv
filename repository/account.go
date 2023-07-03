package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"senao-auth-srv/db"
	"senao-auth-srv/model"
)

const RetrySec = 60

type AccountRepo interface {
	CreateAccount(ctx context.Context, account *model.Account) error
	GetAccounts(ctx context.Context) ([]*model.Account, error)
	GetAccountsByUsername(ctx context.Context, username string) (*model.Account, error)
	UpdateAccount(ctx context.Context, account *model.Account) (*model.Account, error)

	GetAccountRetryById(ctx context.Context, id string) (string, error)
	UpdateAccountRetryById(ctx context.Context, id string) error
}

func NewAccountRepoImpl(store *db.Database) AccountRepo {
	return &AccountRepoImpl{
		store: store,
	}
}

type AccountRepoImpl struct {
	store *db.Database
}

func (repo *AccountRepoImpl) CreateAccount(ctx context.Context, account *model.Account) error {
	account.Id = uuid.New().String()
	bytes, err := json.Marshal(&account)
	if err != nil {
		return err
	}
	repo.store.Client.HSet("accounts", account.Id, bytes)
	if err != nil {
		return err
	}

	return nil
}

func (repo *AccountRepoImpl) GetAccounts(ctx context.Context) ([]*model.Account, error) {
	result, err := repo.store.Client.HGetAll("accounts").Result()
	if err != nil {
		return nil, err
	}
	accounts := []*model.Account{}
	for _, data := range result {
		account := &model.Account{}
		if err = json.Unmarshal([]byte(data), &account); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (repo *AccountRepoImpl) GetAccountsByUsername(ctx context.Context, username string) (*model.Account, error) {
	accounts, err := repo.GetAccounts(ctx)
	if err != nil {
		return nil, err
	}
	for _, account := range accounts {
		if account.Username == username {
			return account, nil
		}
	}

	return nil, fmt.Errorf("%s not found", username)
}

func (repo *AccountRepoImpl) UpdateAccount(ctx context.Context, account *model.Account) (*model.Account, error) {
	bytes, err := json.Marshal(&account)
	if err != nil {
		return nil, err
	}
	repo.store.Client.HSet("accounts", account.Id, bytes)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (repo *AccountRepoImpl) GetAccountRetryById(ctx context.Context, id string) (string, error) {
	return repo.store.Client.Get(fmt.Sprintf("accounts:retry:%s", id)).Result()
}

func (repo *AccountRepoImpl) UpdateAccountRetryById(ctx context.Context, id string) error {
	_, err := repo.store.Client.Set(fmt.Sprintf("accounts:retry:%s", id), "true", RetrySec*time.Second).Result()
	if err != nil {
		return err
	}

	return nil
}
