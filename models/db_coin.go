package models

import (
	"github.com/snakewarhead/chain-gate/utils"
)

type Coin struct {
	ID           int
	Name         string
	Enable       bool
	MainAddress  string
	PublicKey    string
	Password     string
	APIURL       string
	APIWalletURL string
	ConfirmNum   int
}

func GetCoinEnabled() (*Coin, error) {
	c := &Coin{}

	const sql = "SELECT * FROM coin WHERE enable = 1"
	row := utils.DB.QueryRow(sql)
	err := row.Scan(&c.ID, &c.Name, &c.Enable, &c.MainAddress, &c.PublicKey, &c.Password, &c.APIURL, &c.APIWalletURL, &c.ConfirmNum)
	return c, err
}
