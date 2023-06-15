package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"senao-auth-srv/model"
	"senao-auth-srv/util"
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

// CreateAccount godoc
// @Summary create account
// @Schemes
// @Description create account following parameters
// @Tags 	account
// @Accept  json
// @Produce json
// @Param   createAccountRequest body createAccountRequest true "create account parameters"
// @Success 200 {object} creatAccountResponse
// @Failed 	200 {object} creatAccountResponse
// @Router 	/register [post]
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

	isValidatedPassword := util.ValidatePassword(req.Password)
	if !isValidatedPassword {
		res.Success = false
		res.Reason = "Password is not valid"
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		res.Success = false
		res.Reason = "Hashed password failed"
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}
	account := model.Account{
		Username:    req.Username,
		Password:    hashedPassword,
		FailedCount: 0,
	}
	_, err = srv.database.CreateAccount(&account)
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

// VerifyAccount godoc
// @Summary verify account
// @Schemes
// @Description verify account
// @Tags 	account
// @Accept  json
// @Produce json
// @Param   verifyAccountRequest body verifyAccountRequest true "create account parameters"
// @Success 200 {object} verifyAccountResponse
// @Router 	/verify [post]
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
		result, err := srv.database.Client.Get(fmt.Sprintf("accounts:retry:%s", existedAccount.Id)).Result()
		if err != nil {
			if err.Error() != "redis: nil" {
				res.Success = false
				res.Reason = fmt.Sprintf("Redis get accounts:retry:%s", err.Error())
				ctx.JSON(http.StatusUnauthorized, res)
				return
			}
		}

		if result == "true" {
			res.Success = false
			res.Reason = "Pleas try again after one minutes"
			ctx.JSON(http.StatusUnauthorized, res)
			return
		}

		existedAccount.FailedCount = 0
	}

	err = util.CheckPassword(req.Password, existedAccount.Password)
	if err != nil {
		existedAccount.FailedCount++
		srv.database.UpdateAccount(existedAccount)
		srv.database.Client.Set(fmt.Sprintf("accounts:retry:%s", existedAccount.Id), "true", RetrySec*time.Second)
		res.Success = false
		res.Reason = "Unauthorized"
		ctx.JSON(http.StatusUnauthorized, res)
		return
	}

	existedAccount.FailedExpireSec = 0
	existedAccount.FailedCount = 0
	srv.database.UpdateAccount(existedAccount)
	res.Success = true
	ctx.JSON(http.StatusOK, res)
}
