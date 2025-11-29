FROM golang:1.23-alpine

ARG PAYMENT_REST_PORT=8882
ARG PAYMENT_GRPC_PORT=50052
ARG PAYMENT_HOST=0.0.0.0

WORKDIR /pkg
COPY pkg/go.mod .
COPY pkg/go.sum .
RUN go mod download

WORKDIR /payment

COPY payment/go.mod .
COPY payment/go.sum .
RUN go mod download

COPY payment/ .
COPY pkg/ ../pkg

RUN go build -o payment cmd/main.go

EXPOSE ${PAYMENT_REST_PORT}
EXPOSE ${PAYMENT_GRPC_PORT}

ENV PAYMENT_REST_PORT=${PAYMENT_REST_PORT}
ENV PAYMENT_GRPC_PORT=${PAYMENT_GRPC_PORT}
ENV PAYMENT_HOST=${PAYMENT_HOST}

ENTRYPOINT ["./payment"]

