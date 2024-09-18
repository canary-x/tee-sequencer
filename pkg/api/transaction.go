package api

type Transaction struct {
	TxHash  string `json:"tx_hash"`
	Account string `json:"account"`
	Nonce   uint64 `json:"nonce"`
}

type TransactionBatch struct {
	Transactions []Transaction `json:"transactions"`
}

type TransactionBatchSorted struct {
	Transactions []Transaction `json:"transactions"`
}
