package handlers

import (
	"log/slog"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/h4x4d/parking_net/payment/internal/models"
	"github.com/h4x4d/parking_net/payment/internal/restapi/operations/driver"
	"github.com/h4x4d/parking_net/payment/internal/utils"
)

func (handler *Handler) ActivatePromocode(params driver.ActivatePromocodeParams, user *models.User) middleware.Responder {
	defer utils.CatchPanic(nil)

	if params.Object == nil || params.Object.Code == nil || *params.Object.Code == "" {
		errCode := int64(http.StatusBadRequest)
		return &driver.ActivatePromocodeBadRequest{
			Payload: &models.Error{
				ErrorMessage:    "promocode is required",
				ErrorStatusCode: &errCode,
			},
		}
	}

	balance, err := handler.Database.ActivatePromocode(params.HTTPRequest.Context(), user.UserID, *params.Object.Code)
	if err != nil {
		slog.Error("failed to activate promocode", "error", err, "user_id", user.UserID, "code", *params.Object.Code)
		errCode := int64(http.StatusBadRequest)
		if err.Error() == "promocode not found" {
			errCode = int64(http.StatusNotFound)
			return &driver.ActivatePromocodeNotFound{
				Payload: &models.Error{
					ErrorMessage:    "promocode not found",
					ErrorStatusCode: &errCode,
				},
			}
		}
		return &driver.ActivatePromocodeBadRequest{
			Payload: &models.Error{
				ErrorMessage:    err.Error(),
				ErrorStatusCode: &errCode,
			},
		}
	}

	return &driver.ActivatePromocodeOK{
		Payload: balance,
	}
}

