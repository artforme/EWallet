package response

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

type RespTransaction struct {
	TransactionTime string `json:"transactionTime"`
	FromWallet      string `json:"fromWallet"`
	ToWallet        string `json:"toWallet"`
	Amount          string `json:"amount"`
}

type Wallet struct {
	WalletID string `json:"wallet_id,omitempty"`
	Balance  string `json:"balance,omitempty"`
}

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK            = "200"
	StatusError         = "400"
	StatusErrorNotFound = "404"
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

func ResErrorNotFound(msg string) Response {
	return Response{
		Status: StatusErrorNotFound,
		Error:  msg,
	}
}

// ValidationError returns an errors
func ValidationError(errs validator.ValidationErrors) Response {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsgs, ", "),
	}
}
