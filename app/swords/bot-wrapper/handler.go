package botwrapper

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type handler HandlerFunc

var defaultOnError = func(err error, c *Context) {
	if c != nil {
		log.Println(c.UpdateID, err)
	} else {
		log.Println(err)
	}
}

// HandlerFunc is middleware or handler.
type HandlerFunc func(*Context)

func (e *Engine) processUpdate(endpoint string, u tgbotapi.Update) error {
	c := &Context{Update: u, bot: e.bot}
	if handler, ok := e.handlers[endpoint]; ok {
		e.runHandler(handler, c)
	}
	return nil
}

func (e *Engine) runHandler(h handler, c *Context) {
	f := func() { h(c) }
	if e.synchronous {
		f()
	} else {
		go f()
	}
}
