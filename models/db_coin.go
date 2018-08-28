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
}

func GetCoinEnabled() (*Coin, error) {
	c := &Coin{}

	const sql = "select * from coin where enable = 1"
	row := utils.DB.QueryRow(sql)
	if err := row.Scan(&c.ID, &c.Name, &c.Enable, &c.MainAddress, &c.Password); err != nil {
		utils.Logger.Critical("must have one enabled coin! %v", err)
		panic(err)
	}

	return c, nil
}
