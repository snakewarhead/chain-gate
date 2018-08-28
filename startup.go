package main

import (
	"net/http"

	"github.com/snakewarhead/chain-gate/utils"
	"github.com/snakewarhead/chain-gate/models"
)

func hello(resp http.ResponseWriter, req *http.Request) {
	resp.Write([]byte("hello"))
	utils.Logger.Debug("hello")

	coin, err := models.GetCoinEnabled()
	if err != nil {
		return
	}
	coin.ID = 3
}

func main() {

	http.HandleFunc("/hello", hello)

	utils.Logger.Info("Http server startup!!")
	http.ListenAndServe(":8080", nil)
	utils.Logger.Info("Http server shutdown!!")
}
