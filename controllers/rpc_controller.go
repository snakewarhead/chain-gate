package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/shopspring/decimal"

	"github.com/snakewarhead/chain-gate/models"
	"github.com/snakewarhead/chain-gate/services"
	"github.com/snakewarhead/chain-gate/utils"
)

const versionURL = "/v1"

var (
	urlMapper = map[string]func(resp http.ResponseWriter, req *http.Request){
		"/hello": hello,
		versionURL + "/push_transaction": pushTransaction,
		versionURL + "/get_transactions": getTransactions,
		versionURL + "/get_balance": getBalance,
	}
)

func Startup(c chan int) {
	for k, v := range urlMapper {
		http.HandleFunc(k, v)
	}

	utils.Logger.Info("Http server startup!!")
	err := http.ListenAndServe(":8080", nil)
	utils.Logger.Info("Http server shutdown!!")
	if err != nil {
		utils.Logger.Error(err)
		c <- 1
		return
	}

	c <- 0
}

func recoverResponse(resp http.ResponseWriter) {
	if err := recover(); err != nil {
		utils.Logger.Error(err)
		resp.Write(models.HttpResultToJson(501, "server inner error, contact me", "{}"))
	}
}

func hello(resp http.ResponseWriter, req *http.Request) {
	resp.Write([]byte("hello"))
	utils.Logger.Debug("hello")

}

// push a transaction
// -------------------------------
// request params:
// from
// to
// isMain
// symbol
// amount
// fee
// memo
// -------------------------------
// response
// txid
func pushTransaction(resp http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	utils.Logger.Debug(req.Form)

	defer recoverResponse(resp)

	from := req.FormValue("from")
	to := req.FormValue("to")
	isMain := req.FormValue("isMain")
	symbol := req.FormValue("symbol")
	amount := req.FormValue("amount")
	fee := req.FormValue("fee")
	memo := req.FormValue("memo")
	contract := req.FormValue("contract")

	if len(from) == 0 ||
		len(to) == 0 ||
		len(isMain) == 0 ||
		len(symbol) == 0 ||
		len(amount) == 0 {
		resp.Write(models.HttpResultToJson(400, "has not enough params", "{}"))
		return
	}

	isMainN, err := strconv.Atoi(isMain)
	if err != nil {
		resp.Write(models.HttpResultToJson(401, "must be a number", "{}"))
	}

	amountD, err := decimal.NewFromString(amount)
	if err != nil || amountD.LessThanOrEqual(decimal.Zero) {
		resp.Write(models.HttpResultToJson(401, "must be a positive number", "{}"))
		return
	}
	feeD, err := decimal.NewFromString(fee)
	if err != nil || feeD.LessThan(decimal.Zero) {
		resp.Write(models.HttpResultToJson(401, "must be a non-negative number", "{}"))
		return
	}

	txid, err := services.PushTransaction(contract, from, to, memo, symbol, isMainN != 0, amount, fee)
	if err != nil {
		utils.Logger.Error(err)
		resp.Write(models.HttpResultToJson(402, "push transaction error", "{}"))
	}

	// response
	resp.Write(models.HttpResultToJson(200, "success", fmt.Sprintf(`{"txid":"%s"}`, txid)))
}

func getTransactions(resp http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	utils.Logger.Debug(req.Form)

	defer recoverResponse(resp)
}

func getBalance(resp http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	utils.Logger.Debug(req.Form)

	defer recoverResponse(resp)

	contract := req.FormValue("contract")
	symbol := req.FormValue("symbol")
	account := req.FormValue("account")
	if len(contract) == 0 || len(symbol) == 0 || len(account) == 0 {
		resp.Write(models.HttpResultToJson(400, "has not enough params", "{}"))
		return
	}

	balance, err := services.GetBalance(contract, account, symbol)
	if err != nil {
		utils.Logger.Error(err)
		resp.Write(models.HttpResultToJson(402, "get balance error", "{}"))
	}
	
	resp.Write(models.HttpResultToJson(200, "success", fmt.Sprintf(`{"balance":"%s"}`, balance)))
}
