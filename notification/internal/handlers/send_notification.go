package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/h4x4d/parking_net/notification/internal/models"
	"github.com/h4x4d/parking_net/notification/internal/services"
	"log/slog"
	"net/http"
)

func SendNotificationHandler(value []byte) error {
	if len(value) == 0 {
		return fmt.Errorf("notification payload is empty")
	}

	request := models.Notification{}
	err := json.Unmarshal(value, &request)
	if err != nil {
		return fmt.Errorf("failed to unmarshal notification: %w", err)
	}

	if request.Name == "" {
		return fmt.Errorf("notification name is required")
	}
	if request.Text == "" {
		return fmt.Errorf("notification text is required")
	}
	if request.TelegramID <= 0 {
		return fmt.Errorf("invalid telegram ID: must be positive")
	}

	err = services.SendNotification(request)
	if err != nil {
		slog.Error(
			"failed send notification",
			slog.Group("notification-properties",
				slog.String("name", request.Name),
				slog.String("text", request.Text),
				slog.Int("telegram-id", request.TelegramID),
			),
			slog.Int("status_code", http.StatusInternalServerError),
			slog.String("error", err.Error()),
		)
	}
	return err
}
