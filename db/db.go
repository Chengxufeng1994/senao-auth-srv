package db

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"go.uber.org/fx"

	"senao-auth-srv/model"
	"senao-auth-srv/util"
)

type Db interface {
	CreateAccount(ctx context.Context, account *model.Account) error
	GetAccounts(ctx context.Context) ([]*model.Account, error)
	GetAccountsByUsername(ctx context.Context, username string) (*model.Account, error)
	UpdateAccount(ctx context.Context, account *model.Account) (*model.Account, error)

	GetAccountRetryById(ctx context.Context, id string) (string, error)
	SetAccountRetryById(ctx context.Context, id string) error
}

type DbImpl struct {
	client *redis.Client
}

var Module = fx.Options(
	fx.Provide(
		NewDb,
	),
)

func NewDb(config util.Config) *DbImpl {
	addr := fmt.Sprintf("%s:%d", config.RedisHost, config.RedisPort)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: config.RedisPassword,
		DB:       0,
	})

	return &DbImpl{
		client: client,
	}
}

func (db *DbImpl) CreateAccount(ctx context.Context, account *model.Account) error {
	account.Id = uuid.New().String()
	bytes, err := json.Marshal(&account)
	if err != nil {
		return err
	}
	db.client.HSet("accounts", account.Id, bytes)
	if err != nil {
		return err
	}
	return nil
}

func (db *DbImpl) GetAccounts(ctx context.Context) ([]*model.Account, error) {
	result, err := db.client.HGetAll("accounts").Result()
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

func (db *DbImpl) GetAccountsByUsername(ctx context.Context, username string) (*model.Account, error) {
	accounts, err := db.GetAccounts(ctx)
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

func (db *DbImpl) UpdateAccount(ctx context.Context, account *model.Account) (*model.Account, error) {
	bytes, err := json.Marshal(&account)
	if err != nil {
		return nil, err
	}
	db.client.HSet("accounts", account.Id, bytes)
	if err != nil {
		return nil, err
	}
	return account, nil
}

const RetrySec = 60

func (db *DbImpl) GetAccountRetryById(ctx context.Context, id string) (string, error) {
	return db.client.Get(fmt.Sprintf("accounts:retry:%s", id)).Result()
}

func (db *DbImpl) SetAccountRetryById(ctx context.Context, id string) error {
	_, err := db.client.Set(fmt.Sprintf("accounts:retry:%s", id), "true", RetrySec*time.Second).Result()
	if err != nil {
		return err
	}

	return nil
}
