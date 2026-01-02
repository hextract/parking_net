package handlers

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-openapi/runtime/middleware"
	"github.com/h4x4d/parking_net/parking/internal/models"
	"github.com/h4x4d/parking_net/parking/internal/restapi/operations/parking"
	"github.com/h4x4d/parking_net/parking/internal/service"
	"github.com/h4x4d/parking_net/parking/internal/utils"
	"github.com/h4x4d/parking_net/pkg/errors"
	"github.com/h4x4d/parking_net/pkg/jaeger"
	"go.opentelemetry.io/otel/trace"
)

type ServiceHandler struct {
	service *service.ServiceService
	tracer  trace.Tracer
}

func NewServiceHandler(svc *service.ServiceService) (*ServiceHandler, error) {
	tracer, err := jaeger.InitTracer("Parking")
	if err != nil {
		return nil, err
	}

	return &ServiceHandler{
		service: svc,
		tracer:  tracer,
	}, nil
}

func (h *ServiceHandler) CreateService(params parking.CreateParkingServiceParams, principal *models.User) middleware.Responder {
	var responder middleware.Responder
	defer utils.CatchPanic(&responder)

	ctx, span := h.tracer.Start(context.Background(), "create_service")
	defer span.End()
	traceID := fmt.Sprintf("%s", span.SpanContext().TraceID())

	if principal == nil {
		errCode := int64(403)
		slog.Error("failed to create service",
			slog.String("trace_id", traceID),
			slog.Int("status_code", 403),
			slog.String("error", "user not authenticated"),
		)
		return &parking.CreateParkingServiceForbidden{
			Payload: &models.Error{
				ErrorMessage:    "Forbidden",
				ErrorStatusCode: &errCode,
			},
		}
	}

	user := ToDomainUser(principal)
	serviceModel := ToDomainService(params.Object)
	serviceModel.ParkingID = params.ParkingID

	created, appErr := h.service.CreateService(ctx, serviceModel, user)
	if appErr != nil {
		return h.handleCreateError(appErr, "failed to create service", traceID, user.ID)
	}

	return &parking.CreateParkingServiceOK{
		Payload: ToAPIService(created),
	}
}

func (h *ServiceHandler) GetService(params parking.GetParkingServiceParams, principal *models.User) middleware.Responder {
	var responder middleware.Responder
	defer utils.CatchPanic(&responder)

	ctx, span := h.tracer.Start(context.Background(), "get_service")
	defer span.End()
	traceID := fmt.Sprintf("%s", span.SpanContext().TraceID())

	serviceModel, appErr := h.service.GetServiceByID(ctx, params.ServiceID)
	if appErr != nil {
		return h.handleGetError(appErr, "failed to get service", traceID, "")
	}

	return &parking.GetParkingServiceOK{
		Payload: ToAPIService(serviceModel),
	}
}

func (h *ServiceHandler) GetServices(params parking.GetParkingServicesParams, principal *models.User) middleware.Responder {
	var responder middleware.Responder
	defer utils.CatchPanic(&responder)

	ctx, span := h.tracer.Start(context.Background(), "get_services")
	defer span.End()
	traceID := fmt.Sprintf("%s", span.SpanContext().TraceID())

	services, appErr := h.service.GetServicesByParkingID(ctx, params.ParkingID)
	if appErr != nil {
		return h.handleGetError(appErr, "failed to get services", traceID, "")
	}

	return &parking.GetParkingServicesOK{
		Payload: ToAPIServiceList(services),
	}
}

