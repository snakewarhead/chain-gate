package services

import (
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
	pushTransaction(contract, from, to, memo, symbol string, isMain bool, amount, fee string) (string, error)
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

func PushTransaction(contract, from, to, memo, symbol string, isMain bool, amount, fee string) (string, error) {
	utils.Logger.Info("PushTransaction 1 ----------------- contract:%s, from:%s, to:%s, memo:%s, symbol:%s, isMain:%t, amount:%s, fee:%s",
		contract,
		from,
		to,
		memo,
		symbol,
		isMain,
		amount,
		fee)

	txid, err := nodeCurrent.pushTransaction(contract, from, to, memo, symbol, isMain, amount, fee)
	utils.Logger.Info("PushTransaction 2 ----------------- txid:%s, err:%v", txid, err)
	if err != nil {
		return "", err
	}

	// persistent
	errPersistent := models.SaveTransaction(
		nodeCurrent.getBind().ID,
		contract,
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
		utils.Logger.Error("PushTransaction persistent ----------------- txid:%s, err:%v", txid, errPersistent)
	}

	// transaction is success, it must be response, ignore the other errors
	return txid, nil
}
