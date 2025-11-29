package handlers

import (
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"github.com/h4x4d/parking_net/payment/internal/models"
	"github.com/h4x4d/parking_net/payment/internal/restapi/operations/driver"
	"github.com/h4x4d/parking_net/payment/internal/utils"
	"log/slog"
	"net/http"
)

func (handler *Handler) GetBalance(params driver.GetBalanceParams, user *models.User) middleware.Responder {
	defer utils.CatchPanic(nil)

	balance, err := handler.Database.GetBalance(user.UserID)
	if err != nil {
		slog.Error("failed to get balance", "error", err, "user_id", user.UserID)
		errCode := int64(http.StatusInternalServerError)
		return &driver.GetBalanceNotFound{
			Payload: &models.Error{
				ErrorMessage:    fmt.Sprintf("failed to get balance: %v", err),
				ErrorStatusCode: &errCode,
			},
		}
	}

	return &driver.GetBalanceOK{
		Payload: balance,
	}
}

