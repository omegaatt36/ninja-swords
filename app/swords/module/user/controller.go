package user

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	botwrapper "ninja-swords/app/swords/bot-wrapper"
)

// Controller defines controller.
type Controller struct {
}

// NewController returns controller.
func NewController() *Controller {
	return &Controller{}
}

// RegisterHandler register handles.
func (x *Controller) RegisterHandler(e *botwrapper.Engine) {
	e.Apply("start", x.start)
}

func (x *Controller) start(c *botwrapper.Context) {
	msg := tgbotapi.NewMessage(c.Message.Chat.ID, fmt.Sprintf("Welcome %s", c.Message.From.UserName))
	c.Send(msg)
}
