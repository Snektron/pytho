package main

import (
	"log"
	"regexp"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Handler func(*tg.Message)

type Bot struct {
	*tg.BotAPI
	Debug bool
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
	log.Printf("Starting on account @%s", bot.Self.UserName)

	return nil
}

func (bot *Bot) Register(pattern string, handler Handler) error {
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

func (bot *Bot) RegisterCommand(cmd string, handler Handler) error {
	return bot.Register("^\\/" + cmd, handler)
}

func (bot *Bot) Listen() {
	for update := range bot.updates {
		if update.Message == nil {
			continue
		}

		if bot.Debug {
			log.Printf("Incoming message: [%s] %s", update.Message.From.UserName, update.Message.Text)
		}

		bot.handle(update.Message)
	}
}

func (bot *Bot) handle(msg *tg.Message) {
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

func (bot *Bot) QuickSend(origin *tg.Message, text string) {
	bot.Send(tg.NewMessage(origin.Chat.ID, text))
}