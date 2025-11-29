package handlers

import (
	"log/slog"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/h4x4d/parking_net/payment/internal/models"
	"github.com/h4x4d/parking_net/payment/internal/restapi/operations/driver"
	"github.com/h4x4d/parking_net/payment/internal/utils"
)

func (handler *Handler) GeneratePromocode(params driver.GeneratePromocodeParams, user *models.User) middleware.Responder {
	defer utils.CatchPanic(nil)

	if params.Object == nil || params.Object.Amount == nil {
		errCode := int64(http.StatusBadRequest)
		return &driver.GeneratePromocodeBadRequest{
			Payload: &models.Error{
				ErrorMessage:    "amount is required",
				ErrorStatusCode: &errCode,
			},
		}
	}

	if *params.Object.Amount <= 0 {
		errCode := int64(http.StatusBadRequest)
		return &driver.GeneratePromocodeBadRequest{
			Payload: &models.Error{
				ErrorMessage:    "amount must be positive",
				ErrorStatusCode: &errCode,
			},
		}
	}

	result, err := handler.Database.GeneratePromocode(params.HTTPRequest.Context(), user.UserID, *params.Object.Amount)
	if err != nil {
		slog.Error("failed to generate promocode", "error", err, "user_id", user.UserID, "amount", *params.Object.Amount)
		errCode := int64(http.StatusBadRequest)
		if err.Error() == "insufficient funds" {
			return &driver.GeneratePromocodeBadRequest{
				Payload: &models.Error{
					ErrorMessage:    "insufficient funds",
					ErrorStatusCode: &errCode,
				},
			}
		}
		errCode = int64(http.StatusInternalServerError)
		return &driver.GeneratePromocodeBadRequest{
			Payload: &models.Error{
				ErrorMessage:    "failed to generate promocode",
				ErrorStatusCode: &errCode,
			},
		}
	}

	return &driver.GeneratePromocodeOK{
		Payload: result,
	}
}

