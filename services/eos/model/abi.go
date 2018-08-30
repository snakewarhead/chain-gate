package model

import (
	"encoding/json"
	"github.com/snakewarhead/chain-gate/errors"
)

type Abi struct {
	AbiJSON
	AbiBin
}

type AbiJSON struct {
	Code   string                 `json:"code"`
	Action string                 `json:"action"`
	Args   map[string]interface{} `json:"args"`
}

type AbiBin struct {
	Binargs       string        `json:"binargs"`
	RequiredScope []interface{} `json:"required_scope"`
	RequiredAuth  []interface{} `json:"required_auth"`
}

func AbiToBytes(obj *Abi) ([]byte, *errors.AppError) {

	bytes, err := json.Marshal(&obj)

	if err != nil {
		return nil, errors.NewAppError(err, "cannot marshal AbiBin", -1, nil)
	}

	return bytes, nil
}

func AbiJSONToBytes(obj *AbiJSON) ([]byte, *errors.AppError) {

	bytes, err := json.Marshal(&obj)

	if err != nil {
		return nil, errors.NewAppError(err, "cannot marshal AbiJSON", -1, nil)
	}

	return bytes, nil
}
