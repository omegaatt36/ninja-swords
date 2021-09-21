package swords

import (
	"context"
	"fmt"
	"log"
	"ninga-swords/logging"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

// Server is a telebot server.
type Server struct {
}

// RegisterEndpoint installs api representation layer processing function.
func (s *Server) RegisterEndpoint(b *tb.Bot) {
	b.Handle(tb.OnText, func(m *tb.Message) {
		b.Send(m.Sender, fmt.Sprintf("echo \"%s\"", m.Text))
	})

	b.Handle("/hello", func(m *tb.Message) {
		b.Send(m.Sender, "Hello World!")
	})
}

// Start starts telebot server.
func (s *Server) Start(ctx context.Context, botToken string) {
	sword, err := tb.NewBot(tb.Settings{
		Token:  botToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	s.RegisterEndpoint(sword)

	logging.Get().Info("starts serving bot")
	go func() {
		sword.Start()
	}()

	<-ctx.Done()
	logging.Get().Info("stops serving bot")
	sword.Stop()
}