func (h *ServiceHandler) UpdateService(params parking.UpdateParkingServiceParams, principal *models.User) middleware.Responder {
	var responder middleware.Responder
	defer utils.CatchPanic(&responder)

	ctx, span := h.tracer.Start(context.Background(), "update_service")
	defer span.End()
	traceID := fmt.Sprintf("%s", span.SpanContext().TraceID())

	if principal == nil {
		errCode := int64(403)
		slog.Error("failed to update service",
			slog.String("trace_id", traceID),
			slog.Int("status_code", 403),
			slog.String("error", "user not authenticated"),
		)
		return &parking.UpdateParkingServiceForbidden{
			Payload: &models.Error{
				ErrorMessage:    "Forbidden",
				ErrorStatusCode: &errCode,
			},
		}
	}

	user := ToDomainUser(principal)
	serviceModel := ToDomainService(params.Object)

	appErr := h.service.UpdateService(ctx, params.ServiceID, serviceModel, user)
	if appErr != nil {
		return h.handleUpdateError(appErr, "failed to update service", traceID, user.ID)
	}

	updated, appErr := h.service.GetServiceByID(ctx, params.ServiceID)
	if appErr != nil {
		return h.handleUpdateError(appErr, "failed to get updated service", traceID, user.ID)
	}

	return &parking.UpdateParkingServiceOK{
		Payload: ToAPIService(updated),
	}
}

func (h *ServiceHandler) DeleteService(params parking.DeleteParkingServiceParams, principal *models.User) middleware.Responder {
	var responder middleware.Responder
	defer utils.CatchPanic(&responder)

	ctx, span := h.tracer.Start(context.Background(), "delete_service")
	defer span.End()
	traceID := fmt.Sprintf("%s", span.SpanContext().TraceID())

	if principal == nil {
		errCode := int64(403)
		slog.Error("failed to delete service",
			slog.String("trace_id", traceID),
			slog.Int("status_code", 403),
			slog.String("error", "user not authenticated"),
		)
		return &parking.DeleteParkingServiceForbidden{
			Payload: &models.Error{
				ErrorMessage:    "Forbidden",
				ErrorStatusCode: &errCode,
			},
		}
	}

	user := ToDomainUser(principal)

	appErr := h.service.DeleteService(ctx, params.ServiceID, user)
	if appErr != nil {
		return h.handleDeleteError(appErr, "failed to delete service", traceID, user.ID)
	}

	return &parking.DeleteParkingServiceOK{
		Payload: &models.Result{
			Status:  "success",
			Message: "Service deleted successfully",
		},
	}
}

func (h *ServiceHandler) handleCreateError(appErr *errors.AppError, context string, traceID string, userID string) middleware.Responder {
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
		return parking.NewCreateParkingServiceNotFound().WithPayload(errorModel)
	case 400:
		return parking.NewCreateParkingServiceBadRequest().WithPayload(errorModel)
	case 403:
		return parking.NewCreateParkingServiceForbidden().WithPayload(errorModel)
	default:
		return parking.NewCreateParkingServiceBadRequest().WithPayload(errorModel)
	}
}

func (h *ServiceHandler) handleGetError(appErr *errors.AppError, context string, traceID string, userID string) middleware.Responder {
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
		return parking.NewGetParkingServiceNotFound().WithPayload(errorModel)
	default:
		return parking.NewGetParkingServiceNotFound().WithPayload(errorModel)
	}
}

func (h *ServiceHandler) handleUpdateError(appErr *errors.AppError, context string, traceID string, userID string) middleware.Responder {
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
		return parking.NewUpdateParkingServiceNotFound().WithPayload(errorModel)
	case 400:
		return parking.NewUpdateParkingServiceBadRequest().WithPayload(errorModel)
	case 403:
		return parking.NewUpdateParkingServiceForbidden().WithPayload(errorModel)
	default:
		return parking.NewUpdateParkingServiceBadRequest().WithPayload(errorModel)
	}
}

func (h *ServiceHandler) handleDeleteError(appErr *errors.AppError, context string, traceID string, userID string) middleware.Responder {
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
		return parking.NewDeleteParkingServiceNotFound().WithPayload(errorModel)
	case 403:
		return parking.NewDeleteParkingServiceForbidden().WithPayload(errorModel)
	default:
		return parking.NewDeleteParkingServiceForbidden().WithPayload(errorModel)
	}
}

