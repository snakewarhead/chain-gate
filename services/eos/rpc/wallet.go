package rpc

import (
	"strings"

	"github.com/snakewarhead/chain-gate/errors"
	dbmodels "github.com/snakewarhead/chain-gate/models"
	"github.com/snakewarhead/chain-gate/services/eos/model"
	"github.com/snakewarhead/chain-gate/utils"
)

/**
	See if you have eos source {$EOS_SOURCE}/docs/group__eosiorpc.html#walletrpc
	or download from https://eosio.github.io/eos/group__eosiorpc.html
	for detailed specs of:
	- Create
	- Open
	- Lock
	- Unlock
	- LockAll
	- ImportKey
    - ListWallets
	- ListKeys
	- GetPublicKeys
	- SetTimeout
	- SignTrx
*/

func WalletCreate(c *dbmodels.Coin, name string) (string, *errors.AppError) {

	data, err := utils.HTTPPostRawData(c.URL+"/v1/wallet/create", name)

	if err != nil {
		return "", err
	}

	if data == nil {
		return "", errors.NewAppError(nil, "empty response, no key returned by nodeos", -1, nil)
	}

	return string(data), nil
}

func WalletOpen(c *dbmodels.Coin, name string) *errors.AppError {

	_, err := utils.HTTPPostRawData(c.URL+"/v1/wallet/open", name)

	if err != nil {
		return err
	}

	return nil
}

func WalletUnlock(c *dbmodels.Coin, name string, privKey string) *errors.AppError {

	if !strings.HasPrefix(privKey, "\"") {
		privKey = "\"" + privKey + "\""
	}

	_, err := utils.HTTPPostRawData(c.URL+"/v1/wallet/unlock", "[\""+name+"\","+privKey+"]")

	if err != nil {
		return err
	}

	return nil
}

func WalletLock(c *dbmodels.Coin, name string) *errors.AppError {

	_, err := utils.HTTPPostRawData(c.URL+"/v1/wallet/lock", name)

	if err != nil {
		return err
	}

	return nil
}

func WalletLockAll(c *dbmodels.Coin) *errors.AppError {

	_, err := utils.HTTPPostRawData(c.URL+"/v1/wallet/lock_all", "")

	if err != nil {
		return err
	}

	return nil
}

func WalletSignTransaction(c *dbmodels.Coin, trx model.Transaction, pubKeys []string, chainId string) (*model.Transaction, *errors.AppError) {

	trxJson, err := model.TransactionToJSON(&trx)

	// encode keys
	pubKeysRaw := "["

	for i := 0; i < len(pubKeys); i++ {

		pubKeysRaw += "\"" + pubKeys[i] + "\""

		if i+1 < len(pubKeys) {
			pubKeysRaw += ","
		}
	}

	pubKeysRaw += "]"

	raw := "[" + trxJson + "," + pubKeysRaw + ",\"" + chainId + "\"]"

	if err != nil {
		return nil, err
	}

	trxData, err := utils.HTTPPostRawData(c.URL+"/v1/wallet/sign_transaction", raw)

	if err != nil {
		return nil, err
	}

	trxUpdated, err := model.JSONToTransaction(string(trxData))

	if err != nil {
		return nil, err
	}

	return trxUpdated, nil
}
