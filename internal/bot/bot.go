package bot

import (
	"URLShortener/internal/config"
	"URLShortener/internal/service"
	"context"
	"fmt"
	"net/url"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

// ... (constants remain the same)

const (
	msgHelp = `
Hello! I am a URL Shortener bot.
Just send me a valid URL and I will shorten it for you.
For example: https://google.com
`
	msgInvalidURL                = "Please send me a valid URL. Example: https://google.com"
	msgShortenFailed             = "Sorry, I could not shorten this URL. Please try again later."
	msgLinkSuccessfullyShortened = "Here is your shortened link: `%s`" // Updated for markdown
)

// ... (Bot struct remains the same)
type Bot struct {
	api          *tgbotapi.BotAPI
	log          *zap.Logger
	config       *config.Config
	urlShortener *service.URLShortenerService
}

// ... (New method remains the same)
func New(cfg *config.Config, log *zap.Logger, shortener *service.URLShortenerService) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		return nil, err
	}

	log.Info("authorized on account", zap.String("username", api.Self.UserName))

	return &Bot{
		api:          api,
		log:          log,
		config:       cfg,
		urlShortener: shortener,
	}, nil
}

// Start runs the bot's main loop in a separate goroutine.
// It now accepts a context to handle graceful shutdown.
func (b *Bot) Start(ctx context.Context) {
	b.log.Info("starting bot")
	updates := b.getUpdatesChannel()

	// Using a goroutine to not block the main application thread.
	go func() {
		for {
			select {
			case <-ctx.Done(): // If the context is cancelled, stop the bot.
				b.log.Info("stopping bot...")
				b.api.StopReceivingUpdates()
				return
			case update := <-updates: // Process updates as they come.
				b.processUpdate(update)
			}
		}
	}()
}

// processUpdate centralizes the logic for handling an update.
func (b *Bot) processUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	if update.Message.IsCommand() {
		if err := b.handleCommand(update.Message); err != nil {
			b.log.Error("failed to handle command", zap.Error(err))
		}
		return
	}

	if err := b.handleMessage(update.Message); err != nil {
		b.log.Error("failed to handle message", zap.Error(err))
	}
}

// ... (handleCommand remains the same)
func (b *Bot) handleCommand(msg *tgbotapi.Message) error {
	switch msg.Command() {
	case "start":
		return b.sendMessage(msg.Chat.ID, msgHelp)
	default:
		return b.sendMessage(msg.Chat.ID, "I don't know this command.")
	}
}

func (b *Bot) handleMessage(msg *tgbotapi.Message) error {
	log := b.log.With(zap.String("username", msg.From.UserName), zap.Int64("chat_id", msg.Chat.ID))
	log.Info("received message", zap.String("text", msg.Text))

	_, err := url.ParseRequestURI(msg.Text)
	if err != nil {
		log.Debug("message is not a valid URL")
		return b.sendMessage(msg.Chat.ID, msgInvalidURL)
	}

	alias, err := b.urlShortener.Shorten(context.Background(), msg.Text)
	if err != nil {
		log.Error("failed to shorten URL", zap.Error(err))
		return b.sendMessage(msg.Chat.ID, msgShortenFailed)
	}

	// Use the BaseURL from config to construct the full shortened URL
	shortURL := fmt.Sprintf("%s/%s", b.config.HTTPServer.BaseURL, alias)

	// Send the message with markdown formatting for the link
	return b.sendMessage(msg.Chat.ID, fmt.Sprintf(msgLinkSuccessfullyShortened, shortURL))
}

// sendMessage is updated to support markdown.
func (b *Bot) sendMessage(chatID int64, text string) error {
	reply := tgbotapi.NewMessage(chatID, text)
	reply.ParseMode = tgbotapi.ModeMarkdown // Enable markdown parsing
	_, err := b.api.Send(reply)
	return err
}

// ... (getUpdatesChannel remains the same)
func (b *Bot) getUpdatesChannel() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	return b.api.GetUpdatesChan(u)
}
