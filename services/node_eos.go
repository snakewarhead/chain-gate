package services

import (
	"github.com/shopspring/decimal"
	"github.com/snakewarhead/chain-gate/models"
)

type nodeEOS struct {
	coin *models.Coin
}

func (n *nodeEOS) id() int {
	return 1
}

func (n *nodeEOS) bind(c *models.Coin) {
	n.coin = c
}

func (n *nodeEOS) getBind() *models.Coin {
	return n.coin
}

func (n *nodeEOS) pushTransaction(from, to, memo, symbol string, isMain bool, amount, fee decimal.Decimal) (string, error) {
	return "asdfdf", nil
}
