package services

import (
	"fmt"
	"strings"
	"github.com/snakewarhead/chain-gate/utils"
	eos "github.com/snakewarhead/eos-go"
	eostoken "github.com/snakewarhead/eos-go/token"

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

func (n *nodeEOS) initAPI() *eos.API {
	api := eos.New(n.coin.APIURL, n.coin.APIWalletURL)
	api.Signer = eos.NewWalletSigner(api, "default")
	api.Debug = false

	return api
}

// note: the precision of amount must be the same as when the token deployed, like: 100000.0000 HP, so amount is 12.0000
func (n *nodeEOS) pushTransaction(contract, from, to, memo, symbol string, isMain bool, amount, fee string) (string, error) {
	api := n.initAPI()

	api.WalletUnlock("default", n.coin.Password)
	defer func() {
		api.WalletLock("default")
	}()

	quantity, err := eos.NewAsset(amount + " " + symbol)
	if err != nil {
		return "", err
	}

	action := eostoken.NewTransferCommon(eos.AN(contract), eos.AN(from), eos.AN(to), quantity, memo)
	pushed, err := api.SignPushActions(action)
	utils.Logger.Info("pushed action:%v, error:%v", pushed, err)
	if err != nil {
		return "", err
	}

	return pushed.TransactionID, nil
}

func (n *nodeEOS) getBalance(contract, account, symbol string) (string, error) {
	api := n.initAPI()

	assets, err := api.GetCurrencyBalance(eos.AN(account), symbol, eos.AN(contract))
	if err != nil {
		return "", err
	}
	if len(assets) != 1 {
		return "", fmt.Errorf("assets must have one result. %v", assets)
	}
	utils.Logger.Debug("balance : %s", assets[0].String())

	return strings.Split(assets[0].String(), " ")[0], nil
}