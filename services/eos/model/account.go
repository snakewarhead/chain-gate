package model

type Account struct {
	Name              string `json:"account_name"`
	EosBalance        string `json:"eos_balance"`
	StakedBalance     string `json:"staked_balance"`
	UnstakingBalance  string `json:"unstaking_balance"`
	LastUnstakingTime string `json:"last_unstaking_time"`
	Permissions []struct {
		Name   string `json:"name"`
		Parent string `json:"parent"`
		RequiredAuth struct {
			Threshold int `json:"threshold"`
			Keys []struct {
				Key    string `json:"key"`
				Weight int    `json:"weight"`
			} `json:"keys"`
			Accounts []interface{} `json:"accounts"`
		} `json:"required_auth"`
	} `json:"permissions"`
}
