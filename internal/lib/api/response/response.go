package response

type RespTransaction struct {
	TransactionTime string `json:"transactionTime"`
	FromWallet      string `json:"fromWallet"`
	ToWallet        string `json:"toWallet"`
	Amount          string `json:"amount"`
}

type Wallet struct {
	WalletID string `json:"wallet_id"`
	Balance  string `json:"balance"`
}

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
	Wallet
}

const (
	StatusOK    = "200"
	StatusError = "Error"
)

func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

func ResError(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}
