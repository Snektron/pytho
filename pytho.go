package main

import (
	"bytes"
	"strings"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

// All lennies Pytho provides
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
	"chong": "( ͡- ͜ʖ ͡-)",
	"meh": "( ͡°_ʖ ͡° )",
}

// A formatted list to reply to '/lennies'.
var lenniesList string

// Initialize this module, generate lenniesList
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

// Pytho bot.
type Pytho struct {
	// The underlying bot.
	Bot
}

// Initialize Pytho with a Telegram API-token and a timeout in seconds.
func (p *Pytho) Init(token string, timeout int) error {
	err := p.Bot.Init(token, timeout)
	if err != nil {
		return err
	}

	p.RegisterCommand("lennies", p.handleLennies)
	p.RegisterCommand("lenny", p.handleLenny)

	return nil
}

// Handle the '/lennies' command
func (p *Pytho) handleLennies(msg *tg.Message) {
	p.QuickSend(msg, lenniesList);
}

// Handle the '/lenny' command
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