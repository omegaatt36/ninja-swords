package botwrapper

import (
	"math"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// abortIndex represents a typical value used in abort functions.
const abortIndex int8 = math.MaxInt8 >> 1

// Context is a wrapper.
type Context struct {
	tgbotapi.Update

	bot      *tgbotapi.BotAPI
	Endpoint string
	handlers []HandlerFunc
	index    int8
}

// Next should be used only inside middleware.
// It executes the pending handlers in the chain inside the calling handler.
// See example in GitHub.
func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}

// IsAborted returns true if the current context was aborted.
func (c *Context) IsAborted() bool {
	return c.index >= abortIndex
}

// Abort prevents pending handlers from being called. Note that this will not stop the current handler.
// Let's say you have an authorization middleware that validates that the current request is authorized.
// If the authorization fails (ex: the password does not match), call Abort to ensure the remaining handlers
// for this request are not called.
func (c *Context) Abort() {
	c.index = abortIndex
}

// Send is a wrapper.
func (c *Context) Send(chattable tgbotapi.Chattable) {
	if _, err := c.bot.Send(chattable); err != nil {
		c.index = abortIndex
	}
}
