package model

type WeightedKey struct {
	Key    string `json:"key"`
	Weight int    `json:"weight"`
}

type Autorithy struct {
	Threshold int           `json:"threshold"`
	Keys      []WeightedKey `json:"keys"`
	Accounts  []interface{} `json:"accounts"`
}

func NewAuthority(key string, weight int) (*Autorithy) {

	wKey := WeightedKey{
		key,
		weight,
	}

	return &Autorithy{
		1,
		[]WeightedKey{wKey},
		[]interface{}{},
	}
}
