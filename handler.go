package main

import (
	"regexp"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

// An interface for update handlers.
type Handler interface {
	// This function should handle the update, and
	// react accordingly.
	Handle(*tg.Update)
	// Return a description of this handler.
	String() string
}

// A Handler for messages. This function
// gets called when a message is received.
type MessageHandler func(*tg.Message)

func (f MessageHandler) Handle(update *tg.Update) {
	if update.Message != nil {
		f(update.Message)
	}
}

func (f MessageHandler) String() string {
	return "MessageHandler"
}

// A Handler for inline messages. This function gets called
// when an inline query is received.
type InlineQueryHandler func(*tg.InlineQuery)

func (f InlineQueryHandler) Handle(update *tg.Update) {
	if update.InlineQuery != nil {
		f(update.InlineQuery)
	}
}

func (f InlineQueryHandler) String() string {
	return "InlineQueryHandler"
}

type CommandHandler struct {
	MessageHandler
	command string
}

// Create a CommandHandler from a MessageHandler. The handler is called when a message is received
// starting with /command.
func CreateCommandHandler(command string, handler MessageHandler) *CommandHandler {
	return &CommandHandler{handler, command}
}

func (ch *CommandHandler) Handle(update *tg.Update) {
	if update.Message != nil && update.Message.Command() == ch.command {
		ch.MessageHandler(update.Message)
	}
}

func (ch *CommandHandler) String() string {
	return "CommandHandler: /" + ch.command 
}

type RegexHandler struct {
	MessageHandler
	regex *regexp.Regexp
}

// Create a RegexHandler from a MessageHandler. Returns an nil and an error when the regex failed to compile.
// The handler is called when a message is received containing the regex pattern.
func CreateRegexHandler(pattern string, handler MessageHandler) (*RegexHandler, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	return &RegexHandler{handler, re}, nil
}

func (rh *RegexHandler) Handle(update *tg.Update) {
	if update.Message != nil && rh.regex.MatchString(update.Message.Text) {
		rh.MessageHandler(update.Message)
	}
}

func (rh *RegexHandler) String() string {
	return "RegexHandler: " + rh.regex.String() 
}