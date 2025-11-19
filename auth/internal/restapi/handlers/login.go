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

func (h *Handler) LoginHandler(api operations.PostLoginParams) middleware.Responder {
	var responder middleware.Responder
	defer utils.CatchPanic(&responder)

	ctx, span := h.tracer.Start(context.Background(), "login")
	defer span.End()

	traceID := fmt.Sprintf("%s", span.SpanContext().TraceID())

	if api.Body.Login == nil || api.Body.Password == nil {
		errCode := int64(operations.PostLoginUnauthorizedCode)
		slog.Error(
			"failed login user",
			slog.String("method", "POST"),
			slog.String("trace_id", traceID),
			slog.Int("status_code", operations.PostLoginUnauthorizedCode),
			slog.String("error", "missing required fields"),
		)
		responder = new(operations.PostLoginUnauthorized).WithPayload(&models.Error{
			ErrorMessage:    "Invalid request: missing required fields",
			ErrorStatusCode: &errCode,
		})
		return responder
	}

	token, err := impl.LoginUser(ctx, h.Client, api.Body)
	if err != nil {
		errCode := int64(operations.PostLoginUnauthorizedCode)
		slog.Error(
			"failed login user",
			slog.String("method", "POST"),
			slog.String("trace_id", traceID),
			slog.Group("user-properties",
				slog.String("login", *api.Body.Login),
			),
			slog.Int("status_code", operations.PostLoginUnauthorizedCode),
			slog.String("error", "authentication failed"),
		)
		responder = new(operations.PostLoginUnauthorized).WithPayload(&models.Error{
			ErrorMessage:    "Invalid login or password",
			ErrorStatusCode: &errCode,
		})
		return responder
	}

	slog.Info(
		"user login",
		slog.String("method", "POST"),
		slog.String("trace_id", traceID),
		slog.Group("user-properties",
			slog.String("login", *api.Body.Login),
		),
		slog.Int("status_code", operations.PostLoginOKCode),
	)

	responder = new(operations.PostLoginOK).WithPayload(&operations.PostLoginOKBody{
		Token: *token,
	})
	return responder
}
