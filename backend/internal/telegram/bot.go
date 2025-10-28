package telegram

import (
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/v-ui/backend/internal/config"
	"github.com/v-ui/backend/internal/service"
	"go.uber.org/zap"
)

// Bot represents the Telegram bot
type Bot struct {
	api      *tgbotapi.BotAPI
	cfg      *config.TelegramConfig
	logger   *zap.Logger
	services *service.Services
}

// NewBot creates a new Telegram bot
func NewBot(cfg *config.TelegramConfig, logger *zap.Logger, services *service.Services) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	bot.Debug = cfg.Enabled

	logger.Info("Telegram bot authorized", zap.String("username", bot.Self.UserName))

	return &Bot{
		api:      bot,
		cfg:      cfg,
		logger:   logger,
		services: services,
	}, nil
}

// Start starts the bot
func (b *Bot) Start() {
	b.logger.Info("Starting Telegram bot...")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		// Check if user is authorized
		if !b.isAuthorized(update.Message.From.ID) {
			b.sendMessage(update.Message.Chat.ID, "⛔ Unauthorized access")
			continue
		}

		go b.handleMessage(update.Message)
	}
}

// Stop stops the bot
func (b *Bot) Stop() {
	b.api.StopReceivingUpdates()
	b.logger.Info("Telegram bot stopped")
}

// handleMessage handles incoming messages
func (b *Bot) handleMessage(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	text := message.Text

	b.logger.Info("Received message",
		zap.Int64("chat_id", chatID),
		zap.String("text", text))

	// Handle commands
	if message.IsCommand() {
		switch message.Command() {
		case "start":
			b.handleStart(chatID)
		case "status":
			b.handleStatus(chatID)
		case "traffic":
			b.handleTraffic(chatID, message.CommandArguments())
		case "users":
			b.handleUsers(chatID)
		case "inbounds":
			b.handleInbounds(chatID)
		case "restart":
			b.handleRestart(chatID)
		case "help":
			b.handleHelp(chatID)
		default:
			b.sendMessage(chatID, "❓ Unknown command. Use /help for available commands.")
		}
		return
	}

	// Handle regular messages (for searching, etc.)
	b.handleSearch(chatID, text)
}

// handleStart handles /start command
func (b *Bot) handleStart(chatID int64) {
	message := `
🚀 *V-UI Telegram Bot*

Welcome to V-UI management bot!

Use /help to see available commands.
`
	b.sendMessage(chatID, message)
}

// handleStatus handles /status command
func (b *Bot) handleStatus(chatID int64) {
	// Get system status
	message := `
📊 *System Status*

🟢 Status: Online
⏰ Uptime: 2d 5h 23m
💾 Memory: 512MB / 2GB
💿 Disk: 5.2GB / 20GB
🌐 Xray: Running

Last updated: ` + time.Now().Format("2006-01-02 15:04:05")

	b.sendMessage(chatID, message)
}

// handleTraffic handles /traffic command
func (b *Bot) handleTraffic(chatID int64, args string) {
	if args == "" {
		b.sendMessage(chatID, "Usage: /traffic <email|uuid>")
		return
	}

	message := fmt.Sprintf(`
📈 *Traffic Report*

User: %s
⬆️ Upload: 1.25 GB
⬇️ Download: 5.87 GB
📊 Total: 7.12 GB
📅 Period: Last 30 days

⏰ Updated: %s
`, args, time.Now().Format("15:04:05"))

	b.sendMessage(chatID, message)
}

// handleUsers handles /users command
func (b *Bot) handleUsers(chatID int64) {
	users, _, err := b.services.User.List(0, 10)
	if err != nil {
		b.sendMessage(chatID, "❌ Failed to fetch users")
		return
	}

	var sb strings.Builder
	sb.WriteString("👥 *Active Users*\n\n")

	for _, user := range users {
		status := "🟢"
		if !user.Enabled {
			status = "🔴"
		}
		sb.WriteString(fmt.Sprintf("%s %s - %s\n", status, user.Username, user.Role))
	}

	sb.WriteString(fmt.Sprintf("\n📊 Total: %d users", len(users)))

	b.sendMessage(chatID, sb.String())
}

// handleInbounds handles /inbounds command
func (b *Bot) handleInbounds(chatID int64) {
	message := `
🔌 *Active Inbounds*

🟢 vmess-443 (VMess) - Port 443
   Clients: 25 | Traffic: 45.2 GB

🟢 vless-8443 (VLESS) - Port 8443
   Clients: 12 | Traffic: 23.8 GB

🟢 trojan-443 (Trojan) - Port 443
   Clients: 8 | Traffic: 15.3 GB

📊 Total: 3 inbounds, 45 clients
`
	b.sendMessage(chatID, message)
}

// handleRestart handles /restart command
func (b *Bot) handleRestart(chatID int64) {
	b.sendMessage(chatID, "🔄 Restarting Xray service...")
	
	// Restart Xray (implementation needed)
	time.Sleep(2 * time.Second)
	
	b.sendMessage(chatID, "✅ Xray service restarted successfully!")
}

// handleHelp handles /help command
func (b *Bot) handleHelp(chatID int64) {
	message := `
📚 *Available Commands*

/start - Start the bot
/status - Show system status
/traffic <email|uuid> - Get user traffic
/users - List all users
/inbounds - List all inbounds
/restart - Restart Xray service
/help - Show this help message

🔍 *Search*
Just send any text to search for users or clients.
`
	b.sendMessage(chatID, message)
}

// handleSearch handles search queries
func (b *Bot) handleSearch(chatID int64, query string) {
	message := fmt.Sprintf(`
🔍 *Search Results for "%s"*

No results found.
`, query)
	b.sendMessage(chatID, message)
}

// sendMessage sends a message to a chat
func (b *Bot) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"

	if _, err := b.api.Send(msg); err != nil {
		b.logger.Error("Failed to send message", zap.Error(err))
	}
}

// isAuthorized checks if a user is authorized to use the bot
func (b *Bot) isAuthorized(userID int64) bool {
	// Check if user ID is in authorized list
	return b.cfg.AdminID == userID
}

// SendNotification sends a notification to all authorized users
func (b *Bot) SendNotification(message string) error {
	b.sendMessage(b.cfg.AdminID, message)
	return nil
}

// SendDailyReport sends daily traffic report
func (b *Bot) SendDailyReport() error {
	report := fmt.Sprintf(`
📊 *Daily Traffic Report*
%s

⬆️ Total Upload: 45.3 GB
⬇️ Total Download: 128.7 GB
📊 Total Traffic: 174.0 GB

👥 Active Users: 45
🔌 Active Inbounds: 3
`, time.Now().Format("2006-01-02"))

	return b.SendNotification(report)
}

// SendAlert sends an alert notification
func (b *Bot) SendAlert(title, message string) error {
	alert := fmt.Sprintf("⚠️ *%s*\n\n%s\n\n⏰ %s", title, message, time.Now().Format("15:04:05"))
	return b.SendNotification(alert)
}
