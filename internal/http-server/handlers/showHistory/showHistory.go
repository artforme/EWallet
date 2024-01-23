package showHistory

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
	RespTran []response.RespTransaction
}

type HistoryShower interface {
	ShowHistory(walletID string) ([]response.RespTransaction, error)
	CheckIfExists(walletID string) (bool, error)
}

// New creates a handler that shows transactions history of specific wallet
func New(log *slog.Logger, historyShower HistoryShower) http.HandlerFunc {
	//our handler
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.showHistory.showHistory.New"

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
		if exists, err := historyShower.CheckIfExists(walletID); !exists || err != nil {
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

		history, err := historyShower.ShowHistory(walletID)
		if err != nil {
			log.Error("failed to show outgoing walletID", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			render.JSON(w, r, response.ResError("failed request"))
			return
		}

		log.Info("history was found", slog.String("ID", walletID))

		// send to client status 200 and walletID with balance
		render.JSON(w, r, Response{
			RespTran: history,
		})
	}
}
