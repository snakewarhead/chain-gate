package main

import (
	"github.com/snakewarhead/chain-gate/controllers"
	"github.com/snakewarhead/chain-gate/utils"
	"github.com/snakewarhead/chain-gate/services"
)

func main() {
	utils.Logger.Info("startup service 1------------------------------------")
	services.Startup()
	utils.Logger.Info("startup service 2------------------------------------")

	utils.Logger.Info("startup http server 1------------------------------------")
	c := make(chan int)
	go controllers.Startup(c)
	utils.Logger.Info("startup http server 2------------------------------------")

	end := <- c
	utils.Logger.Info("startup end %d------------------------------------", end)
}
