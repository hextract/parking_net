package handlers

import (
	"context"
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"github.com/h4x4d/parking_net/booking/internal/models"
	"github.com/h4x4d/parking_net/booking/internal/restapi/operations/driver"
	"github.com/h4x4d/parking_net/booking/internal/utils"
	"google.golang.org/grpc/metadata"
	"log/slog"
)

func (handler *Handler) GetBookingByID(params driver.GetBookingByIDParams, user *models.User) (responder middleware.Responder) {
	defer utils.CatchPanic(&responder)

	ctx, span := handler.tracer.Start(context.Background(), "get booking by id")
	defer span.End()
	traceId := fmt.Sprintf("%s", span.SpanContext().TraceID())
	ctx = metadata.AppendToOutgoingContext(ctx, "x-trace-id", traceId)

	booking, errGet := handler.Database.GetByID(params.BookingID)
	if errGet != nil {
		return utils.HandleInternalError(errGet)
	}
	if booking == nil {
		userID := "unknown"
		role := "unknown"
		telegramID := 0
		if user != nil {
			userID = user.UserID
			role = user.Role
			telegramID = user.TelegramID
		}
		slog.Error(
			"failed get booking by id",
			slog.String("method", "GET"),
			slog.String("trace_id", traceId),
			slog.Group("user-properties",
				slog.String("user-id", userID),
				slog.String("role", role),
				slog.Int("telegram-id", telegramID),
			),
			slog.Group("booking-properties",
				slog.Int64("booking-id", params.BookingID),
			),
			slog.Int("status_code", driver.GetBookingByIDNotFoundCode),
			slog.String("error", "Booking not found"),
		)

		errCode := int64(driver.GetBookingByIDNotFoundCode)
		result := new(driver.GetBookingByIDNotFound)
		result.SetPayload(&models.Error{
			ErrorMessage:    fmt.Sprintf("Booking with id %d not found", params.BookingID),
			ErrorStatusCode: &errCode,
		})
		return result
	}

	isOwner, err := handler.Database.CheckOwnership(ctx, params.BookingID, user)
	if err != nil {
		return utils.HandleInternalError(err)
	}
	if !isOwner {
		userID := "unknown"
		role := "unknown"
		telegramID := 0
		if user != nil {
			userID = user.UserID
			role = user.Role
			telegramID = user.TelegramID
		}
		slog.Error(
			"failed get booking by id",
			slog.String("method", "GET"),
			slog.String("trace_id", traceId),
			slog.Group("user-properties",
				slog.String("user-id", userID),
				slog.String("role", role),
				slog.Int("telegram-id", telegramID),
			),
			slog.Group("booking-properties",
				slog.Int64("booking-id", params.BookingID),
			),
			slog.Int("status_code", driver.GetBookingByIDForbiddenCode),
			slog.String("error", "Not enough rights"),
		)

		errCode := int64(driver.GetBookingByIDForbiddenCode)
		result := new(driver.GetBookingByIDForbidden)
		result.SetPayload(&models.Error{
			ErrorMessage:    "You don't have permission to get this booking",
			ErrorStatusCode: &errCode,
		})
		return result
	}

	userID := "unknown"
	role := "unknown"
	telegramID := 0
	if user != nil {
		userID = user.UserID
		role = user.Role
		telegramID = user.TelegramID
	}
	slog.Info(
		"get booking by id",
		slog.String("method", "GET"),
		slog.String("trace_id", traceId),
		slog.Group("user-properties",
			slog.String("user-id", userID),
			slog.String("role", role),
			slog.Int("telegram-id", telegramID),
		),
		slog.Group("booking-properties",
			slog.Int64("booking-id", params.BookingID),
		),
		slog.Int("status_code", driver.GetBookingByIDOKCode),
	)
	result := new(driver.GetBookingByIDOK)
	result.SetPayload(booking)
	return result
}
