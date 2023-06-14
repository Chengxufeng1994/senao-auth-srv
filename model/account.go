package model

type Account struct {
	Id              string `json:"id"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	FailedCount     int8   `json:"failed_count"`
	FailedExpireSec int64  `json:"failed_expire_sec"`
}
