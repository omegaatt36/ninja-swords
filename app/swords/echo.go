package swords

import (
	"fmt"

	"github.com/omegaatt36/ninja-swords/logging"
	"golang.org/x/sync/syncmap"
	"gopkg.in/telebot.v3"
)

type echoService struct {
	mSession syncmap.Map
}

func newEchoService() *echoService {
	return &echoService{}
}

func (service *echoService) openEchoSession(m telebot.Context) error {
	userID := m.Message().Sender.ID
	if _, ok := service.mSession.Load(userID); ok {
		return nil
	}

	if err := m.Send(fmt.Sprintf("type something, %s", m.Message().Sender.Username)); err != nil {
		logging.Get().Error(err)
	}

	service.mSession.Store(userID, userID)
	return nil
}

func (service *echoService) doEcho(m telebot.Context) error {
	userID := m.Message().Sender.ID
	if _, ok := service.mSession.Load(userID); !ok {
		return nil
	}
	if err := m.Send(m.Message().Text); err != nil {
		logging.Get().Error(err)
	}

	service.mSession.Delete(userID)

	return nil
}
