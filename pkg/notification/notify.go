package notification

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/h4x4d/parking_net/pkg/models"
	"github.com/segmentio/kafka-go"
)

func (kc KafkaConnection) SendNotification(notification models.Notification) error {
	if kc.Writer == nil {
		return errors.New("kafka connection is not initialized")
	}

	if notification.Name == "" {
		return errors.New("notification name is required")
	}
	if notification.Text == "" {
		return errors.New("notification text is required")
	}
	if notification.TelegramID <= 0 {
		return errors.New("telegram ID must be greater than 0")
	}

	notify, marshalErr := json.Marshal(notification)
	if marshalErr != nil {
		return fmt.Errorf("failed to marshal notification: %w", marshalErr)
	}
	_, err := kc.Writer.WriteMessages(
		kafka.Message{
			Key:   []byte("send_notification"),
			Value: notify,
			Headers: []kafka.Header{
				{
					Key:   "format",
					Value: []byte("json"),
				},
			},
		})
	if err != nil {
		return fmt.Errorf("failed to write message to kafka: %w", err)
	}
	return nil
}
