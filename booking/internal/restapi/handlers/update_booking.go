package handlers

import (
	"context"
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"github.com/h4x4d/parking_net/booking/internal/grpc/client"
	"github.com/h4x4d/parking_net/booking/internal/models"
	"github.com/h4x4d/parking_net/booking/internal/restapi/operations/driver"
	"github.com/h4x4d/parking_net/booking/internal/utils"
	pkg_models "github.com/h4x4d/parking_net/pkg/models"
	"google.golang.org/grpc/metadata"
	"log/slog"
)

func (handler *Handler) UpdateBooking(params driver.UpdateBookingParams, user *models.User) (responder middleware.Responder) {
	defer utils.CatchPanic(&responder)

	ctx, span := handler.tracer.Start(context.Background(), "update booking")
	defer span.End()
	traceId := fmt.Sprintf("%s", span.SpanContext().TraceID())
	ctx = metadata.AppendToOutgoingContext(ctx, "x-trace-id", traceId)

	if params.Object == nil {
		errCode := int64(driver.UpdateBookingBadRequestCode)
		userID := "unknown"
		role := "unknown"
		telegramID := 0
		if user != nil {
			userID = user.UserID
			role = user.Role
			telegramID = user.TelegramID
		}
		slog.Error(
			"failed update booking",
			slog.String("method", "PUT"),
			slog.String("trace_id", traceId),
			slog.Group("user-properties",
				slog.String("user-id", userID),
				slog.String("role", role),
				slog.Int("telegram-id", telegramID),
			),
			slog.Group("booking-properties",
				slog.Int64("booking-id", params.BookingID),
			),
			slog.Int("status_code", driver.UpdateBookingBadRequestCode),
			slog.String("error", "missing request body"),
		)
		return &driver.UpdateBookingBadRequest{
			Payload: &models.Error{
				ErrorMessage:    "Invalid request: missing required fields",
				ErrorStatusCode: &errCode,
			},
		}
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
		parkingPlaceID := int64(0)
		dateFrom := "unknown"
		dateTo := "unknown"
		if params.Object.ParkingPlaceID != nil {
			parkingPlaceID = *params.Object.ParkingPlaceID
		}
		if params.Object.DateFrom != nil {
			dateFrom = *params.Object.DateFrom
		}
		if params.Object.DateTo != nil {
			dateTo = *params.Object.DateTo
		}
		slog.Error(
			"failed update booking",
			slog.String("method", "PUT"),
			slog.String("trace_id", traceId),
			slog.Group("user-properties",
				slog.String("user-id", userID),
				slog.String("role", role),
				slog.Int("telegram-id", telegramID),
			),
			slog.Group("booking-properties",
				slog.Int64("booking-id", params.BookingID),
				slog.Int64("parking-place-id", parkingPlaceID),
				slog.String("date-from", dateFrom),
				slog.String("date-to", dateTo),
				slog.String("status", params.Object.Status),
				slog.Int64("full-cost", params.Object.FullCost),
			),
			slog.Int("status_code", driver.UpdateBookingForbiddenCode),
			slog.String("error", "Not enough rights"),
		)

		errCode := int64(driver.UpdateBookingForbiddenCode)
		result := new(driver.UpdateBookingForbidden)
		result.SetPayload(&models.Error{
			ErrorMessage:    "You don't have permission to update this booking",
			ErrorStatusCode: &errCode,
		})
		return result
	}
	booking, errUpdate := handler.Database.Update(ctx, params.BookingID, params.Object)
	if errUpdate != nil {
		return utils.HandleInternalError(errUpdate)
	}

	if handler.KafkaConn != nil {
		notifyErr := handler.KafkaConn.SendNotification(
			pkg_models.Notification{
				Name: "Booking update",
				Text: fmt.Sprintf("Your booking with booking_id %d was updated successfully",
					params.BookingID),
				TelegramID: user.TelegramID,
			})
		if notifyErr != nil {
			slog.Warn("failed to send notification", "error", notifyErr)
		}
	}
	var tgId int
	if handler.KeyCloak != nil {
		parkingPlace, parkingErr := client.GetParkingPlaceById(ctx, booking.ParkingPlaceID)
		if parkingErr != nil {
			slog.Warn("failed to get parking place for owner notification", "error", parkingErr)
		} else {
			var tgErr error
			tgId, tgErr = handler.KeyCloak.GetTelegramId(ctx, parkingPlace.OwnerID)
			if tgErr != nil {
				slog.Warn("failed to get telegram ID for owner, skipping owner notification", "error", tgErr)
				tgId = 0
			}
		}
	} else {
		slog.Warn("Keycloak client not available, skipping owner notification")
		tgId = 0
	}

	if handler.KafkaConn != nil && tgId > 0 {
		notifyErr2 := handler.KafkaConn.SendNotification(
			pkg_models.Notification{
				Name: "Booking update",
				Text: fmt.Sprintf("Your parking place %d booking with booking_id %d was updated",
					*params.Object.ParkingPlaceID, params.BookingID),
				TelegramID: tgId,
			})
		if notifyErr2 != nil {
			slog.Warn("failed to send notification to owner", "error", notifyErr2)
		}
	}

	userID := "unknown"
	role := "unknown"
	telegramID := 0
	if user != nil {
		userID = user.UserID
		role = user.Role
		telegramID = user.TelegramID
	}
	parkingPlaceID := int64(0)
	dateFrom := "unknown"
	dateTo := "unknown"
	if params.Object.ParkingPlaceID != nil {
		parkingPlaceID = *params.Object.ParkingPlaceID
	}
	if params.Object.DateFrom != nil {
		dateFrom = *params.Object.DateFrom
	}
	if params.Object.DateTo != nil {
		dateTo = *params.Object.DateTo
	}
	slog.Info(
		"update booking",
		slog.String("method", "PUT"),
		slog.String("trace_id", traceId),
		slog.Group("user-properties",
			slog.String("user-id", userID),
			slog.String("role", role),
			slog.Int("telegram-id", telegramID),
		),
		slog.Group("booking-properties",
			slog.Int64("booking-id", params.BookingID),
			slog.Int64("parking-place-id", parkingPlaceID),
			slog.String("date-from", dateFrom),
			slog.String("date-to", dateTo),
			slog.String("status", params.Object.Status),
			slog.Int64("full-cost", params.Object.FullCost),
		),
		slog.Int("status_code", driver.UpdateBookingOKCode),
	)

	result := new(driver.UpdateBookingOK)
	result.SetPayload(booking)
	return result
}
