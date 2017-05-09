package main

import (
	"log"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type Pytho struct {
	bot *tgbotapi.BotAPI
	updates tgbotapi.UpdatesChannel
}

func (p *Pytho) Init(token string, timeout int) error {
	bot, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		return err
	}

	log.Printf("Authorized on account @%s", bot.Self.UserName)
	
	u := tgbotapi.NewUpdate(0)
	u.Timeout = timeout

	updates, err := bot.GetUpdatesChan(u)

	if err != nil {
		return err
	}

		updates.Clear()

	p.bot = bot
	p.updates = updates

	return nil
}

func (p *Pytho) Listen() {
	for update := range p.updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		p.bot.Send(msg)
	}
}