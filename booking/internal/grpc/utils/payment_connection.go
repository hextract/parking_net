package utils

import (
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ConnectToPayment() (*grpc.ClientConn, error) {
	address := os.Getenv("PAYMENT_GRPC_ADDRESS")
	if address == "" {
		address = "payment:50052"
	}

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return conn, nil
}

