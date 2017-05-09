package main

import (
	"log"
	"regexp"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Handler interface {
	handle(*Bot, tg.Message)
}

type Bot struct {
	*tg.BotAPI
	updates tg.UpdatesChannel
	handlers map[*regexp.Regexp]Handler
}

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

	bot.handlers = make(map[*regexp.Regexp]Handler)

	bot.updates.Clear()	
	log.Printf("Authorized on account @%s", bot.Self.UserName)

	return nil
}

func (bot *Bot) Listen() {
	for update := range bot.updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tg.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}

func (bot *Bot) Register(pattern string, handler Handler) error {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}

	bot.handlers[re] = handler
	return nil
}

func (bot *Bot) RegisterCommand(cmd string, handler Handler) error {
	return bot.Register("^\\/" + cmd, handler)
}

func (bot *Bot) Issue(msg tg.Message) {
	for re, handler := range bot.handlers {
		if re.MatchString(msg.Text) {
			go handler.handle(bot, msg)
		}
	}
}