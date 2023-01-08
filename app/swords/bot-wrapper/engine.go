package botwrapper

import (
	"context"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"ninja-swords/logging"
)

// Engine defines bot engine.
type Engine struct {
	bot *tgbotapi.BotAPI

	synchronous     bool
	onError         func(error, *Context)
	handlers        map[string]handler
	middlewareFuncs []HandlerFunc
}

// Use adds middleware to the engine.
func (e *Engine) Use(middleware ...HandlerFunc) {
	e.middlewareFuncs = append(e.middlewareFuncs, middleware...)
}

// NewEngine return engine with bot.
func NewEngine(bot *tgbotapi.BotAPI) *Engine {
	return &Engine{
		bot:      bot,
		onError:  defaultOnError,
		handlers: make(map[string]handler),
	}
}

// Apply applies handler to bot.
func (e *Engine) Apply(endpoint string, handler HandlerFunc, middlewareFuncs ...HandlerFunc) {
	if _, ok := e.handlers[endpoint]; ok {
		panic(fmt.Errorf("endpoint(%s) is set repeatedly", endpoint))
	}

	e.handlers[endpoint] = func(c *Context) {
		if len(e.middlewareFuncs) > 0 {
			middlewareFuncs = append(e.middlewareFuncs, middlewareFuncs...)
		}

		for index := 0; index < len(middlewareFuncs) && !c.IsAborted(); index++ {
			middlewareFuncs[index](c)
		}

		if c.IsAborted() {
			return
		}

		handler(c)
	}
}

// Start starts bot.
func (e *Engine) Start(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := e.bot.GetUpdatesChan(u)
	go func() {
		for ctx.Err() == nil {
			select {
			case <-ctx.Done():
				if err := ctx.Err(); err != nil {
					logging.Get().Error(err)
				}
				return
			case update := <-updates:
				if update.Message == nil {
					continue
				}

				if !update.Message.IsCommand() {
					continue
				}

				log.Printf("\033[1;30;32m[%s] %s\033[0m",
					update.Message.From.UserName,
					update.Message.Command(),
				)
				e.processUpdate(update.Message.Command(), update)
			}
		}
	}()
}
