FROM golang:1.23-alpine

WORKDIR /pkg
COPY pkg/go.mod .
COPY pkg/go.sum .
RUN go mod download

WORKDIR /notification
COPY notification/go.mod .
COPY notification/go.sum .
RUN go mod download

COPY notification/ .
COPY pkg/ ../pkg

RUN go mod tidy
RUN go build -o notification cmd/main.go

ENTRYPOINT ["./notification"]


