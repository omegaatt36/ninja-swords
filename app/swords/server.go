package swords

import (
	"context"
	"log"
	"time"

	botwrapper "github.com/omegaatt36/ninja-swords/app/swords/bot-wrapper"
	"github.com/omegaatt36/ninja-swords/logging"
	"gopkg.in/telebot.v3"
)

// Server is a telebot server.
type Server struct {
}

// Start starts telebot server.
func (s *Server) Start(ctx context.Context, botToken string) {
	sword, err := telebot.NewBot(telebot.Settings{
		Token:  botToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	engine := botwrapper.NewEngine(sword)

	registerHandler(engine)

	logging.Get().Info("starts serving bot")
	engine.Start()

	<-ctx.Done()
	logging.Get().Info("stops serving bot")
	sword.Stop()
}
