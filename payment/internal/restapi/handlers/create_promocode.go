package handlers

import (
	"log/slog"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/h4x4d/parking_net/payment/internal/models"
	"github.com/h4x4d/parking_net/payment/internal/restapi/operations/admin"
	"github.com/h4x4d/parking_net/payment/internal/utils"
)

func (handler *Handler) CreatePromocode(params admin.CreatePromocodeParams, user *models.User) middleware.Responder {
	defer utils.CatchPanic(nil)

	if user.Role != "admin" {
		errCode := int64(http.StatusForbidden)
		return &admin.CreatePromocodeForbidden{
			Payload: &models.Error{
				ErrorMessage:    "admin access required",
				ErrorStatusCode: &errCode,
			},
		}
	}

	if params.Object == nil {
		errCode := int64(http.StatusBadRequest)
		return &admin.CreatePromocodeBadRequest{
			Payload: &models.Error{
				ErrorMessage:    "request body is required",
				ErrorStatusCode: &errCode,
			},
		}
	}

	if params.Object.Amount == nil || params.Object.MaxUses == nil {
		errCode := int64(http.StatusBadRequest)
		return &admin.CreatePromocodeBadRequest{
			Payload: &models.Error{
				ErrorMessage:    "amount and max_uses are required",
				ErrorStatusCode: &errCode,
			},
		}
	}

	if *params.Object.Amount <= 0 {
		errCode := int64(http.StatusBadRequest)
		return &admin.CreatePromocodeBadRequest{
			Payload: &models.Error{
				ErrorMessage:    "amount must be positive",
				ErrorStatusCode: &errCode,
			},
		}
	}

	if *params.Object.MaxUses <= 0 {
		errCode := int64(http.StatusBadRequest)
		return &admin.CreatePromocodeBadRequest{
			Payload: &models.Error{
				ErrorMessage:    "max_uses must be positive",
				ErrorStatusCode: &errCode,
			},
		}
	}

	var codePtr *string
	if params.Object.Code != "" {
		codePtr = &params.Object.Code
	}

	var expiresAtPtr *strfmt.DateTime
	if !params.Object.ExpiresAt.IsZero() {
		expiresAtPtr = &params.Object.ExpiresAt
	}

	result, err := handler.Database.CreatePromocode(
		params.HTTPRequest.Context(),
		user.UserID,
		*params.Object.Amount,
		*params.Object.MaxUses,
		codePtr,
		expiresAtPtr,
	)
	if err != nil {
		slog.Error("failed to create promocode", "error", err, "admin_id", user.UserID)
		errCode := int64(http.StatusBadRequest)
		return &admin.CreatePromocodeBadRequest{
			Payload: &models.Error{
				ErrorMessage:    err.Error(),
				ErrorStatusCode: &errCode,
			},
		}
	}

	return &admin.CreatePromocodeOK{
		Payload: result,
	}
}
