package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/h4x4d/parking_net/notification/internal/models"
)

type SendNotificationRequest struct {
	ChatId string `json:"chat_id"`
	Text   string `json:"text"`
}

func SendNotification(notification models.Notification) error {
	if notification.Name == "" {
		return errors.New("notification name is required")
	}
	if notification.Text == "" {
		return errors.New("notification text is required")
	}
	if notification.TelegramID <= 0 {
		return fmt.Errorf("invalid telegram ID: %d (must be positive)", notification.TelegramID)
	}

	apiKey := os.Getenv("TELEGRAM_API_KEY")
	if apiKey == "" {
		return errors.New("environment variable TELEGRAM_API_KEY is not set")
	}

	telegramAPIURL := fmt.Sprintf("https://api.telegram.org/bot%s", apiKey)

	requestBody, err := json.Marshal(SendNotificationRequest{
		ChatId: strconv.Itoa(notification.TelegramID),
		Text:   fmt.Sprintf("%s\n\n%s", notification.Name, notification.Text),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Post(telegramAPIURL+"/sendMessage", "application/json",
		bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var bodyBytes []byte
		if resp.Body != nil {
			bodyBytes, _ = io.ReadAll(io.LimitReader(resp.Body, 512))
		}
		return fmt.Errorf("telegram API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	slog.Info(
		"send notification",
		slog.Group("notification-properties",
			slog.String("name", notification.Name),
			slog.String("text", notification.Text),
			slog.Int("telegram-id", notification.TelegramID),
		),
		slog.Int("status_code", http.StatusOK),
	)
	return nil
}
