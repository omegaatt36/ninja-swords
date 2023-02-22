package main

import (
	"context"

	"github.com/omegaatt36/ninja-swords/app"
	"github.com/omegaatt36/ninja-swords/app/swords"
	"github.com/omegaatt36/ninja-swords/health"
	"github.com/urfave/cli/v2"
)

// Main starts process in cli.
func Main(ctx context.Context, c *cli.Context) {
	go health.StartServer()

	server := swords.Server{}
	server.Start(ctx, c.String("bot-token"))
}

func main() {
	app := app.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "bot-token",
				EnvVars:  []string{"BOT_TOKEN"},
				Required: true,
			},
		},
		Main: Main,
	}

	app.Run()
}
