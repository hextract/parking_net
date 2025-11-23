FROM golang:1.23-alpine

ARG PARKING_REST_PORT=8888
ARG PARKING_GRPC_PORT=50051
ARG PARKING_HOST=0.0.0.0

WORKDIR /pkg
COPY pkg/go.mod .
COPY pkg/go.sum .
RUN go mod download

WORKDIR /parking
COPY parking/go.mod .
COPY parking/go.sum .
RUN go mod download

COPY parking/ .

COPY pkg/ ../pkg

RUN cd ../pkg && go mod tidy && cd /parking && go mod tidy && go build -o parking cmd/main.go

EXPOSE ${PARKING_REST_PORT}
EXPOSE ${PARKING_GRPC_PORT}

ENV PARKING_REST_PORT=${PARKING_REST_PORT}
ENV PARKING_GRPC_PORT=${PARKING_GRPC_PORT}
ENV PARKING_HOST=${PARKING_HOST}

ENTRYPOINT ["./parking"]


