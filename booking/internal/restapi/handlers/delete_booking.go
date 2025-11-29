package handlers

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-openapi/runtime/middleware"
	payment_client "github.com/h4x4d/parking_net/booking/internal/grpc/client"
	"github.com/h4x4d/parking_net/booking/internal/models"
	"github.com/h4x4d/parking_net/booking/internal/restapi/operations/driver"
	"github.com/h4x4d/parking_net/booking/internal/utils"
	"google.golang.org/grpc/metadata"
)

func (handler *Handler) DeleteBooking(params driver.DeleteBookingParams, user *models.User) (responder middleware.Responder) {
	defer utils.CatchPanic(&responder)

	ctx, span := handler.tracer.Start(context.Background(), "delete booking")
	defer span.End()
	traceId := fmt.Sprintf("%s", span.SpanContext().TraceID())
	ctx = metadata.AppendToOutgoingContext(ctx, "x-trace-id", traceId)

	booking, err := handler.Database.GetByID(params.BookingID)
	if err != nil {
		return utils.HandleInternalError(err)
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
			"failed delete booking",
			slog.String("method", "DELETE"),
			slog.String("trace_id", traceId),
			slog.Group("user-properties",
				slog.String("user-id", userID),
				slog.String("role", role),
				slog.Int("telegram-id", telegramID),
			),
			slog.Group("booking-properties",
				slog.Int64("booking-id", params.BookingID),
			),
			slog.Int("status_code", driver.DeleteBookingNotFoundCode),
			slog.String("error", "Booking not found"),
		)

		errCode := int64(driver.DeleteBookingNotFoundCode)
		result := new(driver.DeleteBookingNotFound)
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
			"failed delete booking",
			slog.String("method", "DELETE"),
			slog.String("trace_id", traceId),
			slog.Group("user-properties",
				slog.String("user-id", userID),
				slog.String("role", role),
				slog.Int("telegram-id", telegramID),
			),
			slog.Group("booking-properties",
				slog.Int64("booking-id", params.BookingID),
			),
			slog.Int("status_code", driver.DeleteBookingForbiddenCode),
			slog.String("error", "Not enough rights"),
		)

		errCode := int64(driver.DeleteBookingForbiddenCode)
		result := new(driver.DeleteBookingForbidden)
		result.SetPayload(&models.Error{
			ErrorMessage:    "You don't have permission to delete this booking",
			ErrorStatusCode: &errCode,
		})
		return result
	}

	if booking.Status == "Confirmed" {
		parkingPlace, parkingErr := payment_client.GetParkingPlaceById(ctx, booking.ParkingPlaceID)
		if parkingErr == nil {
			_, refundErr := handler.PaymentClient.ProcessRefund(ctx, params.BookingID, booking.UserID, parkingPlace.OwnerID, booking.FullCost)
			if refundErr != nil {
				slog.Warn("failed to process refund for canceled booking", "error", refundErr, "booking_id", params.BookingID)
			}
		}
	}

	err = handler.Database.Delete(ctx, params.BookingID)
	if err != nil {
		return utils.HandleInternalError(err)
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
		"delete booking",
		slog.String("method", "DELETE"),
		slog.String("trace_id", traceId),
		slog.Group("user-properties",
			slog.String("user-id", userID),
			slog.String("role", role),
			slog.Int("telegram-id", telegramID),
		),
		slog.Group("booking-properties",
			slog.Int64("booking-id", params.BookingID),
		),
		slog.Int("status_code", driver.DeleteBookingOKCode),
	)

	result := new(driver.DeleteBookingOK)
	result.SetPayload(&models.Result{
		Status:  "success",
		Message: fmt.Sprintf("Booking %d deleted successfully", params.BookingID),
	})
	return result
}
