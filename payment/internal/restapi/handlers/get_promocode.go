package handlers

import (
	"log/slog"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/h4x4d/parking_net/payment/internal/models"
	"github.com/h4x4d/parking_net/payment/internal/restapi/operations/driver"
	"github.com/h4x4d/parking_net/payment/internal/utils"
)

func (handler *Handler) GetPromocode(params driver.GetPromocodeParams, user *models.User) middleware.Responder {
	defer utils.CatchPanic(nil)

	promocode, err := handler.Database.GetPromocode(params.HTTPRequest.Context(), params.Code)
	if err != nil {
		slog.Error("failed to get promocode", "error", err, "code", params.Code)
		errCode := int64(http.StatusNotFound)
		return &driver.GetPromocodeNotFound{
			Payload: &models.Error{
				ErrorMessage:    "failed to get promocode",
				ErrorStatusCode: &errCode,
			},
		}
	}

	if promocode == nil {
		errCode := int64(http.StatusNotFound)
		return &driver.GetPromocodeNotFound{
			Payload: &models.Error{
				ErrorMessage:    "promocode not found",
				ErrorStatusCode: &errCode,
			},
		}
	}

	return &driver.GetPromocodeOK{
		Payload: promocode,
	}
}

