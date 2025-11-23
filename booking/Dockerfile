FROM golang:1.23-alpine

ARG BOOKING_REST_PORT=8880
ARG BOOKING_HOST=0.0.0.0

WORKDIR /pkg
COPY pkg/go.mod .
COPY pkg/go.sum .
RUN go mod download

WORKDIR /booking

COPY booking/go.mod .
COPY booking/go.sum .
RUN go mod download

COPY booking/ .
COPY pkg/ ../pkg

RUN go build -o booking cmd/parking-booking-server/main.go

EXPOSE ${BOOKING_REST_PORT}

ENTRYPOINT ["./booking"]

