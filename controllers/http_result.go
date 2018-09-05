package controllers

import (
	"fmt"
	// "encoding/json"

	// "github.com/snakewarhead/chain-gate/utils"
)

type httpResult struct {
	code int    `json:"code"`
	msg  string `json:"msg"`
	data string `json:"data"`
}

const innerErrorResult = `{"code":600, "msg":"json parse error", "data":{}`

func HttpResultToJson(code int, msg, data string) []byte {
	r := &httpResult{
		code,
		msg,
		data,
	}
	// Data is string, not a jsonobject
	// bytes, err := json.Marshal(r)
	// if err != nil {
	// 	utils.Logger.Error(err)
	// 	return []byte(innerErrorResult)
	// }
	json := fmt.Sprintf(`{"code":%d, "msg":"%s", "data":%s}`, r.code, r.msg, r.data)
	return []byte(json)
}
