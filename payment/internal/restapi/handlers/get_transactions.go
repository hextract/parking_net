package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/h4x4d/parking_net/payment/internal/models"
	"github.com/h4x4d/parking_net/payment/internal/restapi/operations/driver"
	"github.com/h4x4d/parking_net/payment/internal/utils"
	"log/slog"
	"net/http"
)

func (handler *Handler) GetTransactions(params driver.GetTransactionsParams, user *models.User) middleware.Responder {
	defer utils.CatchPanic(nil)

	limit := int64(50)
	if params.Limit != nil {
		limit = *params.Limit
	}

	offset := int64(0)
	if params.Offset != nil {
		offset = *params.Offset
	}

	transactions, err := handler.Database.GetTransactions(params.HTTPRequest.Context(), user.UserID, limit, offset)
	if err != nil {
		slog.Error("failed to get transactions", "error", err, "user_id", user.UserID)
		errCode := int64(http.StatusInternalServerError)
		return &driver.GetTransactionsForbidden{
			Payload: &models.Error{
				ErrorMessage:    "failed to get transactions",
				ErrorStatusCode: &errCode,
			},
		}
	}

	return &driver.GetTransactionsOK{
		Payload: transactions,
	}
}

