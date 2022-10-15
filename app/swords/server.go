package swords

import (
	"context"
	"log"
	"time"

	"gopkg.in/telebot.v3"
	botwrapper "ninja-swords/app/swords/bot-wrapper"
	"ninja-swords/app/swords/module/user"
	"ninja-swords/logging"
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

	userX := user.NewController()
	userX.RegisterHandler(engine)

	logging.Get().Info("starts serving bot")
	engine.Start()

	<-ctx.Done()
	logging.Get().Info("stops serving bot")
	sword.Stop()
}
