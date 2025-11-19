module github.com/h4x4d/parking_net/notification

go 1.23.1

require (
	github.com/h4x4d/parking_net/pkg v0.0.0-00010101000000-000000000000
	github.com/segmentio/kafka-go v0.4.47
	go.opentelemetry.io/otel/trace v1.32.0
)

replace github.com/h4x4d/parking_net/pkg => ../pkg

require (
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
)
