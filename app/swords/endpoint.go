package swords

import (
	"fmt"

	botwrapper "github.com/omegaatt36/ninja-swords/app/swords/bot-wrapper"
	"gopkg.in/telebot.v3"
)

func registerHandler(e *botwrapper.Engine) {
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
}
