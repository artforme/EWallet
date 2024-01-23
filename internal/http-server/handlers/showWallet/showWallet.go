package showWallet

import (
	"EWallet/internal/lib/api/response"
	"database/sql"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type Response struct {
	response.Response
	response.Wallet
}

type WalletShower interface {
	ShowWallet(walletID string) (response.Wallet, error)
	CheckIfExists(walletID string) (bool, error)
}

// New creates a handler that shows wallet and his current balance
func New(log *slog.Logger, walletShower WalletShower) http.HandlerFunc {
	//our handler
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.showWallet.showWallet.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		// get walletID form url
		walletID := chi.URLParam(r, "walletId")
		if walletID == "" {
			log.Info("walletID is empty")

			render.JSON(w, r, response.ResErrorNotFound("invalid request"))

			return
		}
		// check if exists in database
		if exists, err := walletShower.CheckIfExists(walletID); !exists || err != nil {
			log.Error("failed to find outgoing walletID", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			if errors.Is(err, sql.ErrNoRows) || !exists {
				render.JSON(w, r, response.ResErrorNotFound("failed to find outgoing walletID"))
				return
			}
			render.JSON(w, r, response.ResError("failed request"))
			return
		}

		wallet, err := walletShower.ShowWallet(walletID)
		if err != nil {
			log.Error("failed to create wallet", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			// if caught error send to client it
			render.JSON(w, r, response.ResError("failed to create wallet"))

			return
		}

		log.Info("wallet has been found", slog.String("ID", walletID))

		render.JSON(w, r, Response{
			Response: response.OK(),
			Wallet:   wallet,
		})
	}

}
