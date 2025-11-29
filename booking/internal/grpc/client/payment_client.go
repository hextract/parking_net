package client

import (
	"context"
	"fmt"

	"github.com/h4x4d/parking_net/booking/internal/grpc/gen"
	"github.com/h4x4d/parking_net/booking/internal/grpc/utils"
	"go.opentelemetry.io/otel"
)

type PaymentClient struct{}

func NewPaymentClient() *PaymentClient {
	return &PaymentClient{}
}

func (pc *PaymentClient) ProcessTransaction(ctx context.Context, bookingID int64, driverID string, ownerID string, amount int64) (*TransactionResponse, error) {
	conn, err := utils.ConnectToPayment()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to payment service: %w", err)
	}
	defer conn.Close()

	tracer := otel.Tracer("Booking")
	childCtx, span := tracer.Start(ctx, "booking request process transaction")
	defer span.End()

	client := gen.NewPaymentClient(conn)

	req := &gen.TransactionRequest{
		BookingId: bookingID,
		DriverId:  driverID,
		OwnerId:   ownerID,
		Amount:    amount,
	}

	resp, err := client.ProcessTransaction(childCtx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to process transaction: %w", err)
	}

	return &TransactionResponse{
		TransactionID: resp.TransactionId,
		Status:        resp.Status,
		Message:       resp.Message,
	}, nil
}

func (pc *PaymentClient) ProcessRefund(ctx context.Context, bookingID int64, driverID string, ownerID string, amount int64) (*TransactionResponse, error) {
	conn, err := utils.ConnectToPayment()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to payment service: %w", err)
	}
	defer conn.Close()

	tracer := otel.Tracer("Booking")
	childCtx, span := tracer.Start(ctx, "booking request process refund")
	defer span.End()

	client := gen.NewPaymentClient(conn)

	req := &gen.RefundRequest{
		BookingId: bookingID,
		DriverId:  driverID,
		OwnerId:   ownerID,
		Amount:    amount,
	}

	resp, err := client.ProcessRefund(childCtx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to process refund: %w", err)
	}

	return &TransactionResponse{
		TransactionID: resp.TransactionId,
		Status:        resp.Status,
		Message:       resp.Message,
	}, nil
}

type TransactionResponse struct {
	TransactionID int64
	Status        string
	Message       string
}
