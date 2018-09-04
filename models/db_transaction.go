package models

import (
	"database/sql"
	"encoding/json"

	"github.com/snakewarhead/chain-gate/utils"
)

type TransactionStatus int
type TransactionDirection int

const (
	InitTransactionStatus      TransactionStatus = iota // 0
	ConfirmedTransactionStatus                          // 1
)

const (
	_                       TransactionDirection = iota
	InTransactionDirection                       // 1
	OutTransactionDirection                      // 2
)

// this model is a db model as well as a json object which is response in http
type Transaction struct {
	ID         int    `json:"id"`
	CoinID     int    `json:"coin_id"`
	Contract   string `json:"contract"`
	TXID       string `json:"txid"`
	IsMain     int    `json:"is_main"`
	Symbol     string `json:"symbol"`
	Direction  int    `json:"direction"`
	Status     int    `json:"status"`
	From       string `json:"from"`
	To         string `json:"to"`
	Amount     string `json:"amount"`
	Fee        string `json:"fee"`
	Memo       string `json:"memo"`
	CreateTime int64  `json:"create_time"`
	UpdateTime int64  `json:"update_time"`
}

func SaveTransaction(coinID int, contract string, isMain bool, txID, symbol, from, to, memo string, amount, fee string, direction TransactionDirection) error {
	stmt, err := utils.DB.Prepare("INSERT INTO transaction_history(coin_id, contract, tx_id, is_main, symbol, from_address, to_address, amount, fee, memo, create_time, update_time, status, direction) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
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

	_, err = stmt.Exec(coinID, contract, txID, isMainN, symbol, from, to, amount, fee, memo, timeStamp, timeStamp, int(InitTransactionStatus), int(direction))
	if err != nil {
		return err
	}
	return nil
}

func FindTransactionsReceiver(direction TransactionDirection, contract, symbol, account, memo string, pos, offset int) ([]Transaction, error) {
	var (
		rows *sql.Rows
		err  error
	)

	// pos, offset is opposite of the cause of sql(limit offset)
	whereCause := " WHERE direction=? and contract=? and symbol=? and to_address=?"
	limitCause := " LIMIT ? OFFSET ?"
	if len(memo) == 0 {
		rows, err = utils.DB.Query(
			"SELECT * FROM transaction_history"+whereCause+" ORDER BY id DESC"+limitCause,
			int(direction),
			contract,
			symbol,
			account,
			offset,
			pos,
		)
	} else {
		whereCause += " and memo=?"
		rows, err = utils.DB.Query(
			"SELECT * FROM transaction_history"+whereCause+" ORDER BY id DESC"+limitCause,
			int(direction),
			contract,
			symbol,
			account,
			memo,
			offset,
			pos,
		)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	trxs := make([]Transaction, 0)
	for rows.Next() {
		t := Transaction{}

		err = rows.Scan(
			&t.ID,
			&t.CoinID,
			&t.Contract,
			&t.TXID,
			&t.IsMain,
			&t.Symbol,
			&t.Direction,
			&t.Status,
			&t.From,
			&t.To,
			&t.Amount,
			&t.Fee,
			&t.Memo,
			&t.CreateTime,
			&t.UpdateTime,
		)
		if err != nil {
			return nil, err
		}

		trxs = append(trxs, t)
	}
	return trxs, nil
}

func TrxsToJson(trxs []Transaction) ([]byte, error) {
	return json.Marshal(trxs)
}
