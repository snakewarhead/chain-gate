package model

type Code struct {
	Name     string `json:"name"`
	CodeHash string `json:"code_hash"`
	Wast     string `json:"wast"`
	Abi struct {
		Types []struct {
			NewTypeName string `json:"new_type_name"`
			Type        string `json:"type"`
		} `json:"types"`
		Structs []struct {
			Name string `json:"name"`
			Base string `json:"base"`
			Fields []struct {
				Name string `json:"name"`
				Type string `json:"type"`
			} `json:"fields"`
		} `json:"structs"`
		Actions []struct {
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"actions"`
		Tables []struct {
			Name      string   `json:"name"`
			Type      string   `json:"type"`
			IndexType string   `json:"index_type"`
			KeyNames  []string `json:"key_names"`
			KeyTypes  []string `json:"key_types"`
		} `json:"tables"`
	} `json:"abi"`
}
