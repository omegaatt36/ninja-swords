package swords

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	tele "gopkg.in/telebot.v3"
	"ninja-swords/logging"
)

// Server is a telebot server.
type Server struct {
	sword *tele.Bot
}

// RegisterEndpoint installs api representation layer processing function.
func (s *Server) RegisterEndpoint() {
	var (
		// Universal markup builders.
		menu = &tele.ReplyMarkup{
			ResizeKeyboard: true,
		}

		// Reply buttons.
		btnStart    = menu.Text("/start")
		btnRegister = menu.Text("/hello")
	)

	menu.Reply(
		menu.Row(btnStart),
		menu.Row(btnRegister),
	)

	s.sword.Handle("/start", func(m tele.Context) error {
		return m.Send(fmt.Sprintf("Welcome %s", m.Message().Sender.Username), menu)
	})

	s.sword.Handle("/hello", func(m tele.Context) error {
		return m.Send("Hello World!")
	})

	s.sword.Handle(tele.OnAddedToGroup, func(m tele.Context) error {
		return m.Send("Ahoy")
	})

	s.sword.Handle(tele.OnText, func(m tele.Context) error {
		if m.Chat().Type != tele.ChatGroup {
			return nil
		}

		text := m.Text()

		return m.Send(fmt.Sprintf("echo %s", text))
	})

	s.sword.Handle(tele.OnQuery, func(c tele.Context) error {
		functions := []string{
			"/start",
			"/hello",
		}

		results := make(tele.Results, len(functions))
		for i, fn := range functions {
			result := &tele.ArticleResult{
				Title:       fn,
				Text:        fn,
				Description: fn,
			}

			results[i] = result
			results[i].SetResultID(strconv.Itoa(i + 1))
		}

		return c.Answer(&tele.QueryResponse{
			Results:   results,
			CacheTime: 60, // a minute
		})
	})
}

// Start starts telebot server.
func (s *Server) Start(ctx context.Context, botToken string) {
	sword, err := tele.NewBot(tele.Settings{
		Token:  botToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}
	s.sword = sword

	s.RegisterEndpoint()

	logging.Get().Info("starts serving bot")
	go func() {
		sword.Start()
	}()

	<-ctx.Done()
	logging.Get().Info("stops serving bot")
	sword.Stop()
}
