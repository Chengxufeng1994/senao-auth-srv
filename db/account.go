package db

import (
	"fmt"

	"github.com/goccy/go-json"
	"github.com/google/uuid"

	"senao-auth-srv/model"
)

func (db *Database) CreateAccount(account *model.Account) (*model.Account, error) {
	account.Id = uuid.New().String()
	bytes, err := json.Marshal(&account)
	if err != nil {
		return nil, err
	}
	db.Client.HSet("accounts", account.Id, bytes)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (db *Database) GetAccounts() ([]*model.Account, error) {
	result, err := db.Client.HGetAll("accounts").Result()
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

func (db *Database) GetAccountsByUsername(username string) (*model.Account, error) {
	accounts, err := db.GetAccounts()
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

func (db *Database) UpdateAccount(account *model.Account) (*model.Account, error) {
	bytes, err := json.Marshal(&account)
	if err != nil {
		return nil, err
	}
	db.Client.HSet("accounts", account.Id, bytes)
	if err != nil {
		return nil, err
	}
	return account, nil
}
