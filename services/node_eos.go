package services

import (
	"fmt"
	"strings"

	"github.com/snakewarhead/chain-gate/utils"
	eos "github.com/snakewarhead/eos-go"
	eostoken "github.com/snakewarhead/eos-go/token"

	"github.com/snakewarhead/chain-gate/models"
)

const (
	mainContract = "eosio.token"
	actionType   = "transfer"

	offsetQuery  = 9	// [pos, pos + offset]
)

type nodeEOS struct {
	coin *models.Coin
	api  *eos.API
}

func (n *nodeEOS) id() int {
	return 1
}

func (n *nodeEOS) bind(c *models.Coin) {
	n.coin = c
	n.initAPI()
}

func (n *nodeEOS) getBind() *models.Coin {
	return n.coin
}

func (n *nodeEOS) initAPI() {
	api := eos.New(n.coin.APIURL, n.coin.APIWalletURL)
	api.Signer = eos.NewWalletSigner(api, "default")
	api.Debug = false

	n.api = api
}

// note: the precision of amount must be the same as when the token deployed, like: 100000.0000 HP, so amount is 12.0000
func (n *nodeEOS) pushTransaction(contract, from, to, memo, symbol string, isMain bool, amount, fee string) (string, error) {
	api := n.api

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
	api := n.api

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

func (n *nodeEOS) obversing() {
	// must add this, it can't be shut down
	defer utils.RecoverAndLog("node_eos", "obversing")

	api := n.api

	// get actions from begin to end obversed
	resp, err := api.GetActions(eos.GetActionsRequest{eos.AN(n.coin.MainAddress), n.coin.PositionQuery, offsetQuery})
	if err != nil {
		utils.Logger.Error("api.GetActions --- %v", err)
		return
	}

	if len(resp.Actions) == 0 {
		return
	}

	// this is the block hight, it will be checked for the confirmed transactions
	lastIrreversibleBlock := int64(resp.LastIrreversibleBlock)

	// record all transaction about our address
	for _, actionResp := range resp.Actions {
		for {
			action := actionResp.Trace.Action
			if strings.Compare(actionType, string(action.Name)) != 0 {
				break
			}

			blockNum := int64(actionResp.BlockNum)
			trxid := actionResp.Trace.TransactionID.Encode()
			contract := string(action.Account)

			data := action.Data.(map[string]interface{})
			from := data["from"].(string)
			to := data["to"].(string)
			memo := data["memo"].(string)

			// "1.1111 HP"
			quantity := data["quantity"].(string)
			qu := strings.Split(quantity, " ")
			amount := qu[0]
			symbol := qu[1]

			// need to calculate params
			var isMain int
			if contract == mainContract {
				isMain = 1
			} else {
				isMain = 0
			}

			// send to main account is in direction, otherwise is out direction
			var direction models.TransactionDirection
			if from == n.coin.MainAddress {
				direction = models.OutTransactionDirection
			} else if to == n.coin.MainAddress {
				direction = models.InTransactionDirection
			}

			// status of confirmed must grow up over confrimed num
			var status models.TransactionStatus
			if lastIrreversibleBlock-blockNum >= n.coin.ConfirmNum {
				status = models.ConfirmedTransactionStatus
			} else {
				status = models.InitTransactionStatus
			}

			// is already in db ?
			trx, err := models.FindOneTransaction(n.coin.ID, trxid)
			if err != nil {
				utils.Logger.Error("FindOneTransaction has error: %s, %v", trxid, err)
				// retry, because an exception has bean throwed in db
				return
			}
			if trx != nil {
				// utils.Logger.Error("transaction find one --- %s", trxid)
				break
			}

			trx = &models.Transaction{}
			trx.CoinID = n.coin.ID
			trx.Contract = contract
			trx.TXID = trxid
			trx.IsMain = isMain
			trx.Symbol = symbol
			trx.Direction = int(direction)
			trx.Status = int(status)
			trx.From = from
			trx.To = to
			trx.Amount = amount
			trx.Fee = "0"
			trx.Memo = memo
			trx.BlockNum = blockNum

			err = models.SaveTransaction(trx)
			if err != nil {
				utils.Logger.Error("SaveTransaction has error: %s, %v", trxid, err)
				// retry, because an exception has bean throwed in db
				return
			}

			// must break
			break
		}

		// record position
		n.coin.PositionQuery++
		err := models.UpdateCoin(n.coin)
		if err != nil {
			utils.Logger.Error("UpdateCoin --- %v", err)
			// retry, because an exception has bean throwed in db
			n.coin.PositionQuery--
			return
		}
	}
}

func (n *nodeEOS) obversingConfirmed() {
	// must add this, it can't be shut down
	defer utils.RecoverAndLog("node_eos", "obversing")

	// find all initing transactions from db
	trxs, err := models.FindTransactions(n.coin.ID, 0, 1, mainContract, n.coin.Name, "", "", 100, 0)
	if err != nil {
		utils.Logger.Error("FindTransactions has error: %v", err)
		return
	}
	if len(trxs) == 0 {
		return
	}

	info, err := n.api.GetInfo()
	if err != nil {
		utils.Logger.Error("api.GetInfo has error: %v", err)
		return
	}
	lastIrreversibleBlockNum := int64(info.LastIrreversibleBlockNum)

	for _, t := range trxs {
		// check whether transaction has bean confrimed
		if (lastIrreversibleBlockNum - t.BlockNum >= n.coin.ConfirmNum) {
			t.Status = int(models.ConfirmedTransactionStatus)
			// update it
			models.SaveOrUpdateTransaction(t)
		}
	}

}
