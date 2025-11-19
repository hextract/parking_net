package handlers

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-openapi/runtime/middleware"
	"github.com/h4x4d/parking_net/parking/internal/models"
	"github.com/h4x4d/parking_net/parking/internal/repository"
	"github.com/h4x4d/parking_net/parking/internal/restapi/operations/parking"
	"github.com/h4x4d/parking_net/parking/internal/service"
	"github.com/h4x4d/parking_net/parking/internal/utils"
	"github.com/h4x4d/parking_net/pkg/domain"
	"github.com/h4x4d/parking_net/pkg/errors"
	"github.com/h4x4d/parking_net/pkg/jaeger"
	"go.opentelemetry.io/otel/trace"
)

type ParkingHandler struct {
	service *service.ParkingService
	tracer  trace.Tracer
}

func NewParkingHandler(svc *service.ParkingService) (*ParkingHandler, error) {
	tracer, err := jaeger.InitTracer("Parking")
	if err != nil {
		return nil, err
	}

	return &ParkingHandler{
		service: svc,
		tracer:  tracer,
	}, nil
}

func (h *ParkingHandler) CreateParking(params parking.CreateParkingParams, principal *models.User) middleware.Responder {
	var responder middleware.Responder
	defer utils.CatchPanic(&responder)

	ctx, span := h.tracer.Start(context.Background(), "create_parking")
	defer span.End()
	traceID := fmt.Sprintf("%s", span.SpanContext().TraceID())

	if principal == nil {
		errCode := int64(403)
		slog.Error("failed to create parking",
			slog.String("trace_id", traceID),
			slog.Int("status_code", 403),
			slog.String("error", "user not authenticated"),
		)
		responder = parking.NewCreateParkingForbidden().WithPayload(&models.Error{
			ErrorMessage:    "User not authenticated",
			ErrorStatusCode: &errCode,
		})
		return responder
	}

	if params.Object == nil {
		errCode := int64(400)
		userID := "unknown"
		if principal != nil {
			userID = principal.UserID
		}
		slog.Error("failed to create parking",
			slog.String("trace_id", traceID),
			slog.String("user_id", userID),
			slog.Int("status_code", 400),
			slog.String("error", "missing request body"),
		)
		responder = parking.NewCreateParkingBadRequest().WithPayload(&models.Error{
			ErrorMessage:    "Invalid request: missing required fields",
			ErrorStatusCode: &errCode,
		})
		return responder
	}

	domainParking := ToDomainParking(params.Object)
	domainUser := ToDomainUser(principal)

	created, appErr := h.service.CreateParking(ctx, domainParking, domainUser)
	if appErr != nil {
		responder = h.handleError(appErr, "failed to create parking", traceID, domainUser.ID)
		return responder
	}

	slog.Info("parking created",
		slog.String("trace_id", traceID),
		slog.Int64("parking_id", created.ID),
		slog.String("user_id", domainUser.ID),
	)

	responder = parking.NewCreateParkingOK().WithPayload(&parking.CreateParkingOKBody{
		ID: created.ID,
	})
	return responder
}

func (h *ParkingHandler) GetParkingByID(params parking.GetParkingByIDParams) middleware.Responder {
	var responder middleware.Responder
	defer utils.CatchPanic(&responder)

	ctx, span := h.tracer.Start(context.Background(), "get_parking_by_id")
	defer span.End()
	traceID := fmt.Sprintf("%s", span.SpanContext().TraceID())

	id := params.ParkingID

	p, appErr := h.service.GetParkingByID(ctx, id)
	if appErr != nil {
		if appErr.Code == 404 {
			statusCode := int64(404)
			slog.Error("failed to get parking by id",
				slog.String("trace_id", traceID),
				slog.Int64("parking_id", id),
				slog.Int("status_code", 404),
				slog.String("error", appErr.Message),
			)
			errorModel := &models.Error{
				ErrorMessage:    appErr.Message,
				ErrorStatusCode: &statusCode,
			}
			responder = parking.NewGetParkingByIDNotFound().WithPayload(errorModel)
			return responder
		}
		responder = h.handleError(appErr, "failed to get parking", traceID, "")
		return responder
	}

	slog.Info("get parking by id",
		slog.String("trace_id", traceID),
		slog.Int64("parking_id", id),
		slog.Int("status_code", 200),
	)

	responder = parking.NewGetParkingByIDOK().WithPayload(ToAPIParking(p))
	return responder
}

