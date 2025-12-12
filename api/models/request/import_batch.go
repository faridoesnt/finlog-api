package request

type ImportBatchRequest struct {
	Items []ImportBatchItem `json:"items"`
}

type ImportBatchItem struct {
	Ciphertext string `json:"ciphertext"`
	Nonce      string `json:"nonce"`
	Tag        string `json:"tag"`
	OccurredAt string `json:"occurred_at"`
	IsExpense  bool   `json:"is_expense"`
	CategoryID int64  `json:"category_id"`
}
