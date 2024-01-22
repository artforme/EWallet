package createWallet

import (
	"EWallet/internal/lib/api/response"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type Response struct {
	response.Response
	response.Wallet
}

type WalletCreator interface {
	CreateWallet() (string, error)
}

// New creates handler that creates new wallet and response
func New(log *slog.Logger, walletCreator WalletCreator) http.HandlerFunc {
	//our handler
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.createWallet.create.New"

		// add to log info
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// create walletID
		walletID, err := walletCreator.CreateWallet()
		if err != nil {
			log.Error("failed to create wallet", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			// if caught error send to client it
			render.JSON(w, r, response.ResError("failed to create wallet"))

			return
		}

		log.Info("wallet created", slog.String("ID", walletID))

		wallet := response.Wallet{WalletID: walletID, Balance: "100.0"}
		// send to client status 200 and walletID with balance
		render.JSON(w, r, Response{
			Response: response.OK(),
			Wallet:   wallet,
		})
	}

}
