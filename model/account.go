package model

type Account struct {
	Username        string
	Password        string
	FailedCount     int8
	FailedExpireSec int64
}
