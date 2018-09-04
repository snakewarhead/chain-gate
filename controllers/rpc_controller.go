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
		versionURL + "/get_balance":      getBalance,
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

// push a transaction, caller must deal the mapping of coin type and wallet server
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

	if !utils.MustNotEmpty(from, to, isMain, symbol, amount) {
		resp.Write(models.HttpResultToJson(400, "has not enough params", "{}"))
		return
	}

	isMainN, err := strconv.Atoi(isMain)
	if err != nil {
		resp.Write(models.HttpResultToJson(401, "isMain must be a number", "{}"))
		return
	}

	amountD, err := decimal.NewFromString(amount)
	if err != nil || amountD.LessThanOrEqual(decimal.Zero) {
		resp.Write(models.HttpResultToJson(401, "amount must be a positive number", "{}"))
		return
	}
	feeD, err := decimal.NewFromString(fee)
	if err != nil || feeD.LessThan(decimal.Zero) {
		resp.Write(models.HttpResultToJson(401, "fee must be a non-negative number", "{}"))
		return
	}

	txid, err := services.PushTransaction(contract, from, to, memo, symbol, isMainN != 0, amount, fee)
	if err != nil {
		utils.Logger.Error(err)
		resp.Write(models.HttpResultToJson(402, "push transaction error", "{}"))
		return
	}

	// response
	resp.Write(models.HttpResultToJson(200, "success", fmt.Sprintf(`{"txid":"%s"}`, txid)))
}

// get receiver's transactions in DESC, filterd by account, memo, limited by pos and offset
// caller must have the responsibility of dealing the repeat transaction, need to verify the transaction that whether would be dealed
func getTransactions(resp http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	utils.Logger.Debug(req.Form)

	defer recoverResponse(resp)

	direction := req.FormValue("direction")
	contract := req.FormValue("contract")
	symbol := req.FormValue("symbol")
	account := req.FormValue("account")
	memo := req.FormValue("memo") // empty ignore
	pos := req.FormValue("pos")
	offset := req.FormValue("offset")
	if !utils.MustNotEmpty(direction, contract, symbol, account, pos, offset) {
		resp.Write(models.HttpResultToJson(400, "has not enough params", "{}"))
		return
	}

	directionN, err := strconv.Atoi(direction)
	if err != nil || (directionN != 1 && directionN != 2) {
		resp.Write(models.HttpResultToJson(401, "direction must be 1 or 2", "{}"))
		return
	}
	directionE := models.TransactionDirection(directionN)

	posN, err := strconv.Atoi(pos)
	if err != nil || posN < 0 {
		resp.Write(models.HttpResultToJson(401, "pos must be a non-negative number", "{}"))
		return
	}

	offsetN, err := strconv.Atoi(offset)
	if err != nil || offsetN < 1 {
		resp.Write(models.HttpResultToJson(401, "offset must be a non-negative number and more than one", "{}"))
		return
	}

	trxs, err := services.GetTransactionsReceiverFromDB(directionE, contract, symbol, account, memo, posN, offsetN)
	if err != nil {
		utils.Logger.Error(err)
		resp.Write(models.HttpResultToJson(402, "get transactions error", "{}"))
		return
	}

	trxsB, err := models.TrxsToJson(trxs)
	if err != nil {
		utils.Logger.Error(err)
		resp.Write(models.HttpResultToJson(403, "parse transaction to json error", "{}"))
		return
	}

	resp.Write(models.HttpResultToJson(200, "success", string(trxsB)))
}

func getBalance(resp http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	utils.Logger.Debug(req.Form)

	defer recoverResponse(resp)

	contract := req.FormValue("contract")
	symbol := req.FormValue("symbol")
	account := req.FormValue("account")
	if !utils.MustNotEmpty(contract, symbol, account) {
		resp.Write(models.HttpResultToJson(400, "has not enough params", "{}"))
		return
	}

	balance, err := services.GetBalance(contract, account, symbol)
	if err != nil {
		utils.Logger.Error(err)
		resp.Write(models.HttpResultToJson(402, "get balance error", "{}"))
		return
	}

	resp.Write(models.HttpResultToJson(200, "success", fmt.Sprintf(`{"balance":"%s"}`, balance)))
}
