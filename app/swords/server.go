package swords

import (
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	botwrapper "ninja-swords/app/swords/bot-wrapper"
	"ninja-swords/app/swords/module/user"
	"ninja-swords/logging"
)

// Server is a telebot server.
type Server struct {
}

// Start starts telebot server.
func (s *Server) Start(ctx context.Context, botToken string) {
	sword, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
		return
	}
	// sword.Debug = true

	log.Printf("Authorized on account %s", sword.Self.UserName)

	engine := botwrapper.NewEngine(sword)

	userX := user.NewController()
	userX.RegisterHandler(engine)

	logging.Get().Info("starts serving bot")
	engine.Start(ctx)

	<-ctx.Done()
	logging.Get().Info("stops serving bot")
	// sword.Stop()
}
