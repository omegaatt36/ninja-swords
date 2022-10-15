package user

import (
	"fmt"

	"gopkg.in/telebot.v3"
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
	e.Apply(botwrapper.Handler{
		Endpoint: "/start",
		HandlerFunc: func(m telebot.Context) error {
			return m.Send(fmt.Sprintf("Welcome %s", m.Message().Sender.Username))
		},
		Describe: "send welcome.",
		OnQuery:  true,
		OnButton: true,
	})

	e.Apply(botwrapper.Handler{
		Endpoint: "/hello",
		HandlerFunc: func(m telebot.Context) error {
			return m.Send("Hello World!")
		},
		Describe: "send welcome.",
		OnQuery:  true,
		OnButton: true,
	})

	e.Apply(botwrapper.Handler{
		Endpoint: telebot.OnAddedToGroup,
		HandlerFunc: func(m telebot.Context) error {
			return m.Send("Ahoy")
		},
	})

	e.Apply(botwrapper.Handler{
		Endpoint: telebot.OnText,
		HandlerFunc: func(m telebot.Context) error {
			if m.Chat().Type != telebot.ChatGroup {
				return nil
			}

			text := m.Text()

			return m.Send(fmt.Sprintf("echo %s", text))
		},
	})
}
