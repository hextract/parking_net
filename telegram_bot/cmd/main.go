package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"telegram_bot/api_service"
	"telegram_bot/data_representation"
	"telegram_bot/internal/utils"
	"telegram_bot/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type LoginState struct {
	WaitingForLogin    bool
	WaitingForPassword bool
	Login              string
	PasswordMessageID  int
	CreatedAt          time.Time
}

type LoginStateManager struct {
	states map[int64]*LoginState
	mu     sync.RWMutex
}

func NewLoginStateManager() *LoginStateManager {
	manager := &LoginStateManager{
		states: make(map[int64]*LoginState),
	}
	go manager.cleanupExpired()
	return manager
}

func (m *LoginStateManager) Get(telegramID int64) (*LoginState, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	state, exists := m.states[telegramID]
	return state, exists
}

func (m *LoginStateManager) Set(telegramID int64, state *LoginState) {
	m.mu.Lock()
	defer m.mu.Unlock()
	state.CreatedAt = time.Now()
	m.states[telegramID] = state
}

func (m *LoginStateManager) Delete(telegramID int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.states, telegramID)
}

func (m *LoginStateManager) cleanupExpired() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		m.mu.Lock()
		now := time.Now()
		for id, state := range m.states {
			if now.Sub(state.CreatedAt) > 10*time.Minute {
				delete(m.states, id)
			}
		}
		m.mu.Unlock()
	}
}

