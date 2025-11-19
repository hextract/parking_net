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

func (h *Handler) RegisterHandler(api operations.PostRegisterParams) middleware.Responder {
	var responder middleware.Responder
	defer utils.CatchPanic(&responder)

	ctx, span := h.tracer.Start(context.Background(), "register")
	defer span.End()

	traceID := fmt.Sprintf("%s", span.SpanContext().TraceID())

	if api.Body.Login == nil || api.Body.Email == nil || api.Body.Password == nil ||
		api.Body.Role == nil || api.Body.TelegramID == nil {
		errCode := int64(operations.PostRegisterConflictCode)
		slog.Error(
			"failed register new user",
			slog.String("method", "POST"),
			slog.String("trace_id", traceID),
			slog.Int("status_code", operations.PostRegisterConflictCode),
			slog.String("error", "missing required fields"),
		)
		responder = new(operations.PostRegisterConflict).WithPayload(&models.Error{
			ErrorMessage:    "Invalid request: missing required fields",
			ErrorStatusCode: &errCode,
		})
		return responder
	}

	token, err := impl.CreateUser(ctx, h.Client, api.Body)
	if err != nil {
		errorMsg := "Failed to register user"
		if strings.Contains(err.Error(), "already exists") || strings.Contains(err.Error(), "duplicate") {
			errorMsg = "User already exists"
		}

		slog.Error(
			"failed register new user",
			slog.String("method", "POST"),
			slog.String("trace_id", traceID),
			slog.Group("user-properties",
				slog.String("login", *api.Body.Login),
				slog.String("email", *api.Body.Email),
				slog.Int("telegram-id", int(*api.Body.TelegramID)),
			),
			slog.Int("status_code", operations.PostRegisterConflictCode),
			slog.String("error", err.Error()),
		)

		conflict := int64(operations.PostRegisterConflictCode)
		responder = new(operations.PostRegisterConflict).WithPayload(&models.Error{
			ErrorMessage:    errorMsg,
			ErrorStatusCode: &conflict,
		})
		return responder
	}

	slog.Info(
		"register new user",
		slog.String("method", "POST"),
		slog.String("trace_id", traceID),
		slog.Group("user-properties",
			slog.String("login", *api.Body.Login),
			slog.String("email", *api.Body.Email),
			slog.Int("telegram-id", int(*api.Body.TelegramID)),
		),
		slog.Int("status_code", operations.PostRegisterOKCode),
	)

	responder = new(operations.PostRegisterOK).WithPayload(&operations.PostRegisterOKBody{
		Token: *token,
	})
	return responder
}
