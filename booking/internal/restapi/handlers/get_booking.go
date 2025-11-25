package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/h4x4d/parking_net/booking/internal/grpc/client"
	"github.com/h4x4d/parking_net/booking/internal/models"
	"github.com/h4x4d/parking_net/booking/internal/restapi/operations/driver"
	"github.com/h4x4d/parking_net/booking/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (handler *Handler) GetBooking(params driver.GetBookingParams, user *models.User) (responder middleware.Responder) {
	defer utils.CatchPanic(&responder)

	ctx, span := handler.tracer.Start(context.Background(), "get booking")
	defer span.End()
	traceId := fmt.Sprintf("%s", span.SpanContext().TraceID())
	ctx = metadata.AppendToOutgoingContext(ctx, "x-trace-id", traceId)

	if user != nil && user.Role == "admin" {
		bookings, errGet := handler.Database.GetAll(params.ParkingPlaceID, params.UserID)
		if errGet != nil {
			return utils.HandleInternalError(errGet)
		}

		slog.Info(
			"get bookings",
			slog.String("method", "GET"),
			slog.String("trace_id", traceId),
			slog.Group("user-properties",
				slog.String("user-id", user.UserID),
				slog.String("role", user.Role),
				slog.Int("telegram-id", user.TelegramID),
			),
			slog.Int("status_code", driver.GetBookingOKCode),
		)

		result := new(driver.GetBookingOK)
		result.SetPayload(bookings)
		return result
	}

	if user != nil && user.Role == "driver" {
		var userID *string
		if params.UserID != nil {
			if *params.UserID != user.UserID {
				errCode := int64(driver.GetBookingForbiddenCode)
				return &driver.GetBookingForbidden{
					Payload: &models.Error{
						ErrorMessage:    "You can only view your own bookings",
						ErrorStatusCode: &errCode,
					},
				}
			}
			userID = params.UserID
		} else {
			userID = &user.UserID
		}
		bookings, errGet := handler.Database.GetAll(params.ParkingPlaceID, userID)
		if errGet != nil {
			return utils.HandleInternalError(errGet)
		}

		slog.Info(
			"get bookings",
			slog.String("method", "GET"),
			slog.String("trace_id", traceId),
			slog.Group("user-properties",
				slog.String("user-id", user.UserID),
				slog.String("role", user.Role),
				slog.Int("telegram-id", user.TelegramID),
			),
			slog.Int("status_code", driver.GetBookingOKCode),
		)

		result := new(driver.GetBookingOK)
		result.SetPayload(bookings)
		return result
	}

	if user != nil && user.Role == "owner" {
		if params.ParkingPlaceID == nil {
			errCode := int64(driver.GetBookingForbiddenCode)
			return &driver.GetBookingForbidden{
				Payload: &models.Error{
					ErrorMessage:    "parking_place_id is required for owners",
					ErrorStatusCode: &errCode,
				},
			}
		}
		parkingPlace, parkingErr := client.GetParkingPlaceById(ctx, params.ParkingPlaceID)
		if parkingErr != nil {
			if statusCode, ok := status.FromError(parkingErr); ok && statusCode.Code() == codes.NotFound {
				slog.Error(
					"failed get bookings",
					slog.String("method", "GET"),
					slog.String("trace_id", traceId),
					slog.Group("user-properties",
						slog.String("user-id", user.UserID),
						slog.String("role", user.Role),
						slog.Int("telegram-id", user.TelegramID),
					),
					slog.Group("booking-properties",
						slog.Int64("parking-place-id", *params.ParkingPlaceID),
					),
					slog.Int("status_code", http.StatusNotFound),
					slog.String("error", "Not found"),
				)

				code := int64(http.StatusNotFound)
				return &driver.GetBookingNotFound{
					Payload: &models.Error{
						ErrorStatusCode: &code,
						ErrorMessage:    fmt.Sprintf("Parking place with id %d not found", *params.ParkingPlaceID),
					},
				}
			}
			return utils.HandleInternalError(parkingErr)
		}
		if parkingPlace.OwnerID == user.UserID || user.Role == "admin" {
			bookings, errGet := handler.Database.GetAll(params.ParkingPlaceID, nil)
			if errGet != nil {
				return utils.HandleInternalError(errGet)
			}

			slog.Info(
				"get bookings",
				slog.String("method", "GET"),
				slog.String("trace_id", traceId),
				slog.Group("user-properties",
					slog.String("user-id", user.UserID),
					slog.String("role", user.Role),
					slog.Int("telegram-id", user.TelegramID),
				),
				slog.Group("booking-properties",
					slog.Int64("parking-place-id", *params.ParkingPlaceID),
				),
				slog.Int("status_code", driver.GetBookingOKCode),
			)

			result := new(driver.GetBookingOK)
			result.SetPayload(bookings)
			return result
		}
	}
	if user == nil {
		user = &models.User{
			UserID:     "empty",
			Role:       "empty",
			TelegramID: 0,
		}
	}
	slog.Error(
		"failed get bookings",
		slog.String("method", "GET"),
		slog.String("trace_id", traceId),
		slog.Group("user-properties",
			slog.String("user-id", user.UserID),
			slog.String("role", user.Role),
			slog.Int("telegram-id", user.TelegramID),
		),
		slog.Int("status_code", driver.GetBookingForbiddenCode),
		slog.String("error", "Not enough rights"),
	)

	errCode := int64(driver.GetBookingForbiddenCode)
	result := new(driver.GetBookingForbidden)
	result.SetPayload(&models.Error{
		ErrorMessage:    "You don't have permission to get this bookings",
		ErrorStatusCode: &errCode,
	})
	return result
}
