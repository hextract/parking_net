package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/go-openapi/runtime/middleware"
	payment_client "github.com/h4x4d/parking_net/booking/internal/grpc/client"
	"github.com/h4x4d/parking_net/booking/internal/models"
	"github.com/h4x4d/parking_net/booking/internal/restapi/operations/driver"
	"github.com/h4x4d/parking_net/booking/internal/utils"
	pkg_models "github.com/h4x4d/parking_net/pkg/models"
	"google.golang.org/grpc/metadata"
	"log/slog"
	"net/http"
)

func (handler *Handler) CreateBooking(params driver.CreateBookingParams, user *models.User) (responder middleware.Responder) {
	defer utils.CatchPanic(&responder)

	ctx, span := handler.tracer.Start(context.Background(), "create booking")
	defer span.End()
	traceId := fmt.Sprintf("%s", span.SpanContext().TraceID())
	ctx = metadata.AppendToOutgoingContext(ctx, "x-trace-id", traceId)

	if user != nil && user.Role == "driver" {
		if params.Object.DateFrom == nil || params.Object.DateTo == nil ||
			params.Object.ParkingPlaceID == nil {
			errCode := int64(http.StatusBadRequest)
			slog.Error(
				"failed create new booking",
				slog.String("method", "POST"),
				slog.String("trace_id", traceId),
				slog.Group("user-properties",
					slog.String("user-id", user.UserID),
					slog.String("role", user.Role),
					slog.Int("telegram-id", user.TelegramID),
				),
				slog.Int("status_code", http.StatusBadRequest),
				slog.String("error", "missing required fields"),
			)
			return &driver.CreateBookingBadRequest{
				Payload: &models.Error{
					ErrorMessage:    "Invalid request: missing required fields",
					ErrorStatusCode: &errCode,
				},
			}
		}

		bookingId, errCreate := handler.Database.Create(ctx,
			params.Object.DateFrom,
			params.Object.DateTo,
			params.Object.ParkingPlaceID,
			user.UserID,
		)
		if errCreate != nil {
			slog.Error(
				"failed create new booking",
				slog.String("method", "POST"),
				slog.String("trace_id", traceId),
				slog.Group("user-properties",
					slog.String("user-id", user.UserID),
					slog.String("role", user.Role),
					slog.Int("telegram-id", user.TelegramID),
				),
				slog.Group("booking-properties",
					slog.String("date-from", params.Object.DateFrom.String()),
					slog.String("date-to", params.Object.DateTo.String()),
					slog.Int64("parking-place-id", *params.Object.ParkingPlaceID),
				),
				slog.Int("status_code", http.StatusInternalServerError),
				slog.String("error", "failed to create booking"),
			)

			return utils.HandleInternalError(fmt.Errorf("failed to create booking"))
		}

		booking, errGet := handler.Database.GetByID(*bookingId)
		if errGet != nil {
			slog.Error("failed to get booking after creation", "error", errGet, "booking_id", *bookingId)
			return utils.HandleInternalError(errGet)
		}

		parkingPlace, parkingErr := payment_client.GetParkingPlaceById(ctx, params.Object.ParkingPlaceID)
		if parkingErr != nil {
			return utils.HandleInternalError(parkingErr)
		}

		dateFrom := time.Time(*params.Object.DateFrom)
		now := time.Now()
		if dateFrom.Before(now) || dateFrom.Sub(now) < 5*time.Minute {
			paymentResult, paymentErr := handler.PaymentClient.ProcessTransaction(ctx, *bookingId, user.UserID, parkingPlace.OwnerID, booking.FullCost)
			if paymentErr != nil || paymentResult == nil || paymentResult.Status != "completed" {
				booking.Status = "Canceled"
				handler.Database.Update(ctx, *bookingId, booking)
				if paymentErr != nil {
					slog.Error("payment processing failed", "error", paymentErr, "booking_id", *bookingId)
				} else {
					slog.Warn("payment processing failed", "status", paymentResult.Status, "message", paymentResult.Message, "booking_id", *bookingId)
				}
				errCode := int64(http.StatusBadRequest)
				return &driver.CreateBookingBadRequest{
					Payload: &models.Error{
						ErrorMessage:    "payment processing failed: insufficient funds",
						ErrorStatusCode: &errCode,
					},
				}
			}
			booking.Status = "Confirmed"
			handler.Database.Update(ctx, *bookingId, booking)
		}

		if handler.KafkaConn != nil {
			notifyErr := handler.KafkaConn.SendNotification(
				pkg_models.Notification{
					Name: "New booking",
					Text: fmt.Sprintf("Your booking with booking_id %d was created successfully",
						*bookingId),
					TelegramID: user.TelegramID,
				})
			if notifyErr != nil {
				slog.Warn("failed to send notification", "error", notifyErr)
			}
		}
		var tgId int
		if handler.KeyCloak != nil {
			var tgErr error
			tgId, tgErr = handler.KeyCloak.GetTelegramId(ctx, parkingPlace.OwnerID)
			if tgErr != nil {
				slog.Warn("failed to get telegram ID for owner, skipping owner notification", "error", tgErr)
				tgId = 0
			}
		} else {
			slog.Warn("Keycloak client not available, skipping owner notification")
			tgId = 0
		}

		if handler.KafkaConn != nil && tgId > 0 {
			notifyErr2 := handler.KafkaConn.SendNotification(
				pkg_models.Notification{
					Name: "New Booking",
					Text: fmt.Sprintf("Your parking place %d was booked with booking_id %d",
						*params.Object.ParkingPlaceID, *bookingId),
					TelegramID: tgId,
				})
			if notifyErr2 != nil {
				slog.Warn("failed to send notification to owner", "error", notifyErr2)
			}
		}

		slog.Info(
			"create new booking",
			slog.String("method", "POST"),
			slog.String("trace_id", traceId),
			slog.Group("user-properties",
				slog.String("user-id", user.UserID),
				slog.String("role", user.Role),
				slog.Int("telegram-id", user.TelegramID),
			),
			slog.Group("booking-properties",
				slog.String("date-from", params.Object.DateFrom.String()),
				slog.String("date-to", params.Object.DateTo.String()),
				slog.Int64("parking-place-id", *params.Object.ParkingPlaceID),
				slog.Int64("booking-id", *bookingId),
			),
			slog.Int("status_code", driver.CreateBookingOKCode),
		)

		result := new(driver.CreateBookingOK)
		result.SetPayload(&driver.CreateBookingOKBody{BookingID: *bookingId})
		return result
	} else {
		if user == nil {
			user = &models.User{
				UserID:     "empty",
				Role:       "empty",
				TelegramID: 0,
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
		slog.Error(
			"failed create new booking",
			slog.String("method", "POST"),
			slog.String("trace_id", traceId),
			slog.Group("user-properties",
				slog.String("user-id", userID),
				slog.String("role", role),
				slog.Int("telegram-id", telegramID),
			),
			slog.Int("status_code", http.StatusForbidden),
			slog.String("error", "Creation of booking allowed only to drivers"),
		)

		errCode := int64(http.StatusForbidden)
		result := new(driver.CreateBookingForbidden)
		result.SetPayload(&models.Error{
			ErrorMessage:    "You don't have permission to create a booking",
			ErrorStatusCode: &errCode,
		})
		return result
	}
}
