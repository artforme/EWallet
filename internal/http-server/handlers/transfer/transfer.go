package transfer

import (
	"EWallet/internal/lib/api/response"
	"database/sql"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

type Request struct {
	WalletID string `json:"walletId" validate:"required"`
	Amount   string `json:"amount" validate:"required"`
}

type Response struct {
	response.Response
}

type WalletTransfer interface {
	Transfer(fromWallet, toWallet, amount string) error
	CheckIfExists(walletID string) (bool, error)
}

// New creates a handler that send amount from one wallet to another
func New(log *slog.Logger, walletTransfer WalletTransfer) http.HandlerFunc {
	//our handler
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.transfer.transfer.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		// get walletID form url
		FromWalletID := chi.URLParam(r, "walletId")
		if FromWalletID == "" {
			log.Info("walletID is empty")

			render.JSON(w, r, response.ResErrorNotFound("invalid request"))

			return
		}
		// check if exists in database
		if exists, err := walletTransfer.CheckIfExists(FromWalletID); !exists || err != nil {
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

		var req Request
		// decode body request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})

			render.JSON(w, r, response.ResError("failed to decode request"))

			return
		}
		// check if exists in database
		if exists, err := walletTransfer.CheckIfExists(req.WalletID); !exists || err != nil {
			log.Error("failed to find outgoing walletID", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			if errors.Is(err, sql.ErrNoRows) || !exists {
				render.JSON(w, r, response.ResError("failed to find target walletID"))
				return
			}
			render.JSON(w, r, response.ResError("failed request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))
		// check fields of request body
		if err = validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})

			render.JSON(w, r, response.ValidationError(validateErr))

			return
		}
		// do transfer
		err = walletTransfer.Transfer(FromWalletID, req.WalletID, req.Amount)
		if err != nil {
			log.Error("invalid request", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			render.JSON(w, r, response.ResError("failed request"))
			return
		}
		log.Info("successful transfer")

		render.JSON(w, r, response.OK())

	}
}
