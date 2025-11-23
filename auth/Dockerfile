FROM golang:1.23-alpine

ARG AUTH_REST_PORT=8800
ARG AUTH_HOST=0.0.0.0

WORKDIR /pkg
COPY pkg/go.mod .
COPY pkg/go.sum .
RUN go mod download

WORKDIR /auth
COPY auth/go.mod .
COPY auth/go.sum .
RUN go mod download

COPY auth/ .

COPY pkg/ ../pkg

RUN go build -o auth cmd/main.go

EXPOSE ${AUTH_REST_PORT}

ENTRYPOINT ["./auth"]

