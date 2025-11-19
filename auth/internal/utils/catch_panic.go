package utils

import (
	"fmt"
	"log/slog"
	"github.com/go-openapi/runtime/middleware"
	"github.com/h4x4d/parking_net/auth/internal/models"
	"github.com/h4x4d/parking_net/auth/internal/restapi/operations"
	"net/http"
)

func CatchPanic(responder *middleware.Responder) {
	if err := recover(); err != nil {
		errText := fmt.Sprintf("%v", err)
		slog.Error("panic recovered",
			slog.String("error", errText),
		)
		errCode := int64(http.StatusInternalServerError)
		*responder = new(operations.PostLoginUnauthorized).WithPayload(&models.Error{
			ErrorMessage:    "Internal server error",
			ErrorStatusCode: &errCode,
		})
	}
}

