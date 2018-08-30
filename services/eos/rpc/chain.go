package rpc

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/snakewarhead/chain-gate/errors"
	dbmodels "github.com/snakewarhead/chain-gate/models"
	"github.com/snakewarhead/chain-gate/services/eos/model"
	"github.com/snakewarhead/chain-gate/utils"
)

/**
See if you have EOS source {$EOS_SOURCE}/docs/group__eosiorpc.html#chainrpc for detailed specs of:
or download from https://eosio.github.io/eos/group__eosiorpc.html
- GetInfo
- GetBlock
- GetAccount
- GetCode
- GetTableRows
- AbiJSONToBin
- AbiBinToJSON
- PushTransaction
- //TODO: PushTransactions
- GetRequiredKeys

*/

func ChainGetInfo(c *dbmodels.Coin) (*model.ChainInfo, *errors.AppError) {

	data, err := utils.HTTPGet(c.URL+"/v1/chain/get_info", map[string]string{})

	if err != nil {
		return nil, err
	}

	chainInfo := model.ChainInfo{}
	errM := json.Unmarshal(data, &chainInfo)

	if errM != nil {
		return nil, errors.NewAppError(nil, "cannot parse result", -1, nil)
	}

	return &chainInfo, nil
}

func ChainGetBlock(c *dbmodels.Coin, blockNumOrId string) (*model.Block, *errors.AppError) {

	data, err := utils.HTTPPost(c.URL+"/v1/chain/get_block", map[string]interface{}{"block_num_or_id": blockNumOrId}, nil)

	if err != nil {
		return nil, err
	}

	block := model.Block{}
	errM := json.Unmarshal(data, &block)

	if errM != nil {
		return nil, errors.NewAppError(nil, "cannot parse result", -1, nil)
	}

	return &block, nil
}

func ChainAbiJSONToBin(c *dbmodels.Coin, abiJSON *model.AbiJSON) (*model.Abi, *errors.AppError) {

	bin, err := model.AbiJSONToBytes(abiJSON)

	data, err := utils.HTTPPost(c.URL+"/v1/chain/abi_json_to_bin", nil, bin)

	if err != nil {
		return nil, err
	}

	abiBin := model.AbiBin{}
	errM := json.Unmarshal(data, &abiBin)

	if errM != nil {
		return nil, errors.NewAppError(nil, "cannot parse result", -1, nil)
	}

	abi := model.Abi{
		*abiJSON,
		abiBin,
	}

	return &abi, nil
}

func ChainAbiBinToJSON(c *dbmodels.Coin, abi *model.Abi) (*model.Abi, *errors.AppError) {

	bin, err := model.AbiToBytes(abi)

	data, err := utils.HTTPPost(c.URL+"/v1/chain/abi_bin_to_json", nil, bin)

	if err != nil {
		return nil, err
	}

	abiJSON := model.AbiJSON{}

	errM := json.Unmarshal(data, &abiJSON)

	if errM != nil {
		return nil, errors.NewAppError(nil, "cannot parse result", -1, nil)
	}

	abi.AbiJSON = abiJSON

	return abi, nil
}

func ChainPushTransaction(c *dbmodels.Coin, trx model.Transaction, pubKeys []string, chainId string) (*model.Transaction, *errors.AppError) {
	const transactionExpirationDelay = 30

	chainInfo, err := ChainGetInfo(c)
	block, err := ChainGetBlock(c, strconv.Itoa(chainInfo.LastIrreversibleBlockNum))

	trx.RefBlockNum = chainInfo.LastIrreversibleBlockNum
	trx.RefBlockPrefix = int64(block.RefBlockPrefix)

	actions := trx.Actions

	// calculate expiration date
	time := time.Now().UTC().Add(time.Duration(transactionExpirationDelay * 1000 * 1000 * 1000))
	trx.Expiration = time.Format("2006-01-02T15:04:05")

	fmt.Println("time:", trx.Expiration)

	// calculate HEX data for each action
	for i := 0; i < len(trx.Actions); i++ {

		abiJSON := model.AbiJSON{
			trx.Actions[i].Account,
			trx.Actions[i].Name,
			trx.Actions[i].Args,
		}

		data, err := ChainAbiJSONToBin(c, &abiJSON)

		if err != nil {
			return nil, err
		}

		trx.Actions[i].Data = string(data.Binargs)
	}

	// sign transaction
	trxSigned, err := WalletSignTransaction(c, trx, pubKeys, chainId)

	if err != nil {
		return nil, err
	}

	trx.Signatures = trxSigned.Signatures
	trx.Actions = actions

	// encode trx

	raw, err := model.TransactionToJSON(&trx)

	if err != nil {
		return nil, err
	}

	// encode signatures
	signatures := "["

	for i := 0; i < len(trx.Signatures); i++ {

		signatures += "\"" + trx.Signatures[i] + "\""

		if i+1 < len(trx.Signatures) {
			signatures += ","
		}
	}

	signatures += "]"

	trxData, err := utils.HTTPPostRawData(c.URL+"/v1/chain/push_transaction", "{\"signatures\":"+signatures+",\"transaction\":"+raw+",\"compression\":\"none\"}")

	if err != nil {
		return nil, err
	}

	var trxDataMap map[string]interface{}
	dec := json.NewDecoder(strings.NewReader(string(trxData)))
	errD := dec.Decode(&trxDataMap)

	if errD != nil {
		return nil, errors.NewAppError(nil, "cannot parse transaction data from: "+string(trxData), -1, nil)
	}

	trx.ID = trxDataMap["transaction_id"].(string)

	return &trx, nil
}
