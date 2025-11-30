#!/bin/bash


swagger generate server -f booking/api/swagger/booking.yaml -t booking/internal --exclude-main --principal models.User

swagger generate server -f parking/api/swagger/parking.yaml -t parking/internal --exclude-main --principal models.User

swagger generate server -f payment/api/swagger/payment.yaml -t payment/internal --exclude-main --principal models.User

swagger generate server -f auth/api/swagger/auth.yaml -t auth/internal --exclude-main

echo "REGENERATED. NOW TIDYING"

cd booking || exit
go mod tidy

cd ../parking || exit
go mod tidy

cd ../auth || exit
go mod tidy
