package models

import (
	"github.com/snakewarhead/chain-gate/utils"
)

type TransactionStatus int

const (
	Init      TransactionStatus = iota // 0
	Confirmed                          // 1
)

type Transaction struct {
	ID         int
	CoinID     int
	Contract   string
	TXID       int
	IsMain     bool
	Symbol     string
	From       string
	To         string
	Amount     string
	Fee        string
	Memo       string
	CreateTime int64
	UpdateTime int64
	Status     int
}

func SaveTransaction(coinID int, contract string, isMain bool, txID, symbol, from, to, memo string, amount, fee string) error {
	stmt, err := utils.DB.Prepare("INSERT INTO transaction_history(coin_id, contract, tx_id, is_main, symbol, from_address, to_address, amount, fee, memo, create_time, update_time, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	timeStamp := utils.GetCurrentTimestamp()
	var isMainN int
	if isMain {
		isMainN = 1
	} else {
		isMainN = 0
	}

	_, err = stmt.Exec(coinID, contract, txID, isMainN, symbol, from, to, amount, fee, memo, timeStamp, timeStamp, int(Init))
	if err != nil {
		return err
	}
	return nil
}
