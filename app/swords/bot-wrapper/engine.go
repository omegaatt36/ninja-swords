package botwrapper

import (
	"strconv"
	"strings"

	"gopkg.in/telebot.v3"
)

// Engine defines bot engine.
type Engine struct {
	bot *telebot.Bot

	menu *telebot.ReplyMarkup

	buttons []telebot.Btn
	queries []telebot.Result
}

// Handler describe bot handler.
type Handler struct {
	Endpoint       string
	HandlerFunc    telebot.HandlerFunc
	MiddleWareFunc []telebot.MiddlewareFunc
	Describe       string
	OnQuery        bool
	OnButton       bool
}

// NewEngine return engine with bot.
func NewEngine(bot *telebot.Bot) *Engine {
	return &Engine{
		bot: bot,
		menu: &telebot.ReplyMarkup{
			ResizeKeyboard: true,
		},
	}
}

// Apply applies handler to bot.
func (e *Engine) Apply(handler Handler) {
	e.bot.Handle(handler.Endpoint, handler.HandlerFunc, handler.MiddleWareFunc...)

	if !strings.HasPrefix(handler.Endpoint, "/") {
		return
	}

	if handler.OnButton {
		e.buttons = append(e.buttons, e.menu.Text(handler.Endpoint))
	}

	if handler.OnQuery {
		result := &telebot.ArticleResult{
			Title:       handler.Endpoint,
			Text:        handler.Endpoint,
			Description: handler.Describe,
		}

		result.SetResultID(strconv.Itoa(len(e.queries) + 1))
		e.queries = append(e.queries, result)
	}
}

// Start starts bot.
func (e *Engine) Start() {
	if len(e.buttons) > 0 {
		for _, button := range e.buttons {
			e.menu.Reply(e.menu.Row(button))
		}
	}

	if len(e.queries) > 0 {
		e.bot.Handle(telebot.OnQuery, func(c telebot.Context) error {
			return c.Answer(&telebot.QueryResponse{
				Results:   e.queries,
				CacheTime: 60,
			})
		})
	}

	go func() {
		e.bot.Start()
	}()
}
