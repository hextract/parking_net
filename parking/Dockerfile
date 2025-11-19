FROM golang:1.23-alpine

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

EXPOSE $PARKING_REST_PORT
EXPOSE $PARKING_GRPC_PORT

ENTRYPOINT ["./parking"]


