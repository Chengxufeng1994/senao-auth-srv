package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"senao-auth-srv/model"
	"time"
)

const RetrySec = 60

type createAccountRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=8,max=32"`
}

type creatAccountResponse struct {
	Success bool   `json:"success"`
	Reason  string `json:"reason"`
}

func ValidationErrorToText(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "max":
		return fmt.Sprintf("%s cannot be longer than %s", e.Field(), e.Param())
	case "min":
		return fmt.Sprintf("%s must be longer than %s", e.Field(), e.Param())
		//case "email":
		//	return fmt.Sprintf("Invalid email format")
		//case "len":
		//	return fmt.Sprintf("%s must be %s characters long", e.Field(), e.Param())
	}
	return fmt.Sprintf("%s is not valid", e.Field())
}

func (srv *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	var res creatAccountResponse
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ve := err.(validator.ValidationErrors)
		e := ve[0]
		res.Success = false
		res.Reason = ValidationErrorToText(e)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	accounts, _ := srv.database.GetAccounts()
	for _, account := range accounts {
		if account.Username == req.Username {
			res.Success = false
			res.Reason = "Username already exists"
			ctx.JSON(http.StatusBadRequest, res)
			return
		}
	}
	// TODO: Password validator

	account := model.Account{
		Username:    req.Username,
		Password:    req.Password,
		FailedCount: 0,
	}
	_, err := srv.database.CreateAccount(&account)
	if err != nil {
		res.Success = false
		res.Reason = "CreateAccount Failed, " + err.Error()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res.Success = true
	res.Reason = ""
	ctx.JSON(http.StatusOK, res)
}

type verifyAccountRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=8,max=32"`
}

type verifyAccountResponse struct {
	Success bool   `json:"success"`
	Reason  string `json:"reason"`
}

func (srv *Server) verifyAccount(ctx *gin.Context) {
	var req verifyAccountRequest
	var res verifyAccountResponse
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ve := err.(validator.ValidationErrors)
		e := ve[0]
		res.Success = false
		res.Reason = ValidationErrorToText(e)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	existedAccount, err := srv.database.GetAccountsByUsername(req.Username)
	if err != nil {
		res.Success = false
		res.Reason = err.Error()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	if existedAccount.FailedCount >= 5 {
		now := time.Now().Unix()
		if now < existedAccount.FailedExpireSec {
			res.Success = false
			res.Reason = "Pleas try again after one minutes"
			ctx.JSON(http.StatusUnauthorized, res)
			return
		} else {
			existedAccount.FailedCount = 0
		}
	}

	if existedAccount.Password != req.Password {
		existedAccount.FailedCount++
		if existedAccount.FailedCount >= 5 {
			existedAccount.FailedExpireSec = time.Now().Unix() + RetrySec
		}
		srv.database.UpdateAccount(existedAccount)
		res.Success = false
		res.Reason = "Unauthorized"
		ctx.JSON(http.StatusUnauthorized, res)
		return
	}

	existedAccount.FailedExpireSec = 0
	existedAccount.FailedCount = 0
	srv.database.UpdateAccount(existedAccount)
	ctx.JSON(http.StatusOK, existedAccount)
}
