package utils

import (
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
)

func ConnectTo(host *string, port *string) (*grpc.ClientConn, error) {
	return grpc.NewClient(fmt.Sprintf("%s:%s", *host, *port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
}

func ConnectToParking() (*grpc.ClientConn, error) {
	port := os.Getenv("PARKING_GRPC_PORT")
	if port == "" {
		return nil, errors.New("PARKING port is not found")
	}
	host := "parking"
	return ConnectTo(&host, &port)
}
