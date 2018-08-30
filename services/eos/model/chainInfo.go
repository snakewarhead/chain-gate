package model

type ChainInfo struct {
	ServerVersion            string `json:"server_version"`
	HeadBlockNum             int    `json:"head_block_num"`
	LastIrreversibleBlockNum int    `json:"last_irreversible_block_num"`
	HeadBlockID              string `json:"head_block_id"`
	HeadBlockTime            string `json:"head_block_time"`
	HeadBlockProducer        string `json:"head_block_producer"`
}
