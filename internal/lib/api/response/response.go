package response

type RespTransaction struct {
	TransactionTime string `json:"transactionTime,omitempty"`
	FromWallet      string `json:"fromWallet,omitempty"`
	ToWallet        string `json:"toWallet,omitempty"`
	Amount          string `json:"amount,omitempty"`
}

type Wallet struct {
	WalletID string `json:"wallet_id,omitempty"`
	Balance  string `json:"balance,omitempty"`
}