func (h *ParkingHandler) GetParkings(params parking.GetParkingsParams) middleware.Responder {
	var responder middleware.Responder
	defer utils.CatchPanic(&responder)

	ctx, span := h.tracer.Start(context.Background(), "get_parkings")
	defer span.End()
	traceID := fmt.Sprintf("%s", span.SpanContext().TraceID())

	filters := h.buildFilters(params)

	parkings, appErr := h.service.GetParkings(ctx, filters)
	if appErr != nil {
		responder = h.handleError(appErr, "failed to get parkings", traceID, "")
		return responder
	}

	slog.Info("get parkings",
		slog.String("trace_id", traceID),
		slog.Int("count", len(parkings)),
		slog.Int("status_code", 200),
	)

	responder = parking.NewGetParkingsOK().WithPayload(ToAPIParkingList(parkings))
	return responder
}

func (h *ParkingHandler) UpdateParking(params parking.UpdateParkingParams, principal *models.User) middleware.Responder {
	var responder middleware.Responder
	defer utils.CatchPanic(&responder)

	ctx, span := h.tracer.Start(context.Background(), "update_parking")
	defer span.End()
	traceID := fmt.Sprintf("%s", span.SpanContext().TraceID())

	if principal == nil {
		errCode := int64(403)
		slog.Error("failed to update parking",
			slog.String("trace_id", traceID),
			slog.Int64("parking_id", params.ParkingID),
			slog.Int("status_code", 403),
			slog.String("error", "user not authenticated"),
		)
		responder = parking.NewUpdateParkingForbidden().WithPayload(&models.Error{
			ErrorMessage:    "User not authenticated",
			ErrorStatusCode: &errCode,
		})
		return responder
	}

	if params.Object == nil {
		errCode := int64(400)
		slog.Error("failed to update parking",
			slog.String("trace_id", traceID),
			slog.Int64("parking_id", params.ParkingID),
			slog.String("user_id", principal.UserID),
			slog.Int("status_code", 400),
			slog.String("error", "missing request body"),
		)
		responder = parking.NewUpdateParkingBadRequest().WithPayload(&models.Error{
			ErrorMessage:    "Invalid request: missing required fields",
			ErrorStatusCode: &errCode,
		})
		return responder
	}

	id := params.ParkingID
	domainParking := ToDomainParkingUpdate(params.Object)
	domainUser := ToDomainUser(principal)

	appErr := h.service.UpdateParking(ctx, id, domainParking, domainUser)
	if appErr != nil {
		responder = h.handleUpdateError(appErr, "failed to update parking", traceID, domainUser.ID)
		return responder
	}

	slog.Info("parking updated",
		slog.String("trace_id", traceID),
		slog.Int64("parking_id", id),
		slog.String("user_id", domainUser.ID),
	)

	updated, appErr := h.service.GetParkingByID(ctx, id)
	if appErr != nil {
		responder = h.handleUpdateError(appErr, "failed to get updated parking", traceID, domainUser.ID)
		return responder
	}

	responder = parking.NewUpdateParkingOK().WithPayload(ToAPIParking(updated))
	return responder
}

