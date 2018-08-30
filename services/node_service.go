package services

import (
	"github.com/shopspring/decimal"
	"github.com/snakewarhead/chain-gate/models"
	"github.com/snakewarhead/chain-gate/utils"
)

var (
	nodeManager = make(anodeManager)
	nodeCurrent inode
)

type inode interface {
	id() int
	bind(c *models.Coin)
	getBind() *models.Coin
	pushTransaction(from, to, memo, symbol string, isMain bool, amount, fee decimal.Decimal) (string, error)
}

type anodeManager map[int]inode

func init() {
	// 注册所有node
	n := &nodeEOS{}
	nodeManager[n.id()] = n

	// TODO: add others
}

func Startup() {
	var err error
	coin, err := models.GetCoinEnabled()
	if err != nil {
		utils.Logger.Critical("must have one enabled coin! %v", err)
		panic(err)
	}

	// 找到对应node
	nodeCurrent = nodeManager[coin.ID]
	nodeCurrent.bind(coin)
}

func PushTransaction(from, to, memo, symbol string, isMain bool, amount, fee decimal.Decimal) (string, error) {
	utils.Logger.Info("PushTransaction 1 ----------------- from:%s, to:%s, memo:%s, symbol:%s, isMain:%t, amount:%s, fee:%s",
		from,
		to,
		memo,
		symbol,
		isMain,
		amount.String(),
		fee.String())

	txid, err := nodeCurrent.pushTransaction(from, to, memo, symbol, isMain, amount, fee)
	utils.Logger.Info("PushTransaction 2 ----------------- txid:%s, err:%v", txid, err)

	// persistent
	errPersistent := models.SaveTransaction(
		nodeCurrent.getBind().ID,
		isMain,
		txid,
		symbol,
		from,
		to,
		memo,
		amount,
		fee,
	)
	if errPersistent != nil {
		// transaction is success, it must be response, ignore the persistent error
		utils.Logger.Error("PushTransaction persistent ----------------- txid:%s, err:%v", txid, err)
	}

	return txid, err
}
