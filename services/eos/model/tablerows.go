package model

type TableRows struct {
	Rows []struct {
		Account string `json:"account"`
		Balance int    `json:"balance"`
	} `json:"rows"`
	More bool `json:"more"`
}
