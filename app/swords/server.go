package swords

import (
	"context"
	"log"
	"ninga-swords/logging"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

// Server is a telebot server.
type Server struct {
}

// Start starts telebot server.
func (s *Server) Start(ctx context.Context, botToken string) {
	sowrd, err := tb.NewBot(tb.Settings{
		Token:  botToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	sowrd.Handle("/hello", func(m *tb.Message) {
		sowrd.Send(m.Sender, "Hello World!")
	})

	logging.Get().Info("starts serving bot")
	sowrd.Start()

	<-ctx.Done()
	logging.Get().Info("stops serving bot")
	sowrd.Stop()
}