func main() {
	telegramApiKey := os.Getenv("TELEGRAM_API_KEY")
	bot, err := tgbotapi.NewBotAPI(telegramApiKey)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	apiService, err := api_service.NewService()
	if err != nil {
		log.Panic(err)
	}

	ctx := context.Background()
	loginStateManager := NewLoginStateManager()

	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		telegramID := update.Message.From.ID

		if err := utils.ValidateTelegramID(int64(telegramID)); err != nil {
			continue
		}

		msg := tgbotapi.NewMessage(chatID, "")

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				msg.Text = fmt.Sprintf("Hello! 👋\n\nYour Telegram ID: %d\n\nUse this ID when registering on the website.\n\nAvailable commands:\n/help - show all commands", telegramID)

			case "help":
				msg.Text = "Available commands:\n\n"
				msg.Text += "/start - greeting and your Telegram ID\n"
				msg.Text += "/login - authorize (enter login and password)\n"
				msg.Text += "/balance - view balance\n"
				msg.Text += "/bookings - view my bookings (for drivers)\n"
				msg.Text += "/parkings - view my parking places (for owners)\n"
				msg.Text += "/help - show this help"

			case "login":
				if err := utils.ValidateTelegramID(int64(telegramID)); err != nil {
					msg.Text = "❌ Error: invalid Telegram ID."
					break
				}
				loginStateManager.Set(int64(telegramID), &LoginState{
					WaitingForLogin: true,
				})
				msg.Text = "🔐 Enter your login:"

			case "balance":
				if err := utils.ValidateTelegramID(int64(telegramID)); err != nil {
					msg.Text = "❌ Error: invalid Telegram ID."
					break
				}
				userInfo, err := apiService.GetUserByTelegramID(ctx, int64(telegramID))
				if err != nil {
					msg.Text = "❌ User not found. Please register on the website using your Telegram ID: " + strconv.FormatInt(int64(telegramID), 10)
				} else {
					if userInfo.Token == "" {
						msg.Text = "❌ Authorization required. Use the /login command to sign in."
					} else {
						balance, err := apiService.GetBalance(userInfo.Token)
						if err != nil {
							msg.Text = "❌ Error getting balance."
						} else {
							balanceValue := int64(0)
							if balance.Balance != nil {
								balanceValue = *balance.Balance
							}
							currency := "USD"
							if balance.Currency != nil {
								currency = *balance.Currency
							}
							msg.Text = fmt.Sprintf("💰 Your balance: %.2f %s", float64(balanceValue)/100.0, currency)
						}
					}
				}

			case "bookings":
				if err := utils.ValidateTelegramID(int64(telegramID)); err != nil {
					msg.Text = "❌ Error: invalid Telegram ID."
					break
				}
				userInfo, err := apiService.GetUserByTelegramID(ctx, int64(telegramID))
				if err != nil {
					msg.Text = "❌ User not found. Please register on the website using your Telegram ID: " + strconv.FormatInt(int64(telegramID), 10)
				} else {
					if userInfo.Role != "driver" {
						msg.Text = "❌ This command is only available for drivers."
					} else if userInfo.Token == "" {
						msg.Text = "❌ Authorization required. Use the /login command to sign in."
					} else {
						user := &models.User{
							UserID:     userInfo.UserID,
							TelegramID: int64(telegramID),
						}
						bookings, err := apiService.GetUserBookings(user)
						if err != nil {
							msg.Text = "❌ Error getting bookings."
						} else {
							if len(bookings) == 0 {
								msg.Text = "📋 You have no active bookings."
							} else {
								var builder strings.Builder
								builder.WriteString(fmt.Sprintf("📋 Your bookings (%d):\n\n", len(bookings)))
								for i, booking := range bookings {
									if i > 0 {
										builder.WriteString("\n---\n\n")
									}
									bookingText := data_representation.GetBooking(&booking)
									if err := utils.ValidateMessageLength(builder.String() + bookingText); err != nil {
										builder.WriteString("... (message too long)")
										break
									}
									builder.WriteString(bookingText)
								}
								msg.Text = builder.String()
							}
						}
					}
				}

			case "parkings":
				if err := utils.ValidateTelegramID(int64(telegramID)); err != nil {
					msg.Text = "❌ Error: invalid Telegram ID."
					break
				}
				userInfo, err := apiService.GetUserByTelegramID(ctx, int64(telegramID))
				if err != nil {
					msg.Text = "❌ User not found. Please register on the website using your Telegram ID: " + strconv.FormatInt(int64(telegramID), 10)
				} else {
					if userInfo.Role != "owner" {
						msg.Text = "❌ This command is only available for parking owners."
					} else if userInfo.Token == "" {
						msg.Text = "❌ Authorization required. Use the /login command to sign in."
					} else {
						user := &models.User{
							UserID:     userInfo.UserID,
							TelegramID: int64(telegramID),
						}
						parkings, err := apiService.GetUserParkings(user)
						if err != nil {
							msg.Text = "❌ Error getting parking places."
						} else {
							if len(parkings) == 0 {
								msg.Text = "🅿️ You have no parking places."
							} else {
								var builder strings.Builder
								builder.WriteString(fmt.Sprintf("🅿️ Your parking places (%d):\n\n", len(parkings)))
								for i, parking := range parkings {
									if i > 0 {
										builder.WriteString("\n---\n\n")
									}
									parkingText := data_representation.GetParkingPlace(&parking)
									if err := utils.ValidateMessageLength(builder.String() + parkingText); err != nil {
										builder.WriteString("... (message too long)")
										break
									}
									builder.WriteString(parkingText)
								}
								msg.Text = builder.String()
							}
						}
					}
				}

			default:
				msg.Text = "❌ Unknown command. Use /help for a list of commands."
			}
		} else {
			state, inLogin := loginStateManager.Get(int64(telegramID))
			if inLogin {
				if state.WaitingForLogin {
					login := strings.TrimSpace(update.Message.Text)
					if err := utils.ValidateLogin(login); err != nil {
						msg.Text = "❌ Invalid login format. Login must contain only letters, numbers, and underscores (3-50 characters)."
						loginStateManager.Delete(int64(telegramID))
					} else {
						state.Login = login
						state.WaitingForLogin = false
						state.WaitingForPassword = true
						msg.Text = "🔐 Enter your password:"
						sentMsg, err := bot.Send(msg)
						if err == nil {
							state.PasswordMessageID = sentMsg.MessageID
							loginStateManager.Set(int64(telegramID), state)
						}
					}
					continue
				} else if state.WaitingForPassword {
					password := strings.TrimSpace(update.Message.Text)

					if err := utils.ValidatePassword(password); err != nil {
						msg.Text = "❌ Invalid password format. Password must be 8-128 characters and contain uppercase, lowercase letters, and numbers."
						loginStateManager.Delete(int64(telegramID))
						continue
					}

					deleteMsg := tgbotapi.NewDeleteMessage(chatID, update.Message.MessageID)
					bot.Send(deleteMsg)

					if state.PasswordMessageID > 0 {
						deletePasswordMsg := tgbotapi.NewDeleteMessage(chatID, state.PasswordMessageID)
						bot.Send(deletePasswordMsg)
					}

					user := &models.User{
						Login:      &state.Login,
						Password:   &password,
						TelegramID: int64(telegramID),
					}

					success, err := apiService.Login(user)
					if err != nil || !success {
						msg.Text = "❌ Authorization error. Check your login and password."
					} else {
						msg.Text = "✅ Authorization successful! You can now use commands."
					}
					loginStateManager.Delete(int64(telegramID))
				}
			} else {
				msg.Text = "Use commands starting with /. For example, /help"
			}
		}

		if msg.Text != "" {
			if err := utils.ValidateMessageLength(msg.Text); err != nil {
				msg.Text = "❌ Error: message too long."
			}
			if _, err := bot.Send(msg); err != nil {
				log.Printf("Error sending message: %v", err)
			}
		}
	}
}
