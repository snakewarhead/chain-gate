package models

import (
	"github.com/snakewarhead/chain-gate/utils"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	ID         int
	CoinID     int
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
}

func SaveTransaction(coinID int, isMain bool, txID, symbol, from, to, memo string, amount, fee decimal.Decimal) error {
	stmt, err := utils.DB.Prepare("INSERT INTO transaction_history(coin_id, tx_id, is_main, symbol, from_address, to_address, amount, fee, memo, create_time, update_time) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
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

	_, err = stmt.Exec(coinID, txID, isMainN, symbol, from, to, amount.String(), fee.String(), memo, timeStamp, timeStamp)
	if err != nil {
		return err
	}
	return nil
}
