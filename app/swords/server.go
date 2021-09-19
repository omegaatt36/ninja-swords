package swords

import (
	"context"
	"fmt"
	"log"
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
		btnRegister = menu.Text("/register")
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

	s.sword.Handle("/register", func(m tele.Context) error {
		return nil
	})

	s.sword.Handle("/create_auto_reply", func(m tele.Context) error {
		// tele.CallbackEndpoint
		return nil
	})

	s.sword.Handle(tele.OnSticker, func(m tele.Context) error {
		if m.Chat().Type == tele.ChatGroup {
			return nil
		}

		logging.Get().Debugf("%d-%s send %s",
			m.Message().Sender.ID, m.Message().Sender.Username, m.Message().Sticker.Emoji)
		return m.Send(m.Message().Sticker)
	})

	s.sword.Handle(tele.OnText, func(m tele.Context) error {
		if m.Chat().Type != tele.ChatGroup {
			return nil
		}

		logging.Get().Debugf("%d-%s send %s",
			m.Message().Sender.ID, m.Message().Sender.Username, m.Message().Text)
		return m.Send(fmt.Sprintf("echo \"%s\"", m.Text()))
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
