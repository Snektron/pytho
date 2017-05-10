package main

import (
	"regexp"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

// An interface for update handlers.
type Handler interface {
	// Handle the update.
	Handle(*tg.Update)
	// Return a description of this handler.
	String() string
}

// A Handler for messages.
type MessageHandler func(*tg.Message)

// Calls f when this update contains a message.
func (f MessageHandler) Handle(update *tg.Update) {
	if update.Message != nil {
		f(update.Message)
	}
}

func (f MessageHandler) String() string {
	return "MessageHandler"
}

// A Handler for commands.
type CommandHandler struct {
	// The callback function.
	MessageHandler
	// The command to react upon.
	command string
}

// Create a CommandHandler from a MessageHandler.
func CreateCommandHandler(command string, handler MessageHandler) *CommandHandler {
	return &CommandHandler{handler, command}
}

// Call the MessageHandler when the update contains a message and is
// this handler's command.
func (ch *CommandHandler) Handle(update *tg.Update) {
	if update.Message != nil && update.Message.Command() == ch.command {
		ch.MessageHandler(update.Message)
	}
}

func (ch *CommandHandler) String() string {
	return "CommandHandler: /" + ch.command 
}

// A Handler for regexes.
type RegexHandler struct {
	MessageHandler
	regex *regexp.Regexp
}

// Create a RegexHandler from a MessageHandler. Returns an nil and an error when the regex failed to compile.
func CreateRegexHandler(pattern string, handler MessageHandler) (*RegexHandler, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	return &RegexHandler{handler, re}, nil
}

// Call the MessageHandler when the update contains a message and matches this handler's regex.
func (rh *RegexHandler) Handle(update *tg.Update) {
	if update.Message != nil && rh.regex.MatchString(update.Message.Text) {
		rh.MessageHandler(update.Message)
	}
}

func (rh *RegexHandler) String() string {
	return "RegexHandler: " + rh.regex.String() 
}