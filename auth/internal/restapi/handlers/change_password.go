package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/h4x4d/parking_net/auth/internal/impl"
	"github.com/h4x4d/parking_net/auth/internal/models"
	"github.com/h4x4d/parking_net/auth/internal/restapi/operations"
	"github.com/h4x4d/parking_net/auth/internal/utils"
)

func (h *Handler) ChangePasswordHandler(api operations.PostChangePasswordParams) middleware.Responder {
	var responder middleware.Responder
	defer utils.CatchPanic(&responder)

	ctx, span := h.tracer.Start(context.Background(), "change_password")
	defer span.End()

	traceID := fmt.Sprintf("%s", span.SpanContext().TraceID())

	if api.Body.Login == nil || api.Body.OldPassword == nil || api.Body.NewPassword == nil {
		errCode := int64(operations.PostChangePasswordBadRequestCode)
		slog.Error(
			"failed to change password",
			slog.String("method", "POST"),
			slog.String("trace_id", traceID),
			slog.Int("status_code", operations.PostChangePasswordBadRequestCode),
			slog.String("error", "missing required fields"),
		)
		responder = new(operations.PostChangePasswordBadRequest).WithPayload(&models.Error{
			ErrorMessage:    "Invalid request: missing required fields",
			ErrorStatusCode: &errCode,
		})
		return responder
	}

	token, err := impl.ChangePasswordUser(ctx, h.Client, api.Body)
	if err != nil {
		errorMsg := "Failed to change password"
		statusCode := operations.PostChangePasswordBadRequestCode
		errCode := int64(operations.PostChangePasswordBadRequestCode)

		if strings.Contains(err.Error(), "invalid") || strings.Contains(err.Error(), "unauthorized") ||
			strings.Contains(err.Error(), "401") || strings.Contains(err.Error(), "login") {
			errorMsg = "Invalid old password"
			statusCode = operations.PostChangePasswordUnauthorizedCode
			errCode = int64(operations.PostChangePasswordUnauthorizedCode)
		}

		slog.Error(
			"failed to change password",
			slog.String("method", "POST"),
			slog.String("trace_id", traceID),
			slog.Group("user-properties",
				slog.String("login", *api.Body.Login),
			),
			slog.Int("status_code", statusCode),
			slog.String("error", err.Error()),
		)

		if statusCode == operations.PostChangePasswordUnauthorizedCode {
			responder = new(operations.PostChangePasswordUnauthorized).WithPayload(&models.Error{
				ErrorMessage:    errorMsg,
				ErrorStatusCode: &errCode,
			})
			return responder
		}

		responder = new(operations.PostChangePasswordBadRequest).WithPayload(&models.Error{
			ErrorMessage:    errorMsg,
			ErrorStatusCode: &errCode,
		})
		return responder
	}

	slog.Info(
		"user changed password",
		slog.String("method", "POST"),
		slog.String("trace_id", traceID),
		slog.Group("user-properties",
			slog.String("login", *api.Body.Login),
		),
		slog.Int("status_code", operations.PostChangePasswordOKCode),
	)

	responder = new(operations.PostChangePasswordOK).WithPayload(&operations.PostChangePasswordOKBody{
		Token: *token,
	})
	return responder
}
