package main

import (
	"log"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

// A generic Telegram bot which can listen for commands.
type Bot struct {
	// tgbotapi bot.
	*tg.BotAPI
	// Set to true to log debug info.
	Debug bool
	// The channel for incoming updates.
	updates tg.UpdatesChannel
	// A list of registered handlers.
	handlers []Handler
}

// Initialize this bot with a Telegram API-token and a timeout in seconds.
func (bot *Bot) Init(token string, timeout int) error {
	var err error
	bot.BotAPI, err = tg.NewBotAPI(token)

	if err != nil {
		return err
	}

	u := tg.NewUpdate(0)
	u.Timeout = timeout

	bot.updates, err = bot.GetUpdatesChan(u)

	if err != nil {
		return err
	}

	bot.updates.Clear()	
	log.Printf("Starting on account @%s", bot.Self.UserName)

	if bot.Debug {
		bot.Register(MessageHandler(debugMessageHandler))
		bot.Register(MessageHandler(debugCommandHandler))
		bot.Register(InlineQueryHandler(debugInlineQueryHandler))
	}

	return nil
}

// Register a handler to the bot.
func (bot *Bot) Register(handler Handler) {
	bot.handlers = append(bot.handlers, handler)
	log.Printf("Registered a %s", handler.String())
}

// Listen for updates, and handle messages when they are received.
func (bot *Bot) Listen() {
	for update := range bot.updates {
		bot.handleUpdate(&update)
	}
}

// Handle an incoming update. 
func (bot *Bot) handleUpdate(update *tg.Update) {
	go func() {
		for _, handler := range bot.handlers {
			handler.Handle(update)
		}
	}()
}

// Helper function to easily send a message to the chat another message
// appears in.
func (bot *Bot) QuickSend(msg *tg.Message, text string) {
	bot.Send(tg.NewMessage(msg.Chat.ID, text))
}

// Helper function to easily reply to a message
func (bot *Bot) QuickReply(msg *tg.Message, text string) {
	reply := tg.NewMessage(msg.Chat.ID, text)
	reply.ReplyToMessageID = msg.MessageID
	bot.Send(reply)
}

func debugMessageHandler(msg *tg.Message) {
	log.Printf("[%s]: %s", msg.From.UserName, msg.Text)
}

func debugCommandHandler(msg *tg.Message) {
	if msg.IsCommand() {
		log.Printf("Received command from %s: /%s", msg.From.UserName, msg.Command())
	}
}

func debugInlineQueryHandler(iq *tg.InlineQuery) {
	log.Printf("Received inline query from %s: %s", iq.From.UserName, iq.Query)
}