func (h *ParkingHandler) DeleteParking(params parking.DeleteParkingParams, principal *models.User) middleware.Responder {
	var responder middleware.Responder
	defer utils.CatchPanic(&responder)

	ctx, span := h.tracer.Start(context.Background(), "delete_parking")
	defer span.End()
	traceID := fmt.Sprintf("%s", span.SpanContext().TraceID())

	if principal == nil {
		errCode := int64(403)
		slog.Error("failed to delete parking",
			slog.String("trace_id", traceID),
			slog.Int64("parking_id", params.ParkingID),
			slog.Int("status_code", 403),
			slog.String("error", "user not authenticated"),
		)
		responder = parking.NewDeleteParkingForbidden().WithPayload(&models.Error{
			ErrorMessage:    "User not authenticated",
			ErrorStatusCode: &errCode,
		})
		return responder
	}

	id := params.ParkingID
	domainUser := ToDomainUser(principal)

	appErr := h.service.DeleteParking(ctx, id, domainUser)
	if appErr != nil {
		responder = h.handleDeleteError(appErr, "failed to delete parking", traceID, domainUser.ID)
		return responder
	}

	slog.Info("parking deleted",
		slog.String("trace_id", traceID),
		slog.Int64("parking_id", id),
		slog.String("user_id", domainUser.ID),
	)

	responder = parking.NewDeleteParkingOK().WithPayload(&models.Result{
		Status:  "success",
		Message: fmt.Sprintf("Parking place %d deleted successfully", id),
	})
	return responder
}

func (h *ParkingHandler) buildFilters(params parking.GetParkingsParams) repository.ParkingFilters {
	filters := repository.ParkingFilters{}

	if params.City != nil {
		filters.City = params.City
	}
	if params.Name != nil {
		filters.Name = params.Name
	}
	if params.ParkingType != nil {
		pt := domain.ParkingType(*params.ParkingType)
		filters.ParkingType = &pt
	}
	if params.OwnerID != nil {
		filters.OwnerID = params.OwnerID
	}

	return filters
}

func (h *ParkingHandler) handleError(appErr *errors.AppError, context string, traceID string, userID string) middleware.Responder {
	slog.Error(context,
		slog.String("trace_id", traceID),
		slog.String("user_id", userID),
		slog.String("error", appErr.Error()),
		slog.Int("status_code", appErr.Code),
	)

	statusCode := int64(appErr.Code)
	errorModel := &models.Error{
		ErrorMessage:    appErr.Message,
		ErrorStatusCode: &statusCode,
	}

	switch appErr.Code {
	case 404:
		return parking.NewGetParkingByIDNotFound().WithPayload(errorModel)
	case 400:
		return parking.NewCreateParkingBadRequest().WithPayload(errorModel)
	case 403:
		return parking.NewCreateParkingForbidden().WithPayload(errorModel)
	default:
		statusCode = 500
		errorModel.ErrorStatusCode = &statusCode
		return parking.NewGetParkingByIDNotFound().WithPayload(errorModel)
	}
}

func (h *ParkingHandler) handleUpdateError(appErr *errors.AppError, context string, traceID string, userID string) middleware.Responder {
	slog.Error(context,
		slog.String("trace_id", traceID),
		slog.String("user_id", userID),
		slog.String("error", appErr.Error()),
		slog.Int("status_code", appErr.Code),
	)

	statusCode := int64(appErr.Code)
	errorModel := &models.Error{
		ErrorMessage:    appErr.Message,
		ErrorStatusCode: &statusCode,
	}

	switch appErr.Code {
	case 404:
		return parking.NewUpdateParkingNotFound().WithPayload(errorModel)
	case 400:
		return parking.NewUpdateParkingBadRequest().WithPayload(errorModel)
	case 403:
		return parking.NewUpdateParkingForbidden().WithPayload(errorModel)
	default:
		return parking.NewUpdateParkingBadRequest().WithPayload(errorModel)
	}
}

func (h *ParkingHandler) handleDeleteError(appErr *errors.AppError, context string, traceID string, userID string) middleware.Responder {
	slog.Error(context,
		slog.String("trace_id", traceID),
		slog.String("user_id", userID),
		slog.String("error", appErr.Error()),
		slog.Int("status_code", appErr.Code),
	)

	statusCode := int64(appErr.Code)
	errorModel := &models.Error{
		ErrorMessage:    appErr.Message,
		ErrorStatusCode: &statusCode,
	}

	switch appErr.Code {
	case 404:
		return parking.NewDeleteParkingNotFound().WithPayload(errorModel)
	case 403:
		return parking.NewDeleteParkingForbidden().WithPayload(errorModel)
	default:
		return parking.NewDeleteParkingForbidden().WithPayload(errorModel)
	}
}
