package models

import (
	"github.com/snakewarhead/chain-gate/utils"
)

type Coin struct {
	ID          int
	Name        string
	Enable      bool
	MainAddress string
	Password    string
	URL			string
}

func GetCoinEnabled() (*Coin, error) {
	c := &Coin{}

	const sql = "SELECT * FROM coin WHERE enable = 1"
	row := utils.DB.QueryRow(sql)
	err := row.Scan(&c.ID, &c.Name, &c.Enable, &c.MainAddress, &c.Password, &c.URL)
	return c,err 
}
