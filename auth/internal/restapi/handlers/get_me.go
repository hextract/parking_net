// NOTE: This file requires swagger code regeneration first!
// Run: swagger generate server --target ./internal --name ParkingsAuth --spec ./api/swagger/auth.yaml --principal interface{} --exclude-main
//
// After regeneration, uncomment and use this handler

package handlers

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-openapi/runtime/middleware"
	"github.com/h4x4d/parking_net/auth/internal/impl"
	"github.com/h4x4d/parking_net/auth/internal/models"
	"github.com/h4x4d/parking_net/auth/internal/restapi/operations"
	"github.com/h4x4d/parking_net/auth/internal/utils"
)

func (h *Handler) GetMeHandler(api operations.GetAuthMeParams) middleware.Responder {
	var responder middleware.Responder
	defer utils.CatchPanic(&responder)

	ctx, span := h.tracer.Start(context.Background(), "get_me")
	defer span.End()

	traceID := fmt.Sprintf("%s", span.SpanContext().TraceID())

	userInfo, err := impl.GetUserInfo(ctx, h.Client, api)
	if err != nil {
		errCode := int64(operations.GetAuthMeUnauthorizedCode)
		slog.Error(
			"failed to get user info",
			slog.String("method", "GET"),
			slog.String("trace_id", traceID),
			slog.Int("status_code", operations.GetAuthMeUnauthorizedCode),
			slog.String("error", err.Error()),
		)
		responder = new(operations.GetAuthMeUnauthorized).WithPayload(&models.Error{
			ErrorMessage:    "Unauthorized",
			ErrorStatusCode: &errCode,
		})
		return responder
	}

	slog.Info(
		"user info retrieved",
		slog.String("method", "GET"),
		slog.String("trace_id", traceID),
		slog.Int("status_code", operations.GetAuthMeOKCode),
	)

	responder = new(operations.GetAuthMeOK).WithPayload(userInfo)
	return responder
}
