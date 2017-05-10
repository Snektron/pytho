package main

import (
	"log"
	"regexp"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

// A type of a function pointer representing a message handler.
type MessageHandler func(*tg.Message)

type QueryHandler func(*tg.InlineQuery)

// A generic Telegram bot which can listen for commands.
type Bot struct {
	// tgbotapi bot.
	*tg.BotAPI
	// Set to true to log debug info.
	Debug bool
	// The updates channel
	updates tg.UpdatesChannel
	// A map of regexes and their registered message handlers.
	handlers map[*regexp.Regexp]MessageHandler
	queryHandler QueryHandler
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

	bot.handlers = make(map[*regexp.Regexp]MessageHandler)

	bot.updates.Clear()	
	log.Printf("Starting on account @%s", bot.Self.UserName)

	return nil
}

// Register a handler to this bot. The handler is invoked when 
// a message is received which regex-matches pattern.
func (bot *Bot) RegisterMessageHandler(pattern string, handler MessageHandler) error {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}

	bot.handlers[re] = handler

	if (bot.Debug) {
		log.Printf("Installed handler for '%s'", pattern)
	}

	return nil
}

// Register a commandhandler. The handler is invoked when a message is
// received which starts with /command. 
func (bot *Bot) RegisterCommand(command string, handler MessageHandler) error {
	return bot.RegisterMessageHandler("^\\/" + command + "(\\s+|$|@" + bot.Self.UserName + "(\\s+|$))", handler)
}

func (bot *Bot) SetQueryHandler(handler QueryHandler) {
	bot.queryHandler = handler
}

// Listen for updates, and handle messages when they are received.
func (bot *Bot) Listen() {
	for update := range bot.updates {
		bot.handleUpdate(update)
	}
}

// Handle an incoming update
func (bot *Bot) handleUpdate(update tg.Update) {

	if bot.Debug {
		log.Printf("%+v", update)
	}

	switch {
	case update.InlineQuery != nil:
		if bot.Debug {
			log.Printf("Incoming inline query")
		}

		go bot.queryHandler(update.InlineQuery);
	case update.Message != nil:
		if bot.Debug {
			log.Printf("Incoming message: [%s] %s", update.Message.From.UserName, update.Message.Text)
		}

		bot.handleMessage(update.Message)
	}
}

// Handle a single message, calling the appropriate handlers.
func (bot *Bot) handleMessage(msg *tg.Message) {
	go func() {
		for re, handler := range bot.handlers {
			if re.MatchString(msg.Text) {
				if bot.Debug {
					log.Printf("Invoking handler for '%s'", re)
				}
				handler(msg)
			}
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