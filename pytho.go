package main

import (
	"log"
	"gopkg.in/telegram-bot-api.v4"
)

type Pytho struct {
	bot *tgbotapi.BotAPI
	u tgbotapi.UpdateConfig
}

func (p *Pytho) init(token string, timeout int) error {
	bot, err := tgbotapi.NewBotAPI(token)

	log.Printf("Authorized on account @%s", bot.Self.UserName)

	if err != nil {
		return err
	}

	p.bot = bot;
	p.u = tgbotapi.NewUpdate(0)
	p.u.Timeout = timeout

	return nil
}

func (p *Pytho) Listen() error {
	updates, err := p.bot.GetUpdatesChan(p.u)

	if err != nil {
		return err
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		p.bot.Send(msg)
	}

	return nil
}