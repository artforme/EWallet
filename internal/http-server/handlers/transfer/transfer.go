package transfer

import (
	"EWallet/internal/lib/api/response"
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
}

func New(log *slog.Logger, walletTransfer WalletTransfer) http.HandlerFunc {
	//our handler
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.transfer.transfer.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		FromWalletID := chi.URLParam(r, "walletId")
		if FromWalletID == "" {
			log.Info("walletID is empty")

			render.JSON(w, r, response.ResErrorNotFound("invalid request"))

			return
		}

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})

			render.JSON(w, r, response.ResError("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err = validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})

			render.JSON(w, r, response.ValidationError(validateErr))

			return
		}

		err = walletTransfer.Transfer(FromWalletID, req.WalletID, req.Amount)
		if err != nil {
			log.Error("invalid request", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			if err.Error() == "outgoing wallet not found" {
				render.JSON(w, r, response.ResErrorNotFound(err.Error()))
				return
			}
			render.JSON(w, r, response.ResError(err.Error()))
			return
		}
		log.Info("successful transfer")

		render.JSON(w, r, response.OK())

	}
}
