package model

type Region struct {
	Region int `json:"region"`
	CyclesSummary [][]struct {
		ReadLocks    []interface{} `json:"read_locks"`
		WriteLocks   []interface{} `json:"write_locks"`
		Transactions []Transaction `json:"transactions"`
	} `json:"cycles_summary"`
}

type Block struct {
	Previous              string        `json:"previous"`
	Timestamp             string        `json:"timestamp"`
	TransactionMerkleRoot string        `json:"transaction_mroot"`
	BlockMerkleRoot       string        `json:"block_mroot"`
	Producer              string        `json:"producer"`
	ProducerChanges       []interface{} `json:"producer_changes"`
	ProducerSignature     string        `json:"producer_signature"`
	NewProducers          []interface{} `json:"new_producers"`
	Cycles                []interface{} `json:"cycles"`
	ID                    string        `json:"id"`
	BlockNum              int           `json:"block_num"`
	RefBlockPrefix        int           `json:"ref_block_prefix"`
	ScheduleVersion       int           `json:"schedule_version"`
	Regions               []Region      `json:"regions"`
	InputTransactions     []Transaction `json:"input_transactions"`
}
