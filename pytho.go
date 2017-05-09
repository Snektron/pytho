package main

import (
	"bytes"
	"strings"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

var lennies = map[string]string {
	"lenny": "( ͡° ͜ʖ ͡°)",
	"fast": "ᕕ( ͡° ͜ʖ ͡°)ᕗ",
	"faster": "─=≡Σᕕ( ͡° ͜ʖ ͡°)ᕗ",
	"quite": "(ರ ͜ʖ ರೃ)",
	"magic": "╰( ͡° ͜ʖ ͡° )つ──☆*:・ﾟ",
	"almonds": "( ͡☉ ͜ʖ ͡☉)",
	"3uur": "(ง ° ͜ ʖ °)ง",
	"rip": "( ͡° ʖ̯ ͡°)",
	"hug": "(つ ͡° ͜ʖ ͡°)つ",
	"mexico": "┴┬┴┤( ͡° ͜ʖ├┬┴┬",
	"stronk": "ᕦ( ͡° ͜ʖ ͡°)ᕤ",
	"tableflip": "(ノ͡° ͜ʖ ͡°)ノ︵┻┻",
	"deag": "( ͡° ͜ʖ ͡°)=ε/̵͇̿̿/’̿’̿ ̿",
	"australia": "( ͜。 ͡ʖ ͜。)",
}

var lenniesList string

func init() {
	var buffer bytes.Buffer

	buffer.WriteString("Currently available lennies:")

	for k, v := range lennies {
		buffer.WriteString("\n")
		buffer.WriteString(k)
		buffer.WriteString(": ")
		buffer.WriteString(v)
	}

	lenniesList = buffer.String()
}

type Pytho struct {
	Bot
	lennies string
}

func (p *Pytho) Init(token string, timeout int) error {
	err := p.Bot.Init(token, timeout)
	if err != nil {
		return err
	}

	p.RegisterCommand("lennies", p.handleLennies)
	p.RegisterCommand("lenny", p.handleLenny)

	return nil
}

func (p *Pytho) handleLennies(msg *tg.Message) {
	p.QuickSend(msg, lenniesList);
}

func (p *Pytho) handleLenny(msg *tg.Message) {
	args := strings.Split(msg.Text, " ")

	if len(args) >= 2 {
		lenny, ok := lennies[strings.ToLower(args[1])]
		if !ok {
			p.QuickSend(msg, "Error: invalid lenny");
		} else {
			p.QuickSend(msg, lenny);
		}
	} else {
		p.QuickSend(msg, lennies["lenny"]);
	}
